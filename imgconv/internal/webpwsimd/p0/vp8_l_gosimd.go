//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8LBackwardReferencesTraceBackwards(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int32 {
	mBase := m.M
	_ = mBase
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 int64
	_ = v49
	var v50 int32
	_ = v50
	var v57 int64
	_ = v57
	var v58 int32
	_ = v58
	var v59 int64
	_ = v59
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v74 int64
	_ = v74
	var v75 int32
	_ = v75
	var v82 int32
	_ = v82
	var v89 int64
	_ = v89
	var v90 int32
	_ = v90
	var v91 int64
	_ = v91
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v105 int64
	_ = v105
	var v106 int32
	_ = v106
	var v113 int64
	_ = v113
	var v114 int32
	_ = v114
	var v115 int64
	_ = v115
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v141 int32
	_ = v141
	var v147 int32
	_ = v147
	var v155 int32
	_ = v155
	var v163 int32
	_ = v163
	var v168 int32
	_ = v168
	var v171 int32
	_ = v171
	var v179 base.V128
	_ = v179
	var v187 int32
	_ = v187
	var v191 int32
	_ = v191
	var v195 int32
	_ = v195
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v214 int32
	_ = v214
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v229 base.V128
	_ = v229
	var v234 int32
	_ = v234
	var v238 int32
	_ = v238
	var v242 int32
	_ = v242
	var v248 int32
	_ = v248
	var v250 int32
	_ = v250
	var v252 int32
	_ = v252
	var v257 int32
	_ = v257
	var v260 int32
	_ = v260
	var v266 int32
	_ = v266
	var v268 int32
	_ = v268
	var v270 int32
	_ = v270
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v289 int32
	_ = v289
	var v294 int32
	_ = v294
	var v298 int32
	_ = v298
	var v299 int32
	_ = v299
	var v300 base.V128
	_ = v300
	var v302 int32
	_ = v302
	var v319 int32
	_ = v319
	var v322 base.V128
	_ = v322
	var v323 base.V128
	_ = v323
	var v345 base.V128
	_ = v345
	var v346 base.V128
	_ = v346
	var v348 base.V128
	_ = v348
	var v352 int32
	_ = v352
	var v353 base.V128
	_ = v353
	var v355 base.V128
	_ = v355
	var v356 base.V128
	_ = v356
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v363 base.V128
	_ = v363
	var v368 int32
	_ = v368
	var v377 int32
	_ = v377
	var v393 int32
	_ = v393
	var v394 int32
	_ = v394
	var v400 int32
	_ = v400
	var v419 int32
	_ = v419
	var v420 int32
	_ = v420
	var v423 int32
	_ = v423
	var v427 int32
	_ = v427
	var v446 int32
	_ = v446
	var v452 int32
	_ = v452
	var v516 int32
	_ = v516
	var v518 int32
	_ = v518
	var v529 int32
	_ = v529
	var v551 int32
	_ = v551
	var v554 int32
	_ = v554
	var v555 int32
	_ = v555
	var v563 int32
	_ = v563
	var v564 int32
	_ = v564
	var v599 int32
	_ = v599
	var v600 int32
	_ = v600
	var v605 int64
	_ = v605
	var v608 int32
	_ = v608
	var v609 int32
	_ = v609
	var v621 int32
	_ = v621
	var v641 int32
	_ = v641
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
	var v665 int32
	_ = v665
	var v666 int32
	_ = v666
	var v669 int32
	_ = v669
	var v693 int32
	_ = v693
	var v696 int32
	_ = v696
	var v700 int32
	_ = v700
	var v701 int32
	_ = v701
	var v702 int32
	_ = v702
	var v703 int32
	_ = v703
	var v704 int32
	_ = v704
	var v707 int32
	_ = v707
	var v712 int32
	_ = v712
	var v759 base.V128
	_ = v759
	var v778 int32
	_ = v778
	var v782 base.V128
	_ = v782
	var v783 base.V128
	_ = v783
	var v803 int32
	_ = v803
	var v806 int32
	_ = v806
	var v807 base.V128
	_ = v807
	var v811 base.V128
	_ = v811
	var v815 base.V128
	_ = v815
	var v819 base.V128
	_ = v819
	var v823 base.V128
	_ = v823
	var v831 base.V128
	_ = v831
	var v833 int32
	_ = v833
	var v837 int32
	_ = v837
	var v840 base.V128
	_ = v840
	var v972 base.V128
	_ = v972
	var v977 int32
	_ = v977
	var v980 int32
	_ = v980
	var v984 int32
	_ = v984
	var v985 int32
	_ = v985
	var v986 int32
	_ = v986
	var v987 int32
	_ = v987
	var v988 int32
	_ = v988
	var v991 int32
	_ = v991
	var v1034 int32
	_ = v1034
	var v1037 int32
	_ = v1037
	var v1041 int32
	_ = v1041
	var v1042 int32
	_ = v1042
	var v1043 int32
	_ = v1043
	var v1044 int32
	_ = v1044
	var v1045 int32
	_ = v1045
	var v1050 int32
	_ = v1050
	var v1099 base.V128
	_ = v1099
	var v1118 int32
	_ = v1118
	var v1122 base.V128
	_ = v1122
	var v1123 base.V128
	_ = v1123
	var v1143 int32
	_ = v1143
	var v1146 int32
	_ = v1146
	var v1147 base.V128
	_ = v1147
	var v1151 base.V128
	_ = v1151
	var v1155 base.V128
	_ = v1155
	var v1159 base.V128
	_ = v1159
	var v1163 base.V128
	_ = v1163
	var v1171 base.V128
	_ = v1171
	var v1173 int32
	_ = v1173
	var v1177 int32
	_ = v1177
	var v1180 base.V128
	_ = v1180
	var v1312 base.V128
	_ = v1312
	var v1317 int32
	_ = v1317
	var v1320 int32
	_ = v1320
	var v1324 int32
	_ = v1324
	var v1325 int32
	_ = v1325
	var v1326 int32
	_ = v1326
	var v1327 int32
	_ = v1327
	var v1328 int32
	_ = v1328
	var v1331 int32
	_ = v1331
	var v1374 int32
	_ = v1374
	var v1377 int32
	_ = v1377
	var v1381 int32
	_ = v1381
	var v1382 int32
	_ = v1382
	var v1383 int32
	_ = v1383
	var v1384 int32
	_ = v1384
	var v1385 int32
	_ = v1385
	var v1390 int32
	_ = v1390
	var v1439 base.V128
	_ = v1439
	var v1458 int32
	_ = v1458
	var v1462 base.V128
	_ = v1462
	var v1463 base.V128
	_ = v1463
	var v1483 int32
	_ = v1483
	var v1486 int32
	_ = v1486
	var v1487 base.V128
	_ = v1487
	var v1491 base.V128
	_ = v1491
	var v1495 base.V128
	_ = v1495
	var v1499 base.V128
	_ = v1499
	var v1503 base.V128
	_ = v1503
	var v1511 base.V128
	_ = v1511
	var v1513 int32
	_ = v1513
	var v1518 base.V128
	_ = v1518
	var v1650 base.V128
	_ = v1650
	var v1655 int32
	_ = v1655
	var v1658 int32
	_ = v1658
	var v1662 int32
	_ = v1662
	var v1663 int32
	_ = v1663
	var v1664 int32
	_ = v1664
	var v1665 int32
	_ = v1665
	var v1666 int32
	_ = v1666
	var v1669 int32
	_ = v1669
	var v1712 int32
	_ = v1712
	var v1715 int32
	_ = v1715
	var v1719 int32
	_ = v1719
	var v1720 int32
	_ = v1720
	var v1721 int32
	_ = v1721
	var v1722 int32
	_ = v1722
	var v1723 int32
	_ = v1723
	var v1728 int32
	_ = v1728
	var v1775 int32
	_ = v1775
	var v1776 int32
	_ = v1776
	var v1778 int32
	_ = v1778
	var v1779 int32
	_ = v1779
	var v1781 int32
	_ = v1781
	var v1782 int32
	_ = v1782
	var v1784 int32
	_ = v1784
	var v1785 int32
	_ = v1785
	var v1787 base.V128
	_ = v1787
	var v1788 base.V128
	_ = v1788
	var v1795 int32
	_ = v1795
	var v1813 base.V128
	_ = v1813
	var v1814 base.V128
	_ = v1814
	var v1844 base.V128
	_ = v1844
	var v1845 base.V128
	_ = v1845
	var v1875 base.V128
	_ = v1875
	var v1876 base.V128
	_ = v1876
	var v1906 base.V128
	_ = v1906
	var v1907 base.V128
	_ = v1907
	var v1931 base.V128
	_ = v1931
	var v1932 base.V128
	_ = v1932
	var v1958 base.V128
	_ = v1958
	var v1959 base.V128
	_ = v1959
	var v1985 base.V128
	_ = v1985
	var v1986 base.V128
	_ = v1986
	var v2026 int32
	_ = v2026
	var v2030 int32
	_ = v2030
	var v2033 int32
	_ = v2033
	var v2038 int32
	_ = v2038
	var v2041 int32
	_ = v2041
	var v2045 int32
	_ = v2045
	var v2048 int32
	_ = v2048
	var v2054 int32
	_ = v2054
	var v2186 int32
	_ = v2186
	var v2189 int32
	_ = v2189
	var v2192 int32
	_ = v2192
	var v2195 int32
	_ = v2195
	var v2276 int32
	_ = v2276
	var v2279 int32
	_ = v2279
	var v2283 int32
	_ = v2283
	var v2284 int32
	_ = v2284
	var v2285 int32
	_ = v2285
	var v2286 int32
	_ = v2286
	var v2287 int32
	_ = v2287
	var v2292 int32
	_ = v2292
	var v2335 int32
	_ = v2335
	var v2338 int32
	_ = v2338
	var v2342 int32
	_ = v2342
	var v2343 int32
	_ = v2343
	var v2344 int32
	_ = v2344
	var v2345 int32
	_ = v2345
	var v2346 int32
	_ = v2346
	var v2351 int32
	_ = v2351
	var v2398 int32
	_ = v2398
	var v2406 int32
	_ = v2406
	var v2433 int32
	_ = v2433
	var v2440 int32
	_ = v2440
	var v2443 int32
	_ = v2443
	var v2449 int32
	_ = v2449
	var v2453 int32
	_ = v2453
	var v2455 int32
	_ = v2455
	var v2457 int32
	_ = v2457
	var v2470 int32
	_ = v2470
	var v2474 int32
	_ = v2474
	var v2503 int32
	_ = v2503
	var v2504 int32
	_ = v2504
	var v2505 int32
	_ = v2505
	var v2506 int32
	_ = v2506
	var v2509 int32
	_ = v2509
	var v2511 int32
	_ = v2511
	var v2513 int32
	_ = v2513
	var v2518 int32
	_ = v2518
	var v2519 int32
	_ = v2519
	var v2523 int32
	_ = v2523
	var v2526 int64
	_ = v2526
	var v2534 int32
	_ = v2534
	var v2536 int32
	_ = v2536
	var v2544 int32
	_ = v2544
	var v2547 int64
	_ = v2547
	var v2558 int32
	_ = v2558
	var v2571 int32
	_ = v2571
	var v2575 int32
	_ = v2575
	var v2594 int64
	_ = v2594
	var v2602 int64
	_ = v2602
	var v2605 int32
	_ = v2605
	var v2607 int32
	_ = v2607
	var v2608 int64
	_ = v2608
	var v2611 int32
	_ = v2611
	var v2613 int32
	_ = v2613
	var v2617 int32
	_ = v2617
	var v2622 int32
	_ = v2622
	var v2635 int32
	_ = v2635
	var v2658 int64
	_ = v2658
	var v2669 int64
	_ = v2669
	var v2672 int32
	_ = v2672
	var v2674 int32
	_ = v2674
	var v2691 int32
	_ = v2691
	var v2714 int64
	_ = v2714
	var v2720 int32
	_ = v2720
	var v2727 int64
	_ = v2727
	var v2728 int32
	_ = v2728
	var v2729 int64
	_ = v2729
	var v2739 int32
	_ = v2739
	var v2741 int32
	_ = v2741
	var v2748 int64
	_ = v2748
	var v2752 int32
	_ = v2752
	var v2754 int32
	_ = v2754
	var v2765 int32
	_ = v2765
	var v2778 int32
	_ = v2778
	var v2782 int32
	_ = v2782
	var v2801 int64
	_ = v2801
	var v2809 int64
	_ = v2809
	var v2817 int32
	_ = v2817
	var v2819 int32
	_ = v2819
	var v2821 int64
	_ = v2821
	var v2827 int32
	_ = v2827
	var v2837 int32
	_ = v2837
	var v2850 int32
	_ = v2850
	var v2873 int64
	_ = v2873
	var v2886 int64
	_ = v2886
	var v2892 int32
	_ = v2892
	var v2939 int32
	_ = v2939
	var v2946 int64
	_ = v2946
	var v2947 int32
	_ = v2947
	var v2948 int64
	_ = v2948
	var v2958 int32
	_ = v2958
	var v2960 int32
	_ = v2960
	var v2971 int32
	_ = v2971
	var v2973 int32
	_ = v2973
	var v2986 int32
	_ = v2986
	var v3015 base.V128
	_ = v3015
	var v3021 int32
	_ = v3021
	var v3038 int32
	_ = v3038
	var v3071 int32
	_ = v3071
	var v3084 int32
	_ = v3084
	var v3118 int32
	_ = v3118
	var v3162 int32
	_ = v3162
	var v3164 int32
	_ = v3164
	var v3167 int32
	_ = v3167
	var v3170 int32
	_ = v3170
	var v3171 int32
	_ = v3171
	var v3174 int32
	_ = v3174
	var v3177 int32
	_ = v3177
	var v3179 int32
	_ = v3179
	var v3185 int64
	_ = v3185
	var v3193 int32
	_ = v3193
	var v3196 int64
	_ = v3196
	var v3202 int64
	_ = v3202
	var v3204 int32
	_ = v3204
	var v3210 int64
	_ = v3210
	var v3217 int64
	_ = v3217
	var v3223 int64
	_ = v3223
	var v3224 int64
	_ = v3224
	var v3228 int64
	_ = v3228
	var v3230 int32
	_ = v3230
	var v3236 int32
	_ = v3236
	var v3240 int32
	_ = v3240
	var v3252 int64
	_ = v3252
	var v3261 int32
	_ = v3261
	var v3268 int32
	_ = v3268
	var v3276 int32
	_ = v3276
	var v3278 int32
	_ = v3278
	var v3279 int32
	_ = v3279
	var v3282 int32
	_ = v3282
	var v3283 int32
	_ = v3283
	var v3286 int32
	_ = v3286
	var v3288 int32
	_ = v3288
	var v3290 int32
	_ = v3290
	var v3292 int32
	_ = v3292
	var v3293 int32
	_ = v3293
	var v3295 int32
	_ = v3295
	var v3299 int64
	_ = v3299
	var v3302 int32
	_ = v3302
	var v3305 int32
	_ = v3305
	var v3306 int32
	_ = v3306
	var v3308 int32
	_ = v3308
	var v3309 int32
	_ = v3309
	var v3312 int32
	_ = v3312
	var v3314 int32
	_ = v3314
	var v3318 int64
	_ = v3318
	var v3327 int32
	_ = v3327
	var v3330 int64
	_ = v3330
	var v3336 int64
	_ = v3336
	var v3338 int32
	_ = v3338
	var v3344 int64
	_ = v3344
	var v3351 int64
	_ = v3351
	var v3358 int64
	_ = v3358
	var v3360 int32
	_ = v3360
	var v3362 int32
	_ = v3362
	var v3363 int32
	_ = v3363
	var v3364 int64
	_ = v3364
	var v3368 int64
	_ = v3368
	var v3369 int64
	_ = v3369
	var v3372 int32
	_ = v3372
	var v3378 int32
	_ = v3378
	var v3384 int32
	_ = v3384
	var v3386 int32
	_ = v3386
	var v3391 int32
	_ = v3391
	var v3400 int32
	_ = v3400
	var v3412 int32
	_ = v3412
	var v3418 int32
	_ = v3418
	var v3423 int32
	_ = v3423
	var v3426 int32
	_ = v3426
	var v3429 int32
	_ = v3429
	var v3430 int32
	_ = v3430
	var v3431 int32
	_ = v3431
	var v3432 int32
	_ = v3432
	var v3433 int32
	_ = v3433
	var v3436 int32
	_ = v3436
	var v3438 int32
	_ = v3438
	var v3440 int32
	_ = v3440
	var v3445 int32
	_ = v3445
	var v3446 int32
	_ = v3446
	var v3453 int64
	_ = v3453
	var v3454 int64
	_ = v3454
	var v3461 int32
	_ = v3461
	var v3486 int32
	_ = v3486
	var v3489 int32
	_ = v3489
	var v3516 int32
	_ = v3516
	var v3517 int32
	_ = v3517
	var v3522 int32
	_ = v3522
	var v3525 int32
	_ = v3525
	var v3531 int32
	_ = v3531
	var v3570 int32
	_ = v3570
	var v3574 int32
	_ = v3574
	var v3575 int32
	_ = v3575
	var v3580 int32
	_ = v3580
	var v3598 int32
	_ = v3598
	var v3624 int32
	_ = v3624
	var v3626 int32
	_ = v3626
	var v3627 int32
	_ = v3627
	var v3629 int64
	_ = v3629
	var v3630 int64
	_ = v3630
	var v3633 int32
	_ = v3633
	var v3634 int32
	_ = v3634
	var v3637 int32
	_ = v3637
	var v3640 int32
	_ = v3640
	var v3647 int32
	_ = v3647
	var v3663 int32
	_ = v3663
	var v3691 int32
	_ = v3691
	var v3693 int32
	_ = v3693
	var v3694 int32
	_ = v3694
	var v3696 int64
	_ = v3696
	var v3697 int64
	_ = v3697
	var v3700 int32
	_ = v3700
	var v3701 int32
	_ = v3701
	var v3704 int32
	_ = v3704
	var v3707 int32
	_ = v3707
	var v3756 int64
	_ = v3756
	var v3770 int64
	_ = v3770
	var v3800 int32
	_ = v3800
	var v3801 int32
	_ = v3801
	var v3804 int32
	_ = v3804
	var v3821 int32
	_ = v3821
	var v3850 int32
	_ = v3850
	var v3852 int32
	_ = v3852
	var v3853 int32
	_ = v3853
	var v3855 int32
	_ = v3855
	var v3858 int32
	_ = v3858
	var v3865 int32
	_ = v3865
	var v3867 int32
	_ = v3867
	var v3869 int32
	_ = v3869
	var v3871 int32
	_ = v3871
	var v3875 int32
	_ = v3875
	var v3876 int32
	_ = v3876
	var v3877 int64
	_ = v3877
	var v3878 int64
	_ = v3878
	var v3881 int32
	_ = v3881
	var v3882 int32
	_ = v3882
	var v3885 int32
	_ = v3885
	var v3888 int32
	_ = v3888
	var v3935 int32
	_ = v3935
	var v3940 int32
	_ = v3940
	var v3985 int32
	_ = v3985
	var v3988 int32
	_ = v3988
	var v4040 int32
	_ = v4040
	var v4042 int32
	_ = v4042
	var v4044 int32
	_ = v4044
	var v4052 int32
	_ = v4052
	var v4056 int32
	_ = v4056
	var v4065 int32
	_ = v4065
	var v4067 int32
	_ = v4067
	var v4075 int32
	_ = v4075
	var v4079 int32
	_ = v4079
	var v4090 int32
	_ = v4090
	var v4128 int32
	_ = v4128
	var v4173 int32
	_ = v4173
	var v4176 int32
	_ = v4176
	var v4181 int32
	_ = v4181
	var v4229 int32
	_ = v4229
	var v4231 int32
	_ = v4231
	var v4233 int32
	_ = v4233
	var v4241 int32
	_ = v4241
	var v4245 int32
	_ = v4245
	var v4254 int32
	_ = v4254
	var v4256 int32
	_ = v4256
	var v4264 int32
	_ = v4264
	var v4268 int32
	_ = v4268
	var v4279 int32
	_ = v4279
	var v4321 int32
	_ = v4321
	var v4323 int32
	_ = v4323
	var v4326 int32
	_ = v4326
	var v4340 int32
	_ = v4340
	var v4369 int32
	_ = v4369
	var v4370 int32
	_ = v4370
	var v4374 int32
	_ = v4374
	var v4391 int32
	_ = v4391
	var v4420 int32
	_ = v4420
	var v4421 int32
	_ = v4421
	var v4425 int32
	_ = v4425
	var v4431 int32
	_ = v4431
	var v4442 int32
	_ = v4442
	var v4445 int32
	_ = v4445
	var v4447 int32
	_ = v4447
	var v4449 int32
	_ = v4449
	var v4456 int32
	_ = v4456
	var v4459 int32
	_ = v4459
	var v4463 int32
	_ = v4463
	var v4466 int32
	_ = v4466
	var v4468 int32
	_ = v4468
	var v4470 int32
	_ = v4470
	var v4477 int32
	_ = v4477
	var v4480 int32
	_ = v4480
	var v4484 int32
	_ = v4484
	var v4498 int32
	_ = v4498
	var v4502 int32
	_ = v4502
	var v4526 int32
	_ = v4526
	var v4530 int32
	_ = v4530
	var v4532 int32
	_ = v4532
	var v4536 int32
	_ = v4536
	var v4540 int64
	_ = v4540
	var v4546 int32
	_ = v4546
	var v4549 int32
	_ = v4549
	var v4550 int32
	_ = v4550
	var v4553 int32
	_ = v4553
	var v4555 int32
	_ = v4555
	var v4560 int32
	_ = v4560
	var v4561 int32
	_ = v4561
	var v4568 int32
	_ = v4568
	var v4570 int32
	_ = v4570
	var v4571 int32
	_ = v4571
	var v4573 int32
	_ = v4573
	var v4580 int32
	_ = v4580
	var v4581 int32
	_ = v4581
	var v4585 int32
	_ = v4585
	var v4589 int64
	_ = v4589
	var v4596 int64
	_ = v4596
	var v4600 int64
	_ = v4600
	var v4607 int32
	_ = v4607
	var v4610 int32
	_ = v4610
	var v4611 int32
	_ = v4611
	var v4614 int32
	_ = v4614
	var v4616 int32
	_ = v4616
	var v4621 int32
	_ = v4621
	var v4622 int32
	_ = v4622
	var v4629 int32
	_ = v4629
	var v4631 int32
	_ = v4631
	var v4632 int32
	_ = v4632
	var v4634 int32
	_ = v4634
	var v4641 int32
	_ = v4641
	var v4642 int32
	_ = v4642
	var v4646 int32
	_ = v4646
	var v4650 int64
	_ = v4650
	var v4659 int32
	_ = v4659
	var v4677 int32
	_ = v4677
	var v4685 int32
	_ = v4685
	var v4704 int32
	_ = v4704
	var v4707 int32
	_ = v4707
	var v4711 int32
	_ = v4711
	var v4713 int32
	_ = v4713
	var v4715 int32
	_ = v4715
	var v4717 int32
	_ = v4717
	var v4721 int64
	_ = v4721
	var v4727 int32
	_ = v4727
	var v4730 int32
	_ = v4730
	var v4731 int32
	_ = v4731
	var v4734 int32
	_ = v4734
	var v4736 int32
	_ = v4736
	var v4741 int32
	_ = v4741
	var v4742 int32
	_ = v4742
	var v4749 int32
	_ = v4749
	var v4751 int32
	_ = v4751
	var v4752 int32
	_ = v4752
	var v4754 int32
	_ = v4754
	var v4761 int32
	_ = v4761
	var v4762 int32
	_ = v4762
	var v4766 int32
	_ = v4766
	var v4770 int64
	_ = v4770
	var v4781 int32
	_ = v4781
	var v4783 int32
	_ = v4783
	var v4795 int32
	_ = v4795
	var v4825 int32
	_ = v4825
	var v4826 int32
	_ = v4826
	var v4828 int32
	_ = v4828
	var v4830 int32
	_ = v4830
	var v4836 int32
	_ = v4836
	var v4839 int32
	_ = v4839
	var v4846 int32
	_ = v4846
	var v4848 int32
	_ = v4848
	var v4851 int32
	_ = v4851
	var v4854 int32
	_ = v4854
	var v4860 int32
	_ = v4860
	var v4861 int32
	_ = v4861
	var v4863 int32
	_ = v4863
	var v4864 int32
	_ = v4864
	var v4867 int32
	_ = v4867
	var v4868 int32
	_ = v4868
	var v4871 int32
	_ = v4871
	var v4872 int32
	_ = v4872
	var v4879 int32
	_ = v4879
	var v4881 int32
	_ = v4881
	var v4882 int64
	_ = v4882
	var v4886 int64
	_ = v4886
	var v4893 int32
	_ = v4893
	var v4896 int32
	_ = v4896
	var v4897 int32
	_ = v4897
	var v4900 int32
	_ = v4900
	var v4902 int32
	_ = v4902
	var v4907 int32
	_ = v4907
	var v4908 int32
	_ = v4908
	var v4915 int32
	_ = v4915
	var v4917 int32
	_ = v4917
	var v4918 int32
	_ = v4918
	var v4920 int32
	_ = v4920
	var v4927 int32
	_ = v4927
	var v4928 int32
	_ = v4928
	var v4932 int32
	_ = v4932
	var v4936 int64
	_ = v4936
	var v4985 int32
	_ = v4985
	var v5030 int32
	_ = v5030
	var v5032 int32
	_ = v5032
	var v5046 int32
	_ = v5046
	var v5079 int32
	_ = v5079
	var v5082 int32
	_ = v5082
	var v5097 int32
	_ = v5097
	v44 = m.G0
	v46 = v44 - int32(80)
	m.G0 = v46
	v48 = l1 * l0
	v49 = base.I64_extend_i32_s(v48)
	v50 = int32(2)
	if v49 == int64(0) {
		goto L5
	} else {
		goto L6
	}
L1:
	;
	F_free(m, v71)
	mBase = m.M
	goto L499
L2:
	;
	v74 = int64(1)
	v75 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_0)
	if l3 < int32(1) {
		goto L10
	} else {
		goto L11
	}
L3:
	;
	if v71 != 0 {
		goto L2
	} else {
		goto L9
	}
L4:
	;
	goto L3
L5:
	;
	v69 = F_malloc(m, base.I32_wrap_i64(v49)*v50)
	mBase = m.M
	v71 = v69
	goto L4
L6:
	;
	v57 = base.I64_div_u_s(int64(2147418112), v49)
	v58 = int32(0)
	v59 = base.I64_extend_i32_u(v50)
	if base.Ui64(int64(4294967295)) < base.Ui64(v59*v49) {
		v71 = v58
		goto L4
	} else {
		goto L7
	}
L7:
	;
	if base.Ui64(v57) < base.Ui64(v59) {
		v71 = v58
		goto L4
	} else {
		goto L8
	}
L8:
	;
	goto L5
L9:
	;
	v5097 = int32(0)
	goto L1
L10:
	;
	v82 = v75
	goto L12
L11:
	;
	v82 = int32(4)<<(uint(l3)%32) + v75
	goto L12
L12:
	;
	goto L16
L13:
	;
	v105 = int64(1)
	v106 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_1)
	goto L22
L14:
	;
	goto L13
L15:
	;
	v100 = F_calloc(m, base.I32_wrap_i64(v74), v82)
	mBase = m.M
	v102 = v100
	goto L14
L16:
	;
	v89 = base.I64_div_u_s(int64(2147418112), v74)
	v90 = int32(0)
	v91 = base.I64_extend_i32_u(v82)
	if base.Ui64(int64(4294967295)) < base.Ui64(v91*v74) {
		v102 = v90
		goto L14
	} else {
		goto L17
	}
L17:
	;
	if base.Ui64(v89) < base.Ui64(v91) {
		v102 = v90
		goto L14
	} else {
		goto L18
	}
L18:
	;
	goto L15
L19:
	;
	v128 = int32(0)
	if v102 == v128 {
		v4181 = v128
		goto L25
	} else {
		goto L26
	}
L20:
	;
	goto L19
L21:
	;
	v124 = F_calloc(m, base.I32_wrap_i64(v105), v106)
	mBase = m.M
	v126 = v124
	goto L20
L22:
	;
	v113 = base.I64_div_u_s(int64(2147418112), v105)
	v114 = int32(0)
	v115 = base.I64_extend_i32_u(v106)
	if base.Ui64(int64(4294967295)) < base.Ui64(v115*v105) {
		v126 = v114
		goto L20
	} else {
		goto L23
	}
L23:
	;
	if base.Ui64(v113) < base.Ui64(v115) {
		v126 = v114
		goto L20
	} else {
		goto L24
	}
L24:
	;
	goto L21
L25:
	;
	if v126 == int32(0) {
		goto L387
	} else {
		goto L388
	}
L26:
	;
	v131 = int32(0)
	if v126 == v131 {
		v4181 = v131
		goto L25
	} else {
		goto L27
	}
L27:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v102)+3232)) = v102 + int32(3236)
	if l3 < int32(1) {
		goto L28
	} else {
		goto L29
	}
L28:
	;
	v163 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_2)
	if int32(0) < l3 {
		goto L39
	} else {
		goto L40
	}
L29:
	;
	v141 = v46 + int32(68)
	v147 = F_WebPSafeCalloc(m, base.I64_extend_i32_s(int32(1)<<(uint(l3)%32)), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v141))) = v147
	if v147 != 0 {
		goto L31
	} else {
		goto L32
	}
L30:
	;
	if v155 == int32(0) {
		v4181 = int32(0)
		goto L25
	} else {
		goto L33
	}
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v141)+8)) = l3
	*(*int32)(unsafe.Add(mBase, uint32(v141)+4)) = int32(32) - l3
	v155 = int32(1)
	goto L30
L32:
	;
	v155 = int32(0)
	goto L30
L33:
	;
	goto L28
L34:
	;
	if l3 < int32(1) {
		v4181 = v4128
		goto L25
	} else {
		goto L382
	}
L35:
	;
	if v126 == int32(0) {
		goto L362
	} else {
		goto L363
	}
L36:
	;
	v3988 = int32(0)
	F_WebPSafeFree(m, v3988)
	mBase = m.M
	goto L360
L37:
	;
	if v171 == int32(0) {
		goto L36
	} else {
		goto L43
	}
L38:
	;
	goto L37
L39:
	;
	v168 = int32(4)<<(uint(l3)%32) + v163
	goto L41
L40:
	;
	v168 = v163
	goto L41
L41:
	;
	v171 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v168), int32(1))
	mBase = m.M
	if v171 == int32(0) {
		goto L38
	} else {
		goto L42
	}
L42:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v171)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v171)+3236)) = l3
	*(*int32)(unsafe.Add(mBase, uint32(v171)+3304)) = int32(16843009)
	v179 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k0)
	base.Simd_g_v128_store(m, v171, int32(3256), v179)
	*(*int32)(unsafe.Add(mBase, uint32(v171))) = v171 + int32(3312)
	v187 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v171+int32(3248)))) = uint16(v187)
	v191 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v171+int32(3308)))) = uint8(v191)
	v195 = int32(0)
	base.Simd_g_v128_store(m, v171+int32(3272), v195, v179)
	base.Simd_g_v128_store(m, v171+int32(3288), v195, v179)
	goto L38
L43:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v171)+3236)) = l3
	goto L46
L44:
	;
	v248 = m.G116
	v250 = m.G0
	v252 = v250 - int32(16)
	m.G0 = v252
	F_VP8LRefsCursorInit(m, v252+int32(4), l5)
	mBase = m.M
	v257 = *(*int32)(unsafe.Add(mBase, uint32(v252)+4))
	if v257 == int32(0) {
		goto L51
	} else {
		goto L52
	}
L45:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v171)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v171)+3304)) = int32(16843009)
	v229 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k0)
	base.Simd_g_v128_store(m, v171, int32(3256), v229)
	v234 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v171+int32(3248)))) = uint16(v234)
	v238 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v171+int32(3308)))) = uint8(v238)
	v242 = int32(0)
	base.Simd_g_v128_store(m, v171+int32(3272), v242, v229)
	base.Simd_g_v128_store(m, v171+int32(3288), v242, v229)
	goto L44
L46:
	;
	v210 = *(*int32)(unsafe.Add(mBase, uint32(v171)))
	v211 = int32(0)
	v214 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_2)
	if v211 < l3 {
		goto L47
	} else {
		goto L48
	}
L47:
	;
	v219 = int32(4)<<(uint(l3)%32) + v214
	goto L49
L48:
	;
	v219 = v214
	goto L49
L49:
	;
	v220 = F_memset(m, v171, v211, v219)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v220))) = v210
	*(*int32)(unsafe.Add(mBase, uint32(v220)+3236)) = l3
	goto L45
L50:
	;
	v285 = *(*int32)(unsafe.Add(mBase, uint32(v102)+3232))
	v287 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3236))
	v289 = int32(280)
	if int32(0) < v287 {
		goto L61
	} else {
		goto L62
	}
L51:
	;
	m.G0 = v252 + int32(16)
	goto L50
L52:
	;
	v260 = v257
	goto L53
L53:
	;
	F_HistogramAddSinglePixOrCopy(m, v171, v260, v248, l0)
	mBase = m.M
	v266 = *(*int32)(unsafe.Add(mBase, uint32(v252)+4))
	v268 = v266 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v252)+4)) = v268
	v270 = *(*int32)(unsafe.Add(mBase, uint32(v252)+12))
	if v268 != v270 {
		v276 = v268
		goto L55
	} else {
		goto L56
	}
L54:
	;
	goto L51
L55:
	;
	if v276 != 0 {
		v260 = v276
		goto L53
	} else {
		goto L57
	}
L56:
	;
	F_VP8LRefsCursorNextBlock(m, v252+int32(4))
	mBase = m.M
	v275 = *(*int32)(unsafe.Add(mBase, uint32(v252)+4))
	v276 = v275
	goto L55
L57:
	;
	goto L54
L58:
	;
	v759 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k0)
	v778 = int32(0)
	v782 = v759
	v783 = v759
	goto L97
L59:
	;
	if base.Ui32(int32(255)) < base.Ui32(v452) {
		goto L89
	} else {
		goto L90
	}
L60:
	;
	v516 = int32(0)
	v518 = v294 << (uint(int32(2)) % 32)
	if base.Ui32(v518) < base.Ui32(int32(33)) {
		goto L75
	} else {
		goto L76
	}
L61:
	;
	v294 = int32(1)<<(uint(v287)%32) + v289
	goto L63
L62:
	;
	v294 = v289
	goto L63
L63:
	;
	if v294 < int32(1) {
		goto L60
	} else {
		goto L64
	}
L64:
	;
	v298 = v294 & int32(2147483644)
	v299 = *(*int32)(unsafe.Add(mBase, uint32(v171)))
	v300 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k0)
	v302 = v299
	v319 = v298
	v322 = v300
	v323 = v300
	goto L65
L65:
	;
	v345 = base.Simd_g_v128_load(m, v302, int32(0))
	v346 = base.Simd_g_i32x4_add(v345, v322)
	v348 = base.Simd_g_i32x4_sub(v323, base.Simd_g_i32x4_ne(v345, v300))
	v352 = v319 + int32(-4)
	if v352 != 0 {
		v302 = v302 + int32(16)
		v319 = v352
		v322 = v346
		v323 = v348
		goto L65
	} else {
		goto L67
	}
L66:
	;
	v353 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)
	v355 = base.Simd_g_i32x4_add(v348, base.Simd_g_i8x16_swizzle_c(v348, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	v356 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2)
	v359 = int32(0)
	v360 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v355, base.Simd_g_i8x16_swizzle_c(v355, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))
	v363 = base.Simd_g_i32x4_add(v346, base.Simd_g_i8x16_swizzle_c(v346, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	v368 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v363, base.Simd_g_i8x16_swizzle_c(v363, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))
	if v294 == v298 {
		v446 = v360
		v452 = v368
		goto L68
	} else {
		goto L69
	}
L67:
	;
	goto L66
L68:
	;
	if base.Ui32(int32(1)) < base.Ui32(v446) {
		goto L59
	} else {
		goto L73
	}
L69:
	;
	v377 = v299 + v294<<(uint(int32(2))%32)&int32(-16)
	v393 = v294 - v298
	v394 = v360
	v400 = v368
	goto L70
L70:
	;
	v419 = *(*int32)(unsafe.Add(mBase, uint32(v377)))
	v420 = v419 + v400
	v423 = v394 + base.B2i32(v419 != int32(0))
	v427 = v393 + int32(-1)
	if v427 != 0 {
		v377 = v377 + int32(4)
		v393 = v427
		v394 = v423
		v400 = v420
		goto L70
	} else {
		goto L72
	}
L71:
	;
	v446 = v423
	v452 = v420
	goto L68
L72:
	;
	goto L71
L73:
	;
	goto L60
L74:
	;
	goto L58
L75:
	;
	if v518 == int32(0) {
		goto L77
	} else {
		goto L78
	}
L76:
	;
	base.MemoryFill(m, v285, v516, v518)
	goto L74
L77:
	;
	goto L74
L78:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v285))) = uint8(v516)
	v529 = v285 + v518
	*(*uint8)(unsafe.Add(mBase, uint32(v529+int32(-1)))) = uint8(v516)
	if base.Ui32(v518) < base.Ui32(int32(3)) {
		goto L77
	} else {
		goto L79
	}
L79:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v285)+2)) = uint8(v516)
	*(*uint8)(unsafe.Add(mBase, uint32(v285)+1)) = uint8(v516)
	*(*uint8)(unsafe.Add(mBase, uint32(v529+int32(-3)))) = uint8(v516)
	*(*uint8)(unsafe.Add(mBase, uint32(v529+int32(-2)))) = uint8(v516)
	if base.Ui32(v518) < base.Ui32(int32(7)) {
		goto L77
	} else {
		goto L80
	}
L80:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v285)+3)) = uint8(v516)
	*(*uint8)(unsafe.Add(mBase, uint32(v529+int32(-4)))) = uint8(v516)
	if base.Ui32(v518) < base.Ui32(int32(9)) {
		goto L77
	} else {
		goto L81
	}
L81:
	;
	v551 = int32(0)
	v554 = (v551 - v285) & int32(3)
	v555 = v285 + v554
	*(*int32)(unsafe.Add(mBase, uint32(v555))) = v551
	v563 = (v518 - v554) & int32(60)
	v564 = v555 + v563
	*(*int32)(unsafe.Add(mBase, uint32(v564+int32(-4)))) = v551
	if base.Ui32(v563) < base.Ui32(int32(9)) {
		goto L77
	} else {
		goto L82
	}
L82:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v555)+8)) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v555)+4)) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v564+int32(-8)))) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v564+int32(-12)))) = v551
	if base.Ui32(v563) < base.Ui32(int32(25)) {
		goto L77
	} else {
		goto L83
	}
L83:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v555)+24)) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v555)+20)) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v555)+16)) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v555)+12)) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v564+int32(-16)))) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v564+int32(-20)))) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v564+int32(-24)))) = v551
	*(*int32)(unsafe.Add(mBase, uint32(v564+int32(-28)))) = v551
	v599 = v555&int32(4) | int32(24)
	v600 = v563 - v599
	if base.Ui32(v600) < base.Ui32(int32(32)) {
		goto L77
	} else {
		goto L84
	}
L84:
	;
	v605 = base.I64_extend_i32_u(v551) * int64(4294967297)
	v608 = v600
	v609 = v555 + v599
	goto L85
L85:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v609)+24)) = v605
	*(*int64)(unsafe.Add(mBase, uint32(v609)+16)) = v605
	*(*int64)(unsafe.Add(mBase, uint32(v609)+8)) = v605
	*(*int64)(unsafe.Add(mBase, uint32(v609))) = v605
	v621 = v608 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v621) {
		v608 = v621
		v609 = v609 + int32(32)
		goto L85
	} else {
		goto L87
	}
L86:
	;
	goto L77
L87:
	;
	goto L86
L88:
	;
	v665 = v285
	v666 = v294
	v669 = v299
	goto L91
L89:
	;
	v646 = m.G118
	v647 = *(*int32)(unsafe.Add(mBase, uint32(v646)))
	v648 = m.T0[v647].(func(*base.Module, int32) int32)(m, v452)
	mBase = m.M
	v649 = v648
	goto L88
L90:
	;
	v641 = m.G117
	v645 = *(*int32)(unsafe.Add(mBase, uint32(v641+v452<<(uint(int32(2))%32))))
	v649 = v645
	goto L88
L91:
	;
	v693 = *(*int32)(unsafe.Add(mBase, uint32(v669)))
	if base.Ui32(int32(255)) < base.Ui32(v693) {
		goto L94
	} else {
		goto L95
	}
L92:
	;
	goto L58
L93:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v665))) = v649 - v704
	v707 = int32(4)
	v712 = v666 + int32(-1)
	if v712 != 0 {
		v665 = v665 + v707
		v666 = v712
		v669 = v669 + v707
		goto L91
	} else {
		goto L96
	}
L94:
	;
	v701 = m.G118
	v702 = *(*int32)(unsafe.Add(mBase, uint32(v701)))
	v703 = m.T0[v702].(func(*base.Module, int32) int32)(m, v693)
	mBase = m.M
	v704 = v703
	goto L93
L95:
	;
	v696 = m.G117
	v700 = *(*int32)(unsafe.Add(mBase, uint32(v696+v693<<(uint(int32(2))%32))))
	v704 = v700
	goto L93
L96:
	;
	goto L92
L97:
	;
	v803 = v171 + v778
	v806 = int32(0)
	v807 = base.Simd_g_v128_load_rng(m, v803+int32(52), v806, int32(-48), int32(64))
	v811 = base.Simd_g_v128_load_nc(m, v803+int32(36), v806)
	v815 = base.Simd_g_v128_load_nc(m, v803+int32(20), v806)
	v819 = base.Simd_g_v128_load_nc(m, v803+int32(4), v806)
	v823 = base.Simd_g_i32x4_add(v807, base.Simd_g_i32x4_add(v811, base.Simd_g_i32x4_add(v815, base.Simd_g_i32x4_add(v819, v783))))
	v831 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(v782, base.Simd_g_i32x4_ne(v819, v759)), base.Simd_g_i32x4_ne(v815, v759)), base.Simd_g_i32x4_ne(v811, v759)), base.Simd_g_i32x4_ne(v807, v759))
	v833 = v778 + int32(64)
	if v833 != int32(1024) {
		v778 = v833
		v782 = v831
		v783 = v823
		goto L97
	} else {
		goto L99
	}
L98:
	;
	v837 = v102 + int32(1024)
	v840 = base.Simd_g_i32x4_add(v831, base.Simd_g_i8x16_swizzle_c(v831, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	if base.Ui32(int32(1)) < base.Ui32(base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v840, base.Simd_g_i8x16_swizzle_c(v840, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))) {
		goto L101
	} else {
		goto L102
	}
L99:
	;
	goto L98
L100:
	;
	v1099 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k0)
	v1118 = int32(0)
	v1122 = v1099
	v1123 = v1099
	goto L126
L101:
	;
	v972 = base.Simd_g_i32x4_add(v823, base.Simd_g_i8x16_swizzle_c(v823, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	v977 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v972, base.Simd_g_i8x16_swizzle_c(v972, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))
	if base.Ui32(int32(255)) < base.Ui32(v977) {
		goto L118
	} else {
		goto L119
	}
L102:
	;
	goto L105
L103:
	;
	goto L100
L105:
	;
	base.MemoryFill(m, v837, int32(0), int32(1024))
	goto L103
L117:
	;
	v991 = int32(0)
	goto L120
L118:
	;
	v985 = m.G118
	v986 = *(*int32)(unsafe.Add(mBase, uint32(v985)))
	v987 = m.T0[v986].(func(*base.Module, int32) int32)(m, v977)
	mBase = m.M
	v988 = v987
	goto L117
L119:
	;
	v980 = m.G117
	v984 = *(*int32)(unsafe.Add(mBase, uint32(v980+v977<<(uint(int32(2))%32))))
	v988 = v984
	goto L117
L120:
	;
	v1034 = *(*int32)(unsafe.Add(mBase, uint32(v171+int32(4)+v991)))
	if base.Ui32(int32(255)) < base.Ui32(v1034) {
		goto L123
	} else {
		goto L124
	}
L121:
	;
	goto L100
L122:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v837+v991))) = v988 - v1045
	v1050 = v991 + int32(4)
	if v1050 != int32(1024) {
		v991 = v1050
		goto L120
	} else {
		goto L125
	}
L123:
	;
	v1042 = m.G118
	v1043 = *(*int32)(unsafe.Add(mBase, uint32(v1042)))
	v1044 = m.T0[v1043].(func(*base.Module, int32) int32)(m, v1034)
	mBase = m.M
	v1045 = v1044
	goto L122
L124:
	;
	v1037 = m.G117
	v1041 = *(*int32)(unsafe.Add(mBase, uint32(v1037+v1034<<(uint(int32(2))%32))))
	v1045 = v1041
	goto L122
L125:
	;
	goto L121
L126:
	;
	v1143 = v171 + v1118
	v1146 = int32(0)
	v1147 = base.Simd_g_v128_load_rng(m, v1143+int32(1076), v1146, int32(-48), int32(64))
	v1151 = base.Simd_g_v128_load_nc(m, v1143+int32(1060), v1146)
	v1155 = base.Simd_g_v128_load_nc(m, v1143+int32(1044), v1146)
	v1159 = base.Simd_g_v128_load_nc(m, v1143+int32(1028), v1146)
	v1163 = base.Simd_g_i32x4_add(v1147, base.Simd_g_i32x4_add(v1151, base.Simd_g_i32x4_add(v1155, base.Simd_g_i32x4_add(v1159, v1123))))
	v1171 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(v1122, base.Simd_g_i32x4_ne(v1159, v1099)), base.Simd_g_i32x4_ne(v1155, v1099)), base.Simd_g_i32x4_ne(v1151, v1099)), base.Simd_g_i32x4_ne(v1147, v1099))
	v1173 = v1118 + int32(64)
	if v1173 != int32(1024) {
		v1118 = v1173
		v1122 = v1171
		v1123 = v1163
		goto L126
	} else {
		goto L128
	}
L127:
	;
	v1177 = v102 + int32(2048)
	v1180 = base.Simd_g_i32x4_add(v1171, base.Simd_g_i8x16_swizzle_c(v1171, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	if base.Ui32(int32(1)) < base.Ui32(base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v1180, base.Simd_g_i8x16_swizzle_c(v1180, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))) {
		goto L130
	} else {
		goto L131
	}
L128:
	;
	goto L127
L129:
	;
	v1439 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k0)
	v1458 = int32(0)
	v1462 = v1439
	v1463 = v1439
	goto L155
L130:
	;
	v1312 = base.Simd_g_i32x4_add(v1163, base.Simd_g_i8x16_swizzle_c(v1163, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	v1317 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v1312, base.Simd_g_i8x16_swizzle_c(v1312, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))
	if base.Ui32(int32(255)) < base.Ui32(v1317) {
		goto L147
	} else {
		goto L148
	}
L131:
	;
	goto L134
L132:
	;
	goto L129
L134:
	;
	base.MemoryFill(m, v1177, int32(0), int32(1024))
	goto L132
L146:
	;
	v1331 = int32(0)
	goto L149
L147:
	;
	v1325 = m.G118
	v1326 = *(*int32)(unsafe.Add(mBase, uint32(v1325)))
	v1327 = m.T0[v1326].(func(*base.Module, int32) int32)(m, v1317)
	mBase = m.M
	v1328 = v1327
	goto L146
L148:
	;
	v1320 = m.G117
	v1324 = *(*int32)(unsafe.Add(mBase, uint32(v1320+v1317<<(uint(int32(2))%32))))
	v1328 = v1324
	goto L146
L149:
	;
	v1374 = *(*int32)(unsafe.Add(mBase, uint32(v171+int32(1028)+v1331)))
	if base.Ui32(int32(255)) < base.Ui32(v1374) {
		goto L152
	} else {
		goto L153
	}
L150:
	;
	goto L129
L151:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1177+v1331))) = v1328 - v1385
	v1390 = v1331 + int32(4)
	if v1390 != int32(1024) {
		v1331 = v1390
		goto L149
	} else {
		goto L154
	}
L152:
	;
	v1382 = m.G118
	v1383 = *(*int32)(unsafe.Add(mBase, uint32(v1382)))
	v1384 = m.T0[v1383].(func(*base.Module, int32) int32)(m, v1374)
	mBase = m.M
	v1385 = v1384
	goto L151
L153:
	;
	v1377 = m.G117
	v1381 = *(*int32)(unsafe.Add(mBase, uint32(v1377+v1374<<(uint(int32(2))%32))))
	v1385 = v1381
	goto L151
L154:
	;
	goto L150
L155:
	;
	v1483 = v171 + v1458
	v1486 = int32(0)
	v1487 = base.Simd_g_v128_load_rng(m, v1483+int32(2100), v1486, int32(-48), int32(64))
	v1491 = base.Simd_g_v128_load_nc(m, v1483+int32(2084), v1486)
	v1495 = base.Simd_g_v128_load_nc(m, v1483+int32(2068), v1486)
	v1499 = base.Simd_g_v128_load_nc(m, v1483+int32(2052), v1486)
	v1503 = base.Simd_g_i32x4_add(v1487, base.Simd_g_i32x4_add(v1491, base.Simd_g_i32x4_add(v1495, base.Simd_g_i32x4_add(v1499, v1463))))
	v1511 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_sub(v1462, base.Simd_g_i32x4_ne(v1499, v1439)), base.Simd_g_i32x4_ne(v1495, v1439)), base.Simd_g_i32x4_ne(v1491, v1439)), base.Simd_g_i32x4_ne(v1487, v1439))
	v1513 = v1458 + int32(64)
	if v1513 != int32(1024) {
		v1458 = v1513
		v1462 = v1511
		v1463 = v1503
		goto L155
	} else {
		goto L157
	}
L156:
	;
	v1518 = base.Simd_g_i32x4_add(v1511, base.Simd_g_i8x16_swizzle_c(v1511, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	if base.Ui32(int32(1)) < base.Ui32(base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v1518, base.Simd_g_i8x16_swizzle_c(v1518, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))) {
		goto L159
	} else {
		goto L160
	}
L157:
	;
	goto L156
L158:
	;
	v1775 = v102 + int32(3072)
	v1776 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3076))
	v1778 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3084))
	v1779 = int32(1)
	v1781 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3088))
	v1782 = int32(2)
	v1784 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3092))
	v1785 = int32(3)
	v1787 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k0)
	v1788 = base.Simd_g_i32x4_ne(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v1776), v1778), v1781), v1784), v1787)
	v1795 = int32(0)
	v1813 = base.Simd_g_v128_load_rng(m, v171, int32(3096), int32(3096), int32(112))
	v1814 = base.Simd_g_i32x4_ne(v1813, v1787)
	v1844 = base.Simd_g_v128_load_nc(m, v171, int32(3112))
	v1845 = base.Simd_g_i32x4_ne(v1844, v1787)
	v1875 = base.Simd_g_v128_load_nc(m, v171, int32(3128))
	v1876 = base.Simd_g_i32x4_ne(v1875, v1787)
	v1906 = base.Simd_g_v128_load_nc(m, v171, int32(3192))
	v1907 = base.Simd_g_i32x4_ne(v1906, v1787)
	v1931 = base.Simd_g_v128_load_nc(m, v171, int32(3176))
	v1932 = base.Simd_g_i32x4_ne(v1931, v1787)
	v1958 = base.Simd_g_v128_load_nc(m, v171, int32(3160))
	v1959 = base.Simd_g_i32x4_ne(v1958, v1787)
	v1985 = base.Simd_g_v128_load_nc(m, v171, int32(3144))
	v1986 = base.Simd_g_i32x4_ne(v1985, v1787)
	v2026 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3080))
	v2030 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3208))
	v2033 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3212))
	v2038 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3216))
	v2041 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3220))
	v2045 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3224))
	v2048 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3228))
	v2054 = *(*int32)(unsafe.Add(mBase, uint32(v171)+3232))
	if base.Ui32(v1779) < base.Ui32(base.I32_popcnt((base.Simd_g_i32x4_extract_lane_l1(v1788)&v1779<<(uint(v1779)%32)-base.Simd_g_i32x4_extract_lane_l0(v1788)|base.Simd_g_i32x4_extract_lane_l2(v1788)&v1779<<(uint(v1782)%32)|base.Simd_g_i32x4_extract_lane_l3(v1788)&v1779<<(uint(v1785)%32)|base.Simd_g_i32x4_extract_lane_l0(v1814)&v1779<<(uint(int32(4))%32)|base.Simd_g_i32x4_extract_lane_l1(v1814)&v1779<<(uint(int32(5))%32)|base.Simd_g_i32x4_extract_lane_l2(v1814)&v1779<<(uint(int32(6))%32)|base.Simd_g_i32x4_extract_lane_l3(v1814)&v1779<<(uint(int32(7))%32)|base.Simd_g_i32x4_extract_lane_l0(v1845)&v1779<<(uint(int32(8))%32)|base.Simd_g_i32x4_extract_lane_l1(v1845)&v1779<<(uint(int32(9))%32)|base.Simd_g_i32x4_extract_lane_l2(v1845)&v1779<<(uint(int32(10))%32)|base.Simd_g_i32x4_extract_lane_l3(v1845)&v1779<<(uint(int32(11))%32)|base.Simd_g_i32x4_extract_lane_l0(v1876)&v1779<<(uint(int32(12))%32)|base.Simd_g_i32x4_extract_lane_l1(v1876)&v1779<<(uint(int32(13))%32)|base.Simd_g_i32x4_extract_lane_l2(v1876)&v1779<<(uint(int32(14))%32)|base.Simd_g_i32x4_extract_lane_l3(v1876)<<(uint(int32(15))%32))&int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)|(base.Simd_g_i32x4_extract_lane_l3(v1907)<<(uint(int32(31))%32)|(base.Simd_g_i32x4_extract_lane_l2(v1907)&v1779<<(uint(int32(30))%32)|(base.Simd_g_i32x4_extract_lane_l1(v1907)&v1779<<(uint(int32(29))%32)|(base.Simd_g_i32x4_extract_lane_l0(v1907)&v1779<<(uint(int32(28))%32)|(base.Simd_g_i32x4_extract_lane_l3(v1932)&v1779<<(uint(int32(27))%32)|(base.Simd_g_i32x4_extract_lane_l2(v1932)&v1779<<(uint(int32(26))%32)|(base.Simd_g_i32x4_extract_lane_l1(v1932)&v1779<<(uint(int32(25))%32)|(base.Simd_g_i32x4_extract_lane_l0(v1932)&v1779<<(uint(int32(24))%32)|(base.Simd_g_i32x4_extract_lane_l3(v1959)&v1779<<(uint(int32(23))%32)|(base.Simd_g_i32x4_extract_lane_l2(v1959)&v1779<<(uint(int32(22))%32)|(base.Simd_g_i32x4_extract_lane_l1(v1959)&v1779<<(uint(int32(21))%32)|(base.Simd_g_i32x4_extract_lane_l0(v1959)&v1779<<(uint(int32(20))%32)|(base.Simd_g_i32x4_extract_lane_l3(v1986)&v1779<<(uint(int32(19))%32)|(base.Simd_g_i32x4_extract_lane_l2(v1986)&v1779<<(uint(int32(18))%32)|(base.Simd_g_i32x4_extract_lane_l1(v1986)&v1779<<(uint(v1779)%32)-base.Simd_g_i32x4_extract_lane_l0(v1986))<<(uint(int32(16))%32))))))))))))))))+base.B2i32(v2026 != v1795)+(base.B2i32(v2030 != v1795)+base.B2i32(v2033 != v1795))+(base.B2i32(v2038 != v1795)+base.B2i32(v2041 != v1795)+(base.B2i32(v2045 != v1795)+base.B2i32(v2048 != v1795)))+base.B2i32(v2054 != v1795)) {
		goto L185
	} else {
		goto L186
	}
L159:
	;
	v1650 = base.Simd_g_i32x4_add(v1503, base.Simd_g_i8x16_swizzle_c(v1503, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k1)))
	v1655 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v1650, base.Simd_g_i8x16_swizzle_c(v1650, base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k2))))
	if base.Ui32(int32(255)) < base.Ui32(v1655) {
		goto L176
	} else {
		goto L177
	}
L160:
	;
	goto L163
L161:
	;
	goto L158
L163:
	;
	base.MemoryFill(m, v102, int32(0), int32(1024))
	goto L161
L175:
	;
	v1669 = int32(0)
	goto L178
L176:
	;
	v1663 = m.G118
	v1664 = *(*int32)(unsafe.Add(mBase, uint32(v1663)))
	v1665 = m.T0[v1664].(func(*base.Module, int32) int32)(m, v1655)
	mBase = m.M
	v1666 = v1665
	goto L175
L177:
	;
	v1658 = m.G117
	v1662 = *(*int32)(unsafe.Add(mBase, uint32(v1658+v1655<<(uint(int32(2))%32))))
	v1666 = v1662
	goto L175
L178:
	;
	v1712 = *(*int32)(unsafe.Add(mBase, uint32(v171+int32(2052)+v1669)))
	if base.Ui32(int32(255)) < base.Ui32(v1712) {
		goto L181
	} else {
		goto L182
	}
L179:
	;
	goto L158
L180:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v102+v1669))) = v1666 - v1723
	v1728 = v1669 + int32(4)
	if v1728 != int32(1024) {
		v1669 = v1728
		goto L178
	} else {
		goto L183
	}
L181:
	;
	v1720 = m.G118
	v1721 = *(*int32)(unsafe.Add(mBase, uint32(v1720)))
	v1722 = m.T0[v1721].(func(*base.Module, int32) int32)(m, v1712)
	mBase = m.M
	v1723 = v1722
	goto L180
L182:
	;
	v1715 = m.G117
	v1719 = *(*int32)(unsafe.Add(mBase, uint32(v1715+v1712<<(uint(int32(2))%32))))
	v1723 = v1719
	goto L180
L183:
	;
	goto L179
L184:
	;
	F_WebPSafeFree(m, v171)
	mBase = m.M
	goto L210
L185:
	;
	v2186 = int32(0)
	v2189 = int32(1)
	v2192 = int32(2)
	v2195 = int32(3)
	v2276 = v2026 + v1776 + v1778 + v1781 + v1784 + base.Simd_g_i32x4_extract_lane_l0(v1813) + base.Simd_g_i32x4_extract_lane_l1(v1813) + base.Simd_g_i32x4_extract_lane_l2(v1813) + base.Simd_g_i32x4_extract_lane_l3(v1813) + base.Simd_g_i32x4_extract_lane_l0(v1844) + base.Simd_g_i32x4_extract_lane_l1(v1844) + base.Simd_g_i32x4_extract_lane_l2(v1844) + base.Simd_g_i32x4_extract_lane_l3(v1844) + base.Simd_g_i32x4_extract_lane_l0(v1875) + base.Simd_g_i32x4_extract_lane_l1(v1875) + base.Simd_g_i32x4_extract_lane_l2(v1875) + base.Simd_g_i32x4_extract_lane_l3(v1875) + base.Simd_g_i32x4_extract_lane_l0(v1985) + base.Simd_g_i32x4_extract_lane_l1(v1985) + base.Simd_g_i32x4_extract_lane_l2(v1985) + base.Simd_g_i32x4_extract_lane_l3(v1985) + base.Simd_g_i32x4_extract_lane_l0(v1958) + base.Simd_g_i32x4_extract_lane_l1(v1958) + base.Simd_g_i32x4_extract_lane_l2(v1958) + base.Simd_g_i32x4_extract_lane_l3(v1958) + base.Simd_g_i32x4_extract_lane_l0(v1931) + base.Simd_g_i32x4_extract_lane_l1(v1931) + base.Simd_g_i32x4_extract_lane_l2(v1931) + base.Simd_g_i32x4_extract_lane_l3(v1931) + base.Simd_g_i32x4_extract_lane_l0(v1906) + base.Simd_g_i32x4_extract_lane_l1(v1906) + base.Simd_g_i32x4_extract_lane_l2(v1906) + base.Simd_g_i32x4_extract_lane_l3(v1906) + v2030 + v2033 + v2038 + v2041 + v2045 + v2048 + v2054
	if base.Ui32(int32(255)) < base.Ui32(v2276) {
		goto L202
	} else {
		goto L203
	}
L186:
	;
	goto L189
L187:
	;
	goto L184
L189:
	;
	base.MemoryFill(m, v1775, int32(0), int32(160))
	goto L187
L201:
	;
	v2292 = int32(0)
	goto L204
L202:
	;
	v2284 = m.G118
	v2285 = *(*int32)(unsafe.Add(mBase, uint32(v2284)))
	v2286 = m.T0[v2285].(func(*base.Module, int32) int32)(m, v2276)
	mBase = m.M
	v2287 = v2286
	goto L201
L203:
	;
	v2279 = m.G117
	v2283 = *(*int32)(unsafe.Add(mBase, uint32(v2279+v2276<<(uint(int32(2))%32))))
	v2287 = v2283
	goto L201
L204:
	;
	v2335 = *(*int32)(unsafe.Add(mBase, uint32(v171+int32(3076)+v2292)))
	if base.Ui32(int32(255)) < base.Ui32(v2335) {
		goto L207
	} else {
		goto L208
	}
L205:
	;
	goto L184
L206:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1775+v2292))) = v2287 - v2346
	v2351 = v2292 + int32(4)
	if v2351 != int32(160) {
		v2292 = v2351
		goto L204
	} else {
		goto L209
	}
L207:
	;
	v2343 = m.G118
	v2344 = *(*int32)(unsafe.Add(mBase, uint32(v2343)))
	v2345 = m.T0[v2344].(func(*base.Module, int32) int32)(m, v2335)
	mBase = m.M
	v2346 = v2345
	goto L206
L208:
	;
	v2338 = m.G117
	v2342 = *(*int32)(unsafe.Add(mBase, uint32(v2338+v2335<<(uint(int32(2))%32))))
	v2346 = v2342
	goto L206
L209:
	;
	goto L205
L210:
	;
	v2398 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+8)) = v2398
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0]))) = v2398
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1]))) = v2398
	v2406 = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v2406
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[3]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_5)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[4]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_6)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[5]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_7)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[6]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_8)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[7]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_9)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[8]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_10)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[9]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_11)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[10]))) = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)
	v2433 = v126 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[11]))) = v2433
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[12]))) = v2398
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13]))) = v71
	*(*int64)(unsafe.Add(mBase, uint32(v126))) = int64(0)
	v2440 = int32(4095)
	if v48 < v2440 {
		goto L211
	} else {
		goto L212
	}
L211:
	;
	v2443 = v48
	goto L213
L212:
	;
	v2443 = v2440
	goto L213
L213:
	;
	if int32(0) < v48 {
		goto L215
	} else {
		goto L216
	}
L214:
	;
	v2720 = int32(16)
	if v2714 == int64(0) {
		goto L240
	} else {
		goto L241
	}
L215:
	;
	v2449 = *(*int32)(unsafe.Add(mBase, uint32(v102)+3232))
	v2453 = m.G113
	v2455 = v126 + int32(16)
	v2457 = int32(0)
	v2470 = v2453
	v2474 = v2455
	goto L217
L216:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v126)+12)) = int32(1)
	v2691 = v2398
	v2714 = int64(1)
	goto L214
L217:
	;
	if base.Ui32(int32(511)) < base.Ui32(v2457) {
		goto L220
	} else {
		goto L221
	}
L218:
	;
	v2536 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+12)) = v2536
	if v48 != v2536 {
		goto L223
	} else {
		goto L224
	}
L219:
	;
	v2523 = int32(2)
	v2526 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v2449+int32(1024)+v2519<<(uint(v2523)%32)))))
	*(*int64)(unsafe.Add(mBase, uint32(v2474))) = base.I64_extend_i32_s(v2518)<<(uint(int64(23))%64) + v2526
	v2534 = v2457 + int32(1)
	if v2443 != v2534 {
		v2457 = v2534
		v2470 = v2470 + v2523
		v2474 = v2474 + int32(8)
		goto L217
	} else {
		goto L222
	}
L220:
	;
	v2505 = int32(-1)
	v2506 = v2457 + v2505
	v2509 = base.I32_clz(v2506) ^ int32(31)
	v2511 = v2509 + v2505
	v2513 = int32(1)
	v2518 = v2511
	v2519 = int32(base.Ui32(v2506)>>(uint(v2511)%32))&v2513 | v2509<<(uint(v2513)%32)
	goto L219
L221:
	;
	v2503 = int32(*(*int8)(unsafe.Add(mBase, uint32(v2470+int32(1)))))
	v2504 = int32(*(*int8)(unsafe.Add(mBase, uint32(v2470))))
	v2518 = v2503
	v2519 = v2504
	goto L219
L222:
	;
	goto L218
L223:
	;
	v2544 = v2443 + int32(-1)
	v2547 = *(*int64)(unsafe.Add(mBase, uint32(v2455)))
	if v48 != int32(2) {
		goto L226
	} else {
		goto L227
	}
L224:
	;
	v2691 = int32(0)
	v2714 = int64(1)
	goto L214
L225:
	;
	if v2544&int32(1) == int32(0) {
		v2674 = v2635
		goto L235
	} else {
		goto L236
	}
L226:
	;
	v2558 = v126 + int32(32)
	v2571 = int32(1)
	v2575 = int32(0)
	v2594 = v2547
	goto L228
L227:
	;
	v2622 = int32(1)
	v2635 = v2536
	v2658 = v2547
	goto L225
L228:
	;
	v2602 = *(*int64)(unsafe.Add(mBase, uint32(v2558+int32(-8))))
	if v2602 == v2594 {
		v2607 = v2571
		goto L230
	} else {
		goto L231
	}
L229:
	;
	v2622 = v2575 + int32(3)
	v2635 = v2613
	v2658 = v2608
	goto L225
L230:
	;
	v2608 = *(*int64)(unsafe.Add(mBase, uint32(v2558)))
	if v2608 == v2602 {
		v2613 = v2607
		goto L232
	} else {
		goto L233
	}
L231:
	;
	v2605 = v2571 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+12)) = v2605
	v2607 = v2605
	goto L230
L232:
	;
	v2617 = v2575 + int32(2)
	if v2544&int32(-2) != v2617 {
		v2558 = v2558 + int32(16)
		v2571 = v2613
		v2575 = v2617
		v2594 = v2608
		goto L228
	} else {
		goto L234
	}
L233:
	;
	v2611 = v2607 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+12)) = v2611
	v2613 = v2611
	goto L232
L234:
	;
	goto L229
L235:
	;
	v2691 = int32(1)
	v2714 = base.I64_extend_i32_u(v2674)
	goto L214
L236:
	;
	v2669 = *(*int64)(unsafe.Add(mBase, uint32(v2455+v2622<<(uint(int32(3))%32))))
	if v2669 == v2658 {
		v2674 = v2635
		goto L235
	} else {
		goto L237
	}
L237:
	;
	v2672 = v2635 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v126)+12)) = v2672
	v2674 = v2672
	goto L235
L238:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v126)+8)) = v2741
	if v2741 == int32(0) {
		goto L35
	} else {
		goto L244
	}
L239:
	;
	goto L238
L240:
	;
	v2739 = F_malloc(m, base.I32_wrap_i64(v2714)*v2720)
	mBase = m.M
	v2741 = v2739
	goto L239
L241:
	;
	v2727 = base.I64_div_u_s(int64(2147418112), v2714)
	v2728 = int32(0)
	v2729 = base.I64_extend_i32_u(v2720)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2729*v2714) {
		v2741 = v2728
		goto L239
	} else {
		goto L242
	}
L242:
	;
	if base.Ui64(v2727) < base.Ui64(v2729) {
		v2741 = v2728
		goto L239
	} else {
		goto L243
	}
L243:
	;
	goto L240
L244:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v2741)+8)) = int64(4294967296)
	v2748 = *(*int64)(unsafe.Add(mBase, uint32(v126)+16))
	*(*int64)(unsafe.Add(mBase, uint32(v2741))) = v2748
	if v2691 == int32(0) {
		goto L245
	} else {
		goto L246
	}
L245:
	;
	v2939 = int32(8)
	if v49 == int64(0) {
		goto L261
	} else {
		goto L262
	}
L246:
	;
	v2752 = int32(1)
	v2754 = v2443 + int32(-1)
	if v48 == int32(2) {
		v2837 = v2741
		v2850 = v2752
		v2873 = v2748
		goto L247
	} else {
		goto L248
	}
L247:
	;
	if v2754&v2752 == int32(0) {
		goto L245
	} else {
		goto L256
	}
L248:
	;
	v2765 = v2741
	v2778 = v126 + int32(32)
	v2782 = int32(0)
	v2801 = v2748
	goto L249
L249:
	;
	v2809 = *(*int64)(unsafe.Add(mBase, uint32(v2778+int32(-8))))
	if v2809 == v2801 {
		v2817 = v2765
		goto L251
	} else {
		goto L252
	}
L250:
	;
	v2837 = v2827
	v2850 = v2782 + int32(3)
	v2873 = v2821
	goto L247
L251:
	;
	v2819 = v2782 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v2817)+12)) = v2819
	v2821 = *(*int64)(unsafe.Add(mBase, uint32(v2778)))
	if v2821 == v2809 {
		v2827 = v2817
		goto L253
	} else {
		goto L254
	}
L252:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v2765)+16)) = v2809
	*(*int32)(unsafe.Add(mBase, uint32(v2765)+24)) = v2782 + int32(1)
	v2817 = v2765 + int32(16)
	goto L251
L253:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2827)+12)) = v2782 + int32(3)
	if v2754&int32(-2) != v2819 {
		v2765 = v2827
		v2778 = v2778 + int32(16)
		v2782 = v2819
		v2801 = v2821
		goto L249
	} else {
		goto L255
	}
L254:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v2817)+16)) = v2821
	*(*int32)(unsafe.Add(mBase, uint32(v2817)+24)) = v2819
	v2827 = v2817 + int32(16)
	goto L253
L255:
	;
	goto L250
L256:
	;
	v2886 = *(*int64)(unsafe.Add(mBase, uint32(v126+int32(16)+v2850<<(uint(int32(3))%32))))
	if v2886 == v2873 {
		v2892 = v2837
		goto L257
	} else {
		goto L258
	}
L257:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2892)+12)) = v2850 + int32(1)
	goto L245
L258:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v2837)+16)) = v2886
	*(*int32)(unsafe.Add(mBase, uint32(v2837)+24)) = v2850
	v2892 = v2837 + int32(16)
	goto L257
L259:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0]))) = v2960
	if v2960 == int32(0) {
		goto L35
	} else {
		goto L265
	}
L260:
	;
	goto L259
L261:
	;
	v2958 = F_malloc(m, base.I32_wrap_i64(v49)*v2939)
	mBase = m.M
	v2960 = v2958
	goto L260
L262:
	;
	v2946 = base.I64_div_u_s(int64(2147418112), v49)
	v2947 = int32(0)
	v2948 = base.I64_extend_i32_u(v2939)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2948*v49) {
		v2960 = v2947
		goto L260
	} else {
		goto L263
	}
L263:
	;
	if base.Ui64(v2946) < base.Ui64(v2948) {
		v2960 = v2947
		goto L260
	} else {
		goto L264
	}
L264:
	;
	goto L261
L265:
	;
	if v48 < int32(1) {
		goto L266
	} else {
		goto L267
	}
L266:
	;
	v3162 = int32(0)
	*(*uint16)(unsafe.Add(mBase, uint32(v71))) = uint16(v3162)
	v3164 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
	if l3 < int32(1) {
		goto L278
	} else {
		goto L279
	}
L267:
	;
	if v48 == int32(1) {
		v3038 = int32(0)
		goto L268
	} else {
		goto L269
	}
L268:
	;
	v3071 = v2960 + v3038<<(uint(int32(3))%32)
	v3084 = v48 - v3038
	goto L274
L269:
	;
	v2971 = v48 & int32(2147483646)
	v2973 = v2960
	v2986 = v2971
	goto L270
L270:
	;
	v3015 = base.Simd_g_const(&F_VP8LBackwardReferencesTraceBackwards__k3)
	base.Simd_g_v128_store(m, v2973, int32(0), v3015)
	v3021 = v2986 + int32(-2)
	if v3021 != 0 {
		v2973 = v2973 + int32(16)
		v2986 = v3021
		goto L270
	} else {
		goto L272
	}
L271:
	;
	if v48 == v2971 {
		goto L266
	} else {
		goto L273
	}
L272:
	;
	goto L271
L273:
	;
	v3038 = v2971
	goto L268
L274:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3071))) = int64(9223372036854775807)
	v3118 = v3084 + int32(-1)
	if v3118 != 0 {
		v3071 = v3071 + int32(8)
		v3084 = v3118
		goto L274
	} else {
		goto L276
	}
L275:
	;
	goto L266
L276:
	;
	goto L275
L277:
	;
	v3224 = *(*int64)(unsafe.Add(mBase, uint32(v2960)))
	v3228 = base.I64_div_u_s(v3223+int64(50), int64(100))
	if v3224 <= v3228 {
		goto L283
	} else {
		goto L284
	}
L278:
	;
	v3193 = int32(1020)
	v3196 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v837+int32(base.Ui32(v3164)>>(uint(int32(14))%32))&v3193))))
	v3202 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v102+int32(base.Ui32(v3164)>>(uint(int32(22))%32))&v3193))))
	v3204 = *(*int32)(unsafe.Add(mBase, uint32(v102)+3232))
	v3210 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3204+int32(base.Ui32(v3164)>>(uint(int32(6))%32))&v3193))))
	v3217 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1177+v3164&int32(255)<<(uint(int32(2))%32)))))
	v3223 = (v3196 + v3202 + v3210 + v3217) * int64(82)
	goto L277
L279:
	;
	v3167 = *(*int32)(unsafe.Add(mBase, uint32(v46)+68))
	v3170 = *(*int32)(unsafe.Add(mBase, uint32(v46)+72))
	v3171 = int32(base.Ui32(v3164*int32(506832829)) >> (uint(v3170) % 32))
	v3174 = v3167 + v3171<<(uint(int32(2))%32)
	if v3171 < int32(0) {
		goto L280
	} else {
		goto L281
	}
L280:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3174))) = v3164
	goto L278
L281:
	;
	v3177 = *(*int32)(unsafe.Add(mBase, uint32(v3174)))
	if v3177 != v3164 {
		goto L280
	} else {
		goto L282
	}
L282:
	;
	v3179 = *(*int32)(unsafe.Add(mBase, uint32(v102)+3232))
	v3185 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3179+v3171<<(uint(int32(2))%32)+int32(1120)))))
	v3223 = v3185 * int64(68)
	goto L277
L283:
	;
	if v48 < int32(2) {
		goto L285
	} else {
		goto L286
	}
L284:
	;
	v3230 = int32(1)
	*(*uint16)(unsafe.Add(mBase, uint32(v71))) = uint16(v3230)
	*(*int64)(unsafe.Add(mBase, uint32(v2960))) = v3228
	goto L283
L285:
	;
	v3985 = *(*int32)(unsafe.Add(mBase, uint32(l5)+4))
	v4128 = base.B2i32(v3985 == int32(0))
	goto L34
L286:
	;
	v3236 = int32(0)
	v3240 = int32(-1)
	v3252 = int64(-1)
	v3261 = int32(1)
	v3268 = v3240
	v3276 = v3240
	v3278 = v3236
	v3279 = int32(4)
	v3282 = v3236
	v3283 = v3240
	goto L287
L287:
	;
	v3286 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v3288 = v3261 << (uint(int32(2)) % 32)
	v3290 = *(*int32)(unsafe.Add(mBase, uint32(v3286+v3288)))
	v3292 = *(*int32)(unsafe.Add(mBase, uint32(l2+v3288)))
	v3293 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	v3295 = v3261 + int32(-1)
	v3299 = *(*int64)(unsafe.Add(mBase, uint32(v3293+v3295<<(uint(int32(3))%32))))
	if l3 < int32(1) {
		goto L290
	} else {
		goto L291
	}
L288:
	;
	goto L285
L289:
	;
	v3360 = v3290 & int32(4095)
	v3362 = v3261 << (uint(int32(3)) % 32)
	v3363 = v3293 + v3362
	v3364 = *(*int64)(unsafe.Add(mBase, uint32(v3363)))
	v3368 = base.I64_div_u_s(v3358+int64(50), int64(100))
	v3369 = v3368 + v3299
	if v3364 <= v3369 {
		goto L295
	} else {
		goto L296
	}
L290:
	;
	v3327 = int32(1020)
	v3330 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v837+int32(base.Ui32(v3292)>>(uint(int32(14))%32))&v3327))))
	v3336 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v102+int32(base.Ui32(v3292)>>(uint(int32(22))%32))&v3327))))
	v3338 = *(*int32)(unsafe.Add(mBase, uint32(v102)+3232))
	v3344 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3338+int32(base.Ui32(v3292)>>(uint(int32(6))%32))&v3327))))
	v3351 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1177+v3292&int32(255)<<(uint(int32(2))%32)))))
	v3358 = (v3330 + v3336 + v3344 + v3351) * int64(82)
	goto L289
L291:
	;
	v3302 = *(*int32)(unsafe.Add(mBase, uint32(v46)+68))
	v3305 = *(*int32)(unsafe.Add(mBase, uint32(v46)+72))
	v3306 = int32(base.Ui32(v3292*int32(506832829)) >> (uint(v3305) % 32))
	v3308 = v3306 << (uint(int32(2)) % 32)
	v3309 = v3302 + v3308
	if v3306 < int32(0) {
		goto L292
	} else {
		goto L293
	}
L292:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3309))) = v3292
	goto L290
L293:
	;
	v3312 = *(*int32)(unsafe.Add(mBase, uint32(v3309)))
	if v3312 != v3292 {
		goto L292
	} else {
		goto L294
	}
L294:
	;
	v3314 = *(*int32)(unsafe.Add(mBase, uint32(v102)+3232))
	v3318 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3314+v3308+int32(1120)))))
	v3358 = v3318 * int64(68)
	goto L289
L295:
	;
	v3378 = int32(base.Ui32(v3290) >> (uint(int32(12)) % 32))
	if base.Ui32(v3360) < base.Ui32(int32(2)) {
		v3770 = v3252
		v3800 = v3282
		v3801 = v3283
		goto L297
	} else {
		goto L298
	}
L296:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3363))) = v3369
	v3372 = int32(1)
	*(*uint16)(unsafe.Add(mBase, uint32(v71+v3261<<(uint(v3372)%32)))) = uint16(v3372)
	goto L295
L297:
	;
	v3804 = *(*int32)(unsafe.Add(mBase, uint32(v126)))
	if v3804 == int32(0) {
		goto L340
	} else {
		goto L341
	}
L298:
	;
	if v3378 == v3276 {
		goto L299
	} else {
		goto L300
	}
L299:
	;
	if v3283 != 0 {
		goto L311
	} else {
		goto L312
	}
L300:
	;
	v3384 = base.I32_div_s(v3378, l0)
	v3386 = v3378 - v3384*l0
	if int32(7) < v3384 {
		goto L304
	} else {
		goto L305
	}
L301:
	;
	v3453 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1775+v3446<<(uint(int32(2))%32)))))
	v3454 = base.I64_extend_i32_s(v3445)<<(uint(int64(23))%64) + v3453
	F_PushInterval(m, v126, v3454+v3299, v3261, v3360)
	mBase = m.M
	v3770 = v3454
	v3800 = v3282
	v3801 = int32(1)
	goto L297
L302:
	;
	v3432 = int32(-1)
	v3433 = v3423 + v3432
	v3436 = base.I32_clz(v3433) ^ int32(31)
	v3438 = v3436 + v3432
	v3440 = int32(1)
	v3445 = v3438
	v3446 = int32(base.Ui32(v3433)>>(uint(v3438)%32))&v3440 | v3436<<(uint(v3440)%32)
	goto L301
L303:
	;
	if int32(511) < v3423 {
		goto L302
	} else {
		goto L310
	}
L304:
	;
	if int32(6) < v3384 {
		goto L307
	} else {
		goto L308
	}
L305:
	;
	if int32(8) < v3386 {
		goto L304
	} else {
		goto L306
	}
L306:
	;
	v3391 = m.G1
	v3400 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3391+int32(_a_F_VP8LBackwardReferencesTraceBackwards_14)+(v3384<<(uint(int32(4))%32)|int32(8)-v3386)))))
	v3423 = v3400 + int32(1)
	goto L303
L307:
	;
	v3423 = v3378 + int32(120)
	goto L303
L308:
	;
	if v3386 <= l0+int32(-8) {
		goto L307
	} else {
		goto L309
	}
L309:
	;
	v3412 = m.G1
	v3418 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v3384<<(uint(int32(4))%32)-v3386+(v3412+int32(_a_F_VP8LBackwardReferencesTraceBackwards_14))+int32(24)))))
	v3423 = v3418 + int32(1)
	goto L303
L310:
	;
	v3426 = m.G113
	v3429 = v3426 + v3423<<(uint(int32(1))%32)
	v3430 = int32(*(*int8)(unsafe.Add(mBase, uint32(v3429)+1)))
	v3431 = int32(*(*int8)(unsafe.Add(mBase, uint32(v3429))))
	v3445 = v3430
	v3446 = v3431
	goto L301
L311:
	;
	v3461 = v3261 + v3268 + int32(-2)
	goto L313
L312:
	;
	v3461 = v3282
	goto L313
L313:
	;
	if v3360+v3295 <= v3461 {
		v3770 = v3252
		v3800 = v3461
		v3801 = int32(0)
		goto L297
	} else {
		goto L314
	}
L314:
	;
	if v3461 < v3261 {
		v3531 = v3261
		v3570 = int32(0)
		goto L315
	} else {
		goto L316
	}
L315:
	;
	v3574 = v3531 + int32(-1)
	v3575 = *(*int32)(unsafe.Add(mBase, uint32(v126)))
	if v3575 == int32(0) {
		goto L323
	} else {
		goto L324
	}
L316:
	;
	v3486 = v3286 + v3279
	v3489 = int32(0)
	goto L319
L317:
	;
	v3531 = v3461 + int32(1)
	v3570 = v3517 & int32(4095)
	goto L315
L318:
	;
	v3525 = *(*int32)(unsafe.Add(mBase, uint32(v3486)))
	v3531 = v3261 - v3489
	v3570 = v3525 & int32(4095)
	goto L315
L319:
	;
	v3516 = v3486 + int32(4)
	v3517 = *(*int32)(unsafe.Add(mBase, uint32(v3516)))
	if int32(base.Ui32(v3517)>>(uint(int32(12))%32)) != v3276 {
		goto L318
	} else {
		goto L321
	}
L321:
	;
	v3522 = v3489 + int32(-1)
	if v3278-v3461 == v3522 {
		goto L317
	} else {
		goto L322
	}
L322:
	;
	v3486 = v3516
	v3489 = v3522
	goto L319
L323:
	;
	v3756 = *(*int64)(unsafe.Add(mBase, uint32(v3293+v3574<<(uint(int32(3))%32))))
	F_PushInterval(m, v126, v3756+v3252, v3531, v3570)
	mBase = m.M
	v3770 = v3252
	v3800 = v3574 + v3570
	v3801 = int32(0)
	goto L297
L324:
	;
	v3580 = v3293 + v3574<<(uint(int32(3))%32)
	v3598 = v3575
	goto L326
L325:
	;
	v3647 = v3293 + v3531<<(uint(int32(3))%32)
	v3663 = v3575
	goto L333
L326:
	;
	v3624 = *(*int32)(unsafe.Add(mBase, uint32(v3598)+8))
	if v3531 <= v3624 {
		goto L325
	} else {
		goto L328
	}
L327:
	;
	goto L325
L328:
	;
	v3626 = *(*int32)(unsafe.Add(mBase, uint32(v3598)+24))
	v3627 = *(*int32)(unsafe.Add(mBase, uint32(v3598)+12))
	if v3627 < v3531 {
		goto L329
	} else {
		goto L330
	}
L329:
	;
	if v3626 != 0 {
		v3598 = v3626
		goto L326
	} else {
		goto L332
	}
L330:
	;
	v3629 = *(*int64)(unsafe.Add(mBase, uint32(v3580)))
	v3630 = *(*int64)(unsafe.Add(mBase, uint32(v3598)))
	if v3629 <= v3630 {
		goto L329
	} else {
		goto L331
	}
L331:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3580))) = v3630
	v3633 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13])))
	v3634 = int32(1)
	v3637 = *(*int32)(unsafe.Add(mBase, uint32(v3598)+16))
	v3640 = v3574 - v3637 + v3634
	*(*uint16)(unsafe.Add(mBase, uint32(v3633+v3574<<(uint(v3634)%32)))) = uint16(v3640)
	goto L329
L332:
	;
	goto L327
L333:
	;
	v3691 = *(*int32)(unsafe.Add(mBase, uint32(v3663)+8))
	if v3531 < v3691 {
		goto L323
	} else {
		goto L335
	}
L334:
	;
	goto L323
L335:
	;
	v3693 = *(*int32)(unsafe.Add(mBase, uint32(v3663)+24))
	v3694 = *(*int32)(unsafe.Add(mBase, uint32(v3663)+12))
	if v3694 <= v3531 {
		goto L336
	} else {
		goto L337
	}
L336:
	;
	if v3693 != 0 {
		v3663 = v3693
		goto L333
	} else {
		goto L339
	}
L337:
	;
	v3696 = *(*int64)(unsafe.Add(mBase, uint32(v3647)))
	v3697 = *(*int64)(unsafe.Add(mBase, uint32(v3663)))
	if v3696 <= v3697 {
		goto L336
	} else {
		goto L338
	}
L338:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3647))) = v3697
	v3700 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13])))
	v3701 = int32(1)
	v3704 = *(*int32)(unsafe.Add(mBase, uint32(v3663)+16))
	v3707 = v3531 - v3704 + v3701
	*(*uint16)(unsafe.Add(mBase, uint32(v3700+v3531<<(uint(v3701)%32)))) = uint16(v3707)
	goto L336
L339:
	;
	goto L334
L340:
	;
	v3935 = int32(1)
	v3940 = v3261 + v3935
	if v3940 != v48 {
		v3252 = v3770
		v3261 = v3940
		v3268 = v3360
		v3276 = v3378
		v3278 = v3278 + v3935
		v3279 = v3279 + int32(4)
		v3282 = v3800
		v3283 = v3801
		goto L287
	} else {
		goto L359
	}
L341:
	;
	v3821 = v3804
	goto L342
L342:
	;
	v3850 = *(*int32)(unsafe.Add(mBase, uint32(v3821)+8))
	if v3261 < v3850 {
		goto L340
	} else {
		goto L344
	}
L343:
	;
	goto L340
L344:
	;
	v3852 = *(*int32)(unsafe.Add(mBase, uint32(v3821)+24))
	v3853 = *(*int32)(unsafe.Add(mBase, uint32(v3821)+12))
	if v3261 < v3853 {
		goto L346
	} else {
		goto L347
	}
L345:
	;
	if v3852 != 0 {
		v3821 = v3852
		goto L342
	} else {
		goto L358
	}
L346:
	;
	v3875 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	v3876 = v3875 + v3362
	v3877 = *(*int64)(unsafe.Add(mBase, uint32(v3876)))
	v3878 = *(*int64)(unsafe.Add(mBase, uint32(v3821)))
	if v3877 <= v3878 {
		goto L345
	} else {
		goto L357
	}
L347:
	;
	v3855 = *(*int32)(unsafe.Add(mBase, uint32(v3821)+20))
	if v3855 != 0 {
		goto L348
	} else {
		goto L349
	}
L348:
	;
	v3858 = v3855 + int32(24)
	goto L350
L349:
	;
	v3858 = v126
	goto L350
L350:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3858))) = v3852
	if v3852 == int32(0) {
		goto L351
	} else {
		goto L352
	}
L351:
	;
	if base.Ui32(v3821) < base.Ui32(v2433) {
		goto L354
	} else {
		goto L355
	}
L352:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3852)+20)) = v3855
	goto L351
L353:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3821)+24)) = v3869
	v3871 = *(*int32)(unsafe.Add(mBase, uint32(v126)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v126)+4)) = v3871 + int32(-1)
	goto L345
L354:
	;
	v3867 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1])))
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1]))) = v3821
	v3869 = v3867
	goto L353
L355:
	;
	if base.Ui32(v2406) < base.Ui32(v3821) {
		goto L354
	} else {
		goto L356
	}
L356:
	;
	v3865 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v3821
	v3869 = v3865
	goto L353
L357:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3876))) = v3878
	v3881 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13])))
	v3882 = int32(1)
	v3885 = *(*int32)(unsafe.Add(mBase, uint32(v3821)+16))
	v3888 = v3261 - v3885 + v3882
	*(*uint16)(unsafe.Add(mBase, uint32(v3881+v3261<<(uint(v3882)%32)))) = uint16(v3888)
	goto L345
L358:
	;
	goto L343
L359:
	;
	goto L288
L360:
	;
	v4128 = v3988
	goto L34
L361:
	;
	v4128 = int32(0)
	goto L34
L362:
	;
	goto L361
L363:
	;
	v4040 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	F_WebPSafeFree(m, v4040)
	mBase = m.M
	v4042 = *(*int32)(unsafe.Add(mBase, uint32(v126)+8))
	F_WebPSafeFree(m, v4042)
	mBase = m.M
	v4044 = *(*int32)(unsafe.Add(mBase, uint32(v126)))
	if v4044 == int32(0) {
		goto L364
	} else {
		goto L365
	}
L364:
	;
	v4065 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v126))) = v4065
	v4067 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1])))
	if v4067 == v4065 {
		goto L373
	} else {
		goto L374
	}
L365:
	;
	v4052 = v4044
	goto L366
L366:
	;
	v4056 = *(*int32)(unsafe.Add(mBase, uint32(v4052)+24))
	if base.Ui32(v4052) < base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)) {
		goto L369
	} else {
		goto L370
	}
L367:
	;
	goto L364
L368:
	;
	if v4056 != 0 {
		v4052 = v4056
		goto L366
	} else {
		goto L372
	}
L369:
	;
	F_WebPSafeFree(m, v4052)
	mBase = m.M
	goto L368
L370:
	;
	if base.Ui32(v4052) <= base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)) {
		goto L368
	} else {
		goto L371
	}
L371:
	;
	goto L369
L372:
	;
	goto L367
L373:
	;
	v4090 = F_memset(m, v126, int32(0), int32(_a_F_VP8LBackwardReferencesTraceBackwards_1))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[3]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_5)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[4]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_6)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[5]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_7)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[6]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_8)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[7]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_9)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[8]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_10)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[9]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_11)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[10]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)
	*(*int32)(unsafe.Add(mBase, uint32(v4090)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[11]))) = v4090 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)
	goto L362
L374:
	;
	v4075 = v4067
	goto L375
L375:
	;
	v4079 = *(*int32)(unsafe.Add(mBase, uint32(v4075)+24))
	if base.Ui32(v4075) < base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)) {
		goto L378
	} else {
		goto L379
	}
L376:
	;
	goto L373
L377:
	;
	if v4079 != 0 {
		v4075 = v4079
		goto L375
	} else {
		goto L381
	}
L378:
	;
	F_WebPSafeFree(m, v4075)
	mBase = m.M
	goto L377
L379:
	;
	if base.Ui32(v4075) <= base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)) {
		goto L377
	} else {
		goto L380
	}
L380:
	;
	goto L378
L381:
	;
	goto L376
L382:
	;
	v4173 = v46 + int32(68)
	if v4173 == int32(0) {
		goto L384
	} else {
		goto L385
	}
L383:
	;
	v4181 = v4128
	goto L25
L384:
	;
	goto L383
L385:
	;
	v4176 = *(*int32)(unsafe.Add(mBase, uint32(v4173)))
	F_WebPSafeFree(m, v4176)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4173))) = int32(0)
	goto L384
L386:
	;
	F_free(m, v102)
	mBase = m.M
	goto L407
L387:
	;
	goto L386
L388:
	;
	v4229 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	F_WebPSafeFree(m, v4229)
	mBase = m.M
	v4231 = *(*int32)(unsafe.Add(mBase, uint32(v126)+8))
	F_WebPSafeFree(m, v4231)
	mBase = m.M
	v4233 = *(*int32)(unsafe.Add(mBase, uint32(v126)))
	if v4233 == int32(0) {
		goto L389
	} else {
		goto L390
	}
L389:
	;
	v4254 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v126))) = v4254
	v4256 = *(*int32)(unsafe.Add(mBase, uint32(v126)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1])))
	if v4256 == v4254 {
		goto L398
	} else {
		goto L399
	}
L390:
	;
	v4241 = v4233
	goto L391
L391:
	;
	v4245 = *(*int32)(unsafe.Add(mBase, uint32(v4241)+24))
	if base.Ui32(v4241) < base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)) {
		goto L394
	} else {
		goto L395
	}
L392:
	;
	goto L389
L393:
	;
	if v4245 != 0 {
		v4241 = v4245
		goto L391
	} else {
		goto L397
	}
L394:
	;
	F_WebPSafeFree(m, v4241)
	mBase = m.M
	goto L393
L395:
	;
	if base.Ui32(v4241) <= base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)) {
		goto L393
	} else {
		goto L396
	}
L396:
	;
	goto L394
L397:
	;
	goto L392
L398:
	;
	v4279 = F_memset(m, v126, int32(0), int32(_a_F_VP8LBackwardReferencesTraceBackwards_1))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[3]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_5)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[4]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_6)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[5]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_7)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[6]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_8)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[7]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_9)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[8]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_10)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[9]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_11)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[10]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)
	*(*int32)(unsafe.Add(mBase, uint32(v4279)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[11]))) = v4279 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)
	goto L387
L399:
	;
	v4264 = v4256
	goto L400
L400:
	;
	v4268 = *(*int32)(unsafe.Add(mBase, uint32(v4264)+24))
	if base.Ui32(v4264) < base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)) {
		goto L403
	} else {
		goto L404
	}
L401:
	;
	goto L398
L402:
	;
	if v4268 != 0 {
		v4264 = v4268
		goto L400
	} else {
		goto L406
	}
L403:
	;
	F_WebPSafeFree(m, v4264)
	mBase = m.M
	goto L402
L404:
	;
	if base.Ui32(v4264) <= base.Ui32(v126+int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)) {
		goto L402
	} else {
		goto L405
	}
L405:
	;
	goto L403
L406:
	;
	goto L401
L407:
	;
	F_free(m, v126)
	mBase = m.M
	goto L408
L408:
	;
	if v4181 == int32(0) {
		v5097 = int32(0)
		goto L1
	} else {
		goto L409
	}
L409:
	;
	v4321 = v71 + v48<<(uint(int32(1))%32)
	v4323 = v4321 + int32(-2)
	if base.Ui32(v4323) < base.Ui32(v71) {
		v4391 = v4321
		goto L410
	} else {
		goto L411
	}
L410:
	;
	v4420 = int32(1)
	v4421 = (v4321 - v4391) >> (uint(v4420) % 32)
	if l3 < v4420 {
		goto L420
	} else {
		goto L421
	}
L411:
	;
	v4326 = v4323
	v4340 = v4321
	goto L412
L412:
	;
	v4369 = v4340 + int32(-2)
	v4370 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4326))))
	*(*uint16)(unsafe.Add(mBase, uint32(v4369))) = uint16(v4370)
	v4374 = v4326 - v4370<<(uint(int32(1))%32)
	if base.Ui32(v71) <= base.Ui32(v4374) {
		v4326 = v4374
		v4340 = v4369
		goto L412
	} else {
		goto L414
	}
L413:
	;
	v4391 = v4369
	goto L410
L414:
	;
	goto L413
L415:
	;
	v5079 = v46 + int32(68)
	if v5079 == int32(0) {
		goto L497
	} else {
		goto L498
	}
L416:
	;
	v5030 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	v5032 = base.B2i32(v5030 == int32(0))
	if l3 < int32(1) {
		v5097 = v5032
		goto L1
	} else {
		goto L495
	}
L417:
	;
	v4677 = v4477
	v4685 = int32(0)
	goto L460
L418:
	;
	v4484 = v4456
	v4498 = v4391
	v4502 = v4421
	goto L434
L419:
	;
	v4463 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	if v4463 == int32(0) {
		goto L431
	} else {
		goto L432
	}
L420:
	;
	v4442 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	if v4442 == int32(0) {
		goto L427
	} else {
		goto L428
	}
L421:
	;
	v4425 = v46 + int32(68)
	v4431 = F_WebPSafeCalloc(m, base.I64_extend_i32_s(int32(1)<<(uint(l3)%32)), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4425))) = v4431
	if v4431 != 0 {
		goto L423
	} else {
		goto L424
	}
L422:
	;
	if v4431 != 0 {
		goto L419
	} else {
		goto L425
	}
L423:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4425)+8)) = l3
	*(*int32)(unsafe.Add(mBase, uint32(v4425)+4)) = int32(32) - l3
	goto L422
L424:
	;
	goto L422
L425:
	;
	v5097 = int32(0)
	goto L1
L426:
	;
	v4456 = int32(0)
	if v4456 < v4421 {
		goto L418
	} else {
		goto L429
	}
L427:
	;
	v4447 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v4447
	v4449 = *(*int32)(unsafe.Add(mBase, uint32(l6)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v4449
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = l6 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+8)) = v4447
	goto L426
L428:
	;
	v4445 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v4442))) = v4445
	goto L427
L429:
	;
	v4459 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	v5097 = base.B2i32(v4459 == int32(0))
	goto L1
L430:
	;
	v4477 = int32(0)
	if v4477 < v4421 {
		goto L417
	} else {
		goto L433
	}
L431:
	;
	v4468 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v4468
	v4470 = *(*int32)(unsafe.Add(mBase, uint32(l6)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v4470
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = l6 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+8)) = v4468
	goto L430
L432:
	;
	v4466 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v4463))) = v4466
	goto L431
L433:
	;
	v4480 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	v5046 = base.B2i32(v4480 == int32(0))
	goto L415
L434:
	;
	v4526 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4498))))
	if v4526 == int32(1) {
		goto L437
	} else {
		goto L438
	}
L436:
	;
	v4659 = v4502 + int32(-1)
	if v4659 != 0 {
		v4484 = v4484 + v4526
		v4498 = v4498 + int32(2)
		v4502 = v4659
		goto L434
	} else {
		goto L459
	}
L437:
	;
	v4596 = int64(*(*uint32)(unsafe.Add(mBase, uint32(l2+v4484<<(uint(int32(2))%32)))))
	v4600 = v4596<<(uint(int64(32))%64) | int64(65536)
	*(*int64)(unsafe.Add(mBase, uint32(v46)+48)) = v4600
	*(*int64)(unsafe.Add(mBase, uint32(v46)+8)) = v4600
	v4607 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v4607 == int32(0) {
		goto L451
	} else {
		goto L452
	}
L438:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v46)+62)) = uint16(v4526)
	v4530 = int32(2)
	*(*uint8)(unsafe.Add(mBase, uint32(v46)+60)) = uint8(v4530)
	v4532 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v4536 = *(*int32)(unsafe.Add(mBase, uint32(v4532+v4484<<(uint(v4530)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+64)) = int32(base.Ui32(v4536) >> (uint(int32(12)) % 32))
	v4540 = *(*int64)(unsafe.Add(mBase, uint32(v46)+60))
	*(*int64)(unsafe.Add(mBase, uint32(v46)+16)) = v4540
	v4546 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v4546 == int32(0) {
		goto L441
	} else {
		goto L442
	}
L439:
	;
	goto L436
L440:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4580)+8)) = v4581 + int32(1)
	v4585 = *(*int32)(unsafe.Add(mBase, uint32(v4580)+4))
	v4589 = *(*int64)(unsafe.Add(mBase, uint32(v46+int32(16))))
	*(*int64)(unsafe.Add(mBase, uint32(v4585+v4581<<(uint(int32(3))%32)))) = v4589
	goto L439
L441:
	;
	v4553 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v4553 != 0 {
		goto L445
	} else {
		goto L446
	}
L442:
	;
	v4549 = *(*int32)(unsafe.Add(mBase, uint32(v4546)+8))
	v4550 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v4549 != v4550 {
		v4580 = v4546
		v4581 = v4549
		goto L440
	} else {
		goto L443
	}
L443:
	;
	goto L441
L444:
	;
	v4571 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v4571))) = v4570
	v4573 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v4570)+8)) = v4573
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v4570
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v4570
	*(*int32)(unsafe.Add(mBase, uint32(v4570))) = v4573
	v4580 = v4570
	v4581 = v4573
	goto L440
L445:
	;
	v4568 = *(*int32)(unsafe.Add(mBase, uint32(v4553)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v4568
	v4570 = v4553
	goto L444
L446:
	;
	v4555 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v4560 = F_WebPSafeMalloc(m, int64(1), v4555<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v4560 != 0 {
		goto L447
	} else {
		goto L448
	}
L447:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4560)+4)) = v4560 + int32(12)
	v4570 = v4560
	goto L444
L448:
	;
	v4561 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v4561 | int32(1)
	goto L439
L449:
	;
	goto L436
L450:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4641)+8)) = v4642 + int32(1)
	v4646 = *(*int32)(unsafe.Add(mBase, uint32(v4641)+4))
	v4650 = *(*int64)(unsafe.Add(mBase, uint32(v46+int32(8))))
	*(*int64)(unsafe.Add(mBase, uint32(v4646+v4642<<(uint(int32(3))%32)))) = v4650
	goto L449
L451:
	;
	v4614 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v4614 != 0 {
		goto L455
	} else {
		goto L456
	}
L452:
	;
	v4610 = *(*int32)(unsafe.Add(mBase, uint32(v4607)+8))
	v4611 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v4610 != v4611 {
		v4641 = v4607
		v4642 = v4610
		goto L450
	} else {
		goto L453
	}
L453:
	;
	goto L451
L454:
	;
	v4632 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v4632))) = v4631
	v4634 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v4631)+8)) = v4634
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v4631
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v4631
	*(*int32)(unsafe.Add(mBase, uint32(v4631))) = v4634
	v4641 = v4631
	v4642 = v4634
	goto L450
L455:
	;
	v4629 = *(*int32)(unsafe.Add(mBase, uint32(v4614)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v4629
	v4631 = v4614
	goto L454
L456:
	;
	v4616 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v4621 = F_WebPSafeMalloc(m, int64(1), v4616<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v4621 != 0 {
		goto L457
	} else {
		goto L458
	}
L457:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4621)+4)) = v4621 + int32(12)
	v4631 = v4621
	goto L454
L458:
	;
	v4622 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v4622 | int32(1)
	goto L449
L459:
	;
	goto L416
L460:
	;
	v4704 = int32(1)
	v4707 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4391+v4677<<(uint(v4704)%32)))))
	if v4707 == v4704 {
		goto L463
	} else {
		goto L464
	}
L461:
	;
	goto L416
L462:
	;
	v4985 = v4677 + int32(1)
	if v4985 != v4421 {
		v4677 = v4985
		v4685 = v4685 + v4707
		goto L460
	} else {
		goto L494
	}
L463:
	;
	v4860 = *(*int32)(unsafe.Add(mBase, uint32(v46)+68))
	v4861 = int32(2)
	v4863 = l2 + v4685<<(uint(v4861)%32)
	v4864 = *(*int32)(unsafe.Add(mBase, uint32(v4863)))
	v4867 = *(*int32)(unsafe.Add(mBase, uint32(v46)+72))
	v4868 = int32(base.Ui32(v4864*int32(506832829)) >> (uint(v4867) % 32))
	v4871 = v4860 + v4868<<(uint(v4861)%32)
	v4872 = *(*int32)(unsafe.Add(mBase, uint32(v4871)))
	if v4872 != v4864 {
		goto L481
	} else {
		goto L482
	}
L464:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v46)+62)) = uint16(v4707)
	v4711 = int32(2)
	*(*uint8)(unsafe.Add(mBase, uint32(v46)+60)) = uint8(v4711)
	v4713 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v4715 = v4685 << (uint(v4711) % 32)
	v4717 = *(*int32)(unsafe.Add(mBase, uint32(v4713+v4715)))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+64)) = int32(base.Ui32(v4717) >> (uint(int32(12)) % 32))
	v4721 = *(*int64)(unsafe.Add(mBase, uint32(v46)+60))
	*(*int64)(unsafe.Add(mBase, uint32(v46)+32)) = v4721
	v4727 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v4727 == int32(0) {
		goto L467
	} else {
		goto L468
	}
L465:
	;
	if v4707 == int32(0) {
		goto L462
	} else {
		goto L475
	}
L466:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4761)+8)) = v4762 + int32(1)
	v4766 = *(*int32)(unsafe.Add(mBase, uint32(v4761)+4))
	v4770 = *(*int64)(unsafe.Add(mBase, uint32(v46+int32(32))))
	*(*int64)(unsafe.Add(mBase, uint32(v4766+v4762<<(uint(int32(3))%32)))) = v4770
	goto L465
L467:
	;
	v4734 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v4734 != 0 {
		goto L471
	} else {
		goto L472
	}
L468:
	;
	v4730 = *(*int32)(unsafe.Add(mBase, uint32(v4727)+8))
	v4731 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v4730 != v4731 {
		v4761 = v4727
		v4762 = v4730
		goto L466
	} else {
		goto L469
	}
L469:
	;
	goto L467
L470:
	;
	v4752 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v4752))) = v4751
	v4754 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v4751)+8)) = v4754
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v4751
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v4751
	*(*int32)(unsafe.Add(mBase, uint32(v4751))) = v4754
	v4761 = v4751
	v4762 = v4754
	goto L466
L471:
	;
	v4749 = *(*int32)(unsafe.Add(mBase, uint32(v4734)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v4749
	v4751 = v4734
	goto L470
L472:
	;
	v4736 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v4741 = F_WebPSafeMalloc(m, int64(1), v4736<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v4741 != 0 {
		goto L473
	} else {
		goto L474
	}
L473:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4741)+4)) = v4741 + int32(12)
	v4751 = v4741
	goto L470
L474:
	;
	v4742 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v4742 | int32(1)
	goto L465
L475:
	;
	v4781 = *(*int32)(unsafe.Add(mBase, uint32(v46)+68))
	v4783 = l2 + v4715
	v4795 = v4707 & int32(_a_F_VP8LBackwardReferencesTraceBackwards_15)
	goto L476
L476:
	;
	v4825 = *(*int32)(unsafe.Add(mBase, uint32(v4783)))
	v4826 = int32(506832829)
	v4828 = *(*int32)(unsafe.Add(mBase, uint32(v46)+72))
	v4830 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v4781+int32(base.Ui32(v4825*v4826)>>(uint(v4828)%32))<<(uint(v4830)%32)))) = v4825
	v4836 = *(*int32)(unsafe.Add(mBase, uint32(v4783+int32(4))))
	v4839 = *(*int32)(unsafe.Add(mBase, uint32(v46)+72))
	*(*int32)(unsafe.Add(mBase, uint32(v4781+int32(base.Ui32(v4836*v4826)>>(uint(v4839)%32))<<(uint(v4830)%32)))) = v4836
	v4846 = v4783 + int32(8)
	v4848 = v4795 + int32(-2)
	if v4848 != 0 {
		v4783 = v4846
		v4795 = v4848
		goto L476
	} else {
		goto L478
	}
L477:
	;
	if v4707&int32(1) == int32(0) {
		goto L462
	} else {
		goto L479
	}
L478:
	;
	goto L477
L479:
	;
	v4851 = *(*int32)(unsafe.Add(mBase, uint32(v4846)))
	v4854 = *(*int32)(unsafe.Add(mBase, uint32(v46)+72))
	*(*int32)(unsafe.Add(mBase, uint32(v4781+int32(base.Ui32(v4851*int32(506832829))>>(uint(v4854)%32))<<(uint(int32(2))%32)))) = v4851
	goto L462
L480:
	;
	v4886 = base.I64_extend_i32_u(v4881)<<(uint(int64(32))%64) | v4882
	*(*int64)(unsafe.Add(mBase, uint32(v46)+24)) = v4886
	*(*int64)(unsafe.Add(mBase, uint32(v46)+48)) = v4886
	v4893 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v4893 == int32(0) {
		goto L486
	} else {
		goto L487
	}
L481:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4871))) = v4864
	v4879 = *(*int32)(unsafe.Add(mBase, uint32(v4863)))
	v4881 = v4879
	v4882 = int64(65536)
	goto L480
L482:
	;
	if int32(-1) < v4868 {
		v4881 = v4868
		v4882 = int64(65537)
		goto L480
	} else {
		goto L483
	}
L483:
	;
	goto L481
L484:
	;
	goto L462
L485:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4927)+8)) = v4928 + int32(1)
	v4932 = *(*int32)(unsafe.Add(mBase, uint32(v4927)+4))
	v4936 = *(*int64)(unsafe.Add(mBase, uint32(v46+int32(24))))
	*(*int64)(unsafe.Add(mBase, uint32(v4932+v4928<<(uint(int32(3))%32)))) = v4936
	goto L484
L486:
	;
	v4900 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v4900 != 0 {
		goto L490
	} else {
		goto L491
	}
L487:
	;
	v4896 = *(*int32)(unsafe.Add(mBase, uint32(v4893)+8))
	v4897 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v4896 != v4897 {
		v4927 = v4893
		v4928 = v4896
		goto L485
	} else {
		goto L488
	}
L488:
	;
	goto L486
L489:
	;
	v4918 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v4918))) = v4917
	v4920 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v4917)+8)) = v4920
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v4917
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v4917
	*(*int32)(unsafe.Add(mBase, uint32(v4917))) = v4920
	v4927 = v4917
	v4928 = v4920
	goto L485
L490:
	;
	v4915 = *(*int32)(unsafe.Add(mBase, uint32(v4900)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v4915
	v4917 = v4900
	goto L489
L491:
	;
	v4902 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v4907 = F_WebPSafeMalloc(m, int64(1), v4902<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v4907 != 0 {
		goto L492
	} else {
		goto L493
	}
L492:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4907)+4)) = v4907 + int32(12)
	v4917 = v4907
	goto L489
L493:
	;
	v4908 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v4908 | int32(1)
	goto L484
L494:
	;
	goto L461
L495:
	;
	v5046 = v5032
	goto L415
L496:
	;
	v5097 = v5046
	goto L1
L497:
	;
	goto L496
L498:
	;
	v5082 = *(*int32)(unsafe.Add(mBase, uint32(v5079)))
	F_WebPSafeFree(m, v5082)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v5079))) = int32(0)
	goto L497
L499:
	;
	m.G0 = v46 + int32(80)
	return v5097
}

var F_VP8LBackwardReferencesTraceBackwards__k0 = [2]uint64{0x0, 0x0}
var F_VP8LBackwardReferencesTraceBackwards__k1 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_VP8LBackwardReferencesTraceBackwards__k2 = [2]uint64{0x302010007060504, 0x302010003020100}
var F_VP8LBackwardReferencesTraceBackwards__k3 = [2]uint64{0x7fffffffffffffff, 0x7fffffffffffffff}

func F_VP8LBackwardRefsInit(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 base.V128
	_ = v3
	var v13 int32
	_ = v13
	var v16 int32
	_ = v16
	v3 = base.Simd_g_const(&F_VP8LBackwardRefsInit__k0)
	base.Simd_g_v128_store(m, l0, int32(4), v3)
	*(*int32)(unsafe.Add(mBase, uint32(l0+int32(20)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = l0 + int32(8)
	v13 = int32(256)
	if v13 < l1 {
		v16 = l1
	} else {
		v16 = v13
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v16
	return
}

var F_VP8LBackwardRefsInit__k0 = [2]uint64{0x0, 0x0}

func F_VP8LBitWriterInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 base.V128
	_ = v4
	var v11 int64
	_ = v11
	var v15 int32
	_ = v15
	var v22 int64
	_ = v22
	var v23 int32
	_ = v23
	var v24 int64
	_ = v24
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v42 int32
	_ = v42
	v4 = base.Simd_g_const(&F_VP8LBitWriterInit__k0)
	base.Simd_g_v128_store(m, l0, int32(0), v4)
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(16)))) = int64(0)
	v11 = int64(1)
	v15 = l1&int32(-1024) + int32(1024)
	v22 = base.I64_div_u_s(int64(2147418112), v11)
	v23 = int32(0)
	v24 = base.I64_extend_i32_u(v15)
	if base.Ui64(int64(4294967295)) < base.Ui64(v24*v11) {
		v36 = v23
	} else {
		if base.Ui64(v22) < base.Ui64(v24) {
			v36 = v23
		} else {
			v34 = F_malloc(m, base.I32_wrap_i64(v11)*v15)
			mBase = m.M
			v36 = v34
		}
	}
	if v36 != 0 {
		v42 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		F_free(m, v42)
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v36 + v15
		*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v36
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v36
		return int32(1)
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
		return int32(0)
	}
}

var F_VP8LBitWriterInit__k0 = [2]uint64{0x0, 0x0}

func F_VP8LBitWriterSwap(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v9 int64
	_ = v9
	var v10 int64
	_ = v10
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v15 base.V128
	_ = v15
	var v17 int32
	_ = v17
	var v19 base.V128
	_ = v19
	v5 = int32(0)
	v9 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
	v10 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = v10
	v12 = int32(8)
	v13 = l0 + v12
	v15 = base.Simd_g_v128_load(m, v13, v5)
	v17 = l1 + v12
	v19 = base.Simd_g_v128_load(m, v17, v5)
	base.Simd_g_v128_store(m, v13, v5, v19)
	*(*int64)(unsafe.Add(mBase, uint32(l1))) = v9
	base.Simd_g_v128_store(m, v17, v5, v15)
	return
}
func F_VP8LBitWriterWipeOut(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v12 base.V128
	_ = v12
	if l0 == int32(0) {
	} else {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(8))))
		F_free(m, v6)
		mBase = m.M
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(16)))) = int64(0)
		v12 = base.Simd_g_const(&F_VP8LBitWriterWipeOut__k0)
		base.Simd_g_v128_store(m, l0, int32(0), v12)
	}
	return
}

var F_VP8LBitWriterWipeOut__k0 = [2]uint64{0x0, 0x0}

func F_VP8LCreateHuffmanTree(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v50 int32
	_ = v50
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v126 int64
	_ = v126
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v142 int32
	_ = v142
	var v164 int32
	_ = v164
	var v173 int32
	_ = v173
	var v180 int32
	_ = v180
	var v204 int32
	_ = v204
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v210 int32
	_ = v210
	var v214 int32
	_ = v214
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	var v245 int32
	_ = v245
	var v246 int32
	_ = v246
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v267 int32
	_ = v267
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v296 int32
	_ = v296
	var v301 int32
	_ = v301
	var v302 int32
	_ = v302
	var v337 int32
	_ = v337
	var v338 int32
	_ = v338
	var v343 int64
	_ = v343
	var v346 int32
	_ = v346
	var v347 int32
	_ = v347
	var v359 int32
	_ = v359
	var v377 int32
	_ = v377
	var v381 int32
	_ = v381
	var v382 int32
	_ = v382
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v397 int32
	_ = v397
	var v403 int32
	_ = v403
	var v405 int32
	_ = v405
	var v406 int32
	_ = v406
	var v408 int32
	_ = v408
	var v409 int32
	_ = v409
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v430 int32
	_ = v430
	var v435 int32
	_ = v435
	var v439 int32
	_ = v439
	var v440 int32
	_ = v440
	var v442 int32
	_ = v442
	var v454 int32
	_ = v454
	var v457 int32
	_ = v457
	var v458 int32
	_ = v458
	var v460 int32
	_ = v460
	var v462 int32
	_ = v462
	var v465 base.V128
	_ = v465
	var v467 int32
	_ = v467
	var v470 int32
	_ = v470
	var v479 int32
	_ = v479
	var v506 int32
	_ = v506
	var v525 int32
	_ = v525
	var v547 int32
	_ = v547
	var v556 int32
	_ = v556
	var v582 int32
	_ = v582
	var v617 int32
	_ = v617
	var v619 int32
	_ = v619
	var v620 int32
	_ = v620
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v625 int32
	_ = v625
	var v636 int32
	_ = v636
	var v638 int32
	_ = v638
	var v642 int32
	_ = v642
	var v644 int32
	_ = v644
	var v646 int32
	_ = v646
	var v652 int32
	_ = v652
	var v655 int32
	_ = v655
	var v664 int32
	_ = v664
	var v666 int32
	_ = v666
	var v686 int32
	_ = v686
	var v730 int32
	_ = v730
	var v733 int32
	_ = v733
	var v738 int32
	_ = v738
	var v739 base.V128
	_ = v739
	var v742 int32
	_ = v742
	var v751 int32
	_ = v751
	var v758 base.V128
	_ = v758
	var v774 base.V128
	_ = v774
	var v776 base.V128
	_ = v776
	var v780 int32
	_ = v780
	var v783 base.V128
	_ = v783
	var v788 int32
	_ = v788
	var v802 int32
	_ = v802
	var v810 int32
	_ = v810
	var v829 int32
	_ = v829
	var v838 int32
	_ = v838
	var v847 int32
	_ = v847
	var v860 int32
	_ = v860
	var v863 int32
	_ = v863
	var v867 int32
	_ = v867
	var v888 int32
	_ = v888
	var v903 int32
	_ = v903
	var v906 int32
	_ = v906
	var v913 int32
	_ = v913
	var v915 int32
	_ = v915
	var v931 int32
	_ = v931
	var v954 int32
	_ = v954
	var v957 int32
	_ = v957
	var v958 int32
	_ = v958
	var v962 int32
	_ = v962
	var v971 int32
	_ = v971
	var v972 int32
	_ = v972
	var v993 int32
	_ = v993
	var v998 int32
	_ = v998
	var v1003 int32
	_ = v1003
	var v1008 int32
	_ = v1008
	var v1011 int32
	_ = v1011
	var v1016 int32
	_ = v1016
	var v1023 int32
	_ = v1023
	var v1028 int32
	_ = v1028
	var v1032 int32
	_ = v1032
	var v1045 int32
	_ = v1045
	var v1046 int32
	_ = v1046
	var v1072 int32
	_ = v1072
	var v1077 int32
	_ = v1077
	var v1082 int32
	_ = v1082
	var v1087 int32
	_ = v1087
	var v1094 int32
	_ = v1094
	var v1096 int32
	_ = v1096
	var v1098 int32
	_ = v1098
	var v1105 int32
	_ = v1105
	var v1114 int32
	_ = v1114
	var v1121 base.V128
	_ = v1121
	var v1137 base.V128
	_ = v1137
	var v1140 base.V128
	_ = v1140
	var v1144 int32
	_ = v1144
	var v1147 base.V128
	_ = v1147
	var v1152 int32
	_ = v1152
	var v1156 int32
	_ = v1156
	var v1166 int32
	_ = v1166
	var v1191 int32
	_ = v1191
	var v1200 int32
	_ = v1200
	var v1201 int32
	_ = v1201
	var v1222 int32
	_ = v1222
	var v1224 int32
	_ = v1224
	var v1228 int32
	_ = v1228
	var v1231 int32
	_ = v1231
	var v1269 int32
	_ = v1269
	var v1270 int32
	_ = v1270
	var v1279 int32
	_ = v1279
	var v1281 int32
	_ = v1281
	var v1292 int32
	_ = v1292
	var v1318 int32
	_ = v1318
	var v1321 int32
	_ = v1321
	var v1322 int32
	_ = v1322
	var v1326 int32
	_ = v1326
	var v1335 int32
	_ = v1335
	var v1336 int32
	_ = v1336
	var v1357 int32
	_ = v1357
	var v1362 int32
	_ = v1362
	var v1367 int32
	_ = v1367
	var v1371 int32
	_ = v1371
	var v1375 int32
	_ = v1375
	var v1380 int32
	_ = v1380
	var v1387 int32
	_ = v1387
	var v1391 int32
	_ = v1391
	var v1396 int32
	_ = v1396
	var v1409 int32
	_ = v1409
	var v1410 int32
	_ = v1410
	var v1436 int32
	_ = v1436
	var v1441 int32
	_ = v1441
	var v1446 int32
	_ = v1446
	var v1451 int32
	_ = v1451
	var v1459 int32
	_ = v1459
	var v1469 int32
	_ = v1469
	var v1474 int32
	_ = v1474
	var v1490 int32
	_ = v1490
	var v1492 int32
	_ = v1492
	var v1495 int32
	_ = v1495
	var v1496 int64
	_ = v1496
	var v1498 int32
	_ = v1498
	var v1502 int64
	_ = v1502
	var v1505 int32
	_ = v1505
	var v1508 int32
	_ = v1508
	var v1510 int32
	_ = v1510
	var v1513 int32
	_ = v1513
	var v1514 int64
	_ = v1514
	var v1520 int64
	_ = v1520
	var v1524 int32
	_ = v1524
	var v1531 int32
	_ = v1531
	var v1540 int32
	_ = v1540
	var v1562 int32
	_ = v1562
	var v1567 int32
	_ = v1567
	var v1580 int32
	_ = v1580
	var v1604 int32
	_ = v1604
	var v1606 int32
	_ = v1606
	var v1608 int32
	_ = v1608
	var v1611 int32
	_ = v1611
	var v1619 int32
	_ = v1619
	var v1630 int32
	_ = v1630
	var v1638 int32
	_ = v1638
	var v1641 int32
	_ = v1641
	var v1643 int32
	_ = v1643
	var v1650 int32
	_ = v1650
	var v1653 int32
	_ = v1653
	var v1655 int32
	_ = v1655
	var v1662 int32
	_ = v1662
	var v1665 int32
	_ = v1665
	var v1667 int32
	_ = v1667
	var v1674 int32
	_ = v1674
	var v1689 int32
	_ = v1689
	var v1690 int32
	_ = v1690
	var v1692 int32
	_ = v1692
	var v1699 int32
	_ = v1699
	var v1700 int32
	_ = v1700
	var v1702 int32
	_ = v1702
	var v1709 int32
	_ = v1709
	var v1710 int32
	_ = v1710
	var v1712 int32
	_ = v1712
	var v1719 int32
	_ = v1719
	var v1722 int32
	_ = v1722
	var v1724 int32
	_ = v1724
	var v1730 int32
	_ = v1730
	var v1736 int32
	_ = v1736
	var v1745 int32
	_ = v1745
	var v1746 int32
	_ = v1746
	var v1752 int32
	_ = v1752
	var v1755 int32
	_ = v1755
	var v1757 int32
	_ = v1757
	var v1760 int32
	_ = v1760
	var v1773 int32
	_ = v1773
	var v1778 int32
	_ = v1778
	var v1779 int32
	_ = v1779
	var v1781 int32
	_ = v1781
	var v1784 int32
	_ = v1784
	var v1786 int32
	_ = v1786
	var v1790 int32
	_ = v1790
	var v1792 int32
	_ = v1792
	var v1796 int32
	_ = v1796
	var v1798 int32
	_ = v1798
	var v1801 int32
	_ = v1801
	var v1806 int32
	_ = v1806
	var v1814 int32
	_ = v1814
	var v1824 int32
	_ = v1824
	var v1825 int32
	_ = v1825
	var v1830 int32
	_ = v1830
	var v1832 int32
	_ = v1832
	var v1833 int32
	_ = v1833
	var v1835 int32
	_ = v1835
	var v1839 int32
	_ = v1839
	var v1852 int32
	_ = v1852
	var v1856 int32
	_ = v1856
	var v1857 int32
	_ = v1857
	var v1859 int32
	_ = v1859
	var v1862 int32
	_ = v1862
	var v1864 int32
	_ = v1864
	var v1868 int32
	_ = v1868
	var v1870 int32
	_ = v1870
	var v1874 int32
	_ = v1874
	var v1876 int32
	_ = v1876
	var v1879 int32
	_ = v1879
	var v1880 int32
	_ = v1880
	var v1881 int32
	_ = v1881
	var v1882 int32
	_ = v1882
	var v1886 int32
	_ = v1886
	var v1892 int32
	_ = v1892
	var v1899 int32
	_ = v1899
	var v1900 int32
	_ = v1900
	var v1901 int32
	_ = v1901
	var v1905 int32
	_ = v1905
	var v1907 int32
	_ = v1907
	var v1908 int32
	_ = v1908
	var v1910 int32
	_ = v1910
	var v1912 int32
	_ = v1912
	var v1914 int32
	_ = v1914
	var v1916 int32
	_ = v1916
	var v1918 int32
	_ = v1918
	var v1923 int32
	_ = v1923
	var v1925 int32
	_ = v1925
	var v1927 int32
	_ = v1927
	var v1929 int32
	_ = v1929
	var v1931 int32
	_ = v1931
	var v1933 int32
	_ = v1933
	var v1935 int32
	_ = v1935
	var v1937 int32
	_ = v1937
	var v1939 int32
	_ = v1939
	var v1941 int32
	_ = v1941
	var v1943 int32
	_ = v1943
	var v1945 int32
	_ = v1945
	var v1946 int32
	_ = v1946
	var v1948 int32
	_ = v1948
	var v1950 int32
	_ = v1950
	var v1954 int32
	_ = v1954
	var v1956 int32
	_ = v1956
	var v1958 int32
	_ = v1958
	var v1963 int32
	_ = v1963
	var v1967 int32
	_ = v1967
	var v1968 int32
	_ = v1968
	var v1969 int32
	_ = v1969
	var v1973 int32
	_ = v1973
	var v1975 int32
	_ = v1975
	var v1976 int32
	_ = v1976
	var v1978 int32
	_ = v1978
	var v1980 int32
	_ = v1980
	var v1982 int32
	_ = v1982
	var v1984 int32
	_ = v1984
	var v1985 int32
	_ = v1985
	var v1991 int32
	_ = v1991
	var v1993 int32
	_ = v1993
	var v1994 int32
	_ = v1994
	var v1997 int32
	_ = v1997
	var v1999 int32
	_ = v1999
	var v2001 int32
	_ = v2001
	var v2003 int32
	_ = v2003
	var v2005 int32
	_ = v2005
	var v2007 int32
	_ = v2007
	var v2009 int32
	_ = v2009
	var v2011 int32
	_ = v2011
	var v2013 int32
	_ = v2013
	var v2018 int32
	_ = v2018
	var v2028 int32
	_ = v2028
	var v2037 int32
	_ = v2037
	var v2039 int32
	_ = v2039
	var v2042 int32
	_ = v2042
	var v2045 int32
	_ = v2045
	var v2046 int32
	_ = v2046
	var v2047 int32
	_ = v2047
	var v2051 int32
	_ = v2051
	var v2053 int32
	_ = v2053
	var v2056 int32
	_ = v2056
	var v2057 int32
	_ = v2057
	var v2060 int32
	_ = v2060
	var v2063 int32
	_ = v2063
	var v2065 int32
	_ = v2065
	var v2068 int32
	_ = v2068
	var v2075 int32
	_ = v2075
	var v2084 int32
	_ = v2084
	var v2091 base.V128
	_ = v2091
	var v2107 base.V128
	_ = v2107
	var v2110 base.V128
	_ = v2110
	var v2114 int32
	_ = v2114
	var v2117 base.V128
	_ = v2117
	var v2122 int32
	_ = v2122
	var v2126 int32
	_ = v2126
	var v2136 int32
	_ = v2136
	var v2161 int32
	_ = v2161
	var v2170 int32
	_ = v2170
	var v2171 int32
	_ = v2171
	var v2192 int32
	_ = v2192
	var v2194 int32
	_ = v2194
	var v2198 int32
	_ = v2198
	var v2201 int32
	_ = v2201
	var v2254 int32
	_ = v2254
	var v2270 base.V128
	_ = v2270
	var v2271 int32
	_ = v2271
	var v2283 int32
	_ = v2283
	var v2285 int32
	_ = v2285
	var v2288 int32
	_ = v2288
	var v2303 int32
	_ = v2303
	var v2304 int32
	_ = v2304
	var v2322 int32
	_ = v2322
	var v2344 int32
	_ = v2344
	var v2345 int32
	_ = v2345
	var v2346 int32
	_ = v2346
	var v2348 int32
	_ = v2348
	var v2349 int32
	_ = v2349
	var v2350 int32
	_ = v2350
	var v2355 int32
	_ = v2355
	var v2358 int32
	_ = v2358
	var v2359 int32
	_ = v2359
	var v2365 int32
	_ = v2365
	var v2368 int32
	_ = v2368
	var v2369 int32
	_ = v2369
	var v2375 int32
	_ = v2375
	var v2378 int32
	_ = v2378
	var v2379 int32
	_ = v2379
	var v2384 int32
	_ = v2384
	var v2397 int32
	_ = v2397
	var v2424 int32
	_ = v2424
	var v2434 int32
	_ = v2434
	var v2455 int32
	_ = v2455
	var v2458 int32
	_ = v2458
	var v2459 int32
	_ = v2459
	var v2460 int32
	_ = v2460
	var v2466 int32
	_ = v2466
	var v2500 int32
	_ = v2500
	var v2503 int32
	_ = v2503
	var v2504 int32
	_ = v2504
	var v2505 int32
	_ = v2505
	var v2506 int32
	_ = v2506
	var v2507 int32
	_ = v2507
	var v2508 int32
	_ = v2508
	var v2509 int32
	_ = v2509
	var v2510 int32
	_ = v2510
	var v2511 int32
	_ = v2511
	var v2512 int32
	_ = v2512
	var v2513 int32
	_ = v2513
	var v2514 int32
	_ = v2514
	var v2515 int32
	_ = v2515
	var v2516 int32
	_ = v2516
	var v2518 int32
	_ = v2518
	var v2523 int32
	_ = v2523
	var v2524 int32
	_ = v2524
	var v2525 int32
	_ = v2525
	var v2526 int32
	_ = v2526
	var v2527 int32
	_ = v2527
	var v2528 int32
	_ = v2528
	var v2529 int32
	_ = v2529
	var v2530 int32
	_ = v2530
	var v2531 int32
	_ = v2531
	var v2532 int32
	_ = v2532
	var v2536 int32
	_ = v2536
	var v2539 int32
	_ = v2539
	var v2553 int32
	_ = v2553
	var v2554 int32
	_ = v2554
	var v2558 int32
	_ = v2558
	var v2562 int32
	_ = v2562
	var v2566 int32
	_ = v2566
	var v2570 int32
	_ = v2570
	var v2574 int32
	_ = v2574
	var v2578 int32
	_ = v2578
	var v2582 int32
	_ = v2582
	var v2586 int32
	_ = v2586
	var v2590 int32
	_ = v2590
	var v2594 int32
	_ = v2594
	var v2598 int32
	_ = v2598
	var v2604 int32
	_ = v2604
	var v2605 int32
	_ = v2605
	var v2610 int32
	_ = v2610
	var v2643 int32
	_ = v2643
	var v2646 int32
	_ = v2646
	var v2647 int32
	_ = v2647
	var v2651 int32
	_ = v2651
	var v2655 int32
	_ = v2655
	var v2666 int32
	_ = v2666
	var v2670 int32
	_ = v2670
	var v2678 int32
	_ = v2678
	var v2679 int32
	_ = v2679
	var v2680 int32
	_ = v2680
	var v2701 int32
	_ = v2701
	var v2703 int32
	_ = v2703
	var v2706 int32
	_ = v2706
	var v2709 int32
	_ = v2709
	var v2710 int32
	_ = v2710
	var v2716 int32
	_ = v2716
	var v2721 int32
	_ = v2721
	var v2723 int32
	_ = v2723
	var v2725 int32
	_ = v2725
	var v2727 int32
	_ = v2727
	var v2732 int32
	_ = v2732
	var v2741 int32
	_ = v2741
	var v2742 int32
	_ = v2742
	var v2765 int32
	_ = v2765
	var v2771 int32
	_ = v2771
	var v2788 int32
	_ = v2788
	var v2809 int32
	_ = v2809
	var v2814 int32
	_ = v2814
	var v2817 int32
	_ = v2817
	v6 = int32(0)
	v34 = m.G0
	v36 = v34 - int32(128)
	m.G0 = v36
	v39 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	if base.Ui32(v39) < base.Ui32(int32(33)) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	if v39 < int32(0) {
		goto L17
	} else {
		goto L18
	}
L2:
	;
	if v39 == int32(0) {
		goto L4
	} else {
		goto L5
	}
L3:
	;
	base.MemoryFill(m, l2, v6, v39)
	goto L1
L4:
	;
	goto L1
L5:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l2))) = uint8(v6)
	v50 = l2 + v39
	*(*uint8)(unsafe.Add(mBase, uint32(v50+int32(-1)))) = uint8(v6)
	if base.Ui32(v39) < base.Ui32(int32(3)) {
		goto L4
	} else {
		goto L6
	}
L6:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+2)) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+1)) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(v50+int32(-3)))) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(v50+int32(-2)))) = uint8(v6)
	if base.Ui32(v39) < base.Ui32(int32(7)) {
		goto L4
	} else {
		goto L7
	}
L7:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+3)) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(v50+int32(-4)))) = uint8(v6)
	if base.Ui32(v39) < base.Ui32(int32(9)) {
		goto L4
	} else {
		goto L8
	}
L8:
	;
	v72 = int32(0)
	v75 = (v72 - l2) & int32(3)
	v76 = l2 + v75
	*(*int32)(unsafe.Add(mBase, uint32(v76))) = v72
	v84 = (v39 - v75) & int32(60)
	v85 = v76 + v84
	*(*int32)(unsafe.Add(mBase, uint32(v85+int32(-4)))) = v72
	if base.Ui32(v84) < base.Ui32(int32(9)) {
		goto L4
	} else {
		goto L9
	}
L9:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v76)+8)) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v76)+4)) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v85+int32(-8)))) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v85+int32(-12)))) = v72
	if base.Ui32(v84) < base.Ui32(int32(25)) {
		goto L4
	} else {
		goto L10
	}
L10:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v76)+24)) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v76)+20)) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v76)+16)) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v76)+12)) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v85+int32(-16)))) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v85+int32(-20)))) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v85+int32(-24)))) = v72
	*(*int32)(unsafe.Add(mBase, uint32(v85+int32(-28)))) = v72
	v120 = v76&int32(4) | int32(24)
	v121 = v84 - v120
	if base.Ui32(v121) < base.Ui32(int32(32)) {
		goto L4
	} else {
		goto L11
	}
L11:
	;
	v126 = base.I64_extend_i32_u(v72) * int64(4294967297)
	v129 = v121
	v130 = v76 + v120
	goto L12
L12:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v130)+24)) = v126
	*(*int64)(unsafe.Add(mBase, uint32(v130)+16)) = v126
	*(*int64)(unsafe.Add(mBase, uint32(v130)+8)) = v126
	*(*int64)(unsafe.Add(mBase, uint32(v130))) = v126
	v142 = v129 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v142) {
		v129 = v142
		v130 = v130 + int32(32)
		goto L12
	} else {
		goto L14
	}
L13:
	;
	goto L4
L14:
	;
	goto L13
L15:
	;
	v2270 = base.Simd_g_const(&F_VP8LCreateHuffmanTree__k0)
	v2271 = int32(0)
	base.Simd_g_v128_store(m, v36+int32(48), v2271, v2270)
	base.Simd_g_v128_store(m, v36+int32(32), v2271, v2270)
	base.Simd_g_v128_store(m, v36+int32(16), v2271, v2270)
	base.Simd_g_v128_store(m, v36, v2271, v2270)
	v2283 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v2285 = base.B2i32(v2283 < int32(1))
	if v2285 == v2271 {
		goto L283
	} else {
		goto L284
	}
L16:
	;
	v730 = l4 + int32(4)
	if v39 == int32(0) {
		v2254 = v730
		goto L15
	} else {
		goto L93
	}
L17:
	;
	v2254 = l4 + int32(4)
	goto L15
L18:
	;
	if v39 == int32(0) {
		goto L19
	} else {
		goto L20
	}
L19:
	;
	v2254 = l4 + int32(4)
	goto L15
L20:
	;
	v164 = int32(-4)
	v173 = v39<<(uint(int32(2))%32) + l0 + v164
	v180 = v39
	goto L22
L21:
	;
	v209 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v210 = int32(0)
	v214 = v210
	v222 = v209
	v223 = v210
	v224 = l0
	goto L27
L22:
	;
	v204 = *(*int32)(unsafe.Add(mBase, uint32(v173)))
	if v204 != 0 {
		goto L21
	} else {
		goto L24
	}
L24:
	;
	v208 = v180 + int32(-1)
	if v208 != 0 {
		v173 = v173 + int32(-4)
		v180 = v208
		goto L22
	} else {
		goto L25
	}
L25:
	;
	goto L16
L26:
	;
	v391 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v392 = int32(0)
	v397 = v392
	v403 = l0 + v164
	v405 = v392
	v406 = v391
	v408 = v392
	v409 = l0 + int32(-16)
	goto L56
L27:
	;
	v245 = base.B2i32(v180 == v214)
	if v180 == v214 {
		goto L31
	} else {
		goto L32
	}
L29:
	;
	v214 = v214 + int32(1)
	v222 = v381
	v223 = v382
	v224 = v224 + int32(4)
	goto L27
L30:
	;
	v381 = v222
	v382 = v223 + int32(1)
	goto L29
L31:
	;
	if v222 != 0 {
		goto L36
	} else {
		goto L37
	}
L32:
	;
	v246 = *(*int32)(unsafe.Add(mBase, uint32(v224)))
	if v246 == v222 {
		goto L30
	} else {
		goto L33
	}
L33:
	;
	goto L31
L34:
	;
	if v180 == v214 {
		goto L26
	} else {
		goto L55
	}
L35:
	;
	v255 = l2 + (v214 - v223)
	v256 = int32(1)
	if base.Ui32(v223) < base.Ui32(int32(33)) {
		goto L42
	} else {
		goto L43
	}
L36:
	;
	if v222 == int32(0) {
		goto L34
	} else {
		goto L39
	}
L37:
	;
	if int32(4) < v223 {
		goto L35
	} else {
		goto L38
	}
L38:
	;
	goto L36
L39:
	;
	if v223 < int32(7) {
		goto L34
	} else {
		goto L40
	}
L40:
	;
	goto L35
L41:
	;
	goto L34
L42:
	;
	if v223 == int32(0) {
		goto L44
	} else {
		goto L45
	}
L43:
	;
	base.MemoryFill(m, v255, v256, v223)
	goto L41
L44:
	;
	goto L41
L45:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v255))) = uint8(v256)
	v267 = v255 + v223
	*(*uint8)(unsafe.Add(mBase, uint32(v267+int32(-1)))) = uint8(v256)
	if base.Ui32(v223) < base.Ui32(int32(3)) {
		goto L44
	} else {
		goto L46
	}
L46:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v255)+2)) = uint8(v256)
	*(*uint8)(unsafe.Add(mBase, uint32(v255)+1)) = uint8(v256)
	*(*uint8)(unsafe.Add(mBase, uint32(v267+int32(-3)))) = uint8(v256)
	*(*uint8)(unsafe.Add(mBase, uint32(v267+int32(-2)))) = uint8(v256)
	if base.Ui32(v223) < base.Ui32(int32(7)) {
		goto L44
	} else {
		goto L47
	}
L47:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v255)+3)) = uint8(v256)
	*(*uint8)(unsafe.Add(mBase, uint32(v267+int32(-4)))) = uint8(v256)
	if base.Ui32(v223) < base.Ui32(int32(9)) {
		goto L44
	} else {
		goto L48
	}
L48:
	;
	v292 = (int32(0) - v255) & int32(3)
	v293 = v255 + v292
	v296 = int32(16843009)
	*(*int32)(unsafe.Add(mBase, uint32(v293))) = v296
	v301 = (v223 - v292) & int32(60)
	v302 = v293 + v301
	*(*int32)(unsafe.Add(mBase, uint32(v302+int32(-4)))) = v296
	if base.Ui32(v301) < base.Ui32(int32(9)) {
		goto L44
	} else {
		goto L49
	}
L49:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v293)+8)) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v293)+4)) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v302+int32(-8)))) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v302+int32(-12)))) = v296
	if base.Ui32(v301) < base.Ui32(int32(25)) {
		goto L44
	} else {
		goto L50
	}
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v293)+24)) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v293)+20)) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v293)+16)) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v293)+12)) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v302+int32(-16)))) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v302+int32(-20)))) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v302+int32(-24)))) = v296
	*(*int32)(unsafe.Add(mBase, uint32(v302+int32(-28)))) = v296
	v337 = v293&int32(4) | int32(24)
	v338 = v301 - v337
	if base.Ui32(v338) < base.Ui32(int32(32)) {
		goto L44
	} else {
		goto L51
	}
L51:
	;
	v343 = base.I64_extend_i32_u(v296) * int64(4294967297)
	v346 = v338
	v347 = v293 + v337
	goto L52
L52:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v347)+24)) = v343
	*(*int64)(unsafe.Add(mBase, uint32(v347)+16)) = v343
	*(*int64)(unsafe.Add(mBase, uint32(v347)+8)) = v343
	*(*int64)(unsafe.Add(mBase, uint32(v347))) = v343
	v359 = v346 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v359) {
		v346 = v359
		v347 = v347 + int32(32)
		goto L52
	} else {
		goto L54
	}
L53:
	;
	goto L44
L54:
	;
	goto L53
L55:
	;
	v377 = *(*int32)(unsafe.Add(mBase, uint32(v224)))
	v381 = v377
	v382 = int32(1)
	goto L29
L56:
	;
	v428 = base.B2i32(v180 == v405)
	if v180 == v405 {
		goto L60
	} else {
		goto L61
	}
L58:
	;
	v686 = int32(4)
	v397 = v655
	v403 = v403 + v686
	v405 = v405 + int32(1)
	v406 = v664
	v408 = v666
	v409 = v409 + v686
	goto L56
L59:
	;
	v644 = v439 + v397
	v646 = v408 + int32(1)
	if base.Ui32(v646) < base.Ui32(int32(4)) {
		v655 = v644
		v664 = v406
		v666 = v646
		goto L58
	} else {
		goto L92
	}
L60:
	;
	if base.Ui32(int32(3)) < base.Ui32(v408) {
		goto L68
	} else {
		goto L69
	}
L61:
	;
	v429 = l2 + v405
	v430 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v429))))
	if v430 != 0 {
		goto L60
	} else {
		goto L62
	}
L62:
	;
	if v405 == int32(0) {
		goto L63
	} else {
		goto L64
	}
L63:
	;
	v439 = *(*int32)(unsafe.Add(mBase, uint32(l0+v405<<(uint(int32(2))%32))))
	v440 = v439 - v406
	v442 = v440 >> (uint(int32(31)) % 32)
	if base.Ui32(v440^v442-v442) < base.Ui32(int32(4)) {
		goto L59
	} else {
		goto L66
	}
L64:
	;
	v435 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v429+int32(-1)))))
	if v435 != 0 {
		goto L60
	} else {
		goto L65
	}
L65:
	;
	goto L63
L66:
	;
	goto L60
L67:
	;
	if v180+int32(-3) <= v405 {
		goto L88
	} else {
		goto L89
	}
L68:
	;
	v454 = int32(1)
	v457 = int32(base.Ui32(v408)>>(uint(v454)%32)) + v397
	v458 = base.I32_div_u_s(v457, v408)
	if base.Ui32(v457) < base.Ui32(v408) {
		goto L72
	} else {
		goto L73
	}
L69:
	;
	if v408 != int32(3) {
		goto L67
	} else {
		goto L70
	}
L70:
	;
	if v397 != 0 {
		goto L67
	} else {
		goto L71
	}
L71:
	;
	goto L68
L72:
	;
	v460 = v454
	goto L74
L73:
	;
	v460 = v458
	goto L74
L74:
	;
	if v397 != 0 {
		goto L75
	} else {
		goto L76
	}
L75:
	;
	v462 = v460
	goto L77
L76:
	;
	v462 = int32(0)
	goto L77
L77:
	;
	if base.Ui32(v408) < base.Ui32(int32(4)) {
		v525 = int32(0)
		goto L78
	} else {
		goto L79
	}
L78:
	;
	v547 = v403 - v525<<(uint(int32(2))%32)
	v556 = v408 - v525
	goto L84
L79:
	;
	v465 = base.Simd_g_i32x4_splat(v462)
	v467 = v408 & int32(-4)
	v470 = v409
	v479 = v467
	goto L80
L80:
	;
	base.Simd_g_v128_store(m, v470, int32(0), v465)
	v506 = v479 + int32(-4)
	if v506 != 0 {
		v470 = v470 + int32(-16)
		v479 = v506
		goto L80
	} else {
		goto L82
	}
L81:
	;
	if v408 == v467 {
		goto L67
	} else {
		goto L83
	}
L82:
	;
	goto L81
L83:
	;
	v525 = v467
	goto L78
L84:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v547))) = v462
	v582 = v556 + int32(-1)
	if v582 != 0 {
		v547 = v547 + int32(-4)
		v556 = v582
		goto L84
	} else {
		goto L86
	}
L85:
	;
	goto L67
L86:
	;
	goto L85
L87:
	;
	if v180 == v405 {
		goto L16
	} else {
		goto L91
	}
L88:
	;
	if v180 <= v405 {
		v638 = int32(0)
		goto L87
	} else {
		goto L90
	}
L89:
	;
	v617 = int32(2)
	v619 = l0 + v405<<(uint(v617)%32)
	v620 = *(*int32)(unsafe.Add(mBase, uint32(v619)))
	v621 = *(*int32)(unsafe.Add(mBase, uint32(v619)+4))
	v623 = *(*int32)(unsafe.Add(mBase, uint32(v619)+8))
	v625 = *(*int32)(unsafe.Add(mBase, uint32(v619)+12))
	v638 = int32(base.Ui32(v620+v621+v623+v625+v617) >> (uint(v617) % 32))
	goto L87
L90:
	;
	v636 = *(*int32)(unsafe.Add(mBase, uint32(l0+v405<<(uint(int32(2))%32))))
	v638 = v636
	goto L87
L91:
	;
	v642 = *(*int32)(unsafe.Add(mBase, uint32(l0+v405<<(uint(int32(2))%32))))
	v655 = v642
	v664 = v638
	v666 = int32(1)
	goto L58
L92:
	;
	v652 = base.I32_div_u_s(v644+int32(base.Ui32(v646)>>(uint(int32(1))%32)), v646)
	v655 = v644
	v664 = v652
	v666 = v646
	goto L58
L93:
	;
	v733 = int32(0)
	if base.Ui32(v39) < base.Ui32(int32(4)) {
		v802 = v733
		v810 = v733
		goto L95
	} else {
		goto L96
	}
L94:
	;
	if v888 == int32(0) {
		v2254 = v730
		goto L15
	} else {
		goto L104
	}
L95:
	;
	v829 = l0 + v802<<(uint(int32(2))%32)
	v838 = v39 - v802
	v847 = v810
	goto L101
L96:
	;
	v738 = v39 & int32(2147483644)
	v739 = base.Simd_g_const(&F_VP8LCreateHuffmanTree__k0)
	v742 = l0
	v751 = v738
	v758 = v739
	goto L97
L97:
	;
	v774 = base.Simd_g_v128_load(m, v742, int32(0))
	v776 = base.Simd_g_i32x4_sub(v758, base.Simd_g_i32x4_ne(v774, v739))
	v780 = v751 + int32(-4)
	if v780 != 0 {
		v742 = v742 + int32(16)
		v751 = v780
		v758 = v776
		goto L97
	} else {
		goto L99
	}
L98:
	;
	v783 = base.Simd_g_i32x4_add(v776, base.Simd_g_i8x16_swizzle_c(v776, base.Simd_g_const(&F_VP8LCreateHuffmanTree__k1)))
	v788 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v783, base.Simd_g_i8x16_swizzle_c(v783, base.Simd_g_const(&F_VP8LCreateHuffmanTree__k2))))
	if v39 == v738 {
		v888 = v788
		goto L94
	} else {
		goto L100
	}
L99:
	;
	goto L98
L100:
	;
	v802 = v738
	v810 = v788
	goto L95
L101:
	;
	v860 = *(*int32)(unsafe.Add(mBase, uint32(v829)))
	v863 = v847 + base.B2i32(v860 != int32(0))
	v867 = v838 + int32(-1)
	if v867 != 0 {
		v829 = v829 + int32(4)
		v838 = v867
		v847 = v863
		goto L101
	} else {
		goto L103
	}
L102:
	;
	v888 = v863
	goto L94
L103:
	;
	goto L102
L104:
	;
	v903 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
	if int32(1) < v888 {
		goto L105
	} else {
		goto L106
	}
L105:
	;
	v1269 = l3 + v888<<(uint(int32(4))%32)
	v1270 = int32(1)
	v1279 = v39 + int32(-1)
	v1281 = v1279 & int32(-4)
	v1292 = v1270
	goto L148
L106:
	;
	v906 = int32(1)
	v913 = v39 + int32(-1)
	v915 = v913 & int32(-4)
	v931 = v906
	goto L107
L107:
	;
	v954 = int32(0)
	v957 = base.B2i32(v39 == int32(1))
	if v39 == int32(1) {
		v1045 = v954
		v1046 = v954
		goto L109
	} else {
		goto L110
	}
L109:
	;
	if v39&v906 == int32(0) {
		goto L124
	} else {
		goto L125
	}
L110:
	;
	v958 = int32(0)
	v962 = l0
	v971 = v958
	v972 = v958
	goto L111
L111:
	;
	v993 = *(*int32)(unsafe.Add(mBase, uint32(v962)))
	if v993 == int32(0) {
		v1008 = v971
		goto L113
	} else {
		goto L114
	}
L112:
	;
	v1045 = v1028
	v1046 = v1032
	goto L109
L113:
	;
	v1011 = *(*int32)(unsafe.Add(mBase, uint32(v962+int32(4))))
	if v1011 == int32(0) {
		v1028 = v1008
		goto L118
	} else {
		goto L119
	}
L114:
	;
	v998 = l3 + v971<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v998)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v998)+4)) = v972
	if base.Ui32(v931) < base.Ui32(v993) {
		goto L115
	} else {
		goto L116
	}
L115:
	;
	v1003 = v993
	goto L117
L116:
	;
	v1003 = v931
	goto L117
L117:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v998))) = v1003
	v1008 = v971 + int32(1)
	goto L113
L118:
	;
	v1032 = v972 + int32(2)
	if v39&int32(2147483646) != v1032 {
		v962 = v962 + int32(8)
		v971 = v1028
		v972 = v1032
		goto L111
	} else {
		goto L123
	}
L119:
	;
	v1016 = l3 + v1008<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1016)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1016)+4)) = v972 + int32(1)
	if base.Ui32(v931) < base.Ui32(v1011) {
		goto L120
	} else {
		goto L121
	}
L120:
	;
	v1023 = v1011
	goto L122
L121:
	;
	v1023 = v931
	goto L122
L122:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1016))) = v1023
	v1028 = v1008 + int32(1)
	goto L118
L123:
	;
	goto L112
L124:
	;
	v1087 = m.G2
	F___qsort_r(m, l3, v888, int32(16), int32(344), v1087+int32(335))
	mBase = m.M
	goto L130
L125:
	;
	v1072 = *(*int32)(unsafe.Add(mBase, uint32(l0+v1046<<(uint(int32(2))%32))))
	if v1072 == int32(0) {
		goto L124
	} else {
		goto L126
	}
L126:
	;
	v1077 = l3 + v1045<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1077)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1077)+4)) = v1046
	if base.Ui32(v931) < base.Ui32(v1072) {
		goto L127
	} else {
		goto L128
	}
L127:
	;
	v1082 = v1072
	goto L129
L128:
	;
	v1082 = v931
	goto L129
L129:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1077))) = v1082
	goto L124
L130:
	;
	if v888 != int32(1) {
		goto L131
	} else {
		goto L132
	}
L131:
	;
	v1098 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v903))))
	if v39 == int32(1) {
		v1231 = v1098
		goto L133
	} else {
		goto L134
	}
L132:
	;
	v1094 = *(*int32)(unsafe.Add(mBase, uint32(l3)+4))
	v1096 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v903+v1094))) = uint8(v1096)
	goto L131
L133:
	;
	if l1 < v1231 {
		v931 = v931 << (uint(int32(1)) % 32)
		goto L107
	} else {
		goto L147
	}
L134:
	;
	if base.B2i32(base.Ui32(int32(4)) < base.Ui32(v39)) == int32(0) {
		v1156 = v1098
		v1166 = int32(1)
		goto L135
	} else {
		goto L136
	}
L135:
	;
	v1191 = v1156
	v1200 = v903 + v1166
	v1201 = v39 - v1166
	goto L141
L136:
	;
	v1105 = v903 + v906
	v1114 = v915
	v1121 = base.Simd_g_i32x4_splat(v1098)
	goto L137
L137:
	;
	v1137 = base.Simd_g_v128_load32_zero(m, v1105, int32(0))
	v1140 = base.Simd_g_i32x4_max_s(v1121, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v1137)))
	v1144 = v1114 + int32(-4)
	if v1144 != 0 {
		v1105 = v1105 + int32(4)
		v1114 = v1144
		v1121 = v1140
		goto L137
	} else {
		goto L139
	}
L138:
	;
	v1147 = base.Simd_g_i32x4_max_s(v1140, base.Simd_g_i8x16_swizzle_c(v1140, base.Simd_g_const(&F_VP8LCreateHuffmanTree__k1)))
	v1152 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_max_s(v1147, base.Simd_g_i8x16_swizzle_c(v1147, base.Simd_g_const(&F_VP8LCreateHuffmanTree__k2))))
	if v913 == v915 {
		v1231 = v1152
		goto L133
	} else {
		goto L140
	}
L139:
	;
	goto L138
L140:
	;
	v1156 = v1152
	v1166 = v915 | v906
	goto L135
L141:
	;
	v1222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1200))))
	if v1222 < v1191 {
		goto L143
	} else {
		goto L144
	}
L142:
	;
	v1231 = v1224
	goto L133
L143:
	;
	v1224 = v1191
	goto L145
L144:
	;
	v1224 = v1222
	goto L145
L145:
	;
	v1228 = v1201 + int32(-1)
	if v1228 != 0 {
		v1191 = v1224
		v1200 = v1200 + int32(1)
		v1201 = v1228
		goto L141
	} else {
		goto L146
	}
L146:
	;
	goto L142
L147:
	;
	v2254 = v730
	goto L15
L148:
	;
	v1318 = int32(0)
	v1321 = base.B2i32(v39 == int32(1))
	if v39 == int32(1) {
		v1409 = v1318
		v1410 = v1318
		goto L150
	} else {
		goto L151
	}
L149:
	;
	v2254 = v730
	goto L15
L150:
	;
	if v39&v1270 == int32(0) {
		goto L165
	} else {
		goto L166
	}
L151:
	;
	v1322 = int32(0)
	v1326 = l0
	v1335 = v1322
	v1336 = v1322
	goto L152
L152:
	;
	v1357 = *(*int32)(unsafe.Add(mBase, uint32(v1326)))
	if v1357 == int32(0) {
		v1371 = v1335
		goto L154
	} else {
		goto L155
	}
L153:
	;
	v1409 = v1391
	v1410 = v1396
	goto L150
L154:
	;
	v1375 = *(*int32)(unsafe.Add(mBase, uint32(v1326+int32(4))))
	if v1375 == int32(0) {
		v1391 = v1371
		goto L159
	} else {
		goto L160
	}
L155:
	;
	v1362 = l3 + v1335<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1362)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1362)+4)) = v1336
	if base.Ui32(v1292) < base.Ui32(v1357) {
		goto L156
	} else {
		goto L157
	}
L156:
	;
	v1367 = v1357
	goto L158
L157:
	;
	v1367 = v1292
	goto L158
L158:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1362))) = v1367
	v1371 = v1335 + int32(1)
	goto L154
L159:
	;
	v1396 = v1336 + int32(2)
	if v39&int32(2147483646) != v1396 {
		v1326 = v1326 + int32(8)
		v1335 = v1391
		v1336 = v1396
		goto L152
	} else {
		goto L164
	}
L160:
	;
	v1380 = l3 + v1371<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1380)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1380)+4)) = v1336 + int32(1)
	if base.Ui32(v1292) < base.Ui32(v1375) {
		goto L161
	} else {
		goto L162
	}
L161:
	;
	v1387 = v1375
	goto L163
L162:
	;
	v1387 = v1292
	goto L163
L163:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1380))) = v1387
	v1391 = v1371 + int32(1)
	goto L159
L164:
	;
	goto L153
L165:
	;
	v1451 = m.G2
	F___qsort_r(m, l3, v888, int32(16), int32(344), v1451+int32(335))
	mBase = m.M
	goto L171
L166:
	;
	v1436 = *(*int32)(unsafe.Add(mBase, uint32(l0+v1410<<(uint(int32(2))%32))))
	if v1436 == int32(0) {
		goto L165
	} else {
		goto L167
	}
L167:
	;
	v1441 = l3 + v1409<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1441)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1441)+4)) = v1410
	if base.Ui32(v1292) < base.Ui32(v1436) {
		goto L168
	} else {
		goto L169
	}
L168:
	;
	v1446 = v1436
	goto L170
L169:
	;
	v1446 = v1292
	goto L170
L170:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1441))) = v1446
	goto L165
L171:
	;
	v1459 = v888
	v1469 = v888 + int32(-2)
	v1474 = int32(0)
	goto L172
L172:
	;
	v1490 = int32(4)
	v1492 = v1269 + v1474<<(uint(v1490)%32)
	v1495 = l3 + int32(-16) + v1459<<(uint(v1490)%32)
	v1496 = *(*int64)(unsafe.Add(mBase, uint32(v1495)))
	*(*int64)(unsafe.Add(mBase, uint32(v1492))) = v1496
	v1498 = int32(8)
	v1502 = *(*int64)(unsafe.Add(mBase, uint32(v1495+v1498)))
	*(*int64)(unsafe.Add(mBase, uint32(v1492+v1498))) = v1502
	v1505 = v1474 | int32(1)
	v1508 = v1269 + v1505<<(uint(v1490)%32)
	v1510 = v1459 + int32(-2)
	v1513 = l3 + v1510<<(uint(v1490)%32)
	v1514 = *(*int64)(unsafe.Add(mBase, uint32(v1513)))
	*(*int64)(unsafe.Add(mBase, uint32(v1508))) = v1514
	v1520 = *(*int64)(unsafe.Add(mBase, uint32(v1513+v1498)))
	*(*int64)(unsafe.Add(mBase, uint32(v1508+v1498))) = v1520
	v1524 = base.I32_wrap_i64(v1496) + base.I32_wrap_i64(v1514)
	if v1459 < int32(3) {
		v1580 = int32(0)
		goto L174
	} else {
		goto L175
	}
L173:
	;
	v2037 = int32(0)
	v2039 = *(*int32)(unsafe.Add(mBase, uint32(l3)+8))
	if v2039 < v2037 {
		v2060 = l3
		v2063 = v2037
		goto L262
	} else {
		goto L263
	}
L174:
	;
	v1604 = int32(4)
	v1606 = l3 + v1580<<(uint(v1604)%32)
	v1608 = v1606 + int32(16)
	v1611 = (v1510 - v1580) << (uint(v1604) % 32)
	if base.Ui32(int32(33)) <= base.Ui32(v1611) {
		goto L183
	} else {
		goto L184
	}
L175:
	;
	v1531 = l3
	v1540 = int32(0)
	goto L176
L176:
	;
	v1562 = *(*int32)(unsafe.Add(mBase, uint32(v1531)))
	if base.Ui32(v1562) <= base.Ui32(v1524) {
		v1580 = v1540
		goto L174
	} else {
		goto L178
	}
L177:
	;
	v1580 = v1469
	goto L174
L178:
	;
	v1567 = v1540 + int32(1)
	if v1469 != v1567 {
		v1531 = v1531 + int32(16)
		v1540 = v1567
		goto L176
	} else {
		goto L179
	}
L179:
	;
	goto L177
L180:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1606)+12)) = v1474
	*(*int32)(unsafe.Add(mBase, uint32(v1606)+8)) = v1505
	v2028 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1606)+4)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1606))) = v1524
	if int32(2) < v1459 {
		v1459 = v1459 + v2028
		v1469 = v1469 + v2028
		v1474 = v1474 + int32(2)
		goto L172
	} else {
		goto L260
	}
L181:
	;
	goto L180
L182:
	;
	v1630 = (v1606 ^ v1608) & int32(3)
	if base.Ui32(v1606) <= base.Ui32(v1608) {
		goto L189
	} else {
		goto L190
	}
L183:
	;
	base.MemoryCopy(m, v1608, v1606, v1611)
	goto L181
L184:
	;
	if v1608 == v1606 {
		goto L181
	} else {
		goto L185
	}
L185:
	;
	v1619 = v1608 + v1611
	if base.Ui32(int32(0)-v1611<<(uint(int32(1))%32)) < base.Ui32(v1606-v1619) {
		goto L182
	} else {
		goto L186
	}
L186:
	;
	goto L183
L187:
	;
	if v1958 == int32(0) {
		goto L181
	} else {
		goto L249
	}
L188:
	;
	if base.Ui32(v1882) < base.Ui32(int32(4)) {
		v1954 = v1880
		v1956 = v1881
		v1958 = v1882
		goto L187
	} else {
		goto L239
	}
L189:
	;
	if v1630 != 0 {
		v1806 = v1611
		goto L205
	} else {
		goto L206
	}
L190:
	;
	if v1630 == int32(0) {
		goto L191
	} else {
		goto L192
	}
L191:
	;
	if v1608&int32(3) != 0 {
		goto L193
	} else {
		goto L194
	}
L192:
	;
	v1954 = v1606
	v1956 = v1608
	v1958 = v1611
	goto L187
L193:
	;
	if v1611 == int32(0) {
		goto L181
	} else {
		goto L195
	}
L194:
	;
	v1880 = v1606
	v1881 = v1608
	v1882 = v1611
	goto L188
L195:
	;
	v1638 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1608))) = uint8(v1638)
	v1641 = v1611 + int32(-1)
	v1643 = v1606 + int32(17)
	if v1643&int32(3) != 0 {
		goto L196
	} else {
		goto L197
	}
L196:
	;
	if v1641 == int32(0) {
		goto L181
	} else {
		goto L198
	}
L197:
	;
	v1880 = v1606 + int32(1)
	v1881 = v1643
	v1882 = v1641
	goto L188
L198:
	;
	v1650 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606)+1)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1608)+1)) = uint8(v1650)
	v1653 = v1611 + int32(-2)
	v1655 = v1606 + int32(18)
	if v1655&int32(3) != 0 {
		goto L199
	} else {
		goto L200
	}
L199:
	;
	if v1653 == int32(0) {
		goto L181
	} else {
		goto L201
	}
L200:
	;
	v1880 = v1606 + int32(2)
	v1881 = v1655
	v1882 = v1653
	goto L188
L201:
	;
	v1662 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606)+2)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1608)+2)) = uint8(v1662)
	v1665 = v1611 + int32(-3)
	v1667 = v1606 + int32(19)
	if v1667&int32(3) != 0 {
		goto L202
	} else {
		goto L203
	}
L202:
	;
	if v1665 == int32(0) {
		goto L181
	} else {
		goto L204
	}
L203:
	;
	v1880 = v1606 + int32(3)
	v1881 = v1667
	v1882 = v1665
	goto L188
L204:
	;
	v1674 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606)+3)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1608)+3)) = uint8(v1674)
	v1880 = v1606 + int32(4)
	v1881 = v1606 + int32(20)
	v1882 = v1611 + int32(-4)
	goto L188
L205:
	;
	if v1806 == int32(0) {
		goto L181
	} else {
		goto L229
	}
L206:
	;
	if v1619&int32(3) == int32(0) {
		v1724 = v1611
		goto L207
	} else {
		goto L208
	}
L207:
	;
	if base.Ui32(v1724) < base.Ui32(int32(4)) {
		v1806 = v1724
		goto L205
	} else {
		goto L219
	}
L208:
	;
	if v1611 == int32(0) {
		goto L181
	} else {
		goto L209
	}
L209:
	;
	v1689 = v1611 + int32(-1)
	v1690 = v1608 + v1689
	v1692 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606+v1689))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1690))) = uint8(v1692)
	if v1690&int32(3) != 0 {
		goto L210
	} else {
		goto L211
	}
L210:
	;
	if v1689 == int32(0) {
		goto L181
	} else {
		goto L212
	}
L211:
	;
	v1724 = v1689
	goto L207
L212:
	;
	v1699 = v1611 + int32(-2)
	v1700 = v1608 + v1699
	v1702 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606+v1699))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1700))) = uint8(v1702)
	if v1700&int32(3) != 0 {
		goto L213
	} else {
		goto L214
	}
L213:
	;
	if v1699 == int32(0) {
		goto L181
	} else {
		goto L215
	}
L214:
	;
	v1724 = v1699
	goto L207
L215:
	;
	v1709 = v1611 + int32(-3)
	v1710 = v1608 + v1709
	v1712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606+v1709))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1710))) = uint8(v1712)
	if v1710&int32(3) != 0 {
		goto L216
	} else {
		goto L217
	}
L216:
	;
	if v1709 == int32(0) {
		goto L181
	} else {
		goto L218
	}
L217:
	;
	v1724 = v1709
	goto L207
L218:
	;
	v1719 = v1611 + int32(-4)
	v1722 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606+v1719))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1608+v1719))) = uint8(v1722)
	v1724 = v1719
	goto L207
L219:
	;
	v1730 = v1724 + int32(-4)
	v1736 = (int32(base.Ui32(v1730)>>(uint(int32(2))%32)) + int32(1)) & int32(3)
	if v1736 == int32(0) {
		v1760 = v1724
		goto L220
	} else {
		goto L221
	}
L220:
	;
	if base.Ui32(v1730) < base.Ui32(int32(12)) {
		v1806 = v1760
		goto L205
	} else {
		goto L225
	}
L221:
	;
	v1745 = v1724
	v1746 = v1736
	goto L222
L222:
	;
	v1752 = *(*int32)(unsafe.Add(mBase, uint32(v1606+int32(-4)+v1745)))
	*(*int32)(unsafe.Add(mBase, uint32(v1606+int32(12)+v1745))) = v1752
	v1755 = v1745 + int32(-4)
	v1757 = v1746 + int32(-1)
	if v1757 != 0 {
		v1745 = v1755
		v1746 = v1757
		goto L222
	} else {
		goto L224
	}
L223:
	;
	v1760 = v1755
	goto L220
L224:
	;
	goto L223
L225:
	;
	v1773 = v1760
	goto L226
L226:
	;
	v1778 = v1606 + int32(0) + v1773
	v1779 = int32(12)
	v1781 = v1606 + int32(-16) + v1773
	v1784 = *(*int32)(unsafe.Add(mBase, uint32(v1781+v1779)))
	*(*int32)(unsafe.Add(mBase, uint32(v1778+v1779))) = v1784
	v1786 = int32(8)
	v1790 = *(*int32)(unsafe.Add(mBase, uint32(v1781+v1786)))
	*(*int32)(unsafe.Add(mBase, uint32(v1778+v1786))) = v1790
	v1792 = int32(4)
	v1796 = *(*int32)(unsafe.Add(mBase, uint32(v1781+v1792)))
	*(*int32)(unsafe.Add(mBase, uint32(v1778+v1792))) = v1796
	v1798 = *(*int32)(unsafe.Add(mBase, uint32(v1781)))
	*(*int32)(unsafe.Add(mBase, uint32(v1778))) = v1798
	v1801 = v1773 + int32(-16)
	if base.Ui32(int32(3)) < base.Ui32(v1801) {
		v1773 = v1801
		goto L226
	} else {
		goto L228
	}
L227:
	;
	v1806 = v1801
	goto L205
L228:
	;
	goto L227
L229:
	;
	v1814 = v1806 & int32(3)
	if v1814 == int32(0) {
		v1839 = v1806
		goto L230
	} else {
		goto L231
	}
L230:
	;
	if base.Ui32(v1806) < base.Ui32(int32(4)) {
		goto L181
	} else {
		goto L235
	}
L231:
	;
	v1824 = v1806
	v1825 = v1814
	goto L232
L232:
	;
	v1830 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1606+int32(-1)+v1824))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1606+int32(15)+v1824))) = uint8(v1830)
	v1832 = int32(-1)
	v1833 = v1824 + v1832
	v1835 = v1825 + v1832
	if v1835 != 0 {
		v1824 = v1833
		v1825 = v1835
		goto L232
	} else {
		goto L234
	}
L233:
	;
	v1839 = v1833
	goto L230
L234:
	;
	goto L233
L235:
	;
	v1852 = v1839
	goto L236
L236:
	;
	v1856 = v1606 + int32(12) + v1852
	v1857 = int32(3)
	v1859 = v1606 + int32(-4) + v1852
	v1862 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1859+v1857))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1856+v1857))) = uint8(v1862)
	v1864 = int32(2)
	v1868 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1859+v1864))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1856+v1864))) = uint8(v1868)
	v1870 = int32(1)
	v1874 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1859+v1870))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1856+v1870))) = uint8(v1874)
	v1876 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1859))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1856))) = uint8(v1876)
	v1879 = v1852 + int32(-4)
	if v1879 != 0 {
		v1852 = v1879
		goto L236
	} else {
		goto L238
	}
L238:
	;
	goto L181
L239:
	;
	v1886 = v1882 + int32(-4)
	v1892 = (int32(base.Ui32(v1886)>>(uint(int32(2))%32)) + int32(1)) & int32(7)
	if v1892 == int32(0) {
		v1914 = v1880
		v1916 = v1881
		v1918 = v1882
		goto L240
	} else {
		goto L241
	}
L240:
	;
	if base.Ui32(v1886) < base.Ui32(int32(28)) {
		v1954 = v1914
		v1956 = v1916
		v1958 = v1918
		goto L187
	} else {
		goto L245
	}
L241:
	;
	v1899 = v1880
	v1900 = v1892
	v1901 = v1881
	goto L242
L242:
	;
	v1905 = *(*int32)(unsafe.Add(mBase, uint32(v1899)))
	*(*int32)(unsafe.Add(mBase, uint32(v1901))) = v1905
	v1907 = int32(4)
	v1908 = v1899 + v1907
	v1910 = v1901 + v1907
	v1912 = v1900 + int32(-1)
	if v1912 != 0 {
		v1899 = v1908
		v1900 = v1912
		v1901 = v1910
		goto L242
	} else {
		goto L244
	}
L243:
	;
	v1914 = v1908
	v1916 = v1910
	v1918 = v1882 - v1892<<(uint(int32(2))%32)
	goto L240
L244:
	;
	goto L243
L245:
	;
	v1923 = v1914
	v1925 = v1916
	v1927 = v1918
	goto L246
L246:
	;
	v1929 = *(*int32)(unsafe.Add(mBase, uint32(v1923)))
	*(*int32)(unsafe.Add(mBase, uint32(v1925))) = v1929
	v1931 = *(*int32)(unsafe.Add(mBase, uint32(v1923)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1925)+4)) = v1931
	v1933 = *(*int32)(unsafe.Add(mBase, uint32(v1923)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v1925)+8)) = v1933
	v1935 = *(*int32)(unsafe.Add(mBase, uint32(v1923)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v1925)+12)) = v1935
	v1937 = *(*int32)(unsafe.Add(mBase, uint32(v1923)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v1925)+16)) = v1937
	v1939 = *(*int32)(unsafe.Add(mBase, uint32(v1923)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v1925)+20)) = v1939
	v1941 = *(*int32)(unsafe.Add(mBase, uint32(v1923)+24))
	*(*int32)(unsafe.Add(mBase, uint32(v1925)+24)) = v1941
	v1943 = *(*int32)(unsafe.Add(mBase, uint32(v1923)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v1925)+28)) = v1943
	v1945 = int32(32)
	v1946 = v1923 + v1945
	v1948 = v1925 + v1945
	v1950 = v1927 + int32(-32)
	if base.Ui32(int32(3)) < base.Ui32(v1950) {
		v1923 = v1946
		v1925 = v1948
		v1927 = v1950
		goto L246
	} else {
		goto L248
	}
L247:
	;
	v1954 = v1946
	v1956 = v1948
	v1958 = v1950
	goto L187
L248:
	;
	goto L247
L249:
	;
	v1963 = v1958 & int32(7)
	if v1963 != 0 {
		goto L251
	} else {
		goto L252
	}
L250:
	;
	if base.Ui32(v1958) < base.Ui32(int32(8)) {
		goto L181
	} else {
		goto L256
	}
L251:
	;
	v1967 = v1954
	v1968 = v1963
	v1969 = v1956
	goto L253
L252:
	;
	v1982 = v1954
	v1984 = v1956
	v1985 = v1958
	goto L250
L253:
	;
	v1973 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1967))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1969))) = uint8(v1973)
	v1975 = int32(1)
	v1976 = v1969 + v1975
	v1978 = v1967 + v1975
	v1980 = v1968 + int32(-1)
	if v1980 != 0 {
		v1967 = v1978
		v1968 = v1980
		v1969 = v1976
		goto L253
	} else {
		goto L255
	}
L254:
	;
	v1982 = v1978
	v1984 = v1976
	v1985 = v1958 & int32(-8)
	goto L250
L255:
	;
	goto L254
L256:
	;
	v1991 = v1982
	v1993 = v1984
	v1994 = v1985
	goto L257
L257:
	;
	v1997 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993))) = uint8(v1997)
	v1999 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991)+1)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993)+1)) = uint8(v1999)
	v2001 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991)+2)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993)+2)) = uint8(v2001)
	v2003 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991)+3)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993)+3)) = uint8(v2003)
	v2005 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991)+4)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993)+4)) = uint8(v2005)
	v2007 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991)+5)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993)+5)) = uint8(v2007)
	v2009 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991)+6)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993)+6)) = uint8(v2009)
	v2011 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1991)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1993)+7)) = uint8(v2011)
	v2013 = int32(8)
	v2018 = v1994 + int32(-8)
	if v2018 != 0 {
		v1991 = v1991 + v2013
		v1993 = v1993 + v2013
		v1994 = v2018
		goto L257
	} else {
		goto L259
	}
L258:
	;
	goto L181
L259:
	;
	goto L258
L260:
	;
	goto L173
L261:
	;
	v2068 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v903))))
	if v39 == int32(1) {
		v2201 = v2068
		goto L267
	} else {
		goto L268
	}
L262:
	;
	v2065 = *(*int32)(unsafe.Add(mBase, uint32(v2060)+4))
	*(*uint8)(unsafe.Add(mBase, uint32(v903+v2065))) = uint8(v2063)
	goto L261
L263:
	;
	v2042 = l3
	v2045 = v2037
	v2046 = v2039
	goto L264
L264:
	;
	v2047 = int32(4)
	v2051 = v2045 + int32(1)
	F_SetBitDepths(m, v1269+v2046<<(uint(v2047)%32), v1269, v903, v2051)
	mBase = m.M
	v2053 = *(*int32)(unsafe.Add(mBase, uint32(v2042)+12))
	v2056 = v1269 + v2053<<(uint(v2047)%32)
	v2057 = *(*int32)(unsafe.Add(mBase, uint32(v2056)+8))
	if int32(-1) < v2057 {
		v2042 = v2056
		v2045 = v2051
		v2046 = v2057
		goto L264
	} else {
		goto L266
	}
L265:
	;
	v2060 = v2056
	v2063 = v2051
	goto L262
L266:
	;
	goto L265
L267:
	;
	if l1 < v2201 {
		v1292 = v1292 << (uint(int32(1)) % 32)
		goto L148
	} else {
		goto L281
	}
L268:
	;
	if base.Ui32(v39) <= base.Ui32(int32(4)) {
		v2126 = v2068
		v2136 = int32(1)
		goto L269
	} else {
		goto L270
	}
L269:
	;
	v2161 = v2126
	v2170 = v903 + v2136
	v2171 = v39 - v2136
	goto L275
L270:
	;
	v2075 = v903 + v1270
	v2084 = v1281
	v2091 = base.Simd_g_i32x4_splat(v2068)
	goto L271
L271:
	;
	v2107 = base.Simd_g_v128_load32_zero(m, v2075, int32(0))
	v2110 = base.Simd_g_i32x4_max_s(v2091, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v2107)))
	v2114 = v2084 + int32(-4)
	if v2114 != 0 {
		v2075 = v2075 + int32(4)
		v2084 = v2114
		v2091 = v2110
		goto L271
	} else {
		goto L273
	}
L272:
	;
	v2117 = base.Simd_g_i32x4_max_s(v2110, base.Simd_g_i8x16_swizzle_c(v2110, base.Simd_g_const(&F_VP8LCreateHuffmanTree__k1)))
	v2122 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_max_s(v2117, base.Simd_g_i8x16_swizzle_c(v2117, base.Simd_g_const(&F_VP8LCreateHuffmanTree__k2))))
	if v1279 == v1281 {
		v2201 = v2122
		goto L267
	} else {
		goto L274
	}
L273:
	;
	goto L272
L274:
	;
	v2126 = v2122
	v2136 = v1281 | v1270
	goto L269
L275:
	;
	v2192 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2170))))
	if v2192 < v2161 {
		goto L277
	} else {
		goto L278
	}
L276:
	;
	v2201 = v2194
	goto L267
L277:
	;
	v2194 = v2161
	goto L279
L278:
	;
	v2194 = v2192
	goto L279
L279:
	;
	v2198 = v2171 + int32(-1)
	if v2198 != 0 {
		v2161 = v2194
		v2170 = v2170 + int32(1)
		v2171 = v2198
		goto L275
	} else {
		goto L280
	}
L280:
	;
	goto L276
L281:
	;
	goto L149
L282:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v36)+64)) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+72)) = v2518
	v2553 = int32(1)
	v2554 = (v2518 + v2516) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+76)) = v2554
	v2558 = (v2554 + v2536) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+80)) = v2558
	v2562 = (v2558 + v2539) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+84)) = v2562
	v2566 = (v2562 + v2531) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+88)) = v2566
	v2570 = (v2566 + v2523) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+92)) = v2570
	v2574 = (v2570 + v2532) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+96)) = v2574
	v2578 = (v2574 + v2530) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+100)) = v2578
	v2582 = (v2578 + v2524) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+104)) = v2582
	v2586 = (v2582 + v2525) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+108)) = v2586
	v2590 = (v2586 + v2529) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+112)) = v2590
	v2594 = (v2590 + v2526) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+116)) = v2594
	v2598 = (v2594 + v2528) << (uint(v2553) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v36)+120)) = v2598
	*(*int32)(unsafe.Add(mBase, uint32(v36)+124)) = (v2598 + v2527) << (uint(v2553) % 32)
	if v2283 < int32(1) {
		goto L295
	} else {
		goto L296
	}
L283:
	;
	v2303 = v2283 & int32(3)
	v2304 = *(*int32)(unsafe.Add(mBase, uint32(v2254)))
	if base.Ui32(v2283) < base.Ui32(int32(4)) {
		v2397 = int32(0)
		goto L285
	} else {
		goto L286
	}
L284:
	;
	v2288 = int32(0)
	v2516 = v2288
	v2518 = v2288
	v2523 = v2288
	v2524 = v2288
	v2525 = v2288
	v2526 = v2288
	v2527 = v2288
	v2528 = v2288
	v2529 = v2288
	v2530 = v2288
	v2531 = v2288
	v2532 = v2288
	v2536 = v2288
	v2539 = v2288
	goto L282
L285:
	;
	if v2303 == int32(0) {
		goto L290
	} else {
		goto L291
	}
L286:
	;
	v2322 = int32(0)
	goto L287
L287:
	;
	v2344 = v2304 + v2322
	v2345 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2344))))
	v2346 = int32(2)
	v2348 = v36 + v2345<<(uint(v2346)%32)
	v2349 = *(*int32)(unsafe.Add(mBase, uint32(v2348)))
	v2350 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2348))) = v2349 + v2350
	v2355 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2344+v2350))))
	v2358 = v36 + v2355<<(uint(v2346)%32)
	v2359 = *(*int32)(unsafe.Add(mBase, uint32(v2358)))
	*(*int32)(unsafe.Add(mBase, uint32(v2358))) = v2359 + v2350
	v2365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2344+v2346))))
	v2368 = v36 + v2365<<(uint(v2346)%32)
	v2369 = *(*int32)(unsafe.Add(mBase, uint32(v2368)))
	*(*int32)(unsafe.Add(mBase, uint32(v2368))) = v2369 + v2350
	v2375 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2344+int32(3)))))
	v2378 = v36 + v2375<<(uint(v2346)%32)
	v2379 = *(*int32)(unsafe.Add(mBase, uint32(v2378)))
	*(*int32)(unsafe.Add(mBase, uint32(v2378))) = v2379 + v2350
	v2384 = v2322 + int32(4)
	if v2283&int32(2147483644) != v2384 {
		v2322 = v2384
		goto L287
	} else {
		goto L289
	}
L288:
	;
	v2397 = v2384
	goto L285
L289:
	;
	goto L288
L290:
	;
	v2500 = *(*int32)(unsafe.Add(mBase, uint32(v36)+4))
	v2503 = *(*int32)(unsafe.Add(mBase, uint32(v36)+56))
	v2504 = *(*int32)(unsafe.Add(mBase, uint32(v36)+52))
	v2505 = *(*int32)(unsafe.Add(mBase, uint32(v36)+48))
	v2506 = *(*int32)(unsafe.Add(mBase, uint32(v36)+44))
	v2507 = *(*int32)(unsafe.Add(mBase, uint32(v36)+40))
	v2508 = *(*int32)(unsafe.Add(mBase, uint32(v36)+36))
	v2509 = *(*int32)(unsafe.Add(mBase, uint32(v36)+32))
	v2510 = *(*int32)(unsafe.Add(mBase, uint32(v36)+28))
	v2511 = *(*int32)(unsafe.Add(mBase, uint32(v36)+24))
	v2512 = *(*int32)(unsafe.Add(mBase, uint32(v36)+20))
	v2513 = *(*int32)(unsafe.Add(mBase, uint32(v36)+16))
	v2514 = *(*int32)(unsafe.Add(mBase, uint32(v36)+12))
	v2515 = *(*int32)(unsafe.Add(mBase, uint32(v36)+8))
	v2516 = v2515
	v2518 = v2500 << (uint(int32(1)) % 32)
	v2523 = v2511
	v2524 = v2508
	v2525 = v2507
	v2526 = v2505
	v2527 = v2503
	v2528 = v2504
	v2529 = v2506
	v2530 = v2509
	v2531 = v2512
	v2532 = v2510
	v2536 = v2514
	v2539 = v2513
	goto L282
L291:
	;
	v2424 = v2304 + v2397
	v2434 = v2303
	goto L292
L292:
	;
	v2455 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2424))))
	v2458 = v36 + v2455<<(uint(int32(2))%32)
	v2459 = *(*int32)(unsafe.Add(mBase, uint32(v2458)))
	v2460 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2458))) = v2459 + v2460
	v2466 = v2434 + int32(-1)
	if v2466 != 0 {
		v2424 = v2424 + v2460
		v2434 = v2466
		goto L292
	} else {
		goto L294
	}
L293:
	;
	goto L290
L294:
	;
	goto L293
L295:
	;
	m.G0 = v36 + int32(128)
	return
L296:
	;
	v2604 = *(*int32)(unsafe.Add(mBase, uint32(l4)+8))
	v2605 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
	v2610 = int32(0)
	goto L297
L297:
	;
	v2643 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2605+v2610))))
	v2646 = v36 + int32(64) + v2643<<(uint(int32(2))%32)
	v2647 = *(*int32)(unsafe.Add(mBase, uint32(v2646)))
	*(*int32)(unsafe.Add(mBase, uint32(v2646))) = v2647 + int32(1)
	v2651 = int32(0)
	if v2643 == v2651 {
		v2788 = v2651
		goto L299
	} else {
		goto L300
	}
L298:
	;
	goto L295
L299:
	;
	v2809 = int32(1)
	v2814 = int32(base.Ui32(v2788) >> (uint(int32(16)-v2643) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v2604+v2610<<(uint(v2809)%32)))) = uint16(v2814)
	v2817 = v2610 + v2809
	if v2817 != v2283 {
		v2610 = v2817
		goto L297
	} else {
		goto L308
	}
L300:
	;
	v2655 = v2643 + int32(-1)
	if base.Ui32(int32(5)) <= base.Ui32(v2643) {
		goto L302
	} else {
		goto L303
	}
L301:
	;
	if v2655&int32(4) != 0 {
		v2788 = v2742
		goto L299
	} else {
		goto L307
	}
L302:
	;
	v2666 = int32(0)
	v2670 = v2647
	v2678 = (int32(base.Ui32(v2655)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
	v2679 = v2666
	v2680 = v2666
	goto L304
L303:
	;
	v2732 = v2647
	v2741 = int32(0)
	v2742 = int32(0)
	goto L301
L304:
	;
	v2701 = m.G1
	v2703 = v2701 + int32(_a_F_VP8LCreateHuffmanTree_0)
	v2706 = int32(15)
	v2709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2703+int32(base.Ui32(v2670)>>(uint(int32(4))%32))&v2706))))
	v2710 = int32(8)
	v2716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2703+v2670&v2706))))
	v2721 = v2709<<(uint(v2679+v2710)%32) | (v2716<<(uint(v2679+int32(12))%32) | v2680)
	v2723 = v2679 + int32(-8)
	v2725 = int32(base.Ui32(v2670) >> (uint(v2710) % 32))
	v2727 = v2678 + int32(-2)
	if v2727 != 0 {
		v2670 = v2725
		v2678 = v2727
		v2679 = v2723
		v2680 = v2721
		goto L304
	} else {
		goto L306
	}
L305:
	;
	v2732 = v2725
	v2741 = int32(0) - v2723
	v2742 = v2721
	goto L301
L306:
	;
	goto L305
L307:
	;
	v2765 = m.G1
	v2771 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2765+int32(_a_F_VP8LCreateHuffmanTree_0)+v2732&int32(15)))))
	v2788 = v2771<<(uint(int32(12)-v2741)%32) | v2742
	goto L299
L308:
	;
	goto L298
}

var F_VP8LCreateHuffmanTree__k0 = [2]uint64{0x0, 0x0}
var F_VP8LCreateHuffmanTree__k1 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_VP8LCreateHuffmanTree__k2 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_VP8LDspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
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
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v189 int32
	_ = v189
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v227 int32
	_ = v227
	var v253 int32
	_ = v253
	var v254 int32
	_ = v254
	var v269 int32
	_ = v269
	var v270 int32
	_ = v270
	var v277 int32
	_ = v277
	var v278 int32
	_ = v278
	var v282 int32
	_ = v282
	var v284 int32
	_ = v284
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v290 int32
	_ = v290
	var v292 int32
	_ = v292
	var v294 int32
	_ = v294
	var v296 int32
	_ = v296
	var v298 int32
	_ = v298
	var v302 int32
	_ = v302
	var v306 int32
	_ = v306
	var v310 int32
	_ = v310
	var v314 base.V128
	_ = v314
	var v320 base.V128
	_ = v320
	var v324 base.V128
	_ = v324
	var v330 base.V128
	_ = v330
	var v333 int32
	_ = v333
	var v335 int32
	_ = v335
	var v337 int32
	_ = v337
	var v339 int32
	_ = v339
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v350 int32
	_ = v350
	var v351 int32
	_ = v351
	var v353 int32
	_ = v353
	var v355 int32
	_ = v355
	var v357 int32
	_ = v357
	var v359 int32
	_ = v359
	var v363 int32
	_ = v363
	var v365 int32
	_ = v365
	var v367 int32
	_ = v367
	var v370 int32
	_ = v370
	var v371 int32
	_ = v371
	v19 = m.G1
	v22 = *(*int32)(unsafe.Add(mBase, uint32(v19)+uint32(_c_F_VP8LDspInit[0])))
	v23 = m.G13
	v24 = *(*int32)(unsafe.Add(mBase, uint32(v23)))
	if v22 == v24 {
	} else {
		v26 = m.G2
		v27 = m.G1
		v31 = v26 + int32(90)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[1]))) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[2]))) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[3]))) = v26 + int32(91)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[4]))) = v26 + int32(92)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[5]))) = v26 + int32(93)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[6]))) = v26 + int32(94)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[7]))) = v26 + int32(95)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[8]))) = v26 + int32(96)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[9]))) = v26 + int32(97)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[10]))) = v26 + int32(98)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[11]))) = v26 + int32(99)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[12]))) = v26 + int32(100)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[13]))) = v26 + int32(101)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[14]))) = v26 + int32(102)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[15]))) = v26 + int32(103)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[16]))) = v31
		v77 = v26 + int32(104)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[17]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[18]))) = v77
		v81 = v26 + int32(105)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[19]))) = v81
		v84 = v26 + int32(106)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[20]))) = v84
		v87 = v26 + int32(107)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[21]))) = v87
		v90 = v26 + int32(108)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[22]))) = v90
		v93 = v26 + int32(109)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[23]))) = v93
		v96 = v26 + int32(110)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[24]))) = v96
		v99 = v26 + int32(111)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[25]))) = v99
		v102 = v26 + int32(112)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[26]))) = v102
		v105 = v26 + int32(113)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[27]))) = v105
		v108 = v26 + int32(114)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[28]))) = v108
		v111 = v26 + int32(115)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[29]))) = v111
		v114 = v26 + int32(116)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[30]))) = v114
		v117 = v26 + int32(117)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[31]))) = v117
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[32]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[33]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[34]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[35]))) = v81
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[36]))) = v84
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[37]))) = v87
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[38]))) = v90
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[39]))) = v93
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[40]))) = v96
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[41]))) = v99
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[42]))) = v102
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[43]))) = v105
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[44]))) = v108
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[45]))) = v111
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[46]))) = v114
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[47]))) = v117
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[48]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[49]))) = v26 + int32(118)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[50]))) = v26 + int32(119)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[51]))) = v26 + int32(120)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[52]))) = v26 + int32(121)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[53]))) = v26 + int32(122)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[54]))) = v26 + int32(123)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[55]))) = v26 + int32(124)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[56]))) = v26 + int32(125)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[57]))) = v26 + int32(126)
		if v24 == int32(0) {
		} else {
			v185 = int32(0)
			v186 = m.T0[v24].(func(*base.Module, int32) int32)(m, v185)
			mBase = m.M
			if v186 == v185 {
			} else {
				v189 = int32(0)
				v198 = m.G2
				v199 = m.G20
				*(*int32)(unsafe.Add(mBase, uint32(v199)+52)) = v198 + int32(57)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+48)) = v198 + int32(58)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+44)) = v198 + int32(59)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+40)) = v198 + int32(60)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+36)) = v198 + int32(61)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+32)) = v198 + int32(62)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+28)) = v198 + int32(63)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+24)) = v198 + int32(64)
				*(*int32)(unsafe.Add(mBase, uint32(v199)+20)) = v198 + int32(65)
				v227 = m.G21
				*(*int32)(unsafe.Add(mBase, uint32(v227)+12)) = v198 + int32(66)
				*(*int32)(unsafe.Add(mBase, uint32(v227+int32(8)))) = v198 + int32(67)
				*(*int32)(unsafe.Add(mBase, uint32(v227)+4)) = v198 + int32(68)
				*(*int32)(unsafe.Add(mBase, uint32(v227))) = v198 + int32(69)
				*(*int32)(unsafe.Add(mBase, uint32(v227)+28)) = v198 + int32(70)
				*(*int32)(unsafe.Add(mBase, uint32(v227+int32(24)))) = v198 + int32(71)
				*(*int32)(unsafe.Add(mBase, uint32(v227)+20)) = v198 + int32(72)
				v253 = int32(16)
				v254 = v227 + v253
				*(*int32)(unsafe.Add(mBase, uint32(v254))) = v198 + int32(73)
				*(*int32)(unsafe.Add(mBase, uint32(v227)+44)) = v198 + int32(74)
				*(*int32)(unsafe.Add(mBase, uint32(v227+int32(40)))) = v198 + int32(75)
				*(*int32)(unsafe.Add(mBase, uint32(v227)+36)) = v198 + int32(76)
				v269 = int32(32)
				v270 = v227 + v269
				*(*int32)(unsafe.Add(mBase, uint32(v270))) = v198 + int32(77)
				*(*int32)(unsafe.Add(mBase, uint32(v227)+52)) = v198 + int32(78)
				v277 = int32(48)
				v278 = v227 + v277
				*(*int32)(unsafe.Add(mBase, uint32(v278))) = v198 + int32(79)
				v282 = m.G22
				v284 = v198 + int32(80)
				*(*int32)(unsafe.Add(mBase, uint32(v282))) = v284
				v286 = m.G23
				v288 = v198 + int32(81)
				*(*int32)(unsafe.Add(mBase, uint32(v286))) = v288
				v290 = m.G24
				v292 = v198 + int32(82)
				*(*int32)(unsafe.Add(mBase, uint32(v290))) = v292
				v294 = m.G25
				v296 = v198 + int32(83)
				*(*int32)(unsafe.Add(mBase, uint32(v294))) = v296
				v298 = m.G26
				*(*int32)(unsafe.Add(mBase, uint32(v298))) = v198 + int32(84)
				v302 = m.G27
				*(*int32)(unsafe.Add(mBase, uint32(v302))) = v198 + int32(85)
				v306 = m.G28
				*(*int32)(unsafe.Add(mBase, uint32(v306))) = v198 + int32(86)
				v310 = m.G29
				v314 = base.Simd_g_v128_load(m, v270, v189)
				base.Simd_g_v128_store(m, v310+v269, v189, v314)
				v320 = base.Simd_g_v128_load(m, v254, v189)
				base.Simd_g_v128_store(m, v310+v253, v189, v320)
				v324 = base.Simd_g_v128_load(m, v227, v189)
				base.Simd_g_v128_store(m, v310, v189, v324)
				v330 = base.Simd_g_v128_load(m, v278, v189)
				base.Simd_g_v128_store(m, v310+v277, v189, v330)
				v333 = m.G30
				*(*int32)(unsafe.Add(mBase, uint32(v333))) = v284
				v335 = m.G31
				*(*int32)(unsafe.Add(mBase, uint32(v335))) = v288
				v337 = m.G32
				*(*int32)(unsafe.Add(mBase, uint32(v337))) = v292
				v339 = m.G33
				*(*int32)(unsafe.Add(mBase, uint32(v339))) = v296
				v342 = m.G13
				v343 = *(*int32)(unsafe.Add(mBase, uint32(v342)))
				v344 = m.T0[v343].(func(*base.Module, int32) int32)(m, int32(3))
				mBase = m.M
				if v344 == int32(0) {
				} else {
					v350 = m.G2
					v351 = m.G24
					v353 = v350 + int32(87)
					*(*int32)(unsafe.Add(mBase, uint32(v351))) = v353
					v355 = m.G23
					v357 = v350 + int32(88)
					*(*int32)(unsafe.Add(mBase, uint32(v355))) = v357
					v359 = m.G28
					*(*int32)(unsafe.Add(mBase, uint32(v359))) = v350 + int32(89)
					v363 = m.G31
					*(*int32)(unsafe.Add(mBase, uint32(v363))) = v357
					v365 = m.G32
					*(*int32)(unsafe.Add(mBase, uint32(v365))) = v353
				}
			}
		}
		v367 = m.G1
		v370 = m.G13
		v371 = *(*int32)(unsafe.Add(mBase, uint32(v370)))
		*(*int32)(unsafe.Add(mBase, uint32(v367)+uint32(_c_F_VP8LDspInit[0]))) = v371
	}
	return
}
func F_VP8LDspInitSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v1 int32
	_ = v1
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v39 int32
	_ = v39
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v114 int32
	_ = v114
	var v118 int32
	_ = v118
	var v122 int32
	_ = v122
	var v126 base.V128
	_ = v126
	var v132 base.V128
	_ = v132
	var v136 base.V128
	_ = v136
	var v142 base.V128
	_ = v142
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v149 int32
	_ = v149
	var v151 int32
	_ = v151
	v1 = int32(0)
	v10 = m.G2
	v11 = m.G20
	*(*int32)(unsafe.Add(mBase, uint32(v11)+52)) = v10 + int32(57)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+48)) = v10 + int32(58)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+44)) = v10 + int32(59)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+40)) = v10 + int32(60)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+36)) = v10 + int32(61)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+32)) = v10 + int32(62)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+28)) = v10 + int32(63)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+24)) = v10 + int32(64)
	*(*int32)(unsafe.Add(mBase, uint32(v11)+20)) = v10 + int32(65)
	v39 = m.G21
	*(*int32)(unsafe.Add(mBase, uint32(v39)+12)) = v10 + int32(66)
	*(*int32)(unsafe.Add(mBase, uint32(v39+int32(8)))) = v10 + int32(67)
	*(*int32)(unsafe.Add(mBase, uint32(v39)+4)) = v10 + int32(68)
	*(*int32)(unsafe.Add(mBase, uint32(v39))) = v10 + int32(69)
	*(*int32)(unsafe.Add(mBase, uint32(v39)+28)) = v10 + int32(70)
	*(*int32)(unsafe.Add(mBase, uint32(v39+int32(24)))) = v10 + int32(71)
	*(*int32)(unsafe.Add(mBase, uint32(v39)+20)) = v10 + int32(72)
	v65 = int32(16)
	v66 = v39 + v65
	*(*int32)(unsafe.Add(mBase, uint32(v66))) = v10 + int32(73)
	*(*int32)(unsafe.Add(mBase, uint32(v39)+44)) = v10 + int32(74)
	*(*int32)(unsafe.Add(mBase, uint32(v39+int32(40)))) = v10 + int32(75)
	*(*int32)(unsafe.Add(mBase, uint32(v39)+36)) = v10 + int32(76)
	v81 = int32(32)
	v82 = v39 + v81
	*(*int32)(unsafe.Add(mBase, uint32(v82))) = v10 + int32(77)
	*(*int32)(unsafe.Add(mBase, uint32(v39)+52)) = v10 + int32(78)
	v89 = int32(48)
	v90 = v39 + v89
	*(*int32)(unsafe.Add(mBase, uint32(v90))) = v10 + int32(79)
	v94 = m.G22
	v96 = v10 + int32(80)
	*(*int32)(unsafe.Add(mBase, uint32(v94))) = v96
	v98 = m.G23
	v100 = v10 + int32(81)
	*(*int32)(unsafe.Add(mBase, uint32(v98))) = v100
	v102 = m.G24
	v104 = v10 + int32(82)
	*(*int32)(unsafe.Add(mBase, uint32(v102))) = v104
	v106 = m.G25
	v108 = v10 + int32(83)
	*(*int32)(unsafe.Add(mBase, uint32(v106))) = v108
	v110 = m.G26
	*(*int32)(unsafe.Add(mBase, uint32(v110))) = v10 + int32(84)
	v114 = m.G27
	*(*int32)(unsafe.Add(mBase, uint32(v114))) = v10 + int32(85)
	v118 = m.G28
	*(*int32)(unsafe.Add(mBase, uint32(v118))) = v10 + int32(86)
	v122 = m.G29
	v126 = base.Simd_g_v128_load(m, v82, v1)
	base.Simd_g_v128_store(m, v122+v81, v1, v126)
	v132 = base.Simd_g_v128_load(m, v66, v1)
	base.Simd_g_v128_store(m, v122+v65, v1, v132)
	v136 = base.Simd_g_v128_load(m, v39, v1)
	base.Simd_g_v128_store(m, v122, v1, v136)
	v142 = base.Simd_g_v128_load(m, v90, v1)
	base.Simd_g_v128_store(m, v122+v89, v1, v142)
	v145 = m.G30
	*(*int32)(unsafe.Add(mBase, uint32(v145))) = v96
	v147 = m.G31
	*(*int32)(unsafe.Add(mBase, uint32(v147))) = v100
	v149 = m.G32
	*(*int32)(unsafe.Add(mBase, uint32(v149))) = v104
	v151 = m.G33
	*(*int32)(unsafe.Add(mBase, uint32(v151))) = v108
	return
}
func F_VP8LEncDspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v106 int32
	_ = v106
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v167 int32
	_ = v167
	var v168 int32
	_ = v168
	var v171 int32
	_ = v171
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v202 int32
	_ = v202
	var v206 int32
	_ = v206
	var v210 int32
	_ = v210
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v235 int32
	_ = v235
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
	var v264 int32
	_ = v264
	var v265 int32
	_ = v265
	var v276 int32
	_ = v276
	var v277 int32
	_ = v277
	var v281 int32
	_ = v281
	var v283 int32
	_ = v283
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v289 int32
	_ = v289
	var v291 int32
	_ = v291
	var v295 base.V128
	_ = v295
	var v301 base.V128
	_ = v301
	var v305 base.V128
	_ = v305
	var v311 base.V128
	_ = v311
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v317 int32
	_ = v317
	var v323 int32
	_ = v323
	var v324 int32
	_ = v324
	var v326 int32
	_ = v326
	var v328 int32
	_ = v328
	var v332 int32
	_ = v332
	var v334 int32
	_ = v334
	var v336 int32
	_ = v336
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v342 int32
	_ = v342
	var v344 int32
	_ = v344
	var v346 int32
	_ = v346
	var v349 int32
	_ = v349
	var v350 int32
	_ = v350
	v18 = m.G1
	v21 = *(*int32)(unsafe.Add(mBase, uint32(v18)+uint32(_c_F_VP8LEncDspInit[0])))
	v22 = m.G13
	v23 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
	if v21 == v23 {
	} else {
		v25 = m.G2
		v26 = m.G1
		F_VP8LDspInit(m)
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[1]))) = v25 + int32(296)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[2]))) = v25 + int32(297)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[3]))) = v25 + int32(298)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[4]))) = v25 + int32(299)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[5]))) = v25 + int32(300)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[6]))) = v25 + int32(301)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[7]))) = v25 + int32(302)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[8]))) = v25 + int32(303)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[9]))) = v25 + int32(304)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[10]))) = v25 + int32(305)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[11]))) = v25 + int32(306)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[12]))) = v25 + int32(307)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[13]))) = v25 + int32(308)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[14]))) = v25 + int32(309)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[15]))) = v25 + int32(310)
		v106 = v25 + int32(311)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[16]))) = v106
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[17]))) = v106
		v110 = v25 + int32(312)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[18]))) = v110
		v113 = v25 + int32(313)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[19]))) = v113
		v116 = v25 + int32(314)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[20]))) = v116
		v119 = v25 + int32(315)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[21]))) = v119
		v122 = v25 + int32(316)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[22]))) = v122
		v125 = v25 + int32(317)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[23]))) = v125
		v128 = v25 + int32(318)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[24]))) = v128
		v131 = v25 + int32(319)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[25]))) = v131
		v134 = v25 + int32(320)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[26]))) = v134
		v137 = v25 + int32(321)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[27]))) = v137
		v140 = v25 + int32(322)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[28]))) = v140
		v143 = v25 + int32(323)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[29]))) = v143
		v146 = v25 + int32(324)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[30]))) = v146
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[31]))) = v106
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[32]))) = v106
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[33]))) = v106
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[34]))) = v110
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[35]))) = v113
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[36]))) = v116
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[37]))) = v119
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[38]))) = v122
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[39]))) = v125
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[40]))) = v128
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[41]))) = v131
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[42]))) = v134
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[43]))) = v137
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[44]))) = v140
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[45]))) = v143
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[46]))) = v146
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[47]))) = v106
		v167 = m.G13
		v168 = *(*int32)(unsafe.Add(mBase, uint32(v167)))
		if v168 == int32(0) {
		} else {
			v171 = int32(0)
			v172 = m.T0[v168].(func(*base.Module, int32) int32)(m, v171)
			mBase = m.M
			if v172 == v171 {
			} else {
				v175 = int32(0)
				v185 = m.G2
				v186 = m.G93
				v188 = v185 + int32(269)
				*(*int32)(unsafe.Add(mBase, uint32(v186))) = v188
				v190 = m.G94
				v192 = v185 + int32(270)
				*(*int32)(unsafe.Add(mBase, uint32(v190))) = v192
				v194 = m.G95
				v196 = v185 + int32(271)
				*(*int32)(unsafe.Add(mBase, uint32(v194))) = v196
				v198 = m.G96
				v200 = v185 + int32(272)
				*(*int32)(unsafe.Add(mBase, uint32(v198))) = v200
				v202 = m.G97
				*(*int32)(unsafe.Add(mBase, uint32(v202))) = v185 + int32(273)
				v206 = m.G98
				*(*int32)(unsafe.Add(mBase, uint32(v206))) = v185 + int32(274)
				v210 = m.G99
				*(*int32)(unsafe.Add(mBase, uint32(v210))) = v185 + int32(275)
				v214 = m.G100
				*(*int32)(unsafe.Add(mBase, uint32(v214))) = v185 + int32(276)
				v218 = m.G101
				v220 = v185 + int32(277)
				*(*int32)(unsafe.Add(mBase, uint32(v218))) = v220
				v222 = m.G102
				*(*int32)(unsafe.Add(mBase, uint32(v222)+12)) = v185 + int32(278)
				*(*int32)(unsafe.Add(mBase, uint32(v222+int32(8)))) = v185 + int32(279)
				*(*int32)(unsafe.Add(mBase, uint32(v222)+4)) = v185 + int32(280)
				v235 = v185 + int32(281)
				*(*int32)(unsafe.Add(mBase, uint32(v222))) = v235
				*(*int32)(unsafe.Add(mBase, uint32(v222)+28)) = v185 + int32(282)
				*(*int32)(unsafe.Add(mBase, uint32(v222+int32(24)))) = v185 + int32(283)
				*(*int32)(unsafe.Add(mBase, uint32(v222)+20)) = v185 + int32(284)
				v248 = int32(16)
				v249 = v222 + v248
				*(*int32)(unsafe.Add(mBase, uint32(v249))) = v185 + int32(285)
				*(*int32)(unsafe.Add(mBase, uint32(v222)+44)) = v185 + int32(286)
				*(*int32)(unsafe.Add(mBase, uint32(v222+int32(40)))) = v185 + int32(287)
				*(*int32)(unsafe.Add(mBase, uint32(v222)+36)) = v185 + int32(288)
				v264 = int32(32)
				v265 = v222 + v264
				*(*int32)(unsafe.Add(mBase, uint32(v265))) = v185 + int32(289)
				*(*int32)(unsafe.Add(mBase, uint32(v222)+60)) = v235
				*(*int32)(unsafe.Add(mBase, uint32(v222+int32(56)))) = v235
				*(*int32)(unsafe.Add(mBase, uint32(v222)+52)) = v185 + int32(290)
				v276 = int32(48)
				v277 = v222 + v276
				*(*int32)(unsafe.Add(mBase, uint32(v277))) = v185 + int32(291)
				v281 = m.G103
				*(*int32)(unsafe.Add(mBase, uint32(v281))) = v192
				v283 = m.G104
				*(*int32)(unsafe.Add(mBase, uint32(v283))) = v188
				v285 = m.G105
				*(*int32)(unsafe.Add(mBase, uint32(v285))) = v196
				v287 = m.G106
				*(*int32)(unsafe.Add(mBase, uint32(v287))) = v200
				v289 = m.G107
				*(*int32)(unsafe.Add(mBase, uint32(v289))) = v220
				v291 = m.G108
				v295 = base.Simd_g_v128_load(m, v265, v175)
				base.Simd_g_v128_store(m, v291+v264, v175, v295)
				v301 = base.Simd_g_v128_load(m, v249, v175)
				base.Simd_g_v128_store(m, v291+v248, v175, v301)
				v305 = base.Simd_g_v128_load(m, v222, v175)
				base.Simd_g_v128_store(m, v291, v175, v305)
				v311 = base.Simd_g_v128_load(m, v277, v175)
				base.Simd_g_v128_store(m, v291+v276, v175, v311)
				v315 = m.G13
				v316 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
				v317 = m.T0[v316].(func(*base.Module, int32) int32)(m, int32(3))
				mBase = m.M
				if v317 == int32(0) {
				} else {
					v323 = m.G2
					v324 = m.G94
					v326 = v323 + int32(292)
					*(*int32)(unsafe.Add(mBase, uint32(v324))) = v326
					v328 = m.G112
					*(*int32)(unsafe.Add(mBase, uint32(v328))) = v323 + int32(293)
					v332 = m.G95
					v334 = v323 + int32(294)
					*(*int32)(unsafe.Add(mBase, uint32(v332))) = v334
					v336 = m.G96
					v338 = v323 + int32(295)
					*(*int32)(unsafe.Add(mBase, uint32(v336))) = v338
					v340 = m.G103
					*(*int32)(unsafe.Add(mBase, uint32(v340))) = v326
					v342 = m.G105
					*(*int32)(unsafe.Add(mBase, uint32(v342))) = v334
					v344 = m.G106
					*(*int32)(unsafe.Add(mBase, uint32(v344))) = v338
				}
			}
		}
		v346 = m.G1
		v349 = m.G13
		v350 = *(*int32)(unsafe.Add(mBase, uint32(v349)))
		*(*int32)(unsafe.Add(mBase, uint32(v346)+uint32(_c_F_VP8LEncDspInit[0]))) = v350
	}
	return
}
func F_VP8LEncDspInitSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v1 int32
	_ = v1
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v61 int32
	_ = v61
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v117 int32
	_ = v117
	var v121 base.V128
	_ = v121
	var v127 base.V128
	_ = v127
	var v131 base.V128
	_ = v131
	var v137 base.V128
	_ = v137
	v1 = int32(0)
	v11 = m.G2
	v12 = m.G93
	v14 = v11 + int32(269)
	*(*int32)(unsafe.Add(mBase, uint32(v12))) = v14
	v16 = m.G94
	v18 = v11 + int32(270)
	*(*int32)(unsafe.Add(mBase, uint32(v16))) = v18
	v20 = m.G95
	v22 = v11 + int32(271)
	*(*int32)(unsafe.Add(mBase, uint32(v20))) = v22
	v24 = m.G96
	v26 = v11 + int32(272)
	*(*int32)(unsafe.Add(mBase, uint32(v24))) = v26
	v28 = m.G97
	*(*int32)(unsafe.Add(mBase, uint32(v28))) = v11 + int32(273)
	v32 = m.G98
	*(*int32)(unsafe.Add(mBase, uint32(v32))) = v11 + int32(274)
	v36 = m.G99
	*(*int32)(unsafe.Add(mBase, uint32(v36))) = v11 + int32(275)
	v40 = m.G100
	*(*int32)(unsafe.Add(mBase, uint32(v40))) = v11 + int32(276)
	v44 = m.G101
	v46 = v11 + int32(277)
	*(*int32)(unsafe.Add(mBase, uint32(v44))) = v46
	v48 = m.G102
	*(*int32)(unsafe.Add(mBase, uint32(v48)+12)) = v11 + int32(278)
	*(*int32)(unsafe.Add(mBase, uint32(v48+int32(8)))) = v11 + int32(279)
	*(*int32)(unsafe.Add(mBase, uint32(v48)+4)) = v11 + int32(280)
	v61 = v11 + int32(281)
	*(*int32)(unsafe.Add(mBase, uint32(v48))) = v61
	*(*int32)(unsafe.Add(mBase, uint32(v48)+28)) = v11 + int32(282)
	*(*int32)(unsafe.Add(mBase, uint32(v48+int32(24)))) = v11 + int32(283)
	*(*int32)(unsafe.Add(mBase, uint32(v48)+20)) = v11 + int32(284)
	v74 = int32(16)
	v75 = v48 + v74
	*(*int32)(unsafe.Add(mBase, uint32(v75))) = v11 + int32(285)
	*(*int32)(unsafe.Add(mBase, uint32(v48)+44)) = v11 + int32(286)
	*(*int32)(unsafe.Add(mBase, uint32(v48+int32(40)))) = v11 + int32(287)
	*(*int32)(unsafe.Add(mBase, uint32(v48)+36)) = v11 + int32(288)
	v90 = int32(32)
	v91 = v48 + v90
	*(*int32)(unsafe.Add(mBase, uint32(v91))) = v11 + int32(289)
	*(*int32)(unsafe.Add(mBase, uint32(v48)+60)) = v61
	*(*int32)(unsafe.Add(mBase, uint32(v48+int32(56)))) = v61
	*(*int32)(unsafe.Add(mBase, uint32(v48)+52)) = v11 + int32(290)
	v102 = int32(48)
	v103 = v48 + v102
	*(*int32)(unsafe.Add(mBase, uint32(v103))) = v11 + int32(291)
	v107 = m.G103
	*(*int32)(unsafe.Add(mBase, uint32(v107))) = v18
	v109 = m.G104
	*(*int32)(unsafe.Add(mBase, uint32(v109))) = v14
	v111 = m.G105
	*(*int32)(unsafe.Add(mBase, uint32(v111))) = v22
	v113 = m.G106
	*(*int32)(unsafe.Add(mBase, uint32(v113))) = v26
	v115 = m.G107
	*(*int32)(unsafe.Add(mBase, uint32(v115))) = v46
	v117 = m.G108
	v121 = base.Simd_g_v128_load(m, v91, v1)
	base.Simd_g_v128_store(m, v117+v90, v1, v121)
	v127 = base.Simd_g_v128_load(m, v75, v1)
	base.Simd_g_v128_store(m, v117+v74, v1, v127)
	v131 = base.Simd_g_v128_load(m, v48, v1)
	base.Simd_g_v128_store(m, v117, v1, v131)
	v137 = base.Simd_g_v128_load(m, v103, v1)
	base.Simd_g_v128_store(m, v117+v102, v1, v137)
	return
}
func F_VP8LEncodeImage(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v31 base.V128
	_ = v31
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v56 int32
	_ = v56
	var v62 int32
	_ = v62
	var v65 int32
	_ = v65
	var v192 base.V128
	_ = v192
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v199 int32
	_ = v199
	var v200 int32
	_ = v200
	var v212 int32
	_ = v212
	var v213 int32
	_ = v213
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v226 int32
	_ = v226
	var v227 int32
	_ = v227
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v231 int32
	_ = v231
	var v232 int32
	_ = v232
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v238 int32
	_ = v238
	var v239 int32
	_ = v239
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
	var v252 int64
	_ = v252
	var v255 int32
	_ = v255
	var v257 int64
	_ = v257
	var v260 int32
	_ = v260
	var v267 int32
	_ = v267
	var v269 int32
	_ = v269
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v289 int32
	_ = v289
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v302 int32
	_ = v302
	var v305 int32
	_ = v305
	var v307 int32
	_ = v307
	var v311 int32
	_ = v311
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v314 int32
	_ = v314
	var v352 int32
	_ = v352
	var v353 int32
	_ = v353
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v370 int32
	_ = v370
	var v371 int32
	_ = v371
	var v379 int32
	_ = v379
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v382 int32
	_ = v382
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v389 int32
	_ = v389
	var v390 int32
	_ = v390
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v401 int32
	_ = v401
	var v402 int32
	_ = v402
	var v405 int64
	_ = v405
	var v408 int32
	_ = v408
	var v410 int64
	_ = v410
	var v413 int32
	_ = v413
	var v420 int32
	_ = v420
	var v422 int32
	_ = v422
	var v426 int32
	_ = v426
	var v427 int32
	_ = v427
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v438 int32
	_ = v438
	var v439 int32
	_ = v439
	var v440 int32
	_ = v440
	var v442 int32
	_ = v442
	var v447 int32
	_ = v447
	var v448 int32
	_ = v448
	var v455 int32
	_ = v455
	var v458 int32
	_ = v458
	var v460 int32
	_ = v460
	var v464 int32
	_ = v464
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v467 int32
	_ = v467
	var v504 int32
	_ = v504
	var v508 int32
	_ = v508
	var v511 int32
	_ = v511
	var v512 int32
	_ = v512
	var v524 int32
	_ = v524
	var v525 int32
	_ = v525
	var v529 int32
	_ = v529
	var v530 int32
	_ = v530
	var v538 int32
	_ = v538
	var v539 int32
	_ = v539
	var v540 int32
	_ = v540
	var v541 int32
	_ = v541
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v548 int32
	_ = v548
	var v549 int32
	_ = v549
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v560 int32
	_ = v560
	var v561 int32
	_ = v561
	var v564 int64
	_ = v564
	var v567 int32
	_ = v567
	var v569 int64
	_ = v569
	var v572 int32
	_ = v572
	var v579 int32
	_ = v579
	var v581 int32
	_ = v581
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v588 int32
	_ = v588
	var v597 int32
	_ = v597
	var v598 int32
	_ = v598
	var v599 int32
	_ = v599
	var v601 int32
	_ = v601
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v614 int32
	_ = v614
	var v617 int32
	_ = v617
	var v619 int32
	_ = v619
	var v623 int32
	_ = v623
	var v624 int32
	_ = v624
	var v625 int32
	_ = v625
	var v626 int32
	_ = v626
	var v663 int32
	_ = v663
	var v664 int32
	_ = v664
	var v676 int32
	_ = v676
	var v677 int32
	_ = v677
	var v681 int32
	_ = v681
	var v682 int32
	_ = v682
	var v690 int32
	_ = v690
	var v691 int32
	_ = v691
	var v692 int32
	_ = v692
	var v693 int32
	_ = v693
	var v695 int32
	_ = v695
	var v696 int32
	_ = v696
	var v700 int32
	_ = v700
	var v701 int32
	_ = v701
	var v702 int32
	_ = v702
	var v703 int32
	_ = v703
	var v712 int32
	_ = v712
	var v713 int32
	_ = v713
	var v716 int64
	_ = v716
	var v719 int32
	_ = v719
	var v721 int64
	_ = v721
	var v724 int32
	_ = v724
	var v731 int32
	_ = v731
	var v733 int32
	_ = v733
	var v737 int32
	_ = v737
	var v738 int32
	_ = v738
	var v739 int32
	_ = v739
	var v740 int32
	_ = v740
	var v749 int32
	_ = v749
	var v750 int32
	_ = v750
	var v751 int32
	_ = v751
	var v753 int32
	_ = v753
	var v758 int32
	_ = v758
	var v759 int32
	_ = v759
	var v766 int32
	_ = v766
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
	var v815 int32
	_ = v815
	var v819 int32
	_ = v819
	var v825 int32
	_ = v825
	var v828 int32
	_ = v828
	var v834 int32
	_ = v834
	var v839 int32
	_ = v839
	var v845 int32
	_ = v845
	var v848 int32
	_ = v848
	var v851 int32
	_ = v851
	var v853 int32
	_ = v853
	var v857 int32
	_ = v857
	var v860 int32
	_ = v860
	var v861 int32
	_ = v861
	var v863 int32
	_ = v863
	var v869 int32
	_ = v869
	var v880 int32
	_ = v880
	var v902 int32
	_ = v902
	var v905 int32
	_ = v905
	var v906 int32
	_ = v906
	var v914 int32
	_ = v914
	var v915 int32
	_ = v915
	var v950 int32
	_ = v950
	var v951 int32
	_ = v951
	var v956 int64
	_ = v956
	var v959 int32
	_ = v959
	var v960 int32
	_ = v960
	var v972 int32
	_ = v972
	var v992 int32
	_ = v992
	var v997 int32
	_ = v997
	var v1001 int32
	_ = v1001
	var v1008 int32
	_ = v1008
	var v1014 base.V128
	_ = v1014
	var v1017 int32
	_ = v1017
	var v1021 int32
	_ = v1021
	v3 = int32(0)
	v7 = m.G0
	v9 = v7 - int32(32)
	m.G0 = v9
	*(*int32)(unsafe.Add(mBase, uint32(v9)+24)) = v3
	if l1 == v3 {
		v1021 = v3
		goto L1
	} else {
		goto L2
	}
L1:
	;
	m.G0 = v9 + int32(32)
	return v1021
L2:
	;
	if l0 == int32(0) {
		goto L4
	} else {
		goto L5
	}
L3:
	;
	v23 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v24 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v31 = base.Simd_g_const(&F_VP8LEncodeImage__k0)
	base.Simd_g_v128_store(m, v9, int32(0), v31)
	*(*int64)(unsafe.Add(mBase, uint32(v9+int32(16)))) = int64(0)
	v42 = v23*v24<<(uint(base.B2i32(v26 != int32(3)))%32)&int32(-1024) + int32(1024)
	v43 = F_WebPSafeMalloc(m, int64(1), v42)
	mBase = m.M
	if v43 != 0 {
		goto L13
	} else {
		goto L14
	}
L4:
	;
	v20 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v20 != 0 {
		goto L8
	} else {
		goto L9
	}
L5:
	;
	v18 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v18 != 0 {
		goto L3
	} else {
		goto L6
	}
L6:
	;
	goto L4
L7:
	;
	v1021 = int32(0)
	goto L1
L8:
	;
	goto L7
L9:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(3)
	goto L8
L10:
	;
	v997 = *(*int32)(unsafe.Add(mBase, uint32(v9)+20))
	if v997 == int32(0) {
		goto L176
	} else {
		goto L177
	}
L11:
	;
	v62 = F_WebPReportProgress(m, l1, int32(1), v9+int32(24))
	mBase = m.M
	if v62 == int32(0) {
		goto L19
	} else {
		goto L20
	}
L12:
	;
	if v43 != 0 {
		goto L11
	} else {
		goto L15
	}
L13:
	;
	v47 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v47)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v43 + v42
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v43
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v43
	goto L12
L14:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	goto L12
L15:
	;
	v56 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v56 != 0 {
		goto L17
	} else {
		goto L18
	}
L16:
	;
	goto L10
L17:
	;
	goto L16
L18:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L17
L19:
	;
	v992 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v992 != 0 {
		goto L174
	} else {
		goto L175
	}
L20:
	;
	v65 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v65 == int32(0) {
		goto L21
	} else {
		goto L22
	}
L21:
	;
	v196 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v197 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v199 = v197 + int32(-1)
	v200 = int32(14)
	goto L39
L22:
	;
	goto L25
L23:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v65)+20)) = int32(1120272384)
	v192 = base.Simd_g_const(&F_VP8LEncodeImage__k1)
	base.Simd_g_v128_store(m, v65, int32(4), v192)
	goto L21
L25:
	;
	base.MemoryFill(m, v65, int32(0), int32(188))
	goto L23
L37:
	;
	v352 = v196 + int32(-1)
	v353 = int32(14)
	goto L65
L38:
	;
	goto L37
L39:
	;
	v212 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v213 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v213+v200 < int32(32) {
		goto L42
	} else {
		goto L43
	}
L40:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v314 + v312
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v311<<(uint(v314)%32) | v313
	goto L38
L41:
	;
	v231 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v232 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v236 = v228
	v237 = v229
	v238 = v232
	v239 = v231
	goto L45
L42:
	;
	if v213 < int32(16) {
		v311 = v199
		v312 = v200
		v313 = v212
		v314 = v213
		goto L40
	} else {
		goto L44
	}
L43:
	;
	v217 = int32(32)
	v218 = v217 - v213
	v226 = int32(base.Ui32(v199) >> (uint(v218) % 32))
	v227 = v200 - v218
	v228 = v199<<(uint(v213)%32) | v212
	v229 = v217
	goto L41
L44:
	;
	v226 = v199
	v227 = v200
	v228 = v212
	v229 = v213
	goto L41
L45:
	;
	if base.Ui32(v238+int32(2)) <= base.Ui32(v239) {
		v294 = v238
		v295 = v239
		goto L47
	} else {
		goto L48
	}
L46:
	;
	v311 = v226
	v312 = v227
	v313 = v307
	v314 = v305
	goto L40
L47:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v294))) = uint16(v236)
	v302 = v294 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v302
	v305 = v237 + int32(-16)
	v307 = int32(base.Ui32(v236) >> (uint(int32(16)) % 32))
	if int32(31) < v237 {
		v236 = v307
		v237 = v305
		v238 = v302
		v239 = v295
		goto L45
	} else {
		goto L62
	}
L48:
	;
	v248 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v249 = v239 - v248
	v252 = base.I64_extend_i32_s(v249) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v252) {
		v276 = v248
		goto L50
	} else {
		goto L51
	}
L49:
	;
	if v238 == v248 {
		goto L60
	} else {
		goto L61
	}
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v276
	goto L37
L51:
	;
	v255 = v238 - v248
	v257 = v252 + base.I64_extend_i32_u(v255)
	if base.Ui64(int64(4294967295)) < base.Ui64(v257) {
		v276 = v248
		goto L50
	} else {
		goto L52
	}
L52:
	;
	v260 = base.I32_wrap_i64(v257)
	if v239 == v248 {
		goto L53
	} else {
		goto L54
	}
L53:
	;
	v267 = int32(base.Ui32(v249*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v260) < base.Ui32(v267) {
		goto L56
	} else {
		goto L57
	}
L54:
	;
	if base.Ui32(v260) <= base.Ui32(v249) {
		v294 = v238
		v295 = v239
		goto L47
	} else {
		goto L55
	}
L55:
	;
	goto L53
L56:
	;
	v269 = v267
	goto L58
L57:
	;
	v269 = v260
	goto L58
L58:
	;
	v273 = v269&int32(-1024) + int32(1024)
	v274 = F_WebPSafeMalloc(m, int64(1), v273)
	mBase = m.M
	if v274 != 0 {
		goto L49
	} else {
		goto L59
	}
L59:
	;
	v275 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v276 = v275
	goto L50
L60:
	;
	v287 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v287)
	mBase = m.M
	v289 = v274 + v273
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v289
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v274
	v294 = v274 + v255
	v295 = v289
	goto L47
L61:
	;
	v285 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v286 = F_memcpy(m, v274, v285, v255)
	mBase = m.M
	goto L60
L62:
	;
	goto L46
L63:
	;
	v504 = *(*int32)(unsafe.Add(mBase, uint32(v9)+20))
	if v504 == int32(0) {
		goto L89
	} else {
		goto L90
	}
L64:
	;
	goto L63
L65:
	;
	v365 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v366 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v366+v353 < int32(32) {
		goto L68
	} else {
		goto L69
	}
L66:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v467 + v465
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v464<<(uint(v467)%32) | v466
	goto L64
L67:
	;
	v384 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v385 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v389 = v381
	v390 = v382
	v391 = v385
	v392 = v384
	goto L71
L68:
	;
	if v366 < int32(16) {
		v464 = v352
		v465 = v353
		v466 = v365
		v467 = v366
		goto L66
	} else {
		goto L70
	}
L69:
	;
	v370 = int32(32)
	v371 = v370 - v366
	v379 = int32(base.Ui32(v352) >> (uint(v371) % 32))
	v380 = v353 - v371
	v381 = v352<<(uint(v366)%32) | v365
	v382 = v370
	goto L67
L70:
	;
	v379 = v352
	v380 = v353
	v381 = v365
	v382 = v366
	goto L67
L71:
	;
	if base.Ui32(v391+int32(2)) <= base.Ui32(v392) {
		v447 = v391
		v448 = v392
		goto L73
	} else {
		goto L74
	}
L72:
	;
	v464 = v379
	v465 = v380
	v466 = v460
	v467 = v458
	goto L66
L73:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v447))) = uint16(v389)
	v455 = v447 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v455
	v458 = v390 + int32(-16)
	v460 = int32(base.Ui32(v389) >> (uint(int32(16)) % 32))
	if int32(31) < v390 {
		v389 = v460
		v390 = v458
		v391 = v455
		v392 = v448
		goto L71
	} else {
		goto L88
	}
L74:
	;
	v401 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v402 = v392 - v401
	v405 = base.I64_extend_i32_s(v402) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v405) {
		v429 = v401
		goto L76
	} else {
		goto L77
	}
L75:
	;
	if v391 == v401 {
		goto L86
	} else {
		goto L87
	}
L76:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v429
	goto L63
L77:
	;
	v408 = v391 - v401
	v410 = v405 + base.I64_extend_i32_u(v408)
	if base.Ui64(int64(4294967295)) < base.Ui64(v410) {
		v429 = v401
		goto L76
	} else {
		goto L78
	}
L78:
	;
	v413 = base.I32_wrap_i64(v410)
	if v392 == v401 {
		goto L79
	} else {
		goto L80
	}
L79:
	;
	v420 = int32(base.Ui32(v402*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v413) < base.Ui32(v420) {
		goto L82
	} else {
		goto L83
	}
L80:
	;
	if base.Ui32(v413) <= base.Ui32(v402) {
		v447 = v391
		v448 = v392
		goto L73
	} else {
		goto L81
	}
L81:
	;
	goto L79
L82:
	;
	v422 = v420
	goto L84
L83:
	;
	v422 = v413
	goto L84
L84:
	;
	v426 = v422&int32(-1024) + int32(1024)
	v427 = F_WebPSafeMalloc(m, int64(1), v426)
	mBase = m.M
	if v427 != 0 {
		goto L75
	} else {
		goto L85
	}
L85:
	;
	v428 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v429 = v428
	goto L76
L86:
	;
	v440 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v440)
	mBase = m.M
	v442 = v427 + v426
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v442
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v427
	v447 = v427 + v408
	v448 = v442
	goto L73
L87:
	;
	v438 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v439 = F_memcpy(m, v427, v438, v408)
	mBase = m.M
	goto L86
L88:
	;
	goto L72
L89:
	;
	v511 = F_WebPPictureHasTransparency(m, l1)
	mBase = m.M
	v512 = int32(1)
	goto L96
L90:
	;
	v508 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v508 != 0 {
		goto L92
	} else {
		goto L93
	}
L91:
	;
	goto L10
L92:
	;
	goto L91
L93:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L92
L94:
	;
	v663 = int32(0)
	v664 = int32(3)
	goto L122
L95:
	;
	goto L94
L96:
	;
	v524 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v525 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v525+v512 < int32(32) {
		goto L99
	} else {
		goto L100
	}
L97:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v626 + v624
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v623<<(uint(v626)%32) | v625
	goto L95
L98:
	;
	v543 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v544 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v548 = v540
	v549 = v541
	v550 = v544
	v551 = v543
	goto L102
L99:
	;
	if v525 < int32(16) {
		v623 = v511
		v624 = v512
		v625 = v524
		v626 = v525
		goto L97
	} else {
		goto L101
	}
L100:
	;
	v529 = int32(32)
	v530 = v529 - v525
	v538 = int32(base.Ui32(v511) >> (uint(v530) % 32))
	v539 = v512 - v530
	v540 = v511<<(uint(v525)%32) | v524
	v541 = v529
	goto L98
L101:
	;
	v538 = v511
	v539 = v512
	v540 = v524
	v541 = v525
	goto L98
L102:
	;
	if base.Ui32(v550+int32(2)) <= base.Ui32(v551) {
		v606 = v550
		v607 = v551
		goto L104
	} else {
		goto L105
	}
L103:
	;
	v623 = v538
	v624 = v539
	v625 = v619
	v626 = v617
	goto L97
L104:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v606))) = uint16(v548)
	v614 = v606 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v614
	v617 = v549 + int32(-16)
	v619 = int32(base.Ui32(v548) >> (uint(int32(16)) % 32))
	if int32(31) < v549 {
		v548 = v619
		v549 = v617
		v550 = v614
		v551 = v607
		goto L102
	} else {
		goto L119
	}
L105:
	;
	v560 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v561 = v551 - v560
	v564 = base.I64_extend_i32_s(v561) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v564) {
		v588 = v560
		goto L107
	} else {
		goto L108
	}
L106:
	;
	if v550 == v560 {
		goto L117
	} else {
		goto L118
	}
L107:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v588
	goto L94
L108:
	;
	v567 = v550 - v560
	v569 = v564 + base.I64_extend_i32_u(v567)
	if base.Ui64(int64(4294967295)) < base.Ui64(v569) {
		v588 = v560
		goto L107
	} else {
		goto L109
	}
L109:
	;
	v572 = base.I32_wrap_i64(v569)
	if v551 == v560 {
		goto L110
	} else {
		goto L111
	}
L110:
	;
	v579 = int32(base.Ui32(v561*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v572) < base.Ui32(v579) {
		goto L113
	} else {
		goto L114
	}
L111:
	;
	if base.Ui32(v572) <= base.Ui32(v561) {
		v606 = v550
		v607 = v551
		goto L104
	} else {
		goto L112
	}
L112:
	;
	goto L110
L113:
	;
	v581 = v579
	goto L115
L114:
	;
	v581 = v572
	goto L115
L115:
	;
	v585 = v581&int32(-1024) + int32(1024)
	v586 = F_WebPSafeMalloc(m, int64(1), v585)
	mBase = m.M
	if v586 != 0 {
		goto L106
	} else {
		goto L116
	}
L116:
	;
	v587 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v588 = v587
	goto L107
L117:
	;
	v599 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v599)
	mBase = m.M
	v601 = v586 + v585
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v601
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v586
	v606 = v586 + v567
	v607 = v601
	goto L104
L118:
	;
	v597 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v598 = F_memcpy(m, v586, v597, v567)
	mBase = m.M
	goto L117
L119:
	;
	goto L103
L120:
	;
	v815 = *(*int32)(unsafe.Add(mBase, uint32(v9)+20))
	if v815 == int32(0) {
		goto L146
	} else {
		goto L147
	}
L121:
	;
	goto L120
L122:
	;
	v676 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v677 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v677+v664 < int32(32) {
		goto L125
	} else {
		goto L126
	}
L123:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v778 + v776
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v775<<(uint(v778)%32) | v777
	goto L121
L124:
	;
	v695 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v696 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v700 = v692
	v701 = v693
	v702 = v696
	v703 = v695
	goto L128
L125:
	;
	if v677 < int32(16) {
		v775 = v663
		v776 = v664
		v777 = v676
		v778 = v677
		goto L123
	} else {
		goto L127
	}
L126:
	;
	v681 = int32(32)
	v682 = v681 - v677
	v690 = int32(base.Ui32(v663) >> (uint(v682) % 32))
	v691 = v664 - v682
	v692 = v663<<(uint(v677)%32) | v676
	v693 = v681
	goto L124
L127:
	;
	v690 = v663
	v691 = v664
	v692 = v676
	v693 = v677
	goto L124
L128:
	;
	if base.Ui32(v702+int32(2)) <= base.Ui32(v703) {
		v758 = v702
		v759 = v703
		goto L130
	} else {
		goto L131
	}
L129:
	;
	v775 = v690
	v776 = v691
	v777 = v771
	v778 = v769
	goto L123
L130:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v758))) = uint16(v700)
	v766 = v758 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v766
	v769 = v701 + int32(-16)
	v771 = int32(base.Ui32(v700) >> (uint(int32(16)) % 32))
	if int32(31) < v701 {
		v700 = v771
		v701 = v769
		v702 = v766
		v703 = v759
		goto L128
	} else {
		goto L145
	}
L131:
	;
	v712 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v713 = v703 - v712
	v716 = base.I64_extend_i32_s(v713) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v716) {
		v740 = v712
		goto L133
	} else {
		goto L134
	}
L132:
	;
	if v702 == v712 {
		goto L143
	} else {
		goto L144
	}
L133:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v740
	goto L120
L134:
	;
	v719 = v702 - v712
	v721 = v716 + base.I64_extend_i32_u(v719)
	if base.Ui64(int64(4294967295)) < base.Ui64(v721) {
		v740 = v712
		goto L133
	} else {
		goto L135
	}
L135:
	;
	v724 = base.I32_wrap_i64(v721)
	if v703 == v712 {
		goto L136
	} else {
		goto L137
	}
L136:
	;
	v731 = int32(base.Ui32(v713*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v724) < base.Ui32(v731) {
		goto L139
	} else {
		goto L140
	}
L137:
	;
	if base.Ui32(v724) <= base.Ui32(v713) {
		v758 = v702
		v759 = v703
		goto L130
	} else {
		goto L138
	}
L138:
	;
	goto L136
L139:
	;
	v733 = v731
	goto L141
L140:
	;
	v733 = v724
	goto L141
L141:
	;
	v737 = v733&int32(-1024) + int32(1024)
	v738 = F_WebPSafeMalloc(m, int64(1), v737)
	mBase = m.M
	if v738 != 0 {
		goto L132
	} else {
		goto L142
	}
L142:
	;
	v739 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v740 = v739
	goto L133
L143:
	;
	v751 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v751)
	mBase = m.M
	v753 = v738 + v737
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v753
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v738
	v758 = v738 + v719
	v759 = v753
	goto L130
L144:
	;
	v749 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v750 = F_memcpy(m, v738, v749, v719)
	mBase = m.M
	goto L143
L145:
	;
	goto L129
L146:
	;
	v825 = F_WebPReportProgress(m, l1, int32(2), v9+int32(24))
	mBase = m.M
	if v825 == int32(0) {
		goto L19
	} else {
		goto L151
	}
L147:
	;
	v819 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v819 != 0 {
		goto L149
	} else {
		goto L150
	}
L148:
	;
	goto L10
L149:
	;
	goto L148
L150:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L149
L151:
	;
	v828 = F_VP8LEncodeStream(m, l0, l1, v9)
	mBase = m.M
	if v828 == int32(0) {
		goto L10
	} else {
		goto L152
	}
L152:
	;
	v834 = F_WebPReportProgress(m, l1, int32(99), v9+int32(24))
	mBase = m.M
	if v834 == int32(0) {
		goto L19
	} else {
		goto L153
	}
L153:
	;
	v839 = F_WriteImage(m, l1, v9, v9+int32(28))
	mBase = m.M
	if v839 == int32(0) {
		goto L10
	} else {
		goto L154
	}
L154:
	;
	v845 = F_WebPReportProgress(m, l1, int32(100), v9+int32(24))
	mBase = m.M
	if v845 == int32(0) {
		goto L19
	} else {
		goto L155
	}
L155:
	;
	v848 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v848 == int32(0) {
		goto L156
	} else {
		goto L157
	}
L156:
	;
	v857 = *(*int32)(unsafe.Add(mBase, uint32(l1)+84))
	if v857 == int32(0) {
		goto L10
	} else {
		goto L158
	}
L157:
	;
	v851 = *(*int32)(unsafe.Add(mBase, uint32(v9)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v848)+168)) = v851
	v853 = *(*int32)(unsafe.Add(mBase, uint32(v848)))
	*(*int32)(unsafe.Add(mBase, uint32(v848))) = v851 + v853
	goto L156
L158:
	;
	v860 = int32(0)
	v861 = int32(15)
	v863 = int32(4)
	v869 = (v23 + v861) >> (uint(v863) % 32) * ((v24 + v861) >> (uint(v863) % 32))
	if base.Ui32(v869) < base.Ui32(int32(33)) {
		goto L160
	} else {
		goto L161
	}
L159:
	;
	goto L10
L160:
	;
	if v869 == int32(0) {
		goto L162
	} else {
		goto L163
	}
L161:
	;
	base.MemoryFill(m, v857, v860, v869)
	goto L159
L162:
	;
	goto L159
L163:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v857))) = uint8(v860)
	v880 = v857 + v869
	*(*uint8)(unsafe.Add(mBase, uint32(v880+int32(-1)))) = uint8(v860)
	if base.Ui32(v869) < base.Ui32(int32(3)) {
		goto L162
	} else {
		goto L164
	}
L164:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v857)+2)) = uint8(v860)
	*(*uint8)(unsafe.Add(mBase, uint32(v857)+1)) = uint8(v860)
	*(*uint8)(unsafe.Add(mBase, uint32(v880+int32(-3)))) = uint8(v860)
	*(*uint8)(unsafe.Add(mBase, uint32(v880+int32(-2)))) = uint8(v860)
	if base.Ui32(v869) < base.Ui32(int32(7)) {
		goto L162
	} else {
		goto L165
	}
L165:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v857)+3)) = uint8(v860)
	*(*uint8)(unsafe.Add(mBase, uint32(v880+int32(-4)))) = uint8(v860)
	if base.Ui32(v869) < base.Ui32(int32(9)) {
		goto L162
	} else {
		goto L166
	}
L166:
	;
	v902 = int32(0)
	v905 = (v902 - v857) & int32(3)
	v906 = v857 + v905
	*(*int32)(unsafe.Add(mBase, uint32(v906))) = v902
	v914 = (v869 - v905) & int32(60)
	v915 = v906 + v914
	*(*int32)(unsafe.Add(mBase, uint32(v915+int32(-4)))) = v902
	if base.Ui32(v914) < base.Ui32(int32(9)) {
		goto L162
	} else {
		goto L167
	}
L167:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v906)+8)) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v906)+4)) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v915+int32(-8)))) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v915+int32(-12)))) = v902
	if base.Ui32(v914) < base.Ui32(int32(25)) {
		goto L162
	} else {
		goto L168
	}
L168:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v906)+24)) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v906)+20)) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v906)+16)) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v906)+12)) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v915+int32(-16)))) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v915+int32(-20)))) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v915+int32(-24)))) = v902
	*(*int32)(unsafe.Add(mBase, uint32(v915+int32(-28)))) = v902
	v950 = v906&int32(4) | int32(24)
	v951 = v914 - v950
	if base.Ui32(v951) < base.Ui32(int32(32)) {
		goto L162
	} else {
		goto L169
	}
L169:
	;
	v956 = base.I64_extend_i32_u(v902) * int64(4294967297)
	v959 = v951
	v960 = v906 + v950
	goto L170
L170:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v960)+24)) = v956
	*(*int64)(unsafe.Add(mBase, uint32(v960)+16)) = v956
	*(*int64)(unsafe.Add(mBase, uint32(v960)+8)) = v956
	*(*int64)(unsafe.Add(mBase, uint32(v960))) = v956
	v972 = v959 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v972) {
		v959 = v972
		v960 = v960 + int32(32)
		goto L170
	} else {
		goto L172
	}
L171:
	;
	goto L162
L172:
	;
	goto L171
L173:
	;
	goto L10
L174:
	;
	goto L173
L175:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(10)
	goto L174
L176:
	;
	if v9 == int32(0) {
		goto L182
	} else {
		goto L183
	}
L177:
	;
	v1001 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v1001 != 0 {
		goto L179
	} else {
		goto L180
	}
L178:
	;
	goto L176
L179:
	;
	goto L178
L180:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L179
L181:
	;
	v1017 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	v1021 = base.B2i32(v1017 == int32(0))
	goto L1
L182:
	;
	goto L181
L183:
	;
	v1008 = *(*int32)(unsafe.Add(mBase, uint32(v9+int32(8))))
	F_WebPSafeFree(m, v1008)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v9+int32(16)))) = int64(0)
	v1014 = base.Simd_g_const(&F_VP8LEncodeImage__k0)
	base.Simd_g_v128_store(m, v9, int32(0), v1014)
	goto L182
}

var F_VP8LEncodeImage__k0 = [2]uint64{0x0, 0x0}
var F_VP8LEncodeImage__k1 = [2]uint64{0x42c6000042c60000, 0x42c6000042c60000}

func F_VP8LEncodeStream(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v58 int64
	_ = v58
	var v59 int32
	_ = v59
	var v66 int64
	_ = v66
	var v67 int32
	_ = v67
	var v68 int64
	_ = v68
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v94 int32
	_ = v94
	var v98 int32
	_ = v98
	var v101 base.V128
	_ = v101
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v138 int32
	_ = v138
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v161 int32
	_ = v161
	var v164 int32
	_ = v164
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v171 int32
	_ = v171
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v211 int32
	_ = v211
	var v212 int32
	_ = v212
	var v219 int32
	_ = v219
	var v222 int32
	_ = v222
	var v225 int32
	_ = v225
	var v227 int32
	_ = v227
	var v229 int32
	_ = v229
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v247 int32
	_ = v247
	var v255 int32
	_ = v255
	var v259 int32
	_ = v259
	var v270 int32
	_ = v270
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
	var v278 int32
	_ = v278
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
	var v300 int32
	_ = v300
	var v302 int32
	_ = v302
	var v305 int32
	_ = v305
	var v309 int32
	_ = v309
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v316 int32
	_ = v316
	var v319 int32
	_ = v319
	var v320 int32
	_ = v320
	var v323 int32
	_ = v323
	var v327 int32
	_ = v327
	var v328 int32
	_ = v328
	var v332 int32
	_ = v332
	var v333 int32
	_ = v333
	var v340 int32
	_ = v340
	var v341 int32
	_ = v341
	var v347 int32
	_ = v347
	var v353 int32
	_ = v353
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v367 int32
	_ = v367
	var v380 int32
	_ = v380
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v389 int32
	_ = v389
	var v392 int32
	_ = v392
	var v401 int32
	_ = v401
	var v402 int32
	_ = v402
	var v415 int32
	_ = v415
	var v426 int32
	_ = v426
	var v427 int32
	_ = v427
	var v437 int32
	_ = v437
	var v439 int32
	_ = v439
	var v444 int32
	_ = v444
	var v449 int32
	_ = v449
	var v460 int32
	_ = v460
	var v464 int32
	_ = v464
	var v465 int32
	_ = v465
	var v472 int32
	_ = v472
	var v476 int32
	_ = v476
	var v477 int32
	_ = v477
	var v483 int32
	_ = v483
	var v487 int32
	_ = v487
	var v490 int32
	_ = v490
	var v498 int32
	_ = v498
	var v502 int32
	_ = v502
	var v506 int32
	_ = v506
	var v510 int32
	_ = v510
	var v518 int32
	_ = v518
	var v532 int32
	_ = v532
	var v533 int32
	_ = v533
	var v539 int32
	_ = v539
	var v540 int32
	_ = v540
	var v541 int32
	_ = v541
	var v544 int32
	_ = v544
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v549 int32
	_ = v549
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v554 int32
	_ = v554
	var v556 int32
	_ = v556
	var v559 int32
	_ = v559
	var v569 int32
	_ = v569
	var v618 int32
	_ = v618
	var v621 int32
	_ = v621
	var v625 int32
	_ = v625
	var v645 int32
	_ = v645
	var v648 int32
	_ = v648
	var v687 int32
	_ = v687
	var v689 int32
	_ = v689
	var v690 int32
	_ = v690
	var v713 int32
	_ = v713
	var v752 int32
	_ = v752
	var v757 int32
	_ = v757
	var v758 int32
	_ = v758
	var v759 int32
	_ = v759
	var v785 int32
	_ = v785
	var v821 int32
	_ = v821
	var v827 int32
	_ = v827
	var v830 int32
	_ = v830
	var v832 int32
	_ = v832
	var v835 int32
	_ = v835
	var v839 int32
	_ = v839
	var v840 int32
	_ = v840
	var v844 int32
	_ = v844
	var v846 int32
	_ = v846
	var v850 int32
	_ = v850
	var v851 int32
	_ = v851
	var v852 int32
	_ = v852
	var v857 int64
	_ = v857
	var v858 int32
	_ = v858
	var v865 int64
	_ = v865
	var v866 int32
	_ = v866
	var v867 int64
	_ = v867
	var v876 int32
	_ = v876
	var v878 int32
	_ = v878
	var v912 int32
	_ = v912
	var v913 int32
	_ = v913
	var v920 int32
	_ = v920
	var __phi920 int32
	_ = __phi920
	var v930 int32
	_ = v930
	var __phi930 int32
	_ = __phi930
	var v951 int32
	_ = v951
	var __phi951 int32
	_ = __phi951
	var v952 int32
	_ = v952
	var __phi952 int32
	_ = __phi952
	var v974 int32
	_ = v974
	var v984 int32
	_ = v984
	var v985 int32
	_ = v985
	var v1023 int32
	_ = v1023
	var v1026 int32
	_ = v1026
	var v1028 int32
	_ = v1028
	var v1033 int32
	_ = v1033
	var v1035 int32
	_ = v1035
	var v1038 int32
	_ = v1038
	var v1044 int32
	_ = v1044
	var v1046 int32
	_ = v1046
	var v1048 int32
	_ = v1048
	var v1050 int32
	_ = v1050
	var v1051 int32
	_ = v1051
	var v1052 int32
	_ = v1052
	var v1055 int32
	_ = v1055
	var v1059 int32
	_ = v1059
	var v1060 int32
	_ = v1060
	var v1064 int32
	_ = v1064
	var v1068 int32
	_ = v1068
	var v1069 int32
	_ = v1069
	var v1073 int32
	_ = v1073
	var v1075 int32
	_ = v1075
	var v1077 int32
	_ = v1077
	var v1078 int32
	_ = v1078
	var v1086 int32
	_ = v1086
	var v1087 int32
	_ = v1087
	var v1095 int32
	_ = v1095
	var v1096 int32
	_ = v1096
	var v1104 int32
	_ = v1104
	var v1105 int32
	_ = v1105
	var v1113 int32
	_ = v1113
	var v1114 int32
	_ = v1114
	var v1118 int32
	_ = v1118
	var v1120 int32
	_ = v1120
	var v1121 int32
	_ = v1121
	var v1127 int32
	_ = v1127
	var v1128 int32
	_ = v1128
	var v1137 int32
	_ = v1137
	var v1138 int32
	_ = v1138
	var v1145 int32
	_ = v1145
	var v1151 int32
	_ = v1151
	var v1152 int32
	_ = v1152
	var v1161 int32
	_ = v1161
	var v1162 int32
	_ = v1162
	var v1175 int32
	_ = v1175
	var v1176 int32
	_ = v1176
	var v1186 int32
	_ = v1186
	var v1189 int32
	_ = v1189
	var v1244 int32
	_ = v1244
	var v1245 int32
	_ = v1245
	var v1248 int32
	_ = v1248
	var v1252 int32
	_ = v1252
	var v1256 int32
	_ = v1256
	var v1260 int32
	_ = v1260
	var v1264 int32
	_ = v1264
	var v1268 int32
	_ = v1268
	var v1273 int32
	_ = v1273
	var v1275 int32
	_ = v1275
	var v1280 int32
	_ = v1280
	var v1292 int64
	_ = v1292
	var v1295 int64
	_ = v1295
	var v1296 int64
	_ = v1296
	var v1299 int64
	_ = v1299
	var v1302 int64
	_ = v1302
	var v1304 int64
	_ = v1304
	var v1305 int64
	_ = v1305
	var v1308 int32
	_ = v1308
	var v1311 int32
	_ = v1311
	var v1320 int64
	_ = v1320
	var v1323 int64
	_ = v1323
	var v1326 int64
	_ = v1326
	var v1328 int64
	_ = v1328
	var v1329 int64
	_ = v1329
	var v1340 int32
	_ = v1340
	var v1342 int32
	_ = v1342
	var v1347 int32
	_ = v1347
	var v1359 int64
	_ = v1359
	var v1362 int64
	_ = v1362
	var v1363 int64
	_ = v1363
	var v1366 int64
	_ = v1366
	var v1369 int64
	_ = v1369
	var v1371 int64
	_ = v1371
	var v1372 int64
	_ = v1372
	var v1375 int32
	_ = v1375
	var v1378 int32
	_ = v1378
	var v1387 int64
	_ = v1387
	var v1390 int64
	_ = v1390
	var v1393 int64
	_ = v1393
	var v1395 int64
	_ = v1395
	var v1396 int64
	_ = v1396
	var v1407 int32
	_ = v1407
	var v1409 int32
	_ = v1409
	var v1414 int32
	_ = v1414
	var v1426 int64
	_ = v1426
	var v1429 int64
	_ = v1429
	var v1430 int64
	_ = v1430
	var v1433 int64
	_ = v1433
	var v1436 int64
	_ = v1436
	var v1438 int64
	_ = v1438
	var v1439 int64
	_ = v1439
	var v1442 int32
	_ = v1442
	var v1445 int32
	_ = v1445
	var v1454 int64
	_ = v1454
	var v1457 int64
	_ = v1457
	var v1460 int64
	_ = v1460
	var v1462 int64
	_ = v1462
	var v1463 int64
	_ = v1463
	var v1474 int32
	_ = v1474
	var v1476 int32
	_ = v1476
	var v1481 int32
	_ = v1481
	var v1493 int64
	_ = v1493
	var v1496 int64
	_ = v1496
	var v1497 int64
	_ = v1497
	var v1500 int64
	_ = v1500
	var v1503 int64
	_ = v1503
	var v1505 int64
	_ = v1505
	var v1506 int64
	_ = v1506
	var v1509 int32
	_ = v1509
	var v1512 int32
	_ = v1512
	var v1521 int64
	_ = v1521
	var v1524 int64
	_ = v1524
	var v1527 int64
	_ = v1527
	var v1529 int64
	_ = v1529
	var v1530 int64
	_ = v1530
	var v1541 int32
	_ = v1541
	var v1543 int32
	_ = v1543
	var v1548 int32
	_ = v1548
	var v1560 int64
	_ = v1560
	var v1563 int64
	_ = v1563
	var v1564 int64
	_ = v1564
	var v1567 int64
	_ = v1567
	var v1570 int64
	_ = v1570
	var v1572 int64
	_ = v1572
	var v1573 int64
	_ = v1573
	var v1576 int32
	_ = v1576
	var v1579 int32
	_ = v1579
	var v1588 int64
	_ = v1588
	var v1591 int64
	_ = v1591
	var v1594 int64
	_ = v1594
	var v1596 int64
	_ = v1596
	var v1597 int64
	_ = v1597
	var v1608 int32
	_ = v1608
	var v1610 int32
	_ = v1610
	var v1615 int32
	_ = v1615
	var v1627 int64
	_ = v1627
	var v1630 int64
	_ = v1630
	var v1631 int64
	_ = v1631
	var v1634 int64
	_ = v1634
	var v1637 int64
	_ = v1637
	var v1639 int64
	_ = v1639
	var v1640 int64
	_ = v1640
	var v1643 int32
	_ = v1643
	var v1646 int32
	_ = v1646
	var v1655 int64
	_ = v1655
	var v1658 int64
	_ = v1658
	var v1661 int64
	_ = v1661
	var v1663 int64
	_ = v1663
	var v1664 int64
	_ = v1664
	var v1675 int32
	_ = v1675
	var v1677 int32
	_ = v1677
	var v1682 int32
	_ = v1682
	var v1694 int64
	_ = v1694
	var v1697 int64
	_ = v1697
	var v1698 int64
	_ = v1698
	var v1701 int64
	_ = v1701
	var v1704 int64
	_ = v1704
	var v1706 int64
	_ = v1706
	var v1707 int64
	_ = v1707
	var v1710 int32
	_ = v1710
	var v1713 int32
	_ = v1713
	var v1722 int64
	_ = v1722
	var v1725 int64
	_ = v1725
	var v1728 int64
	_ = v1728
	var v1730 int64
	_ = v1730
	var v1731 int64
	_ = v1731
	var v1742 int32
	_ = v1742
	var v1744 int32
	_ = v1744
	var v1749 int32
	_ = v1749
	var v1761 int64
	_ = v1761
	var v1764 int64
	_ = v1764
	var v1765 int64
	_ = v1765
	var v1768 int64
	_ = v1768
	var v1771 int64
	_ = v1771
	var v1773 int64
	_ = v1773
	var v1774 int64
	_ = v1774
	var v1777 int32
	_ = v1777
	var v1780 int32
	_ = v1780
	var v1789 int64
	_ = v1789
	var v1792 int64
	_ = v1792
	var v1795 int64
	_ = v1795
	var v1797 int64
	_ = v1797
	var v1798 int64
	_ = v1798
	var v1809 int32
	_ = v1809
	var v1811 int32
	_ = v1811
	var v1816 int32
	_ = v1816
	var v1828 int64
	_ = v1828
	var v1831 int64
	_ = v1831
	var v1832 int64
	_ = v1832
	var v1835 int64
	_ = v1835
	var v1838 int64
	_ = v1838
	var v1840 int64
	_ = v1840
	var v1841 int64
	_ = v1841
	var v1844 int32
	_ = v1844
	var v1847 int32
	_ = v1847
	var v1856 int64
	_ = v1856
	var v1859 int64
	_ = v1859
	var v1862 int64
	_ = v1862
	var v1864 int64
	_ = v1864
	var v1865 int64
	_ = v1865
	var v1876 int32
	_ = v1876
	var v1878 int32
	_ = v1878
	var v1883 int32
	_ = v1883
	var v1895 int64
	_ = v1895
	var v1898 int64
	_ = v1898
	var v1899 int64
	_ = v1899
	var v1902 int64
	_ = v1902
	var v1905 int64
	_ = v1905
	var v1907 int64
	_ = v1907
	var v1908 int64
	_ = v1908
	var v1911 int32
	_ = v1911
	var v1914 int32
	_ = v1914
	var v1923 int64
	_ = v1923
	var v1926 int64
	_ = v1926
	var v1929 int64
	_ = v1929
	var v1931 int64
	_ = v1931
	var v1932 int64
	_ = v1932
	var v1943 int32
	_ = v1943
	var v1945 int32
	_ = v1945
	var v1950 int32
	_ = v1950
	var v1962 int64
	_ = v1962
	var v1965 int64
	_ = v1965
	var v1966 int64
	_ = v1966
	var v1969 int64
	_ = v1969
	var v1972 int64
	_ = v1972
	var v1974 int64
	_ = v1974
	var v1975 int64
	_ = v1975
	var v1978 int32
	_ = v1978
	var v1981 int32
	_ = v1981
	var v1990 int64
	_ = v1990
	var v1993 int64
	_ = v1993
	var v1996 int64
	_ = v1996
	var v1998 int64
	_ = v1998
	var v1999 int64
	_ = v1999
	var v2010 int32
	_ = v2010
	var v2012 int32
	_ = v2012
	var v2017 int32
	_ = v2017
	var v2029 int64
	_ = v2029
	var v2032 int64
	_ = v2032
	var v2033 int64
	_ = v2033
	var v2036 int64
	_ = v2036
	var v2039 int64
	_ = v2039
	var v2041 int64
	_ = v2041
	var v2042 int64
	_ = v2042
	var v2045 int32
	_ = v2045
	var v2048 int32
	_ = v2048
	var v2057 int64
	_ = v2057
	var v2060 int64
	_ = v2060
	var v2063 int64
	_ = v2063
	var v2065 int64
	_ = v2065
	var v2066 int64
	_ = v2066
	var v2077 int32
	_ = v2077
	var v2079 int32
	_ = v2079
	var v2084 int32
	_ = v2084
	var v2096 int64
	_ = v2096
	var v2099 int64
	_ = v2099
	var v2100 int64
	_ = v2100
	var v2103 int64
	_ = v2103
	var v2106 int64
	_ = v2106
	var v2108 int64
	_ = v2108
	var v2109 int64
	_ = v2109
	var v2112 int32
	_ = v2112
	var v2115 int32
	_ = v2115
	var v2124 int64
	_ = v2124
	var v2127 int64
	_ = v2127
	var v2130 int64
	_ = v2130
	var v2132 int64
	_ = v2132
	var v2133 int64
	_ = v2133
	var v2138 int64
	_ = v2138
	var v2142 int32
	_ = v2142
	var v2144 int32
	_ = v2144
	var v2153 int64
	_ = v2153
	var v2154 int64
	_ = v2154
	var v2156 int64
	_ = v2156
	var v2165 int64
	_ = v2165
	var v2167 int64
	_ = v2167
	var v2169 int64
	_ = v2169
	var v2171 int64
	_ = v2171
	var v2175 int32
	_ = v2175
	var v2176 int32
	_ = v2176
	var v2179 int64
	_ = v2179
	var v2181 int64
	_ = v2181
	var v2186 int32
	_ = v2186
	var v2192 int64
	_ = v2192
	var v2194 int32
	_ = v2194
	var v2203 int64
	_ = v2203
	var v2204 int64
	_ = v2204
	var v2206 int32
	_ = v2206
	var v2207 int32
	_ = v2207
	var v2208 int32
	_ = v2208
	var v2213 int32
	_ = v2213
	var v2214 int32
	_ = v2214
	var v2215 int32
	_ = v2215
	var v2217 int32
	_ = v2217
	var v2226 int32
	_ = v2226
	var v2236 int32
	_ = v2236
	var v2238 int32
	_ = v2238
	var v2274 int32
	_ = v2274
	var v2275 int32
	_ = v2275
	var v2276 int32
	_ = v2276
	var v2278 int32
	_ = v2278
	var v2279 int32
	_ = v2279
	var v2282 int32
	_ = v2282
	var v2284 int32
	_ = v2284
	var v2286 int32
	_ = v2286
	var v2289 int32
	_ = v2289
	var v2291 int32
	_ = v2291
	var v2293 int32
	_ = v2293
	var v2296 int32
	_ = v2296
	var v2298 int32
	_ = v2298
	var v2303 int32
	_ = v2303
	var v2308 int32
	_ = v2308
	var v2329 int32
	_ = v2329
	var v2332 int32
	_ = v2332
	var v2363 int32
	_ = v2363
	var v2366 int32
	_ = v2366
	var v2369 int32
	_ = v2369
	var v2370 float32
	_ = v2370
	var v2400 int32
	_ = v2400
	var v2405 int32
	_ = v2405
	var v2407 int32
	_ = v2407
	var v2414 int32
	_ = v2414
	var v2417 int32
	_ = v2417
	var v2424 int32
	_ = v2424
	var v2425 int32
	_ = v2425
	var v2426 int32
	_ = v2426
	var v2427 int32
	_ = v2427
	var v2433 int32
	_ = v2433
	var v2434 int32
	_ = v2434
	var v2438 int32
	_ = v2438
	var v2439 int32
	_ = v2439
	var v2440 int32
	_ = v2440
	var v2444 int32
	_ = v2444
	var v2447 int32
	_ = v2447
	var v2448 int32
	_ = v2448
	var v2463 int32
	_ = v2463
	var v2473 int32
	_ = v2473
	var v2513 int32
	_ = v2513
	var v2524 int32
	_ = v2524
	var v2537 int32
	_ = v2537
	var v2577 int32
	_ = v2577
	var v2583 int32
	_ = v2583
	var v2601 int32
	_ = v2601
	var v2610 int32
	_ = v2610
	var v2667 int32
	_ = v2667
	var v2680 int32
	_ = v2680
	var v2715 int32
	_ = v2715
	var v2719 int32
	_ = v2719
	var v2726 int32
	_ = v2726
	var v2784 int32
	_ = v2784
	var v2801 int32
	_ = v2801
	var v2807 int32
	_ = v2807
	var v2808 int32
	_ = v2808
	var v2809 int32
	_ = v2809
	var v2842 int32
	_ = v2842
	var v2853 int32
	_ = v2853
	var v2856 int32
	_ = v2856
	var v2865 int32
	_ = v2865
	var v2868 int32
	_ = v2868
	var v2906 int32
	_ = v2906
	var v2907 int32
	_ = v2907
	var v2908 int32
	_ = v2908
	var v2939 int32
	_ = v2939
	var v2940 int32
	_ = v2940
	var v2941 int32
	_ = v2941
	var v2942 int32
	_ = v2942
	var v2943 int32
	_ = v2943
	var v2947 int32
	_ = v2947
	var v3007 int32
	_ = v3007
	var v3012 int32
	_ = v3012
	var v3016 int32
	_ = v3016
	var v3018 int32
	_ = v3018
	var v3019 base.V128
	_ = v3019
	var v3029 int32
	_ = v3029
	var v3032 int32
	_ = v3032
	var v3035 int32
	_ = v3035
	var v3036 base.V128
	_ = v3036
	var v3046 int32
	_ = v3046
	var v3049 int32
	_ = v3049
	var v3052 int32
	_ = v3052
	var v3053 base.V128
	_ = v3053
	var v3063 int32
	_ = v3063
	var v3066 int32
	_ = v3066
	var v3069 int32
	_ = v3069
	var v3070 base.V128
	_ = v3070
	var v3080 int32
	_ = v3080
	var v3083 int32
	_ = v3083
	var v3085 int32
	_ = v3085
	var v3090 int32
	_ = v3090
	var v3101 int32
	_ = v3101
	var v3104 int32
	_ = v3104
	var v3107 int32
	_ = v3107
	var v3109 int32
	_ = v3109
	var v3110 int32
	_ = v3110
	var v3122 int32
	_ = v3122
	var v3128 int32
	_ = v3128
	var v3133 int32
	_ = v3133
	var v3135 int32
	_ = v3135
	var v3138 int32
	_ = v3138
	var v3149 int32
	_ = v3149
	var v3152 int32
	_ = v3152
	var v3153 int32
	_ = v3153
	var v3155 int32
	_ = v3155
	var v3166 int32
	_ = v3166
	var v3167 int32
	_ = v3167
	var v3170 int32
	_ = v3170
	var v3179 int32
	_ = v3179
	var v3182 int32
	_ = v3182
	var v3186 int32
	_ = v3186
	var v3188 int32
	_ = v3188
	var v3189 int32
	_ = v3189
	var v3192 int32
	_ = v3192
	var v3194 int32
	_ = v3194
	var v3196 int32
	_ = v3196
	var v3201 int32
	_ = v3201
	var v3202 int32
	_ = v3202
	var v3205 int32
	_ = v3205
	var v3208 int32
	_ = v3208
	var v3211 int32
	_ = v3211
	var v3215 int32
	_ = v3215
	var v3218 int32
	_ = v3218
	var v3224 int32
	_ = v3224
	var v3225 int32
	_ = v3225
	var v3227 int32
	_ = v3227
	var v3235 int32
	_ = v3235
	var v3236 int32
	_ = v3236
	var v3237 int32
	_ = v3237
	var v3240 int32
	_ = v3240
	var v3253 int32
	_ = v3253
	var v3258 int32
	_ = v3258
	var v3259 int32
	_ = v3259
	var v3265 int32
	_ = v3265
	var v3275 int32
	_ = v3275
	var v3276 int32
	_ = v3276
	var v3277 int32
	_ = v3277
	var v3279 int32
	_ = v3279
	var v3280 int32
	_ = v3280
	var v3281 int32
	_ = v3281
	var v3283 int64
	_ = v3283
	var v3289 int32
	_ = v3289
	var v3290 int32
	_ = v3290
	var v3291 int32
	_ = v3291
	var v3298 int32
	_ = v3298
	var v3300 int32
	_ = v3300
	var v3304 int32
	_ = v3304
	var v3305 int32
	_ = v3305
	var v3310 int32
	_ = v3310
	var v3311 int32
	_ = v3311
	var v3312 int32
	_ = v3312
	var v3319 int32
	_ = v3319
	var v3320 int32
	_ = v3320
	var v3321 int32
	_ = v3321
	var v3322 int32
	_ = v3322
	var v3323 int64
	_ = v3323
	var v3325 int32
	_ = v3325
	var v3327 int32
	_ = v3327
	var v3335 int32
	_ = v3335
	var v3337 int32
	_ = v3337
	var v3343 int64
	_ = v3343
	var v3344 int32
	_ = v3344
	var v3351 int64
	_ = v3351
	var v3352 int32
	_ = v3352
	var v3353 int64
	_ = v3353
	var v3362 int32
	_ = v3362
	var v3364 int32
	_ = v3364
	var v3367 int32
	_ = v3367
	var v3369 int32
	_ = v3369
	var v3380 int32
	_ = v3380
	var v3381 int32
	_ = v3381
	var v3382 int32
	_ = v3382
	var v3383 int32
	_ = v3383
	var v3384 int32
	_ = v3384
	var v3388 int32
	_ = v3388
	var v3396 int32
	_ = v3396
	var v3400 int32
	_ = v3400
	var v3404 int32
	_ = v3404
	var v3406 int32
	_ = v3406
	var v3407 base.V128
	_ = v3407
	var v3417 int32
	_ = v3417
	var v3420 int32
	_ = v3420
	var v3423 int32
	_ = v3423
	var v3424 base.V128
	_ = v3424
	var v3434 int32
	_ = v3434
	var v3437 int32
	_ = v3437
	var v3440 int32
	_ = v3440
	var v3441 base.V128
	_ = v3441
	var v3451 int32
	_ = v3451
	var v3454 int32
	_ = v3454
	var v3457 int32
	_ = v3457
	var v3458 base.V128
	_ = v3458
	var v3468 int32
	_ = v3468
	var v3471 int32
	_ = v3471
	var v3473 int64
	_ = v3473
	var v3475 int32
	_ = v3475
	var v3477 int32
	_ = v3477
	var v3479 int32
	_ = v3479
	var v3483 int32
	_ = v3483
	var v3484 int32
	_ = v3484
	var v3488 int32
	_ = v3488
	var v3492 int32
	_ = v3492
	var v3496 int32
	_ = v3496
	var v3503 int32
	_ = v3503
	var v3509 int32
	_ = v3509
	var v3510 int32
	_ = v3510
	var v3512 int32
	_ = v3512
	var v3515 int32
	_ = v3515
	var v3521 int32
	_ = v3521
	var v3524 int32
	_ = v3524
	var v3528 int32
	_ = v3528
	var v3529 int32
	_ = v3529
	var v3533 int32
	_ = v3533
	var v3534 int32
	_ = v3534
	var v3537 int32
	_ = v3537
	var v3542 int32
	_ = v3542
	var v3543 int32
	_ = v3543
	var v3544 int32
	_ = v3544
	var v3547 int32
	_ = v3547
	var v3551 int32
	_ = v3551
	var v3552 int32
	_ = v3552
	var v3553 int32
	_ = v3553
	var v3556 int32
	_ = v3556
	var v3557 int32
	_ = v3557
	var v3559 int32
	_ = v3559
	var v3560 int32
	_ = v3560
	var v3562 int32
	_ = v3562
	var v3565 int32
	_ = v3565
	var v3566 int32
	_ = v3566
	var v3568 int32
	_ = v3568
	var v3576 int32
	_ = v3576
	var v3579 int32
	_ = v3579
	var v3583 int64
	_ = v3583
	var v3584 int64
	_ = v3584
	var v3587 int32
	_ = v3587
	var v3589 base.V128
	_ = v3589
	var v3591 int32
	_ = v3591
	var v3593 base.V128
	_ = v3593
	var v3599 int32
	_ = v3599
	var v3605 int32
	_ = v3605
	var v3611 int32
	_ = v3611
	var v3660 int32
	_ = v3660
	var v3665 int32
	_ = v3665
	var v3671 base.V128
	_ = v3671
	var v3675 int32
	_ = v3675
	var v3676 int32
	_ = v3676
	var v3681 int32
	_ = v3681
	var v3684 int32
	_ = v3684
	var v3687 int32
	_ = v3687
	var v3689 int32
	_ = v3689
	var v3691 int32
	_ = v3691
	var v3702 int32
	_ = v3702
	var v3703 int32
	_ = v3703
	var v3710 int32
	_ = v3710
	var v3713 int32
	_ = v3713
	var v3716 int32
	_ = v3716
	var v3718 int32
	_ = v3718
	var v3720 int32
	_ = v3720
	var v3731 int32
	_ = v3731
	var v3732 int32
	_ = v3732
	var v3739 int32
	_ = v3739
	var v3742 int32
	_ = v3742
	var v3745 int32
	_ = v3745
	var v3747 int32
	_ = v3747
	var v3749 int32
	_ = v3749
	var v3760 int32
	_ = v3760
	var v3761 int32
	_ = v3761
	var v3768 int32
	_ = v3768
	var v3771 int32
	_ = v3771
	var v3774 int32
	_ = v3774
	var v3776 int32
	_ = v3776
	var v3778 int32
	_ = v3778
	var v3789 int32
	_ = v3789
	var v3790 int32
	_ = v3790
	var v3796 int32
	_ = v3796
	var v3804 int32
	_ = v3804
	var v3805 int32
	_ = v3805
	var v3810 int32
	_ = v3810
	var v3813 int32
	_ = v3813
	var v3816 int32
	_ = v3816
	var v3818 int32
	_ = v3818
	var v3820 int32
	_ = v3820
	var v3831 int32
	_ = v3831
	var v3832 int32
	_ = v3832
	var v3839 int32
	_ = v3839
	var v3842 int32
	_ = v3842
	var v3845 int32
	_ = v3845
	var v3847 int32
	_ = v3847
	var v3849 int32
	_ = v3849
	var v3860 int32
	_ = v3860
	var v3861 int32
	_ = v3861
	var v3868 int32
	_ = v3868
	var v3871 int32
	_ = v3871
	var v3874 int32
	_ = v3874
	var v3876 int32
	_ = v3876
	var v3878 int32
	_ = v3878
	var v3889 int32
	_ = v3889
	var v3890 int32
	_ = v3890
	var v3897 int32
	_ = v3897
	var v3900 int32
	_ = v3900
	var v3903 int32
	_ = v3903
	var v3905 int32
	_ = v3905
	var v3907 int32
	_ = v3907
	var v3918 int32
	_ = v3918
	var v3919 int32
	_ = v3919
	var v3925 int32
	_ = v3925
	var v3930 int32
	_ = v3930
	var v3938 int32
	_ = v3938
	v54 = m.G0
	v56 = v54 - int32(1728)
	m.G0 = v56
	v58 = int64(1)
	v59 = int32(2224)
	goto L8
L1:
	;
	m.G0 = v56 + int32(1728)
	return v3938
L2:
	;
	v259 = v56 + int32(8)
	goto L62
L3:
	;
	v255 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v255 != 0 {
		goto L57
	} else {
		goto L58
	}
L4:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+12)) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v79)+4)) = l1
	*(*int32)(unsafe.Add(mBase, uint32(v79))) = l0
	F_VP8LEncDspInit(m)
	mBase = m.M
	v94 = m.G1
	goto L16
L5:
	;
	if v79 != 0 {
		goto L4
	} else {
		goto L11
	}
L6:
	;
	goto L5
L7:
	;
	v77 = F_calloc(m, base.I32_wrap_i64(v58), v59)
	mBase = m.M
	v79 = v77
	goto L6
L8:
	;
	v66 = base.I64_div_u_s(int64(2147418112), v58)
	v67 = int32(0)
	v68 = base.I64_extend_i32_u(v59)
	if base.Ui64(int64(4294967295)) < base.Ui64(v68*v58) {
		v79 = v67
		goto L6
	} else {
		goto L9
	}
L9:
	;
	if base.Ui64(v66) < base.Ui64(v68) {
		v79 = v67
		goto L6
	} else {
		goto L10
	}
L10:
	;
	goto L7
L11:
	;
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v82 != 0 {
		goto L13
	} else {
		goto L14
	}
L12:
	;
	goto L15
L13:
	;
	goto L12
L14:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L13
L15:
	;
	goto L3
L16:
	;
	v98 = v56 + int32(180)
	v101 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v98, int32(0), v101)
	*(*int64)(unsafe.Add(mBase, uint32(v56+int32(196)))) = int64(0)
	v111 = int32(1024)
	v113 = F_WebPSafeMalloc(m, int64(1), v111)
	mBase = m.M
	if v113 != 0 {
		goto L18
	} else {
		goto L19
	}
L17:
	;
	if v113 != 0 {
		goto L2
	} else {
		goto L20
	}
L18:
	;
	v117 = *(*int32)(unsafe.Add(mBase, uint32(v98)+8))
	F_WebPSafeFree(m, v117)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v98)+16)) = v113 + v111
	*(*int32)(unsafe.Add(mBase, uint32(v98)+12)) = v113
	*(*int32)(unsafe.Add(mBase, uint32(v98)+8)) = v113
	goto L17
L19:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v98)+20)) = int32(1)
	goto L17
L20:
	;
	v126 = v79 + int32(2216)
	v127 = *(*int32)(unsafe.Add(mBase, uint32(v126)))
	F_WebPSafeFree(m, v127)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v126))) = int64(0)
	goto L21
L21:
	;
	v132 = v79 + int32(2120)
	v135 = *(*int32)(unsafe.Add(mBase, uint32(v132)+12))
	if v135 == int32(0) {
		goto L23
	} else {
		goto L24
	}
L22:
	;
	v161 = v79 + int32(2144)
	v164 = *(*int32)(unsafe.Add(mBase, uint32(v161)+12))
	if v164 == int32(0) {
		goto L31
	} else {
		goto L32
	}
L23:
	;
	v140 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v132)+20)) = v140
	v142 = *(*int32)(unsafe.Add(mBase, uint32(v132)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v132)+16)) = v142
	*(*int32)(unsafe.Add(mBase, uint32(v132)+12)) = v79 + int32(2128)
	*(*int32)(unsafe.Add(mBase, uint32(v132)+8)) = v140
	if v142 == v140 {
		goto L25
	} else {
		goto L26
	}
L24:
	;
	v138 = *(*int32)(unsafe.Add(mBase, uint32(v132)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v135))) = v138
	goto L23
L25:
	;
	goto L22
L26:
	;
	v153 = v142
	goto L27
L27:
	;
	v154 = *(*int32)(unsafe.Add(mBase, uint32(v153)))
	F_WebPSafeFree(m, v153)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v132)+16)) = v154
	if v154 != 0 {
		v153 = v154
		goto L27
	} else {
		goto L29
	}
L28:
	;
	goto L25
L29:
	;
	goto L28
L30:
	;
	v190 = v79 + int32(2168)
	v193 = *(*int32)(unsafe.Add(mBase, uint32(v190)+12))
	if v193 == int32(0) {
		goto L39
	} else {
		goto L40
	}
L31:
	;
	v169 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v161)+20)) = v169
	v171 = *(*int32)(unsafe.Add(mBase, uint32(v161)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v161)+16)) = v171
	*(*int32)(unsafe.Add(mBase, uint32(v161)+12)) = v79 + int32(2152)
	*(*int32)(unsafe.Add(mBase, uint32(v161)+8)) = v169
	if v171 == v169 {
		goto L33
	} else {
		goto L34
	}
L32:
	;
	v167 = *(*int32)(unsafe.Add(mBase, uint32(v161)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v164))) = v167
	goto L31
L33:
	;
	goto L30
L34:
	;
	v182 = v171
	goto L35
L35:
	;
	v183 = *(*int32)(unsafe.Add(mBase, uint32(v182)))
	F_WebPSafeFree(m, v182)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v161)+16)) = v183
	if v183 != 0 {
		v182 = v183
		goto L35
	} else {
		goto L37
	}
L36:
	;
	goto L33
L37:
	;
	goto L36
L38:
	;
	v219 = v79 + int32(2192)
	v222 = *(*int32)(unsafe.Add(mBase, uint32(v219)+12))
	if v222 == int32(0) {
		goto L47
	} else {
		goto L48
	}
L39:
	;
	v198 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v190)+20)) = v198
	v200 = *(*int32)(unsafe.Add(mBase, uint32(v190)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v190)+16)) = v200
	*(*int32)(unsafe.Add(mBase, uint32(v190)+12)) = v79 + int32(2176)
	*(*int32)(unsafe.Add(mBase, uint32(v190)+8)) = v198
	if v200 == v198 {
		goto L41
	} else {
		goto L42
	}
L40:
	;
	v196 = *(*int32)(unsafe.Add(mBase, uint32(v190)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v193))) = v196
	goto L39
L41:
	;
	goto L38
L42:
	;
	v211 = v200
	goto L43
L43:
	;
	v212 = *(*int32)(unsafe.Add(mBase, uint32(v211)))
	F_WebPSafeFree(m, v211)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v190)+16)) = v212
	if v212 != 0 {
		v211 = v212
		goto L43
	} else {
		goto L45
	}
L44:
	;
	goto L41
L45:
	;
	goto L44
L46:
	;
	v247 = *(*int32)(unsafe.Add(mBase, uint32(v79)+24))
	F_free(m, v247)
	mBase = m.M
	goto L54
L47:
	;
	v227 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v219)+20)) = v227
	v229 = *(*int32)(unsafe.Add(mBase, uint32(v219)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v219)+16)) = v229
	*(*int32)(unsafe.Add(mBase, uint32(v219)+12)) = v79 + int32(2200)
	*(*int32)(unsafe.Add(mBase, uint32(v219)+8)) = v227
	if v229 == v227 {
		goto L49
	} else {
		goto L50
	}
L48:
	;
	v225 = *(*int32)(unsafe.Add(mBase, uint32(v219)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v222))) = v225
	goto L47
L49:
	;
	goto L46
L50:
	;
	v240 = v229
	goto L51
L51:
	;
	v241 = *(*int32)(unsafe.Add(mBase, uint32(v240)))
	F_WebPSafeFree(m, v240)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v219)+16)) = v241
	if v241 != 0 {
		v240 = v241
		goto L51
	} else {
		goto L53
	}
L52:
	;
	goto L49
L53:
	;
	goto L52
L54:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v79)+24)) = int64(0)
	F_free(m, v79)
	mBase = m.M
	goto L55
L55:
	;
	goto L3
L56:
	;
	v3938 = int32(0)
	goto L1
L57:
	;
	goto L56
L58:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L57
L59:
	;
	v3660 = v56 + int32(180)
	if v3660 == int32(0) {
		goto L589
	} else {
		goto L590
	}
L60:
	;
	goto L64
L61:
	;
	goto L60
L62:
	;
	if v259 == int32(0) {
		goto L61
	} else {
		goto L63
	}
L63:
	;
	v270 = int32(0)
	v272 = F_memset(m, v259, v270, int32(172))
	mBase = m.M
	v273 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v272)+72)) = v273 + int32(1)
	v278 = F_WebPEncodingSetError(m, v272, v270)
	mBase = m.M
	goto L61
L64:
	;
	v283 = *(*int32)(unsafe.Add(mBase, uint32(v79)+4))
	v284 = *(*int32)(unsafe.Add(mBase, uint32(v283)+12))
	v285 = *(*int32)(unsafe.Add(mBase, uint32(v283)+8))
	v286 = *(*int32)(unsafe.Add(mBase, uint32(v79)))
	v287 = *(*int32)(unsafe.Add(mBase, uint32(v286)+8))
	v289 = v79 + int32(1096)
	v290 = int32(0)
	v300 = m.G0
	v302 = v300 - int32(_a_F_VP8LEncodeStream_0)
	m.G0 = v302
	v305 = int32(_a_F_VP8LEncodeStream_1)
	v309 = F_memset(m, v302+v305, v290, int32(1024))
	mBase = m.M
	v312 = F_memset(m, v302, v290, v305)
	mBase = m.M
	v313 = *(*int32)(unsafe.Add(mBase, uint32(v283)+12))
	if v313 < int32(1) {
		v449 = v290
		goto L68
	} else {
		goto L69
	}
L65:
	;
	v532 = base.B2i32(v518 < int32(257))
	if v518 < int32(257) {
		goto L95
	} else {
		goto L96
	}
L66:
	;
	m.G0 = v312 + int32(_a_F_VP8LEncodeStream_0)
	goto L65
L67:
	;
	v518 = int32(257)
	goto L66
L68:
	;
	if v289 == int32(0) {
		v518 = v449
		goto L66
	} else {
		goto L87
	}
L69:
	;
	v316 = *(*int32)(unsafe.Add(mBase, uint32(v283)+8))
	if v316 < int32(1) {
		v449 = v290
		goto L68
	} else {
		goto L70
	}
L70:
	;
	v319 = *(*int32)(unsafe.Add(mBase, uint32(v283)+52))
	v320 = *(*int32)(unsafe.Add(mBase, uint32(v319)))
	v323 = int32(0)
	v327 = v320 ^ int32(-1)
	v328 = v323
	v332 = v319
	v333 = v323
	goto L71
L71:
	;
	v340 = v327
	v341 = v328
	v347 = int32(0)
	goto L73
L72:
	;
	v449 = v427
	goto L68
L73:
	;
	v353 = *(*int32)(unsafe.Add(mBase, uint32(v332+v347<<(uint(int32(2))%32))))
	if v353 == v340 {
		v426 = v340
		v427 = v341
		goto L75
	} else {
		goto L76
	}
L74:
	;
	v439 = *(*int32)(unsafe.Add(mBase, uint32(v283)+56))
	v444 = v333 + int32(1)
	if v444 != v313 {
		v327 = v426
		v328 = v427
		v332 = v332 + v439<<(uint(int32(2))%32)
		v333 = v444
		goto L71
	} else {
		goto L86
	}
L75:
	;
	v437 = v347 + int32(1)
	if v437 != v316 {
		v340 = v426
		v341 = v427
		v347 = v437
		goto L73
	} else {
		goto L85
	}
L76:
	;
	v360 = int32(base.Ui32(v353*int32(506832829)) >> (uint(int32(22)) % 32))
	v361 = v312 + int32(_a_F_VP8LEncodeStream_1) + v360
	v362 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v361))))
	if v362 == int32(0) {
		v392 = v360
		v401 = v361
		goto L78
	} else {
		goto L79
	}
L77:
	;
	v426 = v353
	v427 = v415
	goto L75
L78:
	;
	v402 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v401))) = uint8(v402)
	*(*int32)(unsafe.Add(mBase, uint32(v312+v392<<(uint(int32(2))%32)))) = v353
	if int32(255) < v341 {
		goto L67
	} else {
		goto L84
	}
L79:
	;
	v367 = v360
	goto L80
L80:
	;
	v380 = *(*int32)(unsafe.Add(mBase, uint32(v312+v367<<(uint(int32(2))%32))))
	if v380 == v353 {
		v415 = v341
		goto L77
	} else {
		goto L82
	}
L81:
	;
	v392 = v387
	v401 = v388
	goto L78
L82:
	;
	v387 = (v367 + int32(1)) & int32(1023)
	v388 = v312 + int32(_a_F_VP8LEncodeStream_1) + v387
	v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v388))))
	if v389 != 0 {
		v367 = v387
		goto L80
	} else {
		goto L83
	}
L83:
	;
	goto L81
L84:
	;
	v415 = v341 + int32(1)
	goto L77
L85:
	;
	goto L74
L86:
	;
	goto L72
L87:
	;
	v460 = int32(0)
	v464 = v312
	v465 = v460
	v472 = v460
	goto L88
L88:
	;
	v476 = v312 + int32(_a_F_VP8LEncodeStream_1) + v472
	v477 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v476))))
	if v477 == int32(0) {
		v487 = v465
		goto L90
	} else {
		goto L91
	}
L89:
	;
	v510 = m.G2
	F_qsort(m, v289, v502, int32(4), v510+int32(328))
	mBase = m.M
	v518 = v502
	goto L66
L90:
	;
	v490 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v476+int32(1)))))
	if v490 == int32(0) {
		v502 = v487
		goto L92
	} else {
		goto L93
	}
L91:
	;
	v483 = *(*int32)(unsafe.Add(mBase, uint32(v464)))
	*(*int32)(unsafe.Add(mBase, uint32(v289+v465<<(uint(int32(2))%32)))) = v483
	v487 = v465 + int32(1)
	goto L90
L92:
	;
	v506 = v472 + int32(2)
	if v506 != int32(1024) {
		v464 = v464 + int32(8)
		v465 = v502
		v472 = v506
		goto L88
	} else {
		goto L94
	}
L93:
	;
	v498 = *(*int32)(unsafe.Add(mBase, uint32(v464+int32(4))))
	*(*int32)(unsafe.Add(mBase, uint32(v289+v487<<(uint(int32(2))%32)))) = v498
	v502 = v487 + int32(1)
	goto L92
L94:
	;
	goto L89
L95:
	;
	v533 = v518
	goto L97
L96:
	;
	v533 = int32(0)
	goto L97
L97:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+68)) = v533
	if v518 < int32(257) {
		goto L98
	} else {
		goto L99
	}
L98:
	;
	v539 = int32(9)
	goto L100
L99:
	;
	v539 = int32(7)
	goto L100
L100:
	;
	v540 = v539 - v287
	v541 = int32(9)
	if base.Ui32(v540) < base.Ui32(v541) {
		goto L101
	} else {
		goto L102
	}
L101:
	;
	v544 = v540
	goto L103
L102:
	;
	v544 = v541
	goto L103
L103:
	;
	if v540 < int32(2) {
		goto L104
	} else {
		goto L105
	}
L104:
	;
	v547 = int32(2)
	goto L106
L105:
	;
	v547 = v544
	goto L106
L106:
	;
	v548 = int32(1) << (uint(v547) % 32)
	v549 = *(*int32)(unsafe.Add(mBase, uint32(v283)+8))
	v550 = int32(-1)
	v551 = v549 + v550
	v554 = *(*int32)(unsafe.Add(mBase, uint32(v283)+12))
	v556 = v554 + v550
	v559 = int32(base.Ui32(v548+v551)>>(uint(v547)%32)) * int32(base.Ui32(v548+v556)>>(uint(v547)%32))
	if base.Ui32(int32(8)) < base.Ui32(v547) {
		goto L109
	} else {
		goto L110
	}
L107:
	;
	if base.Ui32(int32(3)) <= base.Ui32(v645) {
		goto L117
	} else {
		goto L118
	}
L108:
	;
	v569 = v547
	goto L112
L109:
	;
	v645 = v547
	v648 = v559
	goto L107
L110:
	;
	if int32(2601) <= v559 {
		goto L108
	} else {
		goto L111
	}
L111:
	;
	goto L109
L112:
	;
	v618 = int32(2) << (uint(v569) % 32)
	v621 = v569 + int32(1)
	v625 = int32(base.Ui32(v618+v551)>>(uint(v621)%32)) * int32(base.Ui32(v618+v556)>>(uint(v621)%32))
	if base.Ui32(int32(7)) < base.Ui32(v569) {
		v645 = v621
		v648 = v625
		goto L107
	} else {
		goto L114
	}
L113:
	;
	v645 = v621
	v648 = v625
	goto L107
L114:
	;
	if int32(2600) < v625 {
		v569 = v621
		goto L112
	} else {
		goto L115
	}
L115:
	;
	goto L113
L116:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+36)) = v785
	v821 = int32(4)
	if v287 == v821 {
		goto L128
	} else {
		goto L129
	}
L117:
	;
	if v648 == int32(1) {
		goto L119
	} else {
		goto L120
	}
L118:
	;
	v785 = v645
	goto L116
L119:
	;
	v687 = int32(1)
	v689 = v645 + int32(-1)
	v690 = v687 << (uint(v689) % 32)
	if int32(base.Ui32(v690+v551)>>(uint(v689)%32))*int32(base.Ui32(v690+v556)>>(uint(v689)%32)) == v687 {
		goto L121
	} else {
		goto L122
	}
L120:
	;
	v785 = v645
	goto L116
L121:
	;
	v713 = v645
	goto L123
L122:
	;
	v785 = v645
	goto L116
L123:
	;
	v752 = v713 + int32(-1)
	if int32(3) <= v752 {
		goto L125
	} else {
		goto L126
	}
L124:
	;
	v785 = v752
	goto L116
L125:
	;
	v757 = v713 + int32(-2)
	v758 = int32(1)
	v759 = v758 << (uint(v757) % 32)
	if int32(base.Ui32(v759+v551)>>(uint(v757)%32))*int32(base.Ui32(v759+v556)>>(uint(v757)%32)) == v758 {
		v713 = v752
		goto L123
	} else {
		goto L127
	}
L126:
	;
	v785 = int32(2)
	goto L116
L127:
	;
	goto L124
L128:
	;
	v827 = int32(5)
	goto L130
L129:
	;
	v827 = v821
	goto L130
L130:
	;
	if v287 < int32(4) {
		goto L131
	} else {
		goto L132
	}
L131:
	;
	v830 = int32(6)
	goto L133
L132:
	;
	v830 = v827
	goto L133
L133:
	;
	if base.Ui32(v830) < base.Ui32(v785) {
		goto L134
	} else {
		goto L135
	}
L134:
	;
	v832 = v830
	goto L136
L135:
	;
	v832 = v785
	goto L136
L136:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+44)) = v832
	*(*int32)(unsafe.Add(mBase, uint32(v79)+40)) = v832
	if v287 != 0 {
		goto L141
	} else {
		goto L142
	}
L137:
	;
	v3012 = v79 + int32(2120)
	v3016 = base.I32_div_s(v2943+int32(-1), int32(16))
	v3018 = v3016 + int32(1)
	v3019 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3012, int32(4), v3019)
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2140)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3012)+12)) = v79 + int32(2128)
	v3029 = int32(256)
	if v3029 < v3018 {
		goto L461
	} else {
		goto L462
	}
L138:
	;
	v3007 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3007 != 0 {
		goto L458
	} else {
		goto L459
	}
L139:
	;
	v2939 = v79 + int32(2216)
	v2940 = *(*int32)(unsafe.Add(mBase, uint32(v79)+4))
	v2941 = *(*int32)(unsafe.Add(mBase, uint32(v2940)+12))
	v2942 = *(*int32)(unsafe.Add(mBase, uint32(v2940)+8))
	v2943 = v2941 * v2942
	v2947 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v2943), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v2939))) = v2947
	if v2947 != 0 {
		goto L454
	} else {
		goto L455
	}
L140:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1292)) = v2801
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1288)) = int32(3)
	v2842 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1304)) = v2842
	if v2807 == v2842 {
		v2906 = v2807
		v2907 = v2808
		v2908 = v2809
		goto L139
	} else {
		goto L450
	}
L141:
	;
	v850 = base.B2i32(v533 < int32(17))
	v851 = *(*int32)(unsafe.Add(mBase, uint32(v283)+56))
	v852 = *(*int32)(unsafe.Add(mBase, uint32(v283)+52))
	if int32(256) < v518 {
		goto L150
	} else {
		goto L151
	}
L142:
	;
	v835 = int32(0)
	v839 = base.B2i32(v518 < int32(257))
	if v518 < int32(257) {
		goto L143
	} else {
		goto L144
	}
L143:
	;
	v840 = v835
	goto L145
L144:
	;
	v840 = int32(3)
	goto L145
L145:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1284)) = v840
	if v518 < int32(257) {
		goto L146
	} else {
		goto L147
	}
L146:
	;
	v844 = int32(4)
	goto L148
L147:
	;
	v844 = int32(3)
	goto L148
L148:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1280)) = v844
	v846 = int32(1)
	v2801 = v835
	v2807 = v846
	v2808 = int32(0)
	v2809 = v846
	goto L140
L149:
	;
	v2363 = int32(1)
	if v533 < int32(17) {
		goto L414
	} else {
		goto L415
	}
L150:
	;
	v857 = int64(1)
	v858 = int32(_a_F_VP8LEncodeStream_2)
	goto L156
L151:
	;
	if v533 < int32(17) {
		v2329 = v821
		v2332 = int32(1)
		goto L149
	} else {
		goto L152
	}
L152:
	;
	goto L150
L153:
	;
	if v878 == int32(0) {
		goto L138
	} else {
		goto L159
	}
L154:
	;
	goto L153
L155:
	;
	v876 = F_calloc(m, base.I32_wrap_i64(v857), v858)
	mBase = m.M
	v878 = v876
	goto L154
L156:
	;
	v865 = base.I64_div_u_s(int64(2147418112), v857)
	v866 = int32(0)
	v867 = base.I64_extend_i32_u(v858)
	if base.Ui64(int64(4294967295)) < base.Ui64(v867*v857) {
		v878 = v866
		goto L154
	} else {
		goto L157
	}
L157:
	;
	if base.Ui64(v865) < base.Ui64(v867) {
		v878 = v866
		goto L154
	} else {
		goto L158
	}
L158:
	;
	goto L155
L159:
	;
	if v284 < int32(1) {
		goto L160
	} else {
		goto L161
	}
L160:
	;
	v1244 = *(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[0])))
	v1245 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[0]))) = v1244 + v1245
	v1248 = *(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[1])))
	*(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[1]))) = v1248 + v1245
	v1252 = *(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[2]))) = v1252 + v1245
	v1256 = *(*int32)(unsafe.Add(mBase, uint32(v878)+3072))
	*(*int32)(unsafe.Add(mBase, uint32(v878)+3072)) = v1256 + v1245
	v1260 = *(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[3])))
	*(*int32)(unsafe.Add(mBase, uint32(v878)+uint32(_c_F_VP8LEncodeStream[3]))) = v1260 + v1245
	v1264 = *(*int32)(unsafe.Add(mBase, uint32(v878)+1024))
	*(*int32)(unsafe.Add(mBase, uint32(v878)+1024)) = v1264 + v1245
	v1268 = m.G117
	v1273 = m.G0
	v1275 = v1273 - int32(32)
	m.G0 = v1275
	F_VP8LBitsEntropyUnrefined(m, v878, int32(256), v1275+int32(8))
	mBase = m.M
	v1280 = *(*int32)(unsafe.Add(mBase, uint32(v1275)+20))
	if v1280 <= int32(4) {
		goto L177
	} else {
		goto L178
	}
L161:
	;
	if v285 < int32(1) {
		goto L160
	} else {
		goto L162
	}
L162:
	;
	v912 = *(*int32)(unsafe.Add(mBase, uint32(v852)))
	v913 = int32(0)
	__phi920 = v912
	__phi930 = v852
	__phi951 = v913
	__phi952 = v913
	v920 = __phi920
	v930 = __phi930
	v951 = __phi951
	v952 = __phi952
	goto L163
L163:
	;
	v974 = v920
	v984 = int32(0)
	v985 = v285
	goto L165
L164:
	;
	goto L160
L165:
	;
	v1023 = *(*int32)(unsafe.Add(mBase, uint32(v930+v984)))
	v1026 = int32(-16711936)
	v1028 = v1023 | int32(16711680) - v974&v1026
	v1033 = int32(16711935)
	v1035 = v1023 | int32(_a_F_VP8LEncodeStream_3) - v974&v1033
	v1038 = v1028&v1026 | v1035&v1033
	if v1038 == int32(0) {
		goto L167
	} else {
		goto L168
	}
L166:
	;
	v1189 = v951 + int32(1)
	if v1189 != v284 {
		__phi920 = v1023
		__phi930 = v930 + v851<<(uint(int32(2))%32)
		__phi951 = v1189
		__phi952 = v930
		v920 = __phi920
		v930 = __phi930
		v951 = __phi951
		v952 = __phi952
		goto L163
	} else {
		goto L173
	}
L167:
	;
	v1186 = v985 + int32(-1)
	if v1186 != 0 {
		v974 = v1023
		v984 = v984 + int32(4)
		v985 = v1186
		goto L165
	} else {
		goto L172
	}
L168:
	;
	if v952 == int32(0) {
		goto L169
	} else {
		goto L170
	}
L169:
	;
	v1046 = int32(22)
	v1048 = int32(1020)
	v1050 = v878 + int32(base.Ui32(v1023)>>(uint(v1046)%32))&v1048
	v1051 = *(*int32)(unsafe.Add(mBase, uint32(v1050)))
	v1052 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v1050))) = v1051 + v1052
	v1055 = int32(14)
	v1059 = v878 + int32(_a_F_VP8LEncodeStream_1) + int32(base.Ui32(v1023)>>(uint(v1055)%32))&v1048
	v1060 = *(*int32)(unsafe.Add(mBase, uint32(v1059)))
	*(*int32)(unsafe.Add(mBase, uint32(v1059))) = v1060 + v1052
	v1064 = int32(6)
	v1068 = v878 + int32(2048) + int32(base.Ui32(v1023)>>(uint(v1064)%32))&v1048
	v1069 = *(*int32)(unsafe.Add(mBase, uint32(v1068)))
	*(*int32)(unsafe.Add(mBase, uint32(v1068))) = v1069 + v1052
	v1073 = int32(255)
	v1075 = int32(2)
	v1077 = v878 + int32(_a_F_VP8LEncodeStream_4) + v1023&v1073<<(uint(v1075)%32)
	v1078 = *(*int32)(unsafe.Add(mBase, uint32(v1077)))
	*(*int32)(unsafe.Add(mBase, uint32(v1077))) = v1078 + v1052
	v1086 = v878 + int32(1024) + int32(base.Ui32(v1028)>>(uint(v1046)%32))&v1048
	v1087 = *(*int32)(unsafe.Add(mBase, uint32(v1086)))
	*(*int32)(unsafe.Add(mBase, uint32(v1086))) = v1087 + v1052
	v1095 = v878 + int32(_a_F_VP8LEncodeStream_0) + int32(base.Ui32(v1038)>>(uint(v1055)%32))&v1048
	v1096 = *(*int32)(unsafe.Add(mBase, uint32(v1095)))
	*(*int32)(unsafe.Add(mBase, uint32(v1095))) = v1096 + v1052
	v1104 = v878 + int32(3072) + int32(base.Ui32(v1028)>>(uint(v1064)%32))&v1048
	v1105 = *(*int32)(unsafe.Add(mBase, uint32(v1104)))
	*(*int32)(unsafe.Add(mBase, uint32(v1104))) = v1105 + v1052
	v1113 = v878 + int32(_a_F_VP8LEncodeStream_5) + v1035&v1073<<(uint(v1075)%32)
	v1114 = *(*int32)(unsafe.Add(mBase, uint32(v1113)))
	*(*int32)(unsafe.Add(mBase, uint32(v1113))) = v1114 + v1052
	v1118 = int32(16)
	v1120 = int32(8)
	v1121 = v1023 >> (uint(v1120) % 32)
	v1127 = v878 + int32(_a_F_VP8LEncodeStream_6) + (int32(base.Ui32(v1023)>>(uint(v1118)%32))-v1121)&v1073<<(uint(v1075)%32)
	v1128 = *(*int32)(unsafe.Add(mBase, uint32(v1127)))
	*(*int32)(unsafe.Add(mBase, uint32(v1127))) = v1128 + v1052
	v1137 = v878 + int32(_a_F_VP8LEncodeStream_7) + (v1023-v1121)&v1073<<(uint(v1075)%32)
	v1138 = *(*int32)(unsafe.Add(mBase, uint32(v1137)))
	*(*int32)(unsafe.Add(mBase, uint32(v1137))) = v1138 + v1052
	v1145 = v1038 >> (uint(v1120) % 32)
	v1151 = v878 + int32(_a_F_VP8LEncodeStream_8) + (int32(base.Ui32(v1038)>>(uint(v1118)%32))-v1145)&v1073<<(uint(v1075)%32)
	v1152 = *(*int32)(unsafe.Add(mBase, uint32(v1151)))
	*(*int32)(unsafe.Add(mBase, uint32(v1151))) = v1152 + v1052
	v1161 = v878 + int32(_a_F_VP8LEncodeStream_9) + (v1035-v1145)&v1073<<(uint(v1075)%32)
	v1162 = *(*int32)(unsafe.Add(mBase, uint32(v1161)))
	*(*int32)(unsafe.Add(mBase, uint32(v1161))) = v1162 + v1052
	v1175 = v878 + int32(_a_F_VP8LEncodeStream_10) + int32(base.Ui32((int32(base.Ui32(v1023)>>(uint(int32(19))%32))+v1023)*int32(969276327))>>(uint(v1046)%32))&v1048
	v1176 = *(*int32)(unsafe.Add(mBase, uint32(v1175)))
	*(*int32)(unsafe.Add(mBase, uint32(v1175))) = v1176 + v1052
	goto L167
L170:
	;
	v1044 = *(*int32)(unsafe.Add(mBase, uint32(v952+v984)))
	if v1023 == v1044 {
		goto L167
	} else {
		goto L171
	}
L171:
	;
	goto L169
L172:
	;
	goto L166
L173:
	;
	goto L164
L174:
	;
	v1340 = m.G0
	v1342 = v1340 - int32(32)
	m.G0 = v1342
	F_VP8LBitsEntropyUnrefined(m, v878+int32(1024), int32(256), v1342+int32(8))
	mBase = m.M
	v1347 = *(*int32)(unsafe.Add(mBase, uint32(v1342)+20))
	if v1347 <= int32(4) {
		goto L194
	} else {
		goto L195
	}
L175:
	;
	m.G0 = v1275 + int32(32)
	goto L174
L176:
	;
	v1305 = *(*int64)(unsafe.Add(mBase, uint32(v1275)+8))
	v1308 = *(*int32)(unsafe.Add(mBase, uint32(v1275)+16))
	v1311 = *(*int32)(unsafe.Add(mBase, uint32(v1275)+24))
	v1320 = v1304*base.I64_extend_i32_u(v1308<<(uint(int32(1))%32)-v1311)<<(uint(int64(23))%64) + v1305*(int64(1000)-v1304)
	if v1320 < int64(0) {
		goto L185
	} else {
		goto L186
	}
L177:
	;
	if v1280 < int32(2) {
		v1329 = int64(0)
		goto L175
	} else {
		goto L179
	}
L178:
	;
	v1304 = int64(627)
	goto L176
L179:
	;
	switch v1280 + int32(-2) {
	case 0:
		goto L181
	case 1:
		v1304 = int64(950)
		goto L176
	default:
		goto L180
	}
L180:
	;
	v1304 = int64(700)
	goto L176
L181:
	;
	v1292 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1275)+16)))
	v1295 = *(*int64)(unsafe.Add(mBase, uint32(v1275)+8))
	v1296 = v1292*int64(830472192) + v1295
	if v1296 < int64(0) {
		goto L182
	} else {
		goto L183
	}
L182:
	;
	v1299 = int64(-50)
	goto L184
L183:
	;
	v1299 = int64(50)
	goto L184
L184:
	;
	v1302 = base.I64_div_s(v1299+v1296, int64(100))
	v1329 = v1302
	goto L175
L185:
	;
	v1323 = int64(-500)
	goto L187
L186:
	;
	v1323 = int64(500)
	goto L187
L187:
	;
	v1326 = base.I64_div_s(v1323+v1320, int64(1000))
	if base.Ui64(v1326) < base.Ui64(v1305) {
		goto L188
	} else {
		goto L189
	}
L188:
	;
	v1328 = v1305
	goto L190
L189:
	;
	v1328 = v1326
	goto L190
L190:
	;
	v1329 = v1328
	goto L175
L191:
	;
	v1407 = m.G0
	v1409 = v1407 - int32(32)
	m.G0 = v1409
	F_VP8LBitsEntropyUnrefined(m, v878+int32(2048), int32(256), v1409+int32(8))
	mBase = m.M
	v1414 = *(*int32)(unsafe.Add(mBase, uint32(v1409)+20))
	if v1414 <= int32(4) {
		goto L211
	} else {
		goto L212
	}
L192:
	;
	m.G0 = v1342 + int32(32)
	goto L191
L193:
	;
	v1372 = *(*int64)(unsafe.Add(mBase, uint32(v1342)+8))
	v1375 = *(*int32)(unsafe.Add(mBase, uint32(v1342)+16))
	v1378 = *(*int32)(unsafe.Add(mBase, uint32(v1342)+24))
	v1387 = v1371*base.I64_extend_i32_u(v1375<<(uint(int32(1))%32)-v1378)<<(uint(int64(23))%64) + v1372*(int64(1000)-v1371)
	if v1387 < int64(0) {
		goto L202
	} else {
		goto L203
	}
L194:
	;
	if v1347 < int32(2) {
		v1396 = int64(0)
		goto L192
	} else {
		goto L196
	}
L195:
	;
	v1371 = int64(627)
	goto L193
L196:
	;
	switch v1347 + int32(-2) {
	case 0:
		goto L198
	case 1:
		v1371 = int64(950)
		goto L193
	default:
		goto L197
	}
L197:
	;
	v1371 = int64(700)
	goto L193
L198:
	;
	v1359 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1342)+16)))
	v1362 = *(*int64)(unsafe.Add(mBase, uint32(v1342)+8))
	v1363 = v1359*int64(830472192) + v1362
	if v1363 < int64(0) {
		goto L199
	} else {
		goto L200
	}
L199:
	;
	v1366 = int64(-50)
	goto L201
L200:
	;
	v1366 = int64(50)
	goto L201
L201:
	;
	v1369 = base.I64_div_s(v1366+v1363, int64(100))
	v1396 = v1369
	goto L192
L202:
	;
	v1390 = int64(-500)
	goto L204
L203:
	;
	v1390 = int64(500)
	goto L204
L204:
	;
	v1393 = base.I64_div_s(v1390+v1387, int64(1000))
	if base.Ui64(v1393) < base.Ui64(v1372) {
		goto L205
	} else {
		goto L206
	}
L205:
	;
	v1395 = v1372
	goto L207
L206:
	;
	v1395 = v1393
	goto L207
L207:
	;
	v1396 = v1395
	goto L192
L208:
	;
	v1474 = m.G0
	v1476 = v1474 - int32(32)
	m.G0 = v1476
	F_VP8LBitsEntropyUnrefined(m, v878+int32(3072), int32(256), v1476+int32(8))
	mBase = m.M
	v1481 = *(*int32)(unsafe.Add(mBase, uint32(v1476)+20))
	if v1481 <= int32(4) {
		goto L228
	} else {
		goto L229
	}
L209:
	;
	m.G0 = v1409 + int32(32)
	goto L208
L210:
	;
	v1439 = *(*int64)(unsafe.Add(mBase, uint32(v1409)+8))
	v1442 = *(*int32)(unsafe.Add(mBase, uint32(v1409)+16))
	v1445 = *(*int32)(unsafe.Add(mBase, uint32(v1409)+24))
	v1454 = v1438*base.I64_extend_i32_u(v1442<<(uint(int32(1))%32)-v1445)<<(uint(int64(23))%64) + v1439*(int64(1000)-v1438)
	if v1454 < int64(0) {
		goto L219
	} else {
		goto L220
	}
L211:
	;
	if v1414 < int32(2) {
		v1463 = int64(0)
		goto L209
	} else {
		goto L213
	}
L212:
	;
	v1438 = int64(627)
	goto L210
L213:
	;
	switch v1414 + int32(-2) {
	case 0:
		goto L215
	case 1:
		v1438 = int64(950)
		goto L210
	default:
		goto L214
	}
L214:
	;
	v1438 = int64(700)
	goto L210
L215:
	;
	v1426 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1409)+16)))
	v1429 = *(*int64)(unsafe.Add(mBase, uint32(v1409)+8))
	v1430 = v1426*int64(830472192) + v1429
	if v1430 < int64(0) {
		goto L216
	} else {
		goto L217
	}
L216:
	;
	v1433 = int64(-50)
	goto L218
L217:
	;
	v1433 = int64(50)
	goto L218
L218:
	;
	v1436 = base.I64_div_s(v1433+v1430, int64(100))
	v1463 = v1436
	goto L209
L219:
	;
	v1457 = int64(-500)
	goto L221
L220:
	;
	v1457 = int64(500)
	goto L221
L221:
	;
	v1460 = base.I64_div_s(v1457+v1454, int64(1000))
	if base.Ui64(v1460) < base.Ui64(v1439) {
		goto L222
	} else {
		goto L223
	}
L222:
	;
	v1462 = v1439
	goto L224
L223:
	;
	v1462 = v1460
	goto L224
L224:
	;
	v1463 = v1462
	goto L209
L225:
	;
	v1541 = m.G0
	v1543 = v1541 - int32(32)
	m.G0 = v1543
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_1), int32(256), v1543+int32(8))
	mBase = m.M
	v1548 = *(*int32)(unsafe.Add(mBase, uint32(v1543)+20))
	if v1548 <= int32(4) {
		goto L245
	} else {
		goto L246
	}
L226:
	;
	m.G0 = v1476 + int32(32)
	goto L225
L227:
	;
	v1506 = *(*int64)(unsafe.Add(mBase, uint32(v1476)+8))
	v1509 = *(*int32)(unsafe.Add(mBase, uint32(v1476)+16))
	v1512 = *(*int32)(unsafe.Add(mBase, uint32(v1476)+24))
	v1521 = v1505*base.I64_extend_i32_u(v1509<<(uint(int32(1))%32)-v1512)<<(uint(int64(23))%64) + v1506*(int64(1000)-v1505)
	if v1521 < int64(0) {
		goto L236
	} else {
		goto L237
	}
L228:
	;
	if v1481 < int32(2) {
		v1530 = int64(0)
		goto L226
	} else {
		goto L230
	}
L229:
	;
	v1505 = int64(627)
	goto L227
L230:
	;
	switch v1481 + int32(-2) {
	case 0:
		goto L232
	case 1:
		v1505 = int64(950)
		goto L227
	default:
		goto L231
	}
L231:
	;
	v1505 = int64(700)
	goto L227
L232:
	;
	v1493 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1476)+16)))
	v1496 = *(*int64)(unsafe.Add(mBase, uint32(v1476)+8))
	v1497 = v1493*int64(830472192) + v1496
	if v1497 < int64(0) {
		goto L233
	} else {
		goto L234
	}
L233:
	;
	v1500 = int64(-50)
	goto L235
L234:
	;
	v1500 = int64(50)
	goto L235
L235:
	;
	v1503 = base.I64_div_s(v1500+v1497, int64(100))
	v1530 = v1503
	goto L226
L236:
	;
	v1524 = int64(-500)
	goto L238
L237:
	;
	v1524 = int64(500)
	goto L238
L238:
	;
	v1527 = base.I64_div_s(v1524+v1521, int64(1000))
	if base.Ui64(v1527) < base.Ui64(v1506) {
		goto L239
	} else {
		goto L240
	}
L239:
	;
	v1529 = v1506
	goto L241
L240:
	;
	v1529 = v1527
	goto L241
L241:
	;
	v1530 = v1529
	goto L226
L242:
	;
	v1608 = m.G0
	v1610 = v1608 - int32(32)
	m.G0 = v1610
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_0), int32(256), v1610+int32(8))
	mBase = m.M
	v1615 = *(*int32)(unsafe.Add(mBase, uint32(v1610)+20))
	if v1615 <= int32(4) {
		goto L262
	} else {
		goto L263
	}
L243:
	;
	m.G0 = v1543 + int32(32)
	goto L242
L244:
	;
	v1573 = *(*int64)(unsafe.Add(mBase, uint32(v1543)+8))
	v1576 = *(*int32)(unsafe.Add(mBase, uint32(v1543)+16))
	v1579 = *(*int32)(unsafe.Add(mBase, uint32(v1543)+24))
	v1588 = v1572*base.I64_extend_i32_u(v1576<<(uint(int32(1))%32)-v1579)<<(uint(int64(23))%64) + v1573*(int64(1000)-v1572)
	if v1588 < int64(0) {
		goto L253
	} else {
		goto L254
	}
L245:
	;
	if v1548 < int32(2) {
		v1597 = int64(0)
		goto L243
	} else {
		goto L247
	}
L246:
	;
	v1572 = int64(627)
	goto L244
L247:
	;
	switch v1548 + int32(-2) {
	case 0:
		goto L249
	case 1:
		v1572 = int64(950)
		goto L244
	default:
		goto L248
	}
L248:
	;
	v1572 = int64(700)
	goto L244
L249:
	;
	v1560 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1543)+16)))
	v1563 = *(*int64)(unsafe.Add(mBase, uint32(v1543)+8))
	v1564 = v1560*int64(830472192) + v1563
	if v1564 < int64(0) {
		goto L250
	} else {
		goto L251
	}
L250:
	;
	v1567 = int64(-50)
	goto L252
L251:
	;
	v1567 = int64(50)
	goto L252
L252:
	;
	v1570 = base.I64_div_s(v1567+v1564, int64(100))
	v1597 = v1570
	goto L243
L253:
	;
	v1591 = int64(-500)
	goto L255
L254:
	;
	v1591 = int64(500)
	goto L255
L255:
	;
	v1594 = base.I64_div_s(v1591+v1588, int64(1000))
	if base.Ui64(v1594) < base.Ui64(v1573) {
		goto L256
	} else {
		goto L257
	}
L256:
	;
	v1596 = v1573
	goto L258
L257:
	;
	v1596 = v1594
	goto L258
L258:
	;
	v1597 = v1596
	goto L243
L259:
	;
	v1675 = m.G0
	v1677 = v1675 - int32(32)
	m.G0 = v1677
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_4), int32(256), v1677+int32(8))
	mBase = m.M
	v1682 = *(*int32)(unsafe.Add(mBase, uint32(v1677)+20))
	if v1682 <= int32(4) {
		goto L279
	} else {
		goto L280
	}
L260:
	;
	m.G0 = v1610 + int32(32)
	goto L259
L261:
	;
	v1640 = *(*int64)(unsafe.Add(mBase, uint32(v1610)+8))
	v1643 = *(*int32)(unsafe.Add(mBase, uint32(v1610)+16))
	v1646 = *(*int32)(unsafe.Add(mBase, uint32(v1610)+24))
	v1655 = v1639*base.I64_extend_i32_u(v1643<<(uint(int32(1))%32)-v1646)<<(uint(int64(23))%64) + v1640*(int64(1000)-v1639)
	if v1655 < int64(0) {
		goto L270
	} else {
		goto L271
	}
L262:
	;
	if v1615 < int32(2) {
		v1664 = int64(0)
		goto L260
	} else {
		goto L264
	}
L263:
	;
	v1639 = int64(627)
	goto L261
L264:
	;
	switch v1615 + int32(-2) {
	case 0:
		goto L266
	case 1:
		v1639 = int64(950)
		goto L261
	default:
		goto L265
	}
L265:
	;
	v1639 = int64(700)
	goto L261
L266:
	;
	v1627 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1610)+16)))
	v1630 = *(*int64)(unsafe.Add(mBase, uint32(v1610)+8))
	v1631 = v1627*int64(830472192) + v1630
	if v1631 < int64(0) {
		goto L267
	} else {
		goto L268
	}
L267:
	;
	v1634 = int64(-50)
	goto L269
L268:
	;
	v1634 = int64(50)
	goto L269
L269:
	;
	v1637 = base.I64_div_s(v1634+v1631, int64(100))
	v1664 = v1637
	goto L260
L270:
	;
	v1658 = int64(-500)
	goto L272
L271:
	;
	v1658 = int64(500)
	goto L272
L272:
	;
	v1661 = base.I64_div_s(v1658+v1655, int64(1000))
	if base.Ui64(v1661) < base.Ui64(v1640) {
		goto L273
	} else {
		goto L274
	}
L273:
	;
	v1663 = v1640
	goto L275
L274:
	;
	v1663 = v1661
	goto L275
L275:
	;
	v1664 = v1663
	goto L260
L276:
	;
	v1742 = m.G0
	v1744 = v1742 - int32(32)
	m.G0 = v1744
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_5), int32(256), v1744+int32(8))
	mBase = m.M
	v1749 = *(*int32)(unsafe.Add(mBase, uint32(v1744)+20))
	if v1749 <= int32(4) {
		goto L296
	} else {
		goto L297
	}
L277:
	;
	m.G0 = v1677 + int32(32)
	goto L276
L278:
	;
	v1707 = *(*int64)(unsafe.Add(mBase, uint32(v1677)+8))
	v1710 = *(*int32)(unsafe.Add(mBase, uint32(v1677)+16))
	v1713 = *(*int32)(unsafe.Add(mBase, uint32(v1677)+24))
	v1722 = v1706*base.I64_extend_i32_u(v1710<<(uint(int32(1))%32)-v1713)<<(uint(int64(23))%64) + v1707*(int64(1000)-v1706)
	if v1722 < int64(0) {
		goto L287
	} else {
		goto L288
	}
L279:
	;
	if v1682 < int32(2) {
		v1731 = int64(0)
		goto L277
	} else {
		goto L281
	}
L280:
	;
	v1706 = int64(627)
	goto L278
L281:
	;
	switch v1682 + int32(-2) {
	case 0:
		goto L283
	case 1:
		v1706 = int64(950)
		goto L278
	default:
		goto L282
	}
L282:
	;
	v1706 = int64(700)
	goto L278
L283:
	;
	v1694 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1677)+16)))
	v1697 = *(*int64)(unsafe.Add(mBase, uint32(v1677)+8))
	v1698 = v1694*int64(830472192) + v1697
	if v1698 < int64(0) {
		goto L284
	} else {
		goto L285
	}
L284:
	;
	v1701 = int64(-50)
	goto L286
L285:
	;
	v1701 = int64(50)
	goto L286
L286:
	;
	v1704 = base.I64_div_s(v1701+v1698, int64(100))
	v1731 = v1704
	goto L277
L287:
	;
	v1725 = int64(-500)
	goto L289
L288:
	;
	v1725 = int64(500)
	goto L289
L289:
	;
	v1728 = base.I64_div_s(v1725+v1722, int64(1000))
	if base.Ui64(v1728) < base.Ui64(v1707) {
		goto L290
	} else {
		goto L291
	}
L290:
	;
	v1730 = v1707
	goto L292
L291:
	;
	v1730 = v1728
	goto L292
L292:
	;
	v1731 = v1730
	goto L277
L293:
	;
	v1809 = m.G0
	v1811 = v1809 - int32(32)
	m.G0 = v1811
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_6), int32(256), v1811+int32(8))
	mBase = m.M
	v1816 = *(*int32)(unsafe.Add(mBase, uint32(v1811)+20))
	if v1816 <= int32(4) {
		goto L313
	} else {
		goto L314
	}
L294:
	;
	m.G0 = v1744 + int32(32)
	goto L293
L295:
	;
	v1774 = *(*int64)(unsafe.Add(mBase, uint32(v1744)+8))
	v1777 = *(*int32)(unsafe.Add(mBase, uint32(v1744)+16))
	v1780 = *(*int32)(unsafe.Add(mBase, uint32(v1744)+24))
	v1789 = v1773*base.I64_extend_i32_u(v1777<<(uint(int32(1))%32)-v1780)<<(uint(int64(23))%64) + v1774*(int64(1000)-v1773)
	if v1789 < int64(0) {
		goto L304
	} else {
		goto L305
	}
L296:
	;
	if v1749 < int32(2) {
		v1798 = int64(0)
		goto L294
	} else {
		goto L298
	}
L297:
	;
	v1773 = int64(627)
	goto L295
L298:
	;
	switch v1749 + int32(-2) {
	case 0:
		goto L300
	case 1:
		v1773 = int64(950)
		goto L295
	default:
		goto L299
	}
L299:
	;
	v1773 = int64(700)
	goto L295
L300:
	;
	v1761 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1744)+16)))
	v1764 = *(*int64)(unsafe.Add(mBase, uint32(v1744)+8))
	v1765 = v1761*int64(830472192) + v1764
	if v1765 < int64(0) {
		goto L301
	} else {
		goto L302
	}
L301:
	;
	v1768 = int64(-50)
	goto L303
L302:
	;
	v1768 = int64(50)
	goto L303
L303:
	;
	v1771 = base.I64_div_s(v1768+v1765, int64(100))
	v1798 = v1771
	goto L294
L304:
	;
	v1792 = int64(-500)
	goto L306
L305:
	;
	v1792 = int64(500)
	goto L306
L306:
	;
	v1795 = base.I64_div_s(v1792+v1789, int64(1000))
	if base.Ui64(v1795) < base.Ui64(v1774) {
		goto L307
	} else {
		goto L308
	}
L307:
	;
	v1797 = v1774
	goto L309
L308:
	;
	v1797 = v1795
	goto L309
L309:
	;
	v1798 = v1797
	goto L294
L310:
	;
	v1876 = m.G0
	v1878 = v1876 - int32(32)
	m.G0 = v1878
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_8), int32(256), v1878+int32(8))
	mBase = m.M
	v1883 = *(*int32)(unsafe.Add(mBase, uint32(v1878)+20))
	if v1883 <= int32(4) {
		goto L330
	} else {
		goto L331
	}
L311:
	;
	m.G0 = v1811 + int32(32)
	goto L310
L312:
	;
	v1841 = *(*int64)(unsafe.Add(mBase, uint32(v1811)+8))
	v1844 = *(*int32)(unsafe.Add(mBase, uint32(v1811)+16))
	v1847 = *(*int32)(unsafe.Add(mBase, uint32(v1811)+24))
	v1856 = v1840*base.I64_extend_i32_u(v1844<<(uint(int32(1))%32)-v1847)<<(uint(int64(23))%64) + v1841*(int64(1000)-v1840)
	if v1856 < int64(0) {
		goto L321
	} else {
		goto L322
	}
L313:
	;
	if v1816 < int32(2) {
		v1865 = int64(0)
		goto L311
	} else {
		goto L315
	}
L314:
	;
	v1840 = int64(627)
	goto L312
L315:
	;
	switch v1816 + int32(-2) {
	case 0:
		goto L317
	case 1:
		v1840 = int64(950)
		goto L312
	default:
		goto L316
	}
L316:
	;
	v1840 = int64(700)
	goto L312
L317:
	;
	v1828 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1811)+16)))
	v1831 = *(*int64)(unsafe.Add(mBase, uint32(v1811)+8))
	v1832 = v1828*int64(830472192) + v1831
	if v1832 < int64(0) {
		goto L318
	} else {
		goto L319
	}
L318:
	;
	v1835 = int64(-50)
	goto L320
L319:
	;
	v1835 = int64(50)
	goto L320
L320:
	;
	v1838 = base.I64_div_s(v1835+v1832, int64(100))
	v1865 = v1838
	goto L311
L321:
	;
	v1859 = int64(-500)
	goto L323
L322:
	;
	v1859 = int64(500)
	goto L323
L323:
	;
	v1862 = base.I64_div_s(v1859+v1856, int64(1000))
	if base.Ui64(v1862) < base.Ui64(v1841) {
		goto L324
	} else {
		goto L325
	}
L324:
	;
	v1864 = v1841
	goto L326
L325:
	;
	v1864 = v1862
	goto L326
L326:
	;
	v1865 = v1864
	goto L311
L327:
	;
	v1943 = m.G0
	v1945 = v1943 - int32(32)
	m.G0 = v1945
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_7), int32(256), v1945+int32(8))
	mBase = m.M
	v1950 = *(*int32)(unsafe.Add(mBase, uint32(v1945)+20))
	if v1950 <= int32(4) {
		goto L347
	} else {
		goto L348
	}
L328:
	;
	m.G0 = v1878 + int32(32)
	goto L327
L329:
	;
	v1908 = *(*int64)(unsafe.Add(mBase, uint32(v1878)+8))
	v1911 = *(*int32)(unsafe.Add(mBase, uint32(v1878)+16))
	v1914 = *(*int32)(unsafe.Add(mBase, uint32(v1878)+24))
	v1923 = v1907*base.I64_extend_i32_u(v1911<<(uint(int32(1))%32)-v1914)<<(uint(int64(23))%64) + v1908*(int64(1000)-v1907)
	if v1923 < int64(0) {
		goto L338
	} else {
		goto L339
	}
L330:
	;
	if v1883 < int32(2) {
		v1932 = int64(0)
		goto L328
	} else {
		goto L332
	}
L331:
	;
	v1907 = int64(627)
	goto L329
L332:
	;
	switch v1883 + int32(-2) {
	case 0:
		goto L334
	case 1:
		v1907 = int64(950)
		goto L329
	default:
		goto L333
	}
L333:
	;
	v1907 = int64(700)
	goto L329
L334:
	;
	v1895 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1878)+16)))
	v1898 = *(*int64)(unsafe.Add(mBase, uint32(v1878)+8))
	v1899 = v1895*int64(830472192) + v1898
	if v1899 < int64(0) {
		goto L335
	} else {
		goto L336
	}
L335:
	;
	v1902 = int64(-50)
	goto L337
L336:
	;
	v1902 = int64(50)
	goto L337
L337:
	;
	v1905 = base.I64_div_s(v1902+v1899, int64(100))
	v1932 = v1905
	goto L328
L338:
	;
	v1926 = int64(-500)
	goto L340
L339:
	;
	v1926 = int64(500)
	goto L340
L340:
	;
	v1929 = base.I64_div_s(v1926+v1923, int64(1000))
	if base.Ui64(v1929) < base.Ui64(v1908) {
		goto L341
	} else {
		goto L342
	}
L341:
	;
	v1931 = v1908
	goto L343
L342:
	;
	v1931 = v1929
	goto L343
L343:
	;
	v1932 = v1931
	goto L328
L344:
	;
	v2010 = m.G0
	v2012 = v2010 - int32(32)
	m.G0 = v2012
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_9), int32(256), v2012+int32(8))
	mBase = m.M
	v2017 = *(*int32)(unsafe.Add(mBase, uint32(v2012)+20))
	if v2017 <= int32(4) {
		goto L364
	} else {
		goto L365
	}
L345:
	;
	m.G0 = v1945 + int32(32)
	goto L344
L346:
	;
	v1975 = *(*int64)(unsafe.Add(mBase, uint32(v1945)+8))
	v1978 = *(*int32)(unsafe.Add(mBase, uint32(v1945)+16))
	v1981 = *(*int32)(unsafe.Add(mBase, uint32(v1945)+24))
	v1990 = v1974*base.I64_extend_i32_u(v1978<<(uint(int32(1))%32)-v1981)<<(uint(int64(23))%64) + v1975*(int64(1000)-v1974)
	if v1990 < int64(0) {
		goto L355
	} else {
		goto L356
	}
L347:
	;
	if v1950 < int32(2) {
		v1999 = int64(0)
		goto L345
	} else {
		goto L349
	}
L348:
	;
	v1974 = int64(627)
	goto L346
L349:
	;
	switch v1950 + int32(-2) {
	case 0:
		goto L351
	case 1:
		v1974 = int64(950)
		goto L346
	default:
		goto L350
	}
L350:
	;
	v1974 = int64(700)
	goto L346
L351:
	;
	v1962 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1945)+16)))
	v1965 = *(*int64)(unsafe.Add(mBase, uint32(v1945)+8))
	v1966 = v1962*int64(830472192) + v1965
	if v1966 < int64(0) {
		goto L352
	} else {
		goto L353
	}
L352:
	;
	v1969 = int64(-50)
	goto L354
L353:
	;
	v1969 = int64(50)
	goto L354
L354:
	;
	v1972 = base.I64_div_s(v1969+v1966, int64(100))
	v1999 = v1972
	goto L345
L355:
	;
	v1993 = int64(-500)
	goto L357
L356:
	;
	v1993 = int64(500)
	goto L357
L357:
	;
	v1996 = base.I64_div_s(v1993+v1990, int64(1000))
	if base.Ui64(v1996) < base.Ui64(v1975) {
		goto L358
	} else {
		goto L359
	}
L358:
	;
	v1998 = v1975
	goto L360
L359:
	;
	v1998 = v1996
	goto L360
L360:
	;
	v1999 = v1998
	goto L345
L361:
	;
	v2077 = m.G0
	v2079 = v2077 - int32(32)
	m.G0 = v2079
	F_VP8LBitsEntropyUnrefined(m, v878+int32(_a_F_VP8LEncodeStream_10), int32(256), v2079+int32(8))
	mBase = m.M
	v2084 = *(*int32)(unsafe.Add(mBase, uint32(v2079)+20))
	if v2084 <= int32(4) {
		goto L381
	} else {
		goto L382
	}
L362:
	;
	m.G0 = v2012 + int32(32)
	goto L361
L363:
	;
	v2042 = *(*int64)(unsafe.Add(mBase, uint32(v2012)+8))
	v2045 = *(*int32)(unsafe.Add(mBase, uint32(v2012)+16))
	v2048 = *(*int32)(unsafe.Add(mBase, uint32(v2012)+24))
	v2057 = v2041*base.I64_extend_i32_u(v2045<<(uint(int32(1))%32)-v2048)<<(uint(int64(23))%64) + v2042*(int64(1000)-v2041)
	if v2057 < int64(0) {
		goto L372
	} else {
		goto L373
	}
L364:
	;
	if v2017 < int32(2) {
		v2066 = int64(0)
		goto L362
	} else {
		goto L366
	}
L365:
	;
	v2041 = int64(627)
	goto L363
L366:
	;
	switch v2017 + int32(-2) {
	case 0:
		goto L368
	case 1:
		v2041 = int64(950)
		goto L363
	default:
		goto L367
	}
L367:
	;
	v2041 = int64(700)
	goto L363
L368:
	;
	v2029 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v2012)+16)))
	v2032 = *(*int64)(unsafe.Add(mBase, uint32(v2012)+8))
	v2033 = v2029*int64(830472192) + v2032
	if v2033 < int64(0) {
		goto L369
	} else {
		goto L370
	}
L369:
	;
	v2036 = int64(-50)
	goto L371
L370:
	;
	v2036 = int64(50)
	goto L371
L371:
	;
	v2039 = base.I64_div_s(v2036+v2033, int64(100))
	v2066 = v2039
	goto L362
L372:
	;
	v2060 = int64(-500)
	goto L374
L373:
	;
	v2060 = int64(500)
	goto L374
L374:
	;
	v2063 = base.I64_div_s(v2060+v2057, int64(1000))
	if base.Ui64(v2063) < base.Ui64(v2042) {
		goto L375
	} else {
		goto L376
	}
L375:
	;
	v2065 = v2042
	goto L377
L376:
	;
	v2065 = v2063
	goto L377
L377:
	;
	v2066 = v2065
	goto L362
L378:
	;
	v2138 = v1530 + v1396
	v2142 = int32(1) << (uint(v832) % 32)
	v2144 = int32(-1)
	v2153 = base.I64_extend_i32_u(int32(base.Ui32(v285+v2142+v2144)>>(uint(v832)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v284+v2142+v2144)>>(uint(v832)%32)))
	v2154 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1268)+96)))
	v2156 = v2066 + (v1932 + v2138) + v2153*v2154
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1704)) = v2156
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1712)) = v2133 + base.I64_extend_i32_s(v533)<<(uint(int64(26))%64)
	v2165 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1268)+56)))
	v2167 = v1798 + (v1664 + v2138) + v2153*v2165
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1688)) = v2167
	v2169 = v1463 + v1329
	v2171 = v1731 + (v1597 + v2169)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1680)) = v2171
	v2175 = base.B2i32(base.Ui64(v2167) < base.Ui64(v2171))
	v2176 = int32(3)
	v2179 = *(*int64)(unsafe.Add(mBase, uint32(v56+int32(1680)|v2175<<(uint(v2176)%32))))
	v2181 = v1999 + (v1865 + v2169)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1696)) = v2181
	if base.Ui64(v2181) < base.Ui64(v2179) {
		goto L395
	} else {
		goto L396
	}
L379:
	;
	m.G0 = v2079 + int32(32)
	goto L378
L380:
	;
	v2109 = *(*int64)(unsafe.Add(mBase, uint32(v2079)+8))
	v2112 = *(*int32)(unsafe.Add(mBase, uint32(v2079)+16))
	v2115 = *(*int32)(unsafe.Add(mBase, uint32(v2079)+24))
	v2124 = v2108*base.I64_extend_i32_u(v2112<<(uint(int32(1))%32)-v2115)<<(uint(int64(23))%64) + v2109*(int64(1000)-v2108)
	if v2124 < int64(0) {
		goto L389
	} else {
		goto L390
	}
L381:
	;
	if v2084 < int32(2) {
		v2133 = int64(0)
		goto L379
	} else {
		goto L383
	}
L382:
	;
	v2108 = int64(627)
	goto L380
L383:
	;
	switch v2084 + int32(-2) {
	case 0:
		goto L385
	case 1:
		v2108 = int64(950)
		goto L380
	default:
		goto L384
	}
L384:
	;
	v2108 = int64(700)
	goto L380
L385:
	;
	v2096 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v2079)+16)))
	v2099 = *(*int64)(unsafe.Add(mBase, uint32(v2079)+8))
	v2100 = v2096*int64(830472192) + v2099
	if v2100 < int64(0) {
		goto L386
	} else {
		goto L387
	}
L386:
	;
	v2103 = int64(-50)
	goto L388
L387:
	;
	v2103 = int64(50)
	goto L388
L388:
	;
	v2106 = base.I64_div_s(v2103+v2100, int64(100))
	v2133 = v2106
	goto L379
L389:
	;
	v2127 = int64(-500)
	goto L391
L390:
	;
	v2127 = int64(500)
	goto L391
L391:
	;
	v2130 = base.I64_div_s(v2127+v2124, int64(1000))
	if base.Ui64(v2130) < base.Ui64(v2109) {
		goto L392
	} else {
		goto L393
	}
L392:
	;
	v2132 = v2109
	goto L394
L393:
	;
	v2132 = v2130
	goto L394
L394:
	;
	v2133 = v2132
	goto L379
L395:
	;
	v2186 = int32(2)
	goto L397
L396:
	;
	v2186 = v2175
	goto L397
L397:
	;
	v2192 = *(*int64)(unsafe.Add(mBase, uint32(v56+int32(1680)+v2186<<(uint(int32(3))%32))))
	if base.Ui64(v2156) < base.Ui64(v2192) {
		goto L398
	} else {
		goto L399
	}
L398:
	;
	v2194 = v2176
	goto L400
L399:
	;
	v2194 = v2186
	goto L400
L400:
	;
	if int32(256) < v518 {
		v2207 = v2194
		goto L401
	} else {
		goto L402
	}
L401:
	;
	v2208 = m.G1
	v2213 = v2208 + int32(_a_F_VP8LEncodeStream_11) + v2207<<(uint(int32(1))%32)
	v2214 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2213))))
	v2215 = int32(10)
	v2217 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2213)+1)))
	v2226 = v2214 << (uint(v2215) % 32)
	v2236 = v2217 << (uint(v2215) % 32)
	v2238 = int32(255)
	goto L407
L402:
	;
	v2203 = *(*int64)(unsafe.Add(mBase, uint32(v56+int32(1680)+v2194<<(uint(int32(3))%32))))
	v2204 = *(*int64)(unsafe.Add(mBase, uint32(v56)+1712))
	if base.Ui64(v2204) < base.Ui64(v2203) {
		goto L403
	} else {
		goto L404
	}
L403:
	;
	v2206 = int32(4)
	goto L405
L404:
	;
	v2206 = v2194
	goto L405
L405:
	;
	v2207 = v2206
	goto L401
L406:
	;
	F_free(m, v878)
	mBase = m.M
	goto L413
L407:
	;
	v2274 = int32(0)
	v2275 = v878 + v2236
	v2276 = int32(4)
	v2278 = *(*int32)(unsafe.Add(mBase, uint32(v2275+v2276)))
	v2279 = v878 + v2226
	v2282 = *(*int32)(unsafe.Add(mBase, uint32(v2279+v2276)))
	if v2278|v2282 != 0 {
		v2308 = v2274
		goto L406
	} else {
		goto L409
	}
L408:
	;
	v2308 = int32(1)
	goto L406
L409:
	;
	v2284 = int32(8)
	v2286 = *(*int32)(unsafe.Add(mBase, uint32(v2275+v2284)))
	v2289 = *(*int32)(unsafe.Add(mBase, uint32(v2279+v2284)))
	if v2286|v2289 != 0 {
		v2308 = v2274
		goto L406
	} else {
		goto L410
	}
L410:
	;
	v2291 = int32(12)
	v2293 = *(*int32)(unsafe.Add(mBase, uint32(v2275+v2291)))
	v2296 = *(*int32)(unsafe.Add(mBase, uint32(v2279+v2291)))
	if v2293|v2296 != 0 {
		v2308 = v2274
		goto L406
	} else {
		goto L411
	}
L411:
	;
	v2298 = int32(12)
	v2303 = v2238 + int32(-3)
	if v2303 != 0 {
		v2226 = v2226 + v2298
		v2236 = v2236 + v2298
		v2238 = v2303
		goto L407
	} else {
		goto L412
	}
L412:
	;
	goto L408
L413:
	;
	v2329 = v2207
	v2332 = v2308
	goto L149
L414:
	;
	v2366 = int32(2)
	goto L416
L415:
	;
	v2366 = v2363
	goto L416
L416:
	;
	if v533 < int32(1) {
		goto L417
	} else {
		goto L418
	}
L417:
	;
	v2369 = v2363
	goto L419
L418:
	;
	v2369 = v2366
	goto L419
L419:
	;
	v2370 = *(*float32)(unsafe.Add(mBase, uint32(v286)+4))
	if v287 != int32(6) {
		goto L423
	} else {
		goto L424
	}
L420:
	;
	if v2369 == int32(1) {
		v2801 = v2438
		v2807 = v2439
		v2808 = v2332
		v2809 = v2440
		goto L140
	} else {
		goto L436
	}
L421:
	;
	v2438 = v2433
	v2439 = v2434
	v2440 = int32(0)
	goto L420
L422:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2425))) = v2427
	*(*int32)(unsafe.Add(mBase, uint32(v2424))) = int32(5)
	v2433 = int32(1)
	v2434 = v2426
	goto L421
L423:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1280)) = v2329
	v2400 = int32(1)
	if v518 < int32(257) {
		goto L428
	} else {
		goto L429
	}
L424:
	;
	if base.F32_ne(v2370, float32(100)) != 0 {
		goto L423
	} else {
		goto L425
	}
L425:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1364)) = int64(12884901891)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1336)) = int64(12884901890)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1308)) = int64(12884901889)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1280)) = int64(12884901888)
	if v518 <= int32(256) {
		goto L426
	} else {
		goto L427
	}
L426:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1448)) = int64(4294967301)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1420)) = int64(8589934596)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1392)) = int64(4294967300)
	v2424 = v56 + int32(1476)
	v2425 = v56 + int32(1480)
	v2426 = int32(8)
	v2427 = int32(2)
	goto L422
L427:
	;
	v2433 = int32(1)
	v2434 = int32(4)
	goto L421
L428:
	;
	v2405 = v2400
	goto L430
L429:
	;
	v2405 = int32(3)
	goto L430
L430:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1284)) = v2405
	v2407 = int32(0)
	if base.F32_ge(v2370, float32(75)) != 0 {
		goto L431
	} else {
		goto L432
	}
L431:
	;
	if v287 != int32(5) {
		v2438 = v2407
		v2439 = int32(1)
		v2440 = v2400
		goto L420
	} else {
		goto L433
	}
L432:
	;
	v2438 = v2407
	v2439 = int32(1)
	v2440 = v2400
	goto L420
L433:
	;
	v2414 = int32(1)
	if v2329 == int32(4) {
		goto L434
	} else {
		goto L435
	}
L434:
	;
	v2424 = v56 + int32(1308)
	v2425 = v56 + int32(1312)
	v2426 = int32(2)
	v2427 = v2414
	goto L422
L435:
	;
	v2417 = int32(1)
	v2438 = v2417
	v2439 = v2417
	v2440 = v2414
	goto L420
L436:
	;
	v2444 = v2369 + int32(-1)
	v2447 = int32(3)
	v2448 = v2444 & v2447
	v2463 = int32(0)
	v2473 = v56 + int32(1280)
	goto L437
L437:
	;
	v2513 = v56 + int32(1280) + v2463*int32(28)
	*(*int32)(unsafe.Add(mBase, uint32(v2513)+12)) = v2438
	*(*int32)(unsafe.Add(mBase, uint32(v2513)+8)) = int32(3)
	if base.Ui32(v2369+int32(-2)) < base.Ui32(v2447) {
		v2610 = int32(1)
		goto L439
	} else {
		goto L440
	}
L439:
	;
	if v2448 == int32(0) {
		goto L444
	} else {
		goto L445
	}
L440:
	;
	v2524 = v2473
	v2537 = int32(0)
	goto L441
L441:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2524+int32(44)))) = v2438
	v2577 = int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v2524+int32(40)))) = v2577
	*(*int32)(unsafe.Add(mBase, uint32(v2524+int32(36)))) = v2438
	v2583 = v2524 + int32(32)
	*(*int32)(unsafe.Add(mBase, uint32(v2583))) = v2577
	*(*int32)(unsafe.Add(mBase, uint32(v2524+int32(28)))) = v2438
	*(*int32)(unsafe.Add(mBase, uint32(v2524+int32(24)))) = v2577
	*(*int32)(unsafe.Add(mBase, uint32(v2524+int32(20)))) = v2438
	*(*int32)(unsafe.Add(mBase, uint32(v2524+int32(16)))) = v2577
	v2601 = v2537 + v2577
	if v2444&int32(-4) != v2601 {
		v2524 = v2583
		v2537 = v2601
		goto L441
	} else {
		goto L443
	}
L442:
	;
	v2610 = v2537 + int32(5)
	goto L439
L443:
	;
	goto L442
L444:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2513)+24)) = v2369
	v2784 = v2463 + int32(1)
	if v2784 != v2439 {
		v2463 = v2784
		v2473 = v2473 + int32(28)
		goto L437
	} else {
		goto L449
	}
L445:
	;
	v2667 = v2610 << (uint(int32(3)) % 32)
	v2680 = v2448
	goto L446
L446:
	;
	v2715 = v2473 + v2667
	*(*int32)(unsafe.Add(mBase, uint32(v2715+int32(12)))) = v2438
	v2719 = int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v2715+v2719))) = int32(4)
	v2726 = v2680 + int32(-1)
	if v2726 != 0 {
		v2667 = v2667 + v2719
		v2680 = v2726
		goto L446
	} else {
		goto L448
	}
L447:
	;
	goto L444
L448:
	;
	goto L447
L449:
	;
	v2906 = v2439
	v2907 = v2332
	v2908 = v2440
	goto L139
L450:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1332)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1320)) = v2801
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1316)) = int32(3)
	if v2807 == int32(2) {
		v2906 = v2807
		v2907 = v2808
		v2908 = v2809
		goto L139
	} else {
		goto L451
	}
L451:
	;
	v2853 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1388)) = v2853
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1376)) = v2801
	v2856 = int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1372)) = v2856
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1360)) = v2853
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1348)) = v2801
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1344)) = v2856
	if v2807 == int32(4) {
		v2906 = v2807
		v2907 = v2808
		v2908 = v2809
		goto L139
	} else {
		goto L452
	}
L452:
	;
	v2865 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1500)) = v2865
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1488)) = v2801
	v2868 = int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1484)) = v2868
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1472)) = v2865
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1460)) = v2801
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1456)) = v2868
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1444)) = v2865
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1432)) = v2801
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1428)) = v2868
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1416)) = v2865
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1404)) = v2801
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1400)) = v2868
	v2906 = v2807
	v2907 = v2808
	v2908 = v2809
	goto L139
L453:
	;
	if v2947 != 0 {
		goto L137
	} else {
		goto L456
	}
L454:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2939)+4)) = v2943
	goto L453
L455:
	;
	goto L453
L456:
	;
	goto L138
L457:
	;
	v3611 = int32(0)
	goto L59
L458:
	;
	goto L457
L459:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L458
L460:
	;
	v3035 = v79 + int32(2144)
	v3036 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3035, int32(4), v3036)
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2164)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3035)+12)) = v79 + int32(2152)
	v3046 = int32(256)
	if v3046 < v3018 {
		goto L465
	} else {
		goto L466
	}
L461:
	;
	v3032 = v3018
	goto L463
L462:
	;
	v3032 = v3029
	goto L463
L463:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3012))) = v3032
	goto L460
L464:
	;
	v3052 = v79 + int32(2168)
	v3053 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3052, int32(4), v3053)
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2188)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3052)+12)) = v79 + int32(2176)
	v3063 = int32(256)
	if v3063 < v3018 {
		goto L469
	} else {
		goto L470
	}
L465:
	;
	v3049 = v3018
	goto L467
L466:
	;
	v3049 = v3046
	goto L467
L467:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3035))) = v3049
	goto L464
L468:
	;
	v3069 = v79 + int32(2192)
	v3070 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3069, int32(4), v3070)
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2212)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3069)+12)) = v79 + int32(2200)
	v3080 = int32(256)
	if v3080 < v3018 {
		goto L473
	} else {
		goto L474
	}
L469:
	;
	v3066 = v3018
	goto L471
L470:
	;
	v3066 = v3063
	goto L471
L471:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3052))) = v3066
	goto L468
L472:
	;
	v3085 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	if int32(1) <= v3085 {
		goto L477
	} else {
		goto L478
	}
L473:
	;
	v3083 = v3018
	goto L475
L474:
	;
	v3083 = v3080
	goto L475
L475:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3069))) = v3083
	goto L472
L476:
	;
	v3110 = v2906 - v3109
	if v3110 < int32(1) {
		goto L484
	} else {
		goto L485
	}
L477:
	;
	v3090 = int32(base.Ui32(v2906) >> (uint(int32(1)) % 32))
	if v2908 != 0 {
		goto L479
	} else {
		goto L480
	}
L478:
	;
	v3109 = int32(0)
	goto L476
L479:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+800)) = v3090
	v3109 = v3090
	goto L476
L480:
	;
	v3101 = int32(1)
	if base.Ui32(v3101) < base.Ui32(v3090) {
		goto L481
	} else {
		goto L482
	}
L481:
	;
	v3104 = v3090
	goto L483
L482:
	;
	v3104 = v3101
	goto L483
L483:
	;
	v3107 = F_memcpy(m, v56+int32(408), v56+int32(1280)+v3090*int32(-28)+v2906*int32(28), v3104*int32(28))
	mBase = m.M
	goto L479
L484:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1224)) = v2907
	*(*int32)(unsafe.Add(mBase, uint32(v56)+812)) = l0
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1220)) = v3110
	*(*int32)(unsafe.Add(mBase, uint32(v56)+820)) = l2
	*(*int32)(unsafe.Add(mBase, uint32(v56)+816)) = l1
	v3128 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1228)) = v3128
	*(*int32)(unsafe.Add(mBase, uint32(v56)+824)) = v79
	v3133 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[4])))
	m.T0[v3133].(func(*base.Module, int32))(m, v56+int32(1256))
	mBase = m.M
	v3135 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1272)) = v3135
	v3138 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1264)) = v3138 + int32(326)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1268)) = v56 + int32(812)
	if v3109 == v3135 {
		v3503 = v3135
		goto L486
	} else {
		goto L487
	}
L485:
	;
	v3122 = F_memcpy(m, v56+int32(828), v56+int32(1280), v3109*int32(-28)+v2906*int32(28))
	mBase = m.M
	goto L484
L486:
	;
	if v3109 == int32(0) {
		goto L567
	} else {
		goto L568
	}
L487:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+804)) = v2907
	*(*int32)(unsafe.Add(mBase, uint32(v56)+392)) = l0
	v3149 = int32(0)
	v3152 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v3153 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v3155 = v56 + int32(8)
	if l1 == v3149 {
		goto L489
	} else {
		goto L490
	}
L488:
	;
	v3253 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+104)) = v3253
	v3258 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v3258 != 0 {
		goto L509
	} else {
		goto L510
	}
L489:
	;
	goto L488
L490:
	;
	if v3155 == int32(0) {
		goto L489
	} else {
		goto L491
	}
L491:
	;
	v3166 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	if v3166 != 0 {
		goto L492
	} else {
		goto L493
	}
L492:
	;
	v3167 = v3149
	goto L494
L493:
	;
	v3167 = int32(0)
	goto L494
L494:
	;
	if v3166 != 0 {
		goto L495
	} else {
		goto L496
	}
L495:
	;
	v3170 = v3149
	goto L497
L496:
	;
	v3170 = int32(0)
	goto L497
L497:
	;
	if v3167|v3170 < int32(0) {
		goto L489
	} else {
		goto L498
	}
L498:
	;
	if v3152 < int32(1) {
		goto L489
	} else {
		goto L499
	}
L499:
	;
	if v3153 < int32(1) {
		goto L489
	} else {
		goto L500
	}
L500:
	;
	v3179 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	if v3179 < v3167+v3152 {
		goto L489
	} else {
		goto L501
	}
L501:
	;
	v3182 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	if v3182 < v3170+v3153 {
		goto L489
	} else {
		goto L502
	}
L502:
	;
	if l1 == v3155 {
		v3189 = v3166
		goto L503
	} else {
		goto L504
	}
L503:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3155)+12)) = v3153
	*(*int32)(unsafe.Add(mBase, uint32(v3155)+8)) = v3152
	if v3189 != 0 {
		goto L506
	} else {
		goto L507
	}
L504:
	;
	v3186 = F_memcpy(m, v3155, l1, int32(172))
	mBase = m.M
	F_WebPPictureResetBuffers(m, v3186)
	mBase = m.M
	v3188 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v3189 = v3188
	goto L503
L505:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3155+v3237))) = v3235
	*(*int32)(unsafe.Add(mBase, uint32(v3155+v3236))) = v3240
	goto L489
L506:
	;
	v3224 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	v3225 = *(*int32)(unsafe.Add(mBase, uint32(l1)+56))
	v3227 = int32(2)
	v3235 = v3225
	v3236 = int32(52)
	v3237 = int32(56)
	v3240 = v3224 + v3225*v3170<<(uint(v3227)%32) + v3167<<(uint(v3227)%32)
	goto L505
L507:
	;
	v3192 = *(*int32)(unsafe.Add(mBase, uint32(l1)+32))
	*(*int32)(unsafe.Add(mBase, uint32(v3155)+32)) = v3192
	v3194 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v3155)+28)) = v3194
	v3196 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3155)+16)) = v3196 + v3194*v3170 + v3167
	v3201 = int32(1)
	v3202 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v3205 = v3192 * int32(base.Ui32(v3170)>>(uint(v3201)%32))
	v3208 = int32(base.Ui32(v3167) >> (uint(v3201) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(v3155)+20)) = v3202 + v3205 + v3208
	v3211 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	*(*int32)(unsafe.Add(mBase, uint32(v3155)+24)) = v3211 + v3205 + v3208
	v3215 = *(*int32)(unsafe.Add(mBase, uint32(l1)+36))
	if v3215 == int32(0) {
		goto L489
	} else {
		goto L508
	}
L508:
	;
	v3218 = *(*int32)(unsafe.Add(mBase, uint32(l1)+40))
	v3235 = v3218
	v3236 = int32(36)
	v3237 = int32(40)
	v3240 = v3215 + v3218*v3170 + v3167
	goto L505
L509:
	;
	v3259 = v56 + int32(204)
	goto L511
L510:
	;
	v3259 = v3253
	goto L511
L511:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+808)) = v3259
	*(*int32)(unsafe.Add(mBase, uint32(v56)+396)) = v56 + int32(8)
	v3265 = v56 + int32(180)
	v3275 = *(*int32)(unsafe.Add(mBase, uint32(v3265)+12))
	v3276 = *(*int32)(unsafe.Add(mBase, uint32(v3265)+8))
	v3277 = v3275 - v3276
	v3279 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	v3280 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v3281 = v3279 - v3280
	v3283 = base.I64_extend_i32_u(v3277) + base.I64_extend_i32_u(v3281)
	if base.Ui64(v3283) < base.Ui64(int64(4294967296)) {
		goto L514
	} else {
		goto L515
	}
L512:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+400)) = v56 + int32(180)
	v3343 = int64(1)
	v3344 = int32(2224)
	goto L537
L513:
	;
	if v3335 != 0 {
		goto L512
	} else {
		goto L527
	}
L514:
	;
	v3289 = base.I32_wrap_i64(v3283)
	v3290 = *(*int32)(unsafe.Add(mBase, uint32(v3265)+16))
	v3291 = v3290 - v3276
	if v3290 == v3276 {
		goto L517
	} else {
		goto L518
	}
L515:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3265)+20)) = int32(1)
	v3335 = int32(0)
	goto L513
L516:
	;
	v3322 = F_memcpy(m, v3321, v3320, v3281)
	mBase = m.M
	v3323 = *(*int64)(unsafe.Add(mBase, uint32(l2)))
	*(*int64)(unsafe.Add(mBase, uint32(v3265))) = v3323
	v3325 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v3265)+20)) = v3325
	v3327 = *(*int32)(unsafe.Add(mBase, uint32(v3265)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3265)+12)) = v3327 + v3281
	v3335 = int32(1)
	goto L513
L517:
	;
	v3298 = int32(base.Ui32(v3291*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v3289) < base.Ui32(v3298) {
		goto L521
	} else {
		goto L522
	}
L518:
	;
	if base.Ui32(v3291) < base.Ui32(v3289) {
		goto L517
	} else {
		goto L519
	}
L519:
	;
	v3320 = v3280
	v3321 = v3276
	goto L516
L520:
	;
	if v3275 == v3276 {
		goto L525
	} else {
		goto L526
	}
L521:
	;
	v3300 = v3298
	goto L523
L522:
	;
	v3300 = v3289
	goto L523
L523:
	;
	v3304 = v3300&int32(-1024) + int32(1024)
	v3305 = F_WebPSafeMalloc(m, int64(1), v3304)
	mBase = m.M
	if v3305 != 0 {
		goto L520
	} else {
		goto L524
	}
L524:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3265)+20)) = int32(1)
	v3335 = int32(0)
	goto L513
L525:
	;
	v3312 = *(*int32)(unsafe.Add(mBase, uint32(v3265)+8))
	F_WebPSafeFree(m, v3312)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3265)+8)) = v3305
	*(*int32)(unsafe.Add(mBase, uint32(v3265)+16)) = v3305 + v3304
	*(*int32)(unsafe.Add(mBase, uint32(v3265)+12)) = v3305 + v3277
	v3319 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v3320 = v3319
	v3321 = v3305
	goto L516
L526:
	;
	v3310 = *(*int32)(unsafe.Add(mBase, uint32(v3265)+8))
	v3311 = F_memcpy(m, v3305, v3310, v3277)
	mBase = m.M
	goto L525
L527:
	;
	v3337 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3337 != 0 {
		goto L529
	} else {
		goto L530
	}
L528:
	;
	v3611 = v3149
	goto L59
L529:
	;
	goto L528
L530:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L529
L531:
	;
	v3400 = v3364 + int32(2120)
	v3404 = base.I32_div_s(v3384+int32(-1), int32(16))
	v3406 = v3404 + int32(1)
	v3407 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3400, int32(4), v3407)
	*(*int32)(unsafe.Add(mBase, uint32(v3364+int32(2140)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3400)+12)) = v3364 + int32(2128)
	v3417 = int32(256)
	if v3417 < v3406 {
		goto L552
	} else {
		goto L553
	}
L532:
	;
	v3396 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3396 != 0 {
		goto L549
	} else {
		goto L550
	}
L533:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3364)+12)) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3364))) = l0
	*(*int32)(unsafe.Add(mBase, uint32(v3364)+4)) = v56 + int32(8)
	F_VP8LEncDspInit(m)
	mBase = m.M
	v3380 = v3364 + int32(2216)
	v3381 = *(*int32)(unsafe.Add(mBase, uint32(v3364)+4))
	v3382 = *(*int32)(unsafe.Add(mBase, uint32(v3381)+12))
	v3383 = *(*int32)(unsafe.Add(mBase, uint32(v3381)+8))
	v3384 = v3382 * v3383
	v3388 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v3384), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3380))) = v3388
	if v3388 != 0 {
		goto L545
	} else {
		goto L546
	}
L534:
	;
	if v3364 != 0 {
		goto L533
	} else {
		goto L540
	}
L535:
	;
	goto L534
L536:
	;
	v3362 = F_calloc(m, base.I32_wrap_i64(v3343), v3344)
	mBase = m.M
	v3364 = v3362
	goto L535
L537:
	;
	v3351 = base.I64_div_u_s(int64(2147418112), v3343)
	v3352 = int32(0)
	v3353 = base.I64_extend_i32_u(v3344)
	if base.Ui64(int64(4294967295)) < base.Ui64(v3353*v3343) {
		v3364 = v3352
		goto L535
	} else {
		goto L538
	}
L538:
	;
	if base.Ui64(v3351) < base.Ui64(v3353) {
		v3364 = v3352
		goto L535
	} else {
		goto L539
	}
L539:
	;
	goto L536
L540:
	;
	v3367 = v56 + int32(8)
	v3369 = *(*int32)(unsafe.Add(mBase, uint32(v3367)+92))
	if v3369 != 0 {
		goto L542
	} else {
		goto L543
	}
L541:
	;
	goto L532
L542:
	;
	goto L541
L543:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3367)+92)) = int32(1)
	goto L542
L544:
	;
	if v3388 != 0 {
		goto L531
	} else {
		goto L547
	}
L545:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3380)+4)) = v3384
	goto L544
L546:
	;
	goto L544
L547:
	;
	goto L532
L548:
	;
	v3611 = v3364
	goto L59
L549:
	;
	goto L548
L550:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L549
L551:
	;
	v3423 = v3364 + int32(2144)
	v3424 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3423, int32(4), v3424)
	*(*int32)(unsafe.Add(mBase, uint32(v3364+int32(2164)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3423)+12)) = v3364 + int32(2152)
	v3434 = int32(256)
	if v3434 < v3406 {
		goto L556
	} else {
		goto L557
	}
L552:
	;
	v3420 = v3406
	goto L554
L553:
	;
	v3420 = v3417
	goto L554
L554:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3400))) = v3420
	goto L551
L555:
	;
	v3440 = v3364 + int32(2168)
	v3441 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3440, int32(4), v3441)
	*(*int32)(unsafe.Add(mBase, uint32(v3364+int32(2188)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3440)+12)) = v3364 + int32(2176)
	v3451 = int32(256)
	if v3451 < v3406 {
		goto L560
	} else {
		goto L561
	}
L556:
	;
	v3437 = v3406
	goto L558
L557:
	;
	v3437 = v3434
	goto L558
L558:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3423))) = v3437
	goto L555
L559:
	;
	v3457 = v3364 + int32(2192)
	v3458 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3457, int32(4), v3458)
	*(*int32)(unsafe.Add(mBase, uint32(v3364+int32(2212)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3457)+12)) = v3364 + int32(2200)
	v3468 = int32(256)
	if v3468 < v3406 {
		goto L564
	} else {
		goto L565
	}
L560:
	;
	v3454 = v3406
	goto L562
L561:
	;
	v3454 = v3451
	goto L562
L562:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3440))) = v3454
	goto L559
L563:
	;
	v3473 = *(*int64)(unsafe.Add(mBase, uint32(v79)+36))
	*(*int64)(unsafe.Add(mBase, uint32(v3364)+36)) = v3473
	v3475 = *(*int32)(unsafe.Add(mBase, uint32(v79)+44))
	*(*int32)(unsafe.Add(mBase, uint32(v3364)+44)) = v3475
	v3477 = *(*int32)(unsafe.Add(mBase, uint32(v79)+68))
	*(*int32)(unsafe.Add(mBase, uint32(v3364)+68)) = v3477
	v3479 = int32(72)
	v3483 = int32(1024)
	v3484 = F_memcpy(m, v3364+v3479, v79+v3479, v3483)
	mBase = m.M
	v3488 = F_memcpy(m, v3364+int32(1096), v289, v3483)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v56)+404)) = v3364
	v3492 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[4])))
	m.T0[v3492].(func(*base.Module, int32))(m, v56+int32(1232))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1248)) = int32(0)
	v3496 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1240)) = v3496 + int32(326)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1244)) = v56 + int32(392)
	v3503 = v3364
	goto L486
L564:
	;
	v3471 = v3406
	goto L566
L565:
	;
	v3471 = v3468
	goto L566
L566:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3457))) = v3471
	goto L563
L567:
	;
	v3528 = v56 + int32(1256)
	v3529 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[5])))
	m.T0[v3529].(func(*base.Module, int32))(m, v3528)
	mBase = m.M
	v3533 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[6])))
	v3534 = m.T0[v3533].(func(*base.Module, int32) int32)(m, v3528)
	mBase = m.M
	v3537 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[7])))
	m.T0[v3537].(func(*base.Module, int32))(m, v3528)
	mBase = m.M
	if v3109 == int32(0) {
		v3611 = v3503
		goto L59
	} else {
		goto L576
	}
L568:
	;
	v3509 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[8])))
	v3510 = m.T0[v3509].(func(*base.Module, int32) int32)(m, v56+int32(1232))
	mBase = m.M
	if v3510 != 0 {
		goto L569
	} else {
		goto L570
	}
L569:
	;
	v3515 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v3515 == int32(0) {
		goto L574
	} else {
		goto L575
	}
L570:
	;
	v3512 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3512 != 0 {
		goto L572
	} else {
		goto L573
	}
L571:
	;
	v3611 = v3503
	goto L59
L572:
	;
	goto L571
L573:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L572
L574:
	;
	v3524 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[9])))
	m.T0[v3524].(func(*base.Module, int32))(m, v56+int32(1232))
	mBase = m.M
	goto L567
L575:
	;
	v3521 = F_memcpy(m, v56+int32(204), v3515, int32(188))
	mBase = m.M
	goto L574
L576:
	;
	v3542 = v56 + int32(1232)
	v3543 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[6])))
	v3544 = m.T0[v3543].(func(*base.Module, int32) int32)(m, v3542)
	mBase = m.M
	v3547 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[7])))
	m.T0[v3547].(func(*base.Module, int32))(m, v3542)
	mBase = m.M
	if v3534 == int32(0) {
		goto L578
	} else {
		goto L579
	}
L577:
	;
	v3556 = *(*int32)(unsafe.Add(mBase, uint32(v56)+192))
	v3557 = *(*int32)(unsafe.Add(mBase, uint32(v56)+188))
	v3559 = *(*int32)(unsafe.Add(mBase, uint32(v56)+184))
	v3560 = int32(7)
	v3562 = int32(3)
	v3565 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	v3566 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v3568 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	if base.Ui32(v3565-v3566+(v3568+v3560)>>(uint(v3562)%32)) <= base.Ui32(v3556-v3557+(v3559+v3560)>>(uint(v3562)%32)) {
		v3611 = v3503
		goto L59
	} else {
		goto L585
	}
L578:
	;
	v3551 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3551 != 0 {
		v3611 = v3503
		goto L59
	} else {
		goto L581
	}
L579:
	;
	if v3544 != 0 {
		goto L577
	} else {
		goto L580
	}
L580:
	;
	goto L578
L581:
	;
	v3552 = *(*int32)(unsafe.Add(mBase, uint32(v56)+100))
	v3553 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3553 != 0 {
		goto L583
	} else {
		goto L584
	}
L582:
	;
	v3611 = v3503
	goto L59
L583:
	;
	goto L582
L584:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = v3552
	goto L583
L585:
	;
	v3576 = v56 + int32(180)
	v3579 = int32(0)
	v3583 = *(*int64)(unsafe.Add(mBase, uint32(l2)))
	v3584 = *(*int64)(unsafe.Add(mBase, uint32(v3576)))
	*(*int64)(unsafe.Add(mBase, uint32(l2))) = v3584
	v3587 = l2 + int32(8)
	v3589 = base.Simd_g_v128_load(m, v3587, v3579)
	v3591 = v56 + int32(188)
	v3593 = base.Simd_g_v128_load(m, v3591, v3579)
	base.Simd_g_v128_store(m, v3587, v3579, v3593)
	*(*int64)(unsafe.Add(mBase, uint32(v3576))) = v3583
	base.Simd_g_v128_store(m, v3591, v3579, v3589)
	goto L586
L586:
	;
	v3599 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v3599 == int32(0) {
		v3611 = v3503
		goto L59
	} else {
		goto L587
	}
L587:
	;
	v3605 = F_memcpy(m, v3599, v56+int32(204), int32(188))
	mBase = m.M
	v3611 = v3503
	goto L59
L588:
	;
	v3675 = v79 + int32(2216)
	v3676 = *(*int32)(unsafe.Add(mBase, uint32(v3675)))
	F_WebPSafeFree(m, v3676)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3675))) = int64(0)
	goto L591
L589:
	;
	goto L588
L590:
	;
	v3665 = *(*int32)(unsafe.Add(mBase, uint32(v56+int32(188))))
	F_WebPSafeFree(m, v3665)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v56+int32(196)))) = int64(0)
	v3671 = base.Simd_g_const(&F_VP8LEncodeStream__k0)
	base.Simd_g_v128_store(m, v3660, int32(0), v3671)
	goto L589
L591:
	;
	v3681 = v79 + int32(2120)
	v3684 = *(*int32)(unsafe.Add(mBase, uint32(v3681)+12))
	if v3684 == int32(0) {
		goto L593
	} else {
		goto L594
	}
L592:
	;
	v3710 = v79 + int32(2144)
	v3713 = *(*int32)(unsafe.Add(mBase, uint32(v3710)+12))
	if v3713 == int32(0) {
		goto L601
	} else {
		goto L602
	}
L593:
	;
	v3689 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3681)+20)) = v3689
	v3691 = *(*int32)(unsafe.Add(mBase, uint32(v3681)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3681)+16)) = v3691
	*(*int32)(unsafe.Add(mBase, uint32(v3681)+12)) = v79 + int32(2128)
	*(*int32)(unsafe.Add(mBase, uint32(v3681)+8)) = v3689
	if v3691 == v3689 {
		goto L595
	} else {
		goto L596
	}
L594:
	;
	v3687 = *(*int32)(unsafe.Add(mBase, uint32(v3681)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3684))) = v3687
	goto L593
L595:
	;
	goto L592
L596:
	;
	v3702 = v3691
	goto L597
L597:
	;
	v3703 = *(*int32)(unsafe.Add(mBase, uint32(v3702)))
	F_WebPSafeFree(m, v3702)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3681)+16)) = v3703
	if v3703 != 0 {
		v3702 = v3703
		goto L597
	} else {
		goto L599
	}
L598:
	;
	goto L595
L599:
	;
	goto L598
L600:
	;
	v3739 = v79 + int32(2168)
	v3742 = *(*int32)(unsafe.Add(mBase, uint32(v3739)+12))
	if v3742 == int32(0) {
		goto L609
	} else {
		goto L610
	}
L601:
	;
	v3718 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3710)+20)) = v3718
	v3720 = *(*int32)(unsafe.Add(mBase, uint32(v3710)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3710)+16)) = v3720
	*(*int32)(unsafe.Add(mBase, uint32(v3710)+12)) = v79 + int32(2152)
	*(*int32)(unsafe.Add(mBase, uint32(v3710)+8)) = v3718
	if v3720 == v3718 {
		goto L603
	} else {
		goto L604
	}
L602:
	;
	v3716 = *(*int32)(unsafe.Add(mBase, uint32(v3710)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3713))) = v3716
	goto L601
L603:
	;
	goto L600
L604:
	;
	v3731 = v3720
	goto L605
L605:
	;
	v3732 = *(*int32)(unsafe.Add(mBase, uint32(v3731)))
	F_WebPSafeFree(m, v3731)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3710)+16)) = v3732
	if v3732 != 0 {
		v3731 = v3732
		goto L605
	} else {
		goto L607
	}
L606:
	;
	goto L603
L607:
	;
	goto L606
L608:
	;
	v3768 = v79 + int32(2192)
	v3771 = *(*int32)(unsafe.Add(mBase, uint32(v3768)+12))
	if v3771 == int32(0) {
		goto L617
	} else {
		goto L618
	}
L609:
	;
	v3747 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3739)+20)) = v3747
	v3749 = *(*int32)(unsafe.Add(mBase, uint32(v3739)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3739)+16)) = v3749
	*(*int32)(unsafe.Add(mBase, uint32(v3739)+12)) = v79 + int32(2176)
	*(*int32)(unsafe.Add(mBase, uint32(v3739)+8)) = v3747
	if v3749 == v3747 {
		goto L611
	} else {
		goto L612
	}
L610:
	;
	v3745 = *(*int32)(unsafe.Add(mBase, uint32(v3739)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3742))) = v3745
	goto L609
L611:
	;
	goto L608
L612:
	;
	v3760 = v3749
	goto L613
L613:
	;
	v3761 = *(*int32)(unsafe.Add(mBase, uint32(v3760)))
	F_WebPSafeFree(m, v3760)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3739)+16)) = v3761
	if v3761 != 0 {
		v3760 = v3761
		goto L613
	} else {
		goto L615
	}
L614:
	;
	goto L611
L615:
	;
	goto L614
L616:
	;
	v3796 = *(*int32)(unsafe.Add(mBase, uint32(v79)+24))
	F_free(m, v3796)
	mBase = m.M
	goto L624
L617:
	;
	v3776 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3768)+20)) = v3776
	v3778 = *(*int32)(unsafe.Add(mBase, uint32(v3768)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3768)+16)) = v3778
	*(*int32)(unsafe.Add(mBase, uint32(v3768)+12)) = v79 + int32(2200)
	*(*int32)(unsafe.Add(mBase, uint32(v3768)+8)) = v3776
	if v3778 == v3776 {
		goto L619
	} else {
		goto L620
	}
L618:
	;
	v3774 = *(*int32)(unsafe.Add(mBase, uint32(v3768)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3771))) = v3774
	goto L617
L619:
	;
	goto L616
L620:
	;
	v3789 = v3778
	goto L621
L621:
	;
	v3790 = *(*int32)(unsafe.Add(mBase, uint32(v3789)))
	F_WebPSafeFree(m, v3789)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3768)+16)) = v3790
	if v3790 != 0 {
		v3789 = v3790
		goto L621
	} else {
		goto L623
	}
L622:
	;
	goto L619
L623:
	;
	goto L622
L624:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v79)+24)) = int64(0)
	F_free(m, v79)
	mBase = m.M
	goto L625
L625:
	;
	if v3611 == int32(0) {
		goto L626
	} else {
		goto L627
	}
L626:
	;
	v3930 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	v3938 = base.B2i32(v3930 == int32(0))
	goto L1
L627:
	;
	v3804 = v3611 + int32(2216)
	v3805 = *(*int32)(unsafe.Add(mBase, uint32(v3804)))
	F_WebPSafeFree(m, v3805)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3804))) = int64(0)
	goto L628
L628:
	;
	v3810 = v3611 + int32(2120)
	v3813 = *(*int32)(unsafe.Add(mBase, uint32(v3810)+12))
	if v3813 == int32(0) {
		goto L630
	} else {
		goto L631
	}
L629:
	;
	v3839 = v3611 + int32(2144)
	v3842 = *(*int32)(unsafe.Add(mBase, uint32(v3839)+12))
	if v3842 == int32(0) {
		goto L638
	} else {
		goto L639
	}
L630:
	;
	v3818 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3810)+20)) = v3818
	v3820 = *(*int32)(unsafe.Add(mBase, uint32(v3810)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3810)+16)) = v3820
	*(*int32)(unsafe.Add(mBase, uint32(v3810)+12)) = v3611 + int32(2128)
	*(*int32)(unsafe.Add(mBase, uint32(v3810)+8)) = v3818
	if v3820 == v3818 {
		goto L632
	} else {
		goto L633
	}
L631:
	;
	v3816 = *(*int32)(unsafe.Add(mBase, uint32(v3810)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3813))) = v3816
	goto L630
L632:
	;
	goto L629
L633:
	;
	v3831 = v3820
	goto L634
L634:
	;
	v3832 = *(*int32)(unsafe.Add(mBase, uint32(v3831)))
	F_WebPSafeFree(m, v3831)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3810)+16)) = v3832
	if v3832 != 0 {
		v3831 = v3832
		goto L634
	} else {
		goto L636
	}
L635:
	;
	goto L632
L636:
	;
	goto L635
L637:
	;
	v3868 = v3611 + int32(2168)
	v3871 = *(*int32)(unsafe.Add(mBase, uint32(v3868)+12))
	if v3871 == int32(0) {
		goto L646
	} else {
		goto L647
	}
L638:
	;
	v3847 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3839)+20)) = v3847
	v3849 = *(*int32)(unsafe.Add(mBase, uint32(v3839)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3839)+16)) = v3849
	*(*int32)(unsafe.Add(mBase, uint32(v3839)+12)) = v3611 + int32(2152)
	*(*int32)(unsafe.Add(mBase, uint32(v3839)+8)) = v3847
	if v3849 == v3847 {
		goto L640
	} else {
		goto L641
	}
L639:
	;
	v3845 = *(*int32)(unsafe.Add(mBase, uint32(v3839)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3842))) = v3845
	goto L638
L640:
	;
	goto L637
L641:
	;
	v3860 = v3849
	goto L642
L642:
	;
	v3861 = *(*int32)(unsafe.Add(mBase, uint32(v3860)))
	F_WebPSafeFree(m, v3860)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3839)+16)) = v3861
	if v3861 != 0 {
		v3860 = v3861
		goto L642
	} else {
		goto L644
	}
L643:
	;
	goto L640
L644:
	;
	goto L643
L645:
	;
	v3897 = v3611 + int32(2192)
	v3900 = *(*int32)(unsafe.Add(mBase, uint32(v3897)+12))
	if v3900 == int32(0) {
		goto L654
	} else {
		goto L655
	}
L646:
	;
	v3876 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3868)+20)) = v3876
	v3878 = *(*int32)(unsafe.Add(mBase, uint32(v3868)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3868)+16)) = v3878
	*(*int32)(unsafe.Add(mBase, uint32(v3868)+12)) = v3611 + int32(2176)
	*(*int32)(unsafe.Add(mBase, uint32(v3868)+8)) = v3876
	if v3878 == v3876 {
		goto L648
	} else {
		goto L649
	}
L647:
	;
	v3874 = *(*int32)(unsafe.Add(mBase, uint32(v3868)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3871))) = v3874
	goto L646
L648:
	;
	goto L645
L649:
	;
	v3889 = v3878
	goto L650
L650:
	;
	v3890 = *(*int32)(unsafe.Add(mBase, uint32(v3889)))
	F_WebPSafeFree(m, v3889)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3868)+16)) = v3890
	if v3890 != 0 {
		v3889 = v3890
		goto L650
	} else {
		goto L652
	}
L651:
	;
	goto L648
L652:
	;
	goto L651
L653:
	;
	v3925 = *(*int32)(unsafe.Add(mBase, uint32(v3611)+24))
	F_free(m, v3925)
	mBase = m.M
	goto L661
L654:
	;
	v3905 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3897)+20)) = v3905
	v3907 = *(*int32)(unsafe.Add(mBase, uint32(v3897)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3897)+16)) = v3907
	*(*int32)(unsafe.Add(mBase, uint32(v3897)+12)) = v3611 + int32(2200)
	*(*int32)(unsafe.Add(mBase, uint32(v3897)+8)) = v3905
	if v3907 == v3905 {
		goto L656
	} else {
		goto L657
	}
L655:
	;
	v3903 = *(*int32)(unsafe.Add(mBase, uint32(v3897)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3900))) = v3903
	goto L654
L656:
	;
	goto L653
L657:
	;
	v3918 = v3907
	goto L658
L658:
	;
	v3919 = *(*int32)(unsafe.Add(mBase, uint32(v3918)))
	F_WebPSafeFree(m, v3918)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3897)+16)) = v3919
	if v3919 != 0 {
		v3918 = v3919
		goto L658
	} else {
		goto L660
	}
L659:
	;
	goto L656
L660:
	;
	goto L659
L661:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3611)+24)) = int64(0)
	F_free(m, v3611)
	mBase = m.M
	goto L662
L662:
	;
	goto L626
}

var F_VP8LEncodeStream__k0 = [2]uint64{0x0, 0x0}

func F_VP8LResidualImage(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32, l14 int32) int32 {
	mBase := m.M
	_ = mBase
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v104 int32
	_ = v104
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v121 int32
	_ = v121
	var v128 int32
	_ = v128
	var v137 int32
	_ = v137
	var v152 int32
	_ = v152
	var v226 base.V128
	_ = v226
	var v232 int32
	_ = v232
	var v256 int32
	_ = v256
	var v343 int32
	_ = v343
	var v358 int32
	_ = v358
	var v437 int32
	_ = v437
	var v538 int32
	_ = v538
	var v539 int32
	_ = v539
	var v541 int32
	_ = v541
	var v544 int32
	_ = v544
	var v545 int32
	_ = v545
	var v552 int32
	_ = v552
	var v583 int32
	_ = v583
	var v587 base.V128
	_ = v587
	var v591 base.V128
	_ = v591
	var v592 int32
	_ = v592
	var v657 int32
	_ = v657
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v664 int32
	_ = v664
	var v669 int32
	_ = v669
	var v670 int32
	_ = v670
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v679 base.V128
	_ = v679
	var v680 base.V128
	_ = v680
	var v700 base.V128
	_ = v700
	var v720 base.V128
	_ = v720
	var v723 base.V128
	_ = v723
	var v727 base.V128
	_ = v727
	var v729 int32
	_ = v729
	var v732 base.V128
	_ = v732
	var v737 int32
	_ = v737
	var v747 int32
	_ = v747
	var v760 int32
	_ = v760
	var v847 int32
	_ = v847
	var v860 int32
	_ = v860
	var v862 int32
	_ = v862
	var v936 int32
	_ = v936
	var v937 int32
	_ = v937
	var v942 int32
	_ = v942
	var v944 int32
	_ = v944
	var v948 int32
	_ = v948
	var v971 int32
	_ = v971
	var v1047 int64
	_ = v1047
	var v1048 int32
	_ = v1048
	var v1055 int64
	_ = v1055
	var v1056 int32
	_ = v1056
	var v1057 int64
	_ = v1057
	var v1067 int32
	_ = v1067
	var v1069 int32
	_ = v1069
	var v1077 int32
	_ = v1077
	var v1082 int32
	_ = v1082
	var v1086 int32
	_ = v1086
	var v1101 int32
	_ = v1101
	var v1114 int32
	_ = v1114
	var v1115 int32
	_ = v1115
	var v1116 int32
	_ = v1116
	var v1126 int32
	_ = v1126
	var v1190 int32
	_ = v1190
	var v1193 int32
	_ = v1193
	var v1195 int32
	_ = v1195
	var v1200 int32
	_ = v1200
	var v1202 int32
	_ = v1202
	var v1211 int32
	_ = v1211
	var v1225 int32
	_ = v1225
	var v1304 int32
	_ = v1304
	var v1316 int32
	_ = v1316
	var v1330 int32
	_ = v1330
	var v1340 int32
	_ = v1340
	var v1341 int32
	_ = v1341
	var v1405 int32
	_ = v1405
	var v1406 int32
	_ = v1406
	var v1408 int32
	_ = v1408
	var v1409 int32
	_ = v1409
	var v1410 int32
	_ = v1410
	var v1412 int32
	_ = v1412
	var v1414 int32
	_ = v1414
	var v1418 int32
	_ = v1418
	var v1421 int32
	_ = v1421
	var v1423 int32
	_ = v1423
	var v1427 int32
	_ = v1427
	var v1430 int32
	_ = v1430
	var v1432 int32
	_ = v1432
	var v1436 int32
	_ = v1436
	var v1439 int32
	_ = v1439
	var v1446 int32
	_ = v1446
	var v1544 int32
	_ = v1544
	var v1546 int32
	_ = v1546
	var v1549 int64
	_ = v1549
	var v1550 int32
	_ = v1550
	var v1557 int64
	_ = v1557
	var v1558 int32
	_ = v1558
	var v1559 int64
	_ = v1559
	var v1568 int32
	_ = v1568
	var v1570 int32
	_ = v1570
	var v1572 int32
	_ = v1572
	var v1573 int32
	_ = v1573
	var v1577 int32
	_ = v1577
	var v1578 int32
	_ = v1578
	var v1580 int32
	_ = v1580
	var v1582 int32
	_ = v1582
	var v1584 int32
	_ = v1584
	var v1589 int32
	_ = v1589
	var v1593 int32
	_ = v1593
	var v1596 int32
	_ = v1596
	var v1598 int32
	_ = v1598
	var v1599 int32
	_ = v1599
	var v1601 int32
	_ = v1601
	var v1602 int32
	_ = v1602
	var v1605 int32
	_ = v1605
	var v1607 int32
	_ = v1607
	var v1609 int32
	_ = v1609
	var v1610 int32
	_ = v1610
	var v1613 int32
	_ = v1613
	var v1622 int32
	_ = v1622
	var v1627 int32
	_ = v1627
	var v1628 int32
	_ = v1628
	var v1631 int32
	_ = v1631
	var v1632 int32
	_ = v1632
	var v1696 int32
	_ = v1696
	var v1697 int32
	_ = v1697
	var v1698 int32
	_ = v1698
	var v1699 int32
	_ = v1699
	var v1700 int32
	_ = v1700
	var v1701 int32
	_ = v1701
	var v1735 int32
	_ = v1735
	var v1737 int32
	_ = v1737
	var v1738 int32
	_ = v1738
	var v1739 int32
	_ = v1739
	var v1740 int32
	_ = v1740
	var v1741 int32
	_ = v1741
	var v1742 int32
	_ = v1742
	var v1748 int32
	_ = v1748
	var v1749 int32
	_ = v1749
	var v1750 int32
	_ = v1750
	var v1751 int32
	_ = v1751
	var v1752 int32
	_ = v1752
	var v1756 int32
	_ = v1756
	var v1757 int32
	_ = v1757
	var v1759 int32
	_ = v1759
	var v1762 int32
	_ = v1762
	var v1763 int32
	_ = v1763
	var v1764 int32
	_ = v1764
	var v1776 int32
	_ = v1776
	var v1792 int32
	_ = v1792
	var v1843 int32
	_ = v1843
	var v1868 int32
	_ = v1868
	var v1871 int32
	_ = v1871
	var v1881 int32
	_ = v1881
	var __phi1881 int32
	_ = __phi1881
	var v1901 int32
	_ = v1901
	var __phi1901 int32
	_ = __phi1901
	var v1948 int32
	_ = v1948
	var __phi1948 int32
	_ = __phi1948
	var v1971 int32
	_ = v1971
	var v1973 int32
	_ = v1973
	var v1975 int32
	_ = v1975
	var v1977 int32
	_ = v1977
	var v1982 int32
	_ = v1982
	var v2005 int32
	_ = v2005
	var v2006 base.V128
	_ = v2006
	var v2010 base.V128
	_ = v2010
	var v2013 base.V128
	_ = v2013
	var v2023 base.V128
	_ = v2023
	var v2028 int32
	_ = v2028
	var v2032 int32
	_ = v2032
	var v2034 int32
	_ = v2034
	var v2040 int32
	_ = v2040
	var __phi2040 int32
	_ = __phi2040
	var v2041 int32
	_ = v2041
	var __phi2041 int32
	_ = __phi2041
	var v2048 int32
	_ = v2048
	var __phi2048 int32
	_ = __phi2048
	var v2049 int32
	_ = v2049
	var __phi2049 int32
	_ = __phi2049
	var v2050 int32
	_ = v2050
	var __phi2050 int32
	_ = __phi2050
	var v2061 int32
	_ = v2061
	var v2062 int32
	_ = v2062
	var v2064 int32
	_ = v2064
	var v2066 int32
	_ = v2066
	var v2069 int32
	_ = v2069
	var v2070 int32
	_ = v2070
	var v2071 int32
	_ = v2071
	var v2073 int32
	_ = v2073
	var v2076 int32
	_ = v2076
	var v2081 int32
	_ = v2081
	var v2085 int32
	_ = v2085
	var v2100 int32
	_ = v2100
	var v2114 int32
	_ = v2114
	var v2115 int32
	_ = v2115
	var v2116 int32
	_ = v2116
	var v2118 int32
	_ = v2118
	var v2119 int32
	_ = v2119
	var v2122 int32
	_ = v2122
	var v2123 int32
	_ = v2123
	var v2124 int32
	_ = v2124
	var v2126 int32
	_ = v2126
	var v2127 int32
	_ = v2127
	var v2129 int32
	_ = v2129
	var v2130 int32
	_ = v2130
	var v2135 int32
	_ = v2135
	var v2137 int32
	_ = v2137
	var v2139 int32
	_ = v2139
	var v2141 int32
	_ = v2141
	var v2142 int32
	_ = v2142
	var v2144 int32
	_ = v2144
	var v2145 int32
	_ = v2145
	var v2150 int32
	_ = v2150
	var v2151 int32
	_ = v2151
	var v2152 int32
	_ = v2152
	var v2154 int32
	_ = v2154
	var v2156 int32
	_ = v2156
	var v2159 int32
	_ = v2159
	var v2161 int32
	_ = v2161
	var v2163 int32
	_ = v2163
	var v2165 int32
	_ = v2165
	var v2167 int32
	_ = v2167
	var v2170 int32
	_ = v2170
	var v2171 int32
	_ = v2171
	var v2172 int32
	_ = v2172
	var v2174 int32
	_ = v2174
	var v2179 int32
	_ = v2179
	var v2181 int32
	_ = v2181
	var v2183 int32
	_ = v2183
	var v2185 int32
	_ = v2185
	var v2188 int32
	_ = v2188
	var v2190 int32
	_ = v2190
	var v2191 int32
	_ = v2191
	var v2192 int32
	_ = v2192
	var v2194 int32
	_ = v2194
	var v2197 int32
	_ = v2197
	var v2199 int32
	_ = v2199
	var v2201 int32
	_ = v2201
	var v2203 int32
	_ = v2203
	var v2205 int32
	_ = v2205
	var v2207 int32
	_ = v2207
	var v2210 int32
	_ = v2210
	var v2211 int32
	_ = v2211
	var v2212 int32
	_ = v2212
	var v2214 int32
	_ = v2214
	var v2219 int32
	_ = v2219
	var v2221 int32
	_ = v2221
	var v2223 int32
	_ = v2223
	var v2225 int32
	_ = v2225
	var v2228 int32
	_ = v2228
	var v2230 int32
	_ = v2230
	var v2231 int32
	_ = v2231
	var v2232 int32
	_ = v2232
	var v2234 int32
	_ = v2234
	var v2237 int32
	_ = v2237
	var v2239 int32
	_ = v2239
	var v2241 int32
	_ = v2241
	var v2243 int32
	_ = v2243
	var v2245 int32
	_ = v2245
	var v2248 int32
	_ = v2248
	var v2249 int32
	_ = v2249
	var v2250 int32
	_ = v2250
	var v2252 int32
	_ = v2252
	var v2257 int32
	_ = v2257
	var v2259 int32
	_ = v2259
	var v2261 int32
	_ = v2261
	var v2263 int32
	_ = v2263
	var v2266 int32
	_ = v2266
	var v2268 int32
	_ = v2268
	var v2269 int32
	_ = v2269
	var v2270 int32
	_ = v2270
	var v2272 int32
	_ = v2272
	var v2275 int32
	_ = v2275
	var v2277 int32
	_ = v2277
	var v2279 int32
	_ = v2279
	var v2281 int32
	_ = v2281
	var v2283 int32
	_ = v2283
	var v2285 int32
	_ = v2285
	var v2287 int32
	_ = v2287
	var v2292 int32
	_ = v2292
	var v2314 int32
	_ = v2314
	var v2339 int32
	_ = v2339
	var v2341 int32
	_ = v2341
	var v2415 int32
	_ = v2415
	var v2418 int32
	_ = v2418
	var v2420 int32
	_ = v2420
	var v2421 int32
	_ = v2421
	var v2422 int32
	_ = v2422
	var v2431 int32
	_ = v2431
	var v2432 int32
	_ = v2432
	var v2442 int32
	_ = v2442
	var v2443 int32
	_ = v2443
	var v2453 int32
	_ = v2453
	var v2454 int32
	_ = v2454
	var v2461 int32
	_ = v2461
	var v2496 int32
	_ = v2496
	var v2566 int32
	_ = v2566
	var v2590 int32
	_ = v2590
	var v2591 int32
	_ = v2591
	var v2666 int32
	_ = v2666
	var v2669 int32
	_ = v2669
	var v2671 int32
	_ = v2671
	var v2672 int32
	_ = v2672
	var v2673 int32
	_ = v2673
	var v2682 int32
	_ = v2682
	var v2683 int32
	_ = v2683
	var v2693 int32
	_ = v2693
	var v2694 int32
	_ = v2694
	var v2704 int32
	_ = v2704
	var v2705 int32
	_ = v2705
	var v2712 int32
	_ = v2712
	var v2714 int32
	_ = v2714
	var v2814 int32
	_ = v2814
	var v2817 int32
	_ = v2817
	var v2823 int32
	_ = v2823
	var v2853 int32
	_ = v2853
	var v2863 int32
	_ = v2863
	var v2864 int32
	_ = v2864
	var v2931 int32
	_ = v2931
	var v2940 int32
	_ = v2940
	var __phi2940 int32
	_ = __phi2940
	var v2954 int32
	_ = v2954
	var __phi2954 int32
	_ = __phi2954
	var v2955 int32
	_ = v2955
	var __phi2955 int32
	_ = __phi2955
	var v2965 int32
	_ = v2965
	var __phi2965 int32
	_ = __phi2965
	var v2966 int32
	_ = v2966
	var __phi2966 int32
	_ = __phi2966
	var v3031 int32
	_ = v3031
	var v3036 int32
	_ = v3036
	var v3042 int32
	_ = v3042
	var v3044 int32
	_ = v3044
	var v3048 int32
	_ = v3048
	var v3077 int32
	_ = v3077
	var v3078 int32
	_ = v3078
	var v3088 int32
	_ = v3088
	var v3156 int32
	_ = v3156
	var v3165 int32
	_ = v3165
	var __phi3165 int32
	_ = __phi3165
	var v3178 int32
	_ = v3178
	var __phi3178 int32
	_ = __phi3178
	var v3179 int32
	_ = v3179
	var __phi3179 int32
	_ = __phi3179
	var v3190 int32
	_ = v3190
	var __phi3190 int32
	_ = __phi3190
	var v3191 int32
	_ = v3191
	var __phi3191 int32
	_ = __phi3191
	var v3255 int32
	_ = v3255
	var v3256 int32
	_ = v3256
	var v3261 int32
	_ = v3261
	var v3282 int32
	_ = v3282
	var v3283 base.V128
	_ = v3283
	var v3287 base.V128
	_ = v3287
	var v3290 base.V128
	_ = v3290
	var v3300 base.V128
	_ = v3300
	var v3305 int32
	_ = v3305
	var v3309 int32
	_ = v3309
	var v3311 int32
	_ = v3311
	var v3317 int32
	_ = v3317
	var __phi3317 int32
	_ = __phi3317
	var v3318 int32
	_ = v3318
	var __phi3318 int32
	_ = __phi3318
	var v3325 int32
	_ = v3325
	var __phi3325 int32
	_ = __phi3325
	var v3326 int32
	_ = v3326
	var __phi3326 int32
	_ = __phi3326
	var v3327 int32
	_ = v3327
	var __phi3327 int32
	_ = __phi3327
	var v3338 int32
	_ = v3338
	var v3339 int32
	_ = v3339
	var v3341 int32
	_ = v3341
	var v3343 int32
	_ = v3343
	var v3346 int32
	_ = v3346
	var v3347 int32
	_ = v3347
	var v3348 int32
	_ = v3348
	var v3350 int32
	_ = v3350
	var v3353 int32
	_ = v3353
	var v3358 int32
	_ = v3358
	var v3362 int32
	_ = v3362
	var v3377 int32
	_ = v3377
	var v3391 int32
	_ = v3391
	var v3392 int32
	_ = v3392
	var v3393 int32
	_ = v3393
	var v3395 int32
	_ = v3395
	var v3396 int32
	_ = v3396
	var v3399 int32
	_ = v3399
	var v3400 int32
	_ = v3400
	var v3401 int32
	_ = v3401
	var v3403 int32
	_ = v3403
	var v3404 int32
	_ = v3404
	var v3406 int32
	_ = v3406
	var v3407 int32
	_ = v3407
	var v3412 int32
	_ = v3412
	var v3414 int32
	_ = v3414
	var v3416 int32
	_ = v3416
	var v3418 int32
	_ = v3418
	var v3419 int32
	_ = v3419
	var v3421 int32
	_ = v3421
	var v3422 int32
	_ = v3422
	var v3427 int32
	_ = v3427
	var v3428 int32
	_ = v3428
	var v3429 int32
	_ = v3429
	var v3431 int32
	_ = v3431
	var v3433 int32
	_ = v3433
	var v3436 int32
	_ = v3436
	var v3438 int32
	_ = v3438
	var v3440 int32
	_ = v3440
	var v3442 int32
	_ = v3442
	var v3444 int32
	_ = v3444
	var v3447 int32
	_ = v3447
	var v3448 int32
	_ = v3448
	var v3449 int32
	_ = v3449
	var v3451 int32
	_ = v3451
	var v3456 int32
	_ = v3456
	var v3458 int32
	_ = v3458
	var v3460 int32
	_ = v3460
	var v3462 int32
	_ = v3462
	var v3465 int32
	_ = v3465
	var v3467 int32
	_ = v3467
	var v3468 int32
	_ = v3468
	var v3469 int32
	_ = v3469
	var v3471 int32
	_ = v3471
	var v3474 int32
	_ = v3474
	var v3476 int32
	_ = v3476
	var v3478 int32
	_ = v3478
	var v3480 int32
	_ = v3480
	var v3482 int32
	_ = v3482
	var v3484 int32
	_ = v3484
	var v3487 int32
	_ = v3487
	var v3488 int32
	_ = v3488
	var v3489 int32
	_ = v3489
	var v3491 int32
	_ = v3491
	var v3496 int32
	_ = v3496
	var v3498 int32
	_ = v3498
	var v3500 int32
	_ = v3500
	var v3502 int32
	_ = v3502
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
	var v3514 int32
	_ = v3514
	var v3516 int32
	_ = v3516
	var v3518 int32
	_ = v3518
	var v3520 int32
	_ = v3520
	var v3522 int32
	_ = v3522
	var v3525 int32
	_ = v3525
	var v3526 int32
	_ = v3526
	var v3527 int32
	_ = v3527
	var v3529 int32
	_ = v3529
	var v3534 int32
	_ = v3534
	var v3536 int32
	_ = v3536
	var v3538 int32
	_ = v3538
	var v3540 int32
	_ = v3540
	var v3543 int32
	_ = v3543
	var v3545 int32
	_ = v3545
	var v3546 int32
	_ = v3546
	var v3547 int32
	_ = v3547
	var v3549 int32
	_ = v3549
	var v3552 int32
	_ = v3552
	var v3554 int32
	_ = v3554
	var v3556 int32
	_ = v3556
	var v3558 int32
	_ = v3558
	var v3560 int32
	_ = v3560
	var v3562 int32
	_ = v3562
	var v3564 int32
	_ = v3564
	var v3569 int32
	_ = v3569
	var v3595 int32
	_ = v3595
	var v3597 int32
	_ = v3597
	var v3603 int32
	_ = v3603
	var v3604 int32
	_ = v3604
	var v3605 int32
	_ = v3605
	var v3606 int32
	_ = v3606
	var v3607 int32
	_ = v3607
	var v3608 int32
	_ = v3608
	var v3609 int32
	_ = v3609
	var v3610 int32
	_ = v3610
	var v3611 int32
	_ = v3611
	var v3612 int32
	_ = v3612
	var v3613 int32
	_ = v3613
	var v3614 int32
	_ = v3614
	var v3615 int32
	_ = v3615
	var v3616 int32
	_ = v3616
	var v3714 int32
	_ = v3714
	var v3794 int32
	_ = v3794
	var v3799 int32
	_ = v3799
	var v3800 int32
	_ = v3800
	var v3801 int32
	_ = v3801
	var v3802 int32
	_ = v3802
	var v3803 int32
	_ = v3803
	var v3817 int32
	_ = v3817
	var v3818 int32
	_ = v3818
	var v3821 int32
	_ = v3821
	var v3822 int32
	_ = v3822
	var v3826 int32
	_ = v3826
	var v3827 int32
	_ = v3827
	var v3829 int32
	_ = v3829
	var v3833 int32
	_ = v3833
	var v3841 int32
	_ = v3841
	var v3842 int32
	_ = v3842
	var v3845 int32
	_ = v3845
	var v3851 int32
	_ = v3851
	var v3857 int32
	_ = v3857
	var v3858 int32
	_ = v3858
	var v3866 int32
	_ = v3866
	var v3867 int32
	_ = v3867
	var v3869 int32
	_ = v3869
	var v3879 int32
	_ = v3879
	var v3891 int32
	_ = v3891
	var v3895 int32
	_ = v3895
	var v3896 int32
	_ = v3896
	var v3899 int32
	_ = v3899
	var v3903 int32
	_ = v3903
	var v3904 int32
	_ = v3904
	var v3905 int32
	_ = v3905
	var v3958 int64
	_ = v3958
	var v3970 int32
	_ = v3970
	var v3971 int64
	_ = v3971
	var v3997 int32
	_ = v3997
	var v3999 int32
	_ = v3999
	var v4064 int64
	_ = v4064
	var v4065 int64
	_ = v4065
	var v4075 int32
	_ = v4075
	var v4077 int32
	_ = v4077
	var v4080 int64
	_ = v4080
	var v4083 int64
	_ = v4083
	var v4086 int64
	_ = v4086
	var v4087 int64
	_ = v4087
	var v4094 int64
	_ = v4094
	var v4099 int64
	_ = v4099
	var v4101 int32
	_ = v4101
	var v4106 int32
	_ = v4106
	var v4107 int32
	_ = v4107
	var v4108 int64
	_ = v4108
	var v4109 int64
	_ = v4109
	var v4135 int32
	_ = v4135
	var v4137 int32
	_ = v4137
	var v4203 int64
	_ = v4203
	var v4204 int64
	_ = v4204
	var v4213 int32
	_ = v4213
	var v4215 int32
	_ = v4215
	var v4218 int64
	_ = v4218
	var v4221 int64
	_ = v4221
	var v4224 int64
	_ = v4224
	var v4225 int64
	_ = v4225
	var v4232 int64
	_ = v4232
	var v4237 int64
	_ = v4237
	var v4239 int32
	_ = v4239
	var v4244 int32
	_ = v4244
	var v4245 int32
	_ = v4245
	var v4246 int64
	_ = v4246
	var v4247 int64
	_ = v4247
	var v4273 int32
	_ = v4273
	var v4275 int32
	_ = v4275
	var v4341 int64
	_ = v4341
	var v4344 int64
	_ = v4344
	var v4351 int32
	_ = v4351
	var v4353 int32
	_ = v4353
	var v4356 int64
	_ = v4356
	var v4359 int64
	_ = v4359
	var v4362 int64
	_ = v4362
	var v4363 int64
	_ = v4363
	var v4370 int64
	_ = v4370
	var v4375 int64
	_ = v4375
	var v4377 int32
	_ = v4377
	var v4382 int32
	_ = v4382
	var v4383 int32
	_ = v4383
	var v4384 int64
	_ = v4384
	var v4385 int64
	_ = v4385
	var v4411 int32
	_ = v4411
	var v4413 int32
	_ = v4413
	var v4479 int64
	_ = v4479
	var v4484 int64
	_ = v4484
	var v4489 int32
	_ = v4489
	var v4491 int32
	_ = v4491
	var v4494 int64
	_ = v4494
	var v4497 int64
	_ = v4497
	var v4500 int64
	_ = v4500
	var v4501 int64
	_ = v4501
	var v4508 int64
	_ = v4508
	var v4513 int64
	_ = v4513
	var v4515 int32
	_ = v4515
	var v4522 int64
	_ = v4522
	var v4525 int64
	_ = v4525
	var v4531 int64
	_ = v4531
	var v4534 int64
	_ = v4534
	var v4541 int64
	_ = v4541
	var v4544 int64
	_ = v4544
	var v4551 int64
	_ = v4551
	var v4554 int64
	_ = v4554
	var v4556 int32
	_ = v4556
	var v4557 int32
	_ = v4557
	var v4558 int64
	_ = v4558
	var v4559 int64
	_ = v4559
	var v4563 int64
	_ = v4563
	var v4567 int64
	_ = v4567
	var v4568 int32
	_ = v4568
	var v4569 int64
	_ = v4569
	var v4570 int32
	_ = v4570
	var v4571 int32
	_ = v4571
	var v4572 int32
	_ = v4572
	var v4583 int32
	_ = v4583
	var v4587 int32
	_ = v4587
	var v4588 int32
	_ = v4588
	var v4591 int32
	_ = v4591
	var v4607 int32
	_ = v4607
	var v4608 int32
	_ = v4608
	var v4614 int32
	_ = v4614
	var v4621 int32
	_ = v4621
	var v4622 int32
	_ = v4622
	var v4625 int32
	_ = v4625
	var v4630 int32
	_ = v4630
	var v4635 int32
	_ = v4635
	var v4645 int32
	_ = v4645
	var v4657 int32
	_ = v4657
	var v4679 int32
	_ = v4679
	var v4682 int32
	_ = v4682
	var v4683 int32
	_ = v4683
	var v4691 int32
	_ = v4691
	var v4692 int32
	_ = v4692
	var v4727 int32
	_ = v4727
	var v4728 int32
	_ = v4728
	var v4733 int64
	_ = v4733
	var v4736 int32
	_ = v4736
	var v4737 int32
	_ = v4737
	var v4749 int32
	_ = v4749
	var v4770 int32
	_ = v4770
	var v4771 int32
	_ = v4771
	var v4774 int32
	_ = v4774
	var v4776 int32
	_ = v4776
	var v4787 int32
	_ = v4787
	var v4809 int32
	_ = v4809
	var v4812 int32
	_ = v4812
	var v4813 int32
	_ = v4813
	var v4821 int32
	_ = v4821
	var v4822 int32
	_ = v4822
	var v4857 int32
	_ = v4857
	var v4858 int32
	_ = v4858
	var v4863 int64
	_ = v4863
	var v4866 int32
	_ = v4866
	var v4867 int32
	_ = v4867
	var v4879 int32
	_ = v4879
	var v4897 int32
	_ = v4897
	var v4898 int32
	_ = v4898
	var v4900 int32
	_ = v4900
	var v4908 int32
	_ = v4908
	var v4911 int32
	_ = v4911
	var v4912 int32
	_ = v4912
	var v4918 int32
	_ = v4918
	var v4919 int32
	_ = v4919
	var v4920 int32
	_ = v4920
	var v4921 int32
	_ = v4921
	var v4923 int32
	_ = v4923
	var v4925 int32
	_ = v4925
	var v4927 int32
	_ = v4927
	var v4929 int32
	_ = v4929
	var v5030 int32
	_ = v5030
	var v5035 int32
	_ = v5035
	var v5042 int32
	_ = v5042
	var v5066 int32
	_ = v5066
	var v5067 int32
	_ = v5067
	var v5078 int32
	_ = v5078
	var v5079 int32
	_ = v5079
	var v5083 int32
	_ = v5083
	var v5134 int64
	_ = v5134
	var v5144 int32
	_ = v5144
	var v5145 int32
	_ = v5145
	var v5146 int64
	_ = v5146
	var v5147 int32
	_ = v5147
	var v5150 int32
	_ = v5150
	var v5151 int32
	_ = v5151
	var v5152 int64
	_ = v5152
	var v5157 int32
	_ = v5157
	var v5158 int64
	_ = v5158
	var v5163 int32
	_ = v5163
	var v5164 int64
	_ = v5164
	var v5169 int32
	_ = v5169
	var v5170 int64
	_ = v5170
	var v5171 int64
	_ = v5171
	var v5175 int32
	_ = v5175
	var v5176 int32
	_ = v5176
	var v5177 int64
	_ = v5177
	var v5185 int32
	_ = v5185
	var v5188 int32
	_ = v5188
	var v5191 int32
	_ = v5191
	var v5192 int32
	_ = v5192
	var v5195 int32
	_ = v5195
	var v5196 int32
	_ = v5196
	var v5197 int32
	_ = v5197
	var v5198 int32
	_ = v5198
	var v5200 int32
	_ = v5200
	var v5201 int32
	_ = v5201
	var v5202 int32
	_ = v5202
	var v5203 int32
	_ = v5203
	var v5204 int32
	_ = v5204
	var v5206 int32
	_ = v5206
	var v5207 int32
	_ = v5207
	var v5208 int32
	_ = v5208
	var v5209 int32
	_ = v5209
	var v5225 int32
	_ = v5225
	var v5228 int32
	_ = v5228
	var v5233 int32
	_ = v5233
	var v5288 int32
	_ = v5288
	var v5304 int32
	_ = v5304
	var v5326 int32
	_ = v5326
	var v5402 int32
	_ = v5402
	var v5486 int32
	_ = v5486
	var v5494 int32
	_ = v5494
	var v5497 int32
	_ = v5497
	var v5498 int32
	_ = v5498
	var v5501 int32
	_ = v5501
	var v5502 int32
	_ = v5502
	var v5504 int32
	_ = v5504
	var v5505 int32
	_ = v5505
	var v5506 int32
	_ = v5506
	var v5507 int32
	_ = v5507
	var v5508 int32
	_ = v5508
	var v5510 int32
	_ = v5510
	var v5511 int32
	_ = v5511
	var v5512 int32
	_ = v5512
	var v5513 int32
	_ = v5513
	var v5514 int32
	_ = v5514
	var v5516 int32
	_ = v5516
	var v5517 int32
	_ = v5517
	var v5518 int32
	_ = v5518
	var v5519 int32
	_ = v5519
	var v5520 int32
	_ = v5520
	var v5521 int32
	_ = v5521
	var v5601 int32
	_ = v5601
	var v5604 int32
	_ = v5604
	var v5607 int32
	_ = v5607
	var v5612 int32
	_ = v5612
	var v5613 int32
	_ = v5613
	var v5615 int32
	_ = v5615
	var v5616 int32
	_ = v5616
	var v5618 int32
	_ = v5618
	var v5623 int32
	_ = v5623
	var v5628 int32
	_ = v5628
	var v5637 int32
	_ = v5637
	var __phi5637 int32
	_ = __phi5637
	var v5639 int32
	_ = v5639
	var __phi5639 int32
	_ = __phi5639
	var v5653 int32
	_ = v5653
	var __phi5653 int32
	_ = __phi5653
	var v5654 int32
	_ = v5654
	var __phi5654 int32
	_ = __phi5654
	var v5728 int32
	_ = v5728
	var v5729 int32
	_ = v5729
	var v5734 int32
	_ = v5734
	var v5736 int32
	_ = v5736
	var v5737 int32
	_ = v5737
	var v5739 int32
	_ = v5739
	var v5745 int32
	_ = v5745
	var v5751 int32
	_ = v5751
	var v5752 int32
	_ = v5752
	var v5762 int32
	_ = v5762
	var __phi5762 int32
	_ = __phi5762
	var v5764 int32
	_ = v5764
	var __phi5764 int32
	_ = __phi5764
	var v5778 int32
	_ = v5778
	var __phi5778 int32
	_ = __phi5778
	var v5788 int32
	_ = v5788
	var __phi5788 int32
	_ = __phi5788
	var v5789 int32
	_ = v5789
	var __phi5789 int32
	_ = __phi5789
	var v5854 int32
	_ = v5854
	var v5856 int32
	_ = v5856
	var v5858 int32
	_ = v5858
	var v5863 int32
	_ = v5863
	var v5872 int32
	_ = v5872
	var v5890 int32
	_ = v5890
	var v5891 base.V128
	_ = v5891
	var v5895 base.V128
	_ = v5895
	var v5898 base.V128
	_ = v5898
	var v5908 base.V128
	_ = v5908
	var v5913 int32
	_ = v5913
	var v5917 int32
	_ = v5917
	var v5919 int32
	_ = v5919
	var v5925 int32
	_ = v5925
	var __phi5925 int32
	_ = __phi5925
	var v5926 int32
	_ = v5926
	var __phi5926 int32
	_ = __phi5926
	var v5933 int32
	_ = v5933
	var __phi5933 int32
	_ = __phi5933
	var v5934 int32
	_ = v5934
	var __phi5934 int32
	_ = __phi5934
	var v5935 int32
	_ = v5935
	var __phi5935 int32
	_ = __phi5935
	var v5946 int32
	_ = v5946
	var v5947 int32
	_ = v5947
	var v5949 int32
	_ = v5949
	var v5951 int32
	_ = v5951
	var v5954 int32
	_ = v5954
	var v5955 int32
	_ = v5955
	var v5956 int32
	_ = v5956
	var v5958 int32
	_ = v5958
	var v5961 int32
	_ = v5961
	var v5966 int32
	_ = v5966
	var v5970 int32
	_ = v5970
	var v5985 int32
	_ = v5985
	var v5999 int32
	_ = v5999
	var v6000 int32
	_ = v6000
	var v6001 int32
	_ = v6001
	var v6003 int32
	_ = v6003
	var v6004 int32
	_ = v6004
	var v6007 int32
	_ = v6007
	var v6008 int32
	_ = v6008
	var v6009 int32
	_ = v6009
	var v6011 int32
	_ = v6011
	var v6012 int32
	_ = v6012
	var v6014 int32
	_ = v6014
	var v6015 int32
	_ = v6015
	var v6020 int32
	_ = v6020
	var v6022 int32
	_ = v6022
	var v6024 int32
	_ = v6024
	var v6026 int32
	_ = v6026
	var v6027 int32
	_ = v6027
	var v6029 int32
	_ = v6029
	var v6030 int32
	_ = v6030
	var v6035 int32
	_ = v6035
	var v6036 int32
	_ = v6036
	var v6037 int32
	_ = v6037
	var v6039 int32
	_ = v6039
	var v6041 int32
	_ = v6041
	var v6044 int32
	_ = v6044
	var v6046 int32
	_ = v6046
	var v6048 int32
	_ = v6048
	var v6050 int32
	_ = v6050
	var v6052 int32
	_ = v6052
	var v6055 int32
	_ = v6055
	var v6056 int32
	_ = v6056
	var v6057 int32
	_ = v6057
	var v6059 int32
	_ = v6059
	var v6064 int32
	_ = v6064
	var v6066 int32
	_ = v6066
	var v6068 int32
	_ = v6068
	var v6070 int32
	_ = v6070
	var v6073 int32
	_ = v6073
	var v6075 int32
	_ = v6075
	var v6076 int32
	_ = v6076
	var v6077 int32
	_ = v6077
	var v6079 int32
	_ = v6079
	var v6082 int32
	_ = v6082
	var v6084 int32
	_ = v6084
	var v6086 int32
	_ = v6086
	var v6088 int32
	_ = v6088
	var v6090 int32
	_ = v6090
	var v6092 int32
	_ = v6092
	var v6095 int32
	_ = v6095
	var v6096 int32
	_ = v6096
	var v6097 int32
	_ = v6097
	var v6099 int32
	_ = v6099
	var v6104 int32
	_ = v6104
	var v6106 int32
	_ = v6106
	var v6108 int32
	_ = v6108
	var v6110 int32
	_ = v6110
	var v6113 int32
	_ = v6113
	var v6115 int32
	_ = v6115
	var v6116 int32
	_ = v6116
	var v6117 int32
	_ = v6117
	var v6119 int32
	_ = v6119
	var v6122 int32
	_ = v6122
	var v6124 int32
	_ = v6124
	var v6126 int32
	_ = v6126
	var v6128 int32
	_ = v6128
	var v6130 int32
	_ = v6130
	var v6133 int32
	_ = v6133
	var v6134 int32
	_ = v6134
	var v6135 int32
	_ = v6135
	var v6137 int32
	_ = v6137
	var v6142 int32
	_ = v6142
	var v6144 int32
	_ = v6144
	var v6146 int32
	_ = v6146
	var v6148 int32
	_ = v6148
	var v6151 int32
	_ = v6151
	var v6153 int32
	_ = v6153
	var v6154 int32
	_ = v6154
	var v6155 int32
	_ = v6155
	var v6157 int32
	_ = v6157
	var v6160 int32
	_ = v6160
	var v6162 int32
	_ = v6162
	var v6164 int32
	_ = v6164
	var v6166 int32
	_ = v6166
	var v6168 int32
	_ = v6168
	var v6170 int32
	_ = v6170
	var v6172 int32
	_ = v6172
	var v6177 int32
	_ = v6177
	var v6198 int32
	_ = v6198
	var v6199 int32
	_ = v6199
	var v6214 int32
	_ = v6214
	var v6307 int32
	_ = v6307
	var v6308 int32
	_ = v6308
	var v6309 int32
	_ = v6309
	var v6310 int32
	_ = v6310
	var v6319 int32
	_ = v6319
	var v6326 int32
	_ = v6326
	var __phi6326 int32
	_ = __phi6326
	var v6327 int32
	_ = v6327
	var __phi6327 int32
	_ = __phi6327
	var v6329 int32
	_ = v6329
	var __phi6329 int32
	_ = __phi6329
	var v6343 int32
	_ = v6343
	var __phi6343 int32
	_ = __phi6343
	var v6344 int32
	_ = v6344
	var __phi6344 int32
	_ = __phi6344
	var v6353 int32
	_ = v6353
	var __phi6353 int32
	_ = __phi6353
	var v6419 int32
	_ = v6419
	var v6422 int32
	_ = v6422
	var v6424 int32
	_ = v6424
	var v6428 int32
	_ = v6428
	var v6446 int32
	_ = v6446
	var v6447 base.V128
	_ = v6447
	var v6451 base.V128
	_ = v6451
	var v6454 base.V128
	_ = v6454
	var v6464 base.V128
	_ = v6464
	var v6469 int32
	_ = v6469
	var v6473 int32
	_ = v6473
	var v6475 int32
	_ = v6475
	var v6481 int32
	_ = v6481
	var __phi6481 int32
	_ = __phi6481
	var v6482 int32
	_ = v6482
	var __phi6482 int32
	_ = __phi6482
	var v6489 int32
	_ = v6489
	var __phi6489 int32
	_ = __phi6489
	var v6490 int32
	_ = v6490
	var __phi6490 int32
	_ = __phi6490
	var v6491 int32
	_ = v6491
	var __phi6491 int32
	_ = __phi6491
	var v6502 int32
	_ = v6502
	var v6503 int32
	_ = v6503
	var v6505 int32
	_ = v6505
	var v6507 int32
	_ = v6507
	var v6510 int32
	_ = v6510
	var v6511 int32
	_ = v6511
	var v6512 int32
	_ = v6512
	var v6514 int32
	_ = v6514
	var v6517 int32
	_ = v6517
	var v6522 int32
	_ = v6522
	var v6526 int32
	_ = v6526
	var v6541 int32
	_ = v6541
	var v6555 int32
	_ = v6555
	var v6556 int32
	_ = v6556
	var v6557 int32
	_ = v6557
	var v6559 int32
	_ = v6559
	var v6560 int32
	_ = v6560
	var v6563 int32
	_ = v6563
	var v6564 int32
	_ = v6564
	var v6565 int32
	_ = v6565
	var v6567 int32
	_ = v6567
	var v6568 int32
	_ = v6568
	var v6570 int32
	_ = v6570
	var v6571 int32
	_ = v6571
	var v6576 int32
	_ = v6576
	var v6578 int32
	_ = v6578
	var v6580 int32
	_ = v6580
	var v6582 int32
	_ = v6582
	var v6583 int32
	_ = v6583
	var v6585 int32
	_ = v6585
	var v6586 int32
	_ = v6586
	var v6591 int32
	_ = v6591
	var v6592 int32
	_ = v6592
	var v6593 int32
	_ = v6593
	var v6595 int32
	_ = v6595
	var v6597 int32
	_ = v6597
	var v6600 int32
	_ = v6600
	var v6602 int32
	_ = v6602
	var v6604 int32
	_ = v6604
	var v6606 int32
	_ = v6606
	var v6608 int32
	_ = v6608
	var v6611 int32
	_ = v6611
	var v6612 int32
	_ = v6612
	var v6613 int32
	_ = v6613
	var v6615 int32
	_ = v6615
	var v6620 int32
	_ = v6620
	var v6622 int32
	_ = v6622
	var v6624 int32
	_ = v6624
	var v6626 int32
	_ = v6626
	var v6629 int32
	_ = v6629
	var v6631 int32
	_ = v6631
	var v6632 int32
	_ = v6632
	var v6633 int32
	_ = v6633
	var v6635 int32
	_ = v6635
	var v6638 int32
	_ = v6638
	var v6640 int32
	_ = v6640
	var v6642 int32
	_ = v6642
	var v6644 int32
	_ = v6644
	var v6646 int32
	_ = v6646
	var v6648 int32
	_ = v6648
	var v6651 int32
	_ = v6651
	var v6652 int32
	_ = v6652
	var v6653 int32
	_ = v6653
	var v6655 int32
	_ = v6655
	var v6660 int32
	_ = v6660
	var v6662 int32
	_ = v6662
	var v6664 int32
	_ = v6664
	var v6666 int32
	_ = v6666
	var v6669 int32
	_ = v6669
	var v6671 int32
	_ = v6671
	var v6672 int32
	_ = v6672
	var v6673 int32
	_ = v6673
	var v6675 int32
	_ = v6675
	var v6678 int32
	_ = v6678
	var v6680 int32
	_ = v6680
	var v6682 int32
	_ = v6682
	var v6684 int32
	_ = v6684
	var v6686 int32
	_ = v6686
	var v6689 int32
	_ = v6689
	var v6690 int32
	_ = v6690
	var v6691 int32
	_ = v6691
	var v6693 int32
	_ = v6693
	var v6698 int32
	_ = v6698
	var v6700 int32
	_ = v6700
	var v6702 int32
	_ = v6702
	var v6704 int32
	_ = v6704
	var v6707 int32
	_ = v6707
	var v6709 int32
	_ = v6709
	var v6710 int32
	_ = v6710
	var v6711 int32
	_ = v6711
	var v6713 int32
	_ = v6713
	var v6716 int32
	_ = v6716
	var v6718 int32
	_ = v6718
	var v6720 int32
	_ = v6720
	var v6722 int32
	_ = v6722
	var v6724 int32
	_ = v6724
	var v6726 int32
	_ = v6726
	var v6728 int32
	_ = v6728
	var v6733 int32
	_ = v6733
	var v6756 int32
	_ = v6756
	var v6776 int32
	_ = v6776
	var v6791 int32
	_ = v6791
	var v6869 int32
	_ = v6869
	var v6871 int32
	_ = v6871
	var v6874 int32
	_ = v6874
	var v6879 int32
	_ = v6879
	var v6890 int32
	_ = v6890
	var v6982 int32
	_ = v6982
	var v6991 int32
	_ = v6991
	var v7090 int32
	_ = v7090
	var v7099 int32
	_ = v7099
	var v7106 int32
	_ = v7106
	v98 = m.G0
	v100 = v98 - int32(2144)
	m.G0 = v100
	v104 = base.I32_div_s(l8, int32(-20))
	v107 = int32(1) << (uint(v104+int32(5)) % 32)
	v108 = *(*int32)(unsafe.Add(mBase, uint32(l13)))
	if l4 == int32(0) {
		goto L3
	} else {
		goto L4
	}
L1:
	;
	m.G0 = v7106 + int32(2144)
	return v7099
L2:
	;
	if v5502 < int32(1) {
		goto L372
	} else {
		goto L373
	}
L3:
	;
	if l3 < l2 {
		v971 = int32(0)
		goto L16
	} else {
		goto L17
	}
L4:
	;
	v111 = int32(-1)
	v112 = l0 + v111
	v113 = int32(1)
	v114 = v113 << (uint(l3) % 32)
	v116 = int32(base.Ui32(v112+v114) >> (uint(l3) % 32))
	v121 = v116 * int32(base.Ui32(l1+v114+v111)>>(uint(l3)%32))
	if v121 < v113 {
		goto L5
	} else {
		goto L6
	}
L5:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l14))) = l3
	v5501 = l0
	v5502 = l1
	v5504 = l3
	v5505 = l4
	v5506 = l5
	v5507 = l6
	v5508 = l7
	v5510 = l9
	v5511 = l10
	v5512 = l11
	v5513 = l12
	v5514 = l13
	v5516 = v100
	v5517 = v107
	v5518 = v108
	v5519 = v112
	v5520 = v114
	v5521 = v116
	goto L2
L6:
	;
	if base.Ui32(v121) < base.Ui32(int32(4)) {
		v256 = int32(0)
		goto L7
	} else {
		goto L8
	}
L7:
	;
	v343 = l7 + v256<<(uint(int32(2))%32)
	v358 = v121 - v256
	goto L13
L8:
	;
	v128 = v121 & int32(2147483644)
	v137 = l7
	v152 = v128
	goto L9
L9:
	;
	v226 = base.Simd_g_const(&F_VP8LResidualImage__k0)
	base.Simd_g_v128_store(m, v137, int32(0), v226)
	v232 = v152 + int32(-4)
	if v232 != 0 {
		v137 = v137 + int32(16)
		v152 = v232
		goto L9
	} else {
		goto L11
	}
L10:
	;
	if v121 == v128 {
		goto L5
	} else {
		goto L12
	}
L11:
	;
	goto L10
L12:
	;
	v256 = v128
	goto L7
L13:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v343))) = int32(-16774400)
	v437 = v358 + int32(-1)
	if v437 != 0 {
		v343 = v343 + int32(4)
		v358 = v437
		goto L13
	} else {
		goto L15
	}
L14:
	;
	goto L5
L15:
	;
	goto L14
L16:
	;
	v1047 = base.I64_extend_i32_u(v971)
	v1048 = int32(4)
	if v1047 == int64(0) {
		goto L31
	} else {
		goto L32
	}
L17:
	;
	v538 = int32(-1)
	v539 = l1 + v538
	v541 = l0 + v538
	v544 = l3 + int32(1)
	v545 = v544 - l2
	if base.Ui32(v545) < base.Ui32(int32(4)) {
		v747 = l2
		v760 = int32(0)
		goto L18
	} else {
		goto L19
	}
L18:
	;
	v847 = v747
	v860 = v760
	v862 = v100 + v747<<(uint(int32(2))%32)
	goto L24
L19:
	;
	v552 = v545 & int32(-4)
	v583 = v100 + l2<<(uint(int32(2))%32)
	v587 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_splat(l2), base.Simd_g_const(&F_VP8LResidualImage__k1))
	v591 = base.Simd_g_const(&F_VP8LResidualImage__k2)
	v592 = v552
	goto L20
L20:
	;
	v657 = int32(1)
	v658 = int32(0)
	v659 = base.Simd_g_i32x4_extract_lane_l0(v587)
	v664 = base.Simd_g_i32x4_extract_lane_l1(v587)
	v669 = int32(2)
	v670 = base.Simd_g_i32x4_extract_lane_l2(v587)
	v675 = int32(3)
	v676 = base.Simd_g_i32x4_extract_lane_l3(v587)
	v679 = base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v657<<(uint(v659)%32)), v657<<(uint(v664)%32)), v657<<(uint(v670)%32)), v657<<(uint(v676)%32))
	v680 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_splat(v541), v679)
	v700 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_splat(v539), v679)
	v720 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(v680))>>(uint(v659)%32))), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l1(v680))>>(uint(v664)%32))), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l2(v680))>>(uint(v670)%32))), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l3(v680))>>(uint(v676)%32))), base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(v700))>>(uint(v659)%32))), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l1(v700))>>(uint(v664)%32))), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l2(v700))>>(uint(v670)%32))), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l3(v700))>>(uint(v676)%32))))
	base.Simd_g_v128_store(m, v583, v658, v720)
	v723 = base.Simd_g_i32x4_add(v720, v591)
	v727 = base.Simd_g_i32x4_add(v587, base.Simd_g_const(&F_VP8LResidualImage__k3))
	v729 = v592 + int32(-4)
	if v729 != 0 {
		v583 = v583 + int32(16)
		v587 = v727
		v591 = v723
		v592 = v729
		goto L20
	} else {
		goto L22
	}
L21:
	;
	v732 = base.Simd_g_i32x4_add(v723, base.Simd_g_i8x16_swizzle_c(v723, base.Simd_g_const(&F_VP8LResidualImage__k4)))
	v737 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v732, base.Simd_g_i8x16_swizzle_c(v732, base.Simd_g_const(&F_VP8LResidualImage__k5))))
	if v545 == v552 {
		v971 = v737
		goto L16
	} else {
		goto L23
	}
L22:
	;
	goto L21
L23:
	;
	v747 = v552 + l2
	v760 = v737
	goto L18
L24:
	;
	v936 = int32(1)
	v937 = v936 << (uint(v847) % 32)
	v942 = int32(base.Ui32(v541+v937)>>(uint(v847)%32)) * int32(base.Ui32(v539+v937)>>(uint(v847)%32))
	*(*int32)(unsafe.Add(mBase, uint32(v862))) = v942
	v944 = v942 + v860
	v948 = v847 + v936
	if v544 != v948 {
		v847 = v948
		v860 = v944
		v862 = v862 + int32(4)
		goto L24
	} else {
		goto L26
	}
L25:
	;
	v971 = v944
	goto L16
L26:
	;
	goto L25
L27:
	;
	v5486 = int32(1) << (uint(v5288) % 32)
	v5494 = F_memcpy(m, v5198, v5225, int32(base.Ui32(v5486+v5233)>>(uint(v5288)%32))*int32(base.Ui32(v5486+v5209)>>(uint(v5288)%32))<<(uint(int32(2))%32))
	mBase = m.M
	F_free(m, v5228)
	mBase = m.M
	goto L371
L28:
	;
	v7099 = int32(0)
	v7106 = v5402
	goto L1
L29:
	;
	if v1069 == int32(0) {
		v5402 = v100
		goto L28
	} else {
		goto L35
	}
L30:
	;
	goto L29
L31:
	;
	v1067 = F_malloc(m, base.I32_wrap_i64(v1047)*v1048)
	mBase = m.M
	v1069 = v1067
	goto L30
L32:
	;
	v1055 = base.I64_div_u_s(int64(2147418112), v1047)
	v1056 = int32(0)
	v1057 = base.I64_extend_i32_u(v1048)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1057*v1047) {
		v1069 = v1056
		goto L30
	} else {
		goto L33
	}
L33:
	;
	if base.Ui64(v1055) < base.Ui64(v1057) {
		v1069 = v1056
		goto L30
	} else {
		goto L34
	}
L34:
	;
	goto L31
L35:
	;
	v1077 = v100 + int32(48) + l2<<(uint(int32(2))%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1077))) = v1069
	if l3 <= l2 {
		goto L36
	} else {
		goto L37
	}
L36:
	;
	v1544 = l3 - l2
	v1546 = v1544 + int32(1)
	v1549 = base.I64_extend_i32_s(v1546 * int32(_a_F_VP8LResidualImage_0))
	v1550 = int32(4)
	if v1549 == int64(0) {
		goto L49
	} else {
		goto L50
	}
L37:
	;
	v1082 = (l3 - l2) & int32(3)
	if v1082 == int32(0) {
		v1211 = v1069
		v1225 = l2
		goto L38
	} else {
		goto L39
	}
L38:
	;
	if base.Ui32(int32(-4)) < base.Ui32(l2-l3) {
		goto L36
	} else {
		goto L43
	}
L39:
	;
	v1086 = l2 << (uint(int32(2)) % 32)
	v1101 = v1069
	v1114 = v100 + v1086
	v1115 = l2
	v1116 = v1086 + (v100 + int32(48)) + int32(4)
	v1126 = v1082
	goto L40
L40:
	;
	v1190 = *(*int32)(unsafe.Add(mBase, uint32(v1114)))
	v1193 = v1101 + v1190<<(uint(int32(2))%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1116))) = v1193
	v1195 = int32(4)
	v1200 = v1115 + int32(1)
	v1202 = v1126 + int32(-1)
	if v1202 != 0 {
		v1101 = v1193
		v1114 = v1114 + v1195
		v1115 = v1200
		v1116 = v1116 + v1195
		v1126 = v1202
		goto L40
	} else {
		goto L42
	}
L41:
	;
	v1211 = v1193
	v1225 = v1200
	goto L38
L42:
	;
	goto L41
L43:
	;
	v1304 = v1225 << (uint(int32(2)) % 32)
	v1316 = v1211
	v1330 = v100 + int32(48)
	v1340 = l3 - v1225
	v1341 = v100
	goto L44
L44:
	;
	v1405 = v1330 + v1304
	v1406 = int32(4)
	v1408 = v1341 + v1304
	v1409 = *(*int32)(unsafe.Add(mBase, uint32(v1408)))
	v1410 = int32(2)
	v1412 = v1316 + v1409<<(uint(v1410)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1405+v1406))) = v1412
	v1414 = int32(8)
	v1418 = *(*int32)(unsafe.Add(mBase, uint32(v1408+v1406)))
	v1421 = v1412 + v1418<<(uint(v1410)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1405+v1414))) = v1421
	v1423 = int32(12)
	v1427 = *(*int32)(unsafe.Add(mBase, uint32(v1408+v1414)))
	v1430 = v1421 + v1427<<(uint(v1410)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1405+v1423))) = v1430
	v1432 = int32(16)
	v1436 = *(*int32)(unsafe.Add(mBase, uint32(v1408+v1423)))
	v1439 = v1430 + v1436<<(uint(v1410)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1405+v1432))) = v1439
	v1446 = v1340 + int32(-4)
	if v1446 != 0 {
		v1316 = v1439
		v1330 = v1330 + v1432
		v1340 = v1446
		v1341 = v1341 + v1432
		goto L44
	} else {
		goto L46
	}
L45:
	;
	goto L36
L46:
	;
	goto L45
L47:
	;
	v1572 = *(*int32)(unsafe.Add(mBase, uint32(l13)))
	v1573 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l14))) = v1573
	if v1570 == v1573 {
		v5304 = v100
		v5326 = v1069
		goto L53
	} else {
		goto L54
	}
L48:
	;
	goto L47
L49:
	;
	v1568 = F_calloc(m, base.I32_wrap_i64(v1549), v1550)
	mBase = m.M
	v1570 = v1568
	goto L48
L50:
	;
	v1557 = base.I64_div_u_s(int64(2147418112), v1549)
	v1558 = int32(0)
	v1559 = base.I64_extend_i32_u(v1550)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1559*v1549) {
		v1570 = v1558
		goto L48
	} else {
		goto L51
	}
L51:
	;
	if base.Ui64(v1557) < base.Ui64(v1559) {
		v1570 = v1558
		goto L48
	} else {
		goto L52
	}
L52:
	;
	goto L49
L53:
	;
	F_free(m, v5326)
	mBase = m.M
	goto L370
L54:
	;
	v1577 = int32(-1)
	v1578 = l0 + v1577
	v1580 = int32(1) << (uint(l2) % 32)
	v1582 = l1 + v1577
	v1584 = int32(base.Ui32(v1580+v1582) >> (uint(l2) % 32))
	if v1584 == int32(0) {
		goto L56
	} else {
		goto L57
	}
L55:
	;
	v5288 = *(*int32)(unsafe.Add(mBase, uint32(l14)))
	if v5288 != 0 {
		goto L27
	} else {
		goto L369
	}
L56:
	;
	v5030 = int32(_a_F_VP8LResidualImage_1)
	v5035 = int32(_a_F_VP8LResidualImage_2)
	v5042 = int32(0)
	v5066 = l3*v5030 - l2*v5030
	v5067 = l3*v5035 - l2*v5035 + v5035
	v5078 = v5042
	v5079 = v5042
	v5083 = v1077
	v5134 = int64(9223372036854775807)
	goto L363
L57:
	;
	v1589 = v1570 + v1546*int32(_a_F_VP8LResidualImage_1)
	v1593 = int32(4)
	if l3 < v1593 {
		goto L58
	} else {
		goto L59
	}
L58:
	;
	v1596 = l3
	goto L60
L59:
	;
	v1596 = v1593
	goto L60
L60:
	;
	if l2 < v1596 {
		goto L61
	} else {
		goto L62
	}
L61:
	;
	v1598 = v1596
	goto L63
L62:
	;
	v1598 = l2
	goto L63
L63:
	;
	v1599 = v1598 - l2
	v1601 = v1599 + int32(1)
	v1602 = int32(2)
	if base.Ui32(v1602) < base.Ui32(v1601) {
		goto L64
	} else {
		goto L65
	}
L64:
	;
	v1605 = v1601
	goto L66
L65:
	;
	v1605 = v1602
	goto L66
L66:
	;
	v1607 = l0 << (uint(int32(2)) % 32)
	v1609 = int32(4)
	v1610 = l6 + v1607 + v1609
	v1613 = v1610 + v1607 + v1609
	v1622 = int32(-1)
	v1627 = int32(base.Ui32(v1580+v1578)>>(uint(l2)%32)) + v1622
	v1628 = int32(_a_F_VP8LResidualImage_1)
	v1631 = v1544*v1628 + v1628
	v1632 = int32(0)
	v1696 = v1632
	v1697 = v1632
	v1698 = v1632
	v1699 = v1632
	v1700 = v1632
	v1701 = v1632
	goto L67
L67:
	;
	v1735 = v1699 << (uint(l2) % 32)
	v1737 = base.B2i32(int32(0) < v1735)
	v1738 = v1735 - v1737
	v1739 = int32(2)
	v1740 = v1738 << (uint(v1739) % 32)
	v1741 = l5 + v1740
	v1742 = v1698 << (uint(l2) % 32)
	v1748 = v1741 + (v1742+int32(-1))*l0<<(uint(v1739)%32)
	v1749 = l0 - v1735
	v1750 = base.B2i32(v1580 < v1749)
	if v1580 < v1749 {
		goto L69
	} else {
		goto L70
	}
L68:
	;
	goto L56
L69:
	;
	v1751 = v1580
	goto L71
L70:
	;
	v1751 = v1749
	goto L71
L71:
	;
	v1752 = v1751 + v1737
	v1756 = v1752<<(uint(int32(2))%32) + int32(4)
	v1757 = l1 - v1742
	if v1580 < v1757 {
		goto L74
	} else {
		goto L75
	}
L72:
	;
	v3714 = int32(2147483646)
	v3794 = int32(0)
	v3799 = v1570 + int32(4092)
	v3800 = v1570 + int32(3068)
	v3801 = v1570 + int32(2044)
	v3802 = v1570
	v3803 = v1570 + int32(1020)
	goto L237
L73:
	;
	if v1742 < int32(1) {
		goto L72
	} else {
		goto L234
	}
L74:
	;
	v1759 = v1580
	goto L76
L75:
	;
	v1759 = v1757
	goto L76
L76:
	;
	if v1759 < int32(1) {
		goto L73
	} else {
		goto L77
	}
L77:
	;
	v1762 = v1613 + v1738
	v1763 = v1751 + v1735
	v1764 = v1752 + v1750
	if v1751 < int32(1) {
		goto L78
	} else {
		goto L79
	}
L78:
	;
	if int32(1) < v107 {
		goto L158
	} else {
		goto L159
	}
L79:
	;
	v1776 = l6
	v1792 = int32(0)
	v1843 = v1610
	goto L80
L80:
	;
	if v1742 < int32(1) {
		goto L82
	} else {
		goto L83
	}
L82:
	;
	v1871 = v1570 + v1792<<(uint(int32(12))%32)
	__phi1881 = v1776
	__phi1901 = int32(0)
	__phi1948 = v1843
	v1881 = __phi1881
	v1901 = __phi1901
	v1948 = __phi1948
	goto L84
L83:
	;
	v1868 = F_memcpy(m, v1843+v1740, v1748, v1756)
	mBase = m.M
	goto L82
L84:
	;
	v1971 = v1901 + v1742
	v1973 = int32(2)
	v1975 = v1741 + v1971*l0<<(uint(v1973)%32)
	v1977 = v1971 + int32(1)
	v1982 = F_memcpy(m, v1881+v1740, v1975, (v1752+base.B2i32(v1977 < l1))<<(uint(v1973)%32))
	mBase = m.M
	if v107 < v1973 {
		goto L86
	} else {
		goto L87
	}
L85:
	;
	v2817 = v1792 + int32(1)
	if v2817 != int32(14) {
		v1776 = v1948
		v1792 = v2817
		v1843 = v1881
		goto L80
	} else {
		goto L157
	}
L86:
	;
	v2314 = v100 + int32(96)
	F_GetResidual(m, l0, l1, v1948, v1881, v1613, v1792, v1735, v1763, v1971, v107, l9, l10, v2314)
	mBase = m.M
	v2339 = v1751
	v2341 = v2314
	goto L145
L87:
	;
	if v1971 < int32(1) {
		goto L86
	} else {
		goto L88
	}
L88:
	;
	if l1 <= v1977 {
		goto L86
	} else {
		goto L89
	}
L89:
	;
	if v1764 < int32(3) {
		goto L91
	} else {
		goto L92
	}
L90:
	;
	goto L86
L91:
	;
	goto L90
L92:
	;
	v2005 = int32(0)
	v2006 = base.Simd_g_v128_load64_zero(m, v1975, v2005)
	if l10 == v2005 {
		v2023 = v2006
		goto L93
	} else {
		goto L94
	}
L93:
	;
	v2028 = int32(1)
	v2032 = int32(0)
	v2034 = l0 << (uint(int32(2)) % 32)
	__phi2040 = base.Simd_g_i32x4_extract_lane_l1(v2023)
	__phi2041 = base.Simd_g_i32x4_extract_lane_l0(v2023)
	__phi2048 = v1975 + int32(4)
	__phi2049 = v1762 + v2028
	__phi2050 = v1764 + int32(-2)
	v2040 = __phi2040
	v2041 = __phi2041
	v2048 = __phi2048
	v2049 = __phi2049
	v2050 = __phi2050
	goto L95
L94:
	;
	v2010 = base.Simd_g_i32x4_shr_u(v2006, int32(8))
	v2013 = base.Simd_g_const(&F_VP8LResidualImage__k6)
	v2023 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v2010, base.Simd_g_const(&F_VP8LResidualImage__k7)), base.Simd_g_v128_and(v2006, v2013)), base.Simd_g_i32x4_shl(v2010, int32(16))), v2013), base.Simd_g_v128_and(v2006, base.Simd_g_const(&F_VP8LResidualImage__k8)))
	goto L93
L95:
	;
	v2061 = v2048 + int32(4)
	v2062 = *(*int32)(unsafe.Add(mBase, uint32(v2061)))
	v2064 = *(*int32)(unsafe.Add(mBase, uint32(v2048+v2034)))
	v2066 = *(*int32)(unsafe.Add(mBase, uint32(v2048+(v2032-v2034))))
	if l10 == int32(0) {
		v2114 = v2062
		v2115 = v2064
		v2116 = v2066
		goto L97
	} else {
		goto L98
	}
L96:
	;
	goto L91
L97:
	;
	v2118 = int32(24)
	v2119 = int32(base.Ui32(v2040) >> (uint(v2118) % 32))
	v2122 = v2119 - int32(base.Ui32(v2116)>>(uint(v2118)%32))
	v2123 = int32(31)
	v2124 = v2122 >> (uint(v2123) % 32)
	v2126 = v2122 ^ v2124 - v2124
	v2127 = int32(16)
	v2129 = int32(255)
	v2130 = int32(base.Ui32(v2040)>>(uint(v2127)%32)) & v2129
	v2135 = v2130 - int32(base.Ui32(v2116)>>(uint(v2127)%32))&v2129
	v2137 = v2135 >> (uint(v2123) % 32)
	v2139 = v2135 ^ v2137 - v2137
	if base.Ui32(v2139) < base.Ui32(v2126) {
		goto L99
	} else {
		goto L100
	}
L98:
	;
	v2069 = int32(8)
	v2070 = int32(base.Ui32(v2062) >> (uint(v2069) % 32))
	v2071 = int32(255)
	v2073 = int32(16711935)
	v2076 = int32(16)
	v2081 = int32(-16711936)
	v2085 = int32(base.Ui32(v2064) >> (uint(v2069) % 32))
	v2100 = int32(base.Ui32(v2066) >> (uint(v2069) % 32))
	v2114 = (v2070&v2071+v2062&v2073+v2070<<(uint(v2076)%32))&v2073 | v2062&v2081
	v2115 = (v2085&v2071+v2064&v2073+v2085<<(uint(v2076)%32))&v2073 | v2064&v2081
	v2116 = (v2100&v2071+v2066&v2073+v2100<<(uint(v2076)%32))&v2073 | v2066&v2081
	goto L97
L99:
	;
	v2141 = v2126
	goto L101
L100:
	;
	v2141 = v2139
	goto L101
L101:
	;
	v2142 = int32(8)
	v2144 = int32(255)
	v2145 = int32(base.Ui32(v2040)>>(uint(v2142)%32)) & v2144
	v2150 = v2145 - int32(base.Ui32(v2116)>>(uint(v2142)%32))&v2144
	v2151 = int32(31)
	v2152 = v2150 >> (uint(v2151) % 32)
	v2154 = v2150 ^ v2152 - v2152
	v2156 = v2040 & v2144
	v2159 = v2156 - v2116&v2144
	v2161 = v2159 >> (uint(v2151) % 32)
	v2163 = v2159 ^ v2161 - v2161
	if base.Ui32(v2163) < base.Ui32(v2154) {
		goto L102
	} else {
		goto L103
	}
L102:
	;
	v2165 = v2154
	goto L104
L103:
	;
	v2165 = v2163
	goto L104
L104:
	;
	if base.Ui32(v2165) < base.Ui32(v2141) {
		goto L105
	} else {
		goto L106
	}
L105:
	;
	v2167 = v2141
	goto L107
L106:
	;
	v2167 = v2165
	goto L107
L107:
	;
	v2170 = v2119 - int32(base.Ui32(v2115)>>(uint(int32(24))%32))
	v2171 = int32(31)
	v2172 = v2170 >> (uint(v2171) % 32)
	v2174 = v2170 ^ v2172 - v2172
	v2179 = v2130 - int32(base.Ui32(v2115)>>(uint(int32(16))%32))&int32(255)
	v2181 = v2179 >> (uint(v2171) % 32)
	v2183 = v2179 ^ v2181 - v2181
	if base.Ui32(v2183) < base.Ui32(v2174) {
		goto L108
	} else {
		goto L109
	}
L108:
	;
	v2185 = v2174
	goto L110
L109:
	;
	v2185 = v2183
	goto L110
L110:
	;
	v2188 = int32(255)
	v2190 = v2145 - int32(base.Ui32(v2115)>>(uint(int32(8))%32))&v2188
	v2191 = int32(31)
	v2192 = v2190 >> (uint(v2191) % 32)
	v2194 = v2190 ^ v2192 - v2192
	v2197 = v2156 - v2115&v2188
	v2199 = v2197 >> (uint(v2191) % 32)
	v2201 = v2197 ^ v2199 - v2199
	if base.Ui32(v2201) < base.Ui32(v2194) {
		goto L111
	} else {
		goto L112
	}
L111:
	;
	v2203 = v2194
	goto L113
L112:
	;
	v2203 = v2201
	goto L113
L113:
	;
	if base.Ui32(v2203) < base.Ui32(v2185) {
		goto L114
	} else {
		goto L115
	}
L114:
	;
	v2205 = v2185
	goto L116
L115:
	;
	v2205 = v2203
	goto L116
L116:
	;
	if base.Ui32(v2205) < base.Ui32(v2167) {
		goto L117
	} else {
		goto L118
	}
L117:
	;
	v2207 = v2167
	goto L119
L118:
	;
	v2207 = v2205
	goto L119
L119:
	;
	v2210 = v2119 - int32(base.Ui32(v2041)>>(uint(int32(24))%32))
	v2211 = int32(31)
	v2212 = v2210 >> (uint(v2211) % 32)
	v2214 = v2210 ^ v2212 - v2212
	v2219 = v2130 - int32(base.Ui32(v2041)>>(uint(int32(16))%32))&int32(255)
	v2221 = v2219 >> (uint(v2211) % 32)
	v2223 = v2219 ^ v2221 - v2221
	if base.Ui32(v2223) < base.Ui32(v2214) {
		goto L120
	} else {
		goto L121
	}
L120:
	;
	v2225 = v2214
	goto L122
L121:
	;
	v2225 = v2223
	goto L122
L122:
	;
	v2228 = int32(255)
	v2230 = v2145 - int32(base.Ui32(v2041)>>(uint(int32(8))%32))&v2228
	v2231 = int32(31)
	v2232 = v2230 >> (uint(v2231) % 32)
	v2234 = v2230 ^ v2232 - v2232
	v2237 = v2156 - v2041&v2228
	v2239 = v2237 >> (uint(v2231) % 32)
	v2241 = v2237 ^ v2239 - v2239
	if base.Ui32(v2241) < base.Ui32(v2234) {
		goto L123
	} else {
		goto L124
	}
L123:
	;
	v2243 = v2234
	goto L125
L124:
	;
	v2243 = v2241
	goto L125
L125:
	;
	if base.Ui32(v2243) < base.Ui32(v2225) {
		goto L126
	} else {
		goto L127
	}
L126:
	;
	v2245 = v2225
	goto L128
L127:
	;
	v2245 = v2243
	goto L128
L128:
	;
	v2248 = v2119 - int32(base.Ui32(v2114)>>(uint(int32(24))%32))
	v2249 = int32(31)
	v2250 = v2248 >> (uint(v2249) % 32)
	v2252 = v2248 ^ v2250 - v2250
	v2257 = v2130 - int32(base.Ui32(v2114)>>(uint(int32(16))%32))&int32(255)
	v2259 = v2257 >> (uint(v2249) % 32)
	v2261 = v2257 ^ v2259 - v2259
	if base.Ui32(v2261) < base.Ui32(v2252) {
		goto L129
	} else {
		goto L130
	}
L129:
	;
	v2263 = v2252
	goto L131
L130:
	;
	v2263 = v2261
	goto L131
L131:
	;
	v2266 = int32(255)
	v2268 = v2145 - int32(base.Ui32(v2114)>>(uint(int32(8))%32))&v2266
	v2269 = int32(31)
	v2270 = v2268 >> (uint(v2269) % 32)
	v2272 = v2268 ^ v2270 - v2270
	v2275 = v2156 - v2114&v2266
	v2277 = v2275 >> (uint(v2269) % 32)
	v2279 = v2275 ^ v2277 - v2277
	if base.Ui32(v2279) < base.Ui32(v2272) {
		goto L132
	} else {
		goto L133
	}
L132:
	;
	v2281 = v2272
	goto L134
L133:
	;
	v2281 = v2279
	goto L134
L134:
	;
	if base.Ui32(v2281) < base.Ui32(v2263) {
		goto L135
	} else {
		goto L136
	}
L135:
	;
	v2283 = v2263
	goto L137
L136:
	;
	v2283 = v2281
	goto L137
L137:
	;
	if base.Ui32(v2283) < base.Ui32(v2245) {
		goto L138
	} else {
		goto L139
	}
L138:
	;
	v2285 = v2245
	goto L140
L139:
	;
	v2285 = v2283
	goto L140
L140:
	;
	if base.Ui32(v2285) < base.Ui32(v2207) {
		goto L141
	} else {
		goto L142
	}
L141:
	;
	v2287 = v2207
	goto L143
L142:
	;
	v2287 = v2285
	goto L143
L143:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v2049))) = uint8(v2287)
	v2292 = v2050 + int32(-1)
	if v2292 != 0 {
		__phi2040 = v2114
		__phi2041 = v2040
		__phi2048 = v2061
		__phi2049 = v2049 + int32(1)
		__phi2050 = v2292
		v2040 = __phi2040
		v2041 = __phi2041
		v2048 = __phi2048
		v2049 = __phi2049
		v2050 = __phi2050
		goto L95
	} else {
		goto L144
	}
L144:
	;
	goto L96
L145:
	;
	v2415 = *(*int32)(unsafe.Add(mBase, uint32(v2341)))
	v2418 = int32(1020)
	v2420 = v1871 + int32(base.Ui32(v2415)>>(uint(int32(22))%32))&v2418
	v2421 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	v2422 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2420))) = v2421 + v2422
	v2431 = v1871 + int32(base.Ui32(v2415)>>(uint(int32(14))%32))&v2418 + int32(1024)
	v2432 = *(*int32)(unsafe.Add(mBase, uint32(v2431)))
	*(*int32)(unsafe.Add(mBase, uint32(v2431))) = v2432 + v2422
	v2442 = v1871 + int32(base.Ui32(v2415)>>(uint(int32(6))%32))&v2418 + int32(2048)
	v2443 = *(*int32)(unsafe.Add(mBase, uint32(v2442)))
	*(*int32)(unsafe.Add(mBase, uint32(v2442))) = v2443 + v2422
	v2453 = v1871 + v2415&int32(255)<<(uint(int32(2))%32) + int32(3072)
	v2454 = *(*int32)(unsafe.Add(mBase, uint32(v2453)))
	*(*int32)(unsafe.Add(mBase, uint32(v2453))) = v2454 + v2422
	v2461 = v2339 + int32(-1)
	if v2461 != 0 {
		v2339 = v2461
		v2341 = v2341 + int32(4)
		goto L145
	} else {
		goto L147
	}
L146:
	;
	if v1596 <= l2 {
		goto L148
	} else {
		goto L149
	}
L147:
	;
	goto L146
L148:
	;
	v2814 = v1901 + int32(1)
	if v2814 != v1759 {
		__phi1881 = v1948
		__phi1901 = v2814
		__phi1948 = v1881
		v1881 = __phi1881
		v1901 = __phi1901
		v1948 = __phi1948
		goto L84
	} else {
		goto L156
	}
L149:
	;
	v2496 = int32(1)
	goto L150
L150:
	;
	v2566 = v1570 + (v2496*int32(14)+v1792)<<(uint(int32(12))%32)
	v2590 = v100 + int32(96)
	v2591 = v1751
	goto L152
L151:
	;
	goto L148
L152:
	;
	v2666 = *(*int32)(unsafe.Add(mBase, uint32(v2590)))
	v2669 = int32(1020)
	v2671 = v2566 + int32(base.Ui32(v2666)>>(uint(int32(22))%32))&v2669
	v2672 = *(*int32)(unsafe.Add(mBase, uint32(v2671)))
	v2673 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2671))) = v2672 + v2673
	v2682 = v2566 + int32(base.Ui32(v2666)>>(uint(int32(14))%32))&v2669 + int32(1024)
	v2683 = *(*int32)(unsafe.Add(mBase, uint32(v2682)))
	*(*int32)(unsafe.Add(mBase, uint32(v2682))) = v2683 + v2673
	v2693 = v2566 + int32(base.Ui32(v2666)>>(uint(int32(6))%32))&v2669 + int32(2048)
	v2694 = *(*int32)(unsafe.Add(mBase, uint32(v2693)))
	*(*int32)(unsafe.Add(mBase, uint32(v2693))) = v2694 + v2673
	v2704 = v2566 + v2666&int32(255)<<(uint(int32(2))%32) + int32(3072)
	v2705 = *(*int32)(unsafe.Add(mBase, uint32(v2704)))
	*(*int32)(unsafe.Add(mBase, uint32(v2704))) = v2705 + v2673
	v2712 = v2591 + int32(-1)
	if v2712 != 0 {
		v2590 = v2590 + int32(4)
		v2591 = v2712
		goto L152
	} else {
		goto L154
	}
L153:
	;
	v2714 = v2496 + int32(1)
	if v2714 != v1605 {
		v2496 = v2714
		goto L150
	} else {
		goto L155
	}
L154:
	;
	goto L153
L155:
	;
	goto L151
L156:
	;
	goto L85
L157:
	;
	goto L72
L158:
	;
	v3048 = int32(2)
	v3077 = v1610
	v3078 = l6
	v3088 = int32(0)
	goto L168
L159:
	;
	v2823 = int32(2)
	v2853 = v1610
	v2863 = int32(0)
	v2864 = l6
	goto L160
L160:
	;
	if v1742 < int32(1) {
		goto L162
	} else {
		goto L163
	}
L162:
	;
	__phi2940 = l5 + (v1607*v1742 + v1735<<(uint(v2823)%32) - v1737<<(uint(v2823)%32))
	__phi2954 = v2853
	__phi2955 = v1742
	__phi2965 = v2864
	__phi2966 = v1759
	v2940 = __phi2940
	v2954 = __phi2954
	v2955 = __phi2955
	v2965 = __phi2965
	v2966 = __phi2966
	goto L164
L163:
	;
	v2931 = F_memcpy(m, v2853+v1740, v1748, v1756)
	mBase = m.M
	goto L162
L164:
	;
	v3031 = v2955 + int32(1)
	v3036 = F_memcpy(m, v2965+v1740, v2940, (v1752+base.B2i32(v3031 < l1))<<(uint(int32(2))%32))
	mBase = m.M
	F_GetResidual(m, l0, l1, v2954, v2965, v1613, v2863, v1735, v1763, v2955, v107, l9, l10, v100+int32(96))
	mBase = m.M
	v3042 = v2966 + int32(-1)
	if v3042 != 0 {
		__phi2940 = v2940 + v1607
		__phi2954 = v2965
		__phi2955 = v3031
		__phi2965 = v2954
		__phi2966 = v3042
		v2940 = __phi2940
		v2954 = __phi2954
		v2955 = __phi2955
		v2965 = __phi2965
		v2966 = __phi2966
		goto L164
	} else {
		goto L166
	}
L165:
	;
	v3044 = v2863 + int32(1)
	if v3044 != int32(14) {
		v2853 = v2965
		v2863 = v3044
		v2864 = v2954
		goto L160
	} else {
		goto L167
	}
L166:
	;
	goto L165
L167:
	;
	goto L72
L168:
	;
	if v1742 < int32(1) {
		goto L170
	} else {
		goto L171
	}
L170:
	;
	__phi3165 = l5 + (v1607*v1742 + v1735<<(uint(v3048)%32) - v1737<<(uint(v3048)%32))
	__phi3178 = v3077
	__phi3179 = v3078
	__phi3190 = v1742
	__phi3191 = v1759
	v3165 = __phi3165
	v3178 = __phi3178
	v3179 = __phi3179
	v3190 = __phi3190
	v3191 = __phi3191
	goto L172
L171:
	;
	v3156 = F_memcpy(m, v3077+v1740, v1748, v1756)
	mBase = m.M
	goto L170
L172:
	;
	v3255 = int32(1)
	v3256 = v3190 + v3255
	v3261 = F_memcpy(m, v3179+v1740, v3165, (v1752+base.B2i32(v3256 < l1))<<(uint(int32(2))%32))
	mBase = m.M
	if v3190 < v3255 {
		goto L174
	} else {
		goto L175
	}
L173:
	;
	v3597 = v3088 + int32(1)
	if v3597 != int32(14) {
		v3077 = v3179
		v3078 = v3178
		v3088 = v3597
		goto L168
	} else {
		goto L233
	}
L174:
	;
	F_GetResidual(m, l0, l1, v3178, v3179, v1613, v3088, v1735, v1763, v3190, v107, l9, l10, v100+int32(96))
	mBase = m.M
	v3595 = v3191 + int32(-1)
	if v3595 != 0 {
		__phi3165 = v3165 + v1607
		__phi3178 = v3179
		__phi3179 = v3178
		__phi3190 = v3256
		__phi3191 = v3595
		v3165 = __phi3165
		v3178 = __phi3178
		v3179 = __phi3179
		v3190 = __phi3190
		v3191 = __phi3191
		goto L172
	} else {
		goto L232
	}
L175:
	;
	if l1 <= v3256 {
		goto L174
	} else {
		goto L176
	}
L176:
	;
	if v1764 < int32(3) {
		goto L178
	} else {
		goto L179
	}
L177:
	;
	goto L174
L178:
	;
	goto L177
L179:
	;
	v3282 = int32(0)
	v3283 = base.Simd_g_v128_load64_zero(m, v3165, v3282)
	if l10 == v3282 {
		v3300 = v3283
		goto L180
	} else {
		goto L181
	}
L180:
	;
	v3305 = int32(1)
	v3309 = int32(0)
	v3311 = l0 << (uint(int32(2)) % 32)
	__phi3317 = base.Simd_g_i32x4_extract_lane_l1(v3300)
	__phi3318 = base.Simd_g_i32x4_extract_lane_l0(v3300)
	__phi3325 = v3165 + int32(4)
	__phi3326 = v1762 + v3305
	__phi3327 = v1764 + int32(-2)
	v3317 = __phi3317
	v3318 = __phi3318
	v3325 = __phi3325
	v3326 = __phi3326
	v3327 = __phi3327
	goto L182
L181:
	;
	v3287 = base.Simd_g_i32x4_shr_u(v3283, int32(8))
	v3290 = base.Simd_g_const(&F_VP8LResidualImage__k6)
	v3300 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v3287, base.Simd_g_const(&F_VP8LResidualImage__k7)), base.Simd_g_v128_and(v3283, v3290)), base.Simd_g_i32x4_shl(v3287, int32(16))), v3290), base.Simd_g_v128_and(v3283, base.Simd_g_const(&F_VP8LResidualImage__k8)))
	goto L180
L182:
	;
	v3338 = v3325 + int32(4)
	v3339 = *(*int32)(unsafe.Add(mBase, uint32(v3338)))
	v3341 = *(*int32)(unsafe.Add(mBase, uint32(v3325+v3311)))
	v3343 = *(*int32)(unsafe.Add(mBase, uint32(v3325+(v3309-v3311))))
	if l10 == int32(0) {
		v3391 = v3339
		v3392 = v3341
		v3393 = v3343
		goto L184
	} else {
		goto L185
	}
L183:
	;
	goto L178
L184:
	;
	v3395 = int32(24)
	v3396 = int32(base.Ui32(v3317) >> (uint(v3395) % 32))
	v3399 = v3396 - int32(base.Ui32(v3393)>>(uint(v3395)%32))
	v3400 = int32(31)
	v3401 = v3399 >> (uint(v3400) % 32)
	v3403 = v3399 ^ v3401 - v3401
	v3404 = int32(16)
	v3406 = int32(255)
	v3407 = int32(base.Ui32(v3317)>>(uint(v3404)%32)) & v3406
	v3412 = v3407 - int32(base.Ui32(v3393)>>(uint(v3404)%32))&v3406
	v3414 = v3412 >> (uint(v3400) % 32)
	v3416 = v3412 ^ v3414 - v3414
	if base.Ui32(v3416) < base.Ui32(v3403) {
		goto L186
	} else {
		goto L187
	}
L185:
	;
	v3346 = int32(8)
	v3347 = int32(base.Ui32(v3339) >> (uint(v3346) % 32))
	v3348 = int32(255)
	v3350 = int32(16711935)
	v3353 = int32(16)
	v3358 = int32(-16711936)
	v3362 = int32(base.Ui32(v3341) >> (uint(v3346) % 32))
	v3377 = int32(base.Ui32(v3343) >> (uint(v3346) % 32))
	v3391 = (v3347&v3348+v3339&v3350+v3347<<(uint(v3353)%32))&v3350 | v3339&v3358
	v3392 = (v3362&v3348+v3341&v3350+v3362<<(uint(v3353)%32))&v3350 | v3341&v3358
	v3393 = (v3377&v3348+v3343&v3350+v3377<<(uint(v3353)%32))&v3350 | v3343&v3358
	goto L184
L186:
	;
	v3418 = v3403
	goto L188
L187:
	;
	v3418 = v3416
	goto L188
L188:
	;
	v3419 = int32(8)
	v3421 = int32(255)
	v3422 = int32(base.Ui32(v3317)>>(uint(v3419)%32)) & v3421
	v3427 = v3422 - int32(base.Ui32(v3393)>>(uint(v3419)%32))&v3421
	v3428 = int32(31)
	v3429 = v3427 >> (uint(v3428) % 32)
	v3431 = v3427 ^ v3429 - v3429
	v3433 = v3317 & v3421
	v3436 = v3433 - v3393&v3421
	v3438 = v3436 >> (uint(v3428) % 32)
	v3440 = v3436 ^ v3438 - v3438
	if base.Ui32(v3440) < base.Ui32(v3431) {
		goto L189
	} else {
		goto L190
	}
L189:
	;
	v3442 = v3431
	goto L191
L190:
	;
	v3442 = v3440
	goto L191
L191:
	;
	if base.Ui32(v3442) < base.Ui32(v3418) {
		goto L192
	} else {
		goto L193
	}
L192:
	;
	v3444 = v3418
	goto L194
L193:
	;
	v3444 = v3442
	goto L194
L194:
	;
	v3447 = v3396 - int32(base.Ui32(v3392)>>(uint(int32(24))%32))
	v3448 = int32(31)
	v3449 = v3447 >> (uint(v3448) % 32)
	v3451 = v3447 ^ v3449 - v3449
	v3456 = v3407 - int32(base.Ui32(v3392)>>(uint(int32(16))%32))&int32(255)
	v3458 = v3456 >> (uint(v3448) % 32)
	v3460 = v3456 ^ v3458 - v3458
	if base.Ui32(v3460) < base.Ui32(v3451) {
		goto L195
	} else {
		goto L196
	}
L195:
	;
	v3462 = v3451
	goto L197
L196:
	;
	v3462 = v3460
	goto L197
L197:
	;
	v3465 = int32(255)
	v3467 = v3422 - int32(base.Ui32(v3392)>>(uint(int32(8))%32))&v3465
	v3468 = int32(31)
	v3469 = v3467 >> (uint(v3468) % 32)
	v3471 = v3467 ^ v3469 - v3469
	v3474 = v3433 - v3392&v3465
	v3476 = v3474 >> (uint(v3468) % 32)
	v3478 = v3474 ^ v3476 - v3476
	if base.Ui32(v3478) < base.Ui32(v3471) {
		goto L198
	} else {
		goto L199
	}
L198:
	;
	v3480 = v3471
	goto L200
L199:
	;
	v3480 = v3478
	goto L200
L200:
	;
	if base.Ui32(v3480) < base.Ui32(v3462) {
		goto L201
	} else {
		goto L202
	}
L201:
	;
	v3482 = v3462
	goto L203
L202:
	;
	v3482 = v3480
	goto L203
L203:
	;
	if base.Ui32(v3482) < base.Ui32(v3444) {
		goto L204
	} else {
		goto L205
	}
L204:
	;
	v3484 = v3444
	goto L206
L205:
	;
	v3484 = v3482
	goto L206
L206:
	;
	v3487 = v3396 - int32(base.Ui32(v3318)>>(uint(int32(24))%32))
	v3488 = int32(31)
	v3489 = v3487 >> (uint(v3488) % 32)
	v3491 = v3487 ^ v3489 - v3489
	v3496 = v3407 - int32(base.Ui32(v3318)>>(uint(int32(16))%32))&int32(255)
	v3498 = v3496 >> (uint(v3488) % 32)
	v3500 = v3496 ^ v3498 - v3498
	if base.Ui32(v3500) < base.Ui32(v3491) {
		goto L207
	} else {
		goto L208
	}
L207:
	;
	v3502 = v3491
	goto L209
L208:
	;
	v3502 = v3500
	goto L209
L209:
	;
	v3505 = int32(255)
	v3507 = v3422 - int32(base.Ui32(v3318)>>(uint(int32(8))%32))&v3505
	v3508 = int32(31)
	v3509 = v3507 >> (uint(v3508) % 32)
	v3511 = v3507 ^ v3509 - v3509
	v3514 = v3433 - v3318&v3505
	v3516 = v3514 >> (uint(v3508) % 32)
	v3518 = v3514 ^ v3516 - v3516
	if base.Ui32(v3518) < base.Ui32(v3511) {
		goto L210
	} else {
		goto L211
	}
L210:
	;
	v3520 = v3511
	goto L212
L211:
	;
	v3520 = v3518
	goto L212
L212:
	;
	if base.Ui32(v3520) < base.Ui32(v3502) {
		goto L213
	} else {
		goto L214
	}
L213:
	;
	v3522 = v3502
	goto L215
L214:
	;
	v3522 = v3520
	goto L215
L215:
	;
	v3525 = v3396 - int32(base.Ui32(v3391)>>(uint(int32(24))%32))
	v3526 = int32(31)
	v3527 = v3525 >> (uint(v3526) % 32)
	v3529 = v3525 ^ v3527 - v3527
	v3534 = v3407 - int32(base.Ui32(v3391)>>(uint(int32(16))%32))&int32(255)
	v3536 = v3534 >> (uint(v3526) % 32)
	v3538 = v3534 ^ v3536 - v3536
	if base.Ui32(v3538) < base.Ui32(v3529) {
		goto L216
	} else {
		goto L217
	}
L216:
	;
	v3540 = v3529
	goto L218
L217:
	;
	v3540 = v3538
	goto L218
L218:
	;
	v3543 = int32(255)
	v3545 = v3422 - int32(base.Ui32(v3391)>>(uint(int32(8))%32))&v3543
	v3546 = int32(31)
	v3547 = v3545 >> (uint(v3546) % 32)
	v3549 = v3545 ^ v3547 - v3547
	v3552 = v3433 - v3391&v3543
	v3554 = v3552 >> (uint(v3546) % 32)
	v3556 = v3552 ^ v3554 - v3554
	if base.Ui32(v3556) < base.Ui32(v3549) {
		goto L219
	} else {
		goto L220
	}
L219:
	;
	v3558 = v3549
	goto L221
L220:
	;
	v3558 = v3556
	goto L221
L221:
	;
	if base.Ui32(v3558) < base.Ui32(v3540) {
		goto L222
	} else {
		goto L223
	}
L222:
	;
	v3560 = v3540
	goto L224
L223:
	;
	v3560 = v3558
	goto L224
L224:
	;
	if base.Ui32(v3560) < base.Ui32(v3522) {
		goto L225
	} else {
		goto L226
	}
L225:
	;
	v3562 = v3522
	goto L227
L226:
	;
	v3562 = v3560
	goto L227
L227:
	;
	if base.Ui32(v3562) < base.Ui32(v3484) {
		goto L228
	} else {
		goto L229
	}
L228:
	;
	v3564 = v3484
	goto L230
L229:
	;
	v3564 = v3562
	goto L230
L230:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3326))) = uint8(v3564)
	v3569 = v3327 + int32(-1)
	if v3569 != 0 {
		__phi3317 = v3391
		__phi3318 = v3317
		__phi3325 = v3338
		__phi3326 = v3326 + int32(1)
		__phi3327 = v3569
		v3317 = __phi3317
		v3318 = __phi3318
		v3325 = __phi3325
		v3326 = __phi3326
		v3327 = __phi3327
		goto L182
	} else {
		goto L231
	}
L231:
	;
	goto L183
L232:
	;
	goto L173
L233:
	;
	goto L72
L234:
	;
	v3603 = F_memcpy(m, v1610+v1740, v1748, v1756)
	mBase = m.M
	v3604 = F_memcpy(m, v3603, v1748, v1756)
	mBase = m.M
	v3605 = F_memcpy(m, v3604, v1748, v1756)
	mBase = m.M
	v3606 = F_memcpy(m, v3605, v1748, v1756)
	mBase = m.M
	v3607 = F_memcpy(m, v3606, v1748, v1756)
	mBase = m.M
	v3608 = F_memcpy(m, v3607, v1748, v1756)
	mBase = m.M
	v3609 = F_memcpy(m, v3608, v1748, v1756)
	mBase = m.M
	v3610 = F_memcpy(m, v3609, v1748, v1756)
	mBase = m.M
	v3611 = F_memcpy(m, v3610, v1748, v1756)
	mBase = m.M
	v3612 = F_memcpy(m, v3611, v1748, v1756)
	mBase = m.M
	v3613 = F_memcpy(m, v3612, v1748, v1756)
	mBase = m.M
	v3614 = F_memcpy(m, v3613, v1748, v1756)
	mBase = m.M
	v3615 = F_memcpy(m, v3614, v1748, v1756)
	mBase = m.M
	v3616 = F_memcpy(m, v3615, v1748, v1756)
	mBase = m.M
	goto L72
L235:
	;
	v4923 = v4920<<(uint(v1544)%32) + v4919
	v4925 = v4918 + v4921<<(uint(v1544)%32)
	if v4925 != 0 {
		goto L358
	} else {
		goto L359
	}
L236:
	;
	v4774 = int32(0)
	v4776 = v3866 + int32(_a_F_VP8LResidualImage_1)
	if base.Ui32(v4776) < base.Ui32(int32(33)) {
		goto L338
	} else {
		goto L339
	}
L237:
	;
	v3817 = int32(1)
	v3818 = v3794 + l2
	v3821 = int32(base.Ui32(v3817<<(uint(v3818)%32)+v1578) >> (uint(v3818) % 32))
	v3822 = int32(base.Ui32(v1698) >> (uint(v3794) % 32))
	v3826 = *(*int32)(unsafe.Add(mBase, uint32(v1077+v3794<<(uint(int32(2))%32))))
	v3827 = int32(255)
	v3829 = int32(base.Ui32(v1699) >> (uint(v3794) % 32))
	if v3829 < v3817 {
		v3842 = v3827
		goto L239
	} else {
		goto L240
	}
L238:
	;
	v4645 = int32(0)
	if base.Ui32(v1631) < base.Ui32(int32(33)) {
		goto L321
	} else {
		goto L322
	}
L239:
	;
	v3845 = v1589 + v3794<<(uint(int32(12))%32)
	if v3822 < int32(1) {
		v3858 = v3827
		goto L241
	} else {
		goto L242
	}
L240:
	;
	v3833 = int32(2)
	v3841 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3826+v3821*v3822<<(uint(v3833)%32)+v3829<<(uint(v3833)%32)+int32(-3)))))
	v3842 = v3841
	goto L239
L241:
	;
	v3866 = v3794 * int32(_a_F_VP8LResidualImage_1)
	v3867 = v1570 + v3866
	v3869 = int32(0)
	v3879 = v3802
	v3891 = v1570
	v3895 = v3801
	v3896 = v3803
	v3899 = v3869
	v3903 = v3800
	v3904 = v3869
	v3905 = v3799
	v3958 = int64(9223372036854775807)
	goto L243
L242:
	;
	v3851 = int32(2)
	v3857 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3826+v3821*(v3822+int32(-1))<<(uint(v3851)%32)+v3829<<(uint(v3851)%32))+1)))
	v3858 = v3857
	goto L241
L243:
	;
	v3970 = v3867 + v3904<<(uint(int32(12))%32)
	v3971 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3970))))
	v3997 = v3896
	v3999 = int32(4)
	v4064 = v3971 << (uint(int64(23)) % 64)
	v4065 = int64(788529152)
	goto L245
L244:
	;
	v4587 = m.G98
	v4588 = *(*int32)(unsafe.Add(mBase, uint32(v4587)))
	m.T0[v4588].(func(*base.Module, int32, int32, int32))(m, v4570, v3845, int32(1024))
	mBase = m.M
	v4591 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v3826+v3821*v3822<<(uint(v4591)%32)+v3829<<(uint(v4591)%32)))) = v4571<<(uint(int32(8))%32) | int32(-16777216)
	v4607 = v1589 + v1546<<(uint(int32(12))%32) + v3794*int32(56) + v4571<<(uint(v4591)%32)
	v4608 = *(*int32)(unsafe.Add(mBase, uint32(v4607)))
	*(*int32)(unsafe.Add(mBase, uint32(v4607))) = v4608 + int32(1)
	if v3794 == v1544 {
		goto L309
	} else {
		goto L310
	}
L245:
	;
	v4075 = *(*int32)(unsafe.Add(mBase, uint32(v3997)))
	v4077 = *(*int32)(unsafe.Add(mBase, uint32(v3879+v3999)))
	v4080 = v4065 * base.I64_extend_i32_u(v4075+v4077)
	if v4080 < int64(0) {
		goto L247
	} else {
		goto L248
	}
L246:
	;
	v4106 = m.G99
	v4107 = *(*int32)(unsafe.Add(mBase, uint32(v4106)))
	v4108 = m.T0[v4107].(func(*base.Module, int32, int32) int64)(m, v3970, v3845)
	mBase = m.M
	v4109 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3970)+1024)))
	v4135 = v3895
	v4137 = int32(1028)
	v4203 = int64(788529152)
	v4204 = v4109 << (uint(int64(23)) % 64)
	goto L254
L247:
	;
	v4083 = int64(-50)
	goto L249
L248:
	;
	v4083 = int64(50)
	goto L249
L249:
	;
	v4086 = base.I64_div_s(v4083+v4080, int64(100))
	v4087 = v4086 + v4064
	if v4065 < int64(0) {
		goto L250
	} else {
		goto L251
	}
L250:
	;
	v4094 = int64(-5)
	goto L252
L251:
	;
	v4094 = int64(5)
	goto L252
L252:
	;
	v4099 = base.I64_div_s(v4094+v4065*int64(6), int64(10))
	v4101 = v3999 + int32(4)
	if v4101 != int32(64) {
		v3997 = v3997 + int32(-4)
		v3999 = v4101
		v4064 = v4087
		v4065 = v4099
		goto L245
	} else {
		goto L253
	}
L253:
	;
	goto L246
L254:
	;
	v4213 = *(*int32)(unsafe.Add(mBase, uint32(v4135)))
	v4215 = *(*int32)(unsafe.Add(mBase, uint32(v3879+v4137)))
	v4218 = v4203 * base.I64_extend_i32_u(v4213+v4215)
	if v4218 < int64(0) {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	v4244 = m.G99
	v4245 = *(*int32)(unsafe.Add(mBase, uint32(v4244)))
	v4246 = m.T0[v4245].(func(*base.Module, int32, int32) int64)(m, v3970+int32(1024), v3845+int32(1024))
	mBase = m.M
	v4247 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3970)+2048)))
	v4273 = v3903
	v4275 = int32(2052)
	v4341 = int64(788529152)
	v4344 = v4247 << (uint(int64(23)) % 64)
	goto L263
L256:
	;
	v4221 = int64(-50)
	goto L258
L257:
	;
	v4221 = int64(50)
	goto L258
L258:
	;
	v4224 = base.I64_div_s(v4221+v4218, int64(100))
	v4225 = v4224 + v4204
	if v4203 < int64(0) {
		goto L259
	} else {
		goto L260
	}
L259:
	;
	v4232 = int64(-5)
	goto L261
L260:
	;
	v4232 = int64(5)
	goto L261
L261:
	;
	v4237 = base.I64_div_s(v4232+v4203*int64(6), int64(10))
	v4239 = v4137 + int32(4)
	if v4239 != int32(1088) {
		v4135 = v4135 + int32(-4)
		v4137 = v4239
		v4203 = v4237
		v4204 = v4225
		goto L254
	} else {
		goto L262
	}
L262:
	;
	goto L255
L263:
	;
	v4351 = *(*int32)(unsafe.Add(mBase, uint32(v4273)))
	v4353 = *(*int32)(unsafe.Add(mBase, uint32(v3879+v4275)))
	v4356 = v4341 * base.I64_extend_i32_u(v4351+v4353)
	if v4356 < int64(0) {
		goto L265
	} else {
		goto L266
	}
L264:
	;
	v4382 = m.G99
	v4383 = *(*int32)(unsafe.Add(mBase, uint32(v4382)))
	v4384 = m.T0[v4383].(func(*base.Module, int32, int32) int64)(m, v3970+int32(2048), v3845+int32(2048))
	mBase = m.M
	v4385 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3970)+3072)))
	v4411 = v3905
	v4413 = int32(3076)
	v4479 = int64(788529152)
	v4484 = v4385 << (uint(int64(23)) % 64)
	goto L272
L265:
	;
	v4359 = int64(-50)
	goto L267
L266:
	;
	v4359 = int64(50)
	goto L267
L267:
	;
	v4362 = base.I64_div_s(v4359+v4356, int64(100))
	v4363 = v4362 + v4344
	if v4341 < int64(0) {
		goto L268
	} else {
		goto L269
	}
L268:
	;
	v4370 = int64(-5)
	goto L270
L269:
	;
	v4370 = int64(5)
	goto L270
L270:
	;
	v4375 = base.I64_div_s(v4370+v4341*int64(6), int64(10))
	v4377 = v4275 + int32(4)
	if v4377 != int32(2112) {
		v4273 = v4273 + int32(-4)
		v4275 = v4377
		v4341 = v4375
		v4344 = v4363
		goto L263
	} else {
		goto L271
	}
L271:
	;
	goto L264
L272:
	;
	v4489 = *(*int32)(unsafe.Add(mBase, uint32(v4411)))
	v4491 = *(*int32)(unsafe.Add(mBase, uint32(v3879+v4413)))
	v4494 = v4479 * base.I64_extend_i32_u(v4489+v4491)
	if v4494 < int64(0) {
		goto L274
	} else {
		goto L275
	}
L273:
	;
	if v4087 < int64(0) {
		goto L281
	} else {
		goto L282
	}
L274:
	;
	v4497 = int64(-50)
	goto L276
L275:
	;
	v4497 = int64(50)
	goto L276
L276:
	;
	v4500 = base.I64_div_s(v4497+v4494, int64(100))
	v4501 = v4500 + v4484
	if v4479 < int64(0) {
		goto L277
	} else {
		goto L278
	}
L277:
	;
	v4508 = int64(-5)
	goto L279
L278:
	;
	v4508 = int64(5)
	goto L279
L279:
	;
	v4513 = base.I64_div_s(v4508+v4479*int64(6), int64(10))
	v4515 = v4413 + int32(4)
	if v4515 != int32(3136) {
		v4411 = v4411 + int32(-4)
		v4413 = v4515
		v4479 = v4513
		v4484 = v4501
		goto L272
	} else {
		goto L280
	}
L280:
	;
	goto L273
L281:
	;
	v4522 = int64(-5)
	goto L283
L282:
	;
	v4522 = int64(5)
	goto L283
L283:
	;
	v4525 = base.I64_div_s(v4522+v4087, int64(-10))
	if v4225 < int64(0) {
		goto L284
	} else {
		goto L285
	}
L284:
	;
	v4531 = int64(-5)
	goto L286
L285:
	;
	v4531 = int64(5)
	goto L286
L286:
	;
	v4534 = base.I64_div_s(v4531+v4225, int64(-10))
	if v4363 < int64(0) {
		goto L287
	} else {
		goto L288
	}
L287:
	;
	v4541 = int64(-5)
	goto L289
L288:
	;
	v4541 = int64(5)
	goto L289
L289:
	;
	v4544 = base.I64_div_s(v4541+v4363, int64(-10))
	if v4501 < int64(0) {
		goto L290
	} else {
		goto L291
	}
L290:
	;
	v4551 = int64(-5)
	goto L292
L291:
	;
	v4551 = int64(5)
	goto L292
L292:
	;
	v4554 = base.I64_div_s(v4551+v4501, int64(-10))
	v4556 = m.G99
	v4557 = *(*int32)(unsafe.Add(mBase, uint32(v4556)))
	v4558 = m.T0[v4557].(func(*base.Module, int32, int32) int64)(m, v3970+int32(3072), v3845+int32(3072))
	mBase = m.M
	v4559 = v4108 + v4525 + v4534 + v4246 + v4544 + v4384 + v4554 + v4558
	if v3904 == v3842 {
		goto L293
	} else {
		goto L294
	}
L293:
	;
	v4563 = v4559 + int64(-125829120)
	goto L295
L294:
	;
	v4563 = v4559
	goto L295
L295:
	;
	if v3904 == v3858 {
		goto L296
	} else {
		goto L297
	}
L296:
	;
	v4567 = v4563 + int64(-125829120)
	goto L298
L297:
	;
	v4567 = v4563
	goto L298
L298:
	;
	v4568 = base.B2i32(v4567 < v3958)
	if v4567 < v3958 {
		goto L299
	} else {
		goto L300
	}
L299:
	;
	v4569 = v4567
	goto L301
L300:
	;
	v4569 = v3958
	goto L301
L301:
	;
	if v4567 < v3958 {
		goto L302
	} else {
		goto L303
	}
L302:
	;
	v4570 = v3970
	goto L304
L303:
	;
	v4570 = v3891
	goto L304
L304:
	;
	if v4567 < v3958 {
		goto L305
	} else {
		goto L306
	}
L305:
	;
	v4571 = v3904
	goto L307
L306:
	;
	v4571 = v3899
	goto L307
L307:
	;
	v4572 = int32(_a_F_VP8LResidualImage_3)
	v4583 = v3904 + int32(1)
	if v4583 != int32(14) {
		v3879 = v3879 + v4572
		v3891 = v4570
		v3895 = v3895 + v4572
		v3896 = v3896 + v4572
		v3899 = v4571
		v3903 = v3903 + v4572
		v3904 = v4583
		v3905 = v3905 + v4572
		v3958 = v4569
		goto L243
	} else {
		goto L308
	}
L308:
	;
	goto L244
L309:
	;
	goto L238
L310:
	;
	v4614 = v3794 + int32(1)
	if base.Ui32(v4614) <= base.Ui32(v1599) {
		goto L311
	} else {
		goto L312
	}
L311:
	;
	if v1699 == v1627 {
		goto L314
	} else {
		goto L315
	}
L312:
	;
	if base.Ui32(v1544) < base.Ui32(v4614) {
		goto L311
	} else {
		goto L313
	}
L313:
	;
	v4621 = m.G98
	v4622 = *(*int32)(unsafe.Add(mBase, uint32(v4621)))
	m.T0[v4622].(func(*base.Module, int32, int32, int32))(m, v3867, v1570+v4614*int32(_a_F_VP8LResidualImage_1), int32(_a_F_VP8LResidualImage_4))
	mBase = m.M
	goto L311
L314:
	;
	if v1698 == v1584+v1622 {
		goto L317
	} else {
		goto L318
	}
L315:
	;
	v4625 = int32(-1)
	if v4625<<(uint(v4614)%32)|(v3714-v1696) != v4625 {
		goto L236
	} else {
		goto L316
	}
L316:
	;
	goto L314
L317:
	;
	v4635 = int32(_a_F_VP8LResidualImage_1)
	v3794 = v4614
	v3799 = v3799 + v4635
	v3800 = v3800 + v4635
	v3801 = v3801 + v4635
	v3802 = v3802 + v4635
	v3803 = v3803 + v4635
	goto L237
L318:
	;
	v4630 = int32(-1)
	if v4630<<(uint(v4614)%32)|(v3714-v1697) != v4630 {
		goto L236
	} else {
		goto L319
	}
L319:
	;
	goto L317
L320:
	;
	v4770 = base.B2i32(v1699 == v1627)
	if v1699 == v1627 {
		goto L334
	} else {
		goto L335
	}
L321:
	;
	if v1631 == int32(0) {
		goto L323
	} else {
		goto L324
	}
L322:
	;
	base.MemoryFill(m, v1570, v4645, v1631)
	goto L320
L323:
	;
	goto L320
L324:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1570))) = uint8(v4645)
	v4657 = v1570 + v1631
	*(*uint8)(unsafe.Add(mBase, uint32(v4657+int32(-1)))) = uint8(v4645)
	if base.Ui32(v1631) < base.Ui32(int32(3)) {
		goto L323
	} else {
		goto L325
	}
L325:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1570)+2)) = uint8(v4645)
	*(*uint8)(unsafe.Add(mBase, uint32(v1570)+1)) = uint8(v4645)
	*(*uint8)(unsafe.Add(mBase, uint32(v4657+int32(-3)))) = uint8(v4645)
	*(*uint8)(unsafe.Add(mBase, uint32(v4657+int32(-2)))) = uint8(v4645)
	if base.Ui32(v1631) < base.Ui32(int32(7)) {
		goto L323
	} else {
		goto L326
	}
L326:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1570)+3)) = uint8(v4645)
	*(*uint8)(unsafe.Add(mBase, uint32(v4657+int32(-4)))) = uint8(v4645)
	if base.Ui32(v1631) < base.Ui32(int32(9)) {
		goto L323
	} else {
		goto L327
	}
L327:
	;
	v4679 = int32(0)
	v4682 = (v4679 - v1570) & int32(3)
	v4683 = v1570 + v4682
	*(*int32)(unsafe.Add(mBase, uint32(v4683))) = v4679
	v4691 = (v1631 - v4682) & int32(60)
	v4692 = v4683 + v4691
	*(*int32)(unsafe.Add(mBase, uint32(v4692+int32(-4)))) = v4679
	if base.Ui32(v4691) < base.Ui32(int32(9)) {
		goto L323
	} else {
		goto L328
	}
L328:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4683)+8)) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4683)+4)) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4692+int32(-8)))) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4692+int32(-12)))) = v4679
	if base.Ui32(v4691) < base.Ui32(int32(25)) {
		goto L323
	} else {
		goto L329
	}
L329:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4683)+24)) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4683)+20)) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4683)+16)) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4683)+12)) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4692+int32(-16)))) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4692+int32(-20)))) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4692+int32(-24)))) = v4679
	*(*int32)(unsafe.Add(mBase, uint32(v4692+int32(-28)))) = v4679
	v4727 = v4683&int32(4) | int32(24)
	v4728 = v4691 - v4727
	if base.Ui32(v4728) < base.Ui32(int32(32)) {
		goto L323
	} else {
		goto L330
	}
L330:
	;
	v4733 = base.I64_extend_i32_u(v4679) * int64(4294967297)
	v4736 = v4728
	v4737 = v4683 + v4727
	goto L331
L331:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4737)+24)) = v4733
	*(*int64)(unsafe.Add(mBase, uint32(v4737)+16)) = v4733
	*(*int64)(unsafe.Add(mBase, uint32(v4737)+8)) = v4733
	*(*int64)(unsafe.Add(mBase, uint32(v4737))) = v4733
	v4749 = v4736 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v4749) {
		v4736 = v4749
		v4737 = v4737 + int32(32)
		goto L331
	} else {
		goto L333
	}
L332:
	;
	goto L323
L333:
	;
	goto L332
L334:
	;
	v4771 = int32(0)
	goto L336
L335:
	;
	v4771 = v1701 + int32(1)
	goto L336
L336:
	;
	v4918 = int32(0)
	v4919 = v4645
	v4920 = v1700 + v4770
	v4921 = v4771
	goto L235
L337:
	;
	v4897 = int32(base.Ui32(v1697) >> (uint(v3794) % 32))
	v4898 = int32(base.Ui32(v1696) >> (uint(v3794) % 32))
	v4900 = v4898 & int32(1)
	if v1699 != v1627 {
		goto L353
	} else {
		goto L354
	}
L338:
	;
	if v4776 == int32(0) {
		goto L340
	} else {
		goto L341
	}
L339:
	;
	base.MemoryFill(m, v1570, v4774, v4776)
	goto L337
L340:
	;
	goto L337
L341:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1570))) = uint8(v4774)
	v4787 = v1570 + v4776
	*(*uint8)(unsafe.Add(mBase, uint32(v4787+int32(-1)))) = uint8(v4774)
	if base.Ui32(v4776) < base.Ui32(int32(3)) {
		goto L340
	} else {
		goto L342
	}
L342:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1570)+2)) = uint8(v4774)
	*(*uint8)(unsafe.Add(mBase, uint32(v1570)+1)) = uint8(v4774)
	*(*uint8)(unsafe.Add(mBase, uint32(v4787+int32(-3)))) = uint8(v4774)
	*(*uint8)(unsafe.Add(mBase, uint32(v4787+int32(-2)))) = uint8(v4774)
	if base.Ui32(v4776) < base.Ui32(int32(7)) {
		goto L340
	} else {
		goto L343
	}
L343:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1570)+3)) = uint8(v4774)
	*(*uint8)(unsafe.Add(mBase, uint32(v4787+int32(-4)))) = uint8(v4774)
	if base.Ui32(v4776) < base.Ui32(int32(9)) {
		goto L340
	} else {
		goto L344
	}
L344:
	;
	v4809 = int32(0)
	v4812 = (v4809 - v1570) & int32(3)
	v4813 = v1570 + v4812
	*(*int32)(unsafe.Add(mBase, uint32(v4813))) = v4809
	v4821 = (v4776 - v4812) & int32(60)
	v4822 = v4813 + v4821
	*(*int32)(unsafe.Add(mBase, uint32(v4822+int32(-4)))) = v4809
	if base.Ui32(v4821) < base.Ui32(int32(9)) {
		goto L340
	} else {
		goto L345
	}
L345:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4813)+8)) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4813)+4)) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4822+int32(-8)))) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4822+int32(-12)))) = v4809
	if base.Ui32(v4821) < base.Ui32(int32(25)) {
		goto L340
	} else {
		goto L346
	}
L346:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4813)+24)) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4813)+20)) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4813)+16)) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4813)+12)) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4822+int32(-16)))) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4822+int32(-20)))) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4822+int32(-24)))) = v4809
	*(*int32)(unsafe.Add(mBase, uint32(v4822+int32(-28)))) = v4809
	v4857 = v4813&int32(4) | int32(24)
	v4858 = v4821 - v4857
	if base.Ui32(v4858) < base.Ui32(int32(32)) {
		goto L340
	} else {
		goto L347
	}
L347:
	;
	v4863 = base.I64_extend_i32_u(v4809) * int64(4294967297)
	v4866 = v4858
	v4867 = v4813 + v4857
	goto L348
L348:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4867)+24)) = v4863
	*(*int64)(unsafe.Add(mBase, uint32(v4867)+16)) = v4863
	*(*int64)(unsafe.Add(mBase, uint32(v4867)+8)) = v4863
	*(*int64)(unsafe.Add(mBase, uint32(v4867))) = v4863
	v4879 = v4866 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v4879) {
		v4866 = v4879
		v4867 = v4867 + int32(32)
		goto L348
	} else {
		goto L350
	}
L349:
	;
	goto L340
L350:
	;
	goto L349
L351:
	;
	v4918 = v4912 << (uint(v3794) % 32)
	v4919 = v4911 << (uint(v3794) % 32)
	v4920 = v1700
	v4921 = v1701
	goto L235
L352:
	;
	v4911 = v4897 + int32(1)
	v4912 = v4908
	goto L351
L353:
	;
	if v4900 != 0 {
		goto L356
	} else {
		goto L357
	}
L354:
	;
	if v4900 == int32(0) {
		v4908 = v4898
		goto L352
	} else {
		goto L355
	}
L355:
	;
	goto L353
L356:
	;
	v4908 = v4898 + int32(-1)
	goto L352
L357:
	;
	v4911 = v4897
	v4912 = v4898 | int32(1)
	goto L351
L358:
	;
	if base.Ui32(v4923) < base.Ui32(v1584) {
		v1696 = v4918
		v1697 = v4919
		v1698 = v4923
		v1699 = v4925
		v1700 = v4920
		v1701 = v4921
		goto L67
	} else {
		goto L362
	}
L359:
	;
	v4927 = base.I32_div_u_s(v4923*l12, v1584)
	v4929 = F_WebPReportProgress(m, l11, v4927+v1572, l13)
	mBase = m.M
	if v4929 != 0 {
		goto L358
	} else {
		goto L360
	}
L360:
	;
	F_free(m, v1570)
	mBase = m.M
	goto L361
L361:
	;
	v5191 = l0
	v5192 = l1
	v5195 = l4
	v5196 = l5
	v5197 = l6
	v5198 = l7
	v5200 = l9
	v5201 = l10
	v5202 = l11
	v5203 = l12
	v5204 = l13
	v5206 = v100
	v5207 = v107
	v5208 = v108
	v5209 = v1578
	v5225 = int32(0)
	v5228 = v1069
	v5233 = v1582
	goto L55
L362:
	;
	goto L68
L363:
	;
	v5144 = m.G119
	v5145 = *(*int32)(unsafe.Add(mBase, uint32(v5144)))
	v5146 = m.T0[v5145].(func(*base.Module, int32, int32) int64)(m, v1570+v5067, int32(14))
	mBase = m.M
	v5147 = v1570 + v5066
	v5150 = int32(256)
	v5151 = *(*int32)(unsafe.Add(mBase, uint32(v5144)))
	v5152 = m.T0[v5151].(func(*base.Module, int32, int32) int64)(m, v5147+int32(_a_F_VP8LResidualImage_1), v5150)
	mBase = m.M
	v5157 = *(*int32)(unsafe.Add(mBase, uint32(v5144)))
	v5158 = m.T0[v5157].(func(*base.Module, int32, int32) int64)(m, v5147+int32(_a_F_VP8LResidualImage_5), v5150)
	mBase = m.M
	v5163 = *(*int32)(unsafe.Add(mBase, uint32(v5144)))
	v5164 = m.T0[v5163].(func(*base.Module, int32, int32) int64)(m, v5147+int32(_a_F_VP8LResidualImage_6), v5150)
	mBase = m.M
	v5169 = *(*int32)(unsafe.Add(mBase, uint32(v5144)))
	v5170 = m.T0[v5169].(func(*base.Module, int32, int32) int64)(m, v5147+int32(_a_F_VP8LResidualImage_7), v5150)
	mBase = m.M
	v5171 = v5146 + v5152 + v5158 + v5164 + v5170
	if v5134 <= v5171 {
		v5176 = v5079
		v5177 = v5134
		goto L365
	} else {
		goto L366
	}
L364:
	;
	F_free(m, v1570)
	mBase = m.M
	goto L368
L365:
	;
	v5185 = v5078 + int32(1)
	if base.Ui32(v5185) <= base.Ui32(v1544) {
		v5066 = v5066 + int32(_a_F_VP8LResidualImage_3)
		v5067 = v5067 + int32(56)
		v5078 = v5185
		v5079 = v5176
		v5083 = v5083 + int32(4)
		v5134 = v5177
		goto L363
	} else {
		goto L367
	}
L366:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l14))) = l2 + v5078
	v5175 = *(*int32)(unsafe.Add(mBase, uint32(v5083)))
	v5176 = v5175
	v5177 = v5171
	goto L365
L367:
	;
	goto L364
L368:
	;
	v5188 = *(*int32)(unsafe.Add(mBase, uint32(l14)))
	F_VP8LOptimizeSampling(m, v5176, l0, l1, v5188, int32(9), l14)
	mBase = m.M
	v5191 = l0
	v5192 = l1
	v5195 = l4
	v5196 = l5
	v5197 = l6
	v5198 = l7
	v5200 = l9
	v5201 = l10
	v5202 = l11
	v5203 = l12
	v5204 = l13
	v5206 = v100
	v5207 = v107
	v5208 = v108
	v5209 = v1578
	v5225 = v5176
	v5228 = v1069
	v5233 = v1582
	goto L55
L369:
	;
	v5304 = v5206
	v5326 = v5228
	goto L53
L370:
	;
	v5402 = v5304
	goto L28
L371:
	;
	v5497 = *(*int32)(unsafe.Add(mBase, uint32(l14)))
	v5498 = int32(1) << (uint(v5497) % 32)
	v5501 = v5191
	v5502 = v5192
	v5504 = v5497
	v5505 = v5195
	v5506 = v5196
	v5507 = v5197
	v5508 = v5198
	v5510 = v5200
	v5511 = v5201
	v5512 = v5202
	v5513 = v5203
	v5514 = v5204
	v5516 = v5206
	v5517 = v5207
	v5518 = v5208
	v5519 = v5209
	v5520 = v5498
	v5521 = int32(base.Ui32(v5498+v5209) >> (uint(v5497) % 32))
	goto L2
L372:
	;
	v7090 = F_WebPReportProgress(m, v5512, v5518+v5513, v5514)
	mBase = m.M
	v7099 = v7090
	v7106 = v5516
	goto L1
L373:
	;
	v5601 = v5501 << (uint(int32(2)) % 32)
	v5604 = v5507 + v5601 + int32(4)
	if v5505 == int32(0) {
		goto L374
	} else {
		goto L375
	}
L374:
	;
	v5751 = v5604 + v5601 + int32(4)
	v5752 = v5751 + v5501
	if v5501 < int32(1) {
		goto L380
	} else {
		goto L381
	}
L375:
	;
	v5607 = int32(1)
	v5612 = F_memcpy(m, v5507, v5506, (base.B2i32(v5502 != v5607)+v5501)<<(uint(int32(2))%32))
	mBase = m.M
	v5613 = int32(0)
	v5615 = m.G102
	v5616 = *(*int32)(unsafe.Add(mBase, uint32(v5615)))
	m.T0[v5616].(func(*base.Module, int32, int32, int32, int32))(m, v5612, v5613, v5607, v5506)
	mBase = m.M
	v5618 = int32(4)
	v5623 = *(*int32)(unsafe.Add(mBase, uint32(v5615)+4))
	m.T0[v5623].(func(*base.Module, int32, int32, int32, int32))(m, v5612+v5618, v5613, v5519, v5506+v5618)
	mBase = m.M
	if v5502 == v5607 {
		goto L372
	} else {
		goto L376
	}
L376:
	;
	v5628 = v5501 << (uint(int32(2)) % 32)
	__phi5637 = v5507
	__phi5639 = v5506 + v5628
	__phi5653 = v5604
	__phi5654 = int32(1)
	v5637 = __phi5637
	v5639 = __phi5639
	v5653 = __phi5653
	v5654 = __phi5654
	goto L377
L377:
	;
	v5728 = int32(1)
	v5729 = v5654 + v5728
	v5734 = F_memcpy(m, v5653, v5639, (base.B2i32(v5729 < v5502)+v5501)<<(uint(int32(2))%32))
	mBase = m.M
	v5736 = m.G102
	v5737 = *(*int32)(unsafe.Add(mBase, uint32(v5736)+8))
	m.T0[v5737].(func(*base.Module, int32, int32, int32, int32))(m, v5734, v5637, v5728, v5639)
	mBase = m.M
	v5739 = int32(4)
	v5745 = *(*int32)(unsafe.Add(mBase, uint32(v5736)+44))
	m.T0[v5745].(func(*base.Module, int32, int32, int32, int32))(m, v5734+v5739, v5637+v5739, v5519, v5639+v5739)
	mBase = m.M
	if v5502 != v5729 {
		__phi5637 = v5734
		__phi5639 = v5639 + v5628
		__phi5653 = v5637
		__phi5654 = v5729
		v5637 = __phi5637
		v5639 = __phi5639
		v5653 = __phi5653
		v5654 = __phi5654
		goto L377
	} else {
		goto L379
	}
L379:
	;
	goto L372
L380:
	;
	if v5517 < int32(2) {
		goto L451
	} else {
		goto L452
	}
L381:
	;
	__phi5762 = v5507
	__phi5764 = v5752
	__phi5778 = v5604
	__phi5788 = v5751
	__phi5789 = int32(0)
	v5762 = __phi5762
	v5764 = __phi5764
	v5778 = __phi5778
	v5788 = __phi5788
	v5789 = __phi5789
	goto L382
L382:
	;
	v5854 = int32(2)
	v5856 = v5506 + v5789*v5501<<(uint(v5854)%32)
	v5858 = v5789 + int32(1)
	v5863 = F_memcpy(m, v5762, v5856, (base.B2i32(v5858 < v5502)+v5501)<<(uint(v5854)%32))
	mBase = m.M
	if v5854 <= v5517 {
		goto L385
	} else {
		goto L386
	}
L384:
	;
	v6214 = int32(0)
	goto L444
L385:
	;
	if v5502 <= v5789+int32(2) {
		goto L387
	} else {
		goto L388
	}
L386:
	;
	v6198 = v5764
	v6199 = v5788
	goto L384
L387:
	;
	v6198 = v5788
	v6199 = v5764
	goto L384
L388:
	;
	v5872 = v5506 + v5858*v5501<<(uint(int32(2))%32)
	if v5501 < int32(3) {
		goto L390
	} else {
		goto L391
	}
L389:
	;
	goto L387
L390:
	;
	goto L389
L391:
	;
	v5890 = int32(0)
	v5891 = base.Simd_g_v128_load64_zero(m, v5872, v5890)
	if v5511 == v5890 {
		v5908 = v5891
		goto L392
	} else {
		goto L393
	}
L392:
	;
	v5913 = int32(1)
	v5917 = int32(0)
	v5919 = v5501 << (uint(int32(2)) % 32)
	__phi5925 = base.Simd_g_i32x4_extract_lane_l1(v5908)
	__phi5926 = base.Simd_g_i32x4_extract_lane_l0(v5908)
	__phi5933 = v5872 + int32(4)
	__phi5934 = v5788 + v5913
	__phi5935 = v5501 + int32(-2)
	v5925 = __phi5925
	v5926 = __phi5926
	v5933 = __phi5933
	v5934 = __phi5934
	v5935 = __phi5935
	goto L394
L393:
	;
	v5895 = base.Simd_g_i32x4_shr_u(v5891, int32(8))
	v5898 = base.Simd_g_const(&F_VP8LResidualImage__k6)
	v5908 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v5895, base.Simd_g_const(&F_VP8LResidualImage__k7)), base.Simd_g_v128_and(v5891, v5898)), base.Simd_g_i32x4_shl(v5895, int32(16))), v5898), base.Simd_g_v128_and(v5891, base.Simd_g_const(&F_VP8LResidualImage__k8)))
	goto L392
L394:
	;
	v5946 = v5933 + int32(4)
	v5947 = *(*int32)(unsafe.Add(mBase, uint32(v5946)))
	v5949 = *(*int32)(unsafe.Add(mBase, uint32(v5933+v5919)))
	v5951 = *(*int32)(unsafe.Add(mBase, uint32(v5933+(v5917-v5919))))
	if v5511 == int32(0) {
		v5999 = v5947
		v6000 = v5949
		v6001 = v5951
		goto L396
	} else {
		goto L397
	}
L395:
	;
	goto L390
L396:
	;
	v6003 = int32(24)
	v6004 = int32(base.Ui32(v5925) >> (uint(v6003) % 32))
	v6007 = v6004 - int32(base.Ui32(v6001)>>(uint(v6003)%32))
	v6008 = int32(31)
	v6009 = v6007 >> (uint(v6008) % 32)
	v6011 = v6007 ^ v6009 - v6009
	v6012 = int32(16)
	v6014 = int32(255)
	v6015 = int32(base.Ui32(v5925)>>(uint(v6012)%32)) & v6014
	v6020 = v6015 - int32(base.Ui32(v6001)>>(uint(v6012)%32))&v6014
	v6022 = v6020 >> (uint(v6008) % 32)
	v6024 = v6020 ^ v6022 - v6022
	if base.Ui32(v6024) < base.Ui32(v6011) {
		goto L398
	} else {
		goto L399
	}
L397:
	;
	v5954 = int32(8)
	v5955 = int32(base.Ui32(v5947) >> (uint(v5954) % 32))
	v5956 = int32(255)
	v5958 = int32(16711935)
	v5961 = int32(16)
	v5966 = int32(-16711936)
	v5970 = int32(base.Ui32(v5949) >> (uint(v5954) % 32))
	v5985 = int32(base.Ui32(v5951) >> (uint(v5954) % 32))
	v5999 = (v5955&v5956+v5947&v5958+v5955<<(uint(v5961)%32))&v5958 | v5947&v5966
	v6000 = (v5970&v5956+v5949&v5958+v5970<<(uint(v5961)%32))&v5958 | v5949&v5966
	v6001 = (v5985&v5956+v5951&v5958+v5985<<(uint(v5961)%32))&v5958 | v5951&v5966
	goto L396
L398:
	;
	v6026 = v6011
	goto L400
L399:
	;
	v6026 = v6024
	goto L400
L400:
	;
	v6027 = int32(8)
	v6029 = int32(255)
	v6030 = int32(base.Ui32(v5925)>>(uint(v6027)%32)) & v6029
	v6035 = v6030 - int32(base.Ui32(v6001)>>(uint(v6027)%32))&v6029
	v6036 = int32(31)
	v6037 = v6035 >> (uint(v6036) % 32)
	v6039 = v6035 ^ v6037 - v6037
	v6041 = v5925 & v6029
	v6044 = v6041 - v6001&v6029
	v6046 = v6044 >> (uint(v6036) % 32)
	v6048 = v6044 ^ v6046 - v6046
	if base.Ui32(v6048) < base.Ui32(v6039) {
		goto L401
	} else {
		goto L402
	}
L401:
	;
	v6050 = v6039
	goto L403
L402:
	;
	v6050 = v6048
	goto L403
L403:
	;
	if base.Ui32(v6050) < base.Ui32(v6026) {
		goto L404
	} else {
		goto L405
	}
L404:
	;
	v6052 = v6026
	goto L406
L405:
	;
	v6052 = v6050
	goto L406
L406:
	;
	v6055 = v6004 - int32(base.Ui32(v6000)>>(uint(int32(24))%32))
	v6056 = int32(31)
	v6057 = v6055 >> (uint(v6056) % 32)
	v6059 = v6055 ^ v6057 - v6057
	v6064 = v6015 - int32(base.Ui32(v6000)>>(uint(int32(16))%32))&int32(255)
	v6066 = v6064 >> (uint(v6056) % 32)
	v6068 = v6064 ^ v6066 - v6066
	if base.Ui32(v6068) < base.Ui32(v6059) {
		goto L407
	} else {
		goto L408
	}
L407:
	;
	v6070 = v6059
	goto L409
L408:
	;
	v6070 = v6068
	goto L409
L409:
	;
	v6073 = int32(255)
	v6075 = v6030 - int32(base.Ui32(v6000)>>(uint(int32(8))%32))&v6073
	v6076 = int32(31)
	v6077 = v6075 >> (uint(v6076) % 32)
	v6079 = v6075 ^ v6077 - v6077
	v6082 = v6041 - v6000&v6073
	v6084 = v6082 >> (uint(v6076) % 32)
	v6086 = v6082 ^ v6084 - v6084
	if base.Ui32(v6086) < base.Ui32(v6079) {
		goto L410
	} else {
		goto L411
	}
L410:
	;
	v6088 = v6079
	goto L412
L411:
	;
	v6088 = v6086
	goto L412
L412:
	;
	if base.Ui32(v6088) < base.Ui32(v6070) {
		goto L413
	} else {
		goto L414
	}
L413:
	;
	v6090 = v6070
	goto L415
L414:
	;
	v6090 = v6088
	goto L415
L415:
	;
	if base.Ui32(v6090) < base.Ui32(v6052) {
		goto L416
	} else {
		goto L417
	}
L416:
	;
	v6092 = v6052
	goto L418
L417:
	;
	v6092 = v6090
	goto L418
L418:
	;
	v6095 = v6004 - int32(base.Ui32(v5926)>>(uint(int32(24))%32))
	v6096 = int32(31)
	v6097 = v6095 >> (uint(v6096) % 32)
	v6099 = v6095 ^ v6097 - v6097
	v6104 = v6015 - int32(base.Ui32(v5926)>>(uint(int32(16))%32))&int32(255)
	v6106 = v6104 >> (uint(v6096) % 32)
	v6108 = v6104 ^ v6106 - v6106
	if base.Ui32(v6108) < base.Ui32(v6099) {
		goto L419
	} else {
		goto L420
	}
L419:
	;
	v6110 = v6099
	goto L421
L420:
	;
	v6110 = v6108
	goto L421
L421:
	;
	v6113 = int32(255)
	v6115 = v6030 - int32(base.Ui32(v5926)>>(uint(int32(8))%32))&v6113
	v6116 = int32(31)
	v6117 = v6115 >> (uint(v6116) % 32)
	v6119 = v6115 ^ v6117 - v6117
	v6122 = v6041 - v5926&v6113
	v6124 = v6122 >> (uint(v6116) % 32)
	v6126 = v6122 ^ v6124 - v6124
	if base.Ui32(v6126) < base.Ui32(v6119) {
		goto L422
	} else {
		goto L423
	}
L422:
	;
	v6128 = v6119
	goto L424
L423:
	;
	v6128 = v6126
	goto L424
L424:
	;
	if base.Ui32(v6128) < base.Ui32(v6110) {
		goto L425
	} else {
		goto L426
	}
L425:
	;
	v6130 = v6110
	goto L427
L426:
	;
	v6130 = v6128
	goto L427
L427:
	;
	v6133 = v6004 - int32(base.Ui32(v5999)>>(uint(int32(24))%32))
	v6134 = int32(31)
	v6135 = v6133 >> (uint(v6134) % 32)
	v6137 = v6133 ^ v6135 - v6135
	v6142 = v6015 - int32(base.Ui32(v5999)>>(uint(int32(16))%32))&int32(255)
	v6144 = v6142 >> (uint(v6134) % 32)
	v6146 = v6142 ^ v6144 - v6144
	if base.Ui32(v6146) < base.Ui32(v6137) {
		goto L428
	} else {
		goto L429
	}
L428:
	;
	v6148 = v6137
	goto L430
L429:
	;
	v6148 = v6146
	goto L430
L430:
	;
	v6151 = int32(255)
	v6153 = v6030 - int32(base.Ui32(v5999)>>(uint(int32(8))%32))&v6151
	v6154 = int32(31)
	v6155 = v6153 >> (uint(v6154) % 32)
	v6157 = v6153 ^ v6155 - v6155
	v6160 = v6041 - v5999&v6151
	v6162 = v6160 >> (uint(v6154) % 32)
	v6164 = v6160 ^ v6162 - v6162
	if base.Ui32(v6164) < base.Ui32(v6157) {
		goto L431
	} else {
		goto L432
	}
L431:
	;
	v6166 = v6157
	goto L433
L432:
	;
	v6166 = v6164
	goto L433
L433:
	;
	if base.Ui32(v6166) < base.Ui32(v6148) {
		goto L434
	} else {
		goto L435
	}
L434:
	;
	v6168 = v6148
	goto L436
L435:
	;
	v6168 = v6166
	goto L436
L436:
	;
	if base.Ui32(v6168) < base.Ui32(v6130) {
		goto L437
	} else {
		goto L438
	}
L437:
	;
	v6170 = v6130
	goto L439
L438:
	;
	v6170 = v6168
	goto L439
L439:
	;
	if base.Ui32(v6170) < base.Ui32(v6092) {
		goto L440
	} else {
		goto L441
	}
L440:
	;
	v6172 = v6092
	goto L442
L441:
	;
	v6172 = v6170
	goto L442
L442:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v5934))) = uint8(v6172)
	v6177 = v5935 + int32(-1)
	if v6177 != 0 {
		__phi5925 = v5999
		__phi5926 = v5925
		__phi5933 = v5946
		__phi5934 = v5934 + int32(1)
		__phi5935 = v6177
		v5925 = __phi5925
		v5926 = __phi5926
		v5933 = __phi5933
		v5934 = __phi5934
		v5935 = __phi5935
		goto L394
	} else {
		goto L443
	}
L443:
	;
	goto L395
L444:
	;
	v6307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5508+int32(base.Ui32(v5789)>>(uint(v5504)%32))*v5521<<(uint(int32(2))%32)+v6214>>(uint(v5504)%32)<<(uint(int32(2))%32))+1)))
	v6308 = v6214 + v5520
	v6309 = base.B2i32(v6308 < v5501)
	if v6308 < v5501 {
		goto L446
	} else {
		goto L447
	}
L445:
	;
	if v5858 != v5502 {
		__phi5762 = v5778
		__phi5764 = v6198
		__phi5778 = v5863
		__phi5788 = v6199
		__phi5789 = v5858
		v5762 = __phi5762
		v5764 = __phi5764
		v5778 = __phi5778
		v5788 = __phi5788
		v5789 = __phi5789
		goto L382
	} else {
		goto L450
	}
L446:
	;
	v6310 = v6308
	goto L448
L447:
	;
	v6310 = v5501
	goto L448
L448:
	;
	F_GetResidual(m, v5501, v5502, v5778, v5863, v6199, v6307, v6214, v6310, v5789, v5517, v5510, v5511, v5856+v6214<<(uint(int32(2))%32))
	mBase = m.M
	if v6308 < v5501 {
		v6214 = v6310
		goto L444
	} else {
		goto L449
	}
L449:
	;
	goto L445
L450:
	;
	goto L372
L451:
	;
	v6756 = int32(1)
	if v5502 == v6756 {
		v6890 = int32(0)
		goto L513
	} else {
		goto L514
	}
L452:
	;
	v6319 = v5501 << (uint(int32(2)) % 32)
	__phi6326 = v5506
	__phi6327 = v5507
	__phi6329 = v5752
	__phi6343 = v5604
	__phi6344 = int32(0)
	__phi6353 = v5751
	v6326 = __phi6326
	v6327 = __phi6327
	v6329 = __phi6329
	v6343 = __phi6343
	v6344 = __phi6344
	v6353 = __phi6353
	goto L453
L453:
	;
	v6419 = v6344 + int32(1)
	v6422 = int32(2)
	v6424 = F_memcpy(m, v6327, v6326, (base.B2i32(v6419 < v5502)+v5501)<<(uint(v6422)%32))
	mBase = m.M
	if v5502 <= v6344+v6422 {
		goto L455
	} else {
		goto L456
	}
L455:
	;
	if v5502 != v6419 {
		__phi6326 = v6326 + v6319
		__phi6327 = v6343
		__phi6329 = v6353
		__phi6343 = v6424
		__phi6344 = v6419
		__phi6353 = v6329
		v6326 = __phi6326
		v6327 = __phi6327
		v6329 = __phi6329
		v6343 = __phi6343
		v6344 = __phi6344
		v6353 = __phi6353
		goto L453
	} else {
		goto L512
	}
L456:
	;
	v6428 = v6326 + v6319
	if v5501 < int32(3) {
		goto L458
	} else {
		goto L459
	}
L457:
	;
	goto L455
L458:
	;
	goto L457
L459:
	;
	v6446 = int32(0)
	v6447 = base.Simd_g_v128_load64_zero(m, v6428, v6446)
	if v5511 == v6446 {
		v6464 = v6447
		goto L460
	} else {
		goto L461
	}
L460:
	;
	v6469 = int32(1)
	v6473 = int32(0)
	v6475 = v5501 << (uint(int32(2)) % 32)
	__phi6481 = base.Simd_g_i32x4_extract_lane_l1(v6464)
	__phi6482 = base.Simd_g_i32x4_extract_lane_l0(v6464)
	__phi6489 = v6428 + int32(4)
	__phi6490 = v6353 + v6469
	__phi6491 = v5501 + int32(-2)
	v6481 = __phi6481
	v6482 = __phi6482
	v6489 = __phi6489
	v6490 = __phi6490
	v6491 = __phi6491
	goto L462
L461:
	;
	v6451 = base.Simd_g_i32x4_shr_u(v6447, int32(8))
	v6454 = base.Simd_g_const(&F_VP8LResidualImage__k6)
	v6464 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v6451, base.Simd_g_const(&F_VP8LResidualImage__k7)), base.Simd_g_v128_and(v6447, v6454)), base.Simd_g_i32x4_shl(v6451, int32(16))), v6454), base.Simd_g_v128_and(v6447, base.Simd_g_const(&F_VP8LResidualImage__k8)))
	goto L460
L462:
	;
	v6502 = v6489 + int32(4)
	v6503 = *(*int32)(unsafe.Add(mBase, uint32(v6502)))
	v6505 = *(*int32)(unsafe.Add(mBase, uint32(v6489+v6475)))
	v6507 = *(*int32)(unsafe.Add(mBase, uint32(v6489+(v6473-v6475))))
	if v5511 == int32(0) {
		v6555 = v6503
		v6556 = v6505
		v6557 = v6507
		goto L464
	} else {
		goto L465
	}
L463:
	;
	goto L458
L464:
	;
	v6559 = int32(24)
	v6560 = int32(base.Ui32(v6481) >> (uint(v6559) % 32))
	v6563 = v6560 - int32(base.Ui32(v6557)>>(uint(v6559)%32))
	v6564 = int32(31)
	v6565 = v6563 >> (uint(v6564) % 32)
	v6567 = v6563 ^ v6565 - v6565
	v6568 = int32(16)
	v6570 = int32(255)
	v6571 = int32(base.Ui32(v6481)>>(uint(v6568)%32)) & v6570
	v6576 = v6571 - int32(base.Ui32(v6557)>>(uint(v6568)%32))&v6570
	v6578 = v6576 >> (uint(v6564) % 32)
	v6580 = v6576 ^ v6578 - v6578
	if base.Ui32(v6580) < base.Ui32(v6567) {
		goto L466
	} else {
		goto L467
	}
L465:
	;
	v6510 = int32(8)
	v6511 = int32(base.Ui32(v6503) >> (uint(v6510) % 32))
	v6512 = int32(255)
	v6514 = int32(16711935)
	v6517 = int32(16)
	v6522 = int32(-16711936)
	v6526 = int32(base.Ui32(v6505) >> (uint(v6510) % 32))
	v6541 = int32(base.Ui32(v6507) >> (uint(v6510) % 32))
	v6555 = (v6511&v6512+v6503&v6514+v6511<<(uint(v6517)%32))&v6514 | v6503&v6522
	v6556 = (v6526&v6512+v6505&v6514+v6526<<(uint(v6517)%32))&v6514 | v6505&v6522
	v6557 = (v6541&v6512+v6507&v6514+v6541<<(uint(v6517)%32))&v6514 | v6507&v6522
	goto L464
L466:
	;
	v6582 = v6567
	goto L468
L467:
	;
	v6582 = v6580
	goto L468
L468:
	;
	v6583 = int32(8)
	v6585 = int32(255)
	v6586 = int32(base.Ui32(v6481)>>(uint(v6583)%32)) & v6585
	v6591 = v6586 - int32(base.Ui32(v6557)>>(uint(v6583)%32))&v6585
	v6592 = int32(31)
	v6593 = v6591 >> (uint(v6592) % 32)
	v6595 = v6591 ^ v6593 - v6593
	v6597 = v6481 & v6585
	v6600 = v6597 - v6557&v6585
	v6602 = v6600 >> (uint(v6592) % 32)
	v6604 = v6600 ^ v6602 - v6602
	if base.Ui32(v6604) < base.Ui32(v6595) {
		goto L469
	} else {
		goto L470
	}
L469:
	;
	v6606 = v6595
	goto L471
L470:
	;
	v6606 = v6604
	goto L471
L471:
	;
	if base.Ui32(v6606) < base.Ui32(v6582) {
		goto L472
	} else {
		goto L473
	}
L472:
	;
	v6608 = v6582
	goto L474
L473:
	;
	v6608 = v6606
	goto L474
L474:
	;
	v6611 = v6560 - int32(base.Ui32(v6556)>>(uint(int32(24))%32))
	v6612 = int32(31)
	v6613 = v6611 >> (uint(v6612) % 32)
	v6615 = v6611 ^ v6613 - v6613
	v6620 = v6571 - int32(base.Ui32(v6556)>>(uint(int32(16))%32))&int32(255)
	v6622 = v6620 >> (uint(v6612) % 32)
	v6624 = v6620 ^ v6622 - v6622
	if base.Ui32(v6624) < base.Ui32(v6615) {
		goto L475
	} else {
		goto L476
	}
L475:
	;
	v6626 = v6615
	goto L477
L476:
	;
	v6626 = v6624
	goto L477
L477:
	;
	v6629 = int32(255)
	v6631 = v6586 - int32(base.Ui32(v6556)>>(uint(int32(8))%32))&v6629
	v6632 = int32(31)
	v6633 = v6631 >> (uint(v6632) % 32)
	v6635 = v6631 ^ v6633 - v6633
	v6638 = v6597 - v6556&v6629
	v6640 = v6638 >> (uint(v6632) % 32)
	v6642 = v6638 ^ v6640 - v6640
	if base.Ui32(v6642) < base.Ui32(v6635) {
		goto L478
	} else {
		goto L479
	}
L478:
	;
	v6644 = v6635
	goto L480
L479:
	;
	v6644 = v6642
	goto L480
L480:
	;
	if base.Ui32(v6644) < base.Ui32(v6626) {
		goto L481
	} else {
		goto L482
	}
L481:
	;
	v6646 = v6626
	goto L483
L482:
	;
	v6646 = v6644
	goto L483
L483:
	;
	if base.Ui32(v6646) < base.Ui32(v6608) {
		goto L484
	} else {
		goto L485
	}
L484:
	;
	v6648 = v6608
	goto L486
L485:
	;
	v6648 = v6646
	goto L486
L486:
	;
	v6651 = v6560 - int32(base.Ui32(v6482)>>(uint(int32(24))%32))
	v6652 = int32(31)
	v6653 = v6651 >> (uint(v6652) % 32)
	v6655 = v6651 ^ v6653 - v6653
	v6660 = v6571 - int32(base.Ui32(v6482)>>(uint(int32(16))%32))&int32(255)
	v6662 = v6660 >> (uint(v6652) % 32)
	v6664 = v6660 ^ v6662 - v6662
	if base.Ui32(v6664) < base.Ui32(v6655) {
		goto L487
	} else {
		goto L488
	}
L487:
	;
	v6666 = v6655
	goto L489
L488:
	;
	v6666 = v6664
	goto L489
L489:
	;
	v6669 = int32(255)
	v6671 = v6586 - int32(base.Ui32(v6482)>>(uint(int32(8))%32))&v6669
	v6672 = int32(31)
	v6673 = v6671 >> (uint(v6672) % 32)
	v6675 = v6671 ^ v6673 - v6673
	v6678 = v6597 - v6482&v6669
	v6680 = v6678 >> (uint(v6672) % 32)
	v6682 = v6678 ^ v6680 - v6680
	if base.Ui32(v6682) < base.Ui32(v6675) {
		goto L490
	} else {
		goto L491
	}
L490:
	;
	v6684 = v6675
	goto L492
L491:
	;
	v6684 = v6682
	goto L492
L492:
	;
	if base.Ui32(v6684) < base.Ui32(v6666) {
		goto L493
	} else {
		goto L494
	}
L493:
	;
	v6686 = v6666
	goto L495
L494:
	;
	v6686 = v6684
	goto L495
L495:
	;
	v6689 = v6560 - int32(base.Ui32(v6555)>>(uint(int32(24))%32))
	v6690 = int32(31)
	v6691 = v6689 >> (uint(v6690) % 32)
	v6693 = v6689 ^ v6691 - v6691
	v6698 = v6571 - int32(base.Ui32(v6555)>>(uint(int32(16))%32))&int32(255)
	v6700 = v6698 >> (uint(v6690) % 32)
	v6702 = v6698 ^ v6700 - v6700
	if base.Ui32(v6702) < base.Ui32(v6693) {
		goto L496
	} else {
		goto L497
	}
L496:
	;
	v6704 = v6693
	goto L498
L497:
	;
	v6704 = v6702
	goto L498
L498:
	;
	v6707 = int32(255)
	v6709 = v6586 - int32(base.Ui32(v6555)>>(uint(int32(8))%32))&v6707
	v6710 = int32(31)
	v6711 = v6709 >> (uint(v6710) % 32)
	v6713 = v6709 ^ v6711 - v6711
	v6716 = v6597 - v6555&v6707
	v6718 = v6716 >> (uint(v6710) % 32)
	v6720 = v6716 ^ v6718 - v6718
	if base.Ui32(v6720) < base.Ui32(v6713) {
		goto L499
	} else {
		goto L500
	}
L499:
	;
	v6722 = v6713
	goto L501
L500:
	;
	v6722 = v6720
	goto L501
L501:
	;
	if base.Ui32(v6722) < base.Ui32(v6704) {
		goto L502
	} else {
		goto L503
	}
L502:
	;
	v6724 = v6704
	goto L504
L503:
	;
	v6724 = v6722
	goto L504
L504:
	;
	if base.Ui32(v6724) < base.Ui32(v6686) {
		goto L505
	} else {
		goto L506
	}
L505:
	;
	v6726 = v6686
	goto L507
L506:
	;
	v6726 = v6724
	goto L507
L507:
	;
	if base.Ui32(v6726) < base.Ui32(v6648) {
		goto L508
	} else {
		goto L509
	}
L508:
	;
	v6728 = v6648
	goto L510
L509:
	;
	v6728 = v6726
	goto L510
L510:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v6490))) = uint8(v6728)
	v6733 = v6491 + int32(-1)
	if v6733 != 0 {
		__phi6481 = v6555
		__phi6482 = v6481
		__phi6489 = v6502
		__phi6490 = v6490 + int32(1)
		__phi6491 = v6733
		v6481 = __phi6481
		v6482 = __phi6482
		v6489 = __phi6489
		v6490 = __phi6490
		v6491 = __phi6491
		goto L462
	} else {
		goto L511
	}
L511:
	;
	goto L463
L512:
	;
	goto L372
L513:
	;
	if v5502&v6756 == int32(0) {
		goto L372
	} else {
		goto L518
	}
L514:
	;
	v6776 = int32(0)
	v6791 = v5506
	goto L515
L515:
	;
	v6869 = int32(2)
	v6871 = F_memcpy(m, v5507, v6791, (base.B2i32(v6776+int32(1) < v5502)+v5501)<<(uint(v6869)%32))
	mBase = m.M
	v6874 = v6776 + v6869
	v6879 = F_memcpy(m, v5604, v6791+v5501<<(uint(int32(2))%32), (base.B2i32(v6874 < v5502)+v5501)<<(uint(v6869)%32))
	mBase = m.M
	if v5502&int32(2147483646) != v6874 {
		v6776 = v6874
		v6791 = v6791 + v5501<<(uint(int32(3))%32)
		goto L515
	} else {
		goto L517
	}
L516:
	;
	v6890 = v6874
	goto L513
L517:
	;
	goto L516
L518:
	;
	v6982 = int32(2)
	v6991 = F_memcpy(m, v5507, v5506+v6890*v5501<<(uint(v6982)%32), (base.B2i32(v6890+int32(1) < v5502)+v5501)<<(uint(v6982)%32))
	mBase = m.M
	goto L372
}

var F_VP8LResidualImage__k0 = [2]uint64{0xff000b00ff000b00, 0xff000b00ff000b00}
var F_VP8LResidualImage__k1 = [2]uint64{0x100000000, 0x300000002}
var F_VP8LResidualImage__k2 = [2]uint64{0x0, 0x0}
var F_VP8LResidualImage__k3 = [2]uint64{0x400000004, 0x400000004}
var F_VP8LResidualImage__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_VP8LResidualImage__k5 = [2]uint64{0x302010007060504, 0x302010003020100}
var F_VP8LResidualImage__k6 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8LResidualImage__k7 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_VP8LResidualImage__k8 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}

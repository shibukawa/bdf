//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_UpsampleRgb565LinePair_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
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
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v203 int32
	_ = v203
	var v206 int32
	_ = v206
	var v209 int32
	_ = v209
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v225 int32
	_ = v225
	var v228 int32
	_ = v228
	var v231 int32
	_ = v231
	var v241 int32
	_ = v241
	var v243 int32
	_ = v243
	var v250 int32
	_ = v250
	var v253 int32
	_ = v253
	var v254 int32
	_ = v254
	var v258 int32
	_ = v258
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v265 int32
	_ = v265
	var v270 int32
	_ = v270
	var v272 int32
	_ = v272
	var v279 int32
	_ = v279
	var v282 int32
	_ = v282
	var v287 int32
	_ = v287
	var v291 int32
	_ = v291
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v306 int32
	_ = v306
	var v309 int32
	_ = v309
	var v312 int32
	_ = v312
	var v322 int32
	_ = v322
	var v324 int32
	_ = v324
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v342 int32
	_ = v342
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v351 int32
	_ = v351
	var v359 int32
	_ = v359
	var v371 int32
	_ = v371
	var v372 int32
	_ = v372
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v393 int32
	_ = v393
	var v394 int32
	_ = v394
	var v395 base.V128
	_ = v395
	var v396 int32
	_ = v396
	var v398 base.V128
	_ = v398
	var v399 int32
	_ = v399
	var v400 base.V128
	_ = v400
	var v401 base.V128
	_ = v401
	var v403 base.V128
	_ = v403
	var v404 base.V128
	_ = v404
	var v406 base.V128
	_ = v406
	var v407 base.V128
	_ = v407
	var v409 base.V128
	_ = v409
	var v411 base.V128
	_ = v411
	var v413 base.V128
	_ = v413
	var v419 base.V128
	_ = v419
	var v420 base.V128
	_ = v420
	var v426 base.V128
	_ = v426
	var v427 base.V128
	_ = v427
	var v428 base.V128
	_ = v428
	var v429 base.V128
	_ = v429
	var v432 base.V128
	_ = v432
	var v433 base.V128
	_ = v433
	var v436 base.V128
	_ = v436
	var v437 base.V128
	_ = v437
	var v439 base.V128
	_ = v439
	var v443 base.V128
	_ = v443
	var v446 int32
	_ = v446
	var v448 base.V128
	_ = v448
	var v449 int32
	_ = v449
	var v451 base.V128
	_ = v451
	var v453 base.V128
	_ = v453
	var v454 base.V128
	_ = v454
	var v456 base.V128
	_ = v456
	var v457 base.V128
	_ = v457
	var v459 base.V128
	_ = v459
	var v460 base.V128
	_ = v460
	var v462 base.V128
	_ = v462
	var v465 base.V128
	_ = v465
	var v471 base.V128
	_ = v471
	var v472 base.V128
	_ = v472
	var v478 base.V128
	_ = v478
	var v479 base.V128
	_ = v479
	var v481 base.V128
	_ = v481
	var v485 base.V128
	_ = v485
	var v488 base.V128
	_ = v488
	var v489 base.V128
	_ = v489
	var v491 base.V128
	_ = v491
	var v495 base.V128
	_ = v495
	var v512 int32
	_ = v512
	var v513 int32
	_ = v513
	var v522 base.V128
	_ = v522
	var v524 int32
	_ = v524
	var v525 base.V128
	_ = v525
	var v526 base.V128
	_ = v526
	var v527 base.V128
	_ = v527
	var v529 base.V128
	_ = v529
	var v533 base.V128
	_ = v533
	var v534 base.V128
	_ = v534
	var v537 base.V128
	_ = v537
	var v539 base.V128
	_ = v539
	var v540 base.V128
	_ = v540
	var v541 base.V128
	_ = v541
	var v543 base.V128
	_ = v543
	var v549 base.V128
	_ = v549
	var v551 base.V128
	_ = v551
	var v552 base.V128
	_ = v552
	var v553 base.V128
	_ = v553
	var v555 base.V128
	_ = v555
	var v563 int32
	_ = v563
	var v564 base.V128
	_ = v564
	var v565 base.V128
	_ = v565
	var v570 base.V128
	_ = v570
	var v579 base.V128
	_ = v579
	var v584 base.V128
	_ = v584
	var v593 base.V128
	_ = v593
	var v595 int32
	_ = v595
	var v605 base.V128
	_ = v605
	var v630 int32
	_ = v630
	var v631 int32
	_ = v631
	var v640 base.V128
	_ = v640
	var v642 int32
	_ = v642
	var v643 base.V128
	_ = v643
	var v644 base.V128
	_ = v644
	var v645 base.V128
	_ = v645
	var v647 base.V128
	_ = v647
	var v651 base.V128
	_ = v651
	var v652 base.V128
	_ = v652
	var v655 base.V128
	_ = v655
	var v657 base.V128
	_ = v657
	var v658 base.V128
	_ = v658
	var v659 base.V128
	_ = v659
	var v661 base.V128
	_ = v661
	var v667 base.V128
	_ = v667
	var v669 base.V128
	_ = v669
	var v670 base.V128
	_ = v670
	var v671 base.V128
	_ = v671
	var v673 base.V128
	_ = v673
	var v681 int32
	_ = v681
	var v682 base.V128
	_ = v682
	var v683 base.V128
	_ = v683
	var v688 base.V128
	_ = v688
	var v697 base.V128
	_ = v697
	var v702 base.V128
	_ = v702
	var v711 base.V128
	_ = v711
	var v713 int32
	_ = v713
	var v723 base.V128
	_ = v723
	var v732 int32
	_ = v732
	var v737 int32
	_ = v737
	var v755 int32
	_ = v755
	var v756 int32
	_ = v756
	var v779 int32
	_ = v779
	var v782 int32
	_ = v782
	var v788 int32
	_ = v788
	var v789 int32
	_ = v789
	var v791 int32
	_ = v791
	var v793 int32
	_ = v793
	var v794 int32
	_ = v794
	var v798 int32
	_ = v798
	var v799 int32
	_ = v799
	var v800 int32
	_ = v800
	var v802 int32
	_ = v802
	var v813 int32
	_ = v813
	var v838 int32
	_ = v838
	var v839 int32
	_ = v839
	var v843 int32
	_ = v843
	var v847 int32
	_ = v847
	var v848 int32
	_ = v848
	var v883 int32
	_ = v883
	var v884 int32
	_ = v884
	var v889 int64
	_ = v889
	var v892 int32
	_ = v892
	var v893 int32
	_ = v893
	var v905 int32
	_ = v905
	var v923 int32
	_ = v923
	var v924 int32
	_ = v924
	var v925 int32
	_ = v925
	var v936 int32
	_ = v936
	var v961 int32
	_ = v961
	var v962 int32
	_ = v962
	var v966 int32
	_ = v966
	var v970 int32
	_ = v970
	var v971 int32
	_ = v971
	var v1006 int32
	_ = v1006
	var v1007 int32
	_ = v1007
	var v1012 int64
	_ = v1012
	var v1015 int32
	_ = v1015
	var v1016 int32
	_ = v1016
	var v1028 int32
	_ = v1028
	var v1046 int32
	_ = v1046
	var v1047 base.V128
	_ = v1047
	var v1048 int32
	_ = v1048
	var v1049 base.V128
	_ = v1049
	var v1051 base.V128
	_ = v1051
	var v1052 base.V128
	_ = v1052
	var v1053 int32
	_ = v1053
	var v1054 base.V128
	_ = v1054
	var v1055 base.V128
	_ = v1055
	var v1057 base.V128
	_ = v1057
	var v1058 base.V128
	_ = v1058
	var v1060 base.V128
	_ = v1060
	var v1062 base.V128
	_ = v1062
	var v1064 base.V128
	_ = v1064
	var v1070 base.V128
	_ = v1070
	var v1071 base.V128
	_ = v1071
	var v1077 base.V128
	_ = v1077
	var v1078 base.V128
	_ = v1078
	var v1079 base.V128
	_ = v1079
	var v1080 base.V128
	_ = v1080
	var v1083 base.V128
	_ = v1083
	var v1084 base.V128
	_ = v1084
	var v1087 base.V128
	_ = v1087
	var v1088 base.V128
	_ = v1088
	var v1090 base.V128
	_ = v1090
	var v1094 base.V128
	_ = v1094
	var v1100 int32
	_ = v1100
	var v1102 int32
	_ = v1102
	var v1103 int32
	_ = v1103
	var v1114 int32
	_ = v1114
	var v1139 int32
	_ = v1139
	var v1140 int32
	_ = v1140
	var v1144 int32
	_ = v1144
	var v1148 int32
	_ = v1148
	var v1149 int32
	_ = v1149
	var v1184 int32
	_ = v1184
	var v1185 int32
	_ = v1185
	var v1190 int64
	_ = v1190
	var v1193 int32
	_ = v1193
	var v1194 int32
	_ = v1194
	var v1206 int32
	_ = v1206
	var v1224 int32
	_ = v1224
	var v1235 int32
	_ = v1235
	var v1260 int32
	_ = v1260
	var v1261 int32
	_ = v1261
	var v1265 int32
	_ = v1265
	var v1269 int32
	_ = v1269
	var v1270 int32
	_ = v1270
	var v1305 int32
	_ = v1305
	var v1306 int32
	_ = v1306
	var v1311 int64
	_ = v1311
	var v1314 int32
	_ = v1314
	var v1315 int32
	_ = v1315
	var v1327 int32
	_ = v1327
	var v1346 base.V128
	_ = v1346
	var v1347 int32
	_ = v1347
	var v1348 base.V128
	_ = v1348
	var v1350 base.V128
	_ = v1350
	var v1351 base.V128
	_ = v1351
	var v1353 base.V128
	_ = v1353
	var v1354 base.V128
	_ = v1354
	var v1356 base.V128
	_ = v1356
	var v1357 base.V128
	_ = v1357
	var v1359 base.V128
	_ = v1359
	var v1362 base.V128
	_ = v1362
	var v1368 base.V128
	_ = v1368
	var v1369 base.V128
	_ = v1369
	var v1375 base.V128
	_ = v1375
	var v1376 base.V128
	_ = v1376
	var v1377 base.V128
	_ = v1377
	var v1378 base.V128
	_ = v1378
	var v1381 base.V128
	_ = v1381
	var v1382 base.V128
	_ = v1382
	var v1385 base.V128
	_ = v1385
	var v1386 base.V128
	_ = v1386
	var v1388 base.V128
	_ = v1388
	var v1392 base.V128
	_ = v1392
	var v1398 int32
	_ = v1398
	var v1399 int32
	_ = v1399
	var v1401 int32
	_ = v1401
	var v1415 int32
	_ = v1415
	var v1416 int32
	_ = v1416
	var v1425 base.V128
	_ = v1425
	var v1427 int32
	_ = v1427
	var v1428 base.V128
	_ = v1428
	var v1429 base.V128
	_ = v1429
	var v1430 base.V128
	_ = v1430
	var v1432 base.V128
	_ = v1432
	var v1436 base.V128
	_ = v1436
	var v1437 base.V128
	_ = v1437
	var v1440 base.V128
	_ = v1440
	var v1442 base.V128
	_ = v1442
	var v1443 base.V128
	_ = v1443
	var v1444 base.V128
	_ = v1444
	var v1446 base.V128
	_ = v1446
	var v1452 base.V128
	_ = v1452
	var v1454 base.V128
	_ = v1454
	var v1455 base.V128
	_ = v1455
	var v1456 base.V128
	_ = v1456
	var v1458 base.V128
	_ = v1458
	var v1466 int32
	_ = v1466
	var v1467 base.V128
	_ = v1467
	var v1468 base.V128
	_ = v1468
	var v1473 base.V128
	_ = v1473
	var v1482 base.V128
	_ = v1482
	var v1487 base.V128
	_ = v1487
	var v1496 base.V128
	_ = v1496
	var v1498 int32
	_ = v1498
	var v1508 base.V128
	_ = v1508
	var v1517 int32
	_ = v1517
	var v1522 int32
	_ = v1522
	var v1526 int32
	_ = v1526
	var v1540 int32
	_ = v1540
	var v1541 int32
	_ = v1541
	var v1550 base.V128
	_ = v1550
	var v1552 int32
	_ = v1552
	var v1553 base.V128
	_ = v1553
	var v1554 base.V128
	_ = v1554
	var v1555 base.V128
	_ = v1555
	var v1557 base.V128
	_ = v1557
	var v1561 base.V128
	_ = v1561
	var v1562 base.V128
	_ = v1562
	var v1565 base.V128
	_ = v1565
	var v1567 base.V128
	_ = v1567
	var v1568 base.V128
	_ = v1568
	var v1569 base.V128
	_ = v1569
	var v1571 base.V128
	_ = v1571
	var v1577 base.V128
	_ = v1577
	var v1579 base.V128
	_ = v1579
	var v1580 base.V128
	_ = v1580
	var v1581 base.V128
	_ = v1581
	var v1583 base.V128
	_ = v1583
	var v1591 int32
	_ = v1591
	var v1592 base.V128
	_ = v1592
	var v1593 base.V128
	_ = v1593
	var v1598 base.V128
	_ = v1598
	var v1607 base.V128
	_ = v1607
	var v1612 base.V128
	_ = v1612
	var v1621 base.V128
	_ = v1621
	var v1623 int32
	_ = v1623
	var v1633 base.V128
	_ = v1633
	var v1647 int32
	_ = v1647
	var v1661 int32
	_ = v1661
	var v1662 int32
	_ = v1662
	var v1671 base.V128
	_ = v1671
	var v1673 int32
	_ = v1673
	var v1674 base.V128
	_ = v1674
	var v1675 base.V128
	_ = v1675
	var v1676 base.V128
	_ = v1676
	var v1678 base.V128
	_ = v1678
	var v1682 base.V128
	_ = v1682
	var v1683 base.V128
	_ = v1683
	var v1686 base.V128
	_ = v1686
	var v1688 base.V128
	_ = v1688
	var v1689 base.V128
	_ = v1689
	var v1690 base.V128
	_ = v1690
	var v1692 base.V128
	_ = v1692
	var v1698 base.V128
	_ = v1698
	var v1700 base.V128
	_ = v1700
	var v1701 base.V128
	_ = v1701
	var v1702 base.V128
	_ = v1702
	var v1704 base.V128
	_ = v1704
	var v1712 int32
	_ = v1712
	var v1713 base.V128
	_ = v1713
	var v1714 base.V128
	_ = v1714
	var v1719 base.V128
	_ = v1719
	var v1728 base.V128
	_ = v1728
	var v1733 base.V128
	_ = v1733
	var v1742 base.V128
	_ = v1742
	var v1744 int32
	_ = v1744
	var v1754 base.V128
	_ = v1754
	var v1763 int32
	_ = v1763
	var v1764 int32
	_ = v1764
	var v1767 int32
	_ = v1767
	var v1768 int32
	_ = v1768
	var v1770 int32
	_ = v1770
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
	v181 = int32(base.Ui32(v177*int32(_a_F_UpsampleRgb565LinePair_SSE2_0)) >> (uint(v175) % 32))
	v182 = int32(base.Ui32(v172*int32(_a_F_UpsampleRgb565LinePair_SSE2_1))>>(uint(v175)%32)) + v181
	v184 = v182 + int32(-14234)
	if base.Ui32(v182) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_2)) {
		v191 = int32(0)
	} else {
		v191 = int32(248)
	}
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_3)) {
		v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
	} else {
		v194 = v191
	}
	v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v198 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v200 = int32(1)
	v203 = int32(base.Ui32(v197+v198)>>(uint(v200)%32)) + v200
	v206 = int32(base.Ui32(v203+v198) >> (uint(v200) % 32))
	v209 = int32(8)
	v216 = v181 - (int32(base.Ui32(v206*int32(_a_F_UpsampleRgb565LinePair_SSE2_4))>>(uint(v209)%32)) + int32(base.Ui32(v172*int32(_a_F_UpsampleRgb565LinePair_SSE2_5))>>(uint(v209)%32)))
	v218 = v216 + int32(_a_F_UpsampleRgb565LinePair_SSE2_6)
	if v216 < int32(-8708) {
		v225 = int32(0)
	} else {
		v225 = int32(255)
	}
	if base.Ui32(v218) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_3)) {
		v228 = int32(base.Ui32(v218) >> (uint(int32(6)) % 32))
	} else {
		v228 = v225
	}
	v231 = v194&int32(248) | int32(base.Ui32(v228)>>(uint(int32(5))%32))
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v231)
	v241 = int32(base.Ui32(v206*int32(_a_F_UpsampleRgb565LinePair_SSE2_7))>>(uint(int32(8))%32)) + v181
	v243 = v241 + int32(-17685)
	if base.Ui32(v241) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_8)) {
		v250 = int32(0)
	} else {
		v250 = int32(31)
	}
	if base.Ui32(v243) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_3)) {
		v253 = int32(base.Ui32(v243) >> (uint(int32(9)) % 32))
	} else {
		v253 = v250
	}
	v254 = v228<<(uint(int32(3))%32)&int32(224) | v253
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v254)
	if l1 == int32(0) {
	} else {
		v258 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v261 = int32(8)
		v262 = int32(base.Ui32(v258*int32(_a_F_UpsampleRgb565LinePair_SSE2_0)) >> (uint(v261) % 32))
		v265 = int32(base.Ui32(v169+v163) >> (uint(int32(1)) % 32))
		v270 = v262 + int32(base.Ui32(v265*int32(_a_F_UpsampleRgb565LinePair_SSE2_1))>>(uint(v261)%32))
		v272 = v270 + int32(-14234)
		if base.Ui32(v270) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_2)) {
			v279 = int32(0)
		} else {
			v279 = int32(248)
		}
		if base.Ui32(v272) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_3)) {
			v282 = int32(base.Ui32(v272) >> (uint(int32(6)) % 32))
		} else {
			v282 = v279
		}
		v287 = int32(8)
		v291 = int32(base.Ui32(v203+v197) >> (uint(int32(1)) % 32))
		v297 = v262 - (int32(base.Ui32(v265*int32(_a_F_UpsampleRgb565LinePair_SSE2_5))>>(uint(v287)%32)) + int32(base.Ui32(v291*int32(_a_F_UpsampleRgb565LinePair_SSE2_4))>>(uint(v287)%32)))
		v299 = v297 + int32(_a_F_UpsampleRgb565LinePair_SSE2_6)
		if v297 < int32(-8708) {
			v306 = int32(0)
		} else {
			v306 = int32(255)
		}
		if base.Ui32(v299) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_3)) {
			v309 = int32(base.Ui32(v299) >> (uint(int32(6)) % 32))
		} else {
			v309 = v306
		}
		v312 = v282&int32(248) | int32(base.Ui32(v309)>>(uint(int32(5))%32))
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v312)
		v322 = v262 + int32(base.Ui32(v291*int32(_a_F_UpsampleRgb565LinePair_SSE2_7))>>(uint(int32(8))%32))
		v324 = v322 + int32(-17685)
		if base.Ui32(v322) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_8)) {
			v331 = int32(0)
		} else {
			v331 = int32(31)
		}
		if base.Ui32(v324) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_SSE2_3)) {
			v334 = int32(base.Ui32(v324) >> (uint(int32(9)) % 32))
		} else {
			v334 = v331
		}
		v335 = v309<<(uint(int32(3))%32)&int32(224) | v334
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v335)
	}
	v342 = v35 + int32(96)
	v344 = v35 + int32(64)
	if l8 < int32(34) {
		v755 = v10
		v756 = v162
	} else {
		v347 = int32(2)
		v351 = int32(1)
		v359 = int32(0)
		v371 = v359
		v372 = l7 + v347
		v376 = v359
		v377 = l6 + v347
		for {
			v393 = l4 + v371
			v394 = int32(0)
			v395 = base.Simd_g_v128_load_rng(m, v393, v394, int32(0), int32(17))
			v396 = l2 + v371
			v398 = base.Simd_g_v128_load_rng(m, v396, v394, int32(0), int32(17))
			v399 = int32(1)
			v400 = base.Simd_g_v128_load_nc(m, v393, v399)
			v401 = base.Simd_g_i8x16_avgr_u(v398, v400)
			v403 = base.Simd_g_v128_load_nc(m, v396, v399)
			v404 = base.Simd_g_i8x16_avgr_u(v395, v403)
			v406 = base.Simd_g_v128_xor(v400, v398)
			v407 = base.Simd_g_v128_xor(v395, v403)
			v409 = base.Simd_g_v128_xor(v401, v404)
			v411 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k0)
			v413 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v401, v404), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v406, v407), v409), v411))
			v419 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v413, v401), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v413, v401), base.Simd_g_v128_and(v409, v406)), v411))
			v420 = base.Simd_g_i8x16_avgr_u(v395, v419)
			v426 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v413, v404), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v413, v404), base.Simd_g_v128_and(v409, v407)), v411))
			v427 = base.Simd_g_i8x16_avgr_u(v400, v426)
			v428 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k1)
			v429 = base.Simd_g_i8x16_shuffle2(v420, v427, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v344, int32(80), v429)
			v432 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
			v433 = base.Simd_g_i8x16_shuffle2(v420, v427, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v344, int32(64), v433)
			v436 = base.Simd_g_i8x16_avgr_u(v398, v426)
			v437 = base.Simd_g_i8x16_avgr_u(v403, v419)
			v439 = base.Simd_g_i8x16_shuffle2(v436, v437, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v344, int32(16), v439)
			v443 = base.Simd_g_i8x16_shuffle2(v436, v437, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v344, v394, v443)
			v446 = l5 + v371
			v448 = base.Simd_g_v128_load_rng(m, v446, v394, int32(0), int32(17))
			v449 = l3 + v371
			v451 = base.Simd_g_v128_load_rng(m, v449, v394, int32(0), int32(17))
			v453 = base.Simd_g_v128_load_nc(m, v446, v399)
			v454 = base.Simd_g_i8x16_avgr_u(v451, v453)
			v456 = base.Simd_g_v128_load_nc(m, v449, v399)
			v457 = base.Simd_g_i8x16_avgr_u(v448, v456)
			v459 = base.Simd_g_v128_xor(v453, v451)
			v460 = base.Simd_g_v128_xor(v448, v456)
			v462 = base.Simd_g_v128_xor(v454, v457)
			v465 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v454, v457), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v459, v460), v462), v411))
			v471 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v465, v454), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v465, v454), base.Simd_g_v128_and(v462, v459)), v411))
			v472 = base.Simd_g_i8x16_avgr_u(v448, v471)
			v478 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v465, v457), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v465, v457), base.Simd_g_v128_and(v462, v460)), v411))
			v479 = base.Simd_g_i8x16_avgr_u(v453, v478)
			v481 = base.Simd_g_i8x16_shuffle2(v472, v479, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v344, int32(112), v481)
			v485 = base.Simd_g_i8x16_shuffle2(v472, v479, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v344, int32(96), v485)
			v488 = base.Simd_g_i8x16_avgr_u(v451, v478)
			v489 = base.Simd_g_i8x16_avgr_u(v456, v471)
			v491 = base.Simd_g_i8x16_shuffle2(v488, v489, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v344, int32(48), v491)
			v495 = base.Simd_g_i8x16_shuffle2(v488, v489, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v344, int32(32), v495)
			v512 = v377
			v513 = v394
			for {
				v522 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k7)
				v524 = int32(0)
				v525 = base.Simd_g_v128_load64_zero(m, l0+v351+v376+v513, v524)
				v526 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
				v527 = base.Simd_g_i8x16_shuffle2(v522, v525, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v529 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k8)
				v533 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k9)
				v534 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v527), v529), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v527), v529), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))
				v537 = base.Simd_g_v128_load64_zero(m, v344+v513, v524)
				v539 = base.Simd_g_i8x16_shuffle2(v522, v537, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v540 = base.Simd_g_i32x4_extend_low_i16x8_u(v539)
				v541 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k12)
				v543 = base.Simd_g_i32x4_extend_high_i16x8_u(v539)
				v549 = base.Simd_g_v128_load64_zero(m, v342+v513, v524)
				v551 = base.Simd_g_i8x16_shuffle2(v522, v549, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v552 = base.Simd_g_i32x4_extend_low_i16x8_u(v551)
				v553 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k13)
				v555 = base.Simd_g_i32x4_extend_high_i16x8_u(v551)
				v563 = int32(6)
				v564 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v534, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v540, v541), base.Simd_g_i32x4_mul(v543, v541), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v552, v553), base.Simd_g_i32x4_mul(v555, v553), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k14)), v563)
				v565 = base.Simd_g_i8x16_narrow_i16x8_u(v564, v564)
				v570 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k15)
				v579 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v534, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v552, v570), base.Simd_g_i32x4_mul(v555, v570), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k16)), v563)
				v584 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k17)
				v593 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v540, v584), base.Simd_g_i32x4_mul(v543, v584), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), v534), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k18)), v563)
				v595 = int32(3)
				v605 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(v565, int32(5)), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k19)), base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(v579, v579), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k20))), base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_narrow_i16x8_u(v593, v593), v595), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k21)), base.Simd_g_v128_and(base.Simd_g_i16x8_shl(v565, v595), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k22))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				base.Simd_g_v128_store(m, v512, v524, v605)
				if base.Ui32(v513) < base.Ui32(int32(24)) {
					v512 = v512 + int32(16)
					v513 = v513 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			if l1 == int32(0) {
			} else {
				v630 = v372
				v631 = int32(0)
				for {
					v640 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k7)
					v642 = int32(0)
					v643 = base.Simd_g_v128_load64_zero(m, l1+v351+v376+v631, v642)
					v644 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
					v645 = base.Simd_g_i8x16_shuffle2(v640, v643, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
					v647 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k8)
					v651 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k9)
					v652 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v645), v647), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v645), v647), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))
					v655 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v631, v642)
					v657 = base.Simd_g_i8x16_shuffle2(v640, v655, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
					v658 = base.Simd_g_i32x4_extend_low_i16x8_u(v657)
					v659 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k12)
					v661 = base.Simd_g_i32x4_extend_high_i16x8_u(v657)
					v667 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v631, v642)
					v669 = base.Simd_g_i8x16_shuffle2(v640, v667, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
					v670 = base.Simd_g_i32x4_extend_low_i16x8_u(v669)
					v671 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k13)
					v673 = base.Simd_g_i32x4_extend_high_i16x8_u(v669)
					v681 = int32(6)
					v682 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v652, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v658, v659), base.Simd_g_i32x4_mul(v661, v659), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v670, v671), base.Simd_g_i32x4_mul(v673, v671), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k14)), v681)
					v683 = base.Simd_g_i8x16_narrow_i16x8_u(v682, v682)
					v688 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k15)
					v697 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v652, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v670, v688), base.Simd_g_i32x4_mul(v673, v688), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k16)), v681)
					v702 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k17)
					v711 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v658, v702), base.Simd_g_i32x4_mul(v661, v702), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), v652), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k18)), v681)
					v713 = int32(3)
					v723 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(v683, int32(5)), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k19)), base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(v697, v697), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k20))), base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_narrow_i16x8_u(v711, v711), v713), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k21)), base.Simd_g_v128_and(base.Simd_g_i16x8_shl(v683, v713), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k22))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
					base.Simd_g_v128_store(m, v630, v642, v723)
					if base.Ui32(v631) < base.Ui32(int32(24)) {
						v630 = v630 + int32(16)
						v631 = v631 + int32(8)
						continue
					} else {
						break
					}
					break
				}
			}
			v732 = int32(64)
			v737 = v371 + int32(16)
			if v376+int32(66) <= l8 {
				v371 = v737
				v372 = v372 + v732
				v376 = v376 + int32(32)
				v377 = v377 + v732
				continue
			} else {
				break
			}
			break
		}
		v755 = v737
		v756 = v376 + int32(33)
	}
	if l8 < int32(2) {
	} else {
		v779 = int32(32)
		v782 = int32(1)
		v788 = int32(base.Ui32(l8+v782)>>(uint(v782)%32)) - int32(base.Ui32(v756)>>(uint(v782)%32))
		v789 = F_memcpy(m, v35+v779, l2+v755, v788)
		mBase = m.M
		v791 = F_memcpy(m, v35, l4+v755, v788)
		mBase = m.M
		v793 = v791 + v779
		v794 = v793 + v788
		v798 = v788 + int32(-1)
		v799 = v793 + v798
		v800 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v799))))
		v802 = int32(17) - v788
		if base.Ui32(v802) < base.Ui32(int32(33)) {
			if v802 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v794))) = uint8(v800)
				v813 = v794 + v802
				*(*uint8)(unsafe.Add(mBase, uint32(v813+int32(-1)))) = uint8(v800)
				if base.Ui32(v802) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v794)+2)) = uint8(v800)
					*(*uint8)(unsafe.Add(mBase, uint32(v794)+1)) = uint8(v800)
					*(*uint8)(unsafe.Add(mBase, uint32(v813+int32(-3)))) = uint8(v800)
					*(*uint8)(unsafe.Add(mBase, uint32(v813+int32(-2)))) = uint8(v800)
					if base.Ui32(v802) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v794)+3)) = uint8(v800)
						*(*uint8)(unsafe.Add(mBase, uint32(v813+int32(-4)))) = uint8(v800)
						if base.Ui32(v802) < base.Ui32(int32(9)) {
						} else {
							v838 = (int32(0) - v794) & int32(3)
							v839 = v794 + v838
							v843 = v800 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v839))) = v843
							v847 = (v802 - v838) & int32(60)
							v848 = v839 + v847
							*(*int32)(unsafe.Add(mBase, uint32(v848+int32(-4)))) = v843
							if base.Ui32(v847) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v839)+8)) = v843
								*(*int32)(unsafe.Add(mBase, uint32(v839)+4)) = v843
								*(*int32)(unsafe.Add(mBase, uint32(v848+int32(-8)))) = v843
								*(*int32)(unsafe.Add(mBase, uint32(v848+int32(-12)))) = v843
								if base.Ui32(v847) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v839)+24)) = v843
									*(*int32)(unsafe.Add(mBase, uint32(v839)+20)) = v843
									*(*int32)(unsafe.Add(mBase, uint32(v839)+16)) = v843
									*(*int32)(unsafe.Add(mBase, uint32(v839)+12)) = v843
									*(*int32)(unsafe.Add(mBase, uint32(v848+int32(-16)))) = v843
									*(*int32)(unsafe.Add(mBase, uint32(v848+int32(-20)))) = v843
									*(*int32)(unsafe.Add(mBase, uint32(v848+int32(-24)))) = v843
									*(*int32)(unsafe.Add(mBase, uint32(v848+int32(-28)))) = v843
									v883 = v839&int32(4) | int32(24)
									v884 = v847 - v883
									if base.Ui32(v884) < base.Ui32(int32(32)) {
									} else {
										v889 = base.I64_extend_i32_u(v843) * int64(4294967297)
										v892 = v884
										v893 = v839 + v883
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v893)+24)) = v889
											*(*int64)(unsafe.Add(mBase, uint32(v893)+16)) = v889
											*(*int64)(unsafe.Add(mBase, uint32(v893)+8)) = v889
											*(*int64)(unsafe.Add(mBase, uint32(v893))) = v889
											v905 = v892 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v905) {
												v892 = v905
												v893 = v893 + int32(32)
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
			base.MemoryFill(m, v794, v800, v802)
		}
		v923 = v791 + v788
		v924 = v791 + v798
		v925 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v924))))
		if base.Ui32(v802) < base.Ui32(int32(33)) {
			if v802 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v923))) = uint8(v925)
				v936 = v923 + v802
				*(*uint8)(unsafe.Add(mBase, uint32(v936+int32(-1)))) = uint8(v925)
				if base.Ui32(v802) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v923)+2)) = uint8(v925)
					*(*uint8)(unsafe.Add(mBase, uint32(v923)+1)) = uint8(v925)
					*(*uint8)(unsafe.Add(mBase, uint32(v936+int32(-3)))) = uint8(v925)
					*(*uint8)(unsafe.Add(mBase, uint32(v936+int32(-2)))) = uint8(v925)
					if base.Ui32(v802) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v923)+3)) = uint8(v925)
						*(*uint8)(unsafe.Add(mBase, uint32(v936+int32(-4)))) = uint8(v925)
						if base.Ui32(v802) < base.Ui32(int32(9)) {
						} else {
							v961 = (int32(0) - v923) & int32(3)
							v962 = v923 + v961
							v966 = v925 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v962))) = v966
							v970 = (v802 - v961) & int32(60)
							v971 = v962 + v970
							*(*int32)(unsafe.Add(mBase, uint32(v971+int32(-4)))) = v966
							if base.Ui32(v970) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v962)+8)) = v966
								*(*int32)(unsafe.Add(mBase, uint32(v962)+4)) = v966
								*(*int32)(unsafe.Add(mBase, uint32(v971+int32(-8)))) = v966
								*(*int32)(unsafe.Add(mBase, uint32(v971+int32(-12)))) = v966
								if base.Ui32(v970) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v962)+24)) = v966
									*(*int32)(unsafe.Add(mBase, uint32(v962)+20)) = v966
									*(*int32)(unsafe.Add(mBase, uint32(v962)+16)) = v966
									*(*int32)(unsafe.Add(mBase, uint32(v962)+12)) = v966
									*(*int32)(unsafe.Add(mBase, uint32(v971+int32(-16)))) = v966
									*(*int32)(unsafe.Add(mBase, uint32(v971+int32(-20)))) = v966
									*(*int32)(unsafe.Add(mBase, uint32(v971+int32(-24)))) = v966
									*(*int32)(unsafe.Add(mBase, uint32(v971+int32(-28)))) = v966
									v1006 = v962&int32(4) | int32(24)
									v1007 = v970 - v1006
									if base.Ui32(v1007) < base.Ui32(int32(32)) {
									} else {
										v1012 = base.I64_extend_i32_u(v966) * int64(4294967297)
										v1015 = v1007
										v1016 = v962 + v1006
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1016)+24)) = v1012
											*(*int64)(unsafe.Add(mBase, uint32(v1016)+16)) = v1012
											*(*int64)(unsafe.Add(mBase, uint32(v1016)+8)) = v1012
											*(*int64)(unsafe.Add(mBase, uint32(v1016))) = v1012
											v1028 = v1015 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1028) {
												v1015 = v1028
												v1016 = v1016 + int32(32)
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
			base.MemoryFill(m, v923, v925, v802)
		}
		v1046 = int32(0)
		v1047 = base.Simd_g_v128_load_rng(m, v791, v1046, int32(0), int32(49))
		v1048 = int32(32)
		v1049 = base.Simd_g_v128_load_nc(m, v791, v1048)
		v1051 = base.Simd_g_v128_load_nc(m, v791, int32(1))
		v1052 = base.Simd_g_i8x16_avgr_u(v1049, v1051)
		v1053 = int32(33)
		v1054 = base.Simd_g_v128_load_nc(m, v791, v1053)
		v1055 = base.Simd_g_i8x16_avgr_u(v1047, v1054)
		v1057 = base.Simd_g_v128_xor(v1051, v1049)
		v1058 = base.Simd_g_v128_xor(v1047, v1054)
		v1060 = base.Simd_g_v128_xor(v1052, v1055)
		v1062 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k0)
		v1064 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1052, v1055), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1057, v1058), v1060), v1062))
		v1070 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1064, v1052), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1064, v1052), base.Simd_g_v128_and(v1060, v1057)), v1062))
		v1071 = base.Simd_g_i8x16_avgr_u(v1047, v1070)
		v1077 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1064, v1055), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1064, v1055), base.Simd_g_v128_and(v1060, v1058)), v1062))
		v1078 = base.Simd_g_i8x16_avgr_u(v1051, v1077)
		v1079 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k1)
		v1080 = base.Simd_g_i8x16_shuffle2(v1071, v1078, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v344, int32(80), v1080)
		v1083 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
		v1084 = base.Simd_g_i8x16_shuffle2(v1071, v1078, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v344, int32(64), v1084)
		v1087 = base.Simd_g_i8x16_avgr_u(v1049, v1077)
		v1088 = base.Simd_g_i8x16_avgr_u(v1054, v1070)
		v1090 = base.Simd_g_i8x16_shuffle2(v1087, v1088, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v344, int32(16), v1090)
		v1094 = base.Simd_g_i8x16_shuffle2(v1087, v1088, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v344, v1046, v1094)
		v1100 = F_memcpy(m, v791+v1048, l3+v755, v788)
		mBase = m.M
		v1102 = F_memcpy(m, v791, l5+v755, v788)
		mBase = m.M
		v1103 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v799))))
		if base.Ui32(v802) < base.Ui32(v1053) {
			if v802 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v794))) = uint8(v1103)
				v1114 = v794 + v802
				*(*uint8)(unsafe.Add(mBase, uint32(v1114+int32(-1)))) = uint8(v1103)
				if base.Ui32(v802) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v794)+2)) = uint8(v1103)
					*(*uint8)(unsafe.Add(mBase, uint32(v794)+1)) = uint8(v1103)
					*(*uint8)(unsafe.Add(mBase, uint32(v1114+int32(-3)))) = uint8(v1103)
					*(*uint8)(unsafe.Add(mBase, uint32(v1114+int32(-2)))) = uint8(v1103)
					if base.Ui32(v802) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v794)+3)) = uint8(v1103)
						*(*uint8)(unsafe.Add(mBase, uint32(v1114+int32(-4)))) = uint8(v1103)
						if base.Ui32(v802) < base.Ui32(int32(9)) {
						} else {
							v1139 = (int32(0) - v794) & int32(3)
							v1140 = v794 + v1139
							v1144 = v1103 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1140))) = v1144
							v1148 = (v802 - v1139) & int32(60)
							v1149 = v1140 + v1148
							*(*int32)(unsafe.Add(mBase, uint32(v1149+int32(-4)))) = v1144
							if base.Ui32(v1148) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1140)+8)) = v1144
								*(*int32)(unsafe.Add(mBase, uint32(v1140)+4)) = v1144
								*(*int32)(unsafe.Add(mBase, uint32(v1149+int32(-8)))) = v1144
								*(*int32)(unsafe.Add(mBase, uint32(v1149+int32(-12)))) = v1144
								if base.Ui32(v1148) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1140)+24)) = v1144
									*(*int32)(unsafe.Add(mBase, uint32(v1140)+20)) = v1144
									*(*int32)(unsafe.Add(mBase, uint32(v1140)+16)) = v1144
									*(*int32)(unsafe.Add(mBase, uint32(v1140)+12)) = v1144
									*(*int32)(unsafe.Add(mBase, uint32(v1149+int32(-16)))) = v1144
									*(*int32)(unsafe.Add(mBase, uint32(v1149+int32(-20)))) = v1144
									*(*int32)(unsafe.Add(mBase, uint32(v1149+int32(-24)))) = v1144
									*(*int32)(unsafe.Add(mBase, uint32(v1149+int32(-28)))) = v1144
									v1184 = v1140&int32(4) | int32(24)
									v1185 = v1148 - v1184
									if base.Ui32(v1185) < base.Ui32(int32(32)) {
									} else {
										v1190 = base.I64_extend_i32_u(v1144) * int64(4294967297)
										v1193 = v1185
										v1194 = v1140 + v1184
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1194)+24)) = v1190
											*(*int64)(unsafe.Add(mBase, uint32(v1194)+16)) = v1190
											*(*int64)(unsafe.Add(mBase, uint32(v1194)+8)) = v1190
											*(*int64)(unsafe.Add(mBase, uint32(v1194))) = v1190
											v1206 = v1193 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1206) {
												v1193 = v1206
												v1194 = v1194 + int32(32)
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
			base.MemoryFill(m, v794, v1103, v802)
		}
		v1224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v924))))
		if base.Ui32(v802) < base.Ui32(int32(33)) {
			if v802 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v923))) = uint8(v1224)
				v1235 = v923 + v802
				*(*uint8)(unsafe.Add(mBase, uint32(v1235+int32(-1)))) = uint8(v1224)
				if base.Ui32(v802) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v923)+2)) = uint8(v1224)
					*(*uint8)(unsafe.Add(mBase, uint32(v923)+1)) = uint8(v1224)
					*(*uint8)(unsafe.Add(mBase, uint32(v1235+int32(-3)))) = uint8(v1224)
					*(*uint8)(unsafe.Add(mBase, uint32(v1235+int32(-2)))) = uint8(v1224)
					if base.Ui32(v802) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v923)+3)) = uint8(v1224)
						*(*uint8)(unsafe.Add(mBase, uint32(v1235+int32(-4)))) = uint8(v1224)
						if base.Ui32(v802) < base.Ui32(int32(9)) {
						} else {
							v1260 = (int32(0) - v923) & int32(3)
							v1261 = v923 + v1260
							v1265 = v1224 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1261))) = v1265
							v1269 = (v802 - v1260) & int32(60)
							v1270 = v1261 + v1269
							*(*int32)(unsafe.Add(mBase, uint32(v1270+int32(-4)))) = v1265
							if base.Ui32(v1269) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1261)+8)) = v1265
								*(*int32)(unsafe.Add(mBase, uint32(v1261)+4)) = v1265
								*(*int32)(unsafe.Add(mBase, uint32(v1270+int32(-8)))) = v1265
								*(*int32)(unsafe.Add(mBase, uint32(v1270+int32(-12)))) = v1265
								if base.Ui32(v1269) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1261)+24)) = v1265
									*(*int32)(unsafe.Add(mBase, uint32(v1261)+20)) = v1265
									*(*int32)(unsafe.Add(mBase, uint32(v1261)+16)) = v1265
									*(*int32)(unsafe.Add(mBase, uint32(v1261)+12)) = v1265
									*(*int32)(unsafe.Add(mBase, uint32(v1270+int32(-16)))) = v1265
									*(*int32)(unsafe.Add(mBase, uint32(v1270+int32(-20)))) = v1265
									*(*int32)(unsafe.Add(mBase, uint32(v1270+int32(-24)))) = v1265
									*(*int32)(unsafe.Add(mBase, uint32(v1270+int32(-28)))) = v1265
									v1305 = v1261&int32(4) | int32(24)
									v1306 = v1269 - v1305
									if base.Ui32(v1306) < base.Ui32(int32(32)) {
									} else {
										v1311 = base.I64_extend_i32_u(v1265) * int64(4294967297)
										v1314 = v1306
										v1315 = v1261 + v1305
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1315)+24)) = v1311
											*(*int64)(unsafe.Add(mBase, uint32(v1315)+16)) = v1311
											*(*int64)(unsafe.Add(mBase, uint32(v1315)+8)) = v1311
											*(*int64)(unsafe.Add(mBase, uint32(v1315))) = v1311
											v1327 = v1314 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1327) {
												v1314 = v1327
												v1315 = v1315 + int32(32)
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
			base.MemoryFill(m, v923, v1224, v802)
		}
		v1346 = base.Simd_g_v128_load_rng(m, v1102, int32(0), int32(0), int32(49))
		v1347 = int32(32)
		v1348 = base.Simd_g_v128_load_nc(m, v1102, v1347)
		v1350 = base.Simd_g_v128_load_nc(m, v1102, int32(1))
		v1351 = base.Simd_g_i8x16_avgr_u(v1348, v1350)
		v1353 = base.Simd_g_v128_load_nc(m, v1102, int32(33))
		v1354 = base.Simd_g_i8x16_avgr_u(v1346, v1353)
		v1356 = base.Simd_g_v128_xor(v1350, v1348)
		v1357 = base.Simd_g_v128_xor(v1346, v1353)
		v1359 = base.Simd_g_v128_xor(v1351, v1354)
		v1362 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1351, v1354), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1356, v1357), v1359), v1062))
		v1368 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1362, v1351), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1362, v1351), base.Simd_g_v128_and(v1359, v1356)), v1062))
		v1369 = base.Simd_g_i8x16_avgr_u(v1346, v1368)
		v1375 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1362, v1354), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1362, v1354), base.Simd_g_v128_and(v1359, v1357)), v1062))
		v1376 = base.Simd_g_i8x16_avgr_u(v1350, v1375)
		v1377 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k1)
		v1378 = base.Simd_g_i8x16_shuffle2(v1369, v1376, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v344, int32(112), v1378)
		v1381 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
		v1382 = base.Simd_g_i8x16_shuffle2(v1369, v1376, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v344, int32(96), v1382)
		v1385 = base.Simd_g_i8x16_avgr_u(v1348, v1375)
		v1386 = base.Simd_g_i8x16_avgr_u(v1353, v1368)
		v1388 = base.Simd_g_i8x16_shuffle2(v1385, v1386, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v344, int32(48), v1388)
		v1392 = base.Simd_g_i8x16_shuffle2(v1385, v1386, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v344, v1347, v1392)
		v1398 = l8 - v756
		v1399 = F_memcpy(m, v35+int32(448), l0+v756, v1398)
		mBase = m.M
		v1401 = v35 + int32(192)
		if l1 != 0 {
			v1526 = F_memcpy(m, v35+int32(480), l1+v756, v1398)
			mBase = m.M
			v1540 = v1401
			v1541 = int32(0)
			for {
				v1550 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k7)
				v1552 = int32(0)
				v1553 = base.Simd_g_v128_load64_zero(m, v1399+v1541, v1552)
				v1554 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
				v1555 = base.Simd_g_i8x16_shuffle2(v1550, v1553, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1557 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k8)
				v1561 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k9)
				v1562 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1555), v1557), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1555), v1557), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))
				v1565 = base.Simd_g_v128_load64_zero(m, v344+v1541, v1552)
				v1567 = base.Simd_g_i8x16_shuffle2(v1550, v1565, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1568 = base.Simd_g_i32x4_extend_low_i16x8_u(v1567)
				v1569 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k12)
				v1571 = base.Simd_g_i32x4_extend_high_i16x8_u(v1567)
				v1577 = base.Simd_g_v128_load64_zero(m, v342+v1541, v1552)
				v1579 = base.Simd_g_i8x16_shuffle2(v1550, v1577, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1580 = base.Simd_g_i32x4_extend_low_i16x8_u(v1579)
				v1581 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k13)
				v1583 = base.Simd_g_i32x4_extend_high_i16x8_u(v1579)
				v1591 = int32(6)
				v1592 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1562, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1568, v1569), base.Simd_g_i32x4_mul(v1571, v1569), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1580, v1581), base.Simd_g_i32x4_mul(v1583, v1581), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k14)), v1591)
				v1593 = base.Simd_g_i8x16_narrow_i16x8_u(v1592, v1592)
				v1598 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k15)
				v1607 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1562, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1580, v1598), base.Simd_g_i32x4_mul(v1583, v1598), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k16)), v1591)
				v1612 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k17)
				v1621 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1568, v1612), base.Simd_g_i32x4_mul(v1571, v1612), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), v1562), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k18)), v1591)
				v1623 = int32(3)
				v1633 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(v1593, int32(5)), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k19)), base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(v1607, v1607), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k20))), base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_narrow_i16x8_u(v1621, v1621), v1623), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k21)), base.Simd_g_v128_and(base.Simd_g_i16x8_shl(v1593, v1623), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k22))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				base.Simd_g_v128_store(m, v1540, v1552, v1633)
				if base.Ui32(v1541) < base.Ui32(int32(24)) {
					v1540 = v1540 + int32(16)
					v1541 = v1541 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1647 = v35 + int32(320)
			v1661 = v1647
			v1662 = int32(0)
			for {
				v1671 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k7)
				v1673 = int32(0)
				v1674 = base.Simd_g_v128_load64_zero(m, v1526+v1662, v1673)
				v1675 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
				v1676 = base.Simd_g_i8x16_shuffle2(v1671, v1674, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1678 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k8)
				v1682 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k9)
				v1683 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1676), v1678), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1676), v1678), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))
				v1686 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v1662, v1673)
				v1688 = base.Simd_g_i8x16_shuffle2(v1671, v1686, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1689 = base.Simd_g_i32x4_extend_low_i16x8_u(v1688)
				v1690 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k12)
				v1692 = base.Simd_g_i32x4_extend_high_i16x8_u(v1688)
				v1698 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v1662, v1673)
				v1700 = base.Simd_g_i8x16_shuffle2(v1671, v1698, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1701 = base.Simd_g_i32x4_extend_low_i16x8_u(v1700)
				v1702 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k13)
				v1704 = base.Simd_g_i32x4_extend_high_i16x8_u(v1700)
				v1712 = int32(6)
				v1713 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1683, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1689, v1690), base.Simd_g_i32x4_mul(v1692, v1690), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1701, v1702), base.Simd_g_i32x4_mul(v1704, v1702), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k14)), v1712)
				v1714 = base.Simd_g_i8x16_narrow_i16x8_u(v1713, v1713)
				v1719 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k15)
				v1728 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1683, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1701, v1719), base.Simd_g_i32x4_mul(v1704, v1719), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k16)), v1712)
				v1733 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k17)
				v1742 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1689, v1733), base.Simd_g_i32x4_mul(v1692, v1733), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), v1683), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k18)), v1712)
				v1744 = int32(3)
				v1754 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(v1714, int32(5)), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k19)), base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(v1728, v1728), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k20))), base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_narrow_i16x8_u(v1742, v1742), v1744), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k21)), base.Simd_g_v128_and(base.Simd_g_i16x8_shl(v1714, v1744), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k22))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				base.Simd_g_v128_store(m, v1661, v1673, v1754)
				if base.Ui32(v1662) < base.Ui32(int32(24)) {
					v1661 = v1661 + int32(16)
					v1662 = v1662 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1763 = int32(1)
			v1764 = v756 << (uint(v1763) % 32)
			v1767 = v1398 << (uint(v1763) % 32)
			v1768 = F_memcpy(m, l6+v1764, v1401, v1767)
			mBase = m.M
			v1770 = F_memcpy(m, l7+v1764, v1647, v1767)
			mBase = m.M
		} else {
			v1415 = v1401
			v1416 = int32(0)
			for {
				v1425 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k7)
				v1427 = int32(0)
				v1428 = base.Simd_g_v128_load64_zero(m, v1399+v1416, v1427)
				v1429 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k4)
				v1430 = base.Simd_g_i8x16_shuffle2(v1425, v1428, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1432 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k8)
				v1436 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k9)
				v1437 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1430), v1432), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1430), v1432), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))
				v1440 = base.Simd_g_v128_load64_zero(m, v344+v1416, v1427)
				v1442 = base.Simd_g_i8x16_shuffle2(v1425, v1440, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1443 = base.Simd_g_i32x4_extend_low_i16x8_u(v1442)
				v1444 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k12)
				v1446 = base.Simd_g_i32x4_extend_high_i16x8_u(v1442)
				v1452 = base.Simd_g_v128_load64_zero(m, v342+v1416, v1427)
				v1454 = base.Simd_g_i8x16_shuffle2(v1425, v1452, base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				v1455 = base.Simd_g_i32x4_extend_low_i16x8_u(v1454)
				v1456 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k13)
				v1458 = base.Simd_g_i32x4_extend_high_i16x8_u(v1454)
				v1466 = int32(6)
				v1467 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1437, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1443, v1444), base.Simd_g_i32x4_mul(v1446, v1444), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1455, v1456), base.Simd_g_i32x4_mul(v1458, v1456), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k14)), v1466)
				v1468 = base.Simd_g_i8x16_narrow_i16x8_u(v1467, v1467)
				v1473 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k15)
				v1482 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1437, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1455, v1473), base.Simd_g_i32x4_mul(v1458, v1473), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k16)), v1466)
				v1487 = base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k17)
				v1496 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1443, v1487), base.Simd_g_i32x4_mul(v1446, v1487), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k11)), v1437), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k18)), v1466)
				v1498 = int32(3)
				v1508 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(v1468, int32(5)), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k19)), base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(v1482, v1482), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k20))), base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_narrow_i16x8_u(v1496, v1496), v1498), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k21)), base.Simd_g_v128_and(base.Simd_g_i16x8_shl(v1468, v1498), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k22))), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgb565LinePair_SSE2__k6))
				base.Simd_g_v128_store(m, v1415, v1427, v1508)
				if base.Ui32(v1416) < base.Ui32(int32(24)) {
					v1415 = v1415 + int32(16)
					v1416 = v1416 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1517 = int32(1)
			v1522 = F_memcpy(m, l6+v756<<(uint(v1517)%32), v1401, v1398<<(uint(v1517)%32))
			mBase = m.M
		}
	}
	m.G0 = v35 + int32(528)
	return
}

var F_UpsampleRgb565LinePair_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleRgb565LinePair_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleRgb565LinePair_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleRgb565LinePair_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleRgb565LinePair_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleRgb565LinePair_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleRgb565LinePair_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleRgb565LinePair_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleRgb565LinePair_SSE2__k8 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleRgb565LinePair_SSE2__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleRgb565LinePair_SSE2__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleRgb565LinePair_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleRgb565LinePair_SSE2__k12 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleRgb565LinePair_SSE2__k13 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleRgb565LinePair_SSE2__k14 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_UpsampleRgb565LinePair_SSE2__k15 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleRgb565LinePair_SSE2__k16 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleRgb565LinePair_SSE2__k17 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleRgb565LinePair_SSE2__k18 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleRgb565LinePair_SSE2__k19 = [2]uint64{0x707070707070707, 0x707070707070707}
var F_UpsampleRgb565LinePair_SSE2__k20 = [2]uint64{0xf8f8f8f8f8f8f8f8, 0xf8f8f8f8f8f8f8f8}
var F_UpsampleRgb565LinePair_SSE2__k21 = [2]uint64{0x1f1f1f1f1f1f1f1f, 0x1f1f1f1f1f1f1f1f}
var F_UpsampleRgb565LinePair_SSE2__k22 = [2]uint64{0xe0e0e0e0e0e0e0e0, 0xe0e0e0e0e0e0e0e0}

func F_UpsampleRgbLinePair_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
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
	var v518 base.V128
	_ = v518
	var v519 base.V128
	_ = v519
	var v521 base.V128
	_ = v521
	var v523 base.V128
	_ = v523
	var v526 base.V128
	_ = v526
	var v528 base.V128
	_ = v528
	var v530 base.V128
	_ = v530
	var v535 base.V128
	_ = v535
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
	var v545 base.V128
	_ = v545
	var v547 base.V128
	_ = v547
	var v552 base.V128
	_ = v552
	var v554 base.V128
	_ = v554
	var v560 base.V128
	_ = v560
	var v565 base.V128
	_ = v565
	var v566 base.V128
	_ = v566
	var v569 base.V128
	_ = v569
	var v571 base.V128
	_ = v571
	var v572 base.V128
	_ = v572
	var v574 base.V128
	_ = v574
	var v579 base.V128
	_ = v579
	var v581 base.V128
	_ = v581
	var v587 base.V128
	_ = v587
	var v592 int32
	_ = v592
	var v593 base.V128
	_ = v593
	var v595 base.V128
	_ = v595
	var v596 base.V128
	_ = v596
	var v598 base.V128
	_ = v598
	var v603 base.V128
	_ = v603
	var v605 base.V128
	_ = v605
	var v611 base.V128
	_ = v611
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
	var v904 base.V128
	_ = v904
	var v905 base.V128
	_ = v905
	var v907 base.V128
	_ = v907
	var v909 base.V128
	_ = v909
	var v912 base.V128
	_ = v912
	var v914 base.V128
	_ = v914
	var v916 base.V128
	_ = v916
	var v921 base.V128
	_ = v921
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
	var v931 base.V128
	_ = v931
	var v933 base.V128
	_ = v933
	var v938 base.V128
	_ = v938
	var v940 base.V128
	_ = v940
	var v946 base.V128
	_ = v946
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
	var v958 base.V128
	_ = v958
	var v960 base.V128
	_ = v960
	var v965 base.V128
	_ = v965
	var v967 base.V128
	_ = v967
	var v973 base.V128
	_ = v973
	var v978 int32
	_ = v978
	var v979 base.V128
	_ = v979
	var v981 base.V128
	_ = v981
	var v982 base.V128
	_ = v982
	var v984 base.V128
	_ = v984
	var v989 base.V128
	_ = v989
	var v991 base.V128
	_ = v991
	var v997 base.V128
	_ = v997
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
	var v1957 base.V128
	_ = v1957
	var v1958 base.V128
	_ = v1958
	var v1960 base.V128
	_ = v1960
	var v1962 base.V128
	_ = v1962
	var v1965 base.V128
	_ = v1965
	var v1967 base.V128
	_ = v1967
	var v1969 base.V128
	_ = v1969
	var v1974 base.V128
	_ = v1974
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
	var v1984 base.V128
	_ = v1984
	var v1986 base.V128
	_ = v1986
	var v1991 base.V128
	_ = v1991
	var v1993 base.V128
	_ = v1993
	var v1999 base.V128
	_ = v1999
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
	var v2011 base.V128
	_ = v2011
	var v2013 base.V128
	_ = v2013
	var v2018 base.V128
	_ = v2018
	var v2020 base.V128
	_ = v2020
	var v2026 base.V128
	_ = v2026
	var v2031 int32
	_ = v2031
	var v2032 base.V128
	_ = v2032
	var v2034 base.V128
	_ = v2034
	var v2035 base.V128
	_ = v2035
	var v2037 base.V128
	_ = v2037
	var v2042 base.V128
	_ = v2042
	var v2044 base.V128
	_ = v2044
	var v2050 base.V128
	_ = v2050
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
	var v2350 base.V128
	_ = v2350
	var v2351 base.V128
	_ = v2351
	var v2353 base.V128
	_ = v2353
	var v2355 base.V128
	_ = v2355
	var v2358 base.V128
	_ = v2358
	var v2360 base.V128
	_ = v2360
	var v2362 base.V128
	_ = v2362
	var v2367 base.V128
	_ = v2367
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
	var v2377 base.V128
	_ = v2377
	var v2379 base.V128
	_ = v2379
	var v2384 base.V128
	_ = v2384
	var v2386 base.V128
	_ = v2386
	var v2392 base.V128
	_ = v2392
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
	var v2404 base.V128
	_ = v2404
	var v2406 base.V128
	_ = v2406
	var v2411 base.V128
	_ = v2411
	var v2413 base.V128
	_ = v2413
	var v2419 base.V128
	_ = v2419
	var v2424 int32
	_ = v2424
	var v2425 base.V128
	_ = v2425
	var v2427 base.V128
	_ = v2427
	var v2428 base.V128
	_ = v2428
	var v2430 base.V128
	_ = v2430
	var v2435 base.V128
	_ = v2435
	var v2437 base.V128
	_ = v2437
	var v2443 base.V128
	_ = v2443
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
	var v2739 base.V128
	_ = v2739
	var v2740 base.V128
	_ = v2740
	var v2742 base.V128
	_ = v2742
	var v2744 base.V128
	_ = v2744
	var v2747 base.V128
	_ = v2747
	var v2749 base.V128
	_ = v2749
	var v2751 base.V128
	_ = v2751
	var v2756 base.V128
	_ = v2756
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
	var v2766 base.V128
	_ = v2766
	var v2768 base.V128
	_ = v2768
	var v2773 base.V128
	_ = v2773
	var v2775 base.V128
	_ = v2775
	var v2781 base.V128
	_ = v2781
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
	var v2793 base.V128
	_ = v2793
	var v2795 base.V128
	_ = v2795
	var v2800 base.V128
	_ = v2800
	var v2802 base.V128
	_ = v2802
	var v2808 base.V128
	_ = v2808
	var v2813 int32
	_ = v2813
	var v2814 base.V128
	_ = v2814
	var v2816 base.V128
	_ = v2816
	var v2817 base.V128
	_ = v2817
	var v2819 base.V128
	_ = v2819
	var v2824 base.V128
	_ = v2824
	var v2826 base.V128
	_ = v2826
	var v2832 base.V128
	_ = v2832
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
	v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v169 = int32(base.Ui32(v163+v164)>>(uint(v162)%32)) + v162
	v172 = int32(base.Ui32(v169+v164) >> (uint(v162) % 32))
	v175 = int32(8)
	v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v181 = int32(base.Ui32(v177*int32(_a_F_UpsampleRgbLinePair_SSE2_0)) >> (uint(v175) % 32))
	v182 = int32(base.Ui32(v172*int32(_a_F_UpsampleRgbLinePair_SSE2_1))>>(uint(v175)%32)) + v181
	v184 = v182 + int32(-17685)
	if base.Ui32(v182) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_2)) {
		v191 = int32(0)
	} else {
		v191 = int32(255)
	}
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_3)) {
		v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
	} else {
		v194 = v191
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v194)
	v196 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v199 = int32(1)
	v202 = int32(base.Ui32(v196+v197)>>(uint(v199)%32)) + v199
	v205 = int32(base.Ui32(v202+v197) >> (uint(v199) % 32))
	v210 = int32(base.Ui32(v205*int32(_a_F_UpsampleRgbLinePair_SSE2_4))>>(uint(int32(8))%32)) + v181
	v212 = v210 + int32(-14234)
	if base.Ui32(v210) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_5)) {
		v219 = int32(0)
	} else {
		v219 = int32(255)
	}
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_3)) {
		v222 = int32(base.Ui32(v212) >> (uint(int32(6)) % 32))
	} else {
		v222 = v219
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v222)
	v226 = int32(8)
	v233 = v181 - (int32(base.Ui32(v172*int32(_a_F_UpsampleRgbLinePair_SSE2_6))>>(uint(v226)%32)) + int32(base.Ui32(v205*int32(_a_F_UpsampleRgbLinePair_SSE2_7))>>(uint(v226)%32)))
	v235 = v233 + int32(_a_F_UpsampleRgbLinePair_SSE2_8)
	if v233 < int32(-8708) {
		v242 = int32(0)
	} else {
		v242 = int32(255)
	}
	if base.Ui32(v235) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_3)) {
		v245 = int32(base.Ui32(v235) >> (uint(int32(6)) % 32))
	} else {
		v245 = v242
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v245)
	if l1 == int32(0) {
	} else {
		v249 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v252 = int32(8)
		v253 = int32(base.Ui32(v249*int32(_a_F_UpsampleRgbLinePair_SSE2_0)) >> (uint(v252) % 32))
		v256 = int32(base.Ui32(v169+v163) >> (uint(int32(1)) % 32))
		v261 = v253 + int32(base.Ui32(v256*int32(_a_F_UpsampleRgbLinePair_SSE2_1))>>(uint(v252)%32))
		v263 = v261 + int32(-17685)
		if base.Ui32(v261) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_2)) {
			v270 = int32(0)
		} else {
			v270 = int32(255)
		}
		if base.Ui32(v263) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_3)) {
			v273 = int32(base.Ui32(v263) >> (uint(int32(6)) % 32))
		} else {
			v273 = v270
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v273)
		v277 = int32(base.Ui32(v202+v196) >> (uint(int32(1)) % 32))
		v282 = v253 + int32(base.Ui32(v277*int32(_a_F_UpsampleRgbLinePair_SSE2_4))>>(uint(int32(8))%32))
		v284 = v282 + int32(-14234)
		if base.Ui32(v282) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_5)) {
			v291 = int32(0)
		} else {
			v291 = int32(255)
		}
		if base.Ui32(v284) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_3)) {
			v294 = int32(base.Ui32(v284) >> (uint(int32(6)) % 32))
		} else {
			v294 = v291
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v294)
		v298 = int32(8)
		v305 = v253 - (int32(base.Ui32(v277*int32(_a_F_UpsampleRgbLinePair_SSE2_7))>>(uint(v298)%32)) + int32(base.Ui32(v256*int32(_a_F_UpsampleRgbLinePair_SSE2_6))>>(uint(v298)%32)))
		v307 = v305 + int32(_a_F_UpsampleRgbLinePair_SSE2_8)
		if v305 < int32(-8708) {
			v314 = int32(0)
		} else {
			v314 = int32(255)
		}
		if base.Ui32(v307) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE2_3)) {
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
			v394 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k0)
			v396 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v384, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v389, v390), v392), v394))
			v402 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v384), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v384), base.Simd_g_v128_and(v392, v389)), v394))
			v403 = base.Simd_g_i8x16_avgr_u(v378, v402)
			v409 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v387), base.Simd_g_v128_and(v392, v390)), v394))
			v410 = base.Simd_g_i8x16_avgr_u(v383, v409)
			v411 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k1)
			v412 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
			v413 = int32(80)
			base.Simd_g_v128_store(m, v327, v413, v412)
			v415 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k4)
			v416 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v417 = int32(64)
			base.Simd_g_v128_store(m, v327, v417, v416)
			v419 = base.Simd_g_i8x16_avgr_u(v381, v409)
			v420 = base.Simd_g_i8x16_avgr_u(v386, v402)
			v422 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
			v423 = int32(16)
			base.Simd_g_v128_store(m, v327, v423, v422)
			v426 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
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
			v464 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v327, int32(112), v464)
			v468 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v327, int32(96), v468)
			v471 = base.Simd_g_i8x16_avgr_u(v434, v461)
			v472 = base.Simd_g_i8x16_avgr_u(v439, v454)
			v474 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
			v475 = int32(48)
			base.Simd_g_v128_store(m, v327, v475, v474)
			v478 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v479 = int32(32)
			base.Simd_g_v128_store(m, v327, v479, v478)
			v481 = l0 + v330 + v359
			v482 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k7)
			v515 = base.Simd_g_v128_load64_zero(m, v327, v377)
			v517 = base.Simd_g_i8x16_shuffle2(v482, v515, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v518 = base.Simd_g_i32x4_extend_low_i16x8_u(v517)
			v519 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k8)
			v521 = base.Simd_g_i32x4_extend_high_i16x8_u(v517)
			v523 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k9)
			v526 = base.Simd_g_v128_load64_zero(m, v481, v377)
			v528 = base.Simd_g_i8x16_shuffle2(v482, v526, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v530 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k10)
			v535 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v528), v530), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v528), v530), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v537 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k13)
			v539 = int32(6)
			v541 = int32(8)
			v542 = base.Simd_g_v128_load64_zero(m, v327, v541)
			v544 = base.Simd_g_i8x16_shuffle2(v482, v542, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v545 = base.Simd_g_i32x4_extend_low_i16x8_u(v544)
			v547 = base.Simd_g_i32x4_extend_high_i16x8_u(v544)
			v552 = base.Simd_g_v128_load64_zero(m, v481, v541)
			v554 = base.Simd_g_i8x16_shuffle2(v482, v552, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v560 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v554), v530), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v554), v530), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v565 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v518, v519), base.Simd_g_i32x4_mul(v521, v519), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v535), v537), v539), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v545, v519), base.Simd_g_i32x4_mul(v547, v519), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v560), v537), v539))
			v566 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k14)
			v569 = base.Simd_g_v128_load64_zero(m, v327, v423)
			v571 = base.Simd_g_i8x16_shuffle2(v482, v569, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v572 = base.Simd_g_i32x4_extend_low_i16x8_u(v571)
			v574 = base.Simd_g_i32x4_extend_high_i16x8_u(v571)
			v579 = base.Simd_g_v128_load64_zero(m, v481, v423)
			v581 = base.Simd_g_i8x16_shuffle2(v482, v579, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v587 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v581), v530), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v581), v530), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v592 = int32(24)
			v593 = base.Simd_g_v128_load64_zero(m, v327, v592)
			v595 = base.Simd_g_i8x16_shuffle2(v482, v593, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v596 = base.Simd_g_i32x4_extend_low_i16x8_u(v595)
			v598 = base.Simd_g_i32x4_extend_high_i16x8_u(v595)
			v603 = base.Simd_g_v128_load64_zero(m, v481, v592)
			v605 = base.Simd_g_i8x16_shuffle2(v482, v603, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v611 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v605), v530), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v605), v530), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v616 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v572, v519), base.Simd_g_i32x4_mul(v574, v519), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v587), v537), v539), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v596, v519), base.Simd_g_i32x4_mul(v598, v519), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v611), v537), v539))
			v618 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v565, v566), base.Simd_g_v128_and(v616, v566))
			v622 = base.Simd_g_v128_load64_zero(m, v325, v377)
			v624 = base.Simd_g_i8x16_shuffle2(v482, v622, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v625 = base.Simd_g_i32x4_extend_low_i16x8_u(v624)
			v626 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k15)
			v628 = base.Simd_g_i32x4_extend_high_i16x8_u(v624)
			v633 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k16)
			v638 = base.Simd_g_v128_load64_zero(m, v325, v541)
			v640 = base.Simd_g_i8x16_shuffle2(v482, v638, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v641 = base.Simd_g_i32x4_extend_low_i16x8_u(v640)
			v643 = base.Simd_g_i32x4_extend_high_i16x8_u(v640)
			v651 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v535, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v625, v626), base.Simd_g_i32x4_mul(v628, v626), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v633), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v560, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v641, v626), base.Simd_g_i32x4_mul(v643, v626), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v633), v539))
			v655 = base.Simd_g_v128_load64_zero(m, v325, v423)
			v657 = base.Simd_g_i8x16_shuffle2(v482, v655, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v658 = base.Simd_g_i32x4_extend_low_i16x8_u(v657)
			v660 = base.Simd_g_i32x4_extend_high_i16x8_u(v657)
			v669 = base.Simd_g_v128_load64_zero(m, v325, v592)
			v671 = base.Simd_g_i8x16_shuffle2(v482, v669, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v672 = base.Simd_g_i32x4_extend_low_i16x8_u(v671)
			v674 = base.Simd_g_i32x4_extend_high_i16x8_u(v671)
			v682 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v587, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v658, v626), base.Simd_g_i32x4_mul(v660, v626), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v633), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v611, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v672, v626), base.Simd_g_i32x4_mul(v674, v626), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v633), v539))
			v685 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v651, v541), base.Simd_g_i16x8_shr_u(v682, v541))
			v688 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v618, v541), base.Simd_g_i16x8_shr_u(v685, v541))
			v690 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k17)
			v695 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k18)
			v702 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k19)
			v719 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v535, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v518, v690), base.Simd_g_i32x4_mul(v521, v690), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v625, v695), base.Simd_g_i32x4_mul(v628, v695), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v702), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v560, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v545, v690), base.Simd_g_i32x4_mul(v547, v690), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v641, v695), base.Simd_g_i32x4_mul(v643, v695), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v702), v539))
			v748 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v587, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v572, v690), base.Simd_g_i32x4_mul(v574, v690), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v658, v695), base.Simd_g_i32x4_mul(v660, v695), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v702), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v611, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v596, v690), base.Simd_g_i32x4_mul(v598, v690), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v672, v695), base.Simd_g_i32x4_mul(v674, v695), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v702), v539))
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
				v868 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k7)
				v900 = int32(0)
				v901 = base.Simd_g_v128_load64_zero(m, v341, v900)
				v902 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k4)
				v903 = base.Simd_g_i8x16_shuffle2(v868, v901, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v904 = base.Simd_g_i32x4_extend_low_i16x8_u(v903)
				v905 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k8)
				v907 = base.Simd_g_i32x4_extend_high_i16x8_u(v903)
				v909 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k9)
				v912 = base.Simd_g_v128_load64_zero(m, v867, v900)
				v914 = base.Simd_g_i8x16_shuffle2(v868, v912, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v916 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k10)
				v921 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v914), v916), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v914), v916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
				v923 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k13)
				v925 = int32(6)
				v927 = int32(8)
				v928 = base.Simd_g_v128_load64_zero(m, v341, v927)
				v930 = base.Simd_g_i8x16_shuffle2(v868, v928, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v931 = base.Simd_g_i32x4_extend_low_i16x8_u(v930)
				v933 = base.Simd_g_i32x4_extend_high_i16x8_u(v930)
				v938 = base.Simd_g_v128_load64_zero(m, v867, v927)
				v940 = base.Simd_g_i8x16_shuffle2(v868, v938, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v946 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v940), v916), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v940), v916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
				v951 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v904, v905), base.Simd_g_i32x4_mul(v907, v905), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v921), v923), v925), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v931, v905), base.Simd_g_i32x4_mul(v933, v905), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v946), v923), v925))
				v952 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k14)
				v954 = int32(16)
				v955 = base.Simd_g_v128_load64_zero(m, v341, v954)
				v957 = base.Simd_g_i8x16_shuffle2(v868, v955, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v958 = base.Simd_g_i32x4_extend_low_i16x8_u(v957)
				v960 = base.Simd_g_i32x4_extend_high_i16x8_u(v957)
				v965 = base.Simd_g_v128_load64_zero(m, v867, v954)
				v967 = base.Simd_g_i8x16_shuffle2(v868, v965, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v973 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v967), v916), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v967), v916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
				v978 = int32(24)
				v979 = base.Simd_g_v128_load64_zero(m, v341, v978)
				v981 = base.Simd_g_i8x16_shuffle2(v868, v979, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v982 = base.Simd_g_i32x4_extend_low_i16x8_u(v981)
				v984 = base.Simd_g_i32x4_extend_high_i16x8_u(v981)
				v989 = base.Simd_g_v128_load64_zero(m, v867, v978)
				v991 = base.Simd_g_i8x16_shuffle2(v868, v989, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v997 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v991), v916), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v991), v916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
				v1002 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v958, v905), base.Simd_g_i32x4_mul(v960, v905), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v973), v923), v925), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v982, v905), base.Simd_g_i32x4_mul(v984, v905), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v997), v923), v925))
				v1004 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v951, v952), base.Simd_g_v128_and(v1002, v952))
				v1008 = base.Simd_g_v128_load64_zero(m, v339, v900)
				v1010 = base.Simd_g_i8x16_shuffle2(v868, v1008, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v1011 = base.Simd_g_i32x4_extend_low_i16x8_u(v1010)
				v1012 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k15)
				v1014 = base.Simd_g_i32x4_extend_high_i16x8_u(v1010)
				v1019 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k16)
				v1024 = base.Simd_g_v128_load64_zero(m, v339, v927)
				v1026 = base.Simd_g_i8x16_shuffle2(v868, v1024, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v1027 = base.Simd_g_i32x4_extend_low_i16x8_u(v1026)
				v1029 = base.Simd_g_i32x4_extend_high_i16x8_u(v1026)
				v1037 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v921, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1011, v1012), base.Simd_g_i32x4_mul(v1014, v1012), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v1019), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v946, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1027, v1012), base.Simd_g_i32x4_mul(v1029, v1012), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v1019), v925))
				v1041 = base.Simd_g_v128_load64_zero(m, v339, v954)
				v1043 = base.Simd_g_i8x16_shuffle2(v868, v1041, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v1044 = base.Simd_g_i32x4_extend_low_i16x8_u(v1043)
				v1046 = base.Simd_g_i32x4_extend_high_i16x8_u(v1043)
				v1055 = base.Simd_g_v128_load64_zero(m, v339, v978)
				v1057 = base.Simd_g_i8x16_shuffle2(v868, v1055, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
				v1058 = base.Simd_g_i32x4_extend_low_i16x8_u(v1057)
				v1060 = base.Simd_g_i32x4_extend_high_i16x8_u(v1057)
				v1068 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v973, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1044, v1012), base.Simd_g_i32x4_mul(v1046, v1012), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v1019), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v997, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1058, v1012), base.Simd_g_i32x4_mul(v1060, v1012), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v1019), v925))
				v1071 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1037, v927), base.Simd_g_i16x8_shr_u(v1068, v927))
				v1074 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1004, v927), base.Simd_g_i16x8_shr_u(v1071, v927))
				v1076 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k17)
				v1081 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k18)
				v1088 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k19)
				v1105 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v921, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v904, v1076), base.Simd_g_i32x4_mul(v907, v1076), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1011, v1081), base.Simd_g_i32x4_mul(v1014, v1081), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v1088), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v946, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v931, v1076), base.Simd_g_i32x4_mul(v933, v1076), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1027, v1081), base.Simd_g_i32x4_mul(v1029, v1081), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v1088), v925))
				v1134 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v973, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v958, v1076), base.Simd_g_i32x4_mul(v960, v1076), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1044, v1081), base.Simd_g_i32x4_mul(v1046, v1081), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v1088), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v997, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v982, v1076), base.Simd_g_i32x4_mul(v984, v1076), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1058, v1081), base.Simd_g_i32x4_mul(v1060, v1081), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v1088), v925))
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
		v1581 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k0)
		v1583 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1571, v1574), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1576, v1577), v1579), v1581))
		v1589 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1583, v1571), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1583, v1571), base.Simd_g_v128_and(v1579, v1576)), v1581))
		v1590 = base.Simd_g_i8x16_avgr_u(v1566, v1589)
		v1596 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1583, v1574), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1583, v1574), base.Simd_g_v128_and(v1579, v1577)), v1581))
		v1597 = base.Simd_g_i8x16_avgr_u(v1570, v1596)
		v1598 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k1)
		v1599 = base.Simd_g_i8x16_shuffle2(v1590, v1597, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(80), v1599)
		v1602 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k4)
		v1603 = base.Simd_g_i8x16_shuffle2(v1590, v1597, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v327, int32(64), v1603)
		v1606 = base.Simd_g_i8x16_avgr_u(v1568, v1596)
		v1607 = base.Simd_g_i8x16_avgr_u(v1573, v1589)
		v1609 = base.Simd_g_i8x16_shuffle2(v1606, v1607, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(16), v1609)
		v1613 = base.Simd_g_i8x16_shuffle2(v1606, v1607, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
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
		v1896 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k1)
		v1897 = base.Simd_g_i8x16_shuffle2(v1888, v1895, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(112), v1897)
		v1900 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k4)
		v1901 = base.Simd_g_i8x16_shuffle2(v1888, v1895, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v327, int32(96), v1901)
		v1904 = base.Simd_g_i8x16_avgr_u(v1867, v1894)
		v1905 = base.Simd_g_i8x16_avgr_u(v1872, v1887)
		v1907 = base.Simd_g_i8x16_shuffle2(v1904, v1905, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(48), v1907)
		v1911 = base.Simd_g_i8x16_shuffle2(v1904, v1905, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v327, v1866, v1911)
		v1917 = l8 - v1275
		v1918 = F_memcpy(m, v35+int32(448), l0+v1275, v1917)
		mBase = m.M
		v1920 = v35 + int32(192)
		if l1 != 0 {
			v2313 = F_memcpy(m, v35+int32(480), l1+v1275, v1917)
			mBase = m.M
			v2314 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k7)
			v2346 = int32(0)
			v2347 = base.Simd_g_v128_load64_zero(m, v327, v2346)
			v2348 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k4)
			v2349 = base.Simd_g_i8x16_shuffle2(v2314, v2347, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2350 = base.Simd_g_i32x4_extend_low_i16x8_u(v2349)
			v2351 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k8)
			v2353 = base.Simd_g_i32x4_extend_high_i16x8_u(v2349)
			v2355 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k9)
			v2358 = base.Simd_g_v128_load64_zero(m, v1918, v2346)
			v2360 = base.Simd_g_i8x16_shuffle2(v2314, v2358, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2362 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k10)
			v2367 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2360), v2362), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2360), v2362), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2369 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k13)
			v2371 = int32(6)
			v2373 = int32(8)
			v2374 = base.Simd_g_v128_load64_zero(m, v327, v2373)
			v2376 = base.Simd_g_i8x16_shuffle2(v2314, v2374, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2377 = base.Simd_g_i32x4_extend_low_i16x8_u(v2376)
			v2379 = base.Simd_g_i32x4_extend_high_i16x8_u(v2376)
			v2384 = base.Simd_g_v128_load64_zero(m, v1918, v2373)
			v2386 = base.Simd_g_i8x16_shuffle2(v2314, v2384, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2392 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2386), v2362), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2386), v2362), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2397 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2350, v2351), base.Simd_g_i32x4_mul(v2353, v2351), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2367), v2369), v2371), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2377, v2351), base.Simd_g_i32x4_mul(v2379, v2351), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2392), v2369), v2371))
			v2398 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k14)
			v2400 = int32(16)
			v2401 = base.Simd_g_v128_load64_zero(m, v327, v2400)
			v2403 = base.Simd_g_i8x16_shuffle2(v2314, v2401, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2404 = base.Simd_g_i32x4_extend_low_i16x8_u(v2403)
			v2406 = base.Simd_g_i32x4_extend_high_i16x8_u(v2403)
			v2411 = base.Simd_g_v128_load64_zero(m, v1918, v2400)
			v2413 = base.Simd_g_i8x16_shuffle2(v2314, v2411, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2419 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2413), v2362), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2413), v2362), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2424 = int32(24)
			v2425 = base.Simd_g_v128_load64_zero(m, v327, v2424)
			v2427 = base.Simd_g_i8x16_shuffle2(v2314, v2425, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2428 = base.Simd_g_i32x4_extend_low_i16x8_u(v2427)
			v2430 = base.Simd_g_i32x4_extend_high_i16x8_u(v2427)
			v2435 = base.Simd_g_v128_load64_zero(m, v1918, v2424)
			v2437 = base.Simd_g_i8x16_shuffle2(v2314, v2435, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2443 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2437), v2362), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2437), v2362), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2448 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2404, v2351), base.Simd_g_i32x4_mul(v2406, v2351), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2419), v2369), v2371), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2428, v2351), base.Simd_g_i32x4_mul(v2430, v2351), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2443), v2369), v2371))
			v2450 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2397, v2398), base.Simd_g_v128_and(v2448, v2398))
			v2454 = base.Simd_g_v128_load64_zero(m, v325, v2346)
			v2456 = base.Simd_g_i8x16_shuffle2(v2314, v2454, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2457 = base.Simd_g_i32x4_extend_low_i16x8_u(v2456)
			v2458 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k15)
			v2460 = base.Simd_g_i32x4_extend_high_i16x8_u(v2456)
			v2465 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k16)
			v2470 = base.Simd_g_v128_load64_zero(m, v325, v2373)
			v2472 = base.Simd_g_i8x16_shuffle2(v2314, v2470, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2473 = base.Simd_g_i32x4_extend_low_i16x8_u(v2472)
			v2475 = base.Simd_g_i32x4_extend_high_i16x8_u(v2472)
			v2483 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2367, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2457, v2458), base.Simd_g_i32x4_mul(v2460, v2458), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2465), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2392, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2473, v2458), base.Simd_g_i32x4_mul(v2475, v2458), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2465), v2371))
			v2487 = base.Simd_g_v128_load64_zero(m, v325, v2400)
			v2489 = base.Simd_g_i8x16_shuffle2(v2314, v2487, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2490 = base.Simd_g_i32x4_extend_low_i16x8_u(v2489)
			v2492 = base.Simd_g_i32x4_extend_high_i16x8_u(v2489)
			v2501 = base.Simd_g_v128_load64_zero(m, v325, v2424)
			v2503 = base.Simd_g_i8x16_shuffle2(v2314, v2501, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2504 = base.Simd_g_i32x4_extend_low_i16x8_u(v2503)
			v2506 = base.Simd_g_i32x4_extend_high_i16x8_u(v2503)
			v2514 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2419, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2490, v2458), base.Simd_g_i32x4_mul(v2492, v2458), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2465), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2443, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2504, v2458), base.Simd_g_i32x4_mul(v2506, v2458), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2465), v2371))
			v2517 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2483, v2373), base.Simd_g_i16x8_shr_u(v2514, v2373))
			v2520 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2450, v2373), base.Simd_g_i16x8_shr_u(v2517, v2373))
			v2522 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k17)
			v2527 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k18)
			v2534 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k19)
			v2551 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2367, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2350, v2522), base.Simd_g_i32x4_mul(v2353, v2522), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2457, v2527), base.Simd_g_i32x4_mul(v2460, v2527), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2534), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2392, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2377, v2522), base.Simd_g_i32x4_mul(v2379, v2522), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2473, v2527), base.Simd_g_i32x4_mul(v2475, v2527), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2534), v2371))
			v2580 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2419, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2404, v2522), base.Simd_g_i32x4_mul(v2406, v2522), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2490, v2527), base.Simd_g_i32x4_mul(v2492, v2527), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2534), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2443, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2428, v2522), base.Simd_g_i32x4_mul(v2430, v2522), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2504, v2527), base.Simd_g_i32x4_mul(v2506, v2527), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2534), v2371))
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
			v2703 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k7)
			v2735 = int32(0)
			v2736 = base.Simd_g_v128_load64_zero(m, v2698, v2735)
			v2737 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k4)
			v2738 = base.Simd_g_i8x16_shuffle2(v2703, v2736, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2739 = base.Simd_g_i32x4_extend_low_i16x8_u(v2738)
			v2740 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k8)
			v2742 = base.Simd_g_i32x4_extend_high_i16x8_u(v2738)
			v2744 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k9)
			v2747 = base.Simd_g_v128_load64_zero(m, v2313, v2735)
			v2749 = base.Simd_g_i8x16_shuffle2(v2703, v2747, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2751 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k10)
			v2756 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2749), v2751), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2749), v2751), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2758 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k13)
			v2760 = int32(6)
			v2762 = int32(8)
			v2763 = base.Simd_g_v128_load64_zero(m, v2698, v2762)
			v2765 = base.Simd_g_i8x16_shuffle2(v2703, v2763, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2766 = base.Simd_g_i32x4_extend_low_i16x8_u(v2765)
			v2768 = base.Simd_g_i32x4_extend_high_i16x8_u(v2765)
			v2773 = base.Simd_g_v128_load64_zero(m, v2313, v2762)
			v2775 = base.Simd_g_i8x16_shuffle2(v2703, v2773, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2781 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2775), v2751), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2775), v2751), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2786 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2739, v2740), base.Simd_g_i32x4_mul(v2742, v2740), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2756), v2758), v2760), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2766, v2740), base.Simd_g_i32x4_mul(v2768, v2740), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2781), v2758), v2760))
			v2787 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k14)
			v2789 = int32(16)
			v2790 = base.Simd_g_v128_load64_zero(m, v2698, v2789)
			v2792 = base.Simd_g_i8x16_shuffle2(v2703, v2790, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2793 = base.Simd_g_i32x4_extend_low_i16x8_u(v2792)
			v2795 = base.Simd_g_i32x4_extend_high_i16x8_u(v2792)
			v2800 = base.Simd_g_v128_load64_zero(m, v2313, v2789)
			v2802 = base.Simd_g_i8x16_shuffle2(v2703, v2800, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2808 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2802), v2751), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2802), v2751), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2813 = int32(24)
			v2814 = base.Simd_g_v128_load64_zero(m, v2698, v2813)
			v2816 = base.Simd_g_i8x16_shuffle2(v2703, v2814, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2817 = base.Simd_g_i32x4_extend_low_i16x8_u(v2816)
			v2819 = base.Simd_g_i32x4_extend_high_i16x8_u(v2816)
			v2824 = base.Simd_g_v128_load64_zero(m, v2313, v2813)
			v2826 = base.Simd_g_i8x16_shuffle2(v2703, v2824, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2832 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2826), v2751), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2826), v2751), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2837 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2793, v2740), base.Simd_g_i32x4_mul(v2795, v2740), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2808), v2758), v2760), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2817, v2740), base.Simd_g_i32x4_mul(v2819, v2740), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2832), v2758), v2760))
			v2839 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2786, v2787), base.Simd_g_v128_and(v2837, v2787))
			v2843 = base.Simd_g_v128_load64_zero(m, v2700, v2735)
			v2845 = base.Simd_g_i8x16_shuffle2(v2703, v2843, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2846 = base.Simd_g_i32x4_extend_low_i16x8_u(v2845)
			v2847 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k15)
			v2849 = base.Simd_g_i32x4_extend_high_i16x8_u(v2845)
			v2854 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k16)
			v2859 = base.Simd_g_v128_load64_zero(m, v2700, v2762)
			v2861 = base.Simd_g_i8x16_shuffle2(v2703, v2859, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2862 = base.Simd_g_i32x4_extend_low_i16x8_u(v2861)
			v2864 = base.Simd_g_i32x4_extend_high_i16x8_u(v2861)
			v2872 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2756, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2846, v2847), base.Simd_g_i32x4_mul(v2849, v2847), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2854), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2781, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2862, v2847), base.Simd_g_i32x4_mul(v2864, v2847), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2854), v2760))
			v2876 = base.Simd_g_v128_load64_zero(m, v2700, v2789)
			v2878 = base.Simd_g_i8x16_shuffle2(v2703, v2876, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2879 = base.Simd_g_i32x4_extend_low_i16x8_u(v2878)
			v2881 = base.Simd_g_i32x4_extend_high_i16x8_u(v2878)
			v2890 = base.Simd_g_v128_load64_zero(m, v2700, v2813)
			v2892 = base.Simd_g_i8x16_shuffle2(v2703, v2890, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2893 = base.Simd_g_i32x4_extend_low_i16x8_u(v2892)
			v2895 = base.Simd_g_i32x4_extend_high_i16x8_u(v2892)
			v2903 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2808, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2879, v2847), base.Simd_g_i32x4_mul(v2881, v2847), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2854), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2832, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2893, v2847), base.Simd_g_i32x4_mul(v2895, v2847), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2854), v2760))
			v2906 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2872, v2762), base.Simd_g_i16x8_shr_u(v2903, v2762))
			v2909 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2839, v2762), base.Simd_g_i16x8_shr_u(v2906, v2762))
			v2911 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k17)
			v2916 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k18)
			v2923 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k19)
			v2940 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2756, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2739, v2911), base.Simd_g_i32x4_mul(v2742, v2911), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2846, v2916), base.Simd_g_i32x4_mul(v2849, v2916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2923), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2781, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2766, v2911), base.Simd_g_i32x4_mul(v2768, v2911), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2862, v2916), base.Simd_g_i32x4_mul(v2864, v2916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2923), v2760))
			v2969 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2808, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2793, v2911), base.Simd_g_i32x4_mul(v2795, v2911), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2879, v2916), base.Simd_g_i32x4_mul(v2881, v2916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2923), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2832, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2817, v2911), base.Simd_g_i32x4_mul(v2819, v2911), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2893, v2916), base.Simd_g_i32x4_mul(v2895, v2916), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2923), v2760))
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
			v1921 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k7)
			v1953 = int32(0)
			v1954 = base.Simd_g_v128_load64_zero(m, v327, v1953)
			v1955 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k4)
			v1956 = base.Simd_g_i8x16_shuffle2(v1921, v1954, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v1957 = base.Simd_g_i32x4_extend_low_i16x8_u(v1956)
			v1958 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k8)
			v1960 = base.Simd_g_i32x4_extend_high_i16x8_u(v1956)
			v1962 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k9)
			v1965 = base.Simd_g_v128_load64_zero(m, v1918, v1953)
			v1967 = base.Simd_g_i8x16_shuffle2(v1921, v1965, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v1969 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k10)
			v1974 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1967), v1969), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1967), v1969), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v1976 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k13)
			v1978 = int32(6)
			v1980 = int32(8)
			v1981 = base.Simd_g_v128_load64_zero(m, v327, v1980)
			v1983 = base.Simd_g_i8x16_shuffle2(v1921, v1981, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v1984 = base.Simd_g_i32x4_extend_low_i16x8_u(v1983)
			v1986 = base.Simd_g_i32x4_extend_high_i16x8_u(v1983)
			v1991 = base.Simd_g_v128_load64_zero(m, v1918, v1980)
			v1993 = base.Simd_g_i8x16_shuffle2(v1921, v1991, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v1999 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1993), v1969), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1993), v1969), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2004 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1957, v1958), base.Simd_g_i32x4_mul(v1960, v1958), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v1974), v1976), v1978), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1984, v1958), base.Simd_g_i32x4_mul(v1986, v1958), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v1999), v1976), v1978))
			v2005 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k14)
			v2007 = int32(16)
			v2008 = base.Simd_g_v128_load64_zero(m, v327, v2007)
			v2010 = base.Simd_g_i8x16_shuffle2(v1921, v2008, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2011 = base.Simd_g_i32x4_extend_low_i16x8_u(v2010)
			v2013 = base.Simd_g_i32x4_extend_high_i16x8_u(v2010)
			v2018 = base.Simd_g_v128_load64_zero(m, v1918, v2007)
			v2020 = base.Simd_g_i8x16_shuffle2(v1921, v2018, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2026 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2020), v1969), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2020), v1969), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2031 = int32(24)
			v2032 = base.Simd_g_v128_load64_zero(m, v327, v2031)
			v2034 = base.Simd_g_i8x16_shuffle2(v1921, v2032, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2035 = base.Simd_g_i32x4_extend_low_i16x8_u(v2034)
			v2037 = base.Simd_g_i32x4_extend_high_i16x8_u(v2034)
			v2042 = base.Simd_g_v128_load64_zero(m, v1918, v2031)
			v2044 = base.Simd_g_i8x16_shuffle2(v1921, v2042, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2050 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2044), v1969), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2044), v1969), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))
			v2055 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2011, v1958), base.Simd_g_i32x4_mul(v2013, v1958), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2026), v1976), v1978), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2035, v1958), base.Simd_g_i32x4_mul(v2037, v1958), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), v2050), v1976), v1978))
			v2057 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2004, v2005), base.Simd_g_v128_and(v2055, v2005))
			v2061 = base.Simd_g_v128_load64_zero(m, v325, v1953)
			v2063 = base.Simd_g_i8x16_shuffle2(v1921, v2061, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2064 = base.Simd_g_i32x4_extend_low_i16x8_u(v2063)
			v2065 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k15)
			v2067 = base.Simd_g_i32x4_extend_high_i16x8_u(v2063)
			v2072 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k16)
			v2077 = base.Simd_g_v128_load64_zero(m, v325, v1980)
			v2079 = base.Simd_g_i8x16_shuffle2(v1921, v2077, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2080 = base.Simd_g_i32x4_extend_low_i16x8_u(v2079)
			v2082 = base.Simd_g_i32x4_extend_high_i16x8_u(v2079)
			v2090 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1974, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2064, v2065), base.Simd_g_i32x4_mul(v2067, v2065), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2072), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1999, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2080, v2065), base.Simd_g_i32x4_mul(v2082, v2065), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2072), v1978))
			v2094 = base.Simd_g_v128_load64_zero(m, v325, v2007)
			v2096 = base.Simd_g_i8x16_shuffle2(v1921, v2094, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2097 = base.Simd_g_i32x4_extend_low_i16x8_u(v2096)
			v2099 = base.Simd_g_i32x4_extend_high_i16x8_u(v2096)
			v2108 = base.Simd_g_v128_load64_zero(m, v325, v2031)
			v2110 = base.Simd_g_i8x16_shuffle2(v1921, v2108, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k6))
			v2111 = base.Simd_g_i32x4_extend_low_i16x8_u(v2110)
			v2113 = base.Simd_g_i32x4_extend_high_i16x8_u(v2110)
			v2121 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2026, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2097, v2065), base.Simd_g_i32x4_mul(v2099, v2065), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2072), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2050, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2111, v2065), base.Simd_g_i32x4_mul(v2113, v2065), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12))), v2072), v1978))
			v2124 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2090, v1980), base.Simd_g_i16x8_shr_u(v2121, v1980))
			v2127 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2057, v1980), base.Simd_g_i16x8_shr_u(v2124, v1980))
			v2129 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k17)
			v2134 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k18)
			v2141 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k19)
			v2158 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1974, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1957, v2129), base.Simd_g_i32x4_mul(v1960, v2129), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2064, v2134), base.Simd_g_i32x4_mul(v2067, v2134), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2141), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1999, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1984, v2129), base.Simd_g_i32x4_mul(v1986, v2129), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2080, v2134), base.Simd_g_i32x4_mul(v2082, v2134), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2141), v1978))
			v2187 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2026, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2011, v2129), base.Simd_g_i32x4_mul(v2013, v2129), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2097, v2134), base.Simd_g_i32x4_mul(v2099, v2134), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2141), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2050, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2035, v2129), base.Simd_g_i32x4_mul(v2037, v2129), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2111, v2134), base.Simd_g_i32x4_mul(v2113, v2134), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE2__k12)))), v2141), v1978))
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

var F_UpsampleRgbLinePair_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleRgbLinePair_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleRgbLinePair_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleRgbLinePair_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleRgbLinePair_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleRgbLinePair_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleRgbLinePair_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleRgbLinePair_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleRgbLinePair_SSE2__k8 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleRgbLinePair_SSE2__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleRgbLinePair_SSE2__k10 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleRgbLinePair_SSE2__k11 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleRgbLinePair_SSE2__k12 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleRgbLinePair_SSE2__k13 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleRgbLinePair_SSE2__k14 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_UpsampleRgbLinePair_SSE2__k15 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleRgbLinePair_SSE2__k16 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleRgbLinePair_SSE2__k17 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleRgbLinePair_SSE2__k18 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleRgbLinePair_SSE2__k19 = [2]uint64{0x2204220422042204, 0x2204220422042204}

func F_UpsampleRgbLinePair_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
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
	var v510 base.V128
	_ = v510
	var v512 base.V128
	_ = v512
	var v514 base.V128
	_ = v514
	var v518 base.V128
	_ = v518
	var v519 base.V128
	_ = v519
	var v521 base.V128
	_ = v521
	var v523 base.V128
	_ = v523
	var v524 base.V128
	_ = v524
	var v525 base.V128
	_ = v525
	var v527 base.V128
	_ = v527
	var v532 base.V128
	_ = v532
	var v534 base.V128
	_ = v534
	var v535 base.V128
	_ = v535
	var v536 base.V128
	_ = v536
	var v538 base.V128
	_ = v538
	var v544 base.V128
	_ = v544
	var v546 int32
	_ = v546
	var v548 int32
	_ = v548
	var v549 base.V128
	_ = v549
	var v551 base.V128
	_ = v551
	var v557 base.V128
	_ = v557
	var v559 base.V128
	_ = v559
	var v561 base.V128
	_ = v561
	var v562 base.V128
	_ = v562
	var v564 base.V128
	_ = v564
	var v569 base.V128
	_ = v569
	var v571 base.V128
	_ = v571
	var v572 base.V128
	_ = v572
	var v574 base.V128
	_ = v574
	var v583 base.V128
	_ = v583
	var v584 base.V128
	_ = v584
	var v586 base.V128
	_ = v586
	var v592 base.V128
	_ = v592
	var v604 base.V128
	_ = v604
	var v605 base.V128
	_ = v605
	var v608 base.V128
	_ = v608
	var v614 base.V128
	_ = v614
	var v626 base.V128
	_ = v626
	var v627 base.V128
	_ = v627
	var v629 base.V128
	_ = v629
	var v632 base.V128
	_ = v632
	var v634 base.V128
	_ = v634
	var v637 base.V128
	_ = v637
	var v639 base.V128
	_ = v639
	var v642 base.V128
	_ = v642
	var v644 base.V128
	_ = v644
	var v647 base.V128
	_ = v647
	var v649 base.V128
	_ = v649
	var v653 base.V128
	_ = v653
	var v655 base.V128
	_ = v655
	var v661 base.V128
	_ = v661
	var v663 base.V128
	_ = v663
	var v665 base.V128
	_ = v665
	var v666 base.V128
	_ = v666
	var v668 base.V128
	_ = v668
	var v673 base.V128
	_ = v673
	var v675 base.V128
	_ = v675
	var v676 base.V128
	_ = v676
	var v678 base.V128
	_ = v678
	var v687 int32
	_ = v687
	var v688 base.V128
	_ = v688
	var v690 base.V128
	_ = v690
	var v696 base.V128
	_ = v696
	var v698 base.V128
	_ = v698
	var v700 base.V128
	_ = v700
	var v701 base.V128
	_ = v701
	var v703 base.V128
	_ = v703
	var v708 base.V128
	_ = v708
	var v710 base.V128
	_ = v710
	var v711 base.V128
	_ = v711
	var v713 base.V128
	_ = v713
	var v722 base.V128
	_ = v722
	var v740 base.V128
	_ = v740
	var v759 base.V128
	_ = v759
	var v761 base.V128
	_ = v761
	var v768 base.V128
	_ = v768
	var v775 base.V128
	_ = v775
	var v780 int32
	_ = v780
	var v781 base.V128
	_ = v781
	var v808 int32
	_ = v808
	var v809 base.V128
	_ = v809
	var v810 base.V128
	_ = v810
	var v811 base.V128
	_ = v811
	var v813 base.V128
	_ = v813
	var v817 base.V128
	_ = v817
	var v818 base.V128
	_ = v818
	var v820 base.V128
	_ = v820
	var v822 base.V128
	_ = v822
	var v823 base.V128
	_ = v823
	var v824 base.V128
	_ = v824
	var v826 base.V128
	_ = v826
	var v831 base.V128
	_ = v831
	var v833 base.V128
	_ = v833
	var v834 base.V128
	_ = v834
	var v835 base.V128
	_ = v835
	var v837 base.V128
	_ = v837
	var v843 base.V128
	_ = v843
	var v845 int32
	_ = v845
	var v847 int32
	_ = v847
	var v848 base.V128
	_ = v848
	var v850 base.V128
	_ = v850
	var v856 base.V128
	_ = v856
	var v858 base.V128
	_ = v858
	var v860 base.V128
	_ = v860
	var v861 base.V128
	_ = v861
	var v863 base.V128
	_ = v863
	var v868 base.V128
	_ = v868
	var v870 base.V128
	_ = v870
	var v871 base.V128
	_ = v871
	var v873 base.V128
	_ = v873
	var v882 base.V128
	_ = v882
	var v883 base.V128
	_ = v883
	var v885 base.V128
	_ = v885
	var v891 base.V128
	_ = v891
	var v903 base.V128
	_ = v903
	var v904 base.V128
	_ = v904
	var v907 base.V128
	_ = v907
	var v913 base.V128
	_ = v913
	var v925 base.V128
	_ = v925
	var v926 base.V128
	_ = v926
	var v928 base.V128
	_ = v928
	var v931 base.V128
	_ = v931
	var v933 base.V128
	_ = v933
	var v936 base.V128
	_ = v936
	var v938 base.V128
	_ = v938
	var v941 base.V128
	_ = v941
	var v943 base.V128
	_ = v943
	var v946 base.V128
	_ = v946
	var v948 base.V128
	_ = v948
	var v951 int32
	_ = v951
	var v952 base.V128
	_ = v952
	var v954 base.V128
	_ = v954
	var v960 base.V128
	_ = v960
	var v962 base.V128
	_ = v962
	var v964 base.V128
	_ = v964
	var v965 base.V128
	_ = v965
	var v967 base.V128
	_ = v967
	var v972 base.V128
	_ = v972
	var v974 base.V128
	_ = v974
	var v975 base.V128
	_ = v975
	var v977 base.V128
	_ = v977
	var v986 int32
	_ = v986
	var v987 base.V128
	_ = v987
	var v989 base.V128
	_ = v989
	var v995 base.V128
	_ = v995
	var v997 base.V128
	_ = v997
	var v999 base.V128
	_ = v999
	var v1000 base.V128
	_ = v1000
	var v1002 base.V128
	_ = v1002
	var v1007 base.V128
	_ = v1007
	var v1009 base.V128
	_ = v1009
	var v1010 base.V128
	_ = v1010
	var v1012 base.V128
	_ = v1012
	var v1021 base.V128
	_ = v1021
	var v1039 base.V128
	_ = v1039
	var v1058 base.V128
	_ = v1058
	var v1060 base.V128
	_ = v1060
	var v1067 base.V128
	_ = v1067
	var v1074 base.V128
	_ = v1074
	var v1077 int32
	_ = v1077
	var v1082 int32
	_ = v1082
	var v1100 int32
	_ = v1100
	var v1101 int32
	_ = v1101
	var v1124 int32
	_ = v1124
	var v1127 int32
	_ = v1127
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
	var v1144 int32
	_ = v1144
	var v1145 int32
	_ = v1145
	var v1147 int32
	_ = v1147
	var v1158 int32
	_ = v1158
	var v1183 int32
	_ = v1183
	var v1184 int32
	_ = v1184
	var v1188 int32
	_ = v1188
	var v1192 int32
	_ = v1192
	var v1193 int32
	_ = v1193
	var v1228 int32
	_ = v1228
	var v1229 int32
	_ = v1229
	var v1234 int64
	_ = v1234
	var v1237 int32
	_ = v1237
	var v1238 int32
	_ = v1238
	var v1250 int32
	_ = v1250
	var v1268 int32
	_ = v1268
	var v1269 int32
	_ = v1269
	var v1270 int32
	_ = v1270
	var v1281 int32
	_ = v1281
	var v1306 int32
	_ = v1306
	var v1307 int32
	_ = v1307
	var v1311 int32
	_ = v1311
	var v1315 int32
	_ = v1315
	var v1316 int32
	_ = v1316
	var v1351 int32
	_ = v1351
	var v1352 int32
	_ = v1352
	var v1357 int64
	_ = v1357
	var v1360 int32
	_ = v1360
	var v1361 int32
	_ = v1361
	var v1373 int32
	_ = v1373
	var v1391 int32
	_ = v1391
	var v1392 base.V128
	_ = v1392
	var v1393 int32
	_ = v1393
	var v1394 base.V128
	_ = v1394
	var v1396 base.V128
	_ = v1396
	var v1397 base.V128
	_ = v1397
	var v1398 int32
	_ = v1398
	var v1399 base.V128
	_ = v1399
	var v1400 base.V128
	_ = v1400
	var v1402 base.V128
	_ = v1402
	var v1403 base.V128
	_ = v1403
	var v1405 base.V128
	_ = v1405
	var v1407 base.V128
	_ = v1407
	var v1409 base.V128
	_ = v1409
	var v1415 base.V128
	_ = v1415
	var v1416 base.V128
	_ = v1416
	var v1422 base.V128
	_ = v1422
	var v1423 base.V128
	_ = v1423
	var v1424 base.V128
	_ = v1424
	var v1425 base.V128
	_ = v1425
	var v1428 base.V128
	_ = v1428
	var v1429 base.V128
	_ = v1429
	var v1432 base.V128
	_ = v1432
	var v1433 base.V128
	_ = v1433
	var v1435 base.V128
	_ = v1435
	var v1439 base.V128
	_ = v1439
	var v1445 int32
	_ = v1445
	var v1447 int32
	_ = v1447
	var v1448 int32
	_ = v1448
	var v1459 int32
	_ = v1459
	var v1484 int32
	_ = v1484
	var v1485 int32
	_ = v1485
	var v1489 int32
	_ = v1489
	var v1493 int32
	_ = v1493
	var v1494 int32
	_ = v1494
	var v1529 int32
	_ = v1529
	var v1530 int32
	_ = v1530
	var v1535 int64
	_ = v1535
	var v1538 int32
	_ = v1538
	var v1539 int32
	_ = v1539
	var v1551 int32
	_ = v1551
	var v1569 int32
	_ = v1569
	var v1580 int32
	_ = v1580
	var v1605 int32
	_ = v1605
	var v1606 int32
	_ = v1606
	var v1610 int32
	_ = v1610
	var v1614 int32
	_ = v1614
	var v1615 int32
	_ = v1615
	var v1650 int32
	_ = v1650
	var v1651 int32
	_ = v1651
	var v1656 int64
	_ = v1656
	var v1659 int32
	_ = v1659
	var v1660 int32
	_ = v1660
	var v1672 int32
	_ = v1672
	var v1691 base.V128
	_ = v1691
	var v1692 int32
	_ = v1692
	var v1693 base.V128
	_ = v1693
	var v1695 base.V128
	_ = v1695
	var v1696 base.V128
	_ = v1696
	var v1698 base.V128
	_ = v1698
	var v1699 base.V128
	_ = v1699
	var v1701 base.V128
	_ = v1701
	var v1702 base.V128
	_ = v1702
	var v1704 base.V128
	_ = v1704
	var v1707 base.V128
	_ = v1707
	var v1713 base.V128
	_ = v1713
	var v1714 base.V128
	_ = v1714
	var v1720 base.V128
	_ = v1720
	var v1721 base.V128
	_ = v1721
	var v1722 base.V128
	_ = v1722
	var v1723 base.V128
	_ = v1723
	var v1726 base.V128
	_ = v1726
	var v1727 base.V128
	_ = v1727
	var v1730 base.V128
	_ = v1730
	var v1731 base.V128
	_ = v1731
	var v1733 base.V128
	_ = v1733
	var v1737 base.V128
	_ = v1737
	var v1743 int32
	_ = v1743
	var v1744 int32
	_ = v1744
	var v1746 int32
	_ = v1746
	var v1747 base.V128
	_ = v1747
	var v1774 int32
	_ = v1774
	var v1775 base.V128
	_ = v1775
	var v1776 base.V128
	_ = v1776
	var v1777 base.V128
	_ = v1777
	var v1779 base.V128
	_ = v1779
	var v1783 base.V128
	_ = v1783
	var v1784 base.V128
	_ = v1784
	var v1786 base.V128
	_ = v1786
	var v1788 base.V128
	_ = v1788
	var v1789 base.V128
	_ = v1789
	var v1790 base.V128
	_ = v1790
	var v1792 base.V128
	_ = v1792
	var v1797 base.V128
	_ = v1797
	var v1799 base.V128
	_ = v1799
	var v1800 base.V128
	_ = v1800
	var v1801 base.V128
	_ = v1801
	var v1803 base.V128
	_ = v1803
	var v1809 base.V128
	_ = v1809
	var v1811 int32
	_ = v1811
	var v1813 int32
	_ = v1813
	var v1814 base.V128
	_ = v1814
	var v1816 base.V128
	_ = v1816
	var v1822 base.V128
	_ = v1822
	var v1824 base.V128
	_ = v1824
	var v1826 base.V128
	_ = v1826
	var v1827 base.V128
	_ = v1827
	var v1829 base.V128
	_ = v1829
	var v1834 base.V128
	_ = v1834
	var v1836 base.V128
	_ = v1836
	var v1837 base.V128
	_ = v1837
	var v1839 base.V128
	_ = v1839
	var v1848 base.V128
	_ = v1848
	var v1849 base.V128
	_ = v1849
	var v1851 base.V128
	_ = v1851
	var v1857 base.V128
	_ = v1857
	var v1869 base.V128
	_ = v1869
	var v1870 base.V128
	_ = v1870
	var v1873 base.V128
	_ = v1873
	var v1879 base.V128
	_ = v1879
	var v1891 base.V128
	_ = v1891
	var v1892 base.V128
	_ = v1892
	var v1894 base.V128
	_ = v1894
	var v1897 base.V128
	_ = v1897
	var v1899 base.V128
	_ = v1899
	var v1902 base.V128
	_ = v1902
	var v1904 base.V128
	_ = v1904
	var v1907 base.V128
	_ = v1907
	var v1909 base.V128
	_ = v1909
	var v1912 base.V128
	_ = v1912
	var v1914 base.V128
	_ = v1914
	var v1917 int32
	_ = v1917
	var v1918 base.V128
	_ = v1918
	var v1920 base.V128
	_ = v1920
	var v1926 base.V128
	_ = v1926
	var v1928 base.V128
	_ = v1928
	var v1930 base.V128
	_ = v1930
	var v1931 base.V128
	_ = v1931
	var v1933 base.V128
	_ = v1933
	var v1938 base.V128
	_ = v1938
	var v1940 base.V128
	_ = v1940
	var v1941 base.V128
	_ = v1941
	var v1943 base.V128
	_ = v1943
	var v1952 int32
	_ = v1952
	var v1953 base.V128
	_ = v1953
	var v1955 base.V128
	_ = v1955
	var v1961 base.V128
	_ = v1961
	var v1963 base.V128
	_ = v1963
	var v1965 base.V128
	_ = v1965
	var v1966 base.V128
	_ = v1966
	var v1968 base.V128
	_ = v1968
	var v1973 base.V128
	_ = v1973
	var v1975 base.V128
	_ = v1975
	var v1976 base.V128
	_ = v1976
	var v1978 base.V128
	_ = v1978
	var v1987 base.V128
	_ = v1987
	var v2005 base.V128
	_ = v2005
	var v2024 base.V128
	_ = v2024
	var v2026 base.V128
	_ = v2026
	var v2033 base.V128
	_ = v2033
	var v2040 base.V128
	_ = v2040
	var v2043 int32
	_ = v2043
	var v2048 int32
	_ = v2048
	var v2052 int32
	_ = v2052
	var v2053 base.V128
	_ = v2053
	var v2080 int32
	_ = v2080
	var v2081 base.V128
	_ = v2081
	var v2082 base.V128
	_ = v2082
	var v2083 base.V128
	_ = v2083
	var v2085 base.V128
	_ = v2085
	var v2089 base.V128
	_ = v2089
	var v2090 base.V128
	_ = v2090
	var v2092 base.V128
	_ = v2092
	var v2094 base.V128
	_ = v2094
	var v2095 base.V128
	_ = v2095
	var v2096 base.V128
	_ = v2096
	var v2098 base.V128
	_ = v2098
	var v2103 base.V128
	_ = v2103
	var v2105 base.V128
	_ = v2105
	var v2106 base.V128
	_ = v2106
	var v2107 base.V128
	_ = v2107
	var v2109 base.V128
	_ = v2109
	var v2115 base.V128
	_ = v2115
	var v2117 int32
	_ = v2117
	var v2119 int32
	_ = v2119
	var v2120 base.V128
	_ = v2120
	var v2122 base.V128
	_ = v2122
	var v2128 base.V128
	_ = v2128
	var v2130 base.V128
	_ = v2130
	var v2132 base.V128
	_ = v2132
	var v2133 base.V128
	_ = v2133
	var v2135 base.V128
	_ = v2135
	var v2140 base.V128
	_ = v2140
	var v2142 base.V128
	_ = v2142
	var v2143 base.V128
	_ = v2143
	var v2145 base.V128
	_ = v2145
	var v2154 base.V128
	_ = v2154
	var v2155 base.V128
	_ = v2155
	var v2157 base.V128
	_ = v2157
	var v2163 base.V128
	_ = v2163
	var v2175 base.V128
	_ = v2175
	var v2176 base.V128
	_ = v2176
	var v2179 base.V128
	_ = v2179
	var v2185 base.V128
	_ = v2185
	var v2197 base.V128
	_ = v2197
	var v2198 base.V128
	_ = v2198
	var v2200 base.V128
	_ = v2200
	var v2203 base.V128
	_ = v2203
	var v2205 base.V128
	_ = v2205
	var v2208 base.V128
	_ = v2208
	var v2210 base.V128
	_ = v2210
	var v2213 base.V128
	_ = v2213
	var v2215 base.V128
	_ = v2215
	var v2218 base.V128
	_ = v2218
	var v2220 base.V128
	_ = v2220
	var v2223 int32
	_ = v2223
	var v2224 base.V128
	_ = v2224
	var v2226 base.V128
	_ = v2226
	var v2232 base.V128
	_ = v2232
	var v2234 base.V128
	_ = v2234
	var v2236 base.V128
	_ = v2236
	var v2237 base.V128
	_ = v2237
	var v2239 base.V128
	_ = v2239
	var v2244 base.V128
	_ = v2244
	var v2246 base.V128
	_ = v2246
	var v2247 base.V128
	_ = v2247
	var v2249 base.V128
	_ = v2249
	var v2258 int32
	_ = v2258
	var v2259 base.V128
	_ = v2259
	var v2261 base.V128
	_ = v2261
	var v2267 base.V128
	_ = v2267
	var v2269 base.V128
	_ = v2269
	var v2271 base.V128
	_ = v2271
	var v2272 base.V128
	_ = v2272
	var v2274 base.V128
	_ = v2274
	var v2279 base.V128
	_ = v2279
	var v2281 base.V128
	_ = v2281
	var v2282 base.V128
	_ = v2282
	var v2284 base.V128
	_ = v2284
	var v2293 base.V128
	_ = v2293
	var v2311 base.V128
	_ = v2311
	var v2330 base.V128
	_ = v2330
	var v2332 base.V128
	_ = v2332
	var v2339 base.V128
	_ = v2339
	var v2346 base.V128
	_ = v2346
	var v2350 int32
	_ = v2350
	var v2352 int32
	_ = v2352
	var v2354 int32
	_ = v2354
	var v2355 base.V128
	_ = v2355
	var v2382 int32
	_ = v2382
	var v2383 base.V128
	_ = v2383
	var v2384 base.V128
	_ = v2384
	var v2385 base.V128
	_ = v2385
	var v2387 base.V128
	_ = v2387
	var v2391 base.V128
	_ = v2391
	var v2392 base.V128
	_ = v2392
	var v2394 base.V128
	_ = v2394
	var v2396 base.V128
	_ = v2396
	var v2397 base.V128
	_ = v2397
	var v2398 base.V128
	_ = v2398
	var v2400 base.V128
	_ = v2400
	var v2405 base.V128
	_ = v2405
	var v2407 base.V128
	_ = v2407
	var v2408 base.V128
	_ = v2408
	var v2409 base.V128
	_ = v2409
	var v2411 base.V128
	_ = v2411
	var v2417 base.V128
	_ = v2417
	var v2419 int32
	_ = v2419
	var v2421 int32
	_ = v2421
	var v2422 base.V128
	_ = v2422
	var v2424 base.V128
	_ = v2424
	var v2430 base.V128
	_ = v2430
	var v2432 base.V128
	_ = v2432
	var v2434 base.V128
	_ = v2434
	var v2435 base.V128
	_ = v2435
	var v2437 base.V128
	_ = v2437
	var v2442 base.V128
	_ = v2442
	var v2444 base.V128
	_ = v2444
	var v2445 base.V128
	_ = v2445
	var v2447 base.V128
	_ = v2447
	var v2456 base.V128
	_ = v2456
	var v2457 base.V128
	_ = v2457
	var v2459 base.V128
	_ = v2459
	var v2465 base.V128
	_ = v2465
	var v2477 base.V128
	_ = v2477
	var v2478 base.V128
	_ = v2478
	var v2481 base.V128
	_ = v2481
	var v2487 base.V128
	_ = v2487
	var v2499 base.V128
	_ = v2499
	var v2500 base.V128
	_ = v2500
	var v2502 base.V128
	_ = v2502
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
	var v2525 int32
	_ = v2525
	var v2526 base.V128
	_ = v2526
	var v2528 base.V128
	_ = v2528
	var v2534 base.V128
	_ = v2534
	var v2536 base.V128
	_ = v2536
	var v2538 base.V128
	_ = v2538
	var v2539 base.V128
	_ = v2539
	var v2541 base.V128
	_ = v2541
	var v2546 base.V128
	_ = v2546
	var v2548 base.V128
	_ = v2548
	var v2549 base.V128
	_ = v2549
	var v2551 base.V128
	_ = v2551
	var v2560 int32
	_ = v2560
	var v2561 base.V128
	_ = v2561
	var v2563 base.V128
	_ = v2563
	var v2569 base.V128
	_ = v2569
	var v2571 base.V128
	_ = v2571
	var v2573 base.V128
	_ = v2573
	var v2574 base.V128
	_ = v2574
	var v2576 base.V128
	_ = v2576
	var v2581 base.V128
	_ = v2581
	var v2583 base.V128
	_ = v2583
	var v2584 base.V128
	_ = v2584
	var v2586 base.V128
	_ = v2586
	var v2595 base.V128
	_ = v2595
	var v2613 base.V128
	_ = v2613
	var v2632 base.V128
	_ = v2632
	var v2634 base.V128
	_ = v2634
	var v2641 base.V128
	_ = v2641
	var v2648 base.V128
	_ = v2648
	var v2651 int32
	_ = v2651
	var v2652 int32
	_ = v2652
	var v2655 int32
	_ = v2655
	var v2656 int32
	_ = v2656
	var v2658 int32
	_ = v2658
	v10 = int32(0)
	v33 = m.G0
	v35 = v33 - int32(528)
	m.G0 = v35
	base.MemoryFill(m, v35+int32(64), v10, int32(463))
	v162 = int32(1)
	v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v169 = int32(base.Ui32(v163+v164)>>(uint(v162)%32)) + v162
	v172 = int32(base.Ui32(v169+v164) >> (uint(v162) % 32))
	v175 = int32(8)
	v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v181 = int32(base.Ui32(v177*int32(_a_F_UpsampleRgbLinePair_SSE41_0)) >> (uint(v175) % 32))
	v182 = int32(base.Ui32(v172*int32(_a_F_UpsampleRgbLinePair_SSE41_1))>>(uint(v175)%32)) + v181
	v184 = v182 + int32(-17685)
	if base.Ui32(v182) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_2)) {
		v191 = int32(0)
	} else {
		v191 = int32(255)
	}
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_3)) {
		v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
	} else {
		v194 = v191
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v194)
	v196 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v199 = int32(1)
	v202 = int32(base.Ui32(v196+v197)>>(uint(v199)%32)) + v199
	v205 = int32(base.Ui32(v202+v197) >> (uint(v199) % 32))
	v210 = int32(base.Ui32(v205*int32(_a_F_UpsampleRgbLinePair_SSE41_4))>>(uint(int32(8))%32)) + v181
	v212 = v210 + int32(-14234)
	if base.Ui32(v210) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_5)) {
		v219 = int32(0)
	} else {
		v219 = int32(255)
	}
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_3)) {
		v222 = int32(base.Ui32(v212) >> (uint(int32(6)) % 32))
	} else {
		v222 = v219
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v222)
	v226 = int32(8)
	v233 = v181 - (int32(base.Ui32(v172*int32(_a_F_UpsampleRgbLinePair_SSE41_6))>>(uint(v226)%32)) + int32(base.Ui32(v205*int32(_a_F_UpsampleRgbLinePair_SSE41_7))>>(uint(v226)%32)))
	v235 = v233 + int32(_a_F_UpsampleRgbLinePair_SSE41_8)
	if v233 < int32(-8708) {
		v242 = int32(0)
	} else {
		v242 = int32(255)
	}
	if base.Ui32(v235) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_3)) {
		v245 = int32(base.Ui32(v235) >> (uint(int32(6)) % 32))
	} else {
		v245 = v242
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v245)
	if l1 == int32(0) {
	} else {
		v249 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v252 = int32(8)
		v253 = int32(base.Ui32(v249*int32(_a_F_UpsampleRgbLinePair_SSE41_0)) >> (uint(v252) % 32))
		v256 = int32(base.Ui32(v169+v163) >> (uint(int32(1)) % 32))
		v261 = v253 + int32(base.Ui32(v256*int32(_a_F_UpsampleRgbLinePair_SSE41_1))>>(uint(v252)%32))
		v263 = v261 + int32(-17685)
		if base.Ui32(v261) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_2)) {
			v270 = int32(0)
		} else {
			v270 = int32(255)
		}
		if base.Ui32(v263) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_3)) {
			v273 = int32(base.Ui32(v263) >> (uint(int32(6)) % 32))
		} else {
			v273 = v270
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v273)
		v277 = int32(base.Ui32(v202+v196) >> (uint(int32(1)) % 32))
		v282 = v253 + int32(base.Ui32(v277*int32(_a_F_UpsampleRgbLinePair_SSE41_4))>>(uint(int32(8))%32))
		v284 = v282 + int32(-14234)
		if base.Ui32(v282) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_5)) {
			v291 = int32(0)
		} else {
			v291 = int32(255)
		}
		if base.Ui32(v284) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_3)) {
			v294 = int32(base.Ui32(v284) >> (uint(int32(6)) % 32))
		} else {
			v294 = v291
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v294)
		v298 = int32(8)
		v305 = v253 - (int32(base.Ui32(v277*int32(_a_F_UpsampleRgbLinePair_SSE41_7))>>(uint(v298)%32)) + int32(base.Ui32(v256*int32(_a_F_UpsampleRgbLinePair_SSE41_6))>>(uint(v298)%32)))
		v307 = v305 + int32(_a_F_UpsampleRgbLinePair_SSE41_8)
		if v305 < int32(-8708) {
			v314 = int32(0)
		} else {
			v314 = int32(255)
		}
		if base.Ui32(v307) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_SSE41_3)) {
			v317 = int32(base.Ui32(v307) >> (uint(int32(6)) % 32))
		} else {
			v317 = v314
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v317)
	}
	v325 = v35 + int32(96)
	v327 = v35 + int32(64)
	if l8 < int32(34) {
		v1100 = v10
		v1101 = v162
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
			v394 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k0)
			v396 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v384, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v389, v390), v392), v394))
			v402 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v384), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v384), base.Simd_g_v128_and(v392, v389)), v394))
			v403 = base.Simd_g_i8x16_avgr_u(v378, v402)
			v409 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v387), base.Simd_g_v128_and(v392, v390)), v394))
			v410 = base.Simd_g_i8x16_avgr_u(v383, v409)
			v411 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k1)
			v412 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
			v413 = int32(80)
			base.Simd_g_v128_store(m, v327, v413, v412)
			v415 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k4)
			v416 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v417 = int32(64)
			base.Simd_g_v128_store(m, v327, v417, v416)
			v419 = base.Simd_g_i8x16_avgr_u(v381, v409)
			v420 = base.Simd_g_i8x16_avgr_u(v386, v402)
			v422 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
			v423 = int32(16)
			base.Simd_g_v128_store(m, v327, v423, v422)
			v426 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
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
			v464 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
			base.Simd_g_v128_store(m, v327, int32(112), v464)
			v468 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			base.Simd_g_v128_store(m, v327, int32(96), v468)
			v471 = base.Simd_g_i8x16_avgr_u(v434, v461)
			v472 = base.Simd_g_i8x16_avgr_u(v439, v454)
			v474 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
			v475 = int32(48)
			base.Simd_g_v128_store(m, v327, v475, v474)
			v478 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v479 = int32(32)
			base.Simd_g_v128_store(m, v327, v479, v478)
			v481 = l0 + v330 + v359
			v482 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k7)
			v510 = base.Simd_g_v128_load64_zero(m, v481, v423)
			v512 = base.Simd_g_i8x16_shuffle2(v482, v510, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v514 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k8)
			v518 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k9)
			v519 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v512), v514), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v512), v514), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v521 = base.Simd_g_v128_load64_zero(m, v327, v423)
			v523 = base.Simd_g_i8x16_shuffle2(v482, v521, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v524 = base.Simd_g_i32x4_extend_low_i16x8_u(v523)
			v525 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k12)
			v527 = base.Simd_g_i32x4_extend_high_i16x8_u(v523)
			v532 = base.Simd_g_v128_load64_zero(m, v325, v423)
			v534 = base.Simd_g_i8x16_shuffle2(v482, v532, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v535 = base.Simd_g_i32x4_extend_low_i16x8_u(v534)
			v536 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k13)
			v538 = base.Simd_g_i32x4_extend_high_i16x8_u(v534)
			v544 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k14)
			v546 = int32(6)
			v548 = int32(24)
			v549 = base.Simd_g_v128_load64_zero(m, v481, v548)
			v551 = base.Simd_g_i8x16_shuffle2(v482, v549, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v557 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v551), v514), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v551), v514), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v559 = base.Simd_g_v128_load64_zero(m, v327, v548)
			v561 = base.Simd_g_i8x16_shuffle2(v482, v559, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v562 = base.Simd_g_i32x4_extend_low_i16x8_u(v561)
			v564 = base.Simd_g_i32x4_extend_high_i16x8_u(v561)
			v569 = base.Simd_g_v128_load64_zero(m, v325, v548)
			v571 = base.Simd_g_i8x16_shuffle2(v482, v569, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v572 = base.Simd_g_i32x4_extend_low_i16x8_u(v571)
			v574 = base.Simd_g_i32x4_extend_high_i16x8_u(v571)
			v583 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v519, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v524, v525), base.Simd_g_i32x4_mul(v527, v525), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v535, v536), base.Simd_g_i32x4_mul(v538, v536), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v544), v546), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v557, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v562, v525), base.Simd_g_i32x4_mul(v564, v525), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v572, v536), base.Simd_g_i32x4_mul(v574, v536), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v544), v546))
			v584 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k15)
			v586 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k16)
			v592 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k17)
			v604 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v519, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v535, v586), base.Simd_g_i32x4_mul(v538, v586), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v592), v546), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v557, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v572, v586), base.Simd_g_i32x4_mul(v574, v586), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v592), v546))
			v605 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k18)
			v608 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k19)
			v614 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k20)
			v626 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v524, v608), base.Simd_g_i32x4_mul(v527, v608), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v519), v614), v546), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v562, v608), base.Simd_g_i32x4_mul(v564, v608), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v557), v614), v546))
			v627 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k21)
			v629 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v583, v584), base.Simd_g_i8x16_swizzle(v604, v605)), base.Simd_g_i8x16_swizzle(v626, v627))
			base.Simd_g_v128_store(m, v360, v413, v629)
			v632 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k22)
			v634 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k23)
			v637 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k24)
			v639 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v583, v632), base.Simd_g_i8x16_swizzle(v604, v634)), base.Simd_g_i8x16_swizzle(v626, v637))
			base.Simd_g_v128_store(m, v360, v417, v639)
			v642 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k25)
			v644 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k26)
			v647 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k27)
			v649 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v583, v642), base.Simd_g_i8x16_swizzle(v604, v644)), base.Simd_g_i8x16_swizzle(v626, v647))
			base.Simd_g_v128_store(m, v360, v475, v649)
			v653 = base.Simd_g_v128_load64_zero(m, v481, v377)
			v655 = base.Simd_g_i8x16_shuffle2(v482, v653, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v661 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v655), v514), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v655), v514), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v663 = base.Simd_g_v128_load64_zero(m, v327, v377)
			v665 = base.Simd_g_i8x16_shuffle2(v482, v663, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v666 = base.Simd_g_i32x4_extend_low_i16x8_u(v665)
			v668 = base.Simd_g_i32x4_extend_high_i16x8_u(v665)
			v673 = base.Simd_g_v128_load64_zero(m, v325, v377)
			v675 = base.Simd_g_i8x16_shuffle2(v482, v673, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v676 = base.Simd_g_i32x4_extend_low_i16x8_u(v675)
			v678 = base.Simd_g_i32x4_extend_high_i16x8_u(v675)
			v687 = int32(8)
			v688 = base.Simd_g_v128_load64_zero(m, v481, v687)
			v690 = base.Simd_g_i8x16_shuffle2(v482, v688, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v696 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v690), v514), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v690), v514), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v698 = base.Simd_g_v128_load64_zero(m, v327, v687)
			v700 = base.Simd_g_i8x16_shuffle2(v482, v698, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v701 = base.Simd_g_i32x4_extend_low_i16x8_u(v700)
			v703 = base.Simd_g_i32x4_extend_high_i16x8_u(v700)
			v708 = base.Simd_g_v128_load64_zero(m, v325, v687)
			v710 = base.Simd_g_i8x16_shuffle2(v482, v708, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v711 = base.Simd_g_i32x4_extend_low_i16x8_u(v710)
			v713 = base.Simd_g_i32x4_extend_high_i16x8_u(v710)
			v722 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v661, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v666, v525), base.Simd_g_i32x4_mul(v668, v525), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v676, v536), base.Simd_g_i32x4_mul(v678, v536), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v544), v546), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v696, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v701, v525), base.Simd_g_i32x4_mul(v703, v525), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v711, v536), base.Simd_g_i32x4_mul(v713, v536), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v544), v546))
			v740 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v661, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v676, v586), base.Simd_g_i32x4_mul(v678, v586), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v592), v546), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v696, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v711, v586), base.Simd_g_i32x4_mul(v713, v586), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v592), v546))
			v759 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v666, v608), base.Simd_g_i32x4_mul(v668, v608), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v661), v614), v546), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v701, v608), base.Simd_g_i32x4_mul(v703, v608), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v696), v614), v546))
			v761 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v722, v584), base.Simd_g_i8x16_swizzle(v740, v605)), base.Simd_g_i8x16_swizzle(v759, v627))
			base.Simd_g_v128_store(m, v360, v479, v761)
			v768 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v722, v632), base.Simd_g_i8x16_swizzle(v740, v634)), base.Simd_g_i8x16_swizzle(v759, v637))
			base.Simd_g_v128_store(m, v360, v423, v768)
			v775 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v722, v642), base.Simd_g_i8x16_swizzle(v740, v644)), base.Simd_g_i8x16_swizzle(v759, v647))
			base.Simd_g_v128_store(m, v360, v377, v775)
			if l1 == int32(0) {
			} else {
				v780 = l1 + v330 + v359
				v781 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k7)
				v808 = int32(16)
				v809 = base.Simd_g_v128_load64_zero(m, v780, v808)
				v810 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k4)
				v811 = base.Simd_g_i8x16_shuffle2(v781, v809, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v813 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k8)
				v817 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k9)
				v818 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v811), v813), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v811), v813), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
				v820 = base.Simd_g_v128_load64_zero(m, v341, v808)
				v822 = base.Simd_g_i8x16_shuffle2(v781, v820, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v823 = base.Simd_g_i32x4_extend_low_i16x8_u(v822)
				v824 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k12)
				v826 = base.Simd_g_i32x4_extend_high_i16x8_u(v822)
				v831 = base.Simd_g_v128_load64_zero(m, v339, v808)
				v833 = base.Simd_g_i8x16_shuffle2(v781, v831, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v834 = base.Simd_g_i32x4_extend_low_i16x8_u(v833)
				v835 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k13)
				v837 = base.Simd_g_i32x4_extend_high_i16x8_u(v833)
				v843 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k14)
				v845 = int32(6)
				v847 = int32(24)
				v848 = base.Simd_g_v128_load64_zero(m, v780, v847)
				v850 = base.Simd_g_i8x16_shuffle2(v781, v848, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v856 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v850), v813), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v850), v813), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
				v858 = base.Simd_g_v128_load64_zero(m, v341, v847)
				v860 = base.Simd_g_i8x16_shuffle2(v781, v858, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v861 = base.Simd_g_i32x4_extend_low_i16x8_u(v860)
				v863 = base.Simd_g_i32x4_extend_high_i16x8_u(v860)
				v868 = base.Simd_g_v128_load64_zero(m, v339, v847)
				v870 = base.Simd_g_i8x16_shuffle2(v781, v868, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v871 = base.Simd_g_i32x4_extend_low_i16x8_u(v870)
				v873 = base.Simd_g_i32x4_extend_high_i16x8_u(v870)
				v882 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v818, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v823, v824), base.Simd_g_i32x4_mul(v826, v824), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v834, v835), base.Simd_g_i32x4_mul(v837, v835), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v843), v845), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v856, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v861, v824), base.Simd_g_i32x4_mul(v863, v824), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v871, v835), base.Simd_g_i32x4_mul(v873, v835), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v843), v845))
				v883 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k15)
				v885 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k16)
				v891 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k17)
				v903 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v818, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v834, v885), base.Simd_g_i32x4_mul(v837, v885), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v891), v845), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v856, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v871, v885), base.Simd_g_i32x4_mul(v873, v885), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v891), v845))
				v904 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k18)
				v907 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k19)
				v913 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k20)
				v925 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v823, v907), base.Simd_g_i32x4_mul(v826, v907), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v818), v913), v845), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v861, v907), base.Simd_g_i32x4_mul(v863, v907), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v856), v913), v845))
				v926 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k21)
				v928 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v882, v883), base.Simd_g_i8x16_swizzle(v903, v904)), base.Simd_g_i8x16_swizzle(v925, v926))
				base.Simd_g_v128_store(m, v355, int32(80), v928)
				v931 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k22)
				v933 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k23)
				v936 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k24)
				v938 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v882, v931), base.Simd_g_i8x16_swizzle(v903, v933)), base.Simd_g_i8x16_swizzle(v925, v936))
				base.Simd_g_v128_store(m, v355, int32(64), v938)
				v941 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k25)
				v943 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k26)
				v946 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k27)
				v948 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v882, v941), base.Simd_g_i8x16_swizzle(v903, v943)), base.Simd_g_i8x16_swizzle(v925, v946))
				base.Simd_g_v128_store(m, v355, int32(48), v948)
				v951 = int32(0)
				v952 = base.Simd_g_v128_load64_zero(m, v780, v951)
				v954 = base.Simd_g_i8x16_shuffle2(v781, v952, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v960 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v954), v813), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v954), v813), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
				v962 = base.Simd_g_v128_load64_zero(m, v341, v951)
				v964 = base.Simd_g_i8x16_shuffle2(v781, v962, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v965 = base.Simd_g_i32x4_extend_low_i16x8_u(v964)
				v967 = base.Simd_g_i32x4_extend_high_i16x8_u(v964)
				v972 = base.Simd_g_v128_load64_zero(m, v339, v951)
				v974 = base.Simd_g_i8x16_shuffle2(v781, v972, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v975 = base.Simd_g_i32x4_extend_low_i16x8_u(v974)
				v977 = base.Simd_g_i32x4_extend_high_i16x8_u(v974)
				v986 = int32(8)
				v987 = base.Simd_g_v128_load64_zero(m, v780, v986)
				v989 = base.Simd_g_i8x16_shuffle2(v781, v987, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v995 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v989), v813), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v989), v813), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
				v997 = base.Simd_g_v128_load64_zero(m, v341, v986)
				v999 = base.Simd_g_i8x16_shuffle2(v781, v997, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v1000 = base.Simd_g_i32x4_extend_low_i16x8_u(v999)
				v1002 = base.Simd_g_i32x4_extend_high_i16x8_u(v999)
				v1007 = base.Simd_g_v128_load64_zero(m, v339, v986)
				v1009 = base.Simd_g_i8x16_shuffle2(v781, v1007, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
				v1010 = base.Simd_g_i32x4_extend_low_i16x8_u(v1009)
				v1012 = base.Simd_g_i32x4_extend_high_i16x8_u(v1009)
				v1021 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v960, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v965, v824), base.Simd_g_i32x4_mul(v967, v824), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v975, v835), base.Simd_g_i32x4_mul(v977, v835), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v843), v845), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v995, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1000, v824), base.Simd_g_i32x4_mul(v1002, v824), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1010, v835), base.Simd_g_i32x4_mul(v1012, v835), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v843), v845))
				v1039 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v960, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v975, v885), base.Simd_g_i32x4_mul(v977, v885), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v891), v845), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v995, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1010, v885), base.Simd_g_i32x4_mul(v1012, v885), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v891), v845))
				v1058 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v965, v907), base.Simd_g_i32x4_mul(v967, v907), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v960), v913), v845), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1000, v907), base.Simd_g_i32x4_mul(v1002, v907), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v995), v913), v845))
				v1060 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1021, v883), base.Simd_g_i8x16_swizzle(v1039, v904)), base.Simd_g_i8x16_swizzle(v1058, v926))
				base.Simd_g_v128_store(m, v355, int32(32), v1060)
				v1067 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1021, v931), base.Simd_g_i8x16_swizzle(v1039, v933)), base.Simd_g_i8x16_swizzle(v1058, v936))
				base.Simd_g_v128_store(m, v355, v808, v1067)
				v1074 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1021, v941), base.Simd_g_i8x16_swizzle(v1039, v943)), base.Simd_g_i8x16_swizzle(v1058, v946))
				base.Simd_g_v128_store(m, v355, v951, v1074)
			}
			v1077 = int32(96)
			v1082 = v354 + int32(16)
			if v359+int32(66) <= l8 {
				v354 = v1082
				v355 = v355 + v1077
				v359 = v359 + int32(32)
				v360 = v360 + v1077
				continue
			} else {
				break
			}
			break
		}
		v1100 = v1082
		v1101 = v359 + int32(33)
	}
	if l8 < int32(2) {
	} else {
		v1124 = int32(32)
		v1127 = int32(1)
		v1133 = int32(base.Ui32(l8+v1127)>>(uint(v1127)%32)) - int32(base.Ui32(v1101)>>(uint(v1127)%32))
		v1134 = F_memcpy(m, v35+v1124, l2+v1100, v1133)
		mBase = m.M
		v1136 = F_memcpy(m, v35, l4+v1100, v1133)
		mBase = m.M
		v1138 = v1136 + v1124
		v1139 = v1138 + v1133
		v1143 = v1133 + int32(-1)
		v1144 = v1138 + v1143
		v1145 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1144))))
		v1147 = int32(17) - v1133
		if base.Ui32(v1147) < base.Ui32(int32(33)) {
			if v1147 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1139))) = uint8(v1145)
				v1158 = v1139 + v1147
				*(*uint8)(unsafe.Add(mBase, uint32(v1158+int32(-1)))) = uint8(v1145)
				if base.Ui32(v1147) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1139)+2)) = uint8(v1145)
					*(*uint8)(unsafe.Add(mBase, uint32(v1139)+1)) = uint8(v1145)
					*(*uint8)(unsafe.Add(mBase, uint32(v1158+int32(-3)))) = uint8(v1145)
					*(*uint8)(unsafe.Add(mBase, uint32(v1158+int32(-2)))) = uint8(v1145)
					if base.Ui32(v1147) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1139)+3)) = uint8(v1145)
						*(*uint8)(unsafe.Add(mBase, uint32(v1158+int32(-4)))) = uint8(v1145)
						if base.Ui32(v1147) < base.Ui32(int32(9)) {
						} else {
							v1183 = (int32(0) - v1139) & int32(3)
							v1184 = v1139 + v1183
							v1188 = v1145 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1184))) = v1188
							v1192 = (v1147 - v1183) & int32(60)
							v1193 = v1184 + v1192
							*(*int32)(unsafe.Add(mBase, uint32(v1193+int32(-4)))) = v1188
							if base.Ui32(v1192) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1184)+8)) = v1188
								*(*int32)(unsafe.Add(mBase, uint32(v1184)+4)) = v1188
								*(*int32)(unsafe.Add(mBase, uint32(v1193+int32(-8)))) = v1188
								*(*int32)(unsafe.Add(mBase, uint32(v1193+int32(-12)))) = v1188
								if base.Ui32(v1192) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1184)+24)) = v1188
									*(*int32)(unsafe.Add(mBase, uint32(v1184)+20)) = v1188
									*(*int32)(unsafe.Add(mBase, uint32(v1184)+16)) = v1188
									*(*int32)(unsafe.Add(mBase, uint32(v1184)+12)) = v1188
									*(*int32)(unsafe.Add(mBase, uint32(v1193+int32(-16)))) = v1188
									*(*int32)(unsafe.Add(mBase, uint32(v1193+int32(-20)))) = v1188
									*(*int32)(unsafe.Add(mBase, uint32(v1193+int32(-24)))) = v1188
									*(*int32)(unsafe.Add(mBase, uint32(v1193+int32(-28)))) = v1188
									v1228 = v1184&int32(4) | int32(24)
									v1229 = v1192 - v1228
									if base.Ui32(v1229) < base.Ui32(int32(32)) {
									} else {
										v1234 = base.I64_extend_i32_u(v1188) * int64(4294967297)
										v1237 = v1229
										v1238 = v1184 + v1228
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1238)+24)) = v1234
											*(*int64)(unsafe.Add(mBase, uint32(v1238)+16)) = v1234
											*(*int64)(unsafe.Add(mBase, uint32(v1238)+8)) = v1234
											*(*int64)(unsafe.Add(mBase, uint32(v1238))) = v1234
											v1250 = v1237 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1250) {
												v1237 = v1250
												v1238 = v1238 + int32(32)
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
			base.MemoryFill(m, v1139, v1145, v1147)
		}
		v1268 = v1136 + v1133
		v1269 = v1136 + v1143
		v1270 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1269))))
		if base.Ui32(v1147) < base.Ui32(int32(33)) {
			if v1147 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1268))) = uint8(v1270)
				v1281 = v1268 + v1147
				*(*uint8)(unsafe.Add(mBase, uint32(v1281+int32(-1)))) = uint8(v1270)
				if base.Ui32(v1147) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1268)+2)) = uint8(v1270)
					*(*uint8)(unsafe.Add(mBase, uint32(v1268)+1)) = uint8(v1270)
					*(*uint8)(unsafe.Add(mBase, uint32(v1281+int32(-3)))) = uint8(v1270)
					*(*uint8)(unsafe.Add(mBase, uint32(v1281+int32(-2)))) = uint8(v1270)
					if base.Ui32(v1147) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1268)+3)) = uint8(v1270)
						*(*uint8)(unsafe.Add(mBase, uint32(v1281+int32(-4)))) = uint8(v1270)
						if base.Ui32(v1147) < base.Ui32(int32(9)) {
						} else {
							v1306 = (int32(0) - v1268) & int32(3)
							v1307 = v1268 + v1306
							v1311 = v1270 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1307))) = v1311
							v1315 = (v1147 - v1306) & int32(60)
							v1316 = v1307 + v1315
							*(*int32)(unsafe.Add(mBase, uint32(v1316+int32(-4)))) = v1311
							if base.Ui32(v1315) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1307)+8)) = v1311
								*(*int32)(unsafe.Add(mBase, uint32(v1307)+4)) = v1311
								*(*int32)(unsafe.Add(mBase, uint32(v1316+int32(-8)))) = v1311
								*(*int32)(unsafe.Add(mBase, uint32(v1316+int32(-12)))) = v1311
								if base.Ui32(v1315) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1307)+24)) = v1311
									*(*int32)(unsafe.Add(mBase, uint32(v1307)+20)) = v1311
									*(*int32)(unsafe.Add(mBase, uint32(v1307)+16)) = v1311
									*(*int32)(unsafe.Add(mBase, uint32(v1307)+12)) = v1311
									*(*int32)(unsafe.Add(mBase, uint32(v1316+int32(-16)))) = v1311
									*(*int32)(unsafe.Add(mBase, uint32(v1316+int32(-20)))) = v1311
									*(*int32)(unsafe.Add(mBase, uint32(v1316+int32(-24)))) = v1311
									*(*int32)(unsafe.Add(mBase, uint32(v1316+int32(-28)))) = v1311
									v1351 = v1307&int32(4) | int32(24)
									v1352 = v1315 - v1351
									if base.Ui32(v1352) < base.Ui32(int32(32)) {
									} else {
										v1357 = base.I64_extend_i32_u(v1311) * int64(4294967297)
										v1360 = v1352
										v1361 = v1307 + v1351
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1361)+24)) = v1357
											*(*int64)(unsafe.Add(mBase, uint32(v1361)+16)) = v1357
											*(*int64)(unsafe.Add(mBase, uint32(v1361)+8)) = v1357
											*(*int64)(unsafe.Add(mBase, uint32(v1361))) = v1357
											v1373 = v1360 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1373) {
												v1360 = v1373
												v1361 = v1361 + int32(32)
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
			base.MemoryFill(m, v1268, v1270, v1147)
		}
		v1391 = int32(0)
		v1392 = base.Simd_g_v128_load_rng(m, v1136, v1391, int32(0), int32(49))
		v1393 = int32(32)
		v1394 = base.Simd_g_v128_load_nc(m, v1136, v1393)
		v1396 = base.Simd_g_v128_load_nc(m, v1136, int32(1))
		v1397 = base.Simd_g_i8x16_avgr_u(v1394, v1396)
		v1398 = int32(33)
		v1399 = base.Simd_g_v128_load_nc(m, v1136, v1398)
		v1400 = base.Simd_g_i8x16_avgr_u(v1392, v1399)
		v1402 = base.Simd_g_v128_xor(v1396, v1394)
		v1403 = base.Simd_g_v128_xor(v1392, v1399)
		v1405 = base.Simd_g_v128_xor(v1397, v1400)
		v1407 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k0)
		v1409 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1397, v1400), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1402, v1403), v1405), v1407))
		v1415 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1409, v1397), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1409, v1397), base.Simd_g_v128_and(v1405, v1402)), v1407))
		v1416 = base.Simd_g_i8x16_avgr_u(v1392, v1415)
		v1422 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1409, v1400), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1409, v1400), base.Simd_g_v128_and(v1405, v1403)), v1407))
		v1423 = base.Simd_g_i8x16_avgr_u(v1396, v1422)
		v1424 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k1)
		v1425 = base.Simd_g_i8x16_shuffle2(v1416, v1423, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(80), v1425)
		v1428 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k4)
		v1429 = base.Simd_g_i8x16_shuffle2(v1416, v1423, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, int32(64), v1429)
		v1432 = base.Simd_g_i8x16_avgr_u(v1394, v1422)
		v1433 = base.Simd_g_i8x16_avgr_u(v1399, v1415)
		v1435 = base.Simd_g_i8x16_shuffle2(v1432, v1433, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(16), v1435)
		v1439 = base.Simd_g_i8x16_shuffle2(v1432, v1433, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, v1391, v1439)
		v1445 = F_memcpy(m, v1136+v1393, l3+v1100, v1133)
		mBase = m.M
		v1447 = F_memcpy(m, v1136, l5+v1100, v1133)
		mBase = m.M
		v1448 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1144))))
		if base.Ui32(v1147) < base.Ui32(v1398) {
			if v1147 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1139))) = uint8(v1448)
				v1459 = v1139 + v1147
				*(*uint8)(unsafe.Add(mBase, uint32(v1459+int32(-1)))) = uint8(v1448)
				if base.Ui32(v1147) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1139)+2)) = uint8(v1448)
					*(*uint8)(unsafe.Add(mBase, uint32(v1139)+1)) = uint8(v1448)
					*(*uint8)(unsafe.Add(mBase, uint32(v1459+int32(-3)))) = uint8(v1448)
					*(*uint8)(unsafe.Add(mBase, uint32(v1459+int32(-2)))) = uint8(v1448)
					if base.Ui32(v1147) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1139)+3)) = uint8(v1448)
						*(*uint8)(unsafe.Add(mBase, uint32(v1459+int32(-4)))) = uint8(v1448)
						if base.Ui32(v1147) < base.Ui32(int32(9)) {
						} else {
							v1484 = (int32(0) - v1139) & int32(3)
							v1485 = v1139 + v1484
							v1489 = v1448 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1485))) = v1489
							v1493 = (v1147 - v1484) & int32(60)
							v1494 = v1485 + v1493
							*(*int32)(unsafe.Add(mBase, uint32(v1494+int32(-4)))) = v1489
							if base.Ui32(v1493) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1485)+8)) = v1489
								*(*int32)(unsafe.Add(mBase, uint32(v1485)+4)) = v1489
								*(*int32)(unsafe.Add(mBase, uint32(v1494+int32(-8)))) = v1489
								*(*int32)(unsafe.Add(mBase, uint32(v1494+int32(-12)))) = v1489
								if base.Ui32(v1493) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1485)+24)) = v1489
									*(*int32)(unsafe.Add(mBase, uint32(v1485)+20)) = v1489
									*(*int32)(unsafe.Add(mBase, uint32(v1485)+16)) = v1489
									*(*int32)(unsafe.Add(mBase, uint32(v1485)+12)) = v1489
									*(*int32)(unsafe.Add(mBase, uint32(v1494+int32(-16)))) = v1489
									*(*int32)(unsafe.Add(mBase, uint32(v1494+int32(-20)))) = v1489
									*(*int32)(unsafe.Add(mBase, uint32(v1494+int32(-24)))) = v1489
									*(*int32)(unsafe.Add(mBase, uint32(v1494+int32(-28)))) = v1489
									v1529 = v1485&int32(4) | int32(24)
									v1530 = v1493 - v1529
									if base.Ui32(v1530) < base.Ui32(int32(32)) {
									} else {
										v1535 = base.I64_extend_i32_u(v1489) * int64(4294967297)
										v1538 = v1530
										v1539 = v1485 + v1529
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1539)+24)) = v1535
											*(*int64)(unsafe.Add(mBase, uint32(v1539)+16)) = v1535
											*(*int64)(unsafe.Add(mBase, uint32(v1539)+8)) = v1535
											*(*int64)(unsafe.Add(mBase, uint32(v1539))) = v1535
											v1551 = v1538 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1551) {
												v1538 = v1551
												v1539 = v1539 + int32(32)
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
			base.MemoryFill(m, v1139, v1448, v1147)
		}
		v1569 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1269))))
		if base.Ui32(v1147) < base.Ui32(int32(33)) {
			if v1147 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1268))) = uint8(v1569)
				v1580 = v1268 + v1147
				*(*uint8)(unsafe.Add(mBase, uint32(v1580+int32(-1)))) = uint8(v1569)
				if base.Ui32(v1147) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1268)+2)) = uint8(v1569)
					*(*uint8)(unsafe.Add(mBase, uint32(v1268)+1)) = uint8(v1569)
					*(*uint8)(unsafe.Add(mBase, uint32(v1580+int32(-3)))) = uint8(v1569)
					*(*uint8)(unsafe.Add(mBase, uint32(v1580+int32(-2)))) = uint8(v1569)
					if base.Ui32(v1147) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1268)+3)) = uint8(v1569)
						*(*uint8)(unsafe.Add(mBase, uint32(v1580+int32(-4)))) = uint8(v1569)
						if base.Ui32(v1147) < base.Ui32(int32(9)) {
						} else {
							v1605 = (int32(0) - v1268) & int32(3)
							v1606 = v1268 + v1605
							v1610 = v1569 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1606))) = v1610
							v1614 = (v1147 - v1605) & int32(60)
							v1615 = v1606 + v1614
							*(*int32)(unsafe.Add(mBase, uint32(v1615+int32(-4)))) = v1610
							if base.Ui32(v1614) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1606)+8)) = v1610
								*(*int32)(unsafe.Add(mBase, uint32(v1606)+4)) = v1610
								*(*int32)(unsafe.Add(mBase, uint32(v1615+int32(-8)))) = v1610
								*(*int32)(unsafe.Add(mBase, uint32(v1615+int32(-12)))) = v1610
								if base.Ui32(v1614) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1606)+24)) = v1610
									*(*int32)(unsafe.Add(mBase, uint32(v1606)+20)) = v1610
									*(*int32)(unsafe.Add(mBase, uint32(v1606)+16)) = v1610
									*(*int32)(unsafe.Add(mBase, uint32(v1606)+12)) = v1610
									*(*int32)(unsafe.Add(mBase, uint32(v1615+int32(-16)))) = v1610
									*(*int32)(unsafe.Add(mBase, uint32(v1615+int32(-20)))) = v1610
									*(*int32)(unsafe.Add(mBase, uint32(v1615+int32(-24)))) = v1610
									*(*int32)(unsafe.Add(mBase, uint32(v1615+int32(-28)))) = v1610
									v1650 = v1606&int32(4) | int32(24)
									v1651 = v1614 - v1650
									if base.Ui32(v1651) < base.Ui32(int32(32)) {
									} else {
										v1656 = base.I64_extend_i32_u(v1610) * int64(4294967297)
										v1659 = v1651
										v1660 = v1606 + v1650
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1660)+24)) = v1656
											*(*int64)(unsafe.Add(mBase, uint32(v1660)+16)) = v1656
											*(*int64)(unsafe.Add(mBase, uint32(v1660)+8)) = v1656
											*(*int64)(unsafe.Add(mBase, uint32(v1660))) = v1656
											v1672 = v1659 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1672) {
												v1659 = v1672
												v1660 = v1660 + int32(32)
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
			base.MemoryFill(m, v1268, v1569, v1147)
		}
		v1691 = base.Simd_g_v128_load_rng(m, v1447, int32(0), int32(0), int32(49))
		v1692 = int32(32)
		v1693 = base.Simd_g_v128_load_nc(m, v1447, v1692)
		v1695 = base.Simd_g_v128_load_nc(m, v1447, int32(1))
		v1696 = base.Simd_g_i8x16_avgr_u(v1693, v1695)
		v1698 = base.Simd_g_v128_load_nc(m, v1447, int32(33))
		v1699 = base.Simd_g_i8x16_avgr_u(v1691, v1698)
		v1701 = base.Simd_g_v128_xor(v1695, v1693)
		v1702 = base.Simd_g_v128_xor(v1691, v1698)
		v1704 = base.Simd_g_v128_xor(v1696, v1699)
		v1707 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1696, v1699), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1701, v1702), v1704), v1407))
		v1713 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1707, v1696), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1707, v1696), base.Simd_g_v128_and(v1704, v1701)), v1407))
		v1714 = base.Simd_g_i8x16_avgr_u(v1691, v1713)
		v1720 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1707, v1699), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1707, v1699), base.Simd_g_v128_and(v1704, v1702)), v1407))
		v1721 = base.Simd_g_i8x16_avgr_u(v1695, v1720)
		v1722 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k1)
		v1723 = base.Simd_g_i8x16_shuffle2(v1714, v1721, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(112), v1723)
		v1726 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k4)
		v1727 = base.Simd_g_i8x16_shuffle2(v1714, v1721, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, int32(96), v1727)
		v1730 = base.Simd_g_i8x16_avgr_u(v1693, v1720)
		v1731 = base.Simd_g_i8x16_avgr_u(v1698, v1713)
		v1733 = base.Simd_g_i8x16_shuffle2(v1730, v1731, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(48), v1733)
		v1737 = base.Simd_g_i8x16_shuffle2(v1730, v1731, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, v1692, v1737)
		v1743 = l8 - v1101
		v1744 = F_memcpy(m, v35+int32(448), l0+v1101, v1743)
		mBase = m.M
		v1746 = v35 + int32(192)
		if l1 != 0 {
			v2052 = F_memcpy(m, v35+int32(480), l1+v1101, v1743)
			mBase = m.M
			v2053 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k7)
			v2080 = int32(16)
			v2081 = base.Simd_g_v128_load64_zero(m, v1744, v2080)
			v2082 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k4)
			v2083 = base.Simd_g_i8x16_shuffle2(v2053, v2081, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2085 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k8)
			v2089 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k9)
			v2090 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2083), v2085), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2083), v2085), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2092 = base.Simd_g_v128_load64_zero(m, v327, v2080)
			v2094 = base.Simd_g_i8x16_shuffle2(v2053, v2092, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2095 = base.Simd_g_i32x4_extend_low_i16x8_u(v2094)
			v2096 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k12)
			v2098 = base.Simd_g_i32x4_extend_high_i16x8_u(v2094)
			v2103 = base.Simd_g_v128_load64_zero(m, v325, v2080)
			v2105 = base.Simd_g_i8x16_shuffle2(v2053, v2103, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2106 = base.Simd_g_i32x4_extend_low_i16x8_u(v2105)
			v2107 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k13)
			v2109 = base.Simd_g_i32x4_extend_high_i16x8_u(v2105)
			v2115 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k14)
			v2117 = int32(6)
			v2119 = int32(24)
			v2120 = base.Simd_g_v128_load64_zero(m, v1744, v2119)
			v2122 = base.Simd_g_i8x16_shuffle2(v2053, v2120, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2128 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2122), v2085), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2122), v2085), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2130 = base.Simd_g_v128_load64_zero(m, v327, v2119)
			v2132 = base.Simd_g_i8x16_shuffle2(v2053, v2130, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2133 = base.Simd_g_i32x4_extend_low_i16x8_u(v2132)
			v2135 = base.Simd_g_i32x4_extend_high_i16x8_u(v2132)
			v2140 = base.Simd_g_v128_load64_zero(m, v325, v2119)
			v2142 = base.Simd_g_i8x16_shuffle2(v2053, v2140, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2143 = base.Simd_g_i32x4_extend_low_i16x8_u(v2142)
			v2145 = base.Simd_g_i32x4_extend_high_i16x8_u(v2142)
			v2154 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2090, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2095, v2096), base.Simd_g_i32x4_mul(v2098, v2096), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2106, v2107), base.Simd_g_i32x4_mul(v2109, v2107), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2115), v2117), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2128, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2133, v2096), base.Simd_g_i32x4_mul(v2135, v2096), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2143, v2107), base.Simd_g_i32x4_mul(v2145, v2107), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2115), v2117))
			v2155 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k15)
			v2157 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k16)
			v2163 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k17)
			v2175 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2090, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2106, v2157), base.Simd_g_i32x4_mul(v2109, v2157), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2163), v2117), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2128, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2143, v2157), base.Simd_g_i32x4_mul(v2145, v2157), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2163), v2117))
			v2176 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k18)
			v2179 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k19)
			v2185 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k20)
			v2197 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2095, v2179), base.Simd_g_i32x4_mul(v2098, v2179), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2090), v2185), v2117), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2133, v2179), base.Simd_g_i32x4_mul(v2135, v2179), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2128), v2185), v2117))
			v2198 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k21)
			v2200 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2154, v2155), base.Simd_g_i8x16_swizzle(v2175, v2176)), base.Simd_g_i8x16_swizzle(v2197, v2198))
			base.Simd_g_v128_store(m, v1746, int32(80), v2200)
			v2203 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k22)
			v2205 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k23)
			v2208 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k24)
			v2210 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2154, v2203), base.Simd_g_i8x16_swizzle(v2175, v2205)), base.Simd_g_i8x16_swizzle(v2197, v2208))
			base.Simd_g_v128_store(m, v1746, int32(64), v2210)
			v2213 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k25)
			v2215 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k26)
			v2218 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k27)
			v2220 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2154, v2213), base.Simd_g_i8x16_swizzle(v2175, v2215)), base.Simd_g_i8x16_swizzle(v2197, v2218))
			base.Simd_g_v128_store(m, v1746, int32(48), v2220)
			v2223 = int32(0)
			v2224 = base.Simd_g_v128_load64_zero(m, v1744, v2223)
			v2226 = base.Simd_g_i8x16_shuffle2(v2053, v2224, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2232 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2226), v2085), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2226), v2085), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2234 = base.Simd_g_v128_load64_zero(m, v327, v2223)
			v2236 = base.Simd_g_i8x16_shuffle2(v2053, v2234, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2237 = base.Simd_g_i32x4_extend_low_i16x8_u(v2236)
			v2239 = base.Simd_g_i32x4_extend_high_i16x8_u(v2236)
			v2244 = base.Simd_g_v128_load64_zero(m, v325, v2223)
			v2246 = base.Simd_g_i8x16_shuffle2(v2053, v2244, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2247 = base.Simd_g_i32x4_extend_low_i16x8_u(v2246)
			v2249 = base.Simd_g_i32x4_extend_high_i16x8_u(v2246)
			v2258 = int32(8)
			v2259 = base.Simd_g_v128_load64_zero(m, v1744, v2258)
			v2261 = base.Simd_g_i8x16_shuffle2(v2053, v2259, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2267 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2261), v2085), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2261), v2085), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2269 = base.Simd_g_v128_load64_zero(m, v327, v2258)
			v2271 = base.Simd_g_i8x16_shuffle2(v2053, v2269, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2272 = base.Simd_g_i32x4_extend_low_i16x8_u(v2271)
			v2274 = base.Simd_g_i32x4_extend_high_i16x8_u(v2271)
			v2279 = base.Simd_g_v128_load64_zero(m, v325, v2258)
			v2281 = base.Simd_g_i8x16_shuffle2(v2053, v2279, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2282 = base.Simd_g_i32x4_extend_low_i16x8_u(v2281)
			v2284 = base.Simd_g_i32x4_extend_high_i16x8_u(v2281)
			v2293 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2232, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2237, v2096), base.Simd_g_i32x4_mul(v2239, v2096), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2247, v2107), base.Simd_g_i32x4_mul(v2249, v2107), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2115), v2117), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2267, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2272, v2096), base.Simd_g_i32x4_mul(v2274, v2096), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2282, v2107), base.Simd_g_i32x4_mul(v2284, v2107), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2115), v2117))
			v2311 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2232, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2247, v2157), base.Simd_g_i32x4_mul(v2249, v2157), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2163), v2117), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2267, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2282, v2157), base.Simd_g_i32x4_mul(v2284, v2157), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2163), v2117))
			v2330 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2237, v2179), base.Simd_g_i32x4_mul(v2239, v2179), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2232), v2185), v2117), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2272, v2179), base.Simd_g_i32x4_mul(v2274, v2179), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2267), v2185), v2117))
			v2332 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2293, v2155), base.Simd_g_i8x16_swizzle(v2311, v2176)), base.Simd_g_i8x16_swizzle(v2330, v2198))
			base.Simd_g_v128_store(m, v1746, int32(32), v2332)
			v2339 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2293, v2203), base.Simd_g_i8x16_swizzle(v2311, v2205)), base.Simd_g_i8x16_swizzle(v2330, v2208))
			base.Simd_g_v128_store(m, v1746, v2080, v2339)
			v2346 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2293, v2213), base.Simd_g_i8x16_swizzle(v2311, v2215)), base.Simd_g_i8x16_swizzle(v2330, v2218))
			base.Simd_g_v128_store(m, v1746, v2223, v2346)
			v2350 = v35 + int32(128)
			v2352 = v35 + int32(160)
			v2354 = v35 + int32(320)
			v2355 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k7)
			v2382 = int32(16)
			v2383 = base.Simd_g_v128_load64_zero(m, v2052, v2382)
			v2384 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k4)
			v2385 = base.Simd_g_i8x16_shuffle2(v2355, v2383, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2387 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k8)
			v2391 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k9)
			v2392 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2385), v2387), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2385), v2387), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2394 = base.Simd_g_v128_load64_zero(m, v2350, v2382)
			v2396 = base.Simd_g_i8x16_shuffle2(v2355, v2394, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2397 = base.Simd_g_i32x4_extend_low_i16x8_u(v2396)
			v2398 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k12)
			v2400 = base.Simd_g_i32x4_extend_high_i16x8_u(v2396)
			v2405 = base.Simd_g_v128_load64_zero(m, v2352, v2382)
			v2407 = base.Simd_g_i8x16_shuffle2(v2355, v2405, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2408 = base.Simd_g_i32x4_extend_low_i16x8_u(v2407)
			v2409 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k13)
			v2411 = base.Simd_g_i32x4_extend_high_i16x8_u(v2407)
			v2417 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k14)
			v2419 = int32(6)
			v2421 = int32(24)
			v2422 = base.Simd_g_v128_load64_zero(m, v2052, v2421)
			v2424 = base.Simd_g_i8x16_shuffle2(v2355, v2422, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2430 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2424), v2387), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2424), v2387), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2432 = base.Simd_g_v128_load64_zero(m, v2350, v2421)
			v2434 = base.Simd_g_i8x16_shuffle2(v2355, v2432, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2435 = base.Simd_g_i32x4_extend_low_i16x8_u(v2434)
			v2437 = base.Simd_g_i32x4_extend_high_i16x8_u(v2434)
			v2442 = base.Simd_g_v128_load64_zero(m, v2352, v2421)
			v2444 = base.Simd_g_i8x16_shuffle2(v2355, v2442, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2445 = base.Simd_g_i32x4_extend_low_i16x8_u(v2444)
			v2447 = base.Simd_g_i32x4_extend_high_i16x8_u(v2444)
			v2456 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2392, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2397, v2398), base.Simd_g_i32x4_mul(v2400, v2398), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2408, v2409), base.Simd_g_i32x4_mul(v2411, v2409), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2417), v2419), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2430, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2435, v2398), base.Simd_g_i32x4_mul(v2437, v2398), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2445, v2409), base.Simd_g_i32x4_mul(v2447, v2409), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2417), v2419))
			v2457 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k15)
			v2459 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k16)
			v2465 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k17)
			v2477 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2392, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2408, v2459), base.Simd_g_i32x4_mul(v2411, v2459), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2465), v2419), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2430, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2445, v2459), base.Simd_g_i32x4_mul(v2447, v2459), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2465), v2419))
			v2478 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k18)
			v2481 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k19)
			v2487 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k20)
			v2499 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2397, v2481), base.Simd_g_i32x4_mul(v2400, v2481), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2392), v2487), v2419), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2435, v2481), base.Simd_g_i32x4_mul(v2437, v2481), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2430), v2487), v2419))
			v2500 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k21)
			v2502 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2456, v2457), base.Simd_g_i8x16_swizzle(v2477, v2478)), base.Simd_g_i8x16_swizzle(v2499, v2500))
			base.Simd_g_v128_store(m, v2354, int32(80), v2502)
			v2505 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k22)
			v2507 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k23)
			v2510 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k24)
			v2512 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2456, v2505), base.Simd_g_i8x16_swizzle(v2477, v2507)), base.Simd_g_i8x16_swizzle(v2499, v2510))
			base.Simd_g_v128_store(m, v2354, int32(64), v2512)
			v2515 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k25)
			v2517 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k26)
			v2520 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k27)
			v2522 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2456, v2515), base.Simd_g_i8x16_swizzle(v2477, v2517)), base.Simd_g_i8x16_swizzle(v2499, v2520))
			base.Simd_g_v128_store(m, v2354, int32(48), v2522)
			v2525 = int32(0)
			v2526 = base.Simd_g_v128_load64_zero(m, v2052, v2525)
			v2528 = base.Simd_g_i8x16_shuffle2(v2355, v2526, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2534 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2528), v2387), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2528), v2387), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2536 = base.Simd_g_v128_load64_zero(m, v2350, v2525)
			v2538 = base.Simd_g_i8x16_shuffle2(v2355, v2536, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2539 = base.Simd_g_i32x4_extend_low_i16x8_u(v2538)
			v2541 = base.Simd_g_i32x4_extend_high_i16x8_u(v2538)
			v2546 = base.Simd_g_v128_load64_zero(m, v2352, v2525)
			v2548 = base.Simd_g_i8x16_shuffle2(v2355, v2546, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2549 = base.Simd_g_i32x4_extend_low_i16x8_u(v2548)
			v2551 = base.Simd_g_i32x4_extend_high_i16x8_u(v2548)
			v2560 = int32(8)
			v2561 = base.Simd_g_v128_load64_zero(m, v2052, v2560)
			v2563 = base.Simd_g_i8x16_shuffle2(v2355, v2561, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2569 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2563), v2387), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2563), v2387), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v2571 = base.Simd_g_v128_load64_zero(m, v2350, v2560)
			v2573 = base.Simd_g_i8x16_shuffle2(v2355, v2571, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2574 = base.Simd_g_i32x4_extend_low_i16x8_u(v2573)
			v2576 = base.Simd_g_i32x4_extend_high_i16x8_u(v2573)
			v2581 = base.Simd_g_v128_load64_zero(m, v2352, v2560)
			v2583 = base.Simd_g_i8x16_shuffle2(v2355, v2581, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v2584 = base.Simd_g_i32x4_extend_low_i16x8_u(v2583)
			v2586 = base.Simd_g_i32x4_extend_high_i16x8_u(v2583)
			v2595 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2534, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2539, v2398), base.Simd_g_i32x4_mul(v2541, v2398), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2549, v2409), base.Simd_g_i32x4_mul(v2551, v2409), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2417), v2419), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2569, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2574, v2398), base.Simd_g_i32x4_mul(v2576, v2398), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2584, v2409), base.Simd_g_i32x4_mul(v2586, v2409), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v2417), v2419))
			v2613 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2534, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2549, v2459), base.Simd_g_i32x4_mul(v2551, v2459), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2465), v2419), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2569, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2584, v2459), base.Simd_g_i32x4_mul(v2586, v2459), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v2465), v2419))
			v2632 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2539, v2481), base.Simd_g_i32x4_mul(v2541, v2481), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2534), v2487), v2419), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2574, v2481), base.Simd_g_i32x4_mul(v2576, v2481), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v2569), v2487), v2419))
			v2634 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2595, v2457), base.Simd_g_i8x16_swizzle(v2613, v2478)), base.Simd_g_i8x16_swizzle(v2632, v2500))
			base.Simd_g_v128_store(m, v2354, int32(32), v2634)
			v2641 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2595, v2505), base.Simd_g_i8x16_swizzle(v2613, v2507)), base.Simd_g_i8x16_swizzle(v2632, v2510))
			base.Simd_g_v128_store(m, v2354, v2382, v2641)
			v2648 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2595, v2515), base.Simd_g_i8x16_swizzle(v2613, v2517)), base.Simd_g_i8x16_swizzle(v2632, v2520))
			base.Simd_g_v128_store(m, v2354, v2525, v2648)
			v2651 = int32(3)
			v2652 = v1101 * v2651
			v2655 = v1743 * v2651
			v2656 = F_memcpy(m, l6+v2652, v1746, v2655)
			mBase = m.M
			v2658 = F_memcpy(m, l7+v2652, v2354, v2655)
			mBase = m.M
		} else {
			v1747 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k7)
			v1774 = int32(16)
			v1775 = base.Simd_g_v128_load64_zero(m, v1744, v1774)
			v1776 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k4)
			v1777 = base.Simd_g_i8x16_shuffle2(v1747, v1775, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1779 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k8)
			v1783 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k9)
			v1784 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1777), v1779), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1777), v1779), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v1786 = base.Simd_g_v128_load64_zero(m, v327, v1774)
			v1788 = base.Simd_g_i8x16_shuffle2(v1747, v1786, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1789 = base.Simd_g_i32x4_extend_low_i16x8_u(v1788)
			v1790 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k12)
			v1792 = base.Simd_g_i32x4_extend_high_i16x8_u(v1788)
			v1797 = base.Simd_g_v128_load64_zero(m, v325, v1774)
			v1799 = base.Simd_g_i8x16_shuffle2(v1747, v1797, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1800 = base.Simd_g_i32x4_extend_low_i16x8_u(v1799)
			v1801 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k13)
			v1803 = base.Simd_g_i32x4_extend_high_i16x8_u(v1799)
			v1809 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k14)
			v1811 = int32(6)
			v1813 = int32(24)
			v1814 = base.Simd_g_v128_load64_zero(m, v1744, v1813)
			v1816 = base.Simd_g_i8x16_shuffle2(v1747, v1814, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1822 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1816), v1779), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1816), v1779), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v1824 = base.Simd_g_v128_load64_zero(m, v327, v1813)
			v1826 = base.Simd_g_i8x16_shuffle2(v1747, v1824, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1827 = base.Simd_g_i32x4_extend_low_i16x8_u(v1826)
			v1829 = base.Simd_g_i32x4_extend_high_i16x8_u(v1826)
			v1834 = base.Simd_g_v128_load64_zero(m, v325, v1813)
			v1836 = base.Simd_g_i8x16_shuffle2(v1747, v1834, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1837 = base.Simd_g_i32x4_extend_low_i16x8_u(v1836)
			v1839 = base.Simd_g_i32x4_extend_high_i16x8_u(v1836)
			v1848 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1784, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1789, v1790), base.Simd_g_i32x4_mul(v1792, v1790), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1800, v1801), base.Simd_g_i32x4_mul(v1803, v1801), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v1809), v1811), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1822, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1827, v1790), base.Simd_g_i32x4_mul(v1829, v1790), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1837, v1801), base.Simd_g_i32x4_mul(v1839, v1801), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v1809), v1811))
			v1849 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k15)
			v1851 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k16)
			v1857 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k17)
			v1869 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1784, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1800, v1851), base.Simd_g_i32x4_mul(v1803, v1851), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v1857), v1811), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1822, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1837, v1851), base.Simd_g_i32x4_mul(v1839, v1851), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v1857), v1811))
			v1870 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k18)
			v1873 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k19)
			v1879 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k20)
			v1891 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1789, v1873), base.Simd_g_i32x4_mul(v1792, v1873), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v1784), v1879), v1811), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1827, v1873), base.Simd_g_i32x4_mul(v1829, v1873), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v1822), v1879), v1811))
			v1892 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k21)
			v1894 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1848, v1849), base.Simd_g_i8x16_swizzle(v1869, v1870)), base.Simd_g_i8x16_swizzle(v1891, v1892))
			base.Simd_g_v128_store(m, v1746, int32(80), v1894)
			v1897 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k22)
			v1899 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k23)
			v1902 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k24)
			v1904 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1848, v1897), base.Simd_g_i8x16_swizzle(v1869, v1899)), base.Simd_g_i8x16_swizzle(v1891, v1902))
			base.Simd_g_v128_store(m, v1746, int32(64), v1904)
			v1907 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k25)
			v1909 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k26)
			v1912 = base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k27)
			v1914 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1848, v1907), base.Simd_g_i8x16_swizzle(v1869, v1909)), base.Simd_g_i8x16_swizzle(v1891, v1912))
			base.Simd_g_v128_store(m, v1746, int32(48), v1914)
			v1917 = int32(0)
			v1918 = base.Simd_g_v128_load64_zero(m, v1744, v1917)
			v1920 = base.Simd_g_i8x16_shuffle2(v1747, v1918, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1926 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1920), v1779), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1920), v1779), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v1928 = base.Simd_g_v128_load64_zero(m, v327, v1917)
			v1930 = base.Simd_g_i8x16_shuffle2(v1747, v1928, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1931 = base.Simd_g_i32x4_extend_low_i16x8_u(v1930)
			v1933 = base.Simd_g_i32x4_extend_high_i16x8_u(v1930)
			v1938 = base.Simd_g_v128_load64_zero(m, v325, v1917)
			v1940 = base.Simd_g_i8x16_shuffle2(v1747, v1938, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1941 = base.Simd_g_i32x4_extend_low_i16x8_u(v1940)
			v1943 = base.Simd_g_i32x4_extend_high_i16x8_u(v1940)
			v1952 = int32(8)
			v1953 = base.Simd_g_v128_load64_zero(m, v1744, v1952)
			v1955 = base.Simd_g_i8x16_shuffle2(v1747, v1953, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1961 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1955), v1779), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1955), v1779), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))
			v1963 = base.Simd_g_v128_load64_zero(m, v327, v1952)
			v1965 = base.Simd_g_i8x16_shuffle2(v1747, v1963, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1966 = base.Simd_g_i32x4_extend_low_i16x8_u(v1965)
			v1968 = base.Simd_g_i32x4_extend_high_i16x8_u(v1965)
			v1973 = base.Simd_g_v128_load64_zero(m, v325, v1952)
			v1975 = base.Simd_g_i8x16_shuffle2(v1747, v1973, base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k6))
			v1976 = base.Simd_g_i32x4_extend_low_i16x8_u(v1975)
			v1978 = base.Simd_g_i32x4_extend_high_i16x8_u(v1975)
			v1987 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1926, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1931, v1790), base.Simd_g_i32x4_mul(v1933, v1790), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1941, v1801), base.Simd_g_i32x4_mul(v1943, v1801), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v1809), v1811), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1961, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1966, v1790), base.Simd_g_i32x4_mul(v1968, v1790), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1976, v1801), base.Simd_g_i32x4_mul(v1978, v1801), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)))), v1809), v1811))
			v2005 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1926, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1941, v1851), base.Simd_g_i32x4_mul(v1943, v1851), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v1857), v1811), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1961, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1976, v1851), base.Simd_g_i32x4_mul(v1978, v1851), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11))), v1857), v1811))
			v2024 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1931, v1873), base.Simd_g_i32x4_mul(v1933, v1873), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v1926), v1879), v1811), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1966, v1873), base.Simd_g_i32x4_mul(v1968, v1873), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleRgbLinePair_SSE41__k11)), v1961), v1879), v1811))
			v2026 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1987, v1849), base.Simd_g_i8x16_swizzle(v2005, v1870)), base.Simd_g_i8x16_swizzle(v2024, v1892))
			base.Simd_g_v128_store(m, v1746, int32(32), v2026)
			v2033 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1987, v1897), base.Simd_g_i8x16_swizzle(v2005, v1899)), base.Simd_g_i8x16_swizzle(v2024, v1902))
			base.Simd_g_v128_store(m, v1746, v1774, v2033)
			v2040 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1987, v1907), base.Simd_g_i8x16_swizzle(v2005, v1909)), base.Simd_g_i8x16_swizzle(v2024, v1912))
			base.Simd_g_v128_store(m, v1746, v1917, v2040)
			v2043 = int32(3)
			v2048 = F_memcpy(m, l6+v1101*v2043, v1746, v1743*v2043)
			mBase = m.M
		}
	}
	m.G0 = v35 + int32(528)
	return
}

var F_UpsampleRgbLinePair_SSE41__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleRgbLinePair_SSE41__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleRgbLinePair_SSE41__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleRgbLinePair_SSE41__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleRgbLinePair_SSE41__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleRgbLinePair_SSE41__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleRgbLinePair_SSE41__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleRgbLinePair_SSE41__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleRgbLinePair_SSE41__k8 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleRgbLinePair_SSE41__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleRgbLinePair_SSE41__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleRgbLinePair_SSE41__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleRgbLinePair_SSE41__k12 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleRgbLinePair_SSE41__k13 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleRgbLinePair_SSE41__k14 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_UpsampleRgbLinePair_SSE41__k15 = [2]uint64{0x8f8f0c8f8f0b8f8f, 0x8f0f8f8f0e8f8f0d}
var F_UpsampleRgbLinePair_SSE41__k16 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleRgbLinePair_SSE41__k17 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleRgbLinePair_SSE41__k18 = [2]uint64{0xd8f8f0c8f8f0b8f, 0x8f8f0f8f8f0e8f8f}
var F_UpsampleRgbLinePair_SSE41__k19 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleRgbLinePair_SSE41__k20 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleRgbLinePair_SSE41__k21 = [2]uint64{0x8f0c8f8f0b8f8f0a, 0xf8f8f0e8f8f0d8f}
var F_UpsampleRgbLinePair_SSE41__k22 = [2]uint64{0x8f078f8f068f8f05, 0xa8f8f098f8f088f}
var F_UpsampleRgbLinePair_SSE41__k23 = [2]uint64{0x8f8f078f8f068f8f, 0x8f0a8f8f098f8f08}
var F_UpsampleRgbLinePair_SSE41__k24 = [2]uint64{0x78f8f068f8f058f, 0x8f8f098f8f088f8f}
var F_UpsampleRgbLinePair_SSE41__k25 = [2]uint64{0x28f8f018f8f008f, 0x8f8f048f8f038f8f}
var F_UpsampleRgbLinePair_SSE41__k26 = [2]uint64{0x8f028f8f018f8f00, 0x58f8f048f8f038f}
var F_UpsampleRgbLinePair_SSE41__k27 = [2]uint64{0x8f8f018f8f008f8f, 0x8f048f8f038f8f02}

func F_UpsampleRgba4444LinePair_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
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
	var v229 int32
	_ = v229
	var v236 int32
	_ = v236
	var v238 int32
	_ = v238
	var v245 int32
	_ = v245
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
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
	var v279 int32
	_ = v279
	var v283 int32
	_ = v283
	var v288 int32
	_ = v288
	var v290 int32
	_ = v290
	var v297 int32
	_ = v297
	var v300 int32
	_ = v300
	var v305 int32
	_ = v305
	var v312 int32
	_ = v312
	var v314 int32
	_ = v314
	var v321 int32
	_ = v321
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v333 int32
	_ = v333
	var v335 int32
	_ = v335
	var v338 int32
	_ = v338
	var v342 int32
	_ = v342
	var v350 int32
	_ = v350
	var v362 int32
	_ = v362
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v386 base.V128
	_ = v386
	var v387 int32
	_ = v387
	var v389 base.V128
	_ = v389
	var v390 int32
	_ = v390
	var v391 base.V128
	_ = v391
	var v392 base.V128
	_ = v392
	var v394 base.V128
	_ = v394
	var v395 base.V128
	_ = v395
	var v397 base.V128
	_ = v397
	var v398 base.V128
	_ = v398
	var v400 base.V128
	_ = v400
	var v402 base.V128
	_ = v402
	var v404 base.V128
	_ = v404
	var v410 base.V128
	_ = v410
	var v411 base.V128
	_ = v411
	var v417 base.V128
	_ = v417
	var v418 base.V128
	_ = v418
	var v419 base.V128
	_ = v419
	var v420 base.V128
	_ = v420
	var v423 base.V128
	_ = v423
	var v424 base.V128
	_ = v424
	var v427 base.V128
	_ = v427
	var v428 base.V128
	_ = v428
	var v430 base.V128
	_ = v430
	var v434 base.V128
	_ = v434
	var v437 int32
	_ = v437
	var v439 base.V128
	_ = v439
	var v440 int32
	_ = v440
	var v442 base.V128
	_ = v442
	var v444 base.V128
	_ = v444
	var v445 base.V128
	_ = v445
	var v447 base.V128
	_ = v447
	var v448 base.V128
	_ = v448
	var v450 base.V128
	_ = v450
	var v451 base.V128
	_ = v451
	var v453 base.V128
	_ = v453
	var v456 base.V128
	_ = v456
	var v462 base.V128
	_ = v462
	var v463 base.V128
	_ = v463
	var v469 base.V128
	_ = v469
	var v470 base.V128
	_ = v470
	var v472 base.V128
	_ = v472
	var v476 base.V128
	_ = v476
	var v479 base.V128
	_ = v479
	var v480 base.V128
	_ = v480
	var v482 base.V128
	_ = v482
	var v486 base.V128
	_ = v486
	var v502 int32
	_ = v502
	var v503 int32
	_ = v503
	var v511 base.V128
	_ = v511
	var v513 int32
	_ = v513
	var v514 base.V128
	_ = v514
	var v515 base.V128
	_ = v515
	var v516 base.V128
	_ = v516
	var v518 base.V128
	_ = v518
	var v522 base.V128
	_ = v522
	var v523 base.V128
	_ = v523
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
	var v539 int32
	_ = v539
	var v543 base.V128
	_ = v543
	var v545 base.V128
	_ = v545
	var v546 base.V128
	_ = v546
	var v547 base.V128
	_ = v547
	var v549 base.V128
	_ = v549
	var v553 base.V128
	_ = v553
	var v564 base.V128
	_ = v564
	var v565 base.V128
	_ = v565
	var v576 base.V128
	_ = v576
	var v587 base.V128
	_ = v587
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v620 base.V128
	_ = v620
	var v622 int32
	_ = v622
	var v623 base.V128
	_ = v623
	var v624 base.V128
	_ = v624
	var v625 base.V128
	_ = v625
	var v627 base.V128
	_ = v627
	var v631 base.V128
	_ = v631
	var v632 base.V128
	_ = v632
	var v635 base.V128
	_ = v635
	var v637 base.V128
	_ = v637
	var v638 base.V128
	_ = v638
	var v639 base.V128
	_ = v639
	var v641 base.V128
	_ = v641
	var v648 int32
	_ = v648
	var v652 base.V128
	_ = v652
	var v654 base.V128
	_ = v654
	var v655 base.V128
	_ = v655
	var v656 base.V128
	_ = v656
	var v658 base.V128
	_ = v658
	var v662 base.V128
	_ = v662
	var v673 base.V128
	_ = v673
	var v674 base.V128
	_ = v674
	var v685 base.V128
	_ = v685
	var v696 base.V128
	_ = v696
	var v705 int32
	_ = v705
	var v710 int32
	_ = v710
	var v728 int32
	_ = v728
	var v729 int32
	_ = v729
	var v752 int32
	_ = v752
	var v755 int32
	_ = v755
	var v761 int32
	_ = v761
	var v762 int32
	_ = v762
	var v764 int32
	_ = v764
	var v766 int32
	_ = v766
	var v767 int32
	_ = v767
	var v771 int32
	_ = v771
	var v772 int32
	_ = v772
	var v773 int32
	_ = v773
	var v775 int32
	_ = v775
	var v786 int32
	_ = v786
	var v811 int32
	_ = v811
	var v812 int32
	_ = v812
	var v816 int32
	_ = v816
	var v820 int32
	_ = v820
	var v821 int32
	_ = v821
	var v856 int32
	_ = v856
	var v857 int32
	_ = v857
	var v862 int64
	_ = v862
	var v865 int32
	_ = v865
	var v866 int32
	_ = v866
	var v878 int32
	_ = v878
	var v896 int32
	_ = v896
	var v897 int32
	_ = v897
	var v898 int32
	_ = v898
	var v909 int32
	_ = v909
	var v934 int32
	_ = v934
	var v935 int32
	_ = v935
	var v939 int32
	_ = v939
	var v943 int32
	_ = v943
	var v944 int32
	_ = v944
	var v979 int32
	_ = v979
	var v980 int32
	_ = v980
	var v985 int64
	_ = v985
	var v988 int32
	_ = v988
	var v989 int32
	_ = v989
	var v1001 int32
	_ = v1001
	var v1019 int32
	_ = v1019
	var v1020 base.V128
	_ = v1020
	var v1021 int32
	_ = v1021
	var v1022 base.V128
	_ = v1022
	var v1024 base.V128
	_ = v1024
	var v1025 base.V128
	_ = v1025
	var v1026 int32
	_ = v1026
	var v1027 base.V128
	_ = v1027
	var v1028 base.V128
	_ = v1028
	var v1030 base.V128
	_ = v1030
	var v1031 base.V128
	_ = v1031
	var v1033 base.V128
	_ = v1033
	var v1035 base.V128
	_ = v1035
	var v1037 base.V128
	_ = v1037
	var v1043 base.V128
	_ = v1043
	var v1044 base.V128
	_ = v1044
	var v1050 base.V128
	_ = v1050
	var v1051 base.V128
	_ = v1051
	var v1052 base.V128
	_ = v1052
	var v1053 base.V128
	_ = v1053
	var v1056 base.V128
	_ = v1056
	var v1057 base.V128
	_ = v1057
	var v1060 base.V128
	_ = v1060
	var v1061 base.V128
	_ = v1061
	var v1063 base.V128
	_ = v1063
	var v1067 base.V128
	_ = v1067
	var v1073 int32
	_ = v1073
	var v1075 int32
	_ = v1075
	var v1076 int32
	_ = v1076
	var v1087 int32
	_ = v1087
	var v1112 int32
	_ = v1112
	var v1113 int32
	_ = v1113
	var v1117 int32
	_ = v1117
	var v1121 int32
	_ = v1121
	var v1122 int32
	_ = v1122
	var v1157 int32
	_ = v1157
	var v1158 int32
	_ = v1158
	var v1163 int64
	_ = v1163
	var v1166 int32
	_ = v1166
	var v1167 int32
	_ = v1167
	var v1179 int32
	_ = v1179
	var v1197 int32
	_ = v1197
	var v1208 int32
	_ = v1208
	var v1233 int32
	_ = v1233
	var v1234 int32
	_ = v1234
	var v1238 int32
	_ = v1238
	var v1242 int32
	_ = v1242
	var v1243 int32
	_ = v1243
	var v1278 int32
	_ = v1278
	var v1279 int32
	_ = v1279
	var v1284 int64
	_ = v1284
	var v1287 int32
	_ = v1287
	var v1288 int32
	_ = v1288
	var v1300 int32
	_ = v1300
	var v1319 base.V128
	_ = v1319
	var v1320 int32
	_ = v1320
	var v1321 base.V128
	_ = v1321
	var v1323 base.V128
	_ = v1323
	var v1324 base.V128
	_ = v1324
	var v1326 base.V128
	_ = v1326
	var v1327 base.V128
	_ = v1327
	var v1329 base.V128
	_ = v1329
	var v1330 base.V128
	_ = v1330
	var v1332 base.V128
	_ = v1332
	var v1335 base.V128
	_ = v1335
	var v1341 base.V128
	_ = v1341
	var v1342 base.V128
	_ = v1342
	var v1348 base.V128
	_ = v1348
	var v1349 base.V128
	_ = v1349
	var v1350 base.V128
	_ = v1350
	var v1351 base.V128
	_ = v1351
	var v1354 base.V128
	_ = v1354
	var v1355 base.V128
	_ = v1355
	var v1358 base.V128
	_ = v1358
	var v1359 base.V128
	_ = v1359
	var v1361 base.V128
	_ = v1361
	var v1365 base.V128
	_ = v1365
	var v1371 int32
	_ = v1371
	var v1372 int32
	_ = v1372
	var v1374 int32
	_ = v1374
	var v1387 int32
	_ = v1387
	var v1388 int32
	_ = v1388
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
	var v1403 base.V128
	_ = v1403
	var v1407 base.V128
	_ = v1407
	var v1408 base.V128
	_ = v1408
	var v1411 base.V128
	_ = v1411
	var v1413 base.V128
	_ = v1413
	var v1414 base.V128
	_ = v1414
	var v1415 base.V128
	_ = v1415
	var v1417 base.V128
	_ = v1417
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
	var v1438 base.V128
	_ = v1438
	var v1449 base.V128
	_ = v1449
	var v1450 base.V128
	_ = v1450
	var v1461 base.V128
	_ = v1461
	var v1472 base.V128
	_ = v1472
	var v1481 int32
	_ = v1481
	var v1486 int32
	_ = v1486
	var v1490 int32
	_ = v1490
	var v1503 int32
	_ = v1503
	var v1504 int32
	_ = v1504
	var v1512 base.V128
	_ = v1512
	var v1514 int32
	_ = v1514
	var v1515 base.V128
	_ = v1515
	var v1516 base.V128
	_ = v1516
	var v1517 base.V128
	_ = v1517
	var v1519 base.V128
	_ = v1519
	var v1523 base.V128
	_ = v1523
	var v1524 base.V128
	_ = v1524
	var v1527 base.V128
	_ = v1527
	var v1529 base.V128
	_ = v1529
	var v1530 base.V128
	_ = v1530
	var v1531 base.V128
	_ = v1531
	var v1533 base.V128
	_ = v1533
	var v1540 int32
	_ = v1540
	var v1544 base.V128
	_ = v1544
	var v1546 base.V128
	_ = v1546
	var v1547 base.V128
	_ = v1547
	var v1548 base.V128
	_ = v1548
	var v1550 base.V128
	_ = v1550
	var v1554 base.V128
	_ = v1554
	var v1565 base.V128
	_ = v1565
	var v1566 base.V128
	_ = v1566
	var v1577 base.V128
	_ = v1577
	var v1588 base.V128
	_ = v1588
	var v1602 int32
	_ = v1602
	var v1615 int32
	_ = v1615
	var v1616 int32
	_ = v1616
	var v1624 base.V128
	_ = v1624
	var v1626 int32
	_ = v1626
	var v1627 base.V128
	_ = v1627
	var v1628 base.V128
	_ = v1628
	var v1629 base.V128
	_ = v1629
	var v1631 base.V128
	_ = v1631
	var v1635 base.V128
	_ = v1635
	var v1636 base.V128
	_ = v1636
	var v1639 base.V128
	_ = v1639
	var v1641 base.V128
	_ = v1641
	var v1642 base.V128
	_ = v1642
	var v1643 base.V128
	_ = v1643
	var v1645 base.V128
	_ = v1645
	var v1652 int32
	_ = v1652
	var v1656 base.V128
	_ = v1656
	var v1658 base.V128
	_ = v1658
	var v1659 base.V128
	_ = v1659
	var v1660 base.V128
	_ = v1660
	var v1662 base.V128
	_ = v1662
	var v1666 base.V128
	_ = v1666
	var v1677 base.V128
	_ = v1677
	var v1678 base.V128
	_ = v1678
	var v1689 base.V128
	_ = v1689
	var v1700 base.V128
	_ = v1700
	var v1709 int32
	_ = v1709
	var v1710 int32
	_ = v1710
	var v1713 int32
	_ = v1713
	var v1714 int32
	_ = v1714
	var v1716 int32
	_ = v1716
	v10 = int32(0)
	v33 = m.G0
	v35 = v33 - int32(528)
	m.G0 = v35
	base.MemoryFill(m, v35+int32(64), v10, int32(463))
	v162 = int32(1)
	v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v169 = int32(base.Ui32(v163+v164)>>(uint(v162)%32)) + v162
	v172 = int32(base.Ui32(v169+v164) >> (uint(v162) % 32))
	v175 = int32(8)
	v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v181 = int32(base.Ui32(v177*int32(_a_F_UpsampleRgba4444LinePair_SSE2_0)) >> (uint(v175) % 32))
	v182 = int32(base.Ui32(v172*int32(_a_F_UpsampleRgba4444LinePair_SSE2_1))>>(uint(v175)%32)) + v181
	v184 = v182 + int32(-17685)
	if base.Ui32(v182) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_2)) {
		v191 = int32(0)
	} else {
		v191 = int32(240)
	}
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_3)) {
		v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
	} else {
		v194 = v191
	}
	v196 = v194 | int32(15)
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v196)
	v198 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v199 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v201 = int32(1)
	v204 = int32(base.Ui32(v198+v199)>>(uint(v201)%32)) + v201
	v207 = int32(base.Ui32(v204+v199) >> (uint(v201) % 32))
	v212 = int32(base.Ui32(v207*int32(_a_F_UpsampleRgba4444LinePair_SSE2_4))>>(uint(int32(8))%32)) + v181
	v214 = v212 + int32(-14234)
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_5)) {
		v221 = int32(0)
	} else {
		v221 = int32(240)
	}
	if base.Ui32(v214) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_3)) {
		v224 = int32(base.Ui32(v214) >> (uint(int32(6)) % 32))
	} else {
		v224 = v221
	}
	v229 = int32(8)
	v236 = v181 - (int32(base.Ui32(v172*int32(_a_F_UpsampleRgba4444LinePair_SSE2_6))>>(uint(v229)%32)) + int32(base.Ui32(v207*int32(_a_F_UpsampleRgba4444LinePair_SSE2_7))>>(uint(v229)%32)))
	v238 = v236 + int32(_a_F_UpsampleRgba4444LinePair_SSE2_8)
	if v236 < int32(-8708) {
		v245 = int32(0)
	} else {
		v245 = int32(15)
	}
	if base.Ui32(v238) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_3)) {
		v248 = int32(base.Ui32(v238) >> (uint(int32(10)) % 32))
	} else {
		v248 = v245
	}
	v249 = v224&int32(240) | v248
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v249)
	if l1 == int32(0) {
	} else {
		v253 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v256 = int32(8)
		v257 = int32(base.Ui32(v253*int32(_a_F_UpsampleRgba4444LinePair_SSE2_0)) >> (uint(v256) % 32))
		v260 = int32(base.Ui32(v169+v163) >> (uint(int32(1)) % 32))
		v265 = v257 + int32(base.Ui32(v260*int32(_a_F_UpsampleRgba4444LinePair_SSE2_1))>>(uint(v256)%32))
		v267 = v265 + int32(-17685)
		if base.Ui32(v265) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_2)) {
			v274 = int32(0)
		} else {
			v274 = int32(240)
		}
		if base.Ui32(v267) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_3)) {
			v277 = int32(base.Ui32(v267) >> (uint(int32(6)) % 32))
		} else {
			v277 = v274
		}
		v279 = v277 | int32(15)
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v279)
		v283 = int32(base.Ui32(v204+v198) >> (uint(int32(1)) % 32))
		v288 = v257 + int32(base.Ui32(v283*int32(_a_F_UpsampleRgba4444LinePair_SSE2_4))>>(uint(int32(8))%32))
		v290 = v288 + int32(-14234)
		if base.Ui32(v288) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_5)) {
			v297 = int32(0)
		} else {
			v297 = int32(240)
		}
		if base.Ui32(v290) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_3)) {
			v300 = int32(base.Ui32(v290) >> (uint(int32(6)) % 32))
		} else {
			v300 = v297
		}
		v305 = int32(8)
		v312 = v257 - (int32(base.Ui32(v283*int32(_a_F_UpsampleRgba4444LinePair_SSE2_7))>>(uint(v305)%32)) + int32(base.Ui32(v260*int32(_a_F_UpsampleRgba4444LinePair_SSE2_6))>>(uint(v305)%32)))
		v314 = v312 + int32(_a_F_UpsampleRgba4444LinePair_SSE2_8)
		if v312 < int32(-8708) {
			v321 = int32(0)
		} else {
			v321 = int32(15)
		}
		if base.Ui32(v314) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_SSE2_3)) {
			v324 = int32(base.Ui32(v314) >> (uint(int32(10)) % 32))
		} else {
			v324 = v321
		}
		v325 = v300&int32(240) | v324
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v325)
	}
	v333 = v35 + int32(96)
	v335 = v35 + int32(64)
	if l8 < int32(34) {
		v728 = v10
		v729 = v162
	} else {
		v338 = int32(2)
		v342 = int32(1)
		v350 = int32(0)
		v362 = v350
		v363 = l7 + v338
		v367 = v350
		v368 = l6 + v338
		for {
			v384 = l4 + v362
			v385 = int32(0)
			v386 = base.Simd_g_v128_load_rng(m, v384, v385, int32(0), int32(17))
			v387 = l2 + v362
			v389 = base.Simd_g_v128_load_rng(m, v387, v385, int32(0), int32(17))
			v390 = int32(1)
			v391 = base.Simd_g_v128_load_nc(m, v384, v390)
			v392 = base.Simd_g_i8x16_avgr_u(v389, v391)
			v394 = base.Simd_g_v128_load_nc(m, v387, v390)
			v395 = base.Simd_g_i8x16_avgr_u(v386, v394)
			v397 = base.Simd_g_v128_xor(v391, v389)
			v398 = base.Simd_g_v128_xor(v386, v394)
			v400 = base.Simd_g_v128_xor(v392, v395)
			v402 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k0)
			v404 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v392, v395), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v397, v398), v400), v402))
			v410 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v404, v392), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v404, v392), base.Simd_g_v128_and(v400, v397)), v402))
			v411 = base.Simd_g_i8x16_avgr_u(v386, v410)
			v417 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v404, v395), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v404, v395), base.Simd_g_v128_and(v400, v398)), v402))
			v418 = base.Simd_g_i8x16_avgr_u(v391, v417)
			v419 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k1)
			v420 = base.Simd_g_i8x16_shuffle2(v411, v418, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v335, int32(80), v420)
			v423 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
			v424 = base.Simd_g_i8x16_shuffle2(v411, v418, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v335, int32(64), v424)
			v427 = base.Simd_g_i8x16_avgr_u(v389, v417)
			v428 = base.Simd_g_i8x16_avgr_u(v394, v410)
			v430 = base.Simd_g_i8x16_shuffle2(v427, v428, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v335, int32(16), v430)
			v434 = base.Simd_g_i8x16_shuffle2(v427, v428, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v335, v385, v434)
			v437 = l5 + v362
			v439 = base.Simd_g_v128_load_rng(m, v437, v385, int32(0), int32(17))
			v440 = l3 + v362
			v442 = base.Simd_g_v128_load_rng(m, v440, v385, int32(0), int32(17))
			v444 = base.Simd_g_v128_load_nc(m, v437, v390)
			v445 = base.Simd_g_i8x16_avgr_u(v442, v444)
			v447 = base.Simd_g_v128_load_nc(m, v440, v390)
			v448 = base.Simd_g_i8x16_avgr_u(v439, v447)
			v450 = base.Simd_g_v128_xor(v444, v442)
			v451 = base.Simd_g_v128_xor(v439, v447)
			v453 = base.Simd_g_v128_xor(v445, v448)
			v456 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v445, v448), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v450, v451), v453), v402))
			v462 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v456, v445), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v456, v445), base.Simd_g_v128_and(v453, v450)), v402))
			v463 = base.Simd_g_i8x16_avgr_u(v439, v462)
			v469 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v456, v448), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v456, v448), base.Simd_g_v128_and(v453, v451)), v402))
			v470 = base.Simd_g_i8x16_avgr_u(v444, v469)
			v472 = base.Simd_g_i8x16_shuffle2(v463, v470, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v335, int32(112), v472)
			v476 = base.Simd_g_i8x16_shuffle2(v463, v470, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v335, int32(96), v476)
			v479 = base.Simd_g_i8x16_avgr_u(v442, v469)
			v480 = base.Simd_g_i8x16_avgr_u(v447, v462)
			v482 = base.Simd_g_i8x16_shuffle2(v479, v480, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v335, int32(48), v482)
			v486 = base.Simd_g_i8x16_shuffle2(v479, v480, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v335, int32(32), v486)
			v502 = v368
			v503 = v385
			for {
				v511 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k7)
				v513 = int32(0)
				v514 = base.Simd_g_v128_load64_zero(m, l0+v342+v367+v503, v513)
				v515 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
				v516 = base.Simd_g_i8x16_shuffle2(v511, v514, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v518 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k8)
				v522 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k9)
				v523 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v516), v518), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v516), v518), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))
				v526 = base.Simd_g_v128_load64_zero(m, v333+v503, v513)
				v528 = base.Simd_g_i8x16_shuffle2(v511, v526, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v529 = base.Simd_g_i32x4_extend_low_i16x8_u(v528)
				v530 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k12)
				v532 = base.Simd_g_i32x4_extend_high_i16x8_u(v528)
				v539 = int32(6)
				v543 = base.Simd_g_v128_load64_zero(m, v335+v503, v513)
				v545 = base.Simd_g_i8x16_shuffle2(v511, v543, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v546 = base.Simd_g_i32x4_extend_low_i16x8_u(v545)
				v547 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k13)
				v549 = base.Simd_g_i32x4_extend_high_i16x8_u(v545)
				v553 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k14)
				v564 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v523, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v529, v530), base.Simd_g_i32x4_mul(v532, v530), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k15)), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v523, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v546, v547), base.Simd_g_i32x4_mul(v549, v547), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v529, v553), base.Simd_g_i32x4_mul(v532, v553), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k16)), v539))
				v565 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k17)
				v576 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v546, v565), base.Simd_g_i32x4_mul(v549, v565), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), v523), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k18)), v539), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k19))
				v587 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v564, v576, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3)), int32(4)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k20)), base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v564, v576, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k21)))
				base.Simd_g_v128_store(m, v502, v513, v587)
				if base.Ui32(v503) < base.Ui32(int32(24)) {
					v502 = v502 + int32(16)
					v503 = v503 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			if l1 == int32(0) {
			} else {
				v611 = v363
				v612 = int32(0)
				for {
					v620 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k7)
					v622 = int32(0)
					v623 = base.Simd_g_v128_load64_zero(m, l1+v342+v367+v612, v622)
					v624 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
					v625 = base.Simd_g_i8x16_shuffle2(v620, v623, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
					v627 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k8)
					v631 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k9)
					v632 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v625), v627), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v625), v627), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))
					v635 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v612, v622)
					v637 = base.Simd_g_i8x16_shuffle2(v620, v635, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
					v638 = base.Simd_g_i32x4_extend_low_i16x8_u(v637)
					v639 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k12)
					v641 = base.Simd_g_i32x4_extend_high_i16x8_u(v637)
					v648 = int32(6)
					v652 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v612, v622)
					v654 = base.Simd_g_i8x16_shuffle2(v620, v652, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
					v655 = base.Simd_g_i32x4_extend_low_i16x8_u(v654)
					v656 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k13)
					v658 = base.Simd_g_i32x4_extend_high_i16x8_u(v654)
					v662 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k14)
					v673 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v632, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v638, v639), base.Simd_g_i32x4_mul(v641, v639), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k15)), v648), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v632, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v655, v656), base.Simd_g_i32x4_mul(v658, v656), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v638, v662), base.Simd_g_i32x4_mul(v641, v662), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k16)), v648))
					v674 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k17)
					v685 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v655, v674), base.Simd_g_i32x4_mul(v658, v674), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), v632), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k18)), v648), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k19))
					v696 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v673, v685, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3)), int32(4)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k20)), base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v673, v685, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k21)))
					base.Simd_g_v128_store(m, v611, v622, v696)
					if base.Ui32(v612) < base.Ui32(int32(24)) {
						v611 = v611 + int32(16)
						v612 = v612 + int32(8)
						continue
					} else {
						break
					}
					break
				}
			}
			v705 = int32(64)
			v710 = v362 + int32(16)
			if v367+int32(66) <= l8 {
				v362 = v710
				v363 = v363 + v705
				v367 = v367 + int32(32)
				v368 = v368 + v705
				continue
			} else {
				break
			}
			break
		}
		v728 = v710
		v729 = v367 + int32(33)
	}
	if l8 < int32(2) {
	} else {
		v752 = int32(32)
		v755 = int32(1)
		v761 = int32(base.Ui32(l8+v755)>>(uint(v755)%32)) - int32(base.Ui32(v729)>>(uint(v755)%32))
		v762 = F_memcpy(m, v35+v752, l2+v728, v761)
		mBase = m.M
		v764 = F_memcpy(m, v35, l4+v728, v761)
		mBase = m.M
		v766 = v764 + v752
		v767 = v766 + v761
		v771 = v761 + int32(-1)
		v772 = v766 + v771
		v773 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v772))))
		v775 = int32(17) - v761
		if base.Ui32(v775) < base.Ui32(int32(33)) {
			if v775 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v767))) = uint8(v773)
				v786 = v767 + v775
				*(*uint8)(unsafe.Add(mBase, uint32(v786+int32(-1)))) = uint8(v773)
				if base.Ui32(v775) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v767)+2)) = uint8(v773)
					*(*uint8)(unsafe.Add(mBase, uint32(v767)+1)) = uint8(v773)
					*(*uint8)(unsafe.Add(mBase, uint32(v786+int32(-3)))) = uint8(v773)
					*(*uint8)(unsafe.Add(mBase, uint32(v786+int32(-2)))) = uint8(v773)
					if base.Ui32(v775) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v767)+3)) = uint8(v773)
						*(*uint8)(unsafe.Add(mBase, uint32(v786+int32(-4)))) = uint8(v773)
						if base.Ui32(v775) < base.Ui32(int32(9)) {
						} else {
							v811 = (int32(0) - v767) & int32(3)
							v812 = v767 + v811
							v816 = v773 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v812))) = v816
							v820 = (v775 - v811) & int32(60)
							v821 = v812 + v820
							*(*int32)(unsafe.Add(mBase, uint32(v821+int32(-4)))) = v816
							if base.Ui32(v820) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v812)+8)) = v816
								*(*int32)(unsafe.Add(mBase, uint32(v812)+4)) = v816
								*(*int32)(unsafe.Add(mBase, uint32(v821+int32(-8)))) = v816
								*(*int32)(unsafe.Add(mBase, uint32(v821+int32(-12)))) = v816
								if base.Ui32(v820) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v812)+24)) = v816
									*(*int32)(unsafe.Add(mBase, uint32(v812)+20)) = v816
									*(*int32)(unsafe.Add(mBase, uint32(v812)+16)) = v816
									*(*int32)(unsafe.Add(mBase, uint32(v812)+12)) = v816
									*(*int32)(unsafe.Add(mBase, uint32(v821+int32(-16)))) = v816
									*(*int32)(unsafe.Add(mBase, uint32(v821+int32(-20)))) = v816
									*(*int32)(unsafe.Add(mBase, uint32(v821+int32(-24)))) = v816
									*(*int32)(unsafe.Add(mBase, uint32(v821+int32(-28)))) = v816
									v856 = v812&int32(4) | int32(24)
									v857 = v820 - v856
									if base.Ui32(v857) < base.Ui32(int32(32)) {
									} else {
										v862 = base.I64_extend_i32_u(v816) * int64(4294967297)
										v865 = v857
										v866 = v812 + v856
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v866)+24)) = v862
											*(*int64)(unsafe.Add(mBase, uint32(v866)+16)) = v862
											*(*int64)(unsafe.Add(mBase, uint32(v866)+8)) = v862
											*(*int64)(unsafe.Add(mBase, uint32(v866))) = v862
											v878 = v865 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v878) {
												v865 = v878
												v866 = v866 + int32(32)
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
			base.MemoryFill(m, v767, v773, v775)
		}
		v896 = v764 + v761
		v897 = v764 + v771
		v898 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v897))))
		if base.Ui32(v775) < base.Ui32(int32(33)) {
			if v775 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v896))) = uint8(v898)
				v909 = v896 + v775
				*(*uint8)(unsafe.Add(mBase, uint32(v909+int32(-1)))) = uint8(v898)
				if base.Ui32(v775) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v896)+2)) = uint8(v898)
					*(*uint8)(unsafe.Add(mBase, uint32(v896)+1)) = uint8(v898)
					*(*uint8)(unsafe.Add(mBase, uint32(v909+int32(-3)))) = uint8(v898)
					*(*uint8)(unsafe.Add(mBase, uint32(v909+int32(-2)))) = uint8(v898)
					if base.Ui32(v775) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v896)+3)) = uint8(v898)
						*(*uint8)(unsafe.Add(mBase, uint32(v909+int32(-4)))) = uint8(v898)
						if base.Ui32(v775) < base.Ui32(int32(9)) {
						} else {
							v934 = (int32(0) - v896) & int32(3)
							v935 = v896 + v934
							v939 = v898 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v935))) = v939
							v943 = (v775 - v934) & int32(60)
							v944 = v935 + v943
							*(*int32)(unsafe.Add(mBase, uint32(v944+int32(-4)))) = v939
							if base.Ui32(v943) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v935)+8)) = v939
								*(*int32)(unsafe.Add(mBase, uint32(v935)+4)) = v939
								*(*int32)(unsafe.Add(mBase, uint32(v944+int32(-8)))) = v939
								*(*int32)(unsafe.Add(mBase, uint32(v944+int32(-12)))) = v939
								if base.Ui32(v943) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v935)+24)) = v939
									*(*int32)(unsafe.Add(mBase, uint32(v935)+20)) = v939
									*(*int32)(unsafe.Add(mBase, uint32(v935)+16)) = v939
									*(*int32)(unsafe.Add(mBase, uint32(v935)+12)) = v939
									*(*int32)(unsafe.Add(mBase, uint32(v944+int32(-16)))) = v939
									*(*int32)(unsafe.Add(mBase, uint32(v944+int32(-20)))) = v939
									*(*int32)(unsafe.Add(mBase, uint32(v944+int32(-24)))) = v939
									*(*int32)(unsafe.Add(mBase, uint32(v944+int32(-28)))) = v939
									v979 = v935&int32(4) | int32(24)
									v980 = v943 - v979
									if base.Ui32(v980) < base.Ui32(int32(32)) {
									} else {
										v985 = base.I64_extend_i32_u(v939) * int64(4294967297)
										v988 = v980
										v989 = v935 + v979
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v989)+24)) = v985
											*(*int64)(unsafe.Add(mBase, uint32(v989)+16)) = v985
											*(*int64)(unsafe.Add(mBase, uint32(v989)+8)) = v985
											*(*int64)(unsafe.Add(mBase, uint32(v989))) = v985
											v1001 = v988 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1001) {
												v988 = v1001
												v989 = v989 + int32(32)
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
			base.MemoryFill(m, v896, v898, v775)
		}
		v1019 = int32(0)
		v1020 = base.Simd_g_v128_load_rng(m, v764, v1019, int32(0), int32(49))
		v1021 = int32(32)
		v1022 = base.Simd_g_v128_load_nc(m, v764, v1021)
		v1024 = base.Simd_g_v128_load_nc(m, v764, int32(1))
		v1025 = base.Simd_g_i8x16_avgr_u(v1022, v1024)
		v1026 = int32(33)
		v1027 = base.Simd_g_v128_load_nc(m, v764, v1026)
		v1028 = base.Simd_g_i8x16_avgr_u(v1020, v1027)
		v1030 = base.Simd_g_v128_xor(v1024, v1022)
		v1031 = base.Simd_g_v128_xor(v1020, v1027)
		v1033 = base.Simd_g_v128_xor(v1025, v1028)
		v1035 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k0)
		v1037 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1025, v1028), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1030, v1031), v1033), v1035))
		v1043 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1037, v1025), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1037, v1025), base.Simd_g_v128_and(v1033, v1030)), v1035))
		v1044 = base.Simd_g_i8x16_avgr_u(v1020, v1043)
		v1050 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1037, v1028), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1037, v1028), base.Simd_g_v128_and(v1033, v1031)), v1035))
		v1051 = base.Simd_g_i8x16_avgr_u(v1024, v1050)
		v1052 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k1)
		v1053 = base.Simd_g_i8x16_shuffle2(v1044, v1051, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v335, int32(80), v1053)
		v1056 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
		v1057 = base.Simd_g_i8x16_shuffle2(v1044, v1051, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v335, int32(64), v1057)
		v1060 = base.Simd_g_i8x16_avgr_u(v1022, v1050)
		v1061 = base.Simd_g_i8x16_avgr_u(v1027, v1043)
		v1063 = base.Simd_g_i8x16_shuffle2(v1060, v1061, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v335, int32(16), v1063)
		v1067 = base.Simd_g_i8x16_shuffle2(v1060, v1061, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v335, v1019, v1067)
		v1073 = F_memcpy(m, v764+v1021, l3+v728, v761)
		mBase = m.M
		v1075 = F_memcpy(m, v764, l5+v728, v761)
		mBase = m.M
		v1076 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v772))))
		if base.Ui32(v775) < base.Ui32(v1026) {
			if v775 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v767))) = uint8(v1076)
				v1087 = v767 + v775
				*(*uint8)(unsafe.Add(mBase, uint32(v1087+int32(-1)))) = uint8(v1076)
				if base.Ui32(v775) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v767)+2)) = uint8(v1076)
					*(*uint8)(unsafe.Add(mBase, uint32(v767)+1)) = uint8(v1076)
					*(*uint8)(unsafe.Add(mBase, uint32(v1087+int32(-3)))) = uint8(v1076)
					*(*uint8)(unsafe.Add(mBase, uint32(v1087+int32(-2)))) = uint8(v1076)
					if base.Ui32(v775) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v767)+3)) = uint8(v1076)
						*(*uint8)(unsafe.Add(mBase, uint32(v1087+int32(-4)))) = uint8(v1076)
						if base.Ui32(v775) < base.Ui32(int32(9)) {
						} else {
							v1112 = (int32(0) - v767) & int32(3)
							v1113 = v767 + v1112
							v1117 = v1076 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1113))) = v1117
							v1121 = (v775 - v1112) & int32(60)
							v1122 = v1113 + v1121
							*(*int32)(unsafe.Add(mBase, uint32(v1122+int32(-4)))) = v1117
							if base.Ui32(v1121) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1113)+8)) = v1117
								*(*int32)(unsafe.Add(mBase, uint32(v1113)+4)) = v1117
								*(*int32)(unsafe.Add(mBase, uint32(v1122+int32(-8)))) = v1117
								*(*int32)(unsafe.Add(mBase, uint32(v1122+int32(-12)))) = v1117
								if base.Ui32(v1121) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1113)+24)) = v1117
									*(*int32)(unsafe.Add(mBase, uint32(v1113)+20)) = v1117
									*(*int32)(unsafe.Add(mBase, uint32(v1113)+16)) = v1117
									*(*int32)(unsafe.Add(mBase, uint32(v1113)+12)) = v1117
									*(*int32)(unsafe.Add(mBase, uint32(v1122+int32(-16)))) = v1117
									*(*int32)(unsafe.Add(mBase, uint32(v1122+int32(-20)))) = v1117
									*(*int32)(unsafe.Add(mBase, uint32(v1122+int32(-24)))) = v1117
									*(*int32)(unsafe.Add(mBase, uint32(v1122+int32(-28)))) = v1117
									v1157 = v1113&int32(4) | int32(24)
									v1158 = v1121 - v1157
									if base.Ui32(v1158) < base.Ui32(int32(32)) {
									} else {
										v1163 = base.I64_extend_i32_u(v1117) * int64(4294967297)
										v1166 = v1158
										v1167 = v1113 + v1157
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1167)+24)) = v1163
											*(*int64)(unsafe.Add(mBase, uint32(v1167)+16)) = v1163
											*(*int64)(unsafe.Add(mBase, uint32(v1167)+8)) = v1163
											*(*int64)(unsafe.Add(mBase, uint32(v1167))) = v1163
											v1179 = v1166 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1179) {
												v1166 = v1179
												v1167 = v1167 + int32(32)
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
			base.MemoryFill(m, v767, v1076, v775)
		}
		v1197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v897))))
		if base.Ui32(v775) < base.Ui32(int32(33)) {
			if v775 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v896))) = uint8(v1197)
				v1208 = v896 + v775
				*(*uint8)(unsafe.Add(mBase, uint32(v1208+int32(-1)))) = uint8(v1197)
				if base.Ui32(v775) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v896)+2)) = uint8(v1197)
					*(*uint8)(unsafe.Add(mBase, uint32(v896)+1)) = uint8(v1197)
					*(*uint8)(unsafe.Add(mBase, uint32(v1208+int32(-3)))) = uint8(v1197)
					*(*uint8)(unsafe.Add(mBase, uint32(v1208+int32(-2)))) = uint8(v1197)
					if base.Ui32(v775) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v896)+3)) = uint8(v1197)
						*(*uint8)(unsafe.Add(mBase, uint32(v1208+int32(-4)))) = uint8(v1197)
						if base.Ui32(v775) < base.Ui32(int32(9)) {
						} else {
							v1233 = (int32(0) - v896) & int32(3)
							v1234 = v896 + v1233
							v1238 = v1197 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1234))) = v1238
							v1242 = (v775 - v1233) & int32(60)
							v1243 = v1234 + v1242
							*(*int32)(unsafe.Add(mBase, uint32(v1243+int32(-4)))) = v1238
							if base.Ui32(v1242) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1234)+8)) = v1238
								*(*int32)(unsafe.Add(mBase, uint32(v1234)+4)) = v1238
								*(*int32)(unsafe.Add(mBase, uint32(v1243+int32(-8)))) = v1238
								*(*int32)(unsafe.Add(mBase, uint32(v1243+int32(-12)))) = v1238
								if base.Ui32(v1242) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1234)+24)) = v1238
									*(*int32)(unsafe.Add(mBase, uint32(v1234)+20)) = v1238
									*(*int32)(unsafe.Add(mBase, uint32(v1234)+16)) = v1238
									*(*int32)(unsafe.Add(mBase, uint32(v1234)+12)) = v1238
									*(*int32)(unsafe.Add(mBase, uint32(v1243+int32(-16)))) = v1238
									*(*int32)(unsafe.Add(mBase, uint32(v1243+int32(-20)))) = v1238
									*(*int32)(unsafe.Add(mBase, uint32(v1243+int32(-24)))) = v1238
									*(*int32)(unsafe.Add(mBase, uint32(v1243+int32(-28)))) = v1238
									v1278 = v1234&int32(4) | int32(24)
									v1279 = v1242 - v1278
									if base.Ui32(v1279) < base.Ui32(int32(32)) {
									} else {
										v1284 = base.I64_extend_i32_u(v1238) * int64(4294967297)
										v1287 = v1279
										v1288 = v1234 + v1278
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1288)+24)) = v1284
											*(*int64)(unsafe.Add(mBase, uint32(v1288)+16)) = v1284
											*(*int64)(unsafe.Add(mBase, uint32(v1288)+8)) = v1284
											*(*int64)(unsafe.Add(mBase, uint32(v1288))) = v1284
											v1300 = v1287 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1300) {
												v1287 = v1300
												v1288 = v1288 + int32(32)
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
			base.MemoryFill(m, v896, v1197, v775)
		}
		v1319 = base.Simd_g_v128_load_rng(m, v1075, int32(0), int32(0), int32(49))
		v1320 = int32(32)
		v1321 = base.Simd_g_v128_load_nc(m, v1075, v1320)
		v1323 = base.Simd_g_v128_load_nc(m, v1075, int32(1))
		v1324 = base.Simd_g_i8x16_avgr_u(v1321, v1323)
		v1326 = base.Simd_g_v128_load_nc(m, v1075, int32(33))
		v1327 = base.Simd_g_i8x16_avgr_u(v1319, v1326)
		v1329 = base.Simd_g_v128_xor(v1323, v1321)
		v1330 = base.Simd_g_v128_xor(v1319, v1326)
		v1332 = base.Simd_g_v128_xor(v1324, v1327)
		v1335 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1324, v1327), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1329, v1330), v1332), v1035))
		v1341 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1335, v1324), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1335, v1324), base.Simd_g_v128_and(v1332, v1329)), v1035))
		v1342 = base.Simd_g_i8x16_avgr_u(v1319, v1341)
		v1348 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1335, v1327), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1335, v1327), base.Simd_g_v128_and(v1332, v1330)), v1035))
		v1349 = base.Simd_g_i8x16_avgr_u(v1323, v1348)
		v1350 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k1)
		v1351 = base.Simd_g_i8x16_shuffle2(v1342, v1349, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v335, int32(112), v1351)
		v1354 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
		v1355 = base.Simd_g_i8x16_shuffle2(v1342, v1349, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v335, int32(96), v1355)
		v1358 = base.Simd_g_i8x16_avgr_u(v1321, v1348)
		v1359 = base.Simd_g_i8x16_avgr_u(v1326, v1341)
		v1361 = base.Simd_g_i8x16_shuffle2(v1358, v1359, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v335, int32(48), v1361)
		v1365 = base.Simd_g_i8x16_shuffle2(v1358, v1359, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v335, v1320, v1365)
		v1371 = l8 - v729
		v1372 = F_memcpy(m, v35+int32(448), l0+v729, v1371)
		mBase = m.M
		v1374 = v35 + int32(192)
		if l1 != 0 {
			v1490 = F_memcpy(m, v35+int32(480), l1+v729, v1371)
			mBase = m.M
			v1503 = v1374
			v1504 = int32(0)
			for {
				v1512 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k7)
				v1514 = int32(0)
				v1515 = base.Simd_g_v128_load64_zero(m, v1372+v1504, v1514)
				v1516 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
				v1517 = base.Simd_g_i8x16_shuffle2(v1512, v1515, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1519 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k8)
				v1523 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k9)
				v1524 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1517), v1519), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1517), v1519), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))
				v1527 = base.Simd_g_v128_load64_zero(m, v333+v1504, v1514)
				v1529 = base.Simd_g_i8x16_shuffle2(v1512, v1527, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1530 = base.Simd_g_i32x4_extend_low_i16x8_u(v1529)
				v1531 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k12)
				v1533 = base.Simd_g_i32x4_extend_high_i16x8_u(v1529)
				v1540 = int32(6)
				v1544 = base.Simd_g_v128_load64_zero(m, v335+v1504, v1514)
				v1546 = base.Simd_g_i8x16_shuffle2(v1512, v1544, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1547 = base.Simd_g_i32x4_extend_low_i16x8_u(v1546)
				v1548 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k13)
				v1550 = base.Simd_g_i32x4_extend_high_i16x8_u(v1546)
				v1554 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k14)
				v1565 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1524, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1530, v1531), base.Simd_g_i32x4_mul(v1533, v1531), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k15)), v1540), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1524, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1547, v1548), base.Simd_g_i32x4_mul(v1550, v1548), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1530, v1554), base.Simd_g_i32x4_mul(v1533, v1554), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k16)), v1540))
				v1566 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k17)
				v1577 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1547, v1566), base.Simd_g_i32x4_mul(v1550, v1566), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), v1524), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k18)), v1540), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k19))
				v1588 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v1565, v1577, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3)), int32(4)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k20)), base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v1565, v1577, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k21)))
				base.Simd_g_v128_store(m, v1503, v1514, v1588)
				if base.Ui32(v1504) < base.Ui32(int32(24)) {
					v1503 = v1503 + int32(16)
					v1504 = v1504 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1602 = v35 + int32(320)
			v1615 = v1602
			v1616 = int32(0)
			for {
				v1624 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k7)
				v1626 = int32(0)
				v1627 = base.Simd_g_v128_load64_zero(m, v1490+v1616, v1626)
				v1628 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
				v1629 = base.Simd_g_i8x16_shuffle2(v1624, v1627, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1631 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k8)
				v1635 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k9)
				v1636 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1629), v1631), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1629), v1631), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))
				v1639 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v1616, v1626)
				v1641 = base.Simd_g_i8x16_shuffle2(v1624, v1639, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1642 = base.Simd_g_i32x4_extend_low_i16x8_u(v1641)
				v1643 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k12)
				v1645 = base.Simd_g_i32x4_extend_high_i16x8_u(v1641)
				v1652 = int32(6)
				v1656 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v1616, v1626)
				v1658 = base.Simd_g_i8x16_shuffle2(v1624, v1656, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1659 = base.Simd_g_i32x4_extend_low_i16x8_u(v1658)
				v1660 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k13)
				v1662 = base.Simd_g_i32x4_extend_high_i16x8_u(v1658)
				v1666 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k14)
				v1677 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1636, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1642, v1643), base.Simd_g_i32x4_mul(v1645, v1643), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k15)), v1652), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1636, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1659, v1660), base.Simd_g_i32x4_mul(v1662, v1660), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1642, v1666), base.Simd_g_i32x4_mul(v1645, v1666), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k16)), v1652))
				v1678 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k17)
				v1689 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1659, v1678), base.Simd_g_i32x4_mul(v1662, v1678), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), v1636), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k18)), v1652), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k19))
				v1700 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v1677, v1689, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3)), int32(4)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k20)), base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v1677, v1689, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k21)))
				base.Simd_g_v128_store(m, v1615, v1626, v1700)
				if base.Ui32(v1616) < base.Ui32(int32(24)) {
					v1615 = v1615 + int32(16)
					v1616 = v1616 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1709 = int32(1)
			v1710 = v729 << (uint(v1709) % 32)
			v1713 = v1371 << (uint(v1709) % 32)
			v1714 = F_memcpy(m, l6+v1710, v1374, v1713)
			mBase = m.M
			v1716 = F_memcpy(m, l7+v1710, v1602, v1713)
			mBase = m.M
		} else {
			v1387 = v1374
			v1388 = int32(0)
			for {
				v1396 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k7)
				v1398 = int32(0)
				v1399 = base.Simd_g_v128_load64_zero(m, v1372+v1388, v1398)
				v1400 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k4)
				v1401 = base.Simd_g_i8x16_shuffle2(v1396, v1399, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1403 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k8)
				v1407 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k9)
				v1408 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1401), v1403), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1401), v1403), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))
				v1411 = base.Simd_g_v128_load64_zero(m, v333+v1388, v1398)
				v1413 = base.Simd_g_i8x16_shuffle2(v1396, v1411, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1414 = base.Simd_g_i32x4_extend_low_i16x8_u(v1413)
				v1415 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k12)
				v1417 = base.Simd_g_i32x4_extend_high_i16x8_u(v1413)
				v1424 = int32(6)
				v1428 = base.Simd_g_v128_load64_zero(m, v335+v1388, v1398)
				v1430 = base.Simd_g_i8x16_shuffle2(v1396, v1428, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6))
				v1431 = base.Simd_g_i32x4_extend_low_i16x8_u(v1430)
				v1432 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k13)
				v1434 = base.Simd_g_i32x4_extend_high_i16x8_u(v1430)
				v1438 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k14)
				v1449 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1408, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1414, v1415), base.Simd_g_i32x4_mul(v1417, v1415), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k15)), v1424), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1408, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1431, v1432), base.Simd_g_i32x4_mul(v1434, v1432), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1414, v1438), base.Simd_g_i32x4_mul(v1417, v1438), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k16)), v1424))
				v1450 = base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k17)
				v1461 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1431, v1450), base.Simd_g_i32x4_mul(v1434, v1450), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k11)), v1408), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k18)), v1424), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k19))
				v1472 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v1449, v1461, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k3)), int32(4)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k20)), base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v1449, v1461, base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k6)), base.Simd_g_const(&F_UpsampleRgba4444LinePair_SSE2__k21)))
				base.Simd_g_v128_store(m, v1387, v1398, v1472)
				if base.Ui32(v1388) < base.Ui32(int32(24)) {
					v1387 = v1387 + int32(16)
					v1388 = v1388 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1481 = int32(1)
			v1486 = F_memcpy(m, l6+v729<<(uint(v1481)%32), v1374, v1371<<(uint(v1481)%32))
			mBase = m.M
		}
	}
	m.G0 = v35 + int32(528)
	return
}

var F_UpsampleRgba4444LinePair_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleRgba4444LinePair_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleRgba4444LinePair_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleRgba4444LinePair_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleRgba4444LinePair_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleRgba4444LinePair_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleRgba4444LinePair_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleRgba4444LinePair_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleRgba4444LinePair_SSE2__k8 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleRgba4444LinePair_SSE2__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleRgba4444LinePair_SSE2__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleRgba4444LinePair_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleRgba4444LinePair_SSE2__k12 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleRgba4444LinePair_SSE2__k13 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleRgba4444LinePair_SSE2__k14 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleRgba4444LinePair_SSE2__k15 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleRgba4444LinePair_SSE2__k16 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_UpsampleRgba4444LinePair_SSE2__k17 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleRgba4444LinePair_SSE2__k18 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleRgba4444LinePair_SSE2__k19 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_UpsampleRgba4444LinePair_SSE2__k20 = [2]uint64{0xf0f0f0f0f0f0f0f, 0xf0f0f0f0f0f0f0f}
var F_UpsampleRgba4444LinePair_SSE2__k21 = [2]uint64{0xf0f0f0f0f0f0f0f0, 0xf0f0f0f0f0f0f0f0}

func F_UpsampleRgbaLinePair_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
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
	var v516 base.V128
	_ = v516
	var v520 base.V128
	_ = v520
	var v521 base.V128
	_ = v521
	var v524 base.V128
	_ = v524
	var v526 base.V128
	_ = v526
	var v527 base.V128
	_ = v527
	var v528 base.V128
	_ = v528
	var v530 base.V128
	_ = v530
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
	var v626 base.V128
	_ = v626
	var v630 base.V128
	_ = v630
	var v631 base.V128
	_ = v631
	var v634 base.V128
	_ = v634
	var v636 base.V128
	_ = v636
	var v637 base.V128
	_ = v637
	var v638 base.V128
	_ = v638
	var v640 base.V128
	_ = v640
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
	var v1403 base.V128
	_ = v1403
	var v1407 base.V128
	_ = v1407
	var v1408 base.V128
	_ = v1408
	var v1411 base.V128
	_ = v1411
	var v1413 base.V128
	_ = v1413
	var v1414 base.V128
	_ = v1414
	var v1415 base.V128
	_ = v1415
	var v1417 base.V128
	_ = v1417
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
	var v1520 base.V128
	_ = v1520
	var v1524 base.V128
	_ = v1524
	var v1525 base.V128
	_ = v1525
	var v1528 base.V128
	_ = v1528
	var v1530 base.V128
	_ = v1530
	var v1531 base.V128
	_ = v1531
	var v1532 base.V128
	_ = v1532
	var v1534 base.V128
	_ = v1534
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
	var v1633 base.V128
	_ = v1633
	var v1637 base.V128
	_ = v1637
	var v1638 base.V128
	_ = v1638
	var v1641 base.V128
	_ = v1641
	var v1643 base.V128
	_ = v1643
	var v1644 base.V128
	_ = v1644
	var v1645 base.V128
	_ = v1645
	var v1647 base.V128
	_ = v1647
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
	v165 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v171 = int32(base.Ui32(v165+v166)>>(uint(v164)%32)) + v164
	v174 = int32(base.Ui32(v171+v166) >> (uint(v164) % 32))
	v177 = int32(8)
	v179 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v183 = int32(base.Ui32(v179*int32(_a_F_UpsampleRgbaLinePair_SSE2_0)) >> (uint(v177) % 32))
	v184 = int32(base.Ui32(v174*int32(_a_F_UpsampleRgbaLinePair_SSE2_1))>>(uint(v177)%32)) + v183
	v186 = v184 + int32(-17685)
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_2)) {
		v193 = int32(0)
	} else {
		v193 = v162
	}
	if base.Ui32(v186) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_3)) {
		v196 = int32(base.Ui32(v186) >> (uint(int32(6)) % 32))
	} else {
		v196 = v193
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v196)
	v198 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v199 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v201 = int32(1)
	v204 = int32(base.Ui32(v198+v199)>>(uint(v201)%32)) + v201
	v207 = int32(base.Ui32(v204+v199) >> (uint(v201) % 32))
	v212 = int32(base.Ui32(v207*int32(_a_F_UpsampleRgbaLinePair_SSE2_4))>>(uint(int32(8))%32)) + v183
	v214 = v212 + int32(-14234)
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_5)) {
		v221 = int32(0)
	} else {
		v221 = int32(255)
	}
	if base.Ui32(v214) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_3)) {
		v224 = int32(base.Ui32(v214) >> (uint(int32(6)) % 32))
	} else {
		v224 = v221
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v224)
	v228 = int32(8)
	v235 = v183 - (int32(base.Ui32(v174*int32(_a_F_UpsampleRgbaLinePair_SSE2_6))>>(uint(v228)%32)) + int32(base.Ui32(v207*int32(_a_F_UpsampleRgbaLinePair_SSE2_7))>>(uint(v228)%32)))
	v237 = v235 + int32(_a_F_UpsampleRgbaLinePair_SSE2_8)
	if v235 < int32(-8708) {
		v244 = int32(0)
	} else {
		v244 = int32(255)
	}
	if base.Ui32(v237) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_3)) {
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
		v257 = int32(base.Ui32(v253*int32(_a_F_UpsampleRgbaLinePair_SSE2_0)) >> (uint(v256) % 32))
		v260 = int32(base.Ui32(v171+v165) >> (uint(int32(1)) % 32))
		v265 = v257 + int32(base.Ui32(v260*int32(_a_F_UpsampleRgbaLinePair_SSE2_1))>>(uint(v256)%32))
		v267 = v265 + int32(-17685)
		if base.Ui32(v265) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_2)) {
			v274 = int32(0)
		} else {
			v274 = v251
		}
		if base.Ui32(v267) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_3)) {
			v277 = int32(base.Ui32(v267) >> (uint(int32(6)) % 32))
		} else {
			v277 = v274
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v277)
		v281 = int32(base.Ui32(v204+v198) >> (uint(int32(1)) % 32))
		v286 = v257 + int32(base.Ui32(v281*int32(_a_F_UpsampleRgbaLinePair_SSE2_4))>>(uint(int32(8))%32))
		v288 = v286 + int32(-14234)
		if base.Ui32(v286) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_5)) {
			v295 = int32(0)
		} else {
			v295 = int32(255)
		}
		if base.Ui32(v288) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_3)) {
			v298 = int32(base.Ui32(v288) >> (uint(int32(6)) % 32))
		} else {
			v298 = v295
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v298)
		v302 = int32(8)
		v309 = v257 - (int32(base.Ui32(v281*int32(_a_F_UpsampleRgbaLinePair_SSE2_7))>>(uint(v302)%32)) + int32(base.Ui32(v260*int32(_a_F_UpsampleRgbaLinePair_SSE2_6))>>(uint(v302)%32)))
		v311 = v309 + int32(_a_F_UpsampleRgbaLinePair_SSE2_8)
		if v309 < int32(-8708) {
			v318 = int32(0)
		} else {
			v318 = int32(255)
		}
		if base.Ui32(v311) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_SSE2_3)) {
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
			v398 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k0)
			v400 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v388, v391), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v393, v394), v396), v398))
			v406 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v400, v388), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v400, v388), base.Simd_g_v128_and(v396, v393)), v398))
			v407 = base.Simd_g_i8x16_avgr_u(v382, v406)
			v413 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v400, v391), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v400, v391), base.Simd_g_v128_and(v396, v394)), v398))
			v414 = base.Simd_g_i8x16_avgr_u(v387, v413)
			v415 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k1)
			v416 = base.Simd_g_i8x16_shuffle2(v407, v414, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(80), v416)
			v419 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
			v420 = base.Simd_g_i8x16_shuffle2(v407, v414, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(64), v420)
			v423 = base.Simd_g_i8x16_avgr_u(v385, v413)
			v424 = base.Simd_g_i8x16_avgr_u(v390, v406)
			v426 = base.Simd_g_i8x16_shuffle2(v423, v424, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(16), v426)
			v430 = base.Simd_g_i8x16_shuffle2(v423, v424, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
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
			v468 = base.Simd_g_i8x16_shuffle2(v459, v466, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(112), v468)
			v472 = base.Simd_g_i8x16_shuffle2(v459, v466, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(96), v472)
			v475 = base.Simd_g_i8x16_avgr_u(v438, v465)
			v476 = base.Simd_g_i8x16_avgr_u(v443, v458)
			v478 = base.Simd_g_i8x16_shuffle2(v475, v476, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(48), v478)
			v482 = base.Simd_g_i8x16_shuffle2(v475, v476, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(32), v482)
			v499 = v364
			v500 = v381
			for {
				v509 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k7)
				v511 = int32(0)
				v512 = base.Simd_g_v128_load64_zero(m, l0+v338+v363+v500, v511)
				v513 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
				v514 = base.Simd_g_i8x16_shuffle2(v509, v512, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v516 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k8)
				v520 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k9)
				v521 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v514), v516), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v514), v516), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))
				v524 = base.Simd_g_v128_load64_zero(m, v329+v500, v511)
				v526 = base.Simd_g_i8x16_shuffle2(v509, v524, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v527 = base.Simd_g_i32x4_extend_low_i16x8_u(v526)
				v528 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k12)
				v530 = base.Simd_g_i32x4_extend_high_i16x8_u(v526)
				v537 = int32(6)
				v541 = base.Simd_g_v128_load64_zero(m, v331+v500, v511)
				v543 = base.Simd_g_i8x16_shuffle2(v509, v541, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v544 = base.Simd_g_i32x4_extend_low_i16x8_u(v543)
				v545 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k13)
				v547 = base.Simd_g_i32x4_extend_high_i16x8_u(v543)
				v556 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v521, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v527, v528), base.Simd_g_i32x4_mul(v530, v528), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k14)), v537), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v544, v545), base.Simd_g_i32x4_mul(v547, v545), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), v521), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k15)), v537))
				v557 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k16)
				v562 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k17)
				v574 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v521, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v544, v557), base.Simd_g_i32x4_mul(v547, v557), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v527, v562), base.Simd_g_i32x4_mul(v530, v562), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k18)), v537), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k19))
				v576 = base.Simd_g_i8x16_shuffle2(v556, v574, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v578 = base.Simd_g_i8x16_shuffle2(v556, v574, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
				v580 = base.Simd_g_i8x16_shuffle2(v576, v578, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v499, int32(16), v580)
				v584 = base.Simd_g_i8x16_shuffle2(v576, v578, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k23))
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
					v619 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k7)
					v621 = int32(0)
					v622 = base.Simd_g_v128_load64_zero(m, l1+v338+v363+v610, v621)
					v623 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
					v624 = base.Simd_g_i8x16_shuffle2(v619, v622, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
					v626 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k8)
					v630 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k9)
					v631 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v624), v626), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v624), v626), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))
					v634 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v610, v621)
					v636 = base.Simd_g_i8x16_shuffle2(v619, v634, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
					v637 = base.Simd_g_i32x4_extend_low_i16x8_u(v636)
					v638 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k12)
					v640 = base.Simd_g_i32x4_extend_high_i16x8_u(v636)
					v647 = int32(6)
					v651 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v610, v621)
					v653 = base.Simd_g_i8x16_shuffle2(v619, v651, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
					v654 = base.Simd_g_i32x4_extend_low_i16x8_u(v653)
					v655 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k13)
					v657 = base.Simd_g_i32x4_extend_high_i16x8_u(v653)
					v666 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v631, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v637, v638), base.Simd_g_i32x4_mul(v640, v638), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k14)), v647), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v654, v655), base.Simd_g_i32x4_mul(v657, v655), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), v631), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k15)), v647))
					v667 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k16)
					v672 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k17)
					v684 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v631, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v654, v667), base.Simd_g_i32x4_mul(v657, v667), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v637, v672), base.Simd_g_i32x4_mul(v640, v672), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k18)), v647), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k19))
					v686 = base.Simd_g_i8x16_shuffle2(v666, v684, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
					v688 = base.Simd_g_i8x16_shuffle2(v666, v684, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
					v690 = base.Simd_g_i8x16_shuffle2(v686, v688, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k21))
					base.Simd_g_v128_store(m, v609, int32(16), v690)
					v694 = base.Simd_g_i8x16_shuffle2(v686, v688, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k23))
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
		v1033 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k0)
		v1035 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1023, v1026), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1028, v1029), v1031), v1033))
		v1041 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1035, v1023), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1035, v1023), base.Simd_g_v128_and(v1031, v1028)), v1033))
		v1042 = base.Simd_g_i8x16_avgr_u(v1018, v1041)
		v1048 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1035, v1026), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1035, v1026), base.Simd_g_v128_and(v1031, v1029)), v1033))
		v1049 = base.Simd_g_i8x16_avgr_u(v1022, v1048)
		v1050 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k1)
		v1051 = base.Simd_g_i8x16_shuffle2(v1042, v1049, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(80), v1051)
		v1054 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
		v1055 = base.Simd_g_i8x16_shuffle2(v1042, v1049, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, int32(64), v1055)
		v1058 = base.Simd_g_i8x16_avgr_u(v1020, v1048)
		v1059 = base.Simd_g_i8x16_avgr_u(v1025, v1041)
		v1061 = base.Simd_g_i8x16_shuffle2(v1058, v1059, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(16), v1061)
		v1065 = base.Simd_g_i8x16_shuffle2(v1058, v1059, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
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
		v1348 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k1)
		v1349 = base.Simd_g_i8x16_shuffle2(v1340, v1347, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(112), v1349)
		v1352 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
		v1353 = base.Simd_g_i8x16_shuffle2(v1340, v1347, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, int32(96), v1353)
		v1356 = base.Simd_g_i8x16_avgr_u(v1319, v1346)
		v1357 = base.Simd_g_i8x16_avgr_u(v1324, v1339)
		v1359 = base.Simd_g_i8x16_shuffle2(v1356, v1357, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(48), v1359)
		v1363 = base.Simd_g_i8x16_shuffle2(v1356, v1357, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
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
				v1513 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k7)
				v1515 = int32(0)
				v1516 = base.Simd_g_v128_load64_zero(m, v1370+v1504, v1515)
				v1517 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
				v1518 = base.Simd_g_i8x16_shuffle2(v1513, v1516, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1520 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k8)
				v1524 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k9)
				v1525 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1518), v1520), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1518), v1520), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))
				v1528 = base.Simd_g_v128_load64_zero(m, v329+v1504, v1515)
				v1530 = base.Simd_g_i8x16_shuffle2(v1513, v1528, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1531 = base.Simd_g_i32x4_extend_low_i16x8_u(v1530)
				v1532 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k12)
				v1534 = base.Simd_g_i32x4_extend_high_i16x8_u(v1530)
				v1541 = int32(6)
				v1545 = base.Simd_g_v128_load64_zero(m, v331+v1504, v1515)
				v1547 = base.Simd_g_i8x16_shuffle2(v1513, v1545, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1548 = base.Simd_g_i32x4_extend_low_i16x8_u(v1547)
				v1549 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k13)
				v1551 = base.Simd_g_i32x4_extend_high_i16x8_u(v1547)
				v1560 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1525, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1531, v1532), base.Simd_g_i32x4_mul(v1534, v1532), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k14)), v1541), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1548, v1549), base.Simd_g_i32x4_mul(v1551, v1549), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), v1525), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k15)), v1541))
				v1561 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k16)
				v1566 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k17)
				v1578 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1525, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1548, v1561), base.Simd_g_i32x4_mul(v1551, v1561), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1531, v1566), base.Simd_g_i32x4_mul(v1534, v1566), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k18)), v1541), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k19))
				v1580 = base.Simd_g_i8x16_shuffle2(v1560, v1578, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1582 = base.Simd_g_i8x16_shuffle2(v1560, v1578, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
				v1584 = base.Simd_g_i8x16_shuffle2(v1580, v1582, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1503, int32(16), v1584)
				v1588 = base.Simd_g_i8x16_shuffle2(v1580, v1582, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k23))
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
				v1626 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k7)
				v1628 = int32(0)
				v1629 = base.Simd_g_v128_load64_zero(m, v1489+v1617, v1628)
				v1630 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
				v1631 = base.Simd_g_i8x16_shuffle2(v1626, v1629, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1633 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k8)
				v1637 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k9)
				v1638 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1631), v1633), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1631), v1633), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))
				v1641 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v1617, v1628)
				v1643 = base.Simd_g_i8x16_shuffle2(v1626, v1641, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1644 = base.Simd_g_i32x4_extend_low_i16x8_u(v1643)
				v1645 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k12)
				v1647 = base.Simd_g_i32x4_extend_high_i16x8_u(v1643)
				v1654 = int32(6)
				v1658 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v1617, v1628)
				v1660 = base.Simd_g_i8x16_shuffle2(v1626, v1658, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1661 = base.Simd_g_i32x4_extend_low_i16x8_u(v1660)
				v1662 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k13)
				v1664 = base.Simd_g_i32x4_extend_high_i16x8_u(v1660)
				v1673 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1638, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1644, v1645), base.Simd_g_i32x4_mul(v1647, v1645), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k14)), v1654), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1661, v1662), base.Simd_g_i32x4_mul(v1664, v1662), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), v1638), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k15)), v1654))
				v1674 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k16)
				v1679 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k17)
				v1691 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1638, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1661, v1674), base.Simd_g_i32x4_mul(v1664, v1674), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1644, v1679), base.Simd_g_i32x4_mul(v1647, v1679), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k18)), v1654), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k19))
				v1693 = base.Simd_g_i8x16_shuffle2(v1673, v1691, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1695 = base.Simd_g_i8x16_shuffle2(v1673, v1691, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
				v1697 = base.Simd_g_i8x16_shuffle2(v1693, v1695, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1616, int32(16), v1697)
				v1701 = base.Simd_g_i8x16_shuffle2(v1693, v1695, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k23))
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
				v1396 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k7)
				v1398 = int32(0)
				v1399 = base.Simd_g_v128_load64_zero(m, v1370+v1387, v1398)
				v1400 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k4)
				v1401 = base.Simd_g_i8x16_shuffle2(v1396, v1399, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1403 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k8)
				v1407 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k9)
				v1408 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1401), v1403), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1401), v1403), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))
				v1411 = base.Simd_g_v128_load64_zero(m, v329+v1387, v1398)
				v1413 = base.Simd_g_i8x16_shuffle2(v1396, v1411, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1414 = base.Simd_g_i32x4_extend_low_i16x8_u(v1413)
				v1415 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k12)
				v1417 = base.Simd_g_i32x4_extend_high_i16x8_u(v1413)
				v1424 = int32(6)
				v1428 = base.Simd_g_v128_load64_zero(m, v331+v1387, v1398)
				v1430 = base.Simd_g_i8x16_shuffle2(v1396, v1428, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1431 = base.Simd_g_i32x4_extend_low_i16x8_u(v1430)
				v1432 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k13)
				v1434 = base.Simd_g_i32x4_extend_high_i16x8_u(v1430)
				v1443 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1408, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1414, v1415), base.Simd_g_i32x4_mul(v1417, v1415), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k14)), v1424), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1431, v1432), base.Simd_g_i32x4_mul(v1434, v1432), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), v1408), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k15)), v1424))
				v1444 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k16)
				v1449 = base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k17)
				v1461 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1408, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1431, v1444), base.Simd_g_i32x4_mul(v1434, v1444), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1414, v1449), base.Simd_g_i32x4_mul(v1417, v1449), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k18)), v1424), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k19))
				v1463 = base.Simd_g_i8x16_shuffle2(v1443, v1461, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k6))
				v1465 = base.Simd_g_i8x16_shuffle2(v1443, v1461, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k3))
				v1467 = base.Simd_g_i8x16_shuffle2(v1463, v1465, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1386, int32(16), v1467)
				v1471 = base.Simd_g_i8x16_shuffle2(v1463, v1465, base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleRgbaLinePair_SSE2__k23))
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

var F_UpsampleRgbaLinePair_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleRgbaLinePair_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleRgbaLinePair_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleRgbaLinePair_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleRgbaLinePair_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleRgbaLinePair_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleRgbaLinePair_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleRgbaLinePair_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleRgbaLinePair_SSE2__k8 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleRgbaLinePair_SSE2__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleRgbaLinePair_SSE2__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleRgbaLinePair_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleRgbaLinePair_SSE2__k12 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleRgbaLinePair_SSE2__k13 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleRgbaLinePair_SSE2__k14 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleRgbaLinePair_SSE2__k15 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleRgbaLinePair_SSE2__k16 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleRgbaLinePair_SSE2__k17 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleRgbaLinePair_SSE2__k18 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_UpsampleRgbaLinePair_SSE2__k19 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_UpsampleRgbaLinePair_SSE2__k20 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_UpsampleRgbaLinePair_SSE2__k21 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_UpsampleRgbaLinePair_SSE2__k22 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_UpsampleRgbaLinePair_SSE2__k23 = [2]uint64{0x302808001008080, 0x706808005048080}

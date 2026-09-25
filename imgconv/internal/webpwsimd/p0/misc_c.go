//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_CollectHistogram_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v137 int32
	_ = v137
	var v146 int32
	_ = v146
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v157 int32
	_ = v157
	var v161 int32
	_ = v161
	var v164 int32
	_ = v164
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v184 int32
	_ = v184
	var v187 int32
	_ = v187
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v197 int32
	_ = v197
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v207 int32
	_ = v207
	var v210 int32
	_ = v210
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v220 int32
	_ = v220
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	var v230 int32
	_ = v230
	var v233 int32
	_ = v233
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
	var v243 int32
	_ = v243
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v266 int32
	_ = v266
	var v269 int32
	_ = v269
	var v270 int32
	_ = v270
	var v276 int32
	_ = v276
	var v279 int32
	_ = v279
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v289 int32
	_ = v289
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v299 int32
	_ = v299
	var v302 int32
	_ = v302
	var v308 int32
	_ = v308
	var v309 int32
	_ = v309
	var v312 int32
	_ = v312
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v322 int32
	_ = v322
	var v325 int32
	_ = v325
	var v331 int32
	_ = v331
	var v332 int32
	_ = v332
	var v335 int32
	_ = v335
	var v338 int32
	_ = v338
	var v339 int32
	_ = v339
	var v345 int32
	_ = v345
	var v348 int32
	_ = v348
	var v354 int32
	_ = v354
	var v355 int32
	_ = v355
	var v358 int32
	_ = v358
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v368 int32
	_ = v368
	var v371 int32
	_ = v371
	var v377 int32
	_ = v377
	var v378 int32
	_ = v378
	var v381 int32
	_ = v381
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v391 int32
	_ = v391
	var v394 int32
	_ = v394
	var v400 int32
	_ = v400
	var v401 int32
	_ = v401
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v408 int32
	_ = v408
	var v414 int32
	_ = v414
	var v417 int32
	_ = v417
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v427 int32
	_ = v427
	var v430 int32
	_ = v430
	var v431 int32
	_ = v431
	var v437 int32
	_ = v437
	var v440 int32
	_ = v440
	var v446 int32
	_ = v446
	var v447 int32
	_ = v447
	var v450 int32
	_ = v450
	var v453 int32
	_ = v453
	var v454 int32
	_ = v454
	var v460 int32
	_ = v460
	var v463 int32
	_ = v463
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v473 int32
	_ = v473
	var v476 int32
	_ = v476
	var v477 int32
	_ = v477
	var v483 int32
	_ = v483
	var v486 int32
	_ = v486
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v496 int32
	_ = v496
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v506 int32
	_ = v506
	var v509 int32
	_ = v509
	var v515 int32
	_ = v515
	var v516 int32
	_ = v516
	var v519 int32
	_ = v519
	var v522 int32
	_ = v522
	var v523 int32
	_ = v523
	var v530 int32
	_ = v530
	var v539 int32
	_ = v539
	var v540 int32
	_ = v540
	var v632 int32
	_ = v632
	var v634 int32
	_ = v634
	var v635 int32
	_ = v635
	var v639 int32
	_ = v639
	var v641 int32
	_ = v641
	var v642 int32
	_ = v642
	var v643 int32
	_ = v643
	var v645 int32
	_ = v645
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v653 int32
	_ = v653
	var v654 int32
	_ = v654
	var v655 int32
	_ = v655
	var v657 int32
	_ = v657
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v661 int32
	_ = v661
	var v662 int32
	_ = v662
	var v663 int32
	_ = v663
	var v665 int32
	_ = v665
	var v666 int32
	_ = v666
	var v667 int32
	_ = v667
	var v669 int32
	_ = v669
	var v670 int32
	_ = v670
	var v671 int32
	_ = v671
	var v673 int32
	_ = v673
	var v674 int32
	_ = v674
	var v675 int32
	_ = v675
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v681 int32
	_ = v681
	var v682 int32
	_ = v682
	var v683 int32
	_ = v683
	var v685 int32
	_ = v685
	var v686 int32
	_ = v686
	var v687 int32
	_ = v687
	var v689 int32
	_ = v689
	var v690 int32
	_ = v690
	var v691 int32
	_ = v691
	var v693 int32
	_ = v693
	var v694 int32
	_ = v694
	var v695 int32
	_ = v695
	var v697 int32
	_ = v697
	var v698 int32
	_ = v698
	var v699 int32
	_ = v699
	var v701 int32
	_ = v701
	var v702 int32
	_ = v702
	var v703 int32
	_ = v703
	var v705 int32
	_ = v705
	var v706 int32
	_ = v706
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v711 int32
	_ = v711
	var v713 int32
	_ = v713
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v717 int32
	_ = v717
	var v718 int32
	_ = v718
	var v719 int32
	_ = v719
	var v721 int32
	_ = v721
	var v722 int32
	_ = v722
	var v723 int32
	_ = v723
	var v725 int32
	_ = v725
	var v726 int32
	_ = v726
	var v727 int32
	_ = v727
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
	var v735 int32
	_ = v735
	var v737 int32
	_ = v737
	var v738 int32
	_ = v738
	var v739 int32
	_ = v739
	var v741 int32
	_ = v741
	var v742 int32
	_ = v742
	var v743 int32
	_ = v743
	var v745 int32
	_ = v745
	var v746 int32
	_ = v746
	var v747 int32
	_ = v747
	var v749 int32
	_ = v749
	var v750 int32
	_ = v750
	var v751 int32
	_ = v751
	var v753 int32
	_ = v753
	var v754 int32
	_ = v754
	var v755 int32
	_ = v755
	var v757 int32
	_ = v757
	var v758 int32
	_ = v758
	var v760 int32
	_ = v760
	var v763 int32
	_ = v763
	var v765 int32
	_ = v765
	var v766 int32
	_ = v766
	var v768 int32
	_ = v768
	var v769 int32
	_ = v769
	var v771 int32
	_ = v771
	var v772 int32
	_ = v772
	var v774 int32
	_ = v774
	var v775 int32
	_ = v775
	var v777 int32
	_ = v777
	var v778 int32
	_ = v778
	var v780 int32
	_ = v780
	var v781 int32
	_ = v781
	var v783 int32
	_ = v783
	var v784 int32
	_ = v784
	var v786 int32
	_ = v786
	var v787 int32
	_ = v787
	var v789 int32
	_ = v789
	var v790 int32
	_ = v790
	var v792 int32
	_ = v792
	var v793 int32
	_ = v793
	var v795 int32
	_ = v795
	var v796 int32
	_ = v796
	var v798 int32
	_ = v798
	var v799 int32
	_ = v799
	var v801 int32
	_ = v801
	var v802 int32
	_ = v802
	var v804 int32
	_ = v804
	var v805 int32
	_ = v805
	var v807 int32
	_ = v807
	var v808 int32
	_ = v808
	var v810 int32
	_ = v810
	var v811 int32
	_ = v811
	var v813 int32
	_ = v813
	var v814 int32
	_ = v814
	var v816 int32
	_ = v816
	var v817 int32
	_ = v817
	var v819 int32
	_ = v819
	var v820 int32
	_ = v820
	var v822 int32
	_ = v822
	var v823 int32
	_ = v823
	var v825 int32
	_ = v825
	var v826 int32
	_ = v826
	var v828 int32
	_ = v828
	var v829 int32
	_ = v829
	var v831 int32
	_ = v831
	var v832 int32
	_ = v832
	var v834 int32
	_ = v834
	var v835 int32
	_ = v835
	var v837 int32
	_ = v837
	var v838 int32
	_ = v838
	var v840 int32
	_ = v840
	var v841 int32
	_ = v841
	var v843 int32
	_ = v843
	var v844 int32
	_ = v844
	var v846 int32
	_ = v846
	var v847 int32
	_ = v847
	var v849 int32
	_ = v849
	var v850 int32
	_ = v850
	var v852 int32
	_ = v852
	var v853 int32
	_ = v853
	var v855 int32
	_ = v855
	var v856 int32
	_ = v856
	v8 = m.G0
	v10 = v8 - int32(160)
	m.G0 = v10
	base.MemoryFill(m, v10+int32(32), int32(0), int32(128))
	if l3 <= l2 {
	} else {
		v137 = m.G1
		v146 = l3 - l2
		v150 = v137 + int32(_a_F_CollectHistogram_C_0) + l2<<(uint(int32(2))%32)
		for {
			v151 = *(*int32)(unsafe.Add(mBase, uint32(v150)))
			v154 = m.G1
			v157 = *(*int32)(unsafe.Add(mBase, uint32(v154)+uint32(_c_F_CollectHistogram_C[0])))
			m.T0[v157].(func(*base.Module, int32, int32, int32))(m, l0+v151, l1+v151, v10)
			mBase = m.M
			v161 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10))))
			v164 = base.I32_extend16_s(v161) >> (uint(int32(15)) % 32)
			v170 = int32(base.Ui32((v161^v164-v164)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v171 = int32(31)
			if base.Ui32(v170) < base.Ui32(v171) {
				v174 = v170
			} else {
				v174 = v171
			}
			v177 = v10 + int32(32) + v174<<(uint(int32(2))%32)
			v178 = *(*int32)(unsafe.Add(mBase, uint32(v177)))
			*(*int32)(unsafe.Add(mBase, uint32(v177))) = v178 + int32(1)
			v184 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+2)))
			v187 = base.I32_extend16_s(v184) >> (uint(int32(15)) % 32)
			v193 = int32(base.Ui32((v184^v187-v187)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v194 = int32(31)
			if base.Ui32(v193) < base.Ui32(v194) {
				v197 = v193
			} else {
				v197 = v194
			}
			v200 = v10 + int32(32) + v197<<(uint(int32(2))%32)
			v201 = *(*int32)(unsafe.Add(mBase, uint32(v200)))
			*(*int32)(unsafe.Add(mBase, uint32(v200))) = v201 + int32(1)
			v207 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+4)))
			v210 = base.I32_extend16_s(v207) >> (uint(int32(15)) % 32)
			v216 = int32(base.Ui32((v207^v210-v210)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v217 = int32(31)
			if base.Ui32(v216) < base.Ui32(v217) {
				v220 = v216
			} else {
				v220 = v217
			}
			v223 = v10 + int32(32) + v220<<(uint(int32(2))%32)
			v224 = *(*int32)(unsafe.Add(mBase, uint32(v223)))
			*(*int32)(unsafe.Add(mBase, uint32(v223))) = v224 + int32(1)
			v230 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+6)))
			v233 = base.I32_extend16_s(v230) >> (uint(int32(15)) % 32)
			v239 = int32(base.Ui32((v230^v233-v233)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v240 = int32(31)
			if base.Ui32(v239) < base.Ui32(v240) {
				v243 = v239
			} else {
				v243 = v240
			}
			v246 = v10 + int32(32) + v243<<(uint(int32(2))%32)
			v247 = *(*int32)(unsafe.Add(mBase, uint32(v246)))
			*(*int32)(unsafe.Add(mBase, uint32(v246))) = v247 + int32(1)
			v253 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+8)))
			v256 = base.I32_extend16_s(v253) >> (uint(int32(15)) % 32)
			v262 = int32(base.Ui32((v253^v256-v256)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v263 = int32(31)
			if base.Ui32(v262) < base.Ui32(v263) {
				v266 = v262
			} else {
				v266 = v263
			}
			v269 = v10 + int32(32) + v266<<(uint(int32(2))%32)
			v270 = *(*int32)(unsafe.Add(mBase, uint32(v269)))
			*(*int32)(unsafe.Add(mBase, uint32(v269))) = v270 + int32(1)
			v276 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+10)))
			v279 = base.I32_extend16_s(v276) >> (uint(int32(15)) % 32)
			v285 = int32(base.Ui32((v276^v279-v279)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v286 = int32(31)
			if base.Ui32(v285) < base.Ui32(v286) {
				v289 = v285
			} else {
				v289 = v286
			}
			v292 = v10 + int32(32) + v289<<(uint(int32(2))%32)
			v293 = *(*int32)(unsafe.Add(mBase, uint32(v292)))
			*(*int32)(unsafe.Add(mBase, uint32(v292))) = v293 + int32(1)
			v299 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+12)))
			v302 = base.I32_extend16_s(v299) >> (uint(int32(15)) % 32)
			v308 = int32(base.Ui32((v299^v302-v302)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v309 = int32(31)
			if base.Ui32(v308) < base.Ui32(v309) {
				v312 = v308
			} else {
				v312 = v309
			}
			v315 = v10 + int32(32) + v312<<(uint(int32(2))%32)
			v316 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
			*(*int32)(unsafe.Add(mBase, uint32(v315))) = v316 + int32(1)
			v322 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+14)))
			v325 = base.I32_extend16_s(v322) >> (uint(int32(15)) % 32)
			v331 = int32(base.Ui32((v322^v325-v325)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v332 = int32(31)
			if base.Ui32(v331) < base.Ui32(v332) {
				v335 = v331
			} else {
				v335 = v332
			}
			v338 = v10 + int32(32) + v335<<(uint(int32(2))%32)
			v339 = *(*int32)(unsafe.Add(mBase, uint32(v338)))
			*(*int32)(unsafe.Add(mBase, uint32(v338))) = v339 + int32(1)
			v345 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+16)))
			v348 = base.I32_extend16_s(v345) >> (uint(int32(15)) % 32)
			v354 = int32(base.Ui32((v345^v348-v348)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v355 = int32(31)
			if base.Ui32(v354) < base.Ui32(v355) {
				v358 = v354
			} else {
				v358 = v355
			}
			v361 = v10 + int32(32) + v358<<(uint(int32(2))%32)
			v362 = *(*int32)(unsafe.Add(mBase, uint32(v361)))
			*(*int32)(unsafe.Add(mBase, uint32(v361))) = v362 + int32(1)
			v368 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+18)))
			v371 = base.I32_extend16_s(v368) >> (uint(int32(15)) % 32)
			v377 = int32(base.Ui32((v368^v371-v371)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v378 = int32(31)
			if base.Ui32(v377) < base.Ui32(v378) {
				v381 = v377
			} else {
				v381 = v378
			}
			v384 = v10 + int32(32) + v381<<(uint(int32(2))%32)
			v385 = *(*int32)(unsafe.Add(mBase, uint32(v384)))
			*(*int32)(unsafe.Add(mBase, uint32(v384))) = v385 + int32(1)
			v391 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+20)))
			v394 = base.I32_extend16_s(v391) >> (uint(int32(15)) % 32)
			v400 = int32(base.Ui32((v391^v394-v394)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v401 = int32(31)
			if base.Ui32(v400) < base.Ui32(v401) {
				v404 = v400
			} else {
				v404 = v401
			}
			v407 = v10 + int32(32) + v404<<(uint(int32(2))%32)
			v408 = *(*int32)(unsafe.Add(mBase, uint32(v407)))
			*(*int32)(unsafe.Add(mBase, uint32(v407))) = v408 + int32(1)
			v414 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+22)))
			v417 = base.I32_extend16_s(v414) >> (uint(int32(15)) % 32)
			v423 = int32(base.Ui32((v414^v417-v417)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v424 = int32(31)
			if base.Ui32(v423) < base.Ui32(v424) {
				v427 = v423
			} else {
				v427 = v424
			}
			v430 = v10 + int32(32) + v427<<(uint(int32(2))%32)
			v431 = *(*int32)(unsafe.Add(mBase, uint32(v430)))
			*(*int32)(unsafe.Add(mBase, uint32(v430))) = v431 + int32(1)
			v437 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+24)))
			v440 = base.I32_extend16_s(v437) >> (uint(int32(15)) % 32)
			v446 = int32(base.Ui32((v437^v440-v440)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v447 = int32(31)
			if base.Ui32(v446) < base.Ui32(v447) {
				v450 = v446
			} else {
				v450 = v447
			}
			v453 = v10 + int32(32) + v450<<(uint(int32(2))%32)
			v454 = *(*int32)(unsafe.Add(mBase, uint32(v453)))
			*(*int32)(unsafe.Add(mBase, uint32(v453))) = v454 + int32(1)
			v460 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+26)))
			v463 = base.I32_extend16_s(v460) >> (uint(int32(15)) % 32)
			v469 = int32(base.Ui32((v460^v463-v463)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v470 = int32(31)
			if base.Ui32(v469) < base.Ui32(v470) {
				v473 = v469
			} else {
				v473 = v470
			}
			v476 = v10 + int32(32) + v473<<(uint(int32(2))%32)
			v477 = *(*int32)(unsafe.Add(mBase, uint32(v476)))
			*(*int32)(unsafe.Add(mBase, uint32(v476))) = v477 + int32(1)
			v483 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+28)))
			v486 = base.I32_extend16_s(v483) >> (uint(int32(15)) % 32)
			v492 = int32(base.Ui32((v483^v486-v486)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v493 = int32(31)
			if base.Ui32(v492) < base.Ui32(v493) {
				v496 = v492
			} else {
				v496 = v493
			}
			v499 = v10 + int32(32) + v496<<(uint(int32(2))%32)
			v500 = *(*int32)(unsafe.Add(mBase, uint32(v499)))
			*(*int32)(unsafe.Add(mBase, uint32(v499))) = v500 + int32(1)
			v506 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v10)+30)))
			v509 = base.I32_extend16_s(v506) >> (uint(int32(15)) % 32)
			v515 = int32(base.Ui32((v506^v509-v509)&int32(_a_F_CollectHistogram_C_1)) >> (uint(int32(3)) % 32))
			v516 = int32(31)
			if base.Ui32(v515) < base.Ui32(v516) {
				v519 = v515
			} else {
				v519 = v516
			}
			v522 = v10 + int32(32) + v519<<(uint(int32(2))%32)
			v523 = *(*int32)(unsafe.Add(mBase, uint32(v522)))
			*(*int32)(unsafe.Add(mBase, uint32(v522))) = v523 + int32(1)
			v530 = v146 + int32(-1)
			if v530 != 0 {
				v146 = v530
				v150 = v150 + int32(4)
				continue
			} else {
				break
			}
			break
		}
	}
	v539 = v10 + int32(32)
	v540 = int32(0)
	v632 = *(*int32)(unsafe.Add(mBase, uint32(v539)+4))
	v634 = base.B2i32(v540 < v632)
	v635 = *(*int32)(unsafe.Add(mBase, uint32(v539)))
	v639 = *(*int32)(unsafe.Add(mBase, uint32(v539)+8))
	v641 = base.B2i32(v540 < v639)
	if v540 < v639 {
		v642 = int32(2)
	} else {
		v642 = v634 | base.B2i32(v635 < int32(1))
	}
	v643 = *(*int32)(unsafe.Add(mBase, uint32(v539)+12))
	v645 = base.B2i32(int32(0) < v643)
	if int32(0) < v643 {
		v646 = int32(3)
	} else {
		v646 = v642
	}
	v647 = *(*int32)(unsafe.Add(mBase, uint32(v539)+16))
	v649 = base.B2i32(int32(0) < v647)
	if int32(0) < v647 {
		v650 = int32(4)
	} else {
		v650 = v646
	}
	v651 = *(*int32)(unsafe.Add(mBase, uint32(v539)+20))
	v653 = base.B2i32(int32(0) < v651)
	if int32(0) < v651 {
		v654 = int32(5)
	} else {
		v654 = v650
	}
	v655 = *(*int32)(unsafe.Add(mBase, uint32(v539)+24))
	v657 = base.B2i32(int32(0) < v655)
	if int32(0) < v655 {
		v658 = int32(6)
	} else {
		v658 = v654
	}
	v659 = *(*int32)(unsafe.Add(mBase, uint32(v539)+28))
	v661 = base.B2i32(int32(0) < v659)
	if int32(0) < v659 {
		v662 = int32(7)
	} else {
		v662 = v658
	}
	v663 = *(*int32)(unsafe.Add(mBase, uint32(v539)+32))
	v665 = base.B2i32(int32(0) < v663)
	if int32(0) < v663 {
		v666 = int32(8)
	} else {
		v666 = v662
	}
	v667 = *(*int32)(unsafe.Add(mBase, uint32(v539)+36))
	v669 = base.B2i32(int32(0) < v667)
	if int32(0) < v667 {
		v670 = int32(9)
	} else {
		v670 = v666
	}
	v671 = *(*int32)(unsafe.Add(mBase, uint32(v539)+40))
	v673 = base.B2i32(int32(0) < v671)
	if int32(0) < v671 {
		v674 = int32(10)
	} else {
		v674 = v670
	}
	v675 = *(*int32)(unsafe.Add(mBase, uint32(v539)+44))
	v677 = base.B2i32(int32(0) < v675)
	if int32(0) < v675 {
		v678 = int32(11)
	} else {
		v678 = v674
	}
	v679 = *(*int32)(unsafe.Add(mBase, uint32(v539)+48))
	v681 = base.B2i32(int32(0) < v679)
	if int32(0) < v679 {
		v682 = int32(12)
	} else {
		v682 = v678
	}
	v683 = *(*int32)(unsafe.Add(mBase, uint32(v539)+52))
	v685 = base.B2i32(int32(0) < v683)
	if int32(0) < v683 {
		v686 = int32(13)
	} else {
		v686 = v682
	}
	v687 = *(*int32)(unsafe.Add(mBase, uint32(v539)+56))
	v689 = base.B2i32(int32(0) < v687)
	if int32(0) < v687 {
		v690 = int32(14)
	} else {
		v690 = v686
	}
	v691 = *(*int32)(unsafe.Add(mBase, uint32(v539)+60))
	v693 = base.B2i32(int32(0) < v691)
	if int32(0) < v691 {
		v694 = int32(15)
	} else {
		v694 = v690
	}
	v695 = *(*int32)(unsafe.Add(mBase, uint32(v539)+64))
	v697 = base.B2i32(int32(0) < v695)
	if int32(0) < v695 {
		v698 = int32(16)
	} else {
		v698 = v694
	}
	v699 = *(*int32)(unsafe.Add(mBase, uint32(v539)+68))
	v701 = base.B2i32(int32(0) < v699)
	if int32(0) < v699 {
		v702 = int32(17)
	} else {
		v702 = v698
	}
	v703 = *(*int32)(unsafe.Add(mBase, uint32(v539)+72))
	v705 = base.B2i32(int32(0) < v703)
	if int32(0) < v703 {
		v706 = int32(18)
	} else {
		v706 = v702
	}
	v707 = *(*int32)(unsafe.Add(mBase, uint32(v539)+76))
	v709 = base.B2i32(int32(0) < v707)
	if int32(0) < v707 {
		v710 = int32(19)
	} else {
		v710 = v706
	}
	v711 = *(*int32)(unsafe.Add(mBase, uint32(v539)+80))
	v713 = base.B2i32(int32(0) < v711)
	if int32(0) < v711 {
		v714 = int32(20)
	} else {
		v714 = v710
	}
	v715 = *(*int32)(unsafe.Add(mBase, uint32(v539)+84))
	v717 = base.B2i32(int32(0) < v715)
	if int32(0) < v715 {
		v718 = int32(21)
	} else {
		v718 = v714
	}
	v719 = *(*int32)(unsafe.Add(mBase, uint32(v539)+88))
	v721 = base.B2i32(int32(0) < v719)
	if int32(0) < v719 {
		v722 = int32(22)
	} else {
		v722 = v718
	}
	v723 = *(*int32)(unsafe.Add(mBase, uint32(v539)+92))
	v725 = base.B2i32(int32(0) < v723)
	if int32(0) < v723 {
		v726 = int32(23)
	} else {
		v726 = v722
	}
	v727 = *(*int32)(unsafe.Add(mBase, uint32(v539)+96))
	v729 = base.B2i32(int32(0) < v727)
	if int32(0) < v727 {
		v730 = int32(24)
	} else {
		v730 = v726
	}
	v731 = *(*int32)(unsafe.Add(mBase, uint32(v539)+100))
	v733 = base.B2i32(int32(0) < v731)
	if int32(0) < v731 {
		v734 = int32(25)
	} else {
		v734 = v730
	}
	v735 = *(*int32)(unsafe.Add(mBase, uint32(v539)+104))
	v737 = base.B2i32(int32(0) < v735)
	if int32(0) < v735 {
		v738 = int32(26)
	} else {
		v738 = v734
	}
	v739 = *(*int32)(unsafe.Add(mBase, uint32(v539)+108))
	v741 = base.B2i32(int32(0) < v739)
	if int32(0) < v739 {
		v742 = int32(27)
	} else {
		v742 = v738
	}
	v743 = *(*int32)(unsafe.Add(mBase, uint32(v539)+112))
	v745 = base.B2i32(int32(0) < v743)
	if int32(0) < v743 {
		v746 = int32(28)
	} else {
		v746 = v742
	}
	v747 = *(*int32)(unsafe.Add(mBase, uint32(v539)+116))
	v749 = base.B2i32(int32(0) < v747)
	if int32(0) < v747 {
		v750 = int32(29)
	} else {
		v750 = v746
	}
	v751 = *(*int32)(unsafe.Add(mBase, uint32(v539)+120))
	v753 = base.B2i32(int32(0) < v751)
	if int32(0) < v751 {
		v754 = int32(30)
	} else {
		v754 = v750
	}
	v755 = *(*int32)(unsafe.Add(mBase, uint32(v539)+124))
	v757 = base.B2i32(int32(0) < v755)
	if int32(0) < v755 {
		v758 = int32(31)
	} else {
		v758 = v754
	}
	*(*int32)(unsafe.Add(mBase, uint32(l4)+4)) = v758
	v760 = int32(0)
	if v760 < v635 {
		v763 = v635
	} else {
		v763 = v760
	}
	if v763 < v632 {
		v765 = v632
	} else {
		v765 = v763
	}
	if v540 < v632 {
		v766 = v765
	} else {
		v766 = v763
	}
	if v766 < v639 {
		v768 = v639
	} else {
		v768 = v766
	}
	if v540 < v639 {
		v769 = v768
	} else {
		v769 = v766
	}
	if v769 < v643 {
		v771 = v643
	} else {
		v771 = v769
	}
	if int32(0) < v643 {
		v772 = v771
	} else {
		v772 = v769
	}
	if v772 < v647 {
		v774 = v647
	} else {
		v774 = v772
	}
	if int32(0) < v647 {
		v775 = v774
	} else {
		v775 = v772
	}
	if v775 < v651 {
		v777 = v651
	} else {
		v777 = v775
	}
	if int32(0) < v651 {
		v778 = v777
	} else {
		v778 = v775
	}
	if v778 < v655 {
		v780 = v655
	} else {
		v780 = v778
	}
	if int32(0) < v655 {
		v781 = v780
	} else {
		v781 = v778
	}
	if v781 < v659 {
		v783 = v659
	} else {
		v783 = v781
	}
	if int32(0) < v659 {
		v784 = v783
	} else {
		v784 = v781
	}
	if v784 < v663 {
		v786 = v663
	} else {
		v786 = v784
	}
	if int32(0) < v663 {
		v787 = v786
	} else {
		v787 = v784
	}
	if v787 < v667 {
		v789 = v667
	} else {
		v789 = v787
	}
	if int32(0) < v667 {
		v790 = v789
	} else {
		v790 = v787
	}
	if v790 < v671 {
		v792 = v671
	} else {
		v792 = v790
	}
	if int32(0) < v671 {
		v793 = v792
	} else {
		v793 = v790
	}
	if v793 < v675 {
		v795 = v675
	} else {
		v795 = v793
	}
	if int32(0) < v675 {
		v796 = v795
	} else {
		v796 = v793
	}
	if v796 < v679 {
		v798 = v679
	} else {
		v798 = v796
	}
	if int32(0) < v679 {
		v799 = v798
	} else {
		v799 = v796
	}
	if v799 < v683 {
		v801 = v683
	} else {
		v801 = v799
	}
	if int32(0) < v683 {
		v802 = v801
	} else {
		v802 = v799
	}
	if v802 < v687 {
		v804 = v687
	} else {
		v804 = v802
	}
	if int32(0) < v687 {
		v805 = v804
	} else {
		v805 = v802
	}
	if v805 < v691 {
		v807 = v691
	} else {
		v807 = v805
	}
	if int32(0) < v691 {
		v808 = v807
	} else {
		v808 = v805
	}
	if v808 < v695 {
		v810 = v695
	} else {
		v810 = v808
	}
	if int32(0) < v695 {
		v811 = v810
	} else {
		v811 = v808
	}
	if v811 < v699 {
		v813 = v699
	} else {
		v813 = v811
	}
	if int32(0) < v699 {
		v814 = v813
	} else {
		v814 = v811
	}
	if v814 < v703 {
		v816 = v703
	} else {
		v816 = v814
	}
	if int32(0) < v703 {
		v817 = v816
	} else {
		v817 = v814
	}
	if v817 < v707 {
		v819 = v707
	} else {
		v819 = v817
	}
	if int32(0) < v707 {
		v820 = v819
	} else {
		v820 = v817
	}
	if v820 < v711 {
		v822 = v711
	} else {
		v822 = v820
	}
	if int32(0) < v711 {
		v823 = v822
	} else {
		v823 = v820
	}
	if v823 < v715 {
		v825 = v715
	} else {
		v825 = v823
	}
	if int32(0) < v715 {
		v826 = v825
	} else {
		v826 = v823
	}
	if v826 < v719 {
		v828 = v719
	} else {
		v828 = v826
	}
	if int32(0) < v719 {
		v829 = v828
	} else {
		v829 = v826
	}
	if v829 < v723 {
		v831 = v723
	} else {
		v831 = v829
	}
	if int32(0) < v723 {
		v832 = v831
	} else {
		v832 = v829
	}
	if v832 < v727 {
		v834 = v727
	} else {
		v834 = v832
	}
	if int32(0) < v727 {
		v835 = v834
	} else {
		v835 = v832
	}
	if v835 < v731 {
		v837 = v731
	} else {
		v837 = v835
	}
	if int32(0) < v731 {
		v838 = v837
	} else {
		v838 = v835
	}
	if v838 < v735 {
		v840 = v735
	} else {
		v840 = v838
	}
	if int32(0) < v735 {
		v841 = v840
	} else {
		v841 = v838
	}
	if v841 < v739 {
		v843 = v739
	} else {
		v843 = v841
	}
	if int32(0) < v739 {
		v844 = v843
	} else {
		v844 = v841
	}
	if v844 < v743 {
		v846 = v743
	} else {
		v846 = v844
	}
	if int32(0) < v743 {
		v847 = v846
	} else {
		v847 = v844
	}
	if v847 < v747 {
		v849 = v747
	} else {
		v849 = v847
	}
	if int32(0) < v747 {
		v850 = v849
	} else {
		v850 = v847
	}
	if v850 < v751 {
		v852 = v751
	} else {
		v852 = v850
	}
	if int32(0) < v751 {
		v853 = v852
	} else {
		v853 = v850
	}
	if v853 < v755 {
		v855 = v755
	} else {
		v855 = v853
	}
	if int32(0) < v755 {
		v856 = v855
	} else {
		v856 = v853
	}
	*(*int32)(unsafe.Add(mBase, uint32(l4))) = v856
	m.G0 = v10 + int32(160)
	return
}
func F_CombinedShannonEntropy_C(m *base.Module, l0 int32, l1 int32) int64 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v17 int64
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v40 int64
	_ = v40
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v45 int64
	_ = v45
	var v46 int64
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 int64
	_ = v49
	var v52 int32
	_ = v52
	var v58 int64
	_ = v58
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v64 int64
	_ = v64
	var v70 int32
	_ = v70
	var v76 int64
	_ = v76
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v81 int64
	_ = v81
	var v82 int64
	_ = v82
	var v85 int64
	_ = v85
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v96 int32
	_ = v96
	var v102 int64
	_ = v102
	var v103 int32
	_ = v103
	var v106 int32
	_ = v106
	var v107 int64
	_ = v107
	var v108 int64
	_ = v108
	var v111 int32
	_ = v111
	var v117 int64
	_ = v117
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v122 int64
	_ = v122
	var v123 int64
	_ = v123
	v4 = int32(0)
	v17 = int64(0)
	v18 = v4
	v19 = v4
	v20 = v4
	for {
		v26 = *(*int32)(unsafe.Add(mBase, uint32(l1+v18)))
		v28 = *(*int32)(unsafe.Add(mBase, uint32(l0+v18)))
		if v28 == int32(0) {
			if v26 == int32(0) {
				v85 = v17
				v86 = v19
				v87 = v20
			} else {
				if base.Ui32(int32(255)) < base.Ui32(v26) {
					v77 = m.G1
					v80 = *(*int32)(unsafe.Add(mBase, uint32(v77)+uint32(_c_F_CombinedShannonEntropy_C[0])))
					v81 = m.T0[v80].(func(*base.Module, int32) int64)(m, v26)
					mBase = m.M
					v82 = v81
				} else {
					v70 = m.G1
					v76 = *(*int64)(unsafe.Add(mBase, uint32(v70+int32(_a_F_CombinedShannonEntropy_C_0)+v26<<(uint(int32(3))%32))))
					v82 = v76
				}
				v85 = v82 + v17
				v86 = v26 + v19
				v87 = v20
			}
		} else {
			v31 = v26 + v28
			if base.Ui32(int32(255)) < base.Ui32(v28) {
				v41 = m.G1
				v44 = *(*int32)(unsafe.Add(mBase, uint32(v41)+uint32(_c_F_CombinedShannonEntropy_C[0])))
				v45 = m.T0[v44].(func(*base.Module, int32) int64)(m, v28)
				mBase = m.M
				v46 = v45
			} else {
				v34 = m.G1
				v40 = *(*int64)(unsafe.Add(mBase, uint32(v34+int32(_a_F_CombinedShannonEntropy_C_0)+v28<<(uint(int32(3))%32))))
				v46 = v40
			}
			v47 = v28 + v20
			v48 = v31 + v19
			v49 = v46 + v17
			if base.Ui32(int32(255)) < base.Ui32(v31) {
				v60 = m.G1
				v63 = *(*int32)(unsafe.Add(mBase, uint32(v60)+uint32(_c_F_CombinedShannonEntropy_C[0])))
				v64 = m.T0[v63].(func(*base.Module, int32) int64)(m, v31)
				mBase = m.M
				v85 = v49 + v64
				v86 = v48
				v87 = v47
			} else {
				v52 = m.G1
				v58 = *(*int64)(unsafe.Add(mBase, uint32(v52+int32(_a_F_CombinedShannonEntropy_C_0)+v31<<(uint(int32(3))%32))))
				v85 = v49 + v58
				v86 = v48
				v87 = v47
			}
		}
		v91 = v18 + int32(4)
		if v91 != int32(1024) {
			v17 = v85
			v18 = v91
			v19 = v86
			v20 = v87
			continue
		} else {
			break
		}
		break
	}
	if base.Ui32(int32(255)) < base.Ui32(v87) {
		v103 = m.G1
		v106 = *(*int32)(unsafe.Add(mBase, uint32(v103)+uint32(_c_F_CombinedShannonEntropy_C[0])))
		v107 = m.T0[v106].(func(*base.Module, int32) int64)(m, v87)
		mBase = m.M
		v108 = v107
	} else {
		v96 = m.G1
		v102 = *(*int64)(unsafe.Add(mBase, uint32(v96+int32(_a_F_CombinedShannonEntropy_C_0)+v87<<(uint(int32(3))%32))))
		v108 = v102
	}
	if base.Ui32(int32(255)) < base.Ui32(v86) {
		v118 = m.G1
		v121 = *(*int32)(unsafe.Add(mBase, uint32(v118)+uint32(_c_F_CombinedShannonEntropy_C[0])))
		v122 = m.T0[v121].(func(*base.Module, int32) int64)(m, v86)
		mBase = m.M
		v123 = v122
	} else {
		v111 = m.G1
		v117 = *(*int64)(unsafe.Add(mBase, uint32(v111+int32(_a_F_CombinedShannonEntropy_C_0)+v86<<(uint(int32(3))%32))))
		v123 = v117
	}
	return v108 - v85 + v123
}
func F_CompareHuffmanTrees(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	if base.Ui32(v5) <= base.Ui32(v6) {
		if base.Ui32(v6) <= base.Ui32(v5) {
			v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v16 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
			if v15 < v16 {
				v18 = int32(-1)
			} else {
				v18 = int32(1)
			}
			return v18
		} else {
			return int32(1)
		}
	} else {
		return int32(-1)
	}
}
func F_ConvertRowsToUV(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v152 int32
	_ = v152
	var v159 int32
	_ = v159
	if l3 < int32(1) {
	} else {
		v16 = l4 + int32(8)
		v17 = l0
		v18 = l1
		v19 = l2
		v20 = l3
		for {
			v29 = *(*int32)(unsafe.Add(mBase, uint32(l4)+228))
			v30 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
			v31 = int32(2)
			v33 = v16 + v30<<(uint(v31)%32)
			v34 = *(*int32)(unsafe.Add(mBase, uint32(v33)))
			v35 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
			v39 = *(*int32)(unsafe.Add(mBase, uint32(v16+v35<<(uint(v31)%32))))
			v40 = v34 - v39
			*(*int32)(unsafe.Add(mBase, uint32(v33))) = v40 & int32(2147483647)
			v45 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
			v47 = v45 + int32(1)
			if v47 == int32(55) {
				v50 = int32(0)
			} else {
				v50 = v47
			}
			*(*int32)(unsafe.Add(mBase, uint32(l4)+4)) = v50
			v53 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
			v55 = v53 + int32(1)
			if v55 == int32(55) {
				v58 = int32(0)
			} else {
				v58 = v55
			}
			*(*int32)(unsafe.Add(mBase, uint32(l4))) = v58
			v60 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v17))))
			v63 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v17)+2)))
			v67 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v17)+4)))
			v82 = (v60*int32(-9719) + v63*int32(-19081) + v67*int32(_a_F_ConvertRowsToUV_0) + v29*(v40<<(uint(int32(1))%32)>>(uint(int32(14))%32))>>(uint(int32(8))%32) + int32(33685504)) >> (uint(int32(18)) % 32)
			v83 = int32(0)
			if v83 < v82 {
				v86 = v82
			} else {
				v86 = v83
			}
			v87 = int32(255)
			if v86 < v87 {
				v90 = v86
			} else {
				v90 = v87
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v18))) = uint8(v90)
			v92 = *(*int32)(unsafe.Add(mBase, uint32(l4)+228))
			v93 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
			v94 = int32(2)
			v96 = v16 + v93<<(uint(v94)%32)
			v97 = *(*int32)(unsafe.Add(mBase, uint32(v96)))
			v98 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
			v102 = *(*int32)(unsafe.Add(mBase, uint32(v16+v98<<(uint(v94)%32))))
			v103 = v97 - v102
			*(*int32)(unsafe.Add(mBase, uint32(v96))) = v103 & int32(2147483647)
			v108 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
			v110 = v108 + int32(1)
			if v110 == int32(55) {
				v113 = int32(0)
			} else {
				v113 = v110
			}
			*(*int32)(unsafe.Add(mBase, uint32(l4))) = v113
			v116 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
			v118 = v116 + int32(1)
			if v118 == int32(55) {
				v121 = int32(0)
			} else {
				v121 = v118
			}
			*(*int32)(unsafe.Add(mBase, uint32(l4)+4)) = v121
			v142 = (v60*int32(_a_F_ConvertRowsToUV_0) + v63*int32(-24116) + v67*int32(-4684) + v92*(v103<<(uint(int32(1))%32)>>(uint(int32(14))%32))>>(uint(int32(8))%32) + int32(33685504)) >> (uint(int32(18)) % 32)
			v143 = int32(0)
			if v143 < v142 {
				v146 = v142
			} else {
				v146 = v143
			}
			v147 = int32(255)
			if v146 < v147 {
				v150 = v146
			} else {
				v150 = v147
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v19))) = uint8(v150)
			v152 = int32(1)
			v159 = v20 + int32(-1)
			if v159 != 0 {
				v17 = v17 + int32(8)
				v18 = v18 + v152
				v19 = v19 + v152
				v20 = v159
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_Copy4x4_C(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	*(*int32)(unsafe.Add(mBase, uint32(l1))) = v3
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+32)) = v5
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+64)) = v7
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+96)) = v9
	return
}
func F_CostManagerClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v20 int32
	_ = v20
	var v24 int32
	_ = v24
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	if l0 == int32(0) {
	} else {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[0])))
		F_free(m, v8)
		mBase = m.M
		v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		F_free(m, v10)
		mBase = m.M
		v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v12 == int32(0) {
		} else {
			v20 = v12
			for {
				v24 = *(*int32)(unsafe.Add(mBase, uint32(v20)+24))
				if base.Ui32(v20) < base.Ui32(l0+int32(_a_F_CostManagerClear_0)) {
					F_free(m, v20)
					mBase = m.M
				} else {
					if base.Ui32(v20) <= base.Ui32(l0+int32(_a_F_CostManagerClear_1)) {
					} else {
						F_free(m, v20)
						mBase = m.M
					}
				}
				if v24 != 0 {
					v20 = v24
					continue
				} else {
					break
				}
				break
			}
		}
		v33 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v33
		v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[1])))
		if v35 == v33 {
		} else {
			v43 = v35
			for {
				v47 = *(*int32)(unsafe.Add(mBase, uint32(v43)+24))
				if base.Ui32(v43) < base.Ui32(l0+int32(_a_F_CostManagerClear_0)) {
					F_free(m, v43)
					mBase = m.M
				} else {
					if base.Ui32(v43) <= base.Ui32(l0+int32(_a_F_CostManagerClear_1)) {
					} else {
						F_free(m, v43)
						mBase = m.M
					}
				}
				if v47 != 0 {
					v43 = v47
					continue
				} else {
					break
				}
				break
			}
		}
		base.MemoryFill(m, l0, int32(0), int32(_a_F_CostManagerClear_2))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[2]))) = l0 + int32(_a_F_CostManagerClear_1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[3]))) = l0 + int32(_a_F_CostManagerClear_3)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[4]))) = l0 + int32(_a_F_CostManagerClear_4)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[5]))) = l0 + int32(_a_F_CostManagerClear_5)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[6]))) = l0 + int32(_a_F_CostManagerClear_6)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[7]))) = l0 + int32(_a_F_CostManagerClear_7)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[8]))) = l0 + int32(_a_F_CostManagerClear_8)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[9]))) = l0 + int32(_a_F_CostManagerClear_9)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[10]))) = l0 + int32(_a_F_CostManagerClear_10)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_CostManagerClear[11]))) = l0 + int32(_a_F_CostManagerClear_0)
	}
	return
}
func F_calloc(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int64
	_ = v8
	var v9 int32
	_ = v9
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v27 int32
	_ = v27
	var v32 int32
	_ = v32
	var v43 int32
	_ = v43
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v119 int64
	_ = v119
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v135 int32
	_ = v135
	if l0 != 0 {
		v8 = base.I64_extend_i32_u(l0) * base.I64_extend_i32_u(l1)
		v9 = base.I32_wrap_i64(v8)
		if base.Ui32(l1|l0) < base.Ui32(int32(_a_F_calloc_0)) {
			v20 = v9
		} else {
			if base.I32_wrap_i64(int64(base.Ui64(v8)>>(uint(int64(32))%64))) != int32(0) {
				v19 = int32(-1)
			} else {
				v19 = v9
			}
			v20 = v19
		}
	} else {
		v20 = int32(0)
	}
	v22 = F_dlmalloc(m, v20)
	mBase = m.M
	if v22 == int32(0) {
	} else {
		v27 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v22+int32(-4)))))
		if v27&int32(3) == int32(0) {
		} else {
			v32 = int32(0)
			if base.Ui32(v20) < base.Ui32(int32(33)) {
				if v20 == int32(0) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v22))) = uint8(v32)
					v43 = v22 + v20
					*(*uint8)(unsafe.Add(mBase, uint32(v43+int32(-1)))) = uint8(v32)
					if base.Ui32(v20) < base.Ui32(int32(3)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v22)+2)) = uint8(v32)
						*(*uint8)(unsafe.Add(mBase, uint32(v22)+1)) = uint8(v32)
						*(*uint8)(unsafe.Add(mBase, uint32(v43+int32(-3)))) = uint8(v32)
						*(*uint8)(unsafe.Add(mBase, uint32(v43+int32(-2)))) = uint8(v32)
						if base.Ui32(v20) < base.Ui32(int32(7)) {
						} else {
							*(*uint8)(unsafe.Add(mBase, uint32(v22)+3)) = uint8(v32)
							*(*uint8)(unsafe.Add(mBase, uint32(v43+int32(-4)))) = uint8(v32)
							if base.Ui32(v20) < base.Ui32(int32(9)) {
							} else {
								v65 = int32(0)
								v68 = (v65 - v22) & int32(3)
								v69 = v22 + v68
								*(*int32)(unsafe.Add(mBase, uint32(v69))) = v65
								v77 = (v20 - v68) & int32(60)
								v78 = v69 + v77
								*(*int32)(unsafe.Add(mBase, uint32(v78+int32(-4)))) = v65
								if base.Ui32(v77) < base.Ui32(int32(9)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v69)+8)) = v65
									*(*int32)(unsafe.Add(mBase, uint32(v69)+4)) = v65
									*(*int32)(unsafe.Add(mBase, uint32(v78+int32(-8)))) = v65
									*(*int32)(unsafe.Add(mBase, uint32(v78+int32(-12)))) = v65
									if base.Ui32(v77) < base.Ui32(int32(25)) {
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(v69)+24)) = v65
										*(*int32)(unsafe.Add(mBase, uint32(v69)+20)) = v65
										*(*int32)(unsafe.Add(mBase, uint32(v69)+16)) = v65
										*(*int32)(unsafe.Add(mBase, uint32(v69)+12)) = v65
										*(*int32)(unsafe.Add(mBase, uint32(v78+int32(-16)))) = v65
										*(*int32)(unsafe.Add(mBase, uint32(v78+int32(-20)))) = v65
										*(*int32)(unsafe.Add(mBase, uint32(v78+int32(-24)))) = v65
										*(*int32)(unsafe.Add(mBase, uint32(v78+int32(-28)))) = v65
										v113 = v69&int32(4) | int32(24)
										v114 = v77 - v113
										if base.Ui32(v114) < base.Ui32(int32(32)) {
										} else {
											v119 = base.I64_extend_i32_u(v65) * int64(4294967297)
											v122 = v114
											v123 = v69 + v113
											for {
												*(*int64)(unsafe.Add(mBase, uint32(v123)+24)) = v119
												*(*int64)(unsafe.Add(mBase, uint32(v123)+16)) = v119
												*(*int64)(unsafe.Add(mBase, uint32(v123)+8)) = v119
												*(*int64)(unsafe.Add(mBase, uint32(v123))) = v119
												v135 = v122 + int32(-32)
												if base.Ui32(int32(31)) < base.Ui32(v135) {
													v122 = v135
													v123 = v123 + int32(32)
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
				base.MemoryFill(m, v22, v32, v20)
			}
		}
	}
	return v22
}
func F_checkint(m *base.Module, l0 int64) int32 {
	var v10 int32
	_ = v10
	var v21 int64
	_ = v21
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	v10 = base.I32_wrap_i64(int64(base.Ui64(l0)>>(uint(int64(52))%64))) & int32(2047)
	if base.Ui32(v10) < base.Ui32(int32(1023)) {
		v33 = int32(0)
	} else {
		if base.Ui32(int32(1075)) < base.Ui32(v10) {
			v33 = int32(2)
		} else {
			v21 = int64(1) << (uint(base.I64_extend_i32_u(int32(1075)-v10)) % 64)
			if (v21+int64(-1))&l0 != int64(0) {
				v33 = int32(0)
			} else {
				if v21&l0 == int64(0) {
					v32 = int32(2)
				} else {
					v32 = int32(1)
				}
				v33 = v32
			}
		}
	}
	return v33
}

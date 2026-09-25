//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8EncDeleteAlpha(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	v4 = int32(1)
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncDeleteAlpha[0])))
	if v5 < v4 {
		v20 = v4
	} else {
		v9 = l0 + int32(384)
		v10 = m.G1
		v13 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8EncDeleteAlpha[1])))
		v14 = m.T0[v13].(func(*base.Module, int32) int32)(m, v9)
		mBase = m.M
		v15 = m.G1
		v18 = *(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_VP8EncDeleteAlpha[2])))
		m.T0[v18].(func(*base.Module, int32))(m, v9)
		mBase = m.M
		v20 = v14
	}
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+376))
	F_free(m, v22)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(l0)+380)) = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+372)) = int64(0)
	return v20
}
func F_VP8EncDspCostInit(m *base.Module) {
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
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	v4 = m.G1
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_VP8EncDspCostInit[0])))
	v8 = m.G13
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8EncDspCostInit[1]))) = v11 + int32(253)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8EncDspCostInit[2]))) = v11 + int32(254)
		if v9 == int32(0) {
		} else {
			v25 = int32(0)
			v26 = m.T0[v9].(func(*base.Module, int32) int32)(m, v25)
			mBase = m.M
			if v26 == v25 {
			} else {
				v30 = m.G2
				v31 = m.G77
				*(*int32)(unsafe.Add(mBase, uint32(v31))) = v30 + int32(251)
				v35 = m.G78
				*(*int32)(unsafe.Add(mBase, uint32(v35))) = v30 + int32(252)
			}
		}
		v39 = m.G1
		v42 = m.G13
		v43 = *(*int32)(unsafe.Add(mBase, uint32(v42)))
		*(*int32)(unsafe.Add(mBase, uint32(v39)+uint32(_c_F_VP8EncDspCostInit[0]))) = v43
	}
	return
}
func F_VP8EncDspCostInitSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	v2 = m.G2
	v3 = m.G77
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(251)
	v7 = m.G78
	*(*int32)(unsafe.Add(mBase, uint32(v7))) = v2 + int32(252)
	return
}
func F_VP8EncDspInitSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v71 int32
	_ = v71
	v2 = m.G2
	v3 = m.G58
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(208)
	v7 = m.G59
	*(*int32)(unsafe.Add(mBase, uint32(v7))) = v2 + int32(209)
	v11 = m.G60
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v2 + int32(210)
	v15 = m.G61
	*(*int32)(unsafe.Add(mBase, uint32(v15))) = v2 + int32(211)
	v19 = m.G62
	*(*int32)(unsafe.Add(mBase, uint32(v19))) = v2 + int32(212)
	v23 = m.G63
	*(*int32)(unsafe.Add(mBase, uint32(v23))) = v2 + int32(213)
	v27 = m.G64
	*(*int32)(unsafe.Add(mBase, uint32(v27))) = v2 + int32(214)
	v31 = m.G65
	*(*int32)(unsafe.Add(mBase, uint32(v31))) = v2 + int32(215)
	v35 = m.G66
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2 + int32(216)
	v39 = m.G67
	*(*int32)(unsafe.Add(mBase, uint32(v39))) = v2 + int32(217)
	v43 = m.G68
	*(*int32)(unsafe.Add(mBase, uint32(v43))) = v2 + int32(218)
	v47 = m.G69
	*(*int32)(unsafe.Add(mBase, uint32(v47))) = v2 + int32(219)
	v51 = m.G70
	*(*int32)(unsafe.Add(mBase, uint32(v51))) = v2 + int32(220)
	v55 = m.G71
	*(*int32)(unsafe.Add(mBase, uint32(v55))) = v2 + int32(221)
	v59 = m.G72
	*(*int32)(unsafe.Add(mBase, uint32(v59))) = v2 + int32(222)
	v63 = m.G73
	*(*int32)(unsafe.Add(mBase, uint32(v63))) = v2 + int32(223)
	v67 = m.G74
	*(*int32)(unsafe.Add(mBase, uint32(v67))) = v2 + int32(224)
	v71 = m.G75
	*(*int32)(unsafe.Add(mBase, uint32(v71))) = v2 + int32(225)
	return
}
func F_VP8EncDspInitSSE41(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	v2 = m.G2
	v3 = m.G62
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(226)
	v7 = m.G59
	*(*int32)(unsafe.Add(mBase, uint32(v7))) = v2 + int32(227)
	v11 = m.G63
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v2 + int32(228)
	v15 = m.G64
	*(*int32)(unsafe.Add(mBase, uint32(v15))) = v2 + int32(229)
	v19 = m.G73
	*(*int32)(unsafe.Add(mBase, uint32(v19))) = v2 + int32(230)
	v23 = m.G74
	*(*int32)(unsafe.Add(mBase, uint32(v23))) = v2 + int32(231)
	return
}
func F_VP8EncFinishAlpha(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v5 int32
	_ = v5
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v23 int32
	_ = v23
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)+372))
	if v2 == int32(0) {
		v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
		v23 = F_WebPReportProgress(m, v17, v18+int32(20), l0+int32(368))
		mBase = m.M
		return v23
	} else {
		v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncFinishAlpha[0])))
		if v5 < int32(1) {
			v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
			v23 = F_WebPReportProgress(m, v17, v18+int32(20), l0+int32(368))
			mBase = m.M
			return v23
		} else {
			v10 = m.G1
			v13 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8EncFinishAlpha[1])))
			v14 = m.T0[v13].(func(*base.Module, int32) int32)(m, l0+int32(384))
			mBase = m.M
			if v14 != 0 {
				v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
				v23 = F_WebPReportProgress(m, v17, v18+int32(20), l0+int32(368))
				mBase = m.M
				return v23
			} else {
				return int32(0)
			}
		}
	}
}
func F_VP8EncInitAlpha(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v9 int32
	_ = v9
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	F_WebPInitAlphaProcessing(m)
	mBase = m.M
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v5 = F_WebPPictureHasTransparency(m, v4)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(l0)+376)) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+372)) = v5
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncInitAlpha[0])))
	if v9 < int32(1) {
	} else {
		v14 = m.G1
		v17 = *(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[1])))
		m.T0[v17].(func(*base.Module, int32))(m, l0+int32(384))
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0)+400)) = int32(0)
		v21 = m.G2
		*(*int32)(unsafe.Add(mBase, uint32(l0)+392)) = v21 + int32(327)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+396)) = l0
	}
	return
}
func F_VP8EncStartAlpha(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v6 int32
	_ = v6
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+372))
	if v3 != 0 {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncStartAlpha[0])))
		if v6 < int32(1) {
			v29 = F_CompressAlphaJob(m, l0, l0)
			mBase = m.M
			return v29
		} else {
			v10 = l0 + int32(384)
			v11 = m.G1
			v14 = *(*int32)(unsafe.Add(mBase, uint32(v11)+uint32(_c_F_VP8EncStartAlpha[1])))
			v15 = m.T0[v14].(func(*base.Module, int32) int32)(m, v10)
			mBase = m.M
			if v15 != 0 {
				v22 = m.G1
				v25 = *(*int32)(unsafe.Add(mBase, uint32(v22)+uint32(_c_F_VP8EncStartAlpha[2])))
				m.T0[v25].(func(*base.Module, int32))(m, v10)
				mBase = m.M
				return int32(1)
			} else {
				v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v18 = *(*int32)(unsafe.Add(mBase, uint32(v16)+92))
				if v18 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v16)+92)) = int32(1)
				}
				return int32(0)
			}
		}
	} else {
		return int32(1)
	}
}

//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_UpdateChroma(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v48 int32
	_ = v48
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v133 int32
	_ = v133
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v156 int32
	_ = v156
	if l4 < int32(13) {
		v25 = int32(2)
	} else {
		v25 = int32(14) - l4
	}
	v26 = v25 + l4
	v27 = int32(1)
	if v27 < l3 {
		v30 = l3
	} else {
		v30 = v27
	}
	v34 = l3 << (uint(int32(2)) % 32)
	v37 = l3 << (uint(int32(3)) % 32)
	v44 = l2
	v45 = int32(0)
	v48 = v30
	for {
		v61 = l1 + v45
		v62 = int32(2)
		v64 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v61+v62))))
		v65 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v61))))
		v66 = l0 + v45
		v69 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v66+v62))))
		v70 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v66))))
		v71 = F_SharpYuvGammaToLinear(m, v70, v26, l5)
		mBase = m.M
		v72 = F_SharpYuvGammaToLinear(m, v69, v26, l5)
		mBase = m.M
		v74 = F_SharpYuvGammaToLinear(m, v65, v26, l5)
		mBase = m.M
		v76 = F_SharpYuvGammaToLinear(m, v64, v26, l5)
		mBase = m.M
		v82 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v71+v72+v74+v76+v62)>>(uint(v62)%32)), v26, l5)
		mBase = m.M
		v83 = l1 + v34 + v45
		v86 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v83+v62))))
		v87 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v83))))
		v88 = l0 + v34 + v45
		v91 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v88+v62))))
		v92 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v88))))
		v93 = F_SharpYuvGammaToLinear(m, v92, v26, l5)
		mBase = m.M
		v94 = F_SharpYuvGammaToLinear(m, v91, v26, l5)
		mBase = m.M
		v96 = F_SharpYuvGammaToLinear(m, v87, v26, l5)
		mBase = m.M
		v98 = F_SharpYuvGammaToLinear(m, v86, v26, l5)
		mBase = m.M
		v104 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v93+v94+v96+v98+v62)>>(uint(v62)%32)), v26, l5)
		mBase = m.M
		v105 = l1 + v37 + v45
		v108 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v105+v62))))
		v109 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v105))))
		v110 = l0 + v37 + v45
		v113 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v110+v62))))
		v121 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v110))))
		v122 = F_SharpYuvGammaToLinear(m, v121, v26, l5)
		mBase = m.M
		v123 = F_SharpYuvGammaToLinear(m, v113, v26, l5)
		mBase = m.M
		v125 = F_SharpYuvGammaToLinear(m, v109, v26, l5)
		mBase = m.M
		v127 = F_SharpYuvGammaToLinear(m, v108, v26, l5)
		mBase = m.M
		v133 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v122+v123+v125+v127+v62)>>(uint(v62)%32)), v26, l5)
		mBase = m.M
		v142 = base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v82)*int64(13933)+base.I64_extend_i32_u(v104)*int64(46871)+base.I64_extend_i32_u(v133)*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64)))
		v143 = v82 - v142
		*(*uint16)(unsafe.Add(mBase, uint32(v44))) = uint16(v143)
		v146 = v104 - v142
		*(*uint16)(unsafe.Add(mBase, uint32(v44+l3<<(uint(int32(1))%32)))) = uint16(v146)
		v149 = v133 - v142
		*(*uint16)(unsafe.Add(mBase, uint32(v44+v34))) = uint16(v149)
		v156 = v48 + int32(-1)
		if v156 != 0 {
			v44 = v44 + v62
			v45 = v45 + int32(4)
			v48 = v156
			continue
		} else {
			break
		}
		break
	}
	return
}

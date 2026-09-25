//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8IteratorSaveBoundary(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v131 int32
	_ = v131
	var v133 base.V128
	_ = v133
	var v134 int32
	_ = v134
	var v136 int32
	_ = v136
	var v138 base.V128
	_ = v138
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)+40))
	if v9+int32(-1) <= v7 {
	} else {
		v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v14 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+15)))
		*(*uint8)(unsafe.Add(mBase, uint32(v13))) = uint8(v14)
		v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v17 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+47)))
		*(*uint8)(unsafe.Add(mBase, uint32(v16)+1)) = uint8(v17)
		v19 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+79)))
		*(*uint8)(unsafe.Add(mBase, uint32(v19)+2)) = uint8(v20)
		v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v23 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+111)))
		*(*uint8)(unsafe.Add(mBase, uint32(v22)+3)) = uint8(v23)
		v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v26 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+143)))
		*(*uint8)(unsafe.Add(mBase, uint32(v25)+4)) = uint8(v26)
		v28 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+175)))
		*(*uint8)(unsafe.Add(mBase, uint32(v28)+5)) = uint8(v29)
		v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+207)))
		*(*uint8)(unsafe.Add(mBase, uint32(v31)+6)) = uint8(v32)
		v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+239)))
		*(*uint8)(unsafe.Add(mBase, uint32(v34)+7)) = uint8(v35)
		v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+271)))
		*(*uint8)(unsafe.Add(mBase, uint32(v37)+8)) = uint8(v38)
		v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v41 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+303)))
		*(*uint8)(unsafe.Add(mBase, uint32(v40)+9)) = uint8(v41)
		v43 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+335)))
		*(*uint8)(unsafe.Add(mBase, uint32(v43)+10)) = uint8(v44)
		v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v47 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+367)))
		*(*uint8)(unsafe.Add(mBase, uint32(v46)+11)) = uint8(v47)
		v49 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+399)))
		*(*uint8)(unsafe.Add(mBase, uint32(v49)+12)) = uint8(v50)
		v52 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+431)))
		*(*uint8)(unsafe.Add(mBase, uint32(v52)+13)) = uint8(v53)
		v55 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+463)))
		*(*uint8)(unsafe.Add(mBase, uint32(v55)+14)) = uint8(v56)
		v58 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+495)))
		*(*uint8)(unsafe.Add(mBase, uint32(v58)+15)) = uint8(v59)
		v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v62 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+23)))
		*(*uint8)(unsafe.Add(mBase, uint32(v61))) = uint8(v62)
		v64 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v65 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+31)))
		*(*uint8)(unsafe.Add(mBase, uint32(v64))) = uint8(v65)
		v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v68 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+55)))
		*(*uint8)(unsafe.Add(mBase, uint32(v67)+1)) = uint8(v68)
		v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v71 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+63)))
		*(*uint8)(unsafe.Add(mBase, uint32(v70)+1)) = uint8(v71)
		v73 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v74 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+87)))
		*(*uint8)(unsafe.Add(mBase, uint32(v73)+2)) = uint8(v74)
		v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v77 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+95)))
		*(*uint8)(unsafe.Add(mBase, uint32(v76)+2)) = uint8(v77)
		v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+119)))
		*(*uint8)(unsafe.Add(mBase, uint32(v79)+3)) = uint8(v80)
		v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+127)))
		*(*uint8)(unsafe.Add(mBase, uint32(v82)+3)) = uint8(v83)
		v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v86 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+151)))
		*(*uint8)(unsafe.Add(mBase, uint32(v85)+4)) = uint8(v86)
		v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v89 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+159)))
		*(*uint8)(unsafe.Add(mBase, uint32(v88)+4)) = uint8(v89)
		v91 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+183)))
		*(*uint8)(unsafe.Add(mBase, uint32(v91)+5)) = uint8(v92)
		v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v95 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+191)))
		*(*uint8)(unsafe.Add(mBase, uint32(v94)+5)) = uint8(v95)
		v97 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+215)))
		*(*uint8)(unsafe.Add(mBase, uint32(v97)+6)) = uint8(v98)
		v100 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+223)))
		*(*uint8)(unsafe.Add(mBase, uint32(v100)+6)) = uint8(v101)
		v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+247)))
		*(*uint8)(unsafe.Add(mBase, uint32(v103)+7)) = uint8(v104)
		v106 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5)+255)))
		*(*uint8)(unsafe.Add(mBase, uint32(v106)+7)) = uint8(v107)
		v109 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v110 = int32(-1)
		v112 = *(*int32)(unsafe.Add(mBase, uint32(l0)+320))
		v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v112)+15)))
		*(*uint8)(unsafe.Add(mBase, uint32(v109+v110))) = uint8(v113)
		v115 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v118 = *(*int32)(unsafe.Add(mBase, uint32(l0)+324))
		v119 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v118)+7)))
		*(*uint8)(unsafe.Add(mBase, uint32(v115+v110))) = uint8(v119)
		v121 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		v124 = *(*int32)(unsafe.Add(mBase, uint32(l0)+324))
		v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v124)+15)))
		*(*uint8)(unsafe.Add(mBase, uint32(v121+v110))) = uint8(v125)
	}
	v127 = *(*int32)(unsafe.Add(mBase, uint32(v8)+44))
	if v127+int32(-1) <= v6 {
	} else {
		v131 = *(*int32)(unsafe.Add(mBase, uint32(l0)+320))
		v133 = base.Simd_g_v128_load(m, v5, int32(480))
		v134 = int32(0)
		base.Simd_g_v128_store(m, v131, v134, v133)
		v136 = *(*int32)(unsafe.Add(mBase, uint32(l0)+324))
		v138 = base.Simd_g_v128_load(m, v5, int32(240))
		base.Simd_g_v128_store(m, v136, v134, v138)
	}
	return
}

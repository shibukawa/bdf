//go:build !bdf_noconv

package webpw

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	_ "unsafe"
)

//go:linkname F__initialize github.com/shibukawa/bdf/imgconv/internal/webpw/p0.F__initialize
func F__initialize(m *base.Module)

//go:linkname F_encode github.com/shibukawa/bdf/imgconv/internal/webpw/p0.F_encode
func F_encode(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32) int32

//go:linkname F_malloc github.com/shibukawa/bdf/imgconv/internal/webpw/p0.F_malloc
func F_malloc(m *base.Module, l0 int32) int32

//go:linkname F_free github.com/shibukawa/bdf/imgconv/internal/webpw/p0.F_free
func F_free(m *base.Module, l0 int32)

//go:linkname InitElemSeg_0_0 github.com/shibukawa/bdf/imgconv/internal/webpw/p0.InitElemSeg_0_0
func InitElemSeg_0_0(m *base.Module)

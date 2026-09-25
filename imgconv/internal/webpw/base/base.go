//go:build !bdf_noconv

package base

import (
	"math"
	"math/bits"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Module struct {
	Memory      []byte
	MaxMem      uint64
	M           unsafe.Pointer
	T0          []any
	G0          int32
	G1          int32
	G2          int32
	G3          int32
	G4          int32
	G5          int32
	G6          int32
	G7          int32
	G8          int32
	G9          int32
	G10         int32
	G11         int32
	G12         int32
	G13         int32
	G14         int32
	G15         int32
	G16         int32
	G17         int32
	G18         int32
	G19         int32
	G20         int32
	G21         int32
	G22         int32
	G23         int32
	G24         int32
	G25         int32
	G26         int32
	G27         int32
	G28         int32
	G29         int32
	G30         int32
	G31         int32
	G32         int32
	G33         int32
	G34         int32
	G35         int32
	G36         int32
	G37         int32
	G38         int32
	G39         int32
	G40         int32
	G41         int32
	G42         int32
	G43         int32
	G44         int32
	G45         int32
	G46         int32
	G47         int32
	G48         int32
	G49         int32
	G50         int32
	G51         int32
	G52         int32
	G53         int32
	G54         int32
	G55         int32
	G56         int32
	G57         int32
	G58         int32
	G59         int32
	G60         int32
	G61         int32
	G62         int32
	G63         int32
	G64         int32
	G65         int32
	G66         int32
	G67         int32
	G68         int32
	G69         int32
	G70         int32
	G71         int32
	G72         int32
	G73         int32
	G74         int32
	G75         int32
	G76         int32
	G77         int32
	G78         int32
	G79         int32
	G80         int32
	G81         int32
	G82         int32
	G83         int32
	G84         int32
	G85         int32
	G86         int32
	MemMu       *sync.Mutex
	MemSize     *atomic.Uint64
	DataEnd     uint32
	MemShared   bool
	Threads     *ThreadPool
	ThreadStart func(*Module, int32, int32)
}

func I32(x int32) int32 { return x }

func I64(x int64) int64 { return x }

// ui32 / ui64 reinterpret a signed integer as its unsigned bit
// equivalent at runtime. Used for the operands of wasm unsigned
// comparisons (i32.lt_u etc.) — emitting `uint32(int32(-N))` directly
// fails Go's compile-time constant rule because the negative typed
// constant isn't representable in uint32; routing through these
// function-call boundaries forces runtime conversion.
func Ui32(x int32) uint32 { return uint32(x) }

func Ui64(x int64) uint64 { return uint64(x) }

// b2i32 materialises a wasm comparison result — an i32 that is 0 or 1 — from
// the Go bool the comparison expression evaluates to.
//
// It exists as a named helper rather than an inline `func() int32 { ... }()`
// because the gcasm backend requires every direct call left in the compiled
// output to be either a package-local FnN or something the Go inliner removed.
// A func literal is normally inlined at its call site, but the inliner gives up
// once the ENCLOSING function grows past its budget — and a single wasm function
// can translate to tens of thousands of lines of Go, as an interpreter's
// bytecode dispatch loop does. The literal is then outlined into a real closure
// symbol (FnN.funcA.funcB), which reaches the assembler as a direct call gcasm
// cannot marshal. A named helper this small is always inlined, and if it ever
// were not, it would fail loudly at its own symbol rather than as a nested
// closure.
func B2i32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func F32(x float32) float32 { runtime.KeepAlive(&x); return x }

func F64(x float64) float64 { runtime.KeepAlive(&x); return x }

//go:noinline
func Wasm_trap_div_zero() { panic("wasm: integer divide by zero") }

//go:noinline
func Wasm_trap_int_overflow() { panic("wasm: integer overflow") }

//go:noinline
func Wasm_trap_invalid_conv() { panic("wasm: invalid conversion to integer") }

//go:noinline
func Wasm_trap_unreachable() { panic("wasm: unreachable") }

//go:noinline
func Wasm_trap_memfill_oob() { panic("wasm: memory.fill out of bounds") }

//go:noinline
func Wasm_trap_memcopy_oob() { panic("wasm: memory.copy out of bounds") }

func I32_div_s(x, y int32) int32 {
	if y == -1 && x == math.MinInt32 {
		Wasm_trap_int_overflow()
	}
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I64_div_s(x, y int64) int64 {
	if y == -1 && x == math.MinInt64 {
		Wasm_trap_int_overflow()
	}
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I32_div_u(x, y uint32) uint32 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I64_div_u(x, y uint64) uint64 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I32_rem_u(x, y uint32) uint32 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x % y
}

func I64_rem_u(x, y uint64) uint64 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x % y
}

func I32_rotl(x, y int32) int32 { return int32(bits.RotateLeft32(uint32(x), int(y&31))) }

func I64_rotl(x, y int64) int64 { return int64(bits.RotateLeft64(uint64(x), int(y&63))) }

func F32_abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}

func F64_abs(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) &^ (1 << 63))
}

func F32_neg(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) ^ (1 << 31))
}

func F64_neg(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) ^ (1 << 63))
}

func I32_trunc_f32_s(x float32) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}

	if !(x > -2147483904.0 && x < 2147483648.0) {
		Wasm_trap_int_overflow()
	}
	return int32(x)
}

func I32_trunc_f32_u(x float32) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 4294967296.0) {
		Wasm_trap_int_overflow()
	}
	return int32(uint32(x))
}

func I32_trunc_f64_s(x float64) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}

	if !(x > -2147483649.0 && x < 2147483648.0) {
		Wasm_trap_int_overflow()
	}
	return int32(x)
}

func I32_trunc_f64_u(x float64) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 4294967296.0) {
		Wasm_trap_int_overflow()
	}
	return int32(uint32(x))
}

func I64_trunc_f64_u(x float64) int64 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 18446744073709551616.0) {
		Wasm_trap_int_overflow()
	}
	return int64(uint64(x))
}

// memorySize returns the current size of m.memory in wasm pages (each
// page is 64 KiB).
func MemorySize(m *Module) int32 {
	return int32(m.MemSize.Load() >> 16)
}

// wasmMemHardCap is the implementation limit on linear-memory size:
// 65534 pages, two short of wasm32's architectural 65536. Growth past
// it fails with -1 like any resource limit (the JS API allows an
// engine to refuse any grow). Keeping memSize strictly below 2^32
// minus a 128 KiB margin is what makes the coalesced SIMD bounds check
// (simd_v128_load_rng) exact: a group whose unwrapped address range
// reaches past memSize can then never be a group whose members all
// individually landed in bounds via u32 wraparound.
//
// A function rather than a const because the helper extractor carries
// only function declarations into the output (it must stay in sync
// with codegen's wasmMemHardCapBytes).
func WasmMemHardCap() uint64 { return (1 << 32) - (1 << 17) }

// memoryGrow grows m.memory by n wasm pages (64 KiB each). Returns the
// previous page count, or -1 if the new size would exceed maxMem or
// wasmMemHardCap. n may be 0, which simply returns the current size.
//
// len(m.memory) must always equal the exact wasm memory size (memory.size
// and every bounds check depend on it), but the backing array is grown
// GEOMETRICALLY: a sequence of small memory.grow calls — which a C++ heap
// does constantly during start-up — would otherwise reallocate and recopy
// the whole linear memory on every page, i.e. O(n^2) total copying. Spare
// capacity makes the common grow a zero-copy reslice and amortizes the
// reallocations to O(n).
func MemoryGrow(m *Module, n int32) int32 {

	m.MemMu.Lock()
	defer m.MemMu.Unlock()
	cur := m.MemSize.Load()
	prev := int32(cur >> 16)
	if n == 0 {
		return prev
	}
	if n < 0 {
		return -1
	}
	want := cur + uint64(n)*65536
	if m.MaxMem != 0 && want > m.MaxMem {
		return -1
	}
	if want > WasmMemHardCap() {
		return -1
	}
	if m.MemShared {

		if want > uint64(len(m.Memory)) {
			return -1
		}
		m.MemSize.Store(want)
		return prev
	}
	if want <= uint64(cap(m.Memory)) {

		m.Memory = m.Memory[:want]
		m.MemSize.Store(want)
		return prev
	}

	newCap := uint64(cap(m.Memory)) * 2
	if newCap < want {
		newCap = want
	}
	if m.MaxMem != 0 && newCap > m.MaxMem {
		newCap = m.MaxMem
	}
	if newCap > WasmMemHardCap() {
		newCap = WasmMemHardCap()
	}
	grown := make([]byte, want, newCap)
	copy(grown, m.Memory)
	m.Memory = grown
	m.MemSize.Store(want)

	m.M = unsafe.Pointer(unsafe.SliceData(m.Memory))
	return prev
}

// accessMemory runs f with the module's current linear memory while
// holding the same lock memoryGrow takes to mutate the memory slice
// header or relocate its backing array. It is the ONE safe way to
// touch linear memory from OUTSIDE the module's execution goroutine —
// e.g. a watchdog goroutine raising CPython's eval-breaker bit while
// an evaluation is running. For the duration of f the memory can
// neither be resliced nor relocated, so f's writes land in the array
// the guest observes; a grow that raced in just before blocks until f
// returns and then copies f's writes forward with the rest of the
// contents. Determinism notes for callers:
//
//   - f MUST NOT call back into the module or into memoryGrow — that
//     would self-deadlock.
//   - f should be short: a running guest blocks inside memory.grow
//     until f returns (ordinary guest loads/stores do not block).
//   - Bytes the guest reads or writes concurrently with f (that is
//     the point of an eval-breaker-style flag) are exchanged with
//     plain single-word accesses; keep such shared words
//     word-aligned and word-sized.
func AccessMemory(m *Module, f func(mem []byte)) {
	m.MemMu.Lock()
	defer m.MemMu.Unlock()
	f(m.Memory)
}

func I32_div_u_s(x, y int32) int32 { return int32(I32_div_u(uint32(x), uint32(y))) }
func I32_rem_u_s(x, y int32) int32 { return int32(I32_rem_u(uint32(x), uint32(y))) }
func I64_div_u_s(x, y int64) int64 { return int64(I64_div_u(uint64(x), uint64(y))) }
func I64_rem_u_s(x, y int64) int64 { return int64(I64_rem_u(uint64(x), uint64(y))) }

// The explicit same-type conversions are NOT redundant: they are
// rounding points. Once these helpers inline, gc is free to fuse a
// multiply feeding an add into a single FMA — legal Go, but wasm
// requires every operation individually rounded, and a fused result
// diverges from every wasm runtime (bitwise, and observably in greedy
// sampling). A float conversion forces the intermediate rounding and
// forbids the fusion (spec: Conversions, "rounds to the precision of
// the target type"; the same rule math.FMA documents).
func F32_add(x, y float32) float32 { return float32(x + y) }
func F32_sub(x, y float32) float32 { return float32(x - y) }
func F32_mul(x, y float32) float32 { return float32(x * y) }
func F32_div(x, y float32) float32 { return float32(x / y) }
func F64_add(x, y float64) float64 { return float64(x + y) }
func F64_sub(x, y float64) float64 { return float64(x - y) }
func F64_mul(x, y float64) float64 { return float64(x * y) }
func F64_div(x, y float64) float64 { return float64(x / y) }

func I32_clz(x int32) int32 { return int32(bits.LeadingZeros32(uint32(x))) }
func I32_ctz(x int32) int32 { return int32(bits.TrailingZeros32(uint32(x))) }

func F32_ceil(x float32) float32 { return float32(math.Ceil(float64(x))) }

func F32_floor(x float32) float32 { return float32(math.Floor(float64(x))) }

func F32_sqrt(x float32) float32 { return float32(math.Sqrt(float64(x))) }

func F32_eq(x, y float32) int32 {
	if x == y {
		return 1
	}
	return 0
}
func F32_ne(x, y float32) int32 {
	if x != y {
		return 1
	}
	return 0
}
func F32_lt(x, y float32) int32 {
	if x < y {
		return 1
	}
	return 0
}
func F32_gt(x, y float32) int32 {
	if x > y {
		return 1
	}
	return 0
}
func F32_le(x, y float32) int32 {
	if x <= y {
		return 1
	}
	return 0
}
func F32_ge(x, y float32) int32 {
	if x >= y {
		return 1
	}
	return 0
}

func F64_eq(x, y float64) int32 {
	if x == y {
		return 1
	}
	return 0
}
func F64_ne(x, y float64) int32 {
	if x != y {
		return 1
	}
	return 0
}
func F64_lt(x, y float64) int32 {
	if x < y {
		return 1
	}
	return 0
}
func F64_gt(x, y float64) int32 {
	if x > y {
		return 1
	}
	return 0
}
func F64_le(x, y float64) int32 {
	if x <= y {
		return 1
	}
	return 0
}
func F64_ge(x, y float64) int32 {
	if x >= y {
		return 1
	}
	return 0
}

func I32_wrap_i64(x int64) int32       { return int32(x) }
func I64_extend_i32_s(x int32) int64   { return int64(x) }
func I64_extend_i32_u(x int32) int64   { return int64(uint32(x)) }
func F32_demote_f64(x float64) float32 { return float32(x) }
func F64_promote_f32(x float32) float64 {

	if math.IsNaN(float64(x)) {

		return float64(x)
	}
	return float64(x)
}

func F32_convert_i32_s(x int32) float32 { return float32(x) }
func F32_convert_i32_u(x int32) float32 { return float32(uint32(x)) }

func F64_convert_i32_s(x int32) float64 { return float64(x) }
func F64_convert_i32_u(x int32) float64 { return float64(uint32(x)) }

func F64_convert_i64_u(x int64) float64 { return float64(uint64(x)) }

func I32_reinterpret_f32(x float32) int32 { return int32(math.Float32bits(x)) }
func I64_reinterpret_f64(x float64) int64 { return int64(math.Float64bits(x)) }
func F32_reinterpret_i32(x int32) float32 { return math.Float32frombits(uint32(x)) }
func F64_reinterpret_i64(x int64) float64 { return math.Float64frombits(uint64(x)) }

func I32_extend8_s(x int32) int32  { return int32(int8(x)) }
func I32_extend16_s(x int32) int32 { return int32(int16(x)) }

func I64_extend32_s(x int64) int64 { return int64(int32(x)) }

func MemoryFill(m *Module, dst int32, val int32, n int32) {
	if n == 0 {
		return
	}
	end := uint64(uint32(dst)) + uint64(uint32(n))
	if end > MemBound(m) {
		Wasm_trap_memfill_oob()
	}
	b := m.Memory[uint32(dst):uint32(end)]
	v := byte(val)

	if v == 0 {
		for k := range b {
			b[k] = 0
		}
		return
	}
	b[0] = v
	for filled := 1; filled < len(b); filled *= 2 {
		copy(b[filled:], b[:filled])
	}
}

func MemoryCopy(m *Module, dst int32, src int32, n int32) {
	if n == 0 {
		return
	}
	srcEnd := uint64(uint32(src)) + uint64(uint32(n))
	dstEnd := uint64(uint32(dst)) + uint64(uint32(n))
	if size := MemBound(m); srcEnd > size || dstEnd > size {
		Wasm_trap_memcopy_oob()
	}
	copy(m.Memory[uint32(dst):uint32(dstEnd)], m.Memory[uint32(src):uint32(srcEnd)])
}

// memBound is the highest address a bounds-checked bulk or atomic access
// may reach. A shared memory's slice spans the whole declared maximum from
// the start, so only memSize says how much of it the guest may touch (and
// reading it atomically keeps growth race-free without a lock). Otherwise
// the slice is the truth: an embedder that maps the whole growable range
// up front and aliases shared segments above the guest-visible size into
// it keeps every byte of the slice addressable, exactly as the unchecked
// load/store paths do.
func MemBound(m *Module) uint64 {
	if m.MemShared {
		return m.MemSize.Load()
	}
	return uint64(len(m.Memory))
}

// forceContendedAtomics makes every atomic helper of m take its LOCKed
// path from now on, as if a wasi thread had been spawned. For an embedder
// that shares part of m's linear memory with OTHER instances (several
// single-threaded modules aliasing one segment, the way processes share
// System V memory): no thread of m ever exists, yet the atomics in the
// shared range race with those instances' goroutines. Host-facing API;
// nothing in the generated code calls it.
func ForceContendedAtomics(m *Module) {
	if m.Threads == nil {
		m.Threads = &ThreadPool{}
	}
	if m.Threads.nextTID.Load() == 0 {
		m.Threads.nextTID.Store(1)
	}
}

var spinAgents int32
var spinOversubscribed uint32

type ThreadPool struct {
	nextTID atomic.Int32
	wg      sync.WaitGroup

	parkMu sync.Mutex
	parked map[uint64][]chan struct{}
}

// wake releases up to count waiters on ea and reports how many it woke.
func (p *ThreadPool) wake(ea uint64, count int32) int32 {
	p.parkMu.Lock()
	defer p.parkMu.Unlock()
	waiters := p.parked[ea]
	n := int32(len(waiters))
	if count >= 0 && count < n {
		n = count
	}
	for _, ch := range waiters[:n] {
		close(ch)
	}
	if int(n) == len(waiters) {
		delete(p.parked, ea)
	} else {
		p.parked[ea] = waiters[n:]
	}
	return n
}

// SaveGlobals returns the module's mutable globals, in a form that can be handed back
// to RestoreGlobals. It is how a snapshot of an instance captures the state that does not
// live in linear memory.
func SaveGlobals(m *Module) []uint64 {
	g := make([]uint64, 1)
	g[0] = uint64(uint32(m.G0))
	return g
}

// RestoreGlobals puts a snapshot's globals back. A snapshot from a different module (or a
// different build of the same one) has a different global count; rather than
// index out of bounds, take what fits and leave the rest at their declared
// initializers.
func RestoreGlobals(m *Module, g []uint64) {
	if len(g) != 1 {
		return
	}
	m.G0 = int32(uint32(g[0]))
}

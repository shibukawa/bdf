#!/bin/sh
# Rebuilds imgconv/internal/webpw: compiles libwebp to wasm (encode only, with
# its name section) using wasi-sdk, then translates it to Go with the wasm2go
# fork, producing per-subject files with symbol names and named data addresses
# so that a rebuild changes as little of the checked-in tree as possible.
#
# SIMD=1 builds the SIMD variant instead: libwebp's SSE2/SSE4.1 kernels
# compiled to wasm SIMD through the emscripten compat headers in
# tools/webp/emcompat, and every v128 function emitted twice, over [2]uint64
# pairs and, for a GOEXPERIMENT=simd build on Go 1.27, over simd/archsimd
# vector registers (wasm2go -simd=go127). It is not what is checked in: see
# docs/design.md 3.2 for the measurements.
#
# Usage: tools/gen-codecs.sh            (WORK=dir keeps the toolchain between runs;
#                                        OUT=dir writes the package elsewhere)
set -eu
ROOT=$(cd "$(dirname "$0")/.." && pwd)
WORK=${WORK:-"$ROOT/../bdf-codec-work"}
OUT=${OUT:-"$ROOT/imgconv/internal/webpw"}
SIMD=${SIMD:-0}
mkdir -p "$WORK"

FORK_REPO=https://github.com/shibukawa/wasm2go-fork
FORK_REF=636a813ef228f05627198ae8b0bebc12788439ce   # pgmem: -simd=go127
WASI_SDK_VER=25
LIBWEBP_VER=v1.6.0

# 1. wasm2go fork (consumed through a replace directive by tools/gen-webp).
if [ ! -d "$WORK/wasm2go-fork" ]; then
  git clone -q "$FORK_REPO" "$WORK/wasm2go-fork"
fi
git -C "$WORK/wasm2go-fork" checkout -q "$FORK_REF"

# 2. wasi-sdk.
SDK="$WORK/wasi-sdk-$WASI_SDK_VER.0-x86_64-linux"
if [ ! -x "$SDK/bin/clang" ]; then
  curl -sSL -o "$WORK/wasi-sdk.tar.gz" "https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-$WASI_SDK_VER/wasi-sdk-$WASI_SDK_VER.0-x86_64-linux.tar.gz"
  tar xzf "$WORK/wasi-sdk.tar.gz" -C "$WORK"
fi

# 3. libwebp static libraries (no threads, encoder + sharpyuv only). With
#    SIMD=1 the SSE kernels are compiled through wasm SIMD: libwebp selects
#    them from the __SSE*__ macros, and the emscripten compat headers map the
#    intrinsics onto wasm_simd128.h (wasi-sdk has no such headers of its own).
if [ ! -d "$WORK/libwebp" ]; then
  git clone -q -b "$LIBWEBP_VER" --depth 1 https://github.com/webmproject/libwebp "$WORK/libwebp"
fi
if [ "$SIMD" = 1 ]; then
  BUILD="$WORK/libwebp/build-sse"
  SIMD_CFLAGS="-msimd128 -D__SSE__=1 -D__SSE2__=1 -D__SSE3__=1 -D__SSSE3__=1 -D__SSE4_1__=1 -DEMSCRIPTEN=1 -I$ROOT/tools/webp/emcompat"
  SIMD_TARGET=go127
else
  BUILD="$WORK/libwebp/build"
  SIMD_CFLAGS=
  SIMD_TARGET=
fi
mkdir -p "$BUILD"
( cd "$BUILD" && cmake .. -G Ninja -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=0 \
    -DCMAKE_C_FLAGS="$SIMD_CFLAGS" -DWEBP_ENABLE_SIMD="$SIMD" -DWEBP_ENABLE_SIMD_DEFAULT="$SIMD" -DWEBP_BUILD_EXTRAS=0 -DWEBP_USE_THREAD=0 -DWEBP_BUILD_ANIM_UTILS=0 \
    -DWEBP_BUILD_CWEBP=0 -DWEBP_BUILD_DWEBP=0 -DWEBP_BUILD_IMG2WEBP=0 -DWEBP_BUILD_WEBPINFO=0 \
    -DWEBP_BUILD_WEBPMUX=0 -DWEBP_BUILD_LIBWEBPMUX=0 -DWEBP_BUILD_GIF2WEBP=0 -DWEBP_BUILD_VWEBP=0 \
    -DCMAKE_TOOLCHAIN_FILE="$SDK/share/cmake/wasi-sdk.cmake" >/dev/null && ninja libwebp.a libsharpyuv.a >/dev/null )

# 4. The encode-only module. wasm-ld keeps the name section, which the
#    translator needs for -symbol-names / -group-files.
"$SDK/bin/clang" --sysroot="$SDK/share/wasi-sysroot" -O3 -Wl,--no-entry \
  -Wl,--export=malloc -Wl,--export=free -Wl,--export=encode -mexec-model=reactor -mnontrapping-fptoint \
  -I"$WORK/libwebp/src" -I"$BUILD/src" -Wall -o "$WORK/webp.wasm" \
  "$ROOT/tools/webp/webp.c" "$BUILD/libwebp.a" "$BUILD/libsharpyuv.a"

# 5. Translate to Go.
MODFILE="$WORK/gen-webp.go.mod"
cp "$ROOT/tools/gen-webp/go.mod" "$MODFILE"
( cd "$ROOT/tools/gen-webp" && go mod edit -modfile="$MODFILE" -replace "github.com/goccy/wasm2go=$WORK/wasm2go-fork" \
  && go mod tidy -modfile="$MODFILE" >/dev/null 2>&1 \
  && go run -modfile="$MODFILE" . -i "$WORK/webp.wasm" -out-dir "$OUT" -simd="$SIMD_TARGET" )
cp "$WORK/libwebp/COPYING" "$OUT/LICENSE.libwebp"
python3 "$ROOT/tools/add-build-tag.py" "$OUT" '!bdf_noconv'
gofmt -w "$OUT"
echo "done: $(find "$OUT" -type f | wc -l) files"

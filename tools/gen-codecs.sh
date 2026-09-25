#!/bin/sh
# Rebuilds imgconv/internal/webpw: compiles libwebp to wasm (encode only, with
# its name section) using wasi-sdk, then translates it to Go with the wasm2go
# fork, producing per-subject files with symbol names and named data addresses
# so that a rebuild changes as little of the checked-in tree as possible.
#
# Usage: tools/gen-codecs.sh            (WORK=dir keeps the toolchain between runs)
set -eu
ROOT=$(cd "$(dirname "$0")/.." && pwd)
WORK=${WORK:-"$ROOT/../bdf-codec-work"}
mkdir -p "$WORK"

FORK_REPO=https://github.com/shibukawa/wasm2go-fork
FORK_REF=ac98bcf00c17d8531f0c071a9836d0b50975e7ff   # v0.5.15-fork.7
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

# 3. libwebp static libraries (no SIMD, no threads, encoder + sharpyuv only).
if [ ! -d "$WORK/libwebp" ]; then
  git clone -q -b "$LIBWEBP_VER" --depth 1 https://github.com/webmproject/libwebp "$WORK/libwebp"
fi
mkdir -p "$WORK/libwebp/build"
( cd "$WORK/libwebp/build" && cmake .. -G Ninja -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=0 \
    -DWEBP_ENABLE_SIMD_DEFAULT=0 -DWEBP_BUILD_EXTRAS=0 -DWEBP_USE_THREAD=0 -DWEBP_BUILD_ANIM_UTILS=0 \
    -DWEBP_BUILD_CWEBP=0 -DWEBP_BUILD_DWEBP=0 -DWEBP_BUILD_IMG2WEBP=0 -DWEBP_BUILD_WEBPINFO=0 \
    -DWEBP_BUILD_WEBPMUX=0 -DWEBP_BUILD_LIBWEBPMUX=0 -DWEBP_BUILD_GIF2WEBP=0 -DWEBP_BUILD_VWEBP=0 \
    -DCMAKE_TOOLCHAIN_FILE="$SDK/share/cmake/wasi-sdk.cmake" >/dev/null && ninja libwebp.a libsharpyuv.a >/dev/null )

# 4. The encode-only module. wasm-ld keeps the name section, which the
#    translator needs for -symbol-names / -group-files.
"$SDK/bin/clang" --sysroot="$SDK/share/wasi-sysroot" -O3 -Wl,--no-entry \
  -Wl,--export=malloc -Wl,--export=free -Wl,--export=encode -mexec-model=reactor -mnontrapping-fptoint \
  -I"$WORK/libwebp/src" -I"$WORK/libwebp/build/src" -Wall -o "$WORK/webp.wasm" \
  "$ROOT/tools/webp/webp.c" "$WORK/libwebp/build/libwebp.a" "$WORK/libwebp/build/libsharpyuv.a"

# 5. Translate to Go.
MODFILE="$WORK/gen-webp.go.mod"
cp "$ROOT/tools/gen-webp/go.mod" "$MODFILE"
( cd "$ROOT/tools/gen-webp" && go mod edit -modfile="$MODFILE" -replace "github.com/goccy/wasm2go=$WORK/wasm2go-fork" \
  && go mod tidy -modfile="$MODFILE" >/dev/null 2>&1 \
  && go run -modfile="$MODFILE" . -i "$WORK/webp.wasm" -out-dir "$ROOT/imgconv/internal/webpw" )
cp "$WORK/libwebp/COPYING" "$ROOT/imgconv/internal/webpw/LICENSE.libwebp"
python3 "$ROOT/tools/add-build-tag.py" "$ROOT/imgconv/internal/webpw" '!bdf_noconv'
gofmt -w "$ROOT/imgconv/internal/webpw"
echo "done: $(find "$ROOT/imgconv/internal/webpw" -type f | wc -l) files"

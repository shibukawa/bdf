#!/bin/bash
# Makes the TIFF test files: converter/internal/tiff/testdata (the same
# pictures stored in every way the decoder reads, with PNG references that
# ImageMagick decoded through libtiff) and converter/tiff/testdata (files
# as scanners, fax software and image editors write them).
#
# Needs ImageMagick 7 and the libtiff tools (tiffcp, tiffset); made with
# ImageMagick 7.1.2 and libtiff 4.7.2. Text is drawn with the M PLUS 1p test
# font, so it does not depend on the fonts of the machine.
set -euo pipefail
root=$(cd "$(dirname "$0")/../.." && pwd)
font="$root/converter/pptx/testdata/fonts/MPLUS1p-Regular-subset.ttf"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
cd "$tmp"

# --- converter/internal/tiff/testdata ---------------------------------------
out="$root/converter/internal/tiff/testdata"
mkdir -p "$out"
rm -f "$out"/*.tif "$out"/*.png
im() { magick "$@"; }

# The pictures: colour, grey, bilevel text, a 16-colour palette, RGBA, CMYK.
im -size 67x45 -seed 5 plasma:red-blue -depth 8 -type TrueColor -density 72 -units PixelsPerInch rgb.tif
im rgb.tif -colorspace Gray -depth 8 -type Grayscale gray.tif
im -size 301x203 xc:white -font "$font" -pointsize 14 -fill black \
  -annotate +5+20 $'TIFF bilevel 1234\nabcdefg HIJKLMN\nThe quick brown fox' \
  -stroke black -draw "line 0,202 300,100" -monochrome -type Bilevel -depth 1 -density 72 -units PixelsPerInch bw.tif
im bw.tif -define quantum:polarity=min-is-white -type Bilevel bw-white.tif
im rgb.tif -colors 16 -type Palette pal.tif
im rgb.tif \( -size 45x67 gradient:white-black -rotate 90 \) -alpha off -compose CopyOpacity -composite -depth 8 rgba.tif
im rgb.tif -colorspace CMYK -depth 8 cmyk.tif
for p in rgb gray bw pal rgba cmyk; do im "$p.tif" -colorspace sRGB -depth 8 -alpha on "PNG32:$out/ref-$p.png"; done

# Every page of these files is the same picture, stored another way.
tiffcp -c none rgb.tif "$out/rgb.tif"
for c in lzw lzw:2 zip zip:2 packbits; do tiffcp -a -c "$c" rgb.tif "$out/rgb.tif"; done
tiffcp -a -c lzw -p separate rgb.tif "$out/rgb.tif"
tiffcp -a -c zip -t -w 16 -l 16 rgb.tif "$out/rgb.tif"
tiffcp -a -c lzw -r 7 rgb.tif "$out/rgb.tif"
im rgb.tif -depth 16 rgb16.tif
tiffcp -a -c lzw:2 rgb16.tif "$out/rgb.tif"
tiffcp -B -c lzw:2 rgb16.tif "$out/rgb-be.tif"
tiffcp -a -c none rgb.tif "$out/rgb-be.tif"
tiffcp -8 -c lzw rgb.tif "$out/rgb-bigtiff.tif"
tiffcp -a -c zip -p separate rgb.tif "$out/rgb-bigtiff.tif"

tiffcp -c none gray.tif "$out/gray.tif"
for c in lzw:2 zip packbits; do tiffcp -a -c "$c" gray.tif "$out/gray.tif"; done
im gray.tif -depth 16 gray16.tif
tiffcp -a -c zip:2 gray16.tif "$out/gray.tif"

tiffcp -c none bw.tif "$out/bw.tif"
for c in packbits lzw zip g3 g3:fill g3:2d g3:2d:fill g4; do tiffcp -a -c "$c" bw.tif "$out/bw.tif"; done
tiffcp -a -c g4 -f lsb2msb bw.tif "$out/bw.tif"
tiffcp -a -c g4 -r 16 bw.tif "$out/bw.tif"
tiffcp -a -c g4 -t -w 64 -l 64 bw.tif "$out/bw.tif"
for c in none g3:2d:fill g4; do tiffcp -a -c "$c" bw-white.tif "$out/bw.tif"; done
# One BlackIsZero G4 page, which golang.org/x/image/tiff reads inverted.
tiffcp -c g4 bw.tif "$out/g4.tif"

tiffcp -c none pal.tif "$out/pal.tif"
tiffcp -a -c packbits pal.tif "$out/pal.tif"
tiffcp -c lzw rgba.tif "$out/rgba.tif"
tiffcp -c none cmyk.tif "$out/cmyk.tif"
tiffcp -a -c zip cmyk.tif "$out/cmyk.tif"
tiffcp -a -c jpeg:95 -r 16 cmyk.tif "$out/cmyk.tif"

# JPEG pages: YCbCr in strips of 16 and 32 rows, RGB, grey, tiles, and the
# RGB strips of ImageMagick (the only page of jpeg-magick.tif).
tiffcp -c jpeg:95 -r 16 rgb.tif "$out/jpeg.tif"
tiffcp -a -c jpeg:95 -r 32 rgb.tif "$out/jpeg.tif"
tiffcp -a -c jpeg:95:r -r 8 rgb.tif "$out/jpeg.tif"
tiffcp -a -c jpeg:95 -r 16 gray.tif "$out/jpeg.tif"
tiffcp -a -c jpeg:95 -t -w 16 -l 16 rgb.tif "$out/jpeg.tif"
im rgb.tif -compress JPEG -quality 95 -define tiff:rows-per-strip=8 "$out/jpeg-magick.tif"

# Two pages with a reduced-resolution copy of the first between them.
tiffcp -c lzw rgb.tif "$out/thumb.tif"
im rgb.tif -resize 50% rgb-small.tif
tiffcp -a -c lzw rgb-small.tif "$out/thumb.tif"
tiffcp -a -c lzw gray.tif "$out/thumb.tif"
tiffset -d 1 -s 254 1 "$out/thumb.tif"

# --- converter/tiff/testdata -------------------------------------------------
out="$root/converter/tiff/testdata"
mkdir -p "$out"
rm -f "$out"/*.tif

# A fax: two pages of 1728 pixels at 204 × 98 dpi, two-dimensional Group 3
# coding with fill bits (what fax software writes).
for n in 1 2; do
  im -size 1728x280 xc:white -font "$font" -pointsize 40 -fill black \
    -annotate +100+120 "Fax page $n of 2" -annotate +100+220 "0123456789 ABCDEFG" \
    -monochrome -type Bilevel -depth 1 -define quantum:polarity=min-is-white \
    -density 204x98 -units PixelsPerInch "fax$n.tif"
done
tiffcp -c g3:2d:fill fax1.tif fax2.tif "$out/fax.tif"

# A scan: a 300 dpi bilevel page (scaled down to the resolution cap), a
# 150 dpi JPEG page (kept as it is), and a page without a resolution, with
# the descriptive tags on the first page.
im -size 1000x1300 xc:white -font "$font" -pointsize 12 -density 300 -fill black \
  -annotate +80+150 $'Scanned page one\nThe quick brown fox jumps over the lazy dog.\n0123456789' \
  -stroke black -strokewidth 3 -draw "rectangle 80,400 920,1200" -draw "line 80,400 920,1200" \
  -monochrome -type Bilevel -depth 1 -units PixelsPerInch scan1.tif
im -size 300x390 -seed 9 plasma:gold-teal -blur 0x4 -font "$font" -pointsize 12 -density 150 -fill black \
  -annotate +30+40 'Scanned page two' -depth 8 -type TrueColor -units PixelsPerInch scan2.tif
im -size 200x150 -seed 3 plasma:white-navy -blur 0x2 -depth 8 -type TrueColor -units Undefined scan3.tif
tiffcp -c g4 scan1.tif "$out/scan.tif"
tiffcp -a -c jpeg:90 -r 16 scan2.tif "$out/scan.tif"
tiffcp -a -c lzw:2 scan3.tif "$out/scan.tif"
for tag in 282 283 296; do tiffset -d 2 -u "$tag" "$out/scan.tif"; done
tiffset -s 269 "Scan test" "$out/scan.tif"
tiffset -s 270 "Three pages of a scanner" "$out/scan.tif"
tiffset -s 315 "Alice; Bob" "$out/scan.tif"
tiffset -s 33432 "(c) 2026 Example" "$out/scan.tif"
tiffset -s 306 "2026:09:26 18:30:00" "$out/scan.tif"

# The eight orientations of one picture, each stored so that it shows the
# same way up.
im -size 96x64 xc:white -fill '#d33' -draw "rectangle 0,0 24,16" -fill black -font "$font" \
  -pointsize 20 -annotate +30+30 'UP' -pointsize 14 -annotate +8+58 'bottom' \
  -depth 8 -type TrueColor -density 72 -units PixelsPerInch up.tif
ops=("" "" "-flop" "-rotate 180" "-flip" "-transpose" "-rotate -90" "-transverse" "-rotate 90")
for k in 1 2 3 4 5 6 7 8; do
  # shellcheck disable=SC2086
  im up.tif ${ops[$k]} +repage -density 72 -units PixelsPerInch "o$k.tif"
  tiffset -s 274 "$k" "o$k.tif"
done
tiffcp -c lzw o1.tif o2.tif o3.tif o4.tif o5.tif o6.tif o7.tif o8.tif "$out/orientation.tif"

ls -l "$root/converter/internal/tiff/testdata" "$root/converter/tiff/testdata"

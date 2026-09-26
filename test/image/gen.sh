#!/bin/sh
# Regenerates the test images of the image converter
# (converter/image/testdata): small pictures in every format, with the
# metadata the converter reads. Requires ImageMagick 7 (magick), exiftool,
# cwebp and avifenc.
set -eu
cd "$(dirname "$0")/../.."
out=converter/image/testdata
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
x() { exiftool -q -overwrite_original "$@"; }

# A sky, a field and a sun in the top left corner: a wrong orientation shows.
magick -size 64x48 gradient:'#5fa8ff'-'#e3f2ff' \
  -fill '#3a9d3a' -draw 'rectangle 0,32 63,47' \
  -fill '#e53935' -draw 'circle 14,12 14,5' "$tmp/scene.png"

# JPEG: turned by EXIF (orientation 6: the browser turns it 90° clockwise),
# with EXIF, XMP (a title in two languages) and IPTC that disagree, to test
# which wins; the camera's placeholder description is ignored.
magick "$tmp/scene.png" -quality 92 -strip "$out/photo.jpg"
x -Orientation#=6 \
  -EXIF:Artist='Taro Yamada; Hanako Sato' -EXIF:Copyright='© 2026 Taro Yamada' \
  -EXIF:ImageDescription='OLYMPUS DIGITAL CAMERA' -EXIF:XPTitle='Windows title' -EXIF:XPKeywords='sun;field' \
  -EXIF:DateTimeOriginal='2026:09:01 10:20:30' -EXIF:OffsetTimeOriginal='+09:00' -EXIF:ModifyDate='2026:09:02 08:00:00' \
  -XMP-dc:Title='Sunny field' -XMP-dc:Title-ja='晴れた野原' -XMP-dc:Description='A red sun over a green field' \
  -XMP-dc:Subject=sun -XMP-dc:Subject=field -XMP-dc:Subject='晴天' -XMP-photoshop:DateCreated='2026-09-01T10:20:30+09:00' \
  -IPTC:CodedCharacterSet=UTF8 -IPTC:ObjectName='IPTC name' -IPTC:By-line='IPTC photographer' -IPTC:Keywords='iptc' \
  "$out/photo.jpg"

# JPEG with IPTC only (UTF-8).
magick -size 8x8 xc:'#808080' -quality 90 -strip "$out/iptc.jpg"
x -IPTC:CodedCharacterSet=UTF8 -IPTC:ObjectName='夕焼け' -IPTC:By-line='Ann' -IPTC:By-line='Bob' \
  -IPTC:Keywords='sky' -IPTC:Keywords='evening' -IPTC:Caption-Abstract='The sky at dusk' \
  -IPTC:CopyrightNotice='Ann and Bob' -IPTC:DateCreated=2026:08:31 -IPTC:TimeCreated=18:45:00+09:00 "$out/iptc.jpg"

# PNG with transparency and its text chunks (a Japanese title goes in iTXt).
magick -size 32x24 xc:none -fill '#1e88e5' -draw 'roundrectangle 2,2 29,21 5,5' -strip "$out/tags.png"
x -PNG:Title='PNG の題名' -PNG:Author='Ann' -PNG:Description='A blue rounded box' -PNG:Copyright='CC0' \
  -PNG:CreationTime#='Tue, 1 Sep 2026 10:20:30 +0900' -PNG:Comment='ignored: Description comes first' "$out/tags.png"

# GIF: two frames and a comment.
magick -delay 50 -size 16x16 xc:'#e53935' xc:'#1e88e5' -loop 0 -set comment 'A blinking square' "$out/anim.gif"

# WebP: an EXIF orientation browsers do not apply, and XMP.
cwebp -quiet -q 90 "$tmp/scene.png" -o "$out/scene.webp"
x -Orientation#=6 -XMP-dc:Title='Scene (WebP)' -XMP-dc:Creator='Ann' "$out/scene.webp"

# AVIF: turned by irot (90° anticlockwise), with Exif and XMP items.
cat > "$tmp/xmp.xml" <<'XMP'
<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmp:CreateDate="2026-09-03T12:00:00Z">
<dc:title><rdf:Alt><rdf:li xml:lang="x-default">Scene (AVIF)</rdf:li></rdf:Alt></dc:title>
</rdf:Description></rdf:RDF></x:xmpmeta>
XMP
magick "$tmp/scene.png" "$tmp/exif.jpg"
x -EXIF:Artist='Bob' "$tmp/exif.jpg"
exiftool -b -exif "$tmp/exif.jpg" > "$tmp/exif.bin"
avifenc -q 80 --irot 1 --xmp "$tmp/xmp.xml" --exif "$tmp/exif.bin" "$tmp/scene.png" "$out/rotated.avif" > /dev/null

# BMP and an icon with two sizes.
magick -size 12x8 xc:'#fdd835' -fill '#6d4c41' -draw 'rectangle 0,0 5,7' BMP3:"$out/flag.bmp"
magick -size 48x48 xc:'#43a047' -fill '#fff' -draw 'circle 24,24 24,8' -define icon:auto-resize=16,32 "$out/icon.ico"

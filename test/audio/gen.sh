#!/bin/sh
# Regenerates the test files of the audio converter (converter/audio/testdata):
# a second of a 440 Hz tone in every format, with the tags and the cover art
# the converter reads. ffmpeg (6 or later, with libmp3lame, libvorbis and
# libopus) encodes the audio and writes the tags it can; the Python library
# mutagen (pip install mutagen) writes the frames ffmpeg writes differently
# or not at all: the lyrics and comment frames of ID3, chapters, several
# artists in one frame, and the pictures of Ogg files.
set -eu
cd "$(dirname "$0")/../.."
out=converter/audio/testdata
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
tone="sine=frequency=440:sample_rate=22050:duration=1"
f() { ffmpeg -y -loglevel error -nostdin "$@"; }

# The cover art: a sky, a field and a sun, as JPEG and as PNG.
f -f lavfi -i "gradients=s=320x320:c0=#2f5f9e:c1=#9fd0ff:x0=0:y0=0:x1=0:y1=320:n=2:d=1" \
  -vf "drawbox=x=0:y=240:w=320:h=80:color=#3a9d3a:t=fill,drawbox=x=200:y=60:w=64:h=64:color=#f6c344:t=fill" -frames:v 1 "$tmp/cover.png"
f -i "$tmp/cover.png" -q:v 4 "$tmp/cover.jpg"

# MP3 with an ID3v2.3 tag (and an ID3v1 tag at the end): a JPEG front cover
# and the usual text frames, one of them in UTF-16; mutagen adds the comment,
# the lyrics and two chapters below.
f -f lavfi -i "$tone" -i "$tmp/cover.jpg" -map 0:a -map 1:v -c:a libmp3lame -b:a 32k -c:v copy -disposition:v attached_pic \
  -id3v2_version 3 -write_id3v1 1 \
  -metadata title="Sunny Field" -metadata artist="Taro Yamada" -metadata album="Four Seasons" -metadata album_artist="Various Artists" \
  -metadata composer="Hanako Sato" -metadata genre="Folk" -metadata date="2026-09-01" -metadata track="3/12" -metadata disc="1/2" \
  -metadata publisher="Example Records" -metadata copyright="© 2026 Example Records" -metadata language="eng" -metadata TSRC="JPA012600001" \
  "$out/tagged.mp3"

# MP3 with an ID3v2.4 tag in UTF-8 and no cover: Japanese text, two artists
# in one frame, lyrics in Japanese (the card shows a placeholder for the
# cover).
f -f lavfi -i "$tone" -c:a libmp3lame -b:a 32k -id3v2_version 4 \
  -metadata title="晴れた野原" -metadata album="四季" -metadata composer="佐藤花子" \
  -metadata genre="フォーク" -metadata date="2026" -metadata track="1" -metadata language="jpn" \
  "$out/japanese.mp3"

# M4A (AAC in an MP4 container) with iTunes-style metadata and a PNG cover.
# The movie box comes after the media data, as the mov muxer writes it.
f -f lavfi -i "$tone" -i "$tmp/cover.png" -map 0:a -map 1:v -c:a aac -b:a 32k -c:v copy -disposition:v attached_pic \
  -metadata title="Sunny Field" -metadata artist="Taro Yamada" -metadata album="Four Seasons" -metadata album_artist="Various Artists" \
  -metadata composer="Hanako Sato" -metadata genre="Folk" -metadata date="2026-09-01" -metadata track="3/12" -metadata disc="1/2" \
  -metadata copyright="© 2026 Example Records" -metadata comment="A song about a field in the sun" -metadata grouping="Summer" \
  -metadata lyrics="Over the field the sun is high
Not a cloud in the summer sky" \
  "$out/tagged.m4a"

# FLAC with a Vorbis comment and a picture block.
f -f lavfi -i "$tone" -i "$tmp/cover.png" -map 0:a -map 1:v -c:a flac -c:v copy -disposition:v attached_pic \
  -metadata title="Sunny Field" -metadata artist="Taro Yamada" -metadata album="Four Seasons" -metadata composer="Hanako Sato" \
  -metadata genre="Folk" -metadata date="2026-09-01" -metadata track="3" -metadata TRACKTOTAL="12" -metadata ISRC="JPA012600001" \
  -metadata LYRICIST="Hanako Sato" -metadata copyright="© 2026 Example Records" -metadata comment="A song about a field in the sun" \
  "$out/tagged.flac"

# Ogg Vorbis and Opus; mutagen adds the picture (METADATA_BLOCK_PICTURE).
f -f lavfi -i "$tone" -c:a libvorbis -q:a 0 \
  -metadata title="Sunny Field" -metadata artist="Taro Yamada" -metadata album="Four Seasons" -metadata date="2026" -metadata track="3" \
  "$out/tagged.ogg"
f -f lavfi -i "$tone" -c:a libopus -b:a 24k \
  -metadata title="Sunny Field" -metadata artist="Taro Yamada" -metadata album="Four Seasons" -metadata date="2026" \
  "$out/tagged.opus"

# WAV with a LIST INFO chunk and a Broadcast Wave extension.
f -f lavfi -i "$tone" -c:a pcm_s16le -ar 8000 -ac 1 -write_bext 1 \
  -metadata title="Sunny Field" -metadata artist="Taro Yamada" -metadata album="Four Seasons" -metadata date="2026-09-01" \
  -metadata genre="Folk" -metadata comment="A song about a field in the sun" -metadata copyright="© 2026 Example Records" \
  -metadata description="A take recorded outdoors" -metadata originator="Field Recorder" -metadata origination_date="2026-09-01" \
  "$out/tagged.wav"

# AIFF with an ID3 chunk (v2.3) holding the tags and the cover.
f -f lavfi -i "$tone" -i "$tmp/cover.jpg" -map 0:a -map 1:v -c:a pcm_s16be -ar 8000 -ac 1 -c:v copy -disposition:v attached_pic \
  -write_id3v2 1 -id3v2_version 3 \
  -metadata title="Sunny Field" -metadata artist="Taro Yamada" -metadata album="Four Seasons" -metadata date="2026" \
  "$out/tagged.aiff"

# Raw AAC (ADTS), which has no tags: the card page shows what the frames say.
f -f lavfi -i "$tone" -c:a aac -b:a 32k -f adts "$out/plain.aac"

python3 - "$out" "$tmp/cover.png" <<'PY'
import base64, sys
from mutagen.id3 import ID3, USLT, COMM, TPE1, CHAP, CTOC, CTOCFlags, TIT2
from mutagen.oggvorbis import OggVorbis
from mutagen.oggopus import OggOpus
from mutagen.flac import Picture

out, cover = sys.argv[1], sys.argv[2]

# ID3v2.3: a comment and lyrics as COMM and USLT (ffmpeg writes them as
# TXXX), and the chapters of a podcast.
t = ID3(out + "/tagged.mp3")
t.delall("TXXX")
t.add(COMM(encoding=1, lang="eng", desc="", text="A song about a field in the sun"))
t.add(COMM(encoding=0, lang="eng", desc="iTunNORM", text=" 00000012 00000000"))
t.add(USLT(encoding=1, lang="eng", desc="", text="Over the field the sun is high\nNot a cloud in the summer sky\n\nWalking home along the lane\nHoping it will never rain"))
t.add(CHAP(element_id="ch1", start_time=0, end_time=500, start_offset=0xFFFFFFFF, end_offset=0xFFFFFFFF, sub_frames=[TIT2(encoding=0, text=["Verse"])]))
t.add(CHAP(element_id="ch2", start_time=500, end_time=1000, start_offset=0xFFFFFFFF, end_offset=0xFFFFFFFF, sub_frames=[TIT2(encoding=0, text=["Chorus"])]))
t.add(CTOC(element_id="toc", flags=CTOCFlags.TOP_LEVEL | CTOCFlags.ORDERED, child_element_ids=["ch1", "ch2"], sub_frames=[]))
t.update_to_v23()
t.save(v2_version=3, v1=2)

# ID3v2.4: two artists in one frame (separated by NUL) and Japanese lyrics.
t = ID3(out + "/japanese.mp3")
t.delall("TXXX")
t.add(TPE1(encoding=3, text=["山田太郎", "佐藤花子"]))
t.add(USLT(encoding=3, lang="jpn", desc="", text="野原の上に太陽\n雲ひとつない夏の空"))
t.save(v2_version=4, v1=0)

# Ogg: the front cover as a FLAC picture block in base64.
pic = Picture()
pic.type, pic.mime, pic.desc = 3, "image/png", "Front cover"
pic.width, pic.height, pic.depth = 64, 64, 24
pic.data = open(cover, "rb").read()
encoded = base64.b64encode(pic.write()).decode("ascii")
for cls, name in ((OggVorbis, "tagged.ogg"), (OggOpus, "tagged.opus")):
    f = cls(out + "/" + name)
    f["metadata_block_picture"] = [encoded]
    f.save()
PY

ls -l "$out"

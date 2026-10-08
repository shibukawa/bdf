# Audio files

An audio file has no pages, but it has a picture and words: the cover art and the tags that players show. `converter/audio` turns an MP3, AAC (M4A, or raw ADTS), FLAC, Ogg Vorbis or Opus, WAV or AIFF file into a document of one page, a card: the cover across the top, and under it the title, the artists and the album, a table of the other fields, the lyrics and the chapters. The tags become the document's Dublin Core, so a file list can search and sort audio as it does documents, and a thumbnail of the document is the cover. The audio itself is not stored.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try [an MP3 with tags, cover art and lyrics](https://shibukawa.github.io/bdf/viewer/?file=samples/tagged.mp3).

## What is read

| Format | Extensions | Tags | Cover art | Length and coding |
|---|---|---|---|---|
| MP3 (and MPEG layers I and II) | .mp3 | ID3v2 (2.2, 2.3 and 2.4), ID3v1 | ID3v2 `APIC` | the frames: Xing/Info/VBRI headers, else every frame |
| AAC in MP4 | .m4a, .m4b, .m4p | iTunes-style `ilst` (including freeform `----` fields and QuickTime's `keys`), Nero chapters (`chpl`), XMP in a `uuid` box | `covr` | `mvhd`, the sample description (`esds`: AAC, HE-AAC, MP3; `alac`, `fLaC`, `Opus`, PCM …) |
| AAC, raw | .aac | an ID3v2 tag in front, if any | | the ADTS frames |
| FLAC | .flac | Vorbis comment (an ID3v2 tag in front fills what it lacks) | `PICTURE` blocks | `STREAMINFO` |
| Ogg | .ogg, .oga, .opus, .spx | Vorbis comment of Vorbis, Opus, FLAC and Speex streams | `METADATA_BLOCK_PICTURE` (and the older `COVERART`) | the identification header and the last page |
| WAV | .wav | RIFF `INFO` list, the Broadcast Wave extension (`bext`), an `id3 ` chunk | the ID3 chunk | `fmt ` and the size of `data` (RF64 too) |
| AIFF, AIFF-C | .aif, .aiff, .aifc | `NAME`, `AUTH`, `(c) `, `ANNO`, an `ID3 ` chunk | the ID3 chunk | `COMM` |

The page is as wide as the cover (a small cover is scaled up to 360 pt, a large one down to 720 pt), and the cover is stored as it is, in the format the file holds it in (JPEG or PNG, usually), as `converter/image` stores images; `-images convert` re-encodes it like the pictures of other documents. A file without a cover gets a placeholder square with a note on it. The text is set in the viewer's sans-serif font, like the text of HTML and Markdown documents, and wrapped on estimated widths. The card marks its structure for search and reading aloud: the title is a heading, the fields a table, the lyrics paragraphs and lines, the chapters a list.

Text in ID3's "ISO-8859-1" frames, ID3v1, RIFF `INFO` and AIFF's text chunks is read as UTF-8 when all of it is valid UTF-8, else as Shift_JIS when all of it is, else as Windows-1252, the way the MIDI converter reads its text. Several artists in one frame are read as several (NUL-separated in ID3v2.4, semicolons elsewhere), the numbered genres of ID3v1 (`(17)`, `17`) by their names, and a date written as year, day-month and time frames (ID3v2.3) as one date.

## Metadata

The tags map onto the Dublin Core elements this way (XMP, when an MP4 file carries it, comes first, as it does for images); `type` is always `Sound` and `format` the file's media type.

| Element | ID3v2 | MP4 | Vorbis comment | RIFF INFO, AIFF |
|---|---|---|---|---|
| `title` | `TIT2` | `©nam` | `TITLE` | `INAM`; `NAME` |
| `creator` | `TPE1` | `©ART` | `ARTIST` | `IART`; `AUTH` |
| `contributor` | `TPE2` (album artist), `TCOM`, `TEXT`, `TPE3`, `TPE4` | `aART`, `©wrt` | `ALBUMARTIST`, `COMPOSER`, `LYRICIST`, `CONDUCTOR`, `ARRANGER`, `REMIXER` | `IMUS`, `IWRI` |
| `subject` | `TCON` | `©gen`, `gnre` | `GENRE` | `IGNR` |
| `description` | `COMM` (not a program's own, such as iTunes') | `©cmt`, `desc` | `COMMENT`, `DESCRIPTION` | `ICMT`, `bext` description; `ANNO` |
| `publisher` | `TPUB` | `----:…:LABEL` | `PUBLISHER`, `LABEL`, `ORGANIZATION` | – |
| `date` | `TDRC`, or `TYER` + `TDAT` + `TIME`, or `TDRL` | `©day` | `DATE` | `ICRD`, `bext` origination date |
| `identifier` | `TSRC` | `----:…:ISRC` | `ISRC` | – |
| `relation` | `TALB` (the album) | `©alb` | `ALBUM` | `IPRD` |
| `language` | `TLAN` (ISO 639-2, written as BCP 47) | the track's `mdhd` | `LANGUAGE` | `ILNG` |
| `rights` | `TCOP` | `cprt` | `COPYRIGHT`, `LICENSE` | `ICOP`; `(c) ` |

The card also shows the subtitle, grouping, track and disc numbers, original date, BPM, the length, and the coding (codec, bits per sample of PCM and lossless files, bitrate, sample rate and channels); and the lyrics (`USLT`, or the text of `SYLT`; `©lyr`; `LYRICS`, `UNSYNCEDLYRICS`) and the chapters (ID3 `CHAP`, `chpl`, `CHAPTER001` … of Vorbis comments). User-defined fields (`TXXX`, `----`) whose names are those of the fields above fill the ones the file does not set otherwise.

There is no `-param` for this format. See [design.md §3.31](../design.md#331-音声ファイル--bdf-変換器converteraudio) for the internals.

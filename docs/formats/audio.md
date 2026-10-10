# Audio files

An audio file has no pages, but it has a picture and words: the cover art and the tags that players show. `converter/audio` turns an MP3, AAC (M4A, or raw ADTS), FLAC, Ogg Vorbis or Opus, WAV or AIFF file into a document of one page, a card: the cover across the top, and under it the title, the artists and the album, a table of the other fields, the lyrics and the chapters. The tags become the document's Dublin Core, so a file list can search and sort audio as it does documents, and a thumbnail of the document is the cover. The card plays: the document carries the file as it is, and the viewer hands it to the browser's own player. When the lyrics say when each line is sung, the viewer marks the line on the card as the song goes.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try [a song with cover art and timed lyrics](https://shibukawa.github.io/bdf/viewer/?file=samples/synced.mp3) and [an MP3 with more tags, lyrics and chapters](https://shibukawa.github.io/bdf/viewer/?file=samples/tagged.mp3). Press ▶ or the space bar.

## Playing

The audio is stored whole and unchanged, as an `audio` part that the view's `play` names with its media type ([spec §4.4](../spec.md#44-演奏play)). Nothing is decoded or re-encoded on the way, so the converter needs no codec and the sound is exactly the file's. The browser decodes it, which decides what plays: MP3, AAC, FLAC, Ogg Vorbis and Opus, and WAV play in the current browsers; AIFF and Apple Lossless play in Safari only. A file the browser cannot play still shows its card, and the viewer says that it does not play.

The demo viewer plays it with the controls it plays scores with: ▶ and ■, the space bar, the speed in percent (the pitch stays), and a slider for the position. The mini viewer of the sample servers shows the browser's own controls under the card.

The audio makes the document as large as the file, and the converter holds it in memory. Three things keep that in hand:

| | Command line | Go |
|---|---|---|
| Leave the audio out | `-param play=false` | `audio.Options.NoPlay` |
| The largest file whose audio is stored (256 MiB by default) | `-param maxaudio=64` (in MiB) | `audio.Options.MaxAudio` (in bytes) |
| A conversion for a thumbnail | `-pages 1` | `converter.Options{Pages: converter.PageList(1)}` |

A file over the limit is still converted: it gets its card without the audio, and a warning. In the browser the default limit is 64 MiB, because the converter module holds a file several times over while it converts; the `params` option of the converter's `convert` and `open` takes the same names (`{ maxaudio: "128" }`, `{ play: "false" }`). A conversion of selected pages leaves the audio out because a thumbnail or a file list has no use for it; `-param play=true` stores it all the same. Audio protected by DRM is never stored.

## Lyrics that follow the song

Lyrics are shown on the card either way. When they carry the time of each line, the view also gets cues ([spec §4.4](../spec.md#44-演奏play)): the rectangle of each line on the card, and the millisecond it starts at. The viewer marks the line being sung, scrolls to keep it in view, and plays from a line you click.

Timed lyrics are read from:

- **ID3 `SYLT`** frames whose content type is lyrics and whose time stamps are in milliseconds. Entries may be lines, or words and syllables with line feeds between the lines.
- **LRC text in a lyrics field**: `USLT`, `©lyr`, and the `LYRICS`, `UNSYNCEDLYRICS` and `SYNCEDLYRICS` Vorbis comments. This is the form most tagging programs write.

```
[ti:Sunny Field]
[offset:+250]
[00:01.50]Over the field the sun is high
[00:04.50]Not a cloud in the summer sky
[00:07.25]
[00:08.00][00:20.00]Walking home along the lane
```

Of an LRC text, the time stamps and the tags that describe the file (`[ti:…]`, `[ar:…]`, `[offset:…]` and the like) are not shown; `offset` shifts the times. A line with several stamps is shown once for each, in the order it is sung. A stamp alone is a break in the singing, and starts a new stanza. The per-word stamps of the enhanced format (`<00:08.40>`) are dropped: lines are followed, not words. A text is taken as LRC only when its stamped lines are at least as many as its other lines, so ordinary lyrics with a heading like `[Chorus: All]` stay as they are.

When a file has both, `SYLT` is used. When it has no timed lyrics but has chapters (audiobooks, podcasts), the chapters get the cues instead: the viewer marks the chapter being played, and a click on a chapter plays from it.

Not read yet: an `.lrc` file beside the audio file, and `SYLT` stamps that count MPEG frames (such lyrics are shown without times).

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

The card also shows the subtitle, grouping, track and disc numbers, original date, BPM, the length, and the coding (codec, bits per sample of PCM and lossless files, bitrate, sample rate and channels); and the lyrics (`SYLT`, `USLT`; `©lyr`; `LYRICS`, `UNSYNCEDLYRICS`, `SYNCEDLYRICS`) and the chapters (ID3 `CHAP`, `chpl`, `CHAPTER001` … of Vorbis comments). User-defined fields (`TXXX`, `----`) whose names are those of the fields above fill the ones the file does not set otherwise.

The `-param` options are `play` and `maxaudio`, described under [Playing](#playing). See [design.md §3.31](../design.md#331-音声ファイル--bdf-変換器converteraudio) for the internals.

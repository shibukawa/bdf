# 音声ファイル

音声ファイルにページはありませんが、絵と文字はあります。プレイヤーが見せるカバーアートとタグです。`converter/audio` は MP3、AAC（M4A、または生の ADTS）、FLAC、Ogg Vorbis・Opus、WAV、AIFF のファイルを 1 ページの文書、つまりカードにします。上端いっぱいにカバーを置き、その下に題名・アーティスト・アルバム、ほかの項目の表、歌詞、章を並べます。タグは文書の Dublin Core になるので、ファイル一覧は文書と同じように音声を検索し並べ替えられ、文書のサムネイルはカバーになります。音声そのものは格納しません。

## 試してみる

ファイルを[ビューア](https://shibukawa.github.io/bdf/viewer/)にドロップするか、[タグ・カバーアート・歌詞付きの MP3](https://shibukawa.github.io/bdf/viewer/?file=samples/tagged.mp3) を試してください。

## 読むもの

| 形式 | 拡張子 | タグ | カバーアート | 長さと符号化 |
|---|---|---|---|---|
| MP3（MPEG レイヤー I・II も） | .mp3 | ID3v2（2.2・2.3・2.4）、ID3v1 | ID3v2 の `APIC` | フレーム。Xing/Info/VBRI ヘッダー、なければ全フレーム |
| MP4 の AAC | .m4a、.m4b、.m4p | iTunes 形式の `ilst`（自由形式の `----` フィールド、QuickTime の `keys` を含む）、Nero の章（`chpl`）、`uuid` ボックスの XMP | `covr` | `mvhd`、サンプル記述（`esds`: AAC・HE-AAC・MP3。`alac`、`fLaC`、`Opus`、PCM …） |
| 生の AAC | .aac | 先頭の ID3v2 タグ（あれば） | | ADTS のフレーム |
| FLAC | .flac | Vorbis コメント（先頭の ID3v2 タグは足りない項目を補う） | `PICTURE` ブロック | `STREAMINFO` |
| Ogg | .ogg、.oga、.opus、.spx | Vorbis・Opus・FLAC・Speex ストリームの Vorbis コメント | `METADATA_BLOCK_PICTURE`（古い `COVERART` も） | 識別ヘッダーと最後のページ |
| WAV | .wav | RIFF の `INFO` リスト、Broadcast Wave 拡張（`bext`）、`id3 ` チャンク | ID3 チャンク | `fmt ` と `data` の大きさ（RF64 も） |
| AIFF、AIFF-C | .aif、.aiff、.aifc | `NAME`、`AUTH`、`(c) `、`ANNO`、`ID3 ` チャンク | ID3 チャンク | `COMM` |

ページの幅はカバーの幅です（小さいカバーは 360 pt まで拡大し、大きいカバーは 720 pt まで縮小します）。カバーは `converter/image` が画像を格納するのと同じく、ファイルが持つ形式（たいてい JPEG か PNG）のまま格納します。`-images convert` はほかの文書の画像と同じく再エンコードします。カバーのないファイルには音符を描いたプレースホルダーの正方形を置きます。文字は HTML・Markdown の文書と同じくビューアのサンセリフのフォントで組み、推定した幅で折り返します。カードは検索と読み上げのための構造を持ちます。題名は見出し、項目は表、歌詞は段落と行、章はリストです。

ID3 の「ISO-8859-1」のフレーム、ID3v1、RIFF の `INFO`、AIFF のテキストチャンクの文字列は、MIDI の変換器と同じ読み方で、すべてが正しい UTF-8 なら UTF-8、すべてが Shift_JIS なら Shift_JIS、そうでなければ Windows-1252 として読みます。1 つのフレームに入った複数のアーティスト（ID3v2.4 は NUL 区切り、ほかはセミコロン区切り）は複数の値に、ID3v1 の番号のジャンル（`(17)`、`17`）は名前に、年・日月・時刻のフレームに分かれた日付（ID3v2.3）は 1 つの日付にします。

## メタデータ

タグは次のように Dublin Core の要素に写します（MP4 ファイルが XMP を持てば、画像と同じくそれが先）。`type` は常に `Sound`、`format` はファイルのメディアタイプです。

| 要素 | ID3v2 | MP4 | Vorbis コメント | RIFF INFO、AIFF |
|---|---|---|---|---|
| `title` | `TIT2` | `©nam` | `TITLE` | `INAM`。`NAME` |
| `creator` | `TPE1` | `©ART` | `ARTIST` | `IART`。`AUTH` |
| `contributor` | `TPE2`（アルバムアーティスト）、`TCOM`、`TEXT`、`TPE3`、`TPE4` | `aART`、`©wrt` | `ALBUMARTIST`、`COMPOSER`、`LYRICIST`、`CONDUCTOR`、`ARRANGER`、`REMIXER` | `IMUS`、`IWRI` |
| `subject` | `TCON` | `©gen`、`gnre` | `GENRE` | `IGNR` |
| `description` | `COMM`（iTunes などプログラム自身のものは除く） | `©cmt`、`desc` | `COMMENT`、`DESCRIPTION` | `ICMT`、`bext` の説明。`ANNO` |
| `publisher` | `TPUB` | `----:…:LABEL` | `PUBLISHER`、`LABEL`、`ORGANIZATION` | – |
| `date` | `TDRC`、または `TYER` + `TDAT` + `TIME`、または `TDRL` | `©day` | `DATE` | `ICRD`、`bext` の作成日 |
| `identifier` | `TSRC` | `----:…:ISRC` | `ISRC` | – |
| `relation` | `TALB`（アルバム） | `©alb` | `ALBUM` | `IPRD` |
| `language` | `TLAN`（ISO 639-2 を BCP 47 で書く） | トラックの `mdhd` | `LANGUAGE` | `ILNG` |
| `rights` | `TCOP` | `cprt` | `COPYRIGHT`、`LICENSE` | `ICOP`。`(c) ` |

カードにはほかに副題、グループ、トラックとディスクの番号、原盤の日付、BPM、長さ、符号化（コーデック、PCM とロスレスのビット深度、ビットレート、サンプリング周波数、チャンネル）、歌詞（`USLT`、なければ `SYLT` の文字。`©lyr`。`LYRICS`、`UNSYNCEDLYRICS`）、章（ID3 の `CHAP`、`chpl`、Vorbis コメントの `CHAPTER001` …）を載せます。上の項目と同じ名前のユーザー定義フィールド（`TXXX`、`----`）は、ファイルがほかで設定していない項目を補います。

この形式に `-param` はありません。内部は [design.md §3.31](../design.md#331-音声ファイル--bdf-変換器converteraudio) を参照してください。

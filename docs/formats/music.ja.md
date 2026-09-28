# 楽譜

MML・MIDI・MusicXML は、いずれも共通の組版エンジン（`converter/internal/music`）が SMuFL のフォント Bravura で A4 のページに五線譜として組む。各 View はその音楽を Standard MIDI File として持つので、デモビューアで演奏できる。演奏の様子は[特徴 → 音楽](../features.ja.html)を参照してほしい。共通なのはそこまでで、MML と MIDI は演奏データから記譜そのものを組み立て直す必要があるのに対し、MusicXML はすでに書かれた記譜をそのまま読む。

## 試す

[ビューア](../../viewer/)にファイルをドロップするか、サンプルを試せる。[MML のチップチューン](../../viewer/?file=samples/frere.mml)、[歌詞付きの MIDI カラオケ](../../viewer/?file=samples/twinkle.kar)、[MusicXML のピアノ譜](../../viewer/?file=samples/minuet.musicxml)。

## MML（converter/mml）

MML と MIDI はどちらも演奏されたとおりのデータ（記譜済みではなく、時刻と長さのある音の並び）なので、組む前に共通の記譜の工程（`Notate`）を通す。音を音価に量子化し（拍や小節線をまたぐタイ、付点、3 連符）、小節と 1 段あたり最大 2 声部に分け、調（ファイルのもの、なければ音から推定したもの）で綴り、音域から音部記号を選ぶ（鍵盤楽器は大譜表になる）。

`converter/mml` は中身から 3 つの方言を判別する。汎用（FlMML と BASIC の `PLAY` 文）、マビノギの `MML@…;`、ファミコンの MCK・PPMCK のトラック行形式である。署名のない `.mml` ファイルで、MML と判別できる行が 2 行に満たないものは、中身に判別の手がかりが無いため拡張子だけで MML と判定される。

| オプション | 値 | 既定 |
|---|---|---|
| `dialect` | `auto`、`generic`（FlMML、MSX BASIC）、`mabinogi`（`MML@…;`）、`ppmck`（ファミコンの MCK・PPMCK） | `auto` |
| `octave` | `auto`、`normal`（`>` でオクターブが上がる）、`reverse`（`<` で上がる） | `auto` — 方言ごとの既定 |
| `time` | 拍子記号。例: `3/4`、`6/8` | `4/4` |
| `key` | 調号。例: `G`、`Bb`、`F#m`、`Em` | 音から推定 |
| `program` | 音色の指定が無いトラックの General MIDI プログラム番号（`0`〜`127`） | 方言ごとの既定 |

## MIDI（converter/midi）

`converter/midi` は Standard MIDI File のフォーマット 0・1・2 と RIFF の `RMID` を読む。トラック（またはトラックとチャンネルの組）ごとにパートになり、名前はトラック名、無ければ General MIDI のプログラム名を使う。カラオケの歌詞（`.kar` の文字、またはメタイベント 05）は、最も合うトラックの音符に付く。MML と同じく、MIDI の音符も演奏されたとおりのデータなので、上と同じ `Notate` の工程を経て記譜される。

| オプション | 値 | 既定 |
|---|---|---|
| `time` | 拍子記号（`3/4`、`6/8`、`C` は 4/4、`cut`） | ファイルのもの |
| `key` | 調号（`G`、`Bb`、`F#m`、または 2 つのフラットを表す `-2`） | ファイルのもの、無ければ音から推定 |

## MusicXML（converter/musicxml）

MusicXML は MML や MIDI と違い、演奏データから記譜し直すのではなく書かれたとおりに読む。複数の段のパート、声部、和音、連桁、連符、装飾音、臨時記号、タイとスラー、アーティキュレーションと装飾記号、強弱記号とヘアピン、節ごとの歌詞、音部記号・調号・拍子の変更、繰り返しと括弧、速度記号とリハーサルマーク、段の区切りは、すべてファイル自身の記述からそのまま読む。partwise・timewise のどちらの文書形式も、圧縮した `.mxl` も読む。演奏は繰り返しと、D.C.・D.S.・Fine・To Coda のような跳躍を展開して行う。

MusicXML に `-param` オプションは無い。楽譜は書かれたとおりに変換される。

共通の組版エンジン（小節の間隔、符幹と連桁、タイとスラー、歌詞の折り返しなど）の内部は [design.md §3.27](../design.html#327-楽譜と演奏convertermmlconvertermidiconvertermusicxmlconverterinternalmusic) を参照。

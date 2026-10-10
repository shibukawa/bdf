# otf: フォントファイルをメモリにマップして開く

`github.com/shibukawa/bdf/contrib/otf` は OpenType のフォントファイル（TrueType 系の .ttf / .ttc も CFF 系の .otf / .otc も）を、少ないメモリで読むために開くパッケージです。依存は標準ライブラリだけです。

フォントファイルの大半はグリフのデータです。和文フォントは 2 万以上のグリフを持ちますが文書が使うのは数百で、コレクション（.ttc）は複数のフォントを持ちますが欲しいのは 1 つです。そのファイルを `os.ReadFile` で読んでレイアウトすると、フォントを持っている間ずっと全体がヒープに残ります。19 MB のコレクションの 1 フェイスに 19 MB かかる計算です。`otf.Open` は、できるプラットフォームではファイルをメモリにマップします。結果はファイル全体の `[]byte` に見えるので、バイト列を受け取るパーサはそのまま使えますが、メモリを使うのは実際に読んだページ（テーブルディレクトリ、小さな表、使ったグリフ）だけで、ファイルが裏付けなのでメモリが足りなくなればカーネルが手放せます。

```go
f, err := otf.Open("/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc")
if err != nil {
	return err
}
// f.Bytes() がファイル全体。これまでどおりパースする。
font, err := sfnt.ParseIndex(f.Bytes(), 0)
// パースしたフォントを使う間は f を保持する（プロセスの寿命でもよい）。
// f.Close() はマップを解除し、その中を指すスライスを全部無効にする。
```

`OpenFS(fsys, name)` は `fs.FS` 向けの `Open` です。ファイルシステムが `*os.File` として開くとき（`os.DirFS` など）はマップし、それ以外（`embed` でプログラムに埋め込んだフォントなど）は全体を読みます。

## マップされるもの、読まれるもの

| プラットフォーム | 空でない通常ファイルの `Open` |
|---|---|
| Unix（Linux、macOS、BSD、Solaris、AIX） | `mmap(2)`、`PROT_READ`、`MAP_SHARED` でマップ |
| Windows | `CreateFileMapping` と `MapViewOfFile` で読み取り専用にマップ |
| js/wasm、wasip1、その他 | ファイルの大きさで 1 回だけ確保して全体を読む |

空のファイル、通常ファイルでないもの、プラットフォームが扱える大きさを超えるもの、マップの呼び出しが拒んだもの（一部のネットワーク・FUSE ファイルシステム）は全体を読みます。どちらになったかは `Mapped()` が返します。ファイルディスクリプタはマップができた時点で閉じ、マップはその後も残ります。

## API

| | |
|---|---|
| `Open(path string) (*File, error)` | パスでファイルを開く |
| `OpenFS(fsys fs.FS, name string) (*File, error)` | `fs.FS` のファイルを開く |
| `(*File).Bytes() []byte` | ファイル全体。読み取り専用。`Close` 後は `nil` |
| `(*File).Size() int` | ファイルの大きさ |
| `(*File).Mapped() bool` | 読み込みでなくマップされているか |
| `(*File).Close() error` | マップを解除する（またはバイト列を捨てる）。2 回呼んでもよい |

バイト列は二重の意味で読み取り専用です。マップは読み取り専用で、読み込みにフォールバックした場合に書いてもディスクには何も反映されません。マップ中に他のプロセスがファイルを切り詰めると、切り詰められたページを読んだ時点でフォールトし、Go はそこから回復できません。システムやフォントディレクトリのフォントはその場で書き換えられないので、このパッケージはそれを前提にしています。書き込み中のファイルはマップしないでください。

## 既存のパッケージとの違い

| パッケージ | 結果の形 | マップできない環境 | `fs.FS` | 備考 |
|---|---|---|---|---|
| `otf`（本パッケージ） | ファイル全体の `[]byte` | 全体を読む | 対応 | 読み取り専用。パスで開く。マップに失敗したら読み込みに切り替える |
| `golang.org/x/exp/mmap` | `*ReaderAt`（`io.ReaderAt`、`Len`、`At`） | 全体を読む | なし | `ReadAt` のたびにコピーする。`[]byte` を取るパーサにはもう 1 度コピーが要る |
| `github.com/edsrzf/mmap-go` | `MMap`（`[]byte`） | 非対応 | なし | 呼ぶ側が `*os.File` を開き `PROT` / `MAP` フラグを選ぶ。読み書き・コピーオンライトのマップ、`Flush`、`Lock` がある |
| `os.ReadFile`、`fs.ReadFile` | `[]byte` | — | `fs.ReadFile` | ファイル全体がヒープに載る |
| `golang.org/x/image/font/sfnt.ParseReaderAt` | パース済みのフォント | — | — | パーサ側で遅延する。`io.ReaderAt` を通して必要な表だけ読む。同じ問題への別の解だが、そのパーサに結び付く |

`otf` が `io.ReaderAt` でなくスライスを返すのは、bdf のパーサ（`internal/sfnt`、`internal/cff`、`font/woff2`）が、多くのフォントパーサと同じくバイト列を添字で読むからです。表はファイルの部分スライス、グリフは表の部分スライスで、コピーしません。`ReaderAt` だと表やグリフをそのたびにコピーするか、パーサをリーダー経由に書き直すことになります。マップならスライスの API を保ったまま、触ったページの分しか払いません。

## bdf の既存パッケージとの関係

`internal/fontdb`（変換器とラスタライザのためのフォントの探索・計測・サブセット化）は、読み込むフォントを `os.ReadFile` で読み、そのバイト列をプロセスの寿命で保持していました。いまは `otf.Open` / `otf.OpenFS` で開き、`*otf.File` を保持します。その先（`internal/sfnt`、`internal/cff`、`woff2`、`raster/imagebdf`）は同じ `[]byte` を受け取るので変更はありません。macOS でサンプル文書を `bdf generate -thumbnail` で変換したとき（システムフォントのヒラギノのコレクションは 7〜19 MB）:

| 入力 | 最大 RSS 前 → 後 | peak footprint 前 → 後 |
|---|---:|---:|
| basic.docx | 73 → 45 MiB | 60 → 28 MiB |
| basic.pptx | 53 → 41 MiB | 39 → 26 MiB |
| basic.xlsx | 54 → 38 MiB | 40 → 23 MiB |

出力はバイト単位で同一です。このパッケージを独立したモジュールにせず bdf モジュールの `contrib/` に置いているのは、`fontdb` が import するためです。独立モジュールにすると、ここで変更を使うたびに先に公開済みのバージョンが必要になります。単独で使う場合（`go get github.com/shibukawa/bdf/contrib/otf`）も、標準ライブラリしか使わないので bdf モジュール以外の依存は増えません。

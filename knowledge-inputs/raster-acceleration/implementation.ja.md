# AA指定と3施策の製品採用

2026-10-15、利用者の採用指示に基づき、AAの明示的無効化と、有望だった3施策を通常のimagebdfへ組み込んだ。AAの未指定・ONは従来の品質を維持し、dpiによる切替は追加していない。前段の[候補評価](results-non-aa.ja.md)は採用前の履歴であり、以下は採用後の実装と検証である。

## 採用した振る舞い

| 施策 | 通常実装 | 制御 |
| --- | --- | --- |
| AA無効化 | 図形・文字・幾何clipの画素中心判定。既存の透明度・合成・影・補間・細線補償を維持 | `Options.Antialias=AntialiasOff`、CLI `-antialias=off`。未指定はAAあり |
| 完成画像の再利用 | `Page` / `PageDPI` / `PageSize` / `Region` / `Renderer.Object` のLRU。保持画像と返却画像を分離し、呼び出し側の画素変更を隔離 | `ImageCacheBytes`。0は64MiB、負数は無効、正数はバイト上限 |
| 独立ページ並列 | `Pages` / `PagesDPI`。callbackは指定順・呼出goroutine。workerは再利用し、mutableなdecode/font/scratchは個別。文字マスクだけをlock付きLRUで共有 | workers=0は最大4を自動選択、1..16は上限。基本canvas合計32M画素を目安に並列数を制限 |
| 文字列全体の被覆率再利用 | fill textの最終変換・Font参照・文字列・サイズ・spacing・方向・描画範囲・AA方式をキーにする。塗り・alpha・clip・影は各回適用 | `TextCacheBytes`。0は32MiB、負数は無効、正数はバイト上限。stroke textは対象外 |

完成画像はレイアウト、領域、画素サイズ、背景、AA方式を区別する。Partバイト列やFontFSを更新したら `Renderer.Invalidate()` で完成画像、文字マスク、decoded資源、警告とworkerを破棄する。FontDirs・システムフォントには従来のプロセス共通scan/load cacheがあり、差し替え可能なフォントにはFontFSを使う。

容量は配列、キー、項目の保守的な見積もりを含む。Go allocator、decodedフォントや画像、layer、描画中の作業領域、callbackが保持する出力は別である。各workerの作業画像保持枠は合計64MiBを分配する。新規描画では保存用の画像コピーが1回増えるため、1回用のCLIと単独Object helperは完成画像保持を省く。persistentなRenderer.Objectは文書の依存資源も扱える。

Rendererと文書はPagesの完了まで変更しない。複数workerに渡すFontFSは並行読み取り対応が必要。キャンセルとcallbackエラーは描画・バッチの区切りで反映し、開始済み描画の終了を待って戻る。帯・タイル分割は採用していない。

## APIとCLI

```go
ren := imagebdf.New(doc, &imagebdf.Options{
    Antialias: imagebdf.AntialiasOff,
    // ImageCacheBytes: -1, // 完成画像を保持しない場合
    // TextCacheBytes: -1,  // 文字マスクを保持しない場合
})
err := ren.PagesDPI(ctx, doc.Views[0], nil, 300, 0,
    func(page int, img *image.RGBA) error {
        return savePage(page, img)
    })
```

`nil`は全ページ、ページ番号は0始まり。CLIは1始まりの `-pages=all` / `-pages=1-3,5` と `{page}` 出力名を使う。

```sh
bdf render -pages=all -workers=4 -antialias=off input.bdf 'page-{page}.png'
```

`-page`との併用、範囲外ページ、シートのbatchはエラー。AA・倍率は既存の指定と同じ。順次PNG等へ書くため、全ページ分の画像を配列に保持しない。

## 採用後の測定

Apple M3 / darwin arm64 / Go 1.27.0 / GOMAXPROCS=8、300dpi AAあり、warm、5×1sの中央値。RGBA返却を含みPNG符号化は除く。完成画像hitと文字マスクだけの効果を分けた。

| 文書 | 両cache無効 | 文字cacheのみ | 完成画像cache hit |
| --- | ---: | ---: | ---: |
| PDF由来 | 143.14ms | 98.61ms（1.45倍） | 3.38ms（42.4倍） |
| PPTX由来 | 186.93ms | 168.02ms（1.11倍） | 3.06ms（61.0倍） |
| draw.io由来 | 31.64ms | 30.61ms（1.03倍） | 0.82ms（38.5倍） |

PDF由来の先頭ページを4回描くbatchは、worker上限1で407.66ms、上限4で253.03ms（1.61倍）。完成画像cacheは無効にし、workerを保持したwarm比較である。32M画素枠で実並列数も制限される。以前の試作の無制限4worker比較とは同一条件ではない。

実行間の変動が大きい。短い予備測定では文字cacheのdraw.ioは遅くなる回もあり、小さな改善を保証しない。完成画像cacheは同じ出力を再利用するhit時の効果で、初回描画や異なるページの処理が38〜61倍になる意味ではない。[生データ・fingerprint](production/benchmark.json)と[出力ログ](production/benchmark.txt)を保存した。

## 確認したこと

- リポジトリ全体のGoテストが通過。
- AA・cache・Pages・CLIを対象にrace検査が通過。
- 3文書、72/300dpi、AA on/offで、文字cache無効／cold／warmの画素が一致。
- 完成画像のcold/warm一致、返却画像を書き換えた後の再利用、容量上限・LRU・無効化・FontFS更新、Regionと非ゼロ原点Objectを検証。
- Pagesの選択順、重複ページ、callbackでの画像変更、共有文書/font、警告集約、キャンセル、callbackエラー、無効な選択、CLIの単ページとの一致を検証。
- 採用前のoverlayとadapterのソースを保存し、過去のbaselineのビルド再現を確認。AA単独benchは両cacheを無効にし、採用後benchは `BenchmarkRenderCaches` / `BenchmarkPages` に分離。

[全体テストのログ](production/tests.txt)、[race検査ログ](production/race.txt)。実機プリンタの画質は今回も未検証である。

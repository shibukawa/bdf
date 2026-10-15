# AA以外の高速化候補の検証

2026-10-14〜15、Apple M3 / darwin arm64 / Go 1.27.0で検証した。結論は、**静的画像の再利用、独立ページの並列化、文字列全体のマスク再利用**を優先すること。ネイティブrendererにも大きな改善余地がある。疎なマスク・SIMD・RGBA8などは、単独のカーネルの速さから文書全体の効果を判断できなかった。

AAの選択は明示指定のまま。dpiによる自動切替や、AA OFFに伴う画素形式・画像補間の変更は加えていない。この評価時点の試作はGo overlayまたは独立ツールへ隔離した。その後、利用者の指示により3施策を採用した。[採用後の実装・検証](implementation.ja.md)を参照。ここに記載する数値は採用前の比較である。

## 比較方法

- 既存のPDF由来・PPTX由来・draw.io由来の3文書を72/300dpi、AA on/offで描画。12出力についてbaselineとの差、cold/warm画素一致、時間・alloc・プロセスCPU時間と最大RSSを記録した。
- Go試作は同じ計測用overlayをbaselineにも適用し、通常は `-cpu=1`、3×200ms。凸形状は5×500ms。小さな速度差には実行順・GC・環境変動が残る。
- 文字列cacheはbaseline→cache→baseline、300dpi AA on、7×500msでも確認。PDF由来文書では前後のbaseline差が約0.2%で、cacheは約1.33倍。PPTXは1.02〜1.07倍、draw.ioは1.08〜1.12倍で、こちらの差は小さい。
- ネイティブ比較は共通の1024×1024図形シーン4種類。AAは有効だがAA方式・線分化誤差は同一ではない。録画・構築済みgeometryからのwarm描画を5×300msで測り、描画のみ／所有する出力画像の取り出し込みを分けた。
- 独立ページ、帯、タイル、画像再利用、PNG出力、PDF生成は公開APIの別ツールで確認。並列ページは同じ先頭ページを4回処理するバッチであり、異なる4ページの代表分布とは区別する。

## 候補ごとの判断

| 候補 | 観測 | 判断 |
| --- | --- | --- |
| グリフ単位cache | 各グリフの被覆率をmax-unionしても、文字列全体の走査と一致しない反例を確認。隣接する半画素の占有情報も失われる | 汎用の置換には使わない。文字列全体cacheが安全な代案 |
| 文字列全体cache | 12出力の画素差0。PDF由来300dpi AA onで約1.33倍を前後baselineで確認 | 同じ文字列・位置・倍率の再描画に優先候補。単文字のページ内再利用とは異なる |
| 行RLE・直接span合成 | 修正済み試作は12出力の画素差0。300dpi AA onの速度比0.75〜1.12倍 | この実装を一律採用しない。行decode、clip、runの数や確保が利益を打ち消す |
| 小A8 mask | 12出力で741画素が変化、最大1段階。速度比0.86〜1.17倍 | 閾値32px/1024pxを固定採用する根拠は得られない |
| 単色融合・不透明store・単純形状 | 単色・不透明の試作は画素差0。既存の矩形専用経路は汎用走査より約3倍速いカーネル例。凸三角形・四角形の巻き数処理省略も画素差0 | 既存の特殊化を維持。追加分は文書全体では小さく、対象頻度を見て判断 |
| ARM64 NEON | full-coverage単色SourceOverを試作。カーネルは大幅に高速化するが文書全体の改善は小さい。float演算のFMAを既存Goに合わせて検証 | 部分被覆・gradient等を含む製品用SIMD pipelineの完成とは扱わない |
| 全面RGBA8／主画面RGBA8＋float層 | プロセス最大RSSはbaseline約524MiBに対し約238／259MiB。ただしPDF/PPTX由来300dpiは約0.53〜0.57倍に低下。通常文書の最大差1でも、低alpha200重合成で最大36段階差 | メモリ策として別途設計。品質互換の既定高速化には採用しない |
| flatten cache／stroke展開cache | 12出力の画素差0。stroke幅・パス変更の無効化を確認。大きな速度向上は安定せず | 変換・許容誤差・dashなどを正確なキーに含め、反復の多い用途へ限定 |
| 局所適応分割 | ループ・変曲点等でサンプル点と線分の距離を検証。文書ではAA onで最大25段階、OFFでは255段階の差 | 幾何誤差が許容内でも画素互換ではない。既定の置換にはしない |
| ページ・帯・タイル並列 | 既存文書の4ページ相当バッチは約1.9〜2.8倍。独立Rendererのページhash一致と72dpi race検査。帯・タイルには座標の丸め差と影の欠け | 独立ページを優先。Regionを分けるだけの実装は採用しない |
| ページ・領域・Object画像cache | 同一画像を保持してコピーする比較では、ページ約89〜195倍、共通シーンのRegion約70〜757倍、単純Object約31〜43倍 | 再描画を省く効果は大きい。初回高速化ではなく、cache hit時の比較。lookup・無効化・LRUの完成実装は対象外 |
| ネイティブCPU backend | Skia・Blend2D・Vello CPUの実行と画素差を確認。共通図形では大幅に速い | 導入価値あり。ただし全BDF機能のadapter、配布、FFI、フォント等は別途実装が必要 |
| GPU・CPU/GPU hybrid | Apple M3のMetal実機でupload・同期・readback込みを測定。静的Sceneでは取り出し込み約1.0〜3.7ms | 静的図形の速度は有望。CPU前処理・Scene構築を毎回含む動的描画の結果ではない |
| ベクタPDF出力 | 共通シーンの生成約0.9〜2.7ms、PDF約25〜45KB。ラスタPNGは約88〜888KB。後段GhostscriptのPNG出力は約73〜445ms | PDFを受け付ける印刷先ならアプリの全ページrasterを省ける。後段の仕事は残り、ローカルPNG生成が常に速くなるわけではない |

時間の比は対象入力に限る。表の高速化候補を組み合わせた結果や、600/1200dpi・別CPU・Wasmでの優位性は未検証。cacheの32MiB制限は試作のpayload上限で、mapやキーを含む厳密な総容量上限ではない。

## 外部rendererの共通図形比較

所有する画像の取り出し込み、中央値ms。Goは既存Rendererを再利用、SkiaはSkPicture、Blend2Dは構築済みpathを再生。Vello CPU freshは各回のcoarse処理も含め、cachedはそれを保持する。GPUは構築済みSceneを再生し、CPU側へ画像を戻す。

| backend | 不透明矩形 | cubic fill | 半透明重なり | cubic stroke |
| --- | ---: | ---: | ---: | ---: |
| imagebdf | 10.94 | 28.28 | 122.58 | 88.99 |
| Skia CPU | 1.03 | 5.90 | 30.41 | 6.37 |
| Blend2D synchronous | 0.70 | 1.88 | 5.95 | 3.90 |
| Blend2D thread_count=4 | 0.42 | 1.11 | 3.27 | 2.00 |
| Vello CPU u8 fresh | 1.14 | 2.58 | 6.33 | 3.32 |
| Vello CPU f32 fresh | 2.08 | 3.72 | 8.79 | 4.06 |
| Vello CPU u8 cached | 0.69 | 0.77 | 3.60 | 0.58 |
| Vello CPU u8 fresh, threads=4 | 0.66 | 1.48 | 3.42 | 1.73 |
| Vello compute GPU | 1.17 | 1.35 | 2.42 | 1.34 |
| Vello hybrid GPU | 1.31 | 1.25 | 3.74 | 1.01 |

異なるAA・flatten・丸めにより出力は一致しない。strokeは全て平端・miter limit 10へ揃えたが、輪郭生成の差は残る。最大差だけで品質の良否は判断しない。Vello CPU/GPUのfilterやmask等には固定commit時点の制限があり、文書全体の置換を実証した比較ではない。GPU試作の初期化値は両GPU rendererを生成する値なので、単独backendの起動時間として比較しない。

取り出した画像はnative側の画素順序のまま。Goのimage.RGBAへの完全な形式変換やGo/cgo境界費用は含まない。Go公開API比較のGOMAXPROCSは8であり、前節の単一CPU overlay比較と計測条件を混同しない。

## 並列化で見つかった品質条件

既存文書の帯・タイルでは、300dpiに最大1段階の差が残る場合がある。公開Region APIのfloat32領域を各タイルに変換し直す方法では、元ページと完全に同じdevice transformを維持できない。

影・soft mask・groupを含む専用入力では、haloなしの帯に72dpiで最大72段階、300dpiで最大80段階の差が出た。固定64px haloでも300dpiのタイルに最大100段階の差が残る。影のblur・offset・変換から必要範囲を計算する必要がある。全ページをhaloにする診断経路は画素差0だが、全ページを重複して描く費用がある。

独立ページはRendererを共有せずに処理する。ページの描画順や領域内合成順を変えず、cache/scratchも各workerが持つ。共有する読み取り専用Documentについてrace検査を行った。

## 保存した証拠と残る範囲

- [Goの生データと画素比較](experiments/go/summary.json)、[文字列cacheの前後baseline比較](experiments/bracket/summary.json)。
- [外部backend・公開APIの結果と画素差](experiments/native/summary.json)。commit、アダプターhash、Rust lockfileを保存。
- [対象一覧](experiments/coverage.yaml)、[多重合成の精度反例](experiments/bracket/before/precision.json)、[検証結果](experiments/verification/verification.json)。
- 再現ツールは `tools/raster-experiments/`。実験用のPNGは一時生成物、数値・hash・生ログはここに保存した。

今回完了したのは各候補の限定評価。製品への採用、一般文書のnative adapter、動的GPU場面のScene構築費用、厳密なcache容量管理、実機プリンタ・網点・QR読取品質は含まない。まずページ並列化と静的画像cache、次に文字列cacheを製品用に仕上げるのが妥当と判断する。

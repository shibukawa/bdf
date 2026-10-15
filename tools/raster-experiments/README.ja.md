# AA以外の高速化施策の比較実験

このディレクトリの候補比較は、3施策を製品採用する前のrendererを対象にした履歴である。採用後の性能は `go test ./raster/imagebdf -run '^$' -bench 'BenchmarkRenderCaches|BenchmarkPages'` で測る。過去のoverlayは `knowledge-inputs/raster-acceleration/experiments/snapshots.json` に保持し、`python3 tools/raster-experiments/replay.py --trial go/baseline` で再実行できる。`run.py` / `run-all.py` の生成器は採用前のコード用であり、採用後のキャッシュが候補比較に混入する実行を拒否する。

通常の描画実装を変えず、Goのoverlayで試作を読み込む。公開APIや既定品質の変更は行わない。実験全体の完了状況は `knowledge-inputs/raster-acceleration/experiments/coverage.yaml` に保存する。

```sh
python3 tools/raster-experiments/run-all.py --output /tmp/bdf-raster-experiments
```

個別に実行する場合:

```sh
python3 tools/raster-experiments/run.py --variant baseline --output /tmp/bdf-raster-experiments/baseline
python3 tools/raster-experiments/run.py --variant sparse-runs --output /tmp/bdf-raster-experiments/sparse-runs --reference /tmp/bdf-raster-experiments/baseline
```

同じ既存3文書、72/300dpi、AA on/offで比較する。入力のfingerprint、実行コマンド、生のベンチ結果、RGBA画素差、cold/warm一致、プロセスのuser/system CPU秒・最大RSSを保存する。`B/op`と最大RSS、カーネル単独と文書全体を区別する。画像出力テストはベンチ時間に含めない。

現時点の試作:

- `solid-fused`: 単色SourceOverの色取得と合成を融合。
- `opaque-fill`: 不透明・完全被覆・単色の一括書き込み。
- `flatten-cache`: 変換と許容誤差・パス更新版を含む32MiB上限の線分化cache。
- `text-run-cache`: グリフの重なりを保つ文字列全体の被覆率cache。グリフ単体cacheとは別。
- `sparse-runs`: 非ゼロ被覆率を行RLEで生成し、単純clipでは直接区間を合成。
- `small-a8`: 小さいパスの被覆率をA8へ量子化。品質差を測る候補。
- `rgba8`: 全描画面とgroup/mask面を4B/画素で保持し、描画ごとに丸める試作。
- `simd`: ARM64 NEONの単色SourceOverカーネル。非ARM64はscalar fallbackでありSIMD結果と扱わない。
- `stroke-cache`: path更新版・変換・line style・範囲を含むstroke展開cache。
- `adaptive`: De Casteljau局所分割。幾何誤差のサンプル検証と画素差を分ける。
- `hybrid-storage`: 主画面RGBA8、group/mask面float32。
- `convex`: 三角形・凸四角形で一般の巻き数処理を省く。一般の星形には使わない。

これらは単独の実験であり、組み合わせの正しさを保証しない。採用前に既存の画像goldenと限界入力・各合成を追加検証する。特にA8とRGBA8は丸めが変わるため、速度だけで採用しない。

外部ライブラリは `/tmp` の専用環境へ取得し、rootのGo依存やグローバルなPython/Rust環境を変更しない。CPU/GPU・FFI・初期化・転送・readbackを分けて測定し、共通の図形シーンの測定を一般文書全体の比較と混同しない。

外部比較の再現には以下の専用環境を用いた。生成シーンは `go run ./tools/raster-experiments/cmd/scene` で作る。

- Python venv `/tmp/bdf-raster-native-venv`: `skia-python==144.0.post2`, `Pillow==12.3.0`, `numpy==2.5.4`。
- Blend2D: `58ca9460b4138af6e793184682e9324e0a8053b2`、AsmJit: `c8d13c48fcd943531d1e26b140ca058dfbc6c0a1`。Release static build。`native-blend2d.cpp` をC++17/O3でリンク。
- Vello: `db90a73ec7309dfa2bd45886b353162a9ea95744`。`native-rust/Cargo.toml.in` の `@VELLO_DIR@` を専用cloneへ置換。`Cargo.lock`は測定結果のnativeディレクトリに保存する。Cargo homeとtargetも `/tmp` へ分離。
- `native-batch.py` は全backendを順次測る。`rerun-vello.py` はVelloだけを更新する。GPU計測にはMetal等の実デバイスが必要。
- `cmd/pipeline` / `run-public-pipelines.py` は公開APIで完成画像再利用、独立ページ、Regionによる帯・タイル、PDF生成を比較する。固定64px haloと全ページhaloは診断用であり、製品のhalo決定法ではない。
- `cmd/parallel-check` は独立Rendererを使ったページ描画の画素hashを検証する。`go run -race ... -dpi 72` で共有文書のraceも確認する。
- `run-bracket.py` はbaseline→text-run-cache→baseline、300dpi AA on、7×500msで環境変動を確認する。
- `collect-native.py` と `summarize.py` で数値と差分を知識ベースへ出す。PNGは比較用の一時生成物である。

外部比較は各backendのAAを有効にしたが、AAアルゴリズムやflatten誤差は同一ではない。SkPictureや事前構築Scene等の再利用費用は別扱いとし、Vello CPUはcoarse処理を含むモードと保持するモードを両方測る。GPU warmは事前処理済みSceneからのupload・描画・同期を含み、owned-outputにはreadbackとCPUコピーを含む。GPU試作の初期化値は両GPU rendererを生成する費用であり、単一backendの起動時間として比較しない。

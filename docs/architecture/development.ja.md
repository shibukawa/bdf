# ビルドとテスト

Go 1.27 以上、そして TypeScript のパッケージとデモサイトのために Node.js が要ります。

## Go

```sh
go test ./...                          # すべてのパッケージ
go run ./cmd/bdf demo out.bdf          # フィクスチャ文書を書き出す
go build -tags bdf_noconv ./...        # 画像・WOFF2 のコーデックを含めないビルド（ブラウザ向けの構成）
GOEXPERIMENT=simd go build ./...       # Go 1.27、amd64/arm64: SIMD 版コーデック（amd64 は AVX2 必須）
```

`bdf_noconv` はブラウザに要らないコーデックを除きます（[design.md §2](../design.md#2-go-と-wasm-について)を参照）。CI はこの構成も別にビルド・テストしています。コーデックを含めた状態でしか動かない変更は、手元でもこの構成を確かめないと、CI まで気づきません。

## TypeScript

```sh
npm ci
npm run build           # tsc -b packages/core packages/render
npm test                # packages/*/test/*.test.mjs を node --test で（先にビルドする）
npm run test:golden     # Chromium でフィクスチャを描画し、testdata/golden/*.png と比較
npm run test:golden:update  # golden の PNG を今の描画結果で置き換える
```

golden テストには固定バージョンの Chromium（`playwright-core` の headless shell）が要ります。フル版の Chromium はテキストのラスタライズがわずかに違うため、golden は作られたのと同じビルドでしか一致しません。`npx playwright-core install chromium` で入れるか、同等の headless shell を `CHROMIUM_PATH` で指定してください。

## フィクスチャの再生成

`npm run testdata` は各 `converter/*/testdata/` のソースファイルから `testdata/` 以下をすべて作り直します（Go が必要）。CI はその結果をコミット済みのものとバイト単位で照合します。各形式のサンプル*入力*（`.pptx`・`.dxf`・`.epub` などのファイル自体）は形式ごとのスクリプトが生成し、多くは Python が要ります。

```sh
npm run test:pptx:gen      # python-pptx, msoffcrypto-tool
npm run test:xlsx:gen      # openpyxl
npm run test:dxf:gen       # ezdxf
npm run test:kicad:gen
npm run test:musicxml:gen  # (converter/musicxml/testdata から実行)
# … 形式ごとに 1 つ。全一覧は package.json の scripts を参照
```

## 手元でビューアを動かす

```sh
npm run demo         # http://127.0.0.1:8765/examples/viewer/.out/ — ビューア単体。?src= でリポジトリ内の任意の .bdf を開く
npm run site         # site/dist をビルド: フルのサイト（トップページ、ビューア、サムネイル・検索テキストのページ、ドキュメント）。Go が必要
npm run site:serve   # npm run site の後、配信もする
node examples/viewer/serve.mjs   # npm run demo と同じ
```

`?src=` はスタンドアロンのビューアにリポジトリ内の `.bdf` のパスを渡します。例: `http://127.0.0.1:8765/examples/viewer/.out/?src=/testdata/drawio/multipage.bdf`。`&layout=spread` で本のように開きます。`BDF_SITE_FONTS=dir1:dir2 npm run site` は既定のテスト用フォントの代わりにそれらのフォントをサイトと一緒に公開します。

## 任意の `.bdf` を描画・確認する

```sh
node test/render.mjs out.bdf pngdir/     # すべてのページを headless Chromium で PNG に
go run ./cmd/bdf ls out.bdf              # Part の一覧
go run ./cmd/bdf render -page 2 out.bdf page2.png  # 1 ページを純 Go で描画（ブラウザ不要）
```

## CI の出力とバイト単位で一致させる

CI の「testdata is up to date」ステップは、`ubuntu-latest`（amd64）上の `bdf generate` の出力をコミット済みのものとバイト単位で比較します。ただ、Go の浮動小数点の融合積和はアーキテクチャによって異なるため、Apple Silicon では同じコマンドでも Object の境界ボックスなどの下位ビットが数箇所変わることがあります。`GOARCH=amd64 go build` でクロスビルドし、その（Rosetta 上の）バイナリを使えば、手元でも CI と同じバイト列を再現でき、差分が本当の回帰かどうかを先に確かめられます。

新しい golden の PNG には、これとは別に、CI と同じ headless shell の Chromium ビルドが要ります。それに合うコンテナ（`mcr.microsoft.com/playwright:v1.56.1-noble`、または Chromium を入れた `node:24-bookworm`）の中で golden テストを動かせば再現できます。

各パッケージ・コマンドが公開するものは [API 一覧](../api.md)を、`bdf` のサブコマンドは [BDF コマンド](cli.ja.md)を参照してください。fixture、描画、ブラウザ変換の確認結果は[変換品質レポート](../quality.md)にまとめています。

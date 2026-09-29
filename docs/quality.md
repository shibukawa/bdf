# 変換品質レポート

golden 画像と一致しても、任意の文書が元ソフトと同じ見た目になるとは限りません。このレポートでは、固定入力の再現性、描画結果、ブラウザ内変換を分けて評価します。ここに示す数値は、2026-09-29 の作業ツリーで実行した検証の記録です。

## 今回確認したこと

| 観点 | 結果 | 確かめている範囲 |
|---|---:|---|
| 登録形式 | 27 形式 | ブラウザ用変換モジュールが公開する形式名を検査 |
| ブラウザ内の入力例 | 34 件すべて成功 | 変換、BDF ヘッダー、サムネイル、検索用テキスト。PDF 3 件はページ単位変換も順方向・逆方向で確認 |
| Go のパッケージテスト | 成功 | `go test ./...` |
| ブラウザ側テスト | 123 / 123 成功 | `npm test` |
| 描画 golden | 115 ケース、今回は未実行 | 固定 Chromium がこの環境になく、画像比較は開始できず |
| fixture の再生成 | amd64 で全件一致 | CI と同じバイト単位の比較。Apple Silicon の標準実行については下記の注記を参照 |

## 検証の読み方

変換品質に一つの合格率は付けていません。変換データの一致、画面上の画素、文字検索やページ構造は別の性質であり、ひとつの数値にまとめると、どこまで確認したのかが見えなくなるためです。

### 形式fixtureの一致

CI の `testdata is up to date` ステップは、代表入力を変換し、生成した BDF をコミット済み fixture とバイト単位で比較します。今回も同じ比較を実行し、`GOARCH=amd64` では全比較が一致しました。PDF、Office、CAD、Web 文書、画像、楽譜、フォントなど、CI が列挙する入力を対象にしています。

この確認は、同じ実装と固定入力から同じ BDF が作られることを示します。元アプリケーションとの意味的・視覚的な同等性を証明するものではありません。

Apple Silicon の標準 `arm64` 実行では、デモ fixture のバイト比較が一致しませんでした。リポジトリの[ビルドとテスト](architecture/development.ja.md#ci-の出力とバイト単位で一致させる)に記載のとおり、浮動小数点の積和演算がアーキテクチャで異なり、境界ボックスなどの最下位ビットが変わる場合があります。amd64 実行で比較し直すと一致しました。

### 描画golden

115 件の固定画面を Chromium で描き、保存済み PNG と比較します。画像ごとに差がある画素を 2% まで許容し、メインスレッド描画と Worker 描画も同じ基準で照合します。ケースには PDF、PowerPoint、Excel、Word、CAD、draw.io、Web 文書、楽譜などのページ・シートが含まれます。

今回のローカル実行では、Playwright が要求する Chromium headless shell が見つからず、比較を開始できませんでした。ブラウザーのインストールも試しましたが、展開が5分以上進まなかったため中断しました。したがって「115」はテストに登録されたケース数であり、今回の比較結果ではありません。CI の golden 実行結果も、この作業からは確認していません。

この比較は、選んだ fixture の描画回帰を検出します。すべての入力形式を網羅するものではなく、フォントや OS が異なる環境の見た目を保証するものでもありません。

### ブラウザ内変換

Pages ビルドの smoke test は、公開サンプル 34 件を WASM 変換器へ渡します。変換後の BDF を開き、サムネイルと検索用テキストを作り、PDF については全体変換とページ単位変換のテキストを比較します。今回の実行では全サンプルが成功しました。

この検査はブラウザ用モジュールの接続とサンプル経路を確かめます。実利用者の任意ファイルや、すべての OS・ブラウザ・フォントの組み合わせを測定した結果ではありません。

## 今回の実行結果

- `go test ./...` — 成功
- `go vet -unreachable=false ./...` — 成功
- `GOEXPERIMENT=simd go test ./imgconv/...` — 成功
- `go test -tags bdf_noconv ./imgconv/ ./converter/... ./woff2/` — 成功
- `npm test` — 123 件成功、失敗 0 件
- `npm run test:golden` — Chromium headless shell がなく比較未実行
- `node test/site.mjs` — 成功
- `npm run site` — 生成成功

Go は 1.27、Node.js は 26.8.1 で実行しました。GitHub Actions は Ubuntu/amd64 と Node.js 24 を使うため、ここでの成功は Actions の実行結果そのものではありません。

## 品質上の境界

変換器は、入力形式が許す全機能を常に再現するわけではありません。省略、代用、縮小がある場合は変換時の警告と形式別の説明を確認してください。たとえば:

- PDF のメッシュシェーディングは平均色、関数型シェーディングの一部は代表色で描きます。[PDF の説明](formats/pdf.ja.md)
- Excel の数式はファイルに保存された計算結果を使います。結果が保存されていないセルは空になり、警告が出ます。現在日時に依存する条件付き書式は評価しません。[表計算の説明](formats/spreadsheet.ja.md)
- アニメーション画像は最初のフレームを表示します。画素数の大きい画像は上限内に縮小されることがあります。[画像の説明](formats/image.ja.md)
- EPUB の固定レイアウトは、文字を含むページをすべての CSS 配置ごと再現するものではありません。[EPUB の説明](formats/ebook.ja.md)
- 公開サンプル用フォントは全 Unicode 文字を含みません。変換ログに「フォントに字形がない」と出た文字は、ビューア側のフォント表示に委ねられます。[フォントの説明](formats/font.ja.md)

ここに挙げたものは境界の例です。形式別の全対応範囲は[形式一覧](formats/index.ja.md)と各形式のページにあります。

## 再確認するコマンド

```sh
go test ./...
go vet -unreachable=false ./...
GOEXPERIMENT=simd go test ./imgconv/...
go test -tags bdf_noconv ./imgconv/ ./converter/... ./woff2/
npm test
npm run test:golden
npm run site
node test/site.mjs
```

fixture のバイト比較は [CI workflow](https://github.com/shibukawa/bdf/blob/main/.github/workflows/ci.yml) の `testdata is up to date` ステップを使います。

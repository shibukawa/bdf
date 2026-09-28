# ドキュメント

bdf はブラウザでのプレビューのために作られた文書フォーマットです。この節では、それが何のためにあるか、何をするか、各部がどうつながっているかを説明します。実際の動きを見るには[デモサイト](https://shibukawa.github.io/bdf/)が一番です。

## 次に読むもの

- **[なぜ bdf か](why.ja.html)** — どんな問題を解決するか、なぜ PDF をそのまま描かないのか
- **[はじめかた](getting-started.ja.html)** — ファイルを変換して開くまで、5 分で
- **[特徴](features.ja.html)** — 文書が保つもの: 検索、選択、リンク、読み上げ用のテキスト、楽譜の演奏
- **[対応形式](formats/index.ja.html)** — すべての入力形式と、それぞれのレイアウト
- **[構成のサンプル](examples/index.ja.html)** — bdf をシステムに組み込む 3 つの方法。それぞれ実際に動く小さなサーバーです
- **[アーキテクチャ](architecture/index.ja.html)** — 変換と描画の仕組み、フォルダ構成、bdf コマンド、ビルドとテスト

## リファレンス（日本語）

詳しいリファレンスは今のところ日本語のみです。

- **[フォーマット仕様](spec.html)** — コンテナ、manifest、命令セット
- **[API 一覧](api.html)** — パッケージとコマンドの一覧
- **[設計メモ](design.html)** — 各変換器がその作りになっている理由

## ソース

ドキュメントは [`docs/`](https://github.com/shibukawa/bdf/tree/main/docs) に Markdown として置かれ、説明する対象のコードのそばにあります。このサイトはそれを描画したものです。訂正や質問は [issue や pull request](https://github.com/shibukawa/bdf) で歓迎します。

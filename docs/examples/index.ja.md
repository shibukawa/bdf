# 構成のサンプル

[`examples/`](https://github.com/shibukawa/bdf/tree/main/examples) 以下に、実際に動く小さな Go サーバーが 3 つあります。どれも同じトレードオフの上で選択が違うだけです: サーバーがあらかじめどこまで変換しておくか、文書を実際に開いたときにブラウザがどこまでするか。3 つとも [`examples/sample-files/`](https://github.com/shibukawa/bdf/tree/main/examples/sample-files) の同じ 6 ファイルを配信し、小さな埋め込みビューア [`examples/miniviewer`](https://github.com/shibukawa/bdf/tree/main/examples/miniviewer)（[ブラウザ側](../architecture/browser.ja.html)で説明）を共有しています。

| | [light-server](light-server.ja.html) | [preview-server](preview-server.ja.html) | [search](search.ja.html) |
|---|---|---|---|
| サーバーが bdf に変換 | しない | 1 回、キャッシュ | 1 回、キャッシュ |
| サーバーがサムネイルを作る | 1 回、キャッシュ | 1 回、キャッシュ | 1 回、キャッシュ |
| どこで変換するか | ブラウザ（開いたとき） | —（すでに bdf） | —（すでに bdf） |
| フロントが持つもの | wasm の変換器＋レンダラ | レンダラのみ | レンダラのみ |
| 追加のサービス | — | — | Meilisearch |
| 向いている場面 | 文書が多く、サーバーを軽くしたい・変換の余力がない | よく開かれる文書。最速で一定した表示時間 | 全文書を横断した全文検索 |

どれが「正解」というものではありません。[なぜ bdf か](../why.ja.html)を読み、実際に作るものに合わせて選んでください: ある文書がどれくらいの頻度で開かれ、どれくらいの頻度で変わるか、サーバーに変換の余力があるか、検索できる必要があるか、といった観点です。

## light-server

サーバーは何も変換しません。文書ごとに 1 度だけサムネイルを描き（そのためだけにメモリ上で変換し、結果は捨てます）、元のファイルをそのまま配信します。開くとブラウザの中で bdf に変換します。使うのは[デモサイト](https://shibukawa.github.io/bdf/)自身が配っているのと同じ wasm の変換モジュールです。これはサーバーの負担が一番小さい構成です（リクエストのたびに変換コードが動くことはありません）。代わりに、ブラウザは wasm モジュールを一度取得し、開くたびに手元で変換する時間がかかります。→ [詳しく読む](light-server.ja.html)

## preview-server

逆のトレードオフです。サーバーがすべての文書をあらかじめ 1 度 bdf に変換し、キャッシュします。文書を開くのはキャッシュ済みの bdf を取得するだけです（`net/http` の `ServeContent` が Range リクエストに自分で応えるので、表示中のページに要る部分だけしか転送しません）。ブラウザ側のフロントは変換器を持たず、レンダラだけです。文書はより速く、より一定した時間で開きますが、サーバーに変換の余力と結果を置く場所が要ります。→ [詳しく読む](preview-server.ja.html)

## search

preview-server のサーバー側変換の上に、検索インデックスを足したものです。各ページのテキスト（`Document.SearchText`）をページごとに 1 件として [Meilisearch](https://www.meilisearch.com/) に入れ、`/search` がクエリを Meilisearch に取り次ぎます。検索結果をクリックすると、見つかったページをそのまま開きます。→ [詳しく読む](search.ja.html)

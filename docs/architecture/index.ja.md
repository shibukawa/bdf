# 処理の流れ

経路は 2 つあります。どちらも同じ bdf の Part を作り、同じレンダラで描画します。

```mermaid
flowchart TB
    SRC["PDF・Excel・CSV・Parquet・PowerPoint・Word・Visio・draw.io・DXF・Jw_cad・SXF・CGM・HP-GL/2・Gerber・KiCad・TIFF<br/>Illustrator・Photoshop・HTML・Markdown・EPUB・MML・MIDI・MusicXML・フォント・画像"]

    subgraph SERVER["Go サーバープロセス"]
        direction TB
        SCONV["converter/pdf<br/>converter/xlsx<br/>converter/pptx<br/>converter/docx …"]
        BUNDLE["bdf バンドル<br/>（パック済み）<br/>manifest JSON<br/>描画命令<br/>画像・フォント"]
        PREVIEW["raster・thumbnail・SearchText<br/>サムネイル画像<br/>検索用のテキスト"]
        SCONV --> BUNDLE
        BUNDLE --> PREVIEW
    end

    subgraph BROWSER["ブラウザ"]
        direction TB
        subgraph CWORKER["変換 Worker（wasm）"]
            WCONV["converter/pdf<br/>converter/xlsx<br/>converter/pptx<br/>converter/docx …"]
        end
        PARTS["bdf 文書<br/>（メモリ上）<br/>manifest JSON<br/>描画命令<br/>画像・フォント"]
        subgraph RWORKER["レンダラ Worker"]
            LOADER["ローダ<br/>fetch・Range<br/>DecompressionStream"]
            STORE["Part キャッシュ"]
            RENDER["レンダラ<br/>OffscreenCanvas<br/>FontFace・Path2D"]
            TEXT["テキスト抽出・検索"]
            LOADER --> STORE
            STORE --> RENDER
            STORE --> TEXT
        end
        UI["ビューア UI<br/>（メインスレッド）<br/>canvas・テキスト選択層"]
        WCONV --> PARTS
    end

    SRC -- "① サーバーで変換" --> SCONV
    SRC -- "② ブラウザ内で変換" --> WCONV
    BUNDLE -- "1 ファイル形式 / 分割形式<br/>HTTP・CDN" --> LOADER
    PARTS -- "postMessage" --> LOADER
    RENDER -- "ImageBitmap" --> UI
    TEXT -- "テキスト run・ヒット矩形" --> UI
```

## ① サーバーで変換

Go のサーバープロセス内で `converter/pdf`、`converter/pptx`、`converter/xlsx` などのパッケージが元ファイルを変換し、manifest・描画命令・画像・フォントをパックした **bdf バンドル**にします。バンドルは 1 ファイル形式（先頭からのストリーミング読み、または Range による Part 単位の取得）か、オブジェクトストレージや CDN にそのまま置ける分割形式で配信します。ブラウザでは Worker 内のレンダラが必要な Part だけを読み込んで `OffscreenCanvas` に描画し、メインスレッドは受け取ったビットマップと透明なテキスト層を配置するだけです。

バンドルと一緒に、文書の一覧や検索エンジンが必要とするものも作れます。作るのは文書（パスワード付きの入力ならパスワードも）を持っている間です。Go のラスタライザ（`raster`、`thumbnail`）で描いたサムネイル画像と、文書のメタデータ・ページごとのテキストの JSON（`Document.SearchText`）です。実際に動くサーバーの例は[構成のサンプル](../examples/index.ja.html)を参照してください。

## ② ブラウザ内で変換

同じ変換パッケージを wasm にしたもの（`cmd/bdfwasm`）が Worker で動き、ユーザーが開いたファイルをメモリ上の bdf 文書に変換して、そのままレンダラに渡します。変換から描画までがブラウザ内で完結し、ファイルは外に出ません。[デモサイト](https://shibukawa.github.io/bdf/)はこの経路で動いています。Office 系の変換器は、サイトと一緒に公開したフリーフォントでテキストをレイアウトします（文書が使うものだけを取得します）。

PDF はページ単位で変換し（`converter.OpenStream`）、ビューアはまずすべてのページの大きさを受け取ってから、表示中のページを優先して変換できたページの中身を順に受け取り、最後に完成した文書と差し替えます。ビューアと Worker の組み立て方は[ブラウザ側](browser.ja.html)を、Go・TypeScript の API は [bdf コマンド](cli.ja.html)や [API 一覧](../api.html)を参照してください。

## 次に読むもの

- **[フォルダ構成](layout.ja.html)** — 各ディレクトリに何が入っているか
- **[bdf コマンド](cli.ja.html)** — CLI のサブコマンドを順に説明
- **[ブラウザ側](browser.ja.html)** — Worker、ビューア、wasm の変換モジュール
- **[ビルドとテスト](development.ja.html)** — Go・TypeScript のテストの実行、フィクスチャと golden の再生成

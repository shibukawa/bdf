# フォントとライセンス

BDF 自体は [MIT ライセンス](https://github.com/shibukawa/bdf/blob/main/LICENSE)です。bdf が含む・配るフォントには、それぞれのライセンスがあります。

## ライブラリに含まれるフォント

| フォント | 含まれる場所 | 出元 | ライセンス |
|---|---|---|---|
| **STIX Two Math** 2.13 b171 | パッケージ [`formula`](rendering.ja.md#数式を描く) に埋め込み（`go:embed`。フォント全体を無改変で）。`Options` で別のフォントを指定しないときの数式フォント | [stipub/stixfonts](https://github.com/stipub/stixfonts) のタグ `v2.13b171`、`fonts/static_otf/STIXTwoMath-Regular.otf` | [SIL Open Font License 1.1](https://github.com/shibukawa/bdf/blob/main/formula/fonts/OFL.txt) |
| **Bravura** 1.482 | 音楽記号の輪郭とメトリクス。`converter/internal/music/smufl` に生成した表として持つ（MML・MIDI・MusicXML の楽譜用） | [steinbergmedia/bravura](https://github.com/steinbergmedia/bravura) | [SIL Open Font License 1.1](https://github.com/shibukawa/bdf/blob/main/converter/internal/music/smufl/OFL.txt) |
| **NewStroke** | KiCad が文字を描く線の字体。`converter/kicad` に表として持つ | vovanium による NewStroke の配布物（ファイルは [`NEWSTROKE.txt`](https://github.com/shibukawa/bdf/blob/main/converter/kicad/NEWSTROKE.txt) に記載） | [CC0 1.0](https://creativecommons.org/publicdomain/zero/1.0/) |
| **DejaVu Sans** の Regular と Bold（ASCII のサブセット） | サンプル文書（`bdf demo`）を作るパッケージ `fixture` | [dejavu-fonts](https://dejavu-fonts.github.io/) | [DejaVu Fonts License](https://github.com/shibukawa/bdf/blob/main/fixture/testdata/fonts/LICENSE.txt) |

STIX Two Math: Copyright 2001-2021 The STIX Fonts Project Authors (https://github.com/stipub/stixfonts), with Reserved Font Name "TM Math". STIX Fonts™ is a trademark of The Institute of Electrical and Electronics Engineers, Inc.

Bravura: Copyright © 2015, Steinberg Media Technologies GmbH (http://www.steinberg.net/), with Reserved Font Name "Bravura".

リポジトリの STIX Two Math のファイルは、配布元のタグのファイルとバイト単位で同じです。SHA-256 と確かめるコマンドは [`formula/fonts/README.md`](https://github.com/shibukawa/bdf/blob/main/formula/fonts/README.md) にあります。

`formula` を import したプログラムは STIX Two Math を含みます。Open Font License は、商用かどうかを問わずソフトウェアにフォントを同梱・埋め込むことを認め、複製に著作権表示とライセンスを残すことを求めます。フォントは name テーブルにその両方を持ち、`OFL.txt` も隣にあります。自前のフォントで組む場合（`formula.Options.Math`）でも、埋め込んだフォントはリンクされます。

変換器は、文書の文字のためのフォントを持ちません。渡されたフォント（`-font-dir`、システムのフォント）で文字を組み、そのサブセットを書き出す文書に埋め込みます。それらのライセンスは、そのフォントのものです。

## このサイトで配っているフォント

デモのページはブラウザの中でファイルを変換し、文書が求めるフォントを `fonts/` から取得します。フォントはライセンスのファイルと一緒に公開しています。

| フォント | 代わりにするもの | ライセンス |
|---|---|---|
| Liberation Sans、Serif、Mono | Arial、Times New Roman、Courier New | SIL Open Font License 1.1 |
| Carlito | Calibri | SIL Open Font License 1.1 |
| Caladea | Cambria | copyright ファイル（`fonts/fonts-crosextra-caladea.copyright.txt`）に書かれたライセンス |
| IPAex ゴシック、IPAex 明朝 | 日本語 | IPA フォントライセンス v1.0 |
| DejaVu Sans、Serif、Mono | 記号 | Bitstream Vera Fonts のライセンス。DejaVu による変更はパブリックドメイン |
| STIX Two Math | Cambria Math、数式 | SIL Open Font License 1.1 |

ファイルは Debian のパッケージ `fonts-liberation`、`fonts-crosextra-carlito`、`fonts-crosextra-caladea`、`fonts-ipaexfont-gothic`、`fonts-ipaexfont-mincho`、`fonts-dejavu-core` と、上の STIX のリリースから取っています。一覧は[サイトをビルドするワークフロー](https://github.com/shibukawa/bdf/blob/main/.github/workflows/pages.yml)にあります。

## テスト用のフォント

リポジトリにあるテスト用のフォント（M PLUS 1p、STIX Two Math、STIX Two Text、DejaVu Sans のサブセット）は、ライセンスと一緒に `testdata` ディレクトリにあります。

# 偽プリンタ 設計案（print-m6）

2026-10-20 の設計案。同日に実装した（`.knowledge/` の `milestone:print-m6`。計画から変えた所と、偽物で見つけた不具合はそこの `facts.delivery`）。2026-10-20 利用者「このプロジェクトのテストのために、フェイクプリンタを作ると良さそうな気がしている。ネット上にラベルプリンタなどは見かけたが。他のプロトコルも」。`.knowledge/` の `requirement:fake-printers`、`milestone:print-m6`、`system:printer-emulators`。印刷の要件は `print-requirements.ja.md` と `requirement:print`。Print Manager（`print-manager.ja.md`）の開発と見せ方の土台にもなる。

**2026-10-20 利用者の決定**：「偽プリンタをまずきちんと作って、テストを全てそれをベースに行うのを先にしたい」。印刷の次の実装は print-m6（この偽プリンタ）で、Print Manager はその後。印刷のテストは全部この偽物の上に移す（§10）。

## 0. いま有るもの

テスト用の偽物はすでに 2 つある。どちらもプロセスの中で Go のテストから使う。

- `print/ipp/ipptest`（670 行）: IPP Everywhere のプリンタの `http.Handler`。受けたジョブと文書を持ち、形式・用紙・両面・色・解像度・PWG / URF / PCLm の能力・状態・理由を欄で変えられる。15 のテストのファイルが使う。
- `print/raw/rawtest`（270 行）: 9100 と LPD。受けたジョブを持ち、途中で切る（`CutAfter`、`CutControl`）、キューを断る（`Refuse`）、9100 で何か返す（`Back`）。
- `print/pdl/*_test.go` の読み手: `readZPL`・`readGFA`・`readSBPL`・`readTPCL`・`readESCPOS`・`readESCP`。書き手の出力を絵に戻して比べる。テストのファイルの中にあり、外からは使えない。
- 外の偽物: CI（`test/print/ci.sh`）は ippeveprinter（CUPS の IPP Everywhere のプリンタ）、cupsd のキュー、Python の 9100 の受け手、ipptool。これらは「bdf 以外の実装」なので、自前の偽物が増えても残す。

足りないのは次のこと。

1. **プロセスの外から使えない**。Print Manager の画面を作る・見せる、手で `bdf print` を試す、利用者が自分の設定を試す、のに「何台もの、種類の違うプリンタ」が要る。
2. **ラベル・レシート・連帳を受けて絵に戻す偽物が無い**（読み手がテストの中に閉じている）。
3. **壊れ方の組み合わせ**（用紙切れの後に再開、再起動で job-id が戻る、証明書が変わる、応答が遅い …）が、テストごとに手で書かれている。
4. **CUPS のサーバーとしての振る舞い**（`CUPS-Get-Printers`、`/printers/<名前>`、クラス）を持つ偽物が無い。bdf の `bdf printers` は CUPS の一覧を読むが、確かめるのは CI の本物の cupsd だけ。

## 1. 何を作るか

**1 つのプロセスで、種類の違う偽プリンタを何台も動かす「偽プリンタの農場」**。各プリンタは次の 4 つの組で決まる。

| 軸 | 値 |
|---|---|
| 搬送路 | `ipp`・`ipps`（1 つのポート）、`socket`（9100）、`lpd`、`serial`（pty。unix だけ）、`cups`（CUPS のサーバーの振る舞い。下の §3） |
| 言語（受けるもの） | PDF、PostScript、PWG Raster、URF、PCLm、JPEG、ZPL、SBPL、TPCL、ESC/POS、ESC/P、何でも（raw のまま持つ） |
| 能力（profile） | IPP の属性の組（用紙、両面、色、解像度、仕上げ、形式）と、ラベル機の dpi・印字幅。組み込みの名前つきのもの（§4）か YAML |
| 振る舞い（faults） | 時刻の筋書き：状態と理由、遅さ、途中で切る、断る、再起動 …（§5） |

受けたジョブは**記録**する：ジョブの属性（IPP）か制御ファイル（LPD）、受けたバイト列のまま、そして**読み戻した絵**（ページごとの 1 ビットか RGB の画像と、印刷側バーコードの命令の一覧）。

使い方は 3 つ。

1. **Go のテスト**: `fp := fakeprinter.Start(t, fakeprinter.Preset("zebra-zd421-203"))`、`fp.URL()`、`fp.Jobs()`、`fp.Pages(job)`、`fp.Fault(...)`。今の ipptest・rawtest と、テストごとの `httptest` の手書きの偽物は、全部この API に移して消す（§10）。
2. **農場のコマンド**: `go run ./tools/fakeprinter -config farm.yaml`（§6）。YAML に並べた偽プリンタを起こし、ジョブの一覧と絵を見る Web のページと、テストやスクリプトが使う JSON の API を出す。
3. **CI**: ci.sh の bdf の側の確かめは全部この偽物（農場）に向ける。外の実装（ippeveprinter、cupsd、ipptool、gs）は「偽物そのものが正しいか」を確かめる側に回す（§10.2）。9100 の Python の受け手は消す。

## 2. 読み戻し（受けたものを絵にする）

偽物の価値の半分は「刷られたはずの絵」を見せること。言語ごとの読み手：

| 言語 | 読み手 |
|---|---|
| PWG Raster、URF | `print/raster` の `PageReader`（m4 にある） |
| PDF、PCLm | `converter/pdf` で .bdf にして `imagebdf` で描く（PCLm は PDF の部分集合） |
| PostScript | 手元に `gs` があれば描く。無ければバイト列のまま（自前で読まない） |
| JPEG、PNG | そのまま |
| ZPL・SBPL・TPCL・ESC/POS・ESC/P | `print/pdl` のテストの読み手を `print/pdl/pdlread`（パッケージ）に出す。bdf が書く命令の部分集合を読み、知らない命令は一覧にして残す |

ラベルの言語の印刷側バーコード（`^BC`、`<BG>`、`{XB` …）は、読み手が `image/barcode` で描いて絵に入れ、命令の一覧にも残す（zxing で読める絵になるかを試せる）。

**比べ方**: 「bdf が送ったもの → 偽物の読み戻し」を、同じページの `imagebdf` の描画と比べる。ラスタの言語は 1 画素も違わないこと、PDF と PostScript は既存の golden と同じ許容。これで m1〜m5 の書き手を、プロトコルを通した後でも確かめられる。

## 3. CUPS のサーバーとしての偽物

`cups` の搬送路は、IPP のポートに CUPS の名前の空間を出す：

- `/printers/<名前>` と `/classes/<名前>`（クラスはメンバーの偽プリンタに流す）
- `CUPS-Get-Printers`、`CUPS-Get-Classes`、`CUPS-Get-Default`、`CUPS-Move-Job`、`CUPS-Accept-Jobs`・`CUPS-Reject-Jobs`
- `printer-type` のビット（クラス、リモート …）、`member-names`

`bdf printers -url ipp://localhost:<port>`（CUPS の一覧）を本物の cupsd 無しで試せる。Print Manager の「CUPS のサーバーとしても振る舞う」（`print-manager.ja.md` §5）を作るときの相手（`lpstat -h`、CUPS の `client.conf` の `ServerName`）を確かめる前に、bdf のクライアントとの往復をここで固める。

## 4. 組み込みの profile（暫定の候補）

| 名前 | 搬送路 | 中身 |
|---|---|---|
| `office-mfp` | ipp | IPP Everywhere。PDF・PWG・URF、A4・A3・B5・はがき、両面、色、ステープル |
| `airprint-only` | ipp | URF と JPEG だけ（PDF を受けない） |
| `pwg-only` | ipp | PWG Raster だけ。copies-supported が 1 |
| `ps-office` | ipp | PostScript と PDF |
| `pclm-inkjet` | ipp | PCLm だけ（Mopria） |
| `legacy-ipp11` | ipp | IPP/1.1 だけ、chunked を受けない（HTTP/1.0 の応答） |
| `cups-server` | cups | 上のうち 3 台とクラス 1 つ |
| `jetdirect-ps` | socket | PostScript（PJL の INFO の応答） |
| `lpd-legacy` | lpd | 何でも受けて持つ |
| `zebra-zd421-203`・`-300` | socket | ZPL、104 mm、`~HS` に答える |
| `sato-cl4nx` | socket | SBPL |
| `tec-bex4t1` | socket | TPCL |
| `epson-tm-t88` | socket | ESC/POS、80 mm、`DLE EOT` に答える |
| `epson-vp-f2000` | serial | ESC/P（pty） |

値は m5 の `print/pdl/profile.go` と同じ表から引く（2 か所に持たない）。

## 5. 振る舞い（faults）

筋書きは「いつ（ジョブの何件目・何ページ目・何バイト目・何秒後）」と「何が起きる」の並び。テストからも YAML からも同じ形で書く。

- **状態**: `stopped` と理由（`media-empty`、`media-jam`、`door-open`、`toner-empty`）、何秒後に `idle` に戻る
- **ジョブ**: `aborted` にする、`completed` にしない、`Get-Jobs` から消える（忘れる）、job-id を 1 に戻す（再起動）
- **搬送**: 応答の遅れ、帯域の上限（背圧）、n バイト目で切る、接続を受けない、`server-error-busy`、`server-error-service-unavailable`
- **IPP**: 形式・属性を断る、`client-error-not-authorized`（Basic を要る）、証明書を作り直す・期限切れにする、Get-Printer-Attributes の応答を巨大に（media-col-database が数千）
- **ラベル機**: `~HS` / `DLE EOT` に「紙が無い」「カバーが開いている」を返す

これでスプールの「二重に刷らない」「止まったプリンタで待つ」「途中で切れたジョブは人に任せる」「指紋の固定」を、決まった順で何度でも試せる。

## 6. 農場のコマンドと画面

```yaml
listen: 127.0.0.1          # 既定は loopback だけ
web: 127.0.0.1:8650        # ジョブの一覧と絵の画面、JSON の API
store: ./fakeprinter-jobs  # 受けたジョブを残す所（無ければメモリ）
dns-sd: false              # 既定は広告しない
printers:
  - name: mfp-3f
    preset: office-mfp
    port: 8701
  - name: label-1
    preset: zebra-zd421-203
    port: 9101
    faults:
      - after: {jobs: 3}
        state: {stopped: media-empty, for: 20s}
  - name: office
    preset: cups-server
    port: 8631
```

画面：偽プリンタの一覧（状態、受けたジョブの数）、ジョブごとの属性・受けたバイト列のダウンロード・読み戻した絵、状態を手で変えるボタン（「紙切れにする」「再起動」）。API：`GET /api/printers`、`GET /api/jobs?printer=…`、`GET /api/jobs/<id>/pages/<n>.png`、`POST /api/printers/<名前>/faults`、`POST /api/reset`。

**DNS-SD の広告は既定で切る**。広告するときは m4 の応答者で、ホスト名は自分のもの（`<機械>-fake.local`）にする。OS のホスト名の A / AAAA は告げない（2026-10-20 に macOS が機械の名前を変えた件。`mdns-hostname-conflict`）。

## 7. 外の偽物（調べたもの）

`system:printer-emulators`。どれも「その言語・プロトコル 1 つ」だけで、壊れ方の筋書きと、bdf の書き手との画素の比べ合わせは持たない。自前で作る理由はそこ。外のものは**別の実装による二重確認**に使う（スクラッチだけ。リポジトリに入れない）。

- ippeveprinter（CUPS / libcups）: IPP Everywhere。CI で使い続ける
- bizzabo/zplserver（Python）: 9100 の ZPL を PNG に。読み手の確認に
- gilbertfl/escpos-netprinter（Docker）: 9100 の ESC/POS を HTML に
- Labelary の Web API: ZPL を絵に。**外に送るので、公開してよい見本だけ**

## 8. 置き場所（暫定）

```
print/fakeprinter/      偽物の本体（device、profile、faults、記録、cups の名前の空間）。テストからも使う公開のパッケージ
print/pdl/pdlread/      ラベル・レシート・連帳の読み手（テストの読み手を出す）
print/ipp/ipptest/      移した後に消す
print/raw/rawtest/      同じ
tools/fakeprinter/      農場のコマンド（リリースのバイナリには入れない）と Web の画面
test/print/farm.yaml    CI と手元の農場
```

`bdf` のコマンドに `bdf fakeprinter` として入れる案もある（利用者が自分の設定を試せる）。リリースのバイナリが大きくなるのと、利用者向けの約束が増えるので、暫定では tools に置く（`requirement:fake-printers` の open_questions `command-home`）。

## 9. 順番（スライス）

1. **S1** `print/fakeprinter` の芯（device = 搬送路＋言語＋profile＋faults、記録）と、偽物そのものの確かめ（§10.2 の ipptool と CUPS のバックエンド）
2. **S2** `pdlread` と読み戻し、bdf print → 偽物 → 絵 の比べ合わせ（m1〜m5 の全部の形式）
3. **S3** faults の筋書きと、スプールのテストをそれで書き直す（再起動の二重印刷、止まって再開、途中で切れる、指紋）
4. **S4** `cups` の名前の空間（CUPS-Get-Printers ほか）と `bdf printers` の往復
5. **S5** `tools/fakeprinter`（YAML、画面、API）、`test/print/farm.yaml`、ci.sh を農場に向ける
6. **S6** テストの移行（§10.1 の全部）。ipptest・rawtest・手書きの偽物を消す。移す前と後で同じテストの数・同じ確かめが残ることを表で確かめる

## 10. テストを全部偽物の上に（2026-10-20 の決定）

### 10.1 移すもの（2026-10-20 に数えた）

| 今の偽物 | 使っているテストのファイル |
|---|---|
| `ipptest`（IPP） | print/print_test.go、formats_test.go、hooks_test.go、device_test.go、print/ipp/client_test.go、print/spool/spool_test.go、run_test.go、print/server/server_test.go、ipp_test.go、infra_test.go、client/client_test.go、cmd/bdf/printers_test.go、jobs_test.go |
| `rawtest`（9100・LPD） | print/rawprint_test.go、device_test.go、print/raw/raw_test.go、print/spool/raw_test.go、print/server/server_test.go、cmd/bdf/printers_test.go |
| テストごとの `httptest` の手書き | print/ipp/client_test.go（TLS のプリンタほか）、print/server/ipp_test.go、server_test.go、client/client_test.go、spool の 2 つ、cmd/bdf の 2 つ |
| pty（シリアル） | print/raw/pty_*_test.go → `serial` の搬送路の偽物（unix）。Windows は今のまま |
| pdl のテストの読み手 | print/pdl/*_test.go → `print/pdl/pdlread` |
| ci.sh の Python の 9100 | `test/print/raw9100.py` → 農場の `socket` |

移さないもの：

- bdf printd を相手にする試験（print/server の IPP サーバー側の試験は、bdf printd が「プリンタとして見せる」側なので偽物ではなく本物。その**先の**実プリンタを偽物にする）。
- Windows の winspool と USB の装置ファイル（OS の口なので偽物を立てられない。今のままの単体の試験）。

移した後の約束：印刷の Go のテストでプリンタ役を作るのは `print/fakeprinter` だけ（`httptest` で IPP を手書きしない）。レビューで分かるように、`print/` と `cmd/bdf` のテストで `httptest.NewServer` に IPP を答えさせているものが無いことを確かめる小さな検査を置く。

### 10.2 偽物そのものが正しいことの確かめ

全部のテストが偽物に乗るので、偽物が間違うと全部が間違う。外の実装で偽物を確かめる：

- **IPP**：ipptool の `ipp-1.1.test`、`ipp-2.0.test`、`ipp-everywhere.test`（PWG 5100.14）を偽物の各 profile に流して通る（m4 で bdf printd に行ったのと同じ）。
- **CUPS のクライアントから**：cupsd に偽物の `ipp`・`socket`・`lpd` のキューを作り（`lpadmin -v ipp://…/ socket://… lpd://…`）、`lp` で投げたジョブが偽物に届いて読み戻せる。`cups` の搬送路は `lpstat -h` と CUPS の `client.conf` の `ServerName` で一覧が出る。
- **差の比べ**：同じ IPP の要求の並びを ippeveprinter と偽物に送り、応答の状態コードと主な属性を比べる（ippeveprinter を CI に残す唯一の役）。
- **ラベル・レシート**：読み戻しの絵を zxing で読み、記号の中身が送ったものと一致する。外の偽物（zplserver、escpos-netprinter）との比べはスクラッチで 1 回。

### 10.3 決まっていない小さなこと（暫定で進める）

- 偽物の Go の API の形（`fakeprinter.Start(t, Preset(...), With...)` の関数オプション）。
- 記録の保存先（テストはメモリ、農場はディレクトリ）。
- 農場の画面の作り（print/server/web と同じ素の JS）。

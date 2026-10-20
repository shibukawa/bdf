# print-server: スプールを通して請求書を印刷する

[`examples/print-server`](https://github.com/shibukawa/bdf/tree/main/examples/print-server) — 請求書を印刷する小さな業務アプリケーションです。ページに入力した内容を請求書のテンプレートに差し込み、組み込んだ[プリントサーバー](../print-server.ja.md)のキューに文書を置き、スプールが送ってプリンタが刷るまでのジョブを見せます。ページはすぐに答え、印刷はその後ろで進み、再起動をまたいでも残ります。

```sh
cd examples/print-server
go run . -fake-printer          # bdf の偽の IPP プリンタ（同じプロセスで動く）に印刷する
open http://127.0.0.1:8083/

go run . -printer ipps://printer.example.local/ipp/print   # ネットワークのプリンタに印刷する
```

## サーバーがすること

- **テンプレートに差し込む。** `POST /invoices` は顧客、品目（1 行に 1 つ: 品名、数量、単価）、備考、利用者を読み、[`converter/xlsx/testdata/form-rows.xlsx`](https://github.com/shibukawa/bdf/tree/main/converter/xlsx/testdata) に差し込みます（[`form.Fill`](../forms.ja.md)）。タグ（`{{no}}`、`{{customer}}`、`{{#items}}` … `{{/items}}`）がデータを受ける Excel のブックで、Excel が印刷するページとして組まれます。結果はメモリの中の `.bdf` です。
- **キューに置く。** アプリケーションは自分で書いた設定（プリンタ 1 台と、両面を既定にしたキュー `invoices`）から [`server.New`](https://github.com/shibukawa/bdf/tree/main/print/server) でプリントサーバーを開き、ページが名乗る利用者のジョブとして `Spool().Submit` で文書を置きます。プリンタは待ちません。
- **API を共有する。** プリントサーバーの `http.Handler` を `/api/` の下に置きます。ページは `/api/v1/jobs` を読み、`/api/v1/events` を追います。`bdf jobs -server http://127.0.0.1:8083` でも同じジョブを扱えます。
- **ジョブを残す。** スプールは `data/spool.db` です。送っている最中にプログラムを止めて動かし直すと、ジョブはプリンタと照らし合わされ、1 回だけ刷られます。

## なぜこの形か

印刷する業務システムは、プリンタが刷り終えるまで要求を開いたまま待つべきではなく、プリンタの電源が切れていたからといってジョブを失うべきでもありません。プリントサーバーを組み込めば、別のプロセスなしにスプールの約束（1 つのファイル、送り直し、二重に刷らない、実績）が手に入ります。大きな現場では [`bdf printd`](../print-server.ja.md) を単独で動かし、[`print/server/client`](https://github.com/shibukawa/bdf/tree/main/print/server/client) や、`curl` のように素の HTTP で話します。

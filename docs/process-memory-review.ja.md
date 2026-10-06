# AC電源での比較ベンチマーク

| 入力 | RSS: LibreOffice / Go (MiB) | 経過時間: LibreOffice / Go (秒) | CPU時間: LibreOffice / Go (秒) |
|---|---:|---:|---:|
| 小さい DOCX | 161.8 / 70.0 | 0.69 / 0.06 | 0.67 / 0.06 |
| 小さい PPTX | 147.3 / 53.1 | 1.01 / 0.05 | 1.02 / 0.05 |
| 小さい XLSX | 123.7 / 52.1 | 0.65 / 0.05 | 0.63 / 0.05 |
| 合成 XLSX（4シート×3,000行） | 288.3 / 53.0 | 1.55 / 0.41 | 1.62 / 0.49 |
| 合成 PPTX（100スライド） | 1918.9 / 42.6 | 56.38 / 0.18 | 68.30 / 0.24 |

推定SoC電力量は LibreOffice **4.267 J/文書**（ラウンド範囲 3.295〜7.155）、Go **0.274 J/文書**（0.272〜0.277）。両実装を同じAC電源で3ラウンド測定しました。

macOSのpeak memory footprintも記録した。RSSとfootprintは異なる指標で、各列の中央値が同じ実行回とは限らない。

| 入力 | footprint: LibreOffice / Go (MiB) |
|---|---:|
| 小さい DOCX | 93.2 / 56.9 |
| 小さい PPTX | 81.0 / 40.0 |
| 小さい XLSX | 53.1 / 39.7 |
| 合成 XLSX（4シート×3,000行） | 210.6 / 39.4 |
| 合成 PPTX（100スライド） | 4114.4 / 29.0 |

時間・メモリは各5回の中央値。GoもLibreOffice / Popplerも同じAC電源で計測する。ビルドとウォームアップは測定に含めず、実行順を交互に入れ替える。新規プロセスで文書を変換し、256 px fit の先頭ページサムネイルを保存する。macOS `/usr/bin/time -l` の最大RSS・peak memory footprint・経過時間・user + systemのCPU時間を使う。LibreOffice / Popplerは逐次実行するため、ピークは2段階の最大値、時間は合計。経過時間の記録は0.01秒単位で、CPU時間は複数スレッドの合計なので経過時間を上回る場合がある。

電力量は1秒間隔のpowermetricsでCPU・GPU・ANEの推定電力を採取する。デフォルトでは各実装を60秒間連続して変換し、入力を一巡する単位で終了して文書構成を揃える。開始待ちと終了用に計12秒の採取余裕を設け、同じ長さのアイドル測定と組み合わせる。各サンプルの実際の採取時間に電力を掛けて積算し、アイドル平均×作業フェーズ採取時間を差し引き、完成した文書数で割る。3ラウンドで実行順を交互にし、中央値と最小・最大値を報告する。作業を長くすると背景負荷の影響を相対的に減らせるが、誤差がなくなるわけではない。端末全体やコンセントの電力量は測っていない。

AC電源に接続し、満充電になって充電負荷が落ち着いてから、画面の明るさと実行中のアプリを揃えて測定する。測定中はビルド・ブラウザ操作などほかの作業を避ける。

```sh
sudo -v
python3 tools/benchmark-all.py --large-input-dir /tmp/bdf-memory-fixtures
```

大きい入力を含める場合は、事前に `tools/generate-memory-fixtures.py --output /tmp/bdf-memory-fixtures` で生成する（xlsxwriter 3.2.9 / python-pptx 1.0.2）。`--large-input-dir` を省略すると小さいDOCX / PPTX / XLSXの3文書を測る。電力量の文書構成はこの小さい3文書に固定し、両実装で同じ入力・同じGoバイナリを使う。

すでにメモリ・CPUの結果がある場合は、その結果を再利用して電力量だけを測れる。

```sh
sudo -v
python3 tools/benchmark-all.py --memory-results /tmp/bdf-comparison-ac/summary.json
```

`benchmark-all.py` は通常ユーザーで動かす。sudoを使うのは内部のpowermetricsだけ。実行完了時に日英の説明と図を更新し、[latest.json](benchmarks/latest.json)に最新の要約だけを保存する。バイナリ・変換出力・生ログ・実行ごとのJSONは `/tmp` の結果ディレクトリに置き、Gitには残さない。

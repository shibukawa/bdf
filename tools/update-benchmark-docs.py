#!/usr/bin/env python3
"""Publish only the latest AC comparison summary; raw logs stay outside Git."""
import argparse
import json
from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
LABELS = {
    'basic.docx': ('小さい DOCX', 'Small DOCX'),
    'basic.pptx': ('小さい PPTX', 'Small PPTX'),
    'basic.xlsx': ('小さい XLSX', 'Small XLSX'),
    'many-sheets.xlsx': ('合成 XLSX（4シート×3,000行）', 'Synthetic XLSX (4 sheets × 3,000 rows)'),
    'many-slides.pptx': ('合成 PPTX（100スライド）', 'Synthetic PPTX (100 slides)'),
}


def compact(memory, power):
    if not all('AC Power' in memory[k] for k in ('power_source_start', 'power_source_end')):
        raise ValueError('memory comparison must start and end on AC power')
    process = {k: memory[k] for k in ('created_utc', 'host', 'rounds', 'versions', 'bdf_binary', 'power_source_start', 'power_source_end')}
    process['inputs'] = [{**{k: p[k] for k in ('bytes', 'sha256')}, 'name': Path(p['path']).name} for p in memory['inputs']]
    process['results'] = {
        name: {w: {**peaks[w], 'elapsed_seconds': memory['median_pipeline_elapsed_seconds'][name][w],
                   'cpu_seconds': memory['median_pipeline_cpu_seconds'][name][w]}
               for w in ('libreoffice', 'bdf')}
        for name, peaks in memory['median_pipeline_peak_bytes'].items()
    }
    data = {'process': process, 'power': None}
    if power is not None:
        if power['bdf_binary']['sha256'] != memory['bdf_binary']['sha256']:
            raise ValueError('use the same bdf binary for process and power measurements')
        if set(power['workflows']) != {'bdf', 'libreoffice'}:
            raise ValueError('power comparison must include both workflows')
        if not all('AC Power' in power[k] for k in ('power_source_start', 'power_source_end')):
            raise ValueError('power comparison must start and end on AC power')
        expected = {p['sha256'] for p in memory['inputs']}
        if not all(p['sha256'] in expected for p in power['inputs']):
            raise ValueError('power inputs must match process comparison inputs')
        data['power'] = {k: power[k] for k in ('created_utc', 'host', 'rounds', 'method', 'work_seconds_target',
                                              'power_source_start', 'power_source_end', 'bdf_binary', 'pairs',
                                              'median_estimated_soc_j_per_document', 'range_estimated_soc_j_per_document')}
        data['power']['inputs'] = [Path(p['path']).name for p in power['inputs']]
        data['power']['phases'] = [{k: phase[k] for k in ('name', 'samples', 'mean_soc_w', 'sampled_seconds', 'sampled_soc_j', 'work_seconds') if k in phase}
                                 for phase in power['phases']]
    return data


def table(data, ja):
    rows = ['| 入力 | RSS: LibreOffice / Go (MiB) | 経過時間: LibreOffice / Go (秒) | CPU時間: LibreOffice / Go (秒) |' if ja else
            '| Input | RSS: LibreOffice / Go (MiB) | Elapsed: LibreOffice / Go (s) | CPU: LibreOffice / Go (s) |',
            '|---|---:|---:|---:|']
    for name, r in data['process']['results'].items():
        label = LABELS.get(name, (name, name))[0 if ja else 1]
        lo, go = r['libreoffice'], r['bdf']
        rows.append(f"| {label} | {lo['max_rss_bytes']/2**20:.1f} / {go['max_rss_bytes']/2**20:.1f} | {lo['elapsed_seconds']:.2f} / {go['elapsed_seconds']:.2f} | {lo['cpu_seconds']:.2f} / {go['cpu_seconds']:.2f} |")
    return '\n'.join(rows)


def energy_text(data, ja):
    p = data['power']
    if p is None:
        return ('電力量は両実装をAC電源で再計測する準備ができています。以前の短時間・異なる電源条件の測定値は比較に使用していません。' if ja else
                'A new AC-powered energy comparison of both workflows is pending. Earlier short workloads and measurements under different power conditions are not used here.')
    med = p['median_estimated_soc_j_per_document']; ranges = p['range_estimated_soc_j_per_document']
    if ja:
        return (f"推定SoC電力量は LibreOffice **{med['libreoffice']:.3f} J/文書**（ラウンド範囲 {ranges['libreoffice'][0]:.3f}〜{ranges['libreoffice'][1]:.3f}）、Go **{med['bdf']:.3f} J/文書**（{ranges['bdf'][0]:.3f}〜{ranges['bdf'][1]:.3f}）。両実装を同じAC電源で{p['rounds']}ラウンド測定しました。")
    return (f"Estimated SoC energy: LibreOffice **{med['libreoffice']:.3f} J/document** (round range {ranges['libreoffice'][0]:.3f}–{ranges['libreoffice'][1]:.3f}), Go **{med['bdf']:.3f} J/document** ({ranges['bdf'][0]:.3f}–{ranges['bdf'][1]:.3f}). Both workflows used AC power for {p['rounds']} rounds.")


def chart(data, ja):
    r = data['process']['results'].get('basic.pptx', next(iter(data['process']['results'].values())))
    title = 'AC電源での比較ベンチ' if ja else 'Comparison benchmark on AC power'
    cards = [('ピークRSS' if ja else 'Peak RSS', 'MiB', [(w, r[w]['max_rss_bytes']/2**20) for w in ('libreoffice', 'bdf')]),
             ('経過時間' if ja else 'Elapsed time', 's', [(w, r[w]['elapsed_seconds']) for w in ('libreoffice', 'bdf')]),
             ('CPU時間 (user + system)' if ja else 'CPU time (user + system)', 's', [(w, r[w]['cpu_seconds']) for w in ('libreoffice', 'bdf')]),
             ('推定SoC電力量' if ja else 'Estimated SoC energy', 'J/doc', [(w, data['power']['median_estimated_soc_j_per_document'][w]) for w in ('libreoffice', 'bdf')] if data['power'] else [])]
    parts = ['<svg xmlns="http://www.w3.org/2000/svg" width="760" height="368" viewBox="0 0 760 368" role="img" aria-labelledby="title description">',
             f'<title id="title">{escape(title)}</title>', '<desc id="description">'+escape('; '.join(h+': '+', '.join(f'{w} {v:.3f} {unit}' for w,v in values) for h,unit,values in cards))+'</desc>', '<rect width="760" height="368" fill="#fff"/>', '<g font-family="-apple-system, BlinkMacSystemFont, Arial, sans-serif">']
    def text(x,y,s,size=12,color='#52627a',weight=400):
        parts.append(f'<text x="{x}" y="{y}" font-size="{size}" font-weight="{weight}" fill="{color}">{escape(s)}</text>')
    text(24,33,title,22,'#182337',700)
    text(24,57, 'Apple M3 · '+data['process']['created_utc'][:10]+' · '+('時間・メモリ: 小さいPPTX' if ja else 'Time / memory: small PPTX'),13)
    for i,(heading,unit,rows) in enumerate(cards):
        x=24+(i%2)*366; y=76+(i//2)*128
        parts.append(f'<rect x="{x}" y="{y}" width="346" height="116" rx="8" fill="#f7f9fc" stroke="#e5eaf0"/>')
        text(x+12,y+25,heading,15,'#182337',700)
        if not rows:
            text(x+12,y+65,'両実装を再計測待ち' if ja else 'Awaiting both-workflow measurement',13)
            continue
        maximum=max(v for _,v in rows)
        for n,(w,v) in enumerate(rows):
            by=y+52+n*26
            text(x+12,by+12,'LO + Poppler' if w=='libreoffice' else 'Go',11)
            width=108*max(v,0)/maximum if maximum>0 else 0
            parts.append(f'<rect x="{x+123}" y="{by}" width="{width:.2f}" height="16" rx="3" fill="'+('#667892' if n==0 else '#087d66')+'"/>')
            text(x+242,by+12,f'{v:.1f}' if unit=='MiB' else f'{v:.3f}' if unit=='J/doc' else f'{v:.2f}',12,'#182337',600)
        text(x+12,y+110,unit+' · '+('中央値' if ja else 'median'),10)
    text(24,347,'電力量: CPU + GPU + ANE の推定値。ラウンド範囲は本文を参照。' if ja else 'Energy estimates CPU + GPU + ANE. Round ranges are shown in the text.',11)
    return '\n'.join(parts+['</g>','</svg>'])+'\n'


def publish(memory, power=None):
    data=compact(memory,power)
    rounds=data['process']['rounds']
    (ROOT/'docs/benchmarks/latest.json').write_text(json.dumps(data,ensure_ascii=False,indent=2)+'\n')
    for ja in (True,False):
        path=ROOT/'docs'/('why.ja.md' if ja else 'why.md')
        old=path.read_text()
        start=old.index('## メモリ・CPU・消費電力の実測' if ja else '## Measured memory, CPU time, and energy')
        end=old.index('## オフィススイートを動かさなくてよい' if ja else '## No office suite to run',start)
        intro=('## メモリ・CPU・消費電力の実測\n\n同じApple M3のMacで、GoとLibreOffice → PDF → PopplerをAC電源で再計測しました。Officeファイルの変換と256 pxの先頭ページサムネイル生成を比較しています。\n\n' if ja else
               '## Measured memory, CPU time, and energy\n\nGo and LibreOffice → PDF → Poppler were remeasured on the same Apple M3 Mac on AC power. Each workflow converts an Office file and creates a 256 px first-page thumbnail.\n\n')
        lang='ja' if ja else 'en'
        alt='AC電源での最新のメモリ・CPU・電力量比較' if ja else 'Latest AC memory, CPU, and energy comparison'
        section=intro+f'[![{alt}](images/why-economy.{lang}.svg)](images/why-economy.{lang}.svg)\n\n'+table(data,ja)+'\n\n'
        section+=(f'時間・メモリはウォームアップ後の各{rounds}回の中央値。RSSは各プロセスの最大常駐メモリ、CPU時間はuser + system、経過時間は起動・変換・保存を含む実時間です。LibreOfficeとPopplerは順番に動くため、ピークメモリは2段階の大きい方、時間は合計です。出力形式と表示結果は異なります。\n\n' if ja else
                  f'Time and memory are medians of {rounds} runs after warmup. RSS is peak resident memory; CPU time is user + system; elapsed time includes startup, conversion, and output. LibreOffice and Poppler run sequentially, so their pipeline peak is the larger stage peak and their times are added. Output formats and visual results differ.\n\n')
        section+=energy_text(data,ja)+'\n\n'
        section+=('[最新結果と測定条件](benchmarks/latest.json)、[再計測手順](process-memory-review.ja.md)。生ログはGitに保存せず、ローカルの結果ディレクトリに置きます。\n\n' if ja else
                  'See the [latest results and conditions](benchmarks/latest.json) and [measurement procedure](process-memory-review.ja.md). Raw logs stay in the local results directory outside Git.\n\n')
        path.write_text(old[:start]+section+old[end:])
        (ROOT/'docs/images'/f'why-economy.{lang}.svg').write_text(chart(data,ja))
    report='# AC電源での比較ベンチマーク\n\n'+table(data,True)+'\n\n'+energy_text(data,True)+'\n\n'
    report+='macOSのpeak memory footprintも記録した。RSSとfootprintは異なる指標で、各列の中央値が同じ実行回とは限らない。\n\n| 入力 | footprint: LibreOffice / Go (MiB) |\n|---|---:|\n'
    for name,r in data['process']['results'].items():
        report+=f"| {LABELS.get(name,(name,name))[0]} | {r['libreoffice']['peak_footprint_bytes']/2**20:.1f} / {r['bdf']['peak_footprint_bytes']/2**20:.1f} |\n"
    report+='\n'
    report+=f'''時間・メモリは各{rounds}回の中央値。GoもLibreOffice / Popplerも同じAC電源で計測する。ビルドとウォームアップは測定に含めず、実行順を交互に入れ替える。新規プロセスで文書を変換し、256 px fit の先頭ページサムネイルを保存する。macOS `/usr/bin/time -l` の最大RSS・peak memory footprint・経過時間・user + systemのCPU時間を使う。LibreOffice / Popplerは逐次実行するため、ピークは2段階の最大値、時間は合計。経過時間の記録は0.01秒単位で、CPU時間は複数スレッドの合計なので経過時間を上回る場合がある。

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
'''
    (ROOT/'docs/process-memory-review.ja.md').write_text(report)
    return data


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--memory',type=Path,required=True)
    p.add_argument('--power',type=Path)
    a=p.parse_args()
    publish(json.loads(a.memory.read_text()), json.loads(a.power.read_text()) if a.power else None)

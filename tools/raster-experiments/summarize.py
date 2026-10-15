#!/usr/bin/env python3
from pathlib import Path
from collections import defaultdict
import argparse, json, re, shutil, statistics
ap=argparse.ArgumentParser();ap.add_argument('--input',type=Path,default=Path('/tmp/bdf-raster-experiments'));ap.add_argument('--output',type=Path,default=Path('knowledge-inputs/raster-acceleration/experiments/go'));args=ap.parse_args();args.output.mkdir(parents=True,exist_ok=True)
variants={}
for work in sorted(args.input.iterdir()):
    if not work.is_dir() or not (work/'resource.json').is_file():continue
    target=args.output/work.name;target.mkdir(parents=True,exist_ok=True)
    for name in ['run.txt','quality.json','precision.json','resource.json','quality-resource.json','inputs.json']:
        if (work/name).is_file():shutil.copyfile(work/name,target/name)
    groups=defaultdict(list)
    for line in (work/'run.txt').read_text().splitlines():
        m=re.match(r'(Benchmark\S+)\s+\d+\s+(\d+) ns/op.*?\s+(\d+) B/op\s+(\d+) allocs/op',line)
        if m:groups[m[1]].append(tuple(map(int,m.groups()[1:])))
    medians={name:{'ns_op':statistics.median(x[0] for x in rows),'bytes_op':statistics.median(x[1] for x in rows),'allocs_op':statistics.median(x[2] for x in rows),'samples':len(rows),'min_ns':min(x[0] for x in rows),'max_ns':max(x[0] for x in rows)} for name,rows in groups.items()}
    quality=json.loads((work/'quality.json').read_text())
    variants[work.name]={'benchmarks':medians,'quality':quality,'resource':json.loads((work/'resource.json').read_text())}
base=variants['baseline']['benchmarks']
for name,data in variants.items():
    for case,row in data['benchmarks'].items():row['speedup_vs_baseline']=base[case]['ns_op']/row['ns_op'] if case in base else None
result={'method':'same Go overlay instrumentation, same fixture/AA/DPI per pair; warm rendering; short samples; separate process resource data; convex uses 5 x 500ms, other variants 3 x 200ms','limitations':['ordered sequential suite, baseline drift possible; use bracket results for text-run conclusions','microbenchmarks do not prove page speedups','RSS covers each benchmark process, not one isolated page','prototype cache byte counts are payload limits, not complete map/key accounting'],'variants':variants}
(args.output/'summary.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
for name,data in variants.items():
    rows=[r for k,r in data['benchmarks'].items() if '/300dpi/on' in k]
    changes=sum(x['changed_pixels'] for x in data['quality'])
    err=max(x['max_channel_error'] for x in data['quality'])
    print(name, '300dpi AA on speedups', [round(x['speedup_vs_baseline'],2) for x in rows], 'changed pixels',changes,'max error',err,'RSS MiB',round(data['resource']['max_rss']/1048576,1))

#!/usr/bin/env python3
from pathlib import Path
import argparse, json, subprocess, sys
here=Path(__file__).resolve().parent
ap=argparse.ArgumentParser();ap.add_argument('--output',type=Path,default=Path('/tmp/bdf-raster-experiments'));ap.add_argument('--variants',nargs='*',default=['baseline','solid-fused','opaque-fill','flatten-cache','text-run-cache','sparse-runs','small-a8','rgba8','simd','stroke-cache','adaptive','hybrid-storage','convex']);args=ap.parse_args()
args.output=args.output.resolve();args.output.mkdir(parents=True,exist_ok=True)
results={}
for variant in args.variants:
    work=args.output/variant
    print('Running',variant,flush=True)
    command=[sys.executable,str(here/'run.py'),'--variant',variant,'--output',str(work),'--benchtime','200ms','--count','3']
    if variant!='baseline':command+=['--reference',str(args.output/'baseline')]
    p=subprocess.run(command,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
    work.mkdir(parents=True,exist_ok=True);(work/'driver.txt').write_text(p.stdout)
    results[variant]={'exit_code':p.returncode};(args.output/'batch-status.json').write_text(json.dumps(results,indent=2)+'\n')
    print(variant,'exit',p.returncode,flush=True)
    if p.returncode:print(p.stdout[-2500:],flush=True)

raise SystemExit(1 if any(row["exit_code"] for row in results.values()) else 0)

#!/usr/bin/env python3
"""Replay the exact saved pre-adoption overlay on the current checkout."""
from pathlib import Path
import argparse,json,os,subprocess
root=Path(__file__).resolve().parents[2];base=root/'knowledge-inputs/raster-acceleration/experiments'
ap=argparse.ArgumentParser();ap.add_argument('--trial',default='go/baseline');ap.add_argument('--output',type=Path,default=Path('/tmp/bdf-raster-replay'));ap.add_argument('--bench',default='BenchmarkRenderAntialias');ap.add_argument('--benchtime',default='200ms');ap.add_argument('--count',type=int,default=3);args=ap.parse_args()
manifest=json.loads((base/'snapshots.json').read_text());trial=manifest['trials'][args.trial];work=args.output.resolve();work.mkdir(parents=True,exist_ok=True)
mapping={str(root/'raster/imagebdf'/name):str(base/path)for name,path in trial['files'].items()};mapping[str(root/'raster/imagebdf/object.go')]=str(base/manifest['baseline_object'])
mapping[str(root/'raster/imagebdf/antialias_test.go')]=str(base/manifest['baseline_antialias_test'])
empty=work/'empty.go.txt';empty.write_text('package imagebdf\n')
for name in ['cache.go','pages.go','cache_test.go','pages_test.go','acceleration_test.go']:mapping[str(root/'raster/imagebdf'/name)]=str(empty)
overlay=work/'overlay.json';overlay.write_text(json.dumps({'Replace':mapping},indent=2)+'\n');env=dict(os.environ,GOCACHE='/tmp/bdf-aa-go-cache',BDF_RASTER_EXPERIMENT=trial['variant']);binary=work/'raster.test'
subprocess.run(['go','test','-c','-overlay',str(overlay),'-o',str(binary),'./raster/imagebdf'],cwd=root,env=env,check=True)
subprocess.run([str(binary),'-test.run=^$','-test.bench='+args.bench,'-test.benchtime='+args.benchtime,'-test.count='+str(args.count),'-test.cpu=1'],cwd=root/'raster/imagebdf',env=env,check=True)

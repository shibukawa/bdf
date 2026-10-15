#!/usr/bin/env python3
"""Bracket the repeated-text candidate with baseline runs to expose drift."""
import json,os,subprocess,time
from pathlib import Path
root=Path(__file__).resolve().parents[2];out=Path('/tmp/bdf-raster-bracket');out.mkdir(exist_ok=True)
for label,variant in [('before','baseline'),('text','text-run-cache'),('after','baseline')]:
 folder=out/label;folder.mkdir(exist_ok=True)
 if label!='after':
  subprocess.run(['python3',str(root/'tools/raster-experiments/run.py'),'--variant',variant,'--output',str(folder),'--reference','/tmp/bdf-raster-final/baseline','--test-only'],cwd=root,check=True)
 binary=(out/'before'if label=='after'else folder)/'raster.test';env=dict(os.environ,BDF_RASTER_EXPERIMENT=variant)
 cmd=[str(binary),'-test.run=^$','-test.bench=^BenchmarkRenderAntialias/(pdf|pptx|drawio)/.*/300dpi/on$','-test.benchtime=500ms','-test.count=7','-test.cpu=1']
 print(label,flush=True)
 with (folder/'run.txt').open('w')as stream:
  t=time.perf_counter();p=subprocess.Popen(cmd,cwd=root/'raster/imagebdf',env=env,stdout=stream,stderr=subprocess.STDOUT);_,status,u=os.wait4(p.pid,0);p.returncode=os.waitstatus_to_exitcode(status)
 (folder/'resource.json').write_text(json.dumps({'exit_code':p.returncode,'wall_seconds':time.perf_counter()-t,'user_seconds':u.ru_utime,'system_seconds':u.ru_stime,'max_rss_bytes':u.ru_maxrss},indent=2)+'\n')
 if p.returncode:raise RuntimeError(f"{label} failed: {p.returncode}")

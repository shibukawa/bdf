#!/usr/bin/env python3
import json,os,subprocess,time
from pathlib import Path
root=Path(__file__).resolve().parents[2];out=Path('/tmp/bdf-native-final');scenes=Path('/tmp/bdf-native-scenes')
def run(command,folder):
 folder.mkdir(parents=True,exist_ok=True);start=time.perf_counter()
 with (folder/'driver.txt').open('w')as stream:
  p=subprocess.Popen(command,cwd=root,stdout=stream,stderr=subprocess.STDOUT);_,status,usage=os.wait4(p.pid,0);p.returncode=os.waitstatus_to_exitcode(status)
 (folder/'resource.json').write_text(json.dumps({'wall_seconds':time.perf_counter()-start,'user_seconds':usage.ru_utime,'system_seconds':usage.ru_stime,'max_rss_bytes':usage.ru_maxrss,'exit_code':p.returncode,'scope':'whole experimental process, retains multiple renderer caches and measures all modes sequentially; not per-mode peak'},indent=2)+'\n')
 if p.returncode:raise RuntimeError((folder/'driver.txt').read_text())
for name in ['opaque-rectangles','cubic-paths','translucent-overlap','stroked-curves']:
 print('Refresh native public API',name,flush=True);folder=out/'imagebdf'/name
 run(['/tmp/bdf-pipeline-bench','-input',str(scenes/(name+'.bdf')),'-output',str(folder)],folder)
 for threads in [0,4]:
  folder=out/f'blend2d-{threads}'/name
  p=subprocess.run(['/tmp/bdf-native-blend2d-bench',str(scenes/(name+'.txt')),str(folder/'page.png'),str(threads)],capture_output=True,text=True,check=True);(folder/'result.json').write_text(p.stdout)
for name,path in [('pdf-chrome','testdata/pdf/chrome-doc.bdf'),('pptx-features','testdata/pptx/features.bdf'),('drawio-showcase','testdata/drawio/showcase.bdf'),('effects-boundary',str(scenes/'effects-boundary.bdf'))]:
 for dpi in [72,300]:
  print(name,dpi,flush=True);folder=out/'public-pipeline'/f'{name}-{dpi}'
  run(['/tmp/bdf-pipeline-bench','-input',path,'-output',str(folder),'-dpi',str(dpi),'-samples','3','-sample-ms','200'],folder)
print('Completed public pipelines',flush=True)

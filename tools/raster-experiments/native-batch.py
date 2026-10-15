#!/usr/bin/env python3
"""Run native and public-API prototypes sequentially, avoiding CPU contention."""
import argparse, hashlib, json, os, statistics, subprocess, time
from pathlib import Path
ap=argparse.ArgumentParser();ap.add_argument('--output',type=Path,default=Path('/tmp/bdf-native-results'));ap.add_argument('--scenes',type=Path,default=Path('/tmp/bdf-native-scenes'));args=ap.parse_args()
here=Path(__file__).resolve().parent
def run(command,folder):
    folder.mkdir(parents=True,exist_ok=True);start=time.perf_counter();p=subprocess.run(command,capture_output=True,text=True)
    (folder/'stdout.txt').write_text(p.stdout);(folder/'stderr.txt').write_text(p.stderr)
    (folder/'execution.json').write_text(json.dumps({'command':command,'exit_code':p.returncode,'wall_seconds':time.perf_counter()-start},indent=2)+'\n')
    if p.returncode:raise RuntimeError(p.stderr)
    if p.stdout.strip().startswith('{'):(folder/'result.json').write_text(p.stdout)
for name in ['opaque-rectangles','cubic-paths','translucent-overlap','stroked-curves']:
    scene=args.scenes/(name+'.json');print('Scene',name,flush=True)
    run(['/tmp/bdf-pipeline-bench','-input',str(scene.with_suffix('.bdf')),'-output',str(args.output/'imagebdf'/name)],args.output/'imagebdf'/name)
    run(['/tmp/bdf-raster-native-venv/bin/python',str(here/'native-skia.py'),str(scene),str(args.output/'skia'/name)],args.output/'skia'/name)
    for threads in [0,4]:
        folder=args.output/f'blend2d-{threads}'/name
        run(['/tmp/bdf-native-blend2d-bench',str(scene.with_suffix('.txt')),str(folder/'page.png'),str(threads)],folder)
    for threads,mode,cached in [(0,'u8','fresh'),(0,'f32','fresh'),(4,'u8','fresh'),(0,'u8','cached')]:
        folder=args.output/f'vello-cpu-{threads}-{mode}-{cached}'/name
        run(['/tmp/bdf-raster-native-target/release/bdf-raster-native-bench',str(scene),str(folder/'page.png'),str(threads),mode,cached],folder)
    for kind in ['compute','hybrid']:
        folder=args.output/('vello-'+kind)/name
        run(['/tmp/bdf-raster-native-target/release/bdf-raster-gpu-bench',str(scene),str(folder/'page.png'),kind],folder)
    folder=args.output/'imagebdf'/name;pdf=folder/'vector.pdf';samples=[]
    for _ in range(5):
        start=time.perf_counter_ns();p=subprocess.run(['gs','-q','-dSAFER','-dBATCH','-dNOPAUSE','-sDEVICE=png16m','-dGraphicsAlphaBits=4','-r72',f'-sOutputFile={folder/"pdf-raster.png"}',str(pdf)],capture_output=True,text=True);p.check_returncode();samples.append(time.perf_counter_ns()-start)
    (folder/'pdf-downstream.json').write_text(json.dumps({'samples_ns':samples,'median_ns':statistics.median(samples),'includes':'Ghostscript process startup, PDF decoding/rasterization and PNG encoding; not comparable to raw image Page alone'},indent=2)+'\n')
print('Finished',flush=True)

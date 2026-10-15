#!/usr/bin/env python3
"""Export bounded numeric native results and pixel errors; omit generated PNGs."""
from pathlib import Path
import argparse,hashlib,json,statistics,subprocess
import numpy as np
from PIL import Image
ap=argparse.ArgumentParser();ap.add_argument('--input',type=Path,default=Path('/tmp/bdf-native-final'));ap.add_argument('--output',type=Path,default=Path('knowledge-inputs/raster-acceleration/experiments/native'));args=ap.parse_args();args.output.mkdir(parents=True,exist_ok=True)
allrows=[]
for p in sorted(args.input.glob('*/*/result.json')):
 d=json.loads(p.read_text());backend,name=p.parent.parent.name,p.parent.name
 baseline=args.input/'imagebdf'/name/'page.png';images=list(p.parent.glob('*.png'));image=p.parent/'page.png'
 if not image.exists():image=p.parent/(name+'.png')
 if baseline.exists() and image.exists():
  a=np.asarray(Image.open(baseline).convert('RGBA')).astype(np.int16);b=np.asarray(Image.open(image).convert('RGBA')).astype(np.int16)
  delta=np.abs(a-b);d['pixel_comparison']={'reference':'imagebdf','changed_pixels':int(np.count_nonzero(np.any(delta!=0,axis=2))),'max_channel_error':int(delta.max()),'mean_channel_error':float(delta.mean()),'note':'AA/flattening/rounding differ by renderer; difference alone is not a quality verdict'}
 for k,v in d.get('results',{}).items():v['median_ns']=statistics.median(v['samples_ns'])
 if backend in ['imagebdf','public-pipeline']:
  d['input']='generated:'+name if backend=='imagebdf' or name.startswith('effects-') else 'fixture:'+name
  extra=p.parent/'pdf-downstream.json';d['downstream_pdf']=json.loads(extra.read_text()) if extra.exists() else None
  resource=p.parent/'resource.json'
  if resource.exists():d['process_resource']=json.loads(resource.read_text())
 # GPU harness currently creates both renderer implementations at startup.
 if backend in ['vello-compute','vello-hybrid']:d['init_both_backends_and_scene_ns']=d.pop('init_and_scene_ns');d['initialization_note']='includes both GPU renderer instances; not backend-specific startup cost'
 target=args.output/backend;target.mkdir(exist_ok=True);(target/(name+'.json')).write_text(json.dumps(d,ensure_ascii=False,indent=2)+'\n');allrows.append({'backend':backend,'scene':name,'result':d})
sources=['native-skia.py','native-blend2d.cpp','native-rust/src/main.rs','native-rust/src/gpu.rs','native-rust/Cargo.toml.in','cmd/scene/main.go','cmd/pipeline/main.go','cmd/effects/main.go']
here=Path(__file__).parent
meta={'machine':'Apple M3 / darwin arm64','skia_python':'144.0.post2','blend2d_commit':'58ca9460b4138af6e793184682e9324e0a8053b2','asmjit_commit':'c8d13c48fcd943531d1e26b140ca058dfbc6c0a1','vello_commit':'db90a73ec7309dfa2bd45886b353162a9ea95744','ghostscript':'10.07.1','sources':{str(p):hashlib.sha256((here/p).read_bytes()).hexdigest()for p in sources},'method':'sequential runs; 5 x 300ms per warm mode; prebuilt geometry or recorded scenes; owned output separated; no full BDF native integration','limitations':['AA algorithms differ','GPU warm measurement uses static preprocessed scenes; dynamic scene conversion excluded','native peak RSS not measured','single machine; not universal backend ranking'],'rows':allrows}
(args.output/'summary.json').write_text(json.dumps(meta,ensure_ascii=False,indent=2)+'\n')

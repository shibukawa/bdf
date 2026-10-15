#!/usr/bin/env python3
"""Refresh only Vello runs after changing its scene adapter."""
from pathlib import Path
import subprocess
root=Path('/tmp/bdf-native-final');scenes=Path('/tmp/bdf-native-scenes')
for name in ['opaque-rectangles','cubic-paths','translucent-overlap','stroked-curves']:
 print(name,flush=True)
 for threads,mode,cached in [(0,'u8','fresh'),(0,'f32','fresh'),(4,'u8','fresh'),(0,'u8','cached')]:
  folder=root/f'vello-cpu-{threads}-{mode}-{cached}'/name
  p=subprocess.run(['/tmp/bdf-raster-native-target/release/bdf-raster-native-bench',str(scenes/(name+'.json')),str(folder/'page.png'),str(threads),mode,cached],capture_output=True,text=True,check=True);(folder/'result.json').write_text(p.stdout)
 for kind in ['compute','hybrid']:
  folder=root/('vello-'+kind)/name
  p=subprocess.run(['/tmp/bdf-raster-native-target/release/bdf-raster-gpu-bench',str(scenes/(name+'.json')),str(folder/'page.png'),kind],capture_output=True,text=True,check=True);(folder/'result.json').write_text(p.stdout)

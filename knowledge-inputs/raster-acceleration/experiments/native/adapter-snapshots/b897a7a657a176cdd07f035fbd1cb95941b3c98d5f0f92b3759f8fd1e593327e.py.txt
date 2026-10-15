#!/usr/bin/env python3
"""Common path workloads through Skia CPU. Conversion is outside warm timing."""
import argparse,json,statistics,time
from pathlib import Path
import skia
ap=argparse.ArgumentParser();ap.add_argument('scene',type=Path);ap.add_argument('output',type=Path);args=ap.parse_args()
s=json.loads(args.scene.read_text());width,height=s['Width'],s['Height'];started=time.perf_counter_ns()
recorder=skia.PictureRecorder();record=recorder.beginRecording(width,height)
for op in s['operations']:
    path=skia.Path();path.setFillType(skia.PathFillType.kEvenOdd if op['rule']=='evenodd' else skia.PathFillType.kWinding)
    for seg in op['segments']:
        kind=int(seg[0]);v=seg[1:]
        if kind==0:path.moveTo(*v)
        elif kind==1:path.lineTo(*v)
        elif kind==2:path.quadTo(*v)
        elif kind==3:path.cubicTo(*v)
        elif kind==4:path.close()
        else:raise ValueError(kind)
    r,g,b,a=op['color'];paint=skia.Paint(Color=skia.ColorSetARGB(a,r,g,b),AntiAlias=True)
    if op['stroke']>0:
        paint.setStyle(skia.Paint.kStroke_Style);paint.setStrokeWidth(op['stroke']);paint.setStrokeMiter(10)
    record.drawPath(path,paint)
picture=recorder.finishRecordingAsPicture();record_ns=time.perf_counter_ns()-started
surface=skia.Surface(width,height);canvas=surface.getCanvas()
def render(owned):
    canvas.clear(skia.ColorWHITE);canvas.drawPicture(picture)
    if owned:return surface.makeImageSnapshot().tobytes()
started=time.perf_counter_ns();render(True);cold_ns=time.perf_counter_ns()-started
results={}
for owned in [False,True]:
    samples=[]
    for _ in range(5):
        start=time.perf_counter_ns();n=0
        while time.perf_counter_ns()-start<300_000_000:render(owned);n+=1
        samples.append((time.perf_counter_ns()-start)/n)
    results['owned-output' if owned else 'render-only']={'samples_ns':samples,'median_ns':statistics.median(samples)}
args.output.mkdir(parents=True,exist_ok=True);surface.makeImageSnapshot().save(str(args.output/(args.scene.stem+'.png')),skia.kPNG)
result={'backend':'skia-cpu','version':skia.__version__,'scene':args.scene.stem,'width':width,'height':height,'native_picture_record_ns':record_ns,'first_render_ns':cold_ns,'results':results,'scope':'SkPicture replay of prebuilt path geometry; JSON parsing/FFI recording excluded from warm; AA enabled'}
(args.output/(args.scene.stem+'.json')).write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result))

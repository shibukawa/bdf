#!/usr/bin/env python3
"""Run isolated raster prototypes using Go overlays; never patch production files."""
from pathlib import Path
import argparse, hashlib, json, os, subprocess, sys, time
ROOT=Path(__file__).resolve().parents[2]
HERE=Path(__file__).resolve().parent
VARIANTS=['baseline','solid-fused','opaque-fill','flatten-cache','text-run-cache','sparse-runs','small-a8','rgba8','simd','stroke-cache','adaptive','hybrid-storage','convex']
def replace(s,old,new,count=None):
    n=s.count(old)
    if n != (count if count is not None else 1):
        raise RuntimeError(f'Overlay anchor changed: expected {count or 1}, got {n}: {old[:100]!r}')
    return s.replace(old,new)
def prepare(variant,work):
    sources={name:(ROOT/'raster/imagebdf'/name).read_text() for name in ['imagebdf.go','scan.go','surface.go','draw.go','path.go']}
    s=sources['imagebdf.go']
    s=replace(s,'type Renderer struct {','type Renderer struct {\n experiment *experimentState')
    s=replace(s,'func (r *Renderer) newSurface(w, h int) *surface {','func (r *Renderer) newSurface(w, h int) *surface {\n if expMode == "rgba8" { return expNewSurface8(r,w,h) }')
    s=replace(s,'func (r *Renderer) release(s *surface) {','func (r *Renderer) release(s *surface) {\n if s != nil && s.packed != nil { expRelease8(r,s); return }')
    s=replace(s,'t.pix[4*(y*t.w+cov.r.Min.X):4*(y*t.w+cov.r.Max.X)]','expSurfaceRow(t,y,cov.r.Min.X,cov.r.Max.X)',2)
    s=replace(s,'if cov.a == nil {','if cov.a == nil && cov.encoded == nil {')
    s=replace(s,'s := r.newSurface(w, h)','s := experimentCanvas(r,w,h)')
    sources['imagebdf.go']=s
    s=sources['surface.go']
    s=replace(s,'type surface struct {','type surface struct {\n packed []byte')
    anchor='\tif unbounded(mode) || plain {'
    s=replace(s,anchor,'\tif experimentFill(s,cov,sh,alpha,mode,clip,r) { return }\n'+anchor)
    s=replace(s,'func (s *surface) fillBackground(cov *mask, c solid) {','func (s *surface) fillBackground(cov *mask, c solid) {\n if s.packed != nil { s.fill(cov,c,1,bdf.BlendSourceOver,&mask{r:s.bounds()}); return }')
    s=replace(s,'func (s *surface) toRGBA() *image.RGBA {','func (s *surface) toRGBA() *image.RGBA {\n if s.packed != nil { out:=image.NewRGBA(s.bounds()); copy(out.Pix,s.packed); return out }')
    s=replace(s,'copy(dst[:4*n], ss.s.pix[4*(sy*ss.s.w+sx):])','copy(dst[:4*n], expSurfaceRow(ss.s,sy,sx,sx+n))')
    s=replace(s,'copy(dst[4*i:4*i+4], ss.s.pix[4*(sy*ss.s.w+sx):])','copy(dst[4*i:4*i+4], expSurfaceRow(ss.s,sy,sx,sx+1))')
    s=replace(s,'plain || cov == nil || cov.a != nil','plain || cov == nil || cov.a != nil || cov.encoded != nil')
    sources['surface.go']=s
    s=sources['scan.go']
    s=replace(s,'type mask struct {','type mask struct {\n encoded *experimentCoverage')
    s=replace(s,'func (m *mask) bytes() int { return 4 * len(m.a) }','func (m *mask) bytes() int { if m.encoded != nil { return m.encoded.bytes() }; return 4*len(m.a) }')
    s=replace(s,'func (m *mask) at(x, y int) float32 {','func (m *mask) at(x, y int) float32 {\n if m.encoded != nil { return m.encoded.at(m.r,x,y) }')
    s=replace(s,'func (m *mask) row(y, x0, x1 int) []float32 {','func (m *mask) row(y, x0, x1 int) []float32 {\n if m.encoded != nil { return m.encoded.row(m.r,y,x0,x1) }')
    s=replace(s,'m := &mask{r: r, a: make([]float32, w*r.Dy())}','m := experimentMask(r)\n _ = w',2)
    s=replace(s,'row := m.a[(y-r.Min.Y)*w : (y-r.Min.Y+1)*w]','row := experimentWriteRow(m,y)',2)
    s=replace(s,'\t\tif next >= len(edges) && len(active) == 0 {','\t\texperimentFinishRow(m,y)\n\t\tif next >= len(edges) && len(active) == 0 {',2)
    s=replace(s,'a.a == nil && b.a == nil','a.a == nil && a.encoded == nil && b.a == nil && b.encoded == nil')
    s=replace(s,'\tif sc.noAA {','\texpIsConvex := expMode=="convex" && experimentConvex(polys)\n\tif sc.noAA {')
    # The AA scan has a span closure; binary scanning retains the existing route.
    assert s.count('\t\t\tsortCrossings(xs, active)')==2
    s=s.replace('\t\t\tsortCrossings(xs, active)','\t\t\tif expIsConvex && len(xs)==2 { span(math.Min(xs[0].x,xs[1].x),math.Max(xs[0].x,xs[1].x)); continue }\n\t\t\tsortCrossings(xs, active)',1)
    sources['scan.go']=s
    s=sources['path.go']
    s=replace(s,'type path struct {','type path struct {\n experimentVersion uint64')
    n=s.count('\tp.boxed = false'); s=replace(s,'\tp.boxed = false','\tp.boxed = false\n p.experimentVersion++',n)
    s=replace(s,'func (p *path) close() {','func (p *path) close() {\n p.experimentVersion++')
    s=replace(s,'func (p *path) flattenTo(sc *scratch, out []polyline, m matrix, tol float64) []polyline {','func (p *path) flattenTo(sc *scratch, out []polyline, m matrix, tol float64) []polyline {\n return experimentFlatten(p,sc,out,m,tol)\n}\nfunc (p *path) flattenUncached(sc *scratch, out []polyline, m matrix, tol float64) []polyline {')
    sources['path.go']=s
    if variant=='adaptive':
        s=replace(s,'func flattenQuad(out []point, x0, y0, x1, y1, x2, y2, tol float64) []point {','func flattenQuad(out []point, x0, y0, x1, y1, x2, y2, tol float64) []point {\n return experimentAdaptive(out,point{x0,y0},point{x0+(x1-x0)*2/3,y0+(y1-y0)*2/3},point{x2+(x1-x2)*2/3,y2+(y1-y2)*2/3},point{x2,y2},tol,0)\n}\nfunc unusedFlattenQuad(out []point, x0, y0, x1, y1, x2, y2, tol float64) []point {')
        s=replace(s,'func flattenCubic(out []point, x0, y0, x1, y1, x2, y2, x3, y3, tol float64) []point {','func flattenCubic(out []point, x0, y0, x1, y1, x2, y2, x3, y3, tol float64) []point {\n return experimentAdaptive(out,point{x0,y0},point{x1,y1},point{x2,y2},point{x3,y3},tol,0)\n}\nfunc unusedFlattenCubic(out []point, x0, y0, x1, y1, x2, y2, x3, y3, tol float64) []point {')
        sources['path.go']=s
    s=sources['draw.go']
    s=replace(s,'\tms := m.maxScale()\n\tsc.polys = p.flattenTo', '\tms := m.maxScale()\n if outline,ok:=experimentStrokeLookup(sc,p,m,ls,bounds);ok { return sc.rasterize(outline,bdf.NonZero,bounds),alpha,false }\n\tsc.polys = p.flattenTo')
    s=replace(s,'\treturn sc.rasterize(outline, bdf.NonZero, bounds), alpha, s.cut','\tif !s.cut { experimentStrokeStore(sc,p,m,ls,bounds,outline) }\n\treturn sc.rasterize(outline, bdf.NonZero, bounds), alpha, s.cut')
    anchor='\tsc, bounds := d.r.scratch(), d.bounds()\n\t// synthesized bold'
    s=replace(s,anchor,'\texpKey := experimentTextKey{ch:ch,text:s,m:m,size:st.size,spacing:st.spacing,dir:st.dir,bounds:d.bounds(),noAA:d.r.opts.Antialias==AntialiasOff}\n if cov:=experimentTextLookup(d.r,expKey); cov!=nil { d.paint(cov,d.shader(&st.fill,m),1); return }\n'+anchor)
    anchor='\td.paint(cov, d.shader(&st.fill, m), 1)\n}\n\n// strokeText'
    s=replace(s,anchor,'\texperimentTextStore(d.r,expKey,cov)\n'+anchor)
    s=replace(s,'copy(mc.target.pix[4*i:], c[:])','expStorePixel(mc.target,i,c)')
    s=replace(s,'p := ms.pix[4*i : 4*i+4]','p := expReadPixel(ms,i)')
    s=replace(s,'q := gs.pix[4*i : 4*i+4]','q := expReadPixel(gs,i)')
    # maskEnd mutates q; persist those changes into the packed store too.
    anchor='\t\tq[0], q[1], q[2], q[3] = q[0]*k, q[1]*k, q[2]*k, q[3]*k'
    s=replace(s,anchor,anchor+'\n expStorePixel(gs,i,q)')
    s=replace(s,'old.a != nil || cov.a != nil','old.a != nil || old.encoded != nil || cov.a != nil || cov.encoded != nil')
    s=replace(s,'cov.a == nil && r == old.r','cov.a == nil && cov.encoded == nil && r == old.r')
    s=replace(s,'old.a == nil && r == cov.r','old.a == nil && old.encoded == nil && r == cov.r')
    sources['draw.go']=s
    mapping={}; fingerprints={}
    for name,source in sources.items():
        p=work/name;p.write_text(source)
        mapping[str(ROOT/'raster/imagebdf'/name)]=str(p)
        fingerprints[name]=hashlib.sha256((ROOT/'raster/imagebdf'/name).read_bytes()).hexdigest()
    for name in ['experiment_support.go','experiment_support_test.go','experiment_simd_arm64.s','experiment_simd_arm64.go','experiment_simd_other.go']:
        snapshot=work/name
        snapshot.write_bytes((HERE/(name+'.txt')).read_bytes())
        mapping[str(ROOT/'raster/imagebdf'/name)]=str(snapshot)
    overlay=work/'overlay.json';overlay.write_text(json.dumps({'Replace':mapping},indent=2)+'\n')
    fixtures={str(p.relative_to(ROOT)):hashlib.sha256(p.read_bytes()).hexdigest() for pattern in ['testdata/pdf/chrome-doc.bdf','testdata/pptx/features.bdf','testdata/drawio/showcase.bdf','raster/imagebdf/antialias_test.go','raster/imagebdf/imagebdf_test.go','converter/pptx/testdata/fonts/*','converter/docx/testdata/fonts/*','fixture/testdata/fonts/*'] for p in ROOT.glob(pattern) if p.is_file()}
    prototype_hashes={name:hashlib.sha256((work/name).read_bytes()).hexdigest() for name in [*sources,'experiment_support.go','experiment_support_test.go','experiment_simd_arm64.s','experiment_simd_arm64.go','experiment_simd_other.go']}
    (work/'inputs.json').write_text(json.dumps({'variant':variant,'source_hashes':fingerprints,'fixture_hashes':fixtures,'prototype_hashes':prototype_hashes,'runner_hash':hashlib.sha256(Path(__file__).read_bytes()).hexdigest()},indent=2)+'\n')
    return overlay

def main():
    if (ROOT/'raster/imagebdf/cache.go').exists():
        raise SystemExit('This overlay generator targets the pre-adoption renderer. Replay saved trials with tools/raster-experiments/replay.py, or benchmark production with BenchmarkRenderCaches and BenchmarkPages.')
    ap=argparse.ArgumentParser();ap.add_argument('--variant',choices=VARIANTS,required=True)
    ap.add_argument('--output',type=Path,required=True);ap.add_argument('--reference',type=Path)
    ap.add_argument('--benchtime',default='200ms');ap.add_argument('--count',default=3,type=int);ap.add_argument('--test-only',action='store_true')
    args=ap.parse_args();args.output=args.output.resolve();args.output.mkdir(parents=True,exist_ok=True)
    overlay=prepare(args.variant,args.output)
    if args.test_only:
        for name in ['run.txt','resource.json']:(args.output/name).unlink(missing_ok=True)
    env=dict(os.environ,GOCACHE='/tmp/bdf-aa-go-cache',BDF_RASTER_EXPERIMENT=args.variant,BDF_EXPERIMENT_OUTPUT=str(args.output))
    if args.reference:env['BDF_EXPERIMENT_REFERENCE']=str(args.reference.resolve())
    binary=args.output/'raster.test'
    compile_command=['go','test','-c','-overlay',str(overlay),'-o',str(binary),'./raster/imagebdf']
    subprocess.run(compile_command,cwd=ROOT,env=env,check=True)
    run_pattern='TestRasterExperimentOutputs|TestExperimentSIMD'
    run_pattern+='|TestExperimentGlyphUnion|TestExperimentPrecision|TestExperimentGeometry'
    if args.variant not in ['small-a8','rgba8','adaptive','hybrid-storage']:run_pattern+='|TestAntialiasDefaultGolden'
    commands=[[str(binary),'-test.run='+run_pattern,'-test.count=1']]
    if not args.test_only:
        commands.append([str(binary),'-test.run=^$','-test.bench=^BenchmarkRenderAntialias$|^BenchmarkExperimentPixelKernel$|^BenchmarkExperimentRectangle$','-test.benchtime='+args.benchtime,'-test.count='+str(args.count),'-test.cpu=1'])
    (args.output/'command.json').write_text(json.dumps({'compile':compile_command,'run':commands},indent=2)+'\n')
    for index,command in enumerate(commands):
        with (args.output/('quality-run.txt' if index==0 else 'run.txt')).open('w') as output:
            started=time.perf_counter()
            proc=subprocess.Popen(command,cwd=ROOT/'raster/imagebdf',env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
            for line in proc.stdout:
                output.write(line);output.flush();sys.stdout.write(line);sys.stdout.flush()
            _,status,usage=os.wait4(proc.pid,0)
            code=os.waitstatus_to_exitcode(status);proc.returncode=code
            resources={'wall_seconds':time.perf_counter()-started,'user_seconds':usage.ru_utime,'system_seconds':usage.ru_stime,'max_rss':usage.ru_maxrss,'max_rss_unit':'bytes' if sys.platform=='darwin' else 'KiB','exit_code':code}
            (args.output/('quality-resource.json' if index==0 else 'resource.json')).write_text(json.dumps(resources,indent=2)+'\n')
        if code:return code
    return 0

if __name__=='__main__':raise SystemExit(main())

# Raster acceleration evidence

Manually checked on 2026-10-14. This is a bounded research note, not a complete audit or a benchmark of candidate backends. Repository observations, primary-source descriptions, and proposed applications are distinct. External links may change; unchanged local notes do not establish upstream freshness. Explicit-AA implementation and warm benchmarks were subsequently completed as recorded below. Non-AA candidates were subsequently evaluated in isolated prototypes and actual native CPU/GPU runs on 2026-10-14–15; see the bounded results below.

## Repository observations before implementation

- `raster/imagebdf/imagebdf.go`, `Options`: font lookup and background settings; no AA selector. `PageDPI` converts dpi to scale by dividing by 72. `Region` maps a document rectangle to independent pixel width and height. `Renderer` is not safe for concurrent use.
- `docs/spec.md`, manifest `unit`: BDF page units are points (1/72 inch).
- `raster/imagebdf/scan.go`, `subSamples`, `scratch.rasterize`: 16 Y samples per pixel row; X span coverage is fractional; the path bounding rectangle gets a fresh float32 coverage array. Active-edge order is reused across sample lines; work limits are explicit.
- `raster/imagebdf/scan.go`, `rectMask`: integer aligned rectangles use implicit full coverage; fractional boundaries allocate area coverage.
- `raster/imagebdf/surface.go`, `surface`, `surface.fill`: float32 premultiplied RGBA canvas (16 bytes/pixel), specialized source-over paths, final conversion to RGBA8.
- `raster/imagebdf/text.go`, `fontFace.outline`: original glyph outline cache. `raster/imagebdf/draw.go`, `fillText`: transforms/flattens glyph paths and rasterizes their coverage per drawing; no persistent glyph pixel-mask cache in this route.
- `raster/imagebdf/draw.go`, `strokeMask`, `strokeText`: subpixel-width strokes may expand to one pixel and reduce alpha. This preservation policy is separate from geometric AA.
- `raster/imagebdf/draw.go`, `shadowMask`, `maskEnd`; `raster/imagebdf/image.go`, `picture.levelFor`, `picture.sample`: blur, soft masks and image sampling have separate semantics.
- `raster/imagebdf/svg.go`: nested SVG includes direct `rasterize` and `rectMask` calls; quality propagation must cover them. `SVGImage` dimensions use CSS pixels, unlike BDF page points.
- `packages/viewer/src/index.ts`, render options: zoom multiplied by devicePixelRatio. A scale-only threshold can misclassify screen drawing.
- `raster/ebitenginebdf/ebitenginebdf.go`, package description: dedicated path-only backend; document text/images/gradients/clips/groups/masks are not supported as general document rendering.
- `docs/cpu-security-review.ja.md`, rasterizer section: prior Apple M3 single-CPU warm-render medians after buffer/sort/blend improvements: PDF 13.2 ms, PPTX 6.2 ms, draw.io 2.22 ms. These are historical repository measurements, not new measurements or speedups for the proposed techniques.

## Skia CPU geometry and storage

Primary sources:

- [Skia CPU Backend Architecture](https://skia.googlesource.com/skia/+/refs/heads/main/docs/architecture/CPU.md)
- [SkScan_AAAPath.cpp](https://github.com/google/skia/blob/main/src/core/SkScan_AAAPath.cpp)
- [SkScan_AntiPath.cpp](https://github.com/google/skia/blob/main/src/core/SkScan_AntiPath.cpp)

The analytic scan converter divides flattened geometry at integer Y boundaries, edge endpoints and intersections, then computes trapezoid coverage. Small regions use a fixed A8 mask; larger ones can use run-based rows. Convex and nonconvex paths take different accumulation routes. Simple rectangles can bypass small-mask preparation. The current general CPU AA fill entry calls AAAFillPath; old AAA-versus-SAA heuristic descriptions should not be presented as the current main implementation. This describes specific CPU paths, not every Skia operation.

## Skia pixel processing and glyphs

Primary sources:

- [SkRasterPipelineBlitter.cpp](https://github.com/google/skia/blob/main/src/core/SkRasterPipelineBlitter.cpp)
- [SkGlyphRunPainter.cpp](https://github.com/google/skia/blob/main/src/core/SkGlyphRunPainter.cpp)
- [SkGlyph.cpp](https://github.com/google/skia/blob/main/src/core/SkGlyph.cpp)
- [SkPaint reference](https://api.skia.org/classSkPaint.html)
- [SkFont reference](https://api.skia.org/classSkFont.html)

SkRasterPipeline combines precompiled stages and SIMD execution. Opaque SrcOver can reduce to Src; eligible constant Src fills use bulk memory writes. CPU text has a strike-backed mask drawing path, with generated glyph images reused. Shape AA and font edging have explicit controls. These facts do not establish a measured advantage on BDF documents.

## Vello sparse coverage and GPU reuse

Primary sources:

- [Vello Architecture](https://github.com/linebender/vello/blob/main/ARCHITECTURE.md)
- [Sparse strip path rendering](https://github.com/linebender/vello/issues/670)
- [Skia Graphite RasterPathAtlas](https://skia.googlesource.com/skia/+/refs/heads/main/src/gpu/graphite/RasterPathAtlas.h)
- [Skia Graphite ComputePathAtlas](https://skia.googlesource.com/skia/+/refs/heads/main/src/gpu/graphite/ComputePathAtlas.cpp)

Vello sparse strips keep detailed coverage near boundaries and compactly represent covered interiors. CPU SIMD/multithreaded and GPU consumers share preprocessing structures. Graphite also has a CPU-rasterized path-mask atlas cache and a conditional compute implementation using Vello shaders. GPU capability, configuration, readback and API stability must be evaluated for an actual backend; the presence of a source file is not proof that a default build selects it.

## Blend2D pipeline specialization

Primary sources:

- [Blend2D](https://blend2d.com/)
- [Blend2D presentation, Helsinki 2023](https://blend2d.com/research/blend2d_presentation_helsinki_2023.pdf)

Blend2D uses analytic AA, JIT-generated CPU-specific pixel pipelines and multithreaded rendering. Its presentation describes keeping coverage, paint fetch and compositing in registers. This differs from Skia RasterPipeline's precompiled stages. Library-authored performance comparisons are workload-specific and do not establish a universal fastest renderer.

## Resolution and print controls

Primary sources:

- [CSS Values and Units: Absolute Lengths](https://www.w3.org/TR/css-values-4/#absolute-lengths)
- [Using Ghostscript: TextAlphaBits and GraphicsAlphaBits](https://ghostscript.readthedocs.io/en/latest/Use.html#dtextalphabits-n)

CSS absolute units relate 1 inch to 96 CSS px and 72 pt; screen CSS inches need not match physical inches. Ghostscript exposes separate text and graphics AA controls and documents image interpolation separately. Neither source defines 100dpi as a universal threshold below which AA is necessary or above which it is unnecessary. Physical pitch is 25.4/dpi millimeters, but visible quality also depends on the final raster size, content, printer and scaling.

## Implemented explicit AA and measurements

`Options.Antialias` now supports AntialiasDefault, AntialiasOn and AntialiasOff. Default and On retain the prior AA route; Off uses one center-sampled scanline per row and integer X spans. Non-AA rectangles use implicit full coverage. Paths, strokes, text, geometric clips, page/region clips and embedded SVG image edges use the renderer policy. SVG raster cache keys include the resolved AA choice. Dense float32 masks remain for general paths; direct sparse-span compositing is still a candidate.

`raster/imagebdf/antialias_test.go` verifies public entry points, a pre-change AA pixel checksum, independent winding/center geometry, effects/opacity/image sampling, SVG caching and complexity limits. `cmd/bdf/preview_test.go` verifies CLI defaults and validation. Tests cover semantics and do not claim physical print quality.

`benchmark.txt` contains actual BenchmarkRenderAntialias output; `benchmark.json` contains median summaries. Apple M3, Go 1.27.0, darwin/arm64, one CPU, 200ms per benchmark, three runs, warm first-page rendering at 72 and 300dpi. Off/on speedups range from 1.17 to 2.13 overall, and 1.42 to 1.88 at 300dpi. Allocations are reported separately from peak memory. PNG encoding, cold starts and physical printing were not measured. One PDF 300dpi/off run was an outlier. These measurements concern this implementation and these fixtures, not a comparison with Skia or other backends.


## Non-AA bounded experiments

The completed evaluation is documented in [results-non-aa.ja.md](results-non-aa.ja.md), with a scope ledger in `experiments/coverage.yaml`. The production imagebdf implementation has not adopted these prototypes. AA stays explicitly controlled and independent of DPI, storage precision, image sampling and backend selection.

Go overlays evaluated run-mask caching, sparse row coverage/direct span compositing, small A8 masks, solid/opaque/convex specialization, ARM64 NEON, RGBA8 and mixed float layers, flattened geometry/stroke outline caching and adaptive subdivision. Three existing documents were rendered at 72/300 dpi with AA on/off. The corrected exact-storage prototypes matched their twelve reference outputs. Quantization and adaptive geometry changed pixels. Targeted regression checks covered compositing, shadows, limits and float behavior. NEON uses fused multiply-add to match this Go/arm64 compiler's scalar behavior; separate multiply/add initially exposed a low-level rounding mismatch.

Run-mask cache was bracketed by baseline runs: PDF-derived 300dpi AA-on rendering improved about 1.33x with baseline drift about 0.2%. PPTX and draw.io changes were smaller. A glyph-mask max-union counterexample demonstrates lost subpixel occupancy; this does not disprove all richer glyph-cache designs. RGBA8 had at most one channel step difference on the ordinary documents, but a 200-layer low-opacity counterexample reached 36 steps. Thus ordinary golden fixtures alone did not establish precision equivalence.

Public-API experiments evaluated Page/Region/Object image reuse and copy cost, independent-page batching and Region partitions. Independent page outputs matched hashes at 72/300 dpi, with race checks at 72 dpi. Fixed-halo partitions showed shadow and rounding differences; whole-page halos matched but duplicated work. Cache copy benchmarks model a static hit, excluding production lookup/eviction/invalidation. Cache payload limits do not include all map/key overhead.

Actual native runs used skia-python 144.0.post2, Blend2D commit 58ca9460b4138af6e793184682e9324e0a8053b2 plus AsmJit c8d13c48fcd943531d1e26b140ca058dfbc6c0a1, and Vello commit db90a73ec7309dfa2bd45886b353162a9ea95744. Four shared 1024x1024 path scenes exercised rectangles, cubic fills, translucent self-crossing fills with both fill rules, and butt-cap cubic strokes. Warm native pictures/paths/scenes were replayed; Vello CPU measured rebuilt and retained coarse geometry, u8/f32 and single/multiple threads. Owned image extraction was separate from render-only. Native pixel order and API boundaries are not normalized into a full Go image.RGBA BDF adapter.

Apple M3 Metal ran Vello compute and CPU-preprocessed sparse-strip hybrid renderers, including submission/upload and synchronization; owned output included readback and CPU copy. These were static preprocessed scenes, not dynamic end-to-end scene conversion. The GPU harness initializes both renderers, so its startup field is explicitly not a single-backend cold-start measurement. Native outputs differ in AA/flattening/rounding and are not pixel-compatible replacements. Full BDF font/image/clip/group/mask integration and native peak RSS were not evaluated.

Local pinned primary implementation observations:

- [Vello CPU current-state and pipeline documentation](https://github.com/linebender/vello/blob/db90a73ec7309dfa2bd45886b353162a9ea95744/vello_cpu/src/lib.rs): u8 and f32 pipelines, multithreading restrictions for filters, and experimental glyph caching.
- [Vello GPU current-state and architecture](https://github.com/linebender/vello/blob/db90a73ec7309dfa2bd45886b353162a9ea95744/vello_gpu/src/lib.rs): CPU preprocessing plus GPU raster/compositing, headless support, mask/filter/non-isolated-blend limitations, and rename from vello_hybrid. These describe the inspected commit, not all releases.

Existing vector/pdfbdf output was timed and sized; a separate Ghostscript 10.07.1 process rasterized it to PNG, including startup and PNG encoding. Vector generation avoids local page rasterization for PDF-capable destinations; downstream rasterization remains. Physical printers, halftones and code-reading quality were not tested.

Measurements and input/source hashes are retained under `experiments/go`, `experiments/bracket`, `experiments/native` and `experiments/verification`. Results are single-machine, bounded-input observations. Proposed production priorities remain inferred recommendations, not measured universal rankings.


## Production adoption on 2026-10-15

At the user's request, explicit AA OFF and the three selected techniques are now adopted in the normal imagebdf code: completed-image LRU (default 64MiB), whole fill-text coverage LRU (default 32MiB), and persistent independent-page workers through Pages/PagesDPI and CLI -pages/-workers. Zero cache limits select defaults, negatives disable caches. Images are privately retained and copied on return. Page/view layout, target geometry, AA and background distinguish completed images; text masks distinguish the final transform, font reference, text and layout state. Mutable effects are applied after mask lookup.

Page workers retain independent font/geometry/scratch state and share only immutable text masks through a locked LRU. Callback order is deterministic and output ownership is independent. Worker concurrency is capped by a 32M-pixel base-canvas budget; retained working surfaces share a total 64MiB allowance. Invalidate releases renderer resources, caches and worker state after Part/FontFS changes. Directory/system font scan/load caches retain their existing process lifetime semantics.

Production tests, targeted race checks and post-adoption benchmarks are saved in `production/`. See [implementation.ja.md](implementation.ja.md). The three-document 72/300dpi AA on/off comparisons matched when text masks were enabled; independent pages also matched serial outputs. Capacity, eviction, output mutation, resource invalidation, FontFS refresh, ordering, warnings, cancellation and CLI validation were tested. All repository Go packages passed.

Production warm cache measurements are distinct from the earlier candidate experiments. Five 1s samples on M3/GOMAXPROCS=8 showed completed-image hit medians of roughly 0.8–3.4ms, and a PDF-derived text-only cache median of 98.6ms versus 143.1ms without caches. A four-request page batch with the memory cap improved from 407.7ms to 253.0ms. Execution variance was substantial; smaller improvements are not guaranteed. Neither these numbers nor the earlier native geometry comparison establish first-render or physical-printer performance. Native geometry tooling disables the adopted caches; historical overlays and adapter sources remain saved for replay.

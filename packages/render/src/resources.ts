import { type BdfDocument, type ObjectPart, type PathData, type PartEntry, type Hash, type Font, Verb, FontKind, FONT_STYLES } from "@bdf/core";
import { VectorImage, domSvgRasterizer, isSvg, type Raster, type SvgRasterizer } from "./svg.js";

/** Build a Path2D from path data. */
export function buildPath2D(p: PathData): Path2D {
  const path = new Path2D();
  const a = p.args;
  let i = 0;
  for (const v of p.verbs) {
    switch (v) {
      case Verb.MOVE: path.moveTo(a[i], a[i + 1]); i += 2; break;
      case Verb.LINE: path.lineTo(a[i], a[i + 1]); i += 2; break;
      case Verb.QUAD: path.quadraticCurveTo(a[i], a[i + 1], a[i + 2], a[i + 3]); i += 4; break;
      case Verb.CUBIC: path.bezierCurveTo(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4], a[i + 5]); i += 6; break;
      case Verb.CLOSE: path.closePath(); break;
      case Verb.RECT: path.rect(a[i], a[i + 1], a[i + 2], a[i + 3]); i += 4; break;
      case Verb.ELLIPSE: path.ellipse(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4], a[i + 5], a[i + 6], a[i + 7] !== 0); i += 8; break;
      case Verb.ARC_TO: path.arcTo(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4]); i += 5; break;
      case Verb.ROUND_RECT: path.roundRect(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4]); i += 5; break;
    }
  }
  return path;
}

/** Font family name used for an embedded font part. */
export function embeddedFamily(hash: Hash): string {
  return `bdf-${hash}`;
}

/** CSS font string for a FONT instruction. */
export function fontString(f: Font, size: number): string {
  let family: string;
  if (f.kind === FontKind.EMBEDDED && f.hash) {
    family = `"${embeddedFamily(f.hash)}"`;
    if (f.family) family += `, ${f.family}`;
  } else {
    family = f.family || "sans-serif";
  }
  return `${FONT_STYLES[f.style] ?? "normal"} ${f.weight || 400} ${size}px ${family}`;
}

/**
 * Caches of browser objects derived from parts: Path2D, ImageBitmap, FontFace.
 * One cache can serve many contexts; it is bound to a document.
 *
 * SVG images are drawn by a rasterizer at the sizes draws ask for (see
 * vector() and settle()): by default with an image element, which workers
 * do not have; the worker asks the page (BdfWorkerClient).
 */
export class ResourceCache {
  private inlinePaths = new WeakMap<ObjectPart, (Path2D | undefined)[]>();
  private extPaths = new Map<Hash, Path2D[]>();
  private images = new Map<Hash, ImageBitmap>();
  private vectors = new Map<Hash, VectorImage>();
  private fonts = new Map<Hash, FontFace>();
  private fontSet: FontFaceSet | undefined;
  private rasterize: SvgRasterizer | undefined;
  /** SVG rasters the draws since takeMisses() did not find, by image and scale. */
  private misses = new Map<string, { vec: VectorImage; k: number }>();

  constructor(readonly doc: BdfDocument, fontSet?: FontFaceSet, rasterize?: SvgRasterizer) {
    this.fontSet = fontSet ?? (globalThis as { fonts?: FontFaceSet }).fonts ?? (globalThis as { document?: { fonts?: FontFaceSet } }).document?.fonts;
    this.rasterize = rasterize ?? domSvgRasterizer();
  }

  /** Load an object, its children, and every font/image/path part they use. */
  async prepare(hash: Hash): Promise<ObjectPart> {
    return this.doc.ensure(hash, (e, bytes) => this.load(e, bytes));
  }

  private async load(e: PartEntry, bytes: Uint8Array): Promise<void> {
    switch (e.t) {
      case "img": {
        if (this.images.has(e.h) || this.vectors.has(e.h)) return;
        if (isSvg(bytes)) {
          this.vectors.set(e.h, new VectorImage(e.h, bytes));
          return;
        }
        const bmp = await createImageBitmap(new Blob([bytes as BlobPart]));
        this.images.set(e.h, bmp);
        return;
      }
      case "font": {
        if (this.fonts.has(e.h)) return;
        const buffer = bytes.buffer.slice(bytes.byteOffset, bytes.byteOffset + bytes.byteLength) as ArrayBuffer;
        const face = new FontFace(embeddedFamily(e.h), buffer);
        await face.load();
        this.fontSet?.add(face);
        this.fonts.set(e.h, face);
        return;
      }
      case "path": {
        if (this.extPaths.has(e.h)) return;
        const paths = await this.doc.pathCollection(e.h);
        this.extPaths.set(e.h, paths.map(buildPath2D));
        return;
      }
    }
  }

  path(o: ObjectPart, index: number): Path2D {
    const entry = o.paths[index];
    if (!entry) throw new Error(`bdf: bad path ref ${index}`);
    if ("inline" in entry) {
      let list = this.inlinePaths.get(o);
      if (!list) {
        list = new Array(o.paths.length);
        this.inlinePaths.set(o, list);
      }
      let p = list[index];
      if (!p) {
        p = buildPath2D(entry.inline);
        list[index] = p;
      }
      return p;
    }
    const set = this.extPaths.get(entry.hash);
    const p = set?.[entry.index];
    if (!p) throw new Error(`bdf: path collection ${entry.hash} not loaded`);
    return p;
  }

  /** A bitmap image; for an SVG image, the raster drawn last. */
  image(hash: Hash): ImageBitmap {
    const img = this.images.get(hash) ?? this.vectors.get(hash)?.latest();
    if (!img) throw new Error(`bdf: image ${hash} not loaded`);
    return img;
  }

  /** An SVG image, or undefined for a bitmap. */
  vector(hash: Hash): VectorImage | undefined {
    return this.vectors.get(hash);
  }

  /**
   * A raster of an SVG image for drawing at scale (device pixels per image
   * pixel). When none fits, the closest one (or none) is returned and the
   * miss is noted for settle().
   */
  svgRaster(vec: VectorImage, scale: number): Raster | undefined {
    const { raster, missing } = vec.raster(scale);
    if (missing !== undefined) this.misses.set(`${vec.hash} ${missing}`, { vec, k: missing });
    return raster;
  }

  /**
   * The SVG rasters the draws since the last call did not find. Draws are
   * synchronous, so taking them before and after a draw gives its own.
   */
  takeMisses(): { vec: VectorImage; k: number }[] {
    const m = [...this.misses.values()];
    this.misses.clear();
    return m;
  }

  /** Draw the rasters a draw missed; true when any was drawn, and the draw should be run again. */
  async settle(misses: { vec: VectorImage; k: number }[]): Promise<boolean> {
    const drawn = await Promise.all(misses.map(({ vec, k }) => vec.draw(k, this.rasterize)));
    return drawn.includes(true);
  }

  object(hash: Hash): ObjectPart {
    return this.doc.objectSync(hash);
  }
}

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

/** The decoded images a cache keeps by default: 256 MiB, counted as width × height × 4 bytes. */
export const DEFAULT_IMAGE_BUDGET = 256 * 1024 * 1024;

export interface ResourceOptions {
  /**
   * Bytes of decoded images to keep (width × height × 4 each; default
   * DEFAULT_IMAGE_BUDGET). The least recently prepared images beyond it are
   * closed, except those a render holds; they are decoded again from their
   * part when needed.
   */
  imageBudget?: number;
  /** Decodes an image part (default createImageBitmap). */
  decodeImage?: (bytes: Uint8Array) => Promise<ImageBitmap>;
  /**
   * Draws SVG images, which createImageBitmap does not decode (default: an
   * image element, which workers do not have; the worker asks the page).
   */
  rasterizeSvg?: SvgRasterizer;
}

/**
 * The images one render uses, from its prepare() calls until release():
 * the cache does not close them in between, whatever else it loads.
 */
export class ImageHold {
  readonly images = new Set<Hash>();
}

const bitmapBytes = (b: ImageBitmap) => b.width * b.height * 4;

/**
 * Caches of browser objects derived from parts: Path2D, ImageBitmap, FontFace.
 * One cache can serve many contexts; it is bound to a document.
 *
 * Decoded images take far more memory than anything else (an A4 scan at
 * 192 dpi is 14 MB as RGBA), so they are kept within a budget, least
 * recently used first out. A render takes a hold() before preparing and
 * releases it after drawing, so an image cannot be closed between the two
 * even when another render finishes in the meantime.
 *
 * SVG images are drawn by a rasterizer at the sizes draws ask for (see
 * svgRaster() and settle()): by default with an image element, which workers
 * do not have; the worker asks the page (BdfWorkerClient).
 */
export class ResourceCache {
  private inlinePaths = new WeakMap<ObjectPart, (Path2D | undefined)[]>();
  private extPaths = new Map<Hash, Path2D[]>();
  /** Decoded images, least recently prepared first. */
  private images = new Map<Hash, ImageBitmap>();
  private decoding = new Map<Hash, Promise<void>>();
  private holds = new Set<ImageHold>();
  private bytes = 0;
  private disposed = false;
  private vectors = new Map<Hash, VectorImage>();
  private fonts = new Map<Hash, FontFace>();
  private fontSet: FontFaceSet | undefined;
  readonly imageBudget: number;
  private decodeImage: (bytes: Uint8Array) => Promise<ImageBitmap>;
  private rasterize: SvgRasterizer | undefined;
  /** SVG rasters the draws since takeMisses() did not find, by image and scale. */
  private misses = new Map<string, { vec: VectorImage; k: number }>();

  constructor(readonly doc: BdfDocument, fontSet?: FontFaceSet, opts: ResourceOptions = {}) {
    this.fontSet = fontSet ?? (globalThis as { fonts?: FontFaceSet }).fonts ?? (globalThis as { document?: { fonts?: FontFaceSet } }).document?.fonts;
    this.imageBudget = opts.imageBudget ?? DEFAULT_IMAGE_BUDGET;
    this.decodeImage = opts.decodeImage ?? ((bytes) => createImageBitmap(new Blob([bytes as BlobPart])));
    this.rasterize = opts.rasterizeSvg ?? domSvgRasterizer();
  }

  /**
   * Load an object, its children, and every font/image/path part they use.
   * The images go into hold, when given, and stay decoded until it is released.
   */
  async prepare(hash: Hash, hold?: ImageHold): Promise<ObjectPart> {
    return this.doc.ensure(hash, (e, bytes) => this.load(e, bytes, hold));
  }

  /**
   * Load an object with its children, fonts and path collections but not
   * its images: what extracting its text needs, without decoding pictures
   * that no render asked for.
   */
  async prepareText(hash: Hash): Promise<ObjectPart> {
    return this.doc.ensure(hash, (e, bytes) => (e.t === "img" ? undefined : this.load(e, bytes)));
  }

  /** Start holding the images of a render (see prepare). */
  hold(): ImageHold {
    const h = new ImageHold();
    this.holds.add(h);
    return h;
  }

  /** End a hold, and close the images over the budget that nothing holds. */
  release(hold: ImageHold): void {
    this.holds.delete(hold);
    this.trim();
  }

  /** The decoded images and their size in bytes (width × height × 4). */
  get imageStats(): { count: number; bytes: number } {
    return { count: this.images.size, bytes: this.bytes };
  }

  /**
   * Close the decoded images; the cache is not used for new renders
   * afterwards. Those a render in progress holds are closed when it
   * releases them, and those still decoding when they are done.
   */
  dispose(): void {
    this.disposed = true;
    this.trim();
    for (const v of this.vectors.values()) v.dispose();
  }

  private async load(e: PartEntry, bytes: Uint8Array, hold?: ImageHold): Promise<void> {
    switch (e.t) {
      case "img": {
        if (isSvg(bytes)) {
          // drawn at the sizes draws ask for (svgRaster, settle), not decoded here
          if (!this.vectors.has(e.h) && !this.disposed) this.vectors.set(e.h, new VectorImage(e.h, bytes));
          return;
        }
        hold?.images.add(e.h); // held from now, before any await
        const bmp = this.images.get(e.h);
        if (bmp) {
          this.images.delete(e.h); // most recently prepared: to the end
          this.images.set(e.h, bmp);
          return;
        }
        // Renders that need the same image at once share one decoding.
        let p = this.decoding.get(e.h);
        if (!p) {
          p = this.decodeImage(bytes)
            .then((bmp) => {
              if (this.disposed) {
                bmp.close();
                return;
              }
              this.images.set(e.h, bmp);
              this.bytes += bitmapBytes(bmp);
              this.trim();
            })
            .finally(() => this.decoding.delete(e.h));
          this.decoding.set(e.h, p);
        }
        return p;
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

  /** Close the least recently prepared images that nothing holds until the rest fit the budget. */
  private trim(): void {
    const budget = this.disposed ? 0 : this.imageBudget;
    if (this.bytes <= budget) return;
    const held = new Set<Hash>();
    for (const h of this.holds) for (const i of h.images) held.add(i);
    for (const [hash, bmp] of this.images) {
      if (this.bytes <= budget) break;
      if (held.has(hash)) continue;
      this.images.delete(hash);
      this.bytes -= bitmapBytes(bmp);
      bmp.close();
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

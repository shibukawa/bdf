import { type BdfDocument, type ObjectPart, type PathData, type PartEntry, type Hash, type Font, BdfFormatError, Verb, FontKind, FONT_STYLES } from "@bdf/core";
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

/**
 * Images of more pixels than this are not decoded, by default: 1 GiB
 * decoded, more than the photos of cameras have. The header of an image
 * states its size, and a few bytes can state any: a PNG of 30 KB decodes
 * into a gigabyte.
 */
export const DEFAULT_MAX_IMAGE_PIXELS = 1 << 28;

/** The decoded images one render holds at most, by default: 4 GiB, counted as width × height × 4 bytes. */
export const DEFAULT_HOLD_LIMIT = 4 * 1024 * 1024 * 1024;

/**
 * The size in pixels that the header of a PNG, JPEG, GIF, WebP or BMP image
 * states (what decoding it allocates); undefined for other images and for
 * headers cut short.
 */
export function imageSize(b: Uint8Array): { width: number; height: number } | undefined {
  const dv = new DataView(b.buffer, b.byteOffset, b.byteLength);
  const tag = (at: number, s: string) => at + s.length <= b.length && [...s].every((c, i) => b[at + i] === c.charCodeAt(0));
  if (tag(0, "\x89PNG\r\n\x1a\n")) {
    return b.length >= 24 && tag(12, "IHDR") ? { width: dv.getUint32(16), height: dv.getUint32(20) } : undefined;
  }
  if (tag(0, "GIF87a") || tag(0, "GIF89a")) {
    return b.length >= 10 ? { width: dv.getUint16(6, true), height: dv.getUint16(8, true) } : undefined;
  }
  if (tag(0, "BM")) {
    if (b.length < 26) return undefined;
    // the OS/2 header has 16 bit sizes; a negative height is a top-down image
    if (dv.getUint32(14, true) === 12) return { width: dv.getUint16(18, true), height: dv.getUint16(20, true) };
    return { width: Math.abs(dv.getInt32(18, true)), height: Math.abs(dv.getInt32(22, true)) };
  }
  if (tag(0, "RIFF") && tag(8, "WEBP")) {
    if (b.length < 30) return undefined;
    if (tag(12, "VP8X")) return { width: 1 + (b[24] | (b[25] << 8) | (b[26] << 16)), height: 1 + (b[27] | (b[28] << 8) | (b[29] << 16)) };
    if (tag(12, "VP8 ")) return { width: dv.getUint16(26, true) & 0x3fff, height: dv.getUint16(28, true) & 0x3fff };
    if (tag(12, "VP8L")) {
      const v = dv.getUint32(21, true);
      return { width: 1 + (v & 0x3fff), height: 1 + ((v >>> 14) & 0x3fff) };
    }
    return undefined;
  }
  if (b[0] === 0xff && b[1] === 0xd8) {
    // the first frame header (SOF), after the segments before it
    for (let at = 2; at + 4 <= b.length; ) {
      if (b[at] !== 0xff) return undefined;
      const m = b[at + 1];
      if (m === 0xff) { at++; continue; } // fill
      if (m === 0x01 || (m >= 0xd0 && m <= 0xd8)) { at += 2; continue; } // no length
      if (m === 0xd9 || m === 0xda) return undefined; // the image data, with no frame header before it
      if (m >= 0xc0 && m <= 0xcf && m !== 0xc4 && m !== 0xc8 && m !== 0xcc) {
        return at + 9 <= b.length ? { width: dv.getUint16(at + 7), height: dv.getUint16(at + 5) } : undefined;
      }
      at += 2 + dv.getUint16(at + 2);
    }
  }
  return undefined;
}

export interface ResourceOptions {
  /**
   * Bytes of decoded images to keep (width × height × 4 each; default
   * DEFAULT_IMAGE_BUDGET). The least recently prepared images beyond it are
   * closed, except those a render holds; they are decoded again from their
   * part when needed.
   */
  imageBudget?: number;
  /**
   * Images of more pixels than this are not decoded, and drawn as nothing
   * (default DEFAULT_MAX_IMAGE_PIXELS).
   */
  maxImagePixels?: number;
  /**
   * Bytes of decoded images one render may hold (default
   * DEFAULT_HOLD_LIMIT): preparing a page whose images take more fails.
   */
  holdLimit?: number;
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
  /** What the images take decoded, as far as it is known. */
  bytes = 0;
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
  /** Images that are not decoded: too large. */
  private refused = new Set<Hash>();
  private vectors = new Map<Hash, VectorImage>();
  private fonts = new Map<Hash, FontFace>();
  private fontSet: FontFaceSet | undefined;
  readonly imageBudget: number;
  readonly maxImagePixels: number;
  readonly holdLimit: number;
  private decodeImage: (bytes: Uint8Array) => Promise<ImageBitmap>;
  private rasterize: SvgRasterizer | undefined;
  /** SVG rasters the draws since takeMisses() did not find, by image and scale. */
  private misses = new Map<string, { vec: VectorImage; k: number }>();

  constructor(readonly doc: BdfDocument, fontSet?: FontFaceSet, opts: ResourceOptions = {}) {
    this.fontSet = fontSet ?? (globalThis as { fonts?: FontFaceSet }).fonts ?? (globalThis as { document?: { fonts?: FontFaceSet } }).document?.fonts;
    this.imageBudget = opts.imageBudget ?? DEFAULT_IMAGE_BUDGET;
    this.maxImagePixels = opts.maxImagePixels ?? DEFAULT_MAX_IMAGE_PIXELS;
    this.holdLimit = opts.holdLimit ?? DEFAULT_HOLD_LIMIT;
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
   * Close the decoded images and take the fonts out of the font set; the
   * cache is not used for new renders afterwards. Images a render in
   * progress holds are closed when it releases them, and those still
   * decoding when they are done; the fonts go at once, so dispose of a
   * cache whose renders may still draw text only once they are done.
   */
  dispose(): void {
    this.disposed = true;
    this.trim();
    for (const v of this.vectors.values()) v.dispose();
    for (const face of this.fonts.values()) this.fontSet?.delete(face);
    this.fonts.clear();
    this.extPaths.clear();
  }

  private async load(e: PartEntry, bytes: Uint8Array, hold?: ImageHold): Promise<void> {
    switch (e.t) {
      case "img": {
        if (isSvg(bytes)) {
          // drawn at the sizes draws ask for (svgRaster, settle), not decoded here
          if (!this.vectors.has(e.h) && !this.disposed) this.vectors.set(e.h, new VectorImage(e.h, bytes));
          return;
        }
        if (this.refused.has(e.h)) return;
        const bmp = this.images.get(e.h);
        // what decoding allocates is known before: the header states it
        const stated = bmp ?? imageSize(bytes);
        if (stated && stated.width * stated.height > this.maxImagePixels) return this.refuse(e.h, stated);
        const held = hold?.images.has(e.h);
        hold?.images.add(e.h); // held from now, before any await
        if (hold && !held && stated) this.pin(hold, stated);
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
              if (bmp.width * bmp.height > this.maxImagePixels) {
                bmp.close();
                return this.refuse(e.h, bmp);
              }
              this.images.set(e.h, bmp);
              this.bytes += bitmapBytes(bmp);
              this.trim();
            })
            .finally(() => this.decoding.delete(e.h));
          this.decoding.set(e.h, p);
        }
        if (!hold || held || stated) return p;
        // an image whose header was not read: what it takes is known now
        return p.then(() => {
          const decoded = this.images.get(e.h);
          if (decoded) this.pin(hold, decoded);
        });
      }
      case "font": {
        if (this.fonts.has(e.h)) return;
        const buffer = bytes.buffer.slice(bytes.byteOffset, bytes.byteOffset + bytes.byteLength) as ArrayBuffer;
        const face = new FontFace(embeddedFamily(e.h), buffer);
        await face.load();
        if (this.disposed) return;
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

  /** An image that is not decoded: draws find nothing to draw. */
  private refuse(hash: Hash, size: { width: number; height: number }): void {
    this.refused.add(hash);
    console.warn(`bdf: image ${hash} of ${size.width} × ${size.height} pixels is not drawn`);
  }

  /** Count an image of a hold; a hold of more than the limit fails the render it is for. */
  private pin(hold: ImageHold, size: { width: number; height: number }): void {
    hold.bytes += size.width * size.height * 4;
    if (hold.bytes > this.holdLimit) throw new BdfFormatError(`the images of a render take more than ${this.holdLimit} bytes`);
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

  /**
   * Whether text in a CSS font is drawn in the font's own faces: none of
   * them in the font set is still to load (the embedded fonts are loaded
   * by prepare; a page's own web fonts may load later).
   */
  fontLoaded(font: string): boolean {
    try {
      return this.fontSet?.check(font) ?? true;
    } catch {
      return false;
    }
  }

  /**
   * A bitmap image; for an SVG image, the raster drawn last. Undefined for
   * an image that was not decoded for its size: it is drawn as nothing.
   */
  image(hash: Hash): ImageBitmap | undefined {
    const img = this.images.get(hash) ?? this.vectors.get(hash)?.latest();
    if (!img && !this.refused.has(hash)) throw new Error(`bdf: image ${hash} not loaded`);
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

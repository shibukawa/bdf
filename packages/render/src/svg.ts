// SVG image parts (docs/spec.md §6.2). createImageBitmap does not decode
// SVG, in a worker or on the page: an SVG image is drawn with an image
// element onto a canvas, at the size it is drawn at, so that it stays sharp
// at every zoom. Workers have no image element; the page draws for them
// (BdfWorkerClient answers the worker's requests).
import type { Hash } from "@bdf/core";

/** Draws an SVG image into a bitmap of width × height pixels. */
export type SvgRasterizer = (hash: Hash, data: Uint8Array, width: number, height: number) => Promise<ImageBitmap>;

/** Whether image bytes are SVG (markup, possibly after a byte order mark) rather than an encoded bitmap. */
export function isSvg(b: Uint8Array): boolean {
  if ((b[0] === 0xff && b[1] === 0xfe) || (b[0] === 0xfe && b[1] === 0xff)) return true; // UTF-16
  let i = b[0] === 0xef && b[1] === 0xbb && b[2] === 0xbf ? 3 : 0;
  while (b[i] === 0x20 || b[i] === 0x09 || b[i] === 0x0a || b[i] === 0x0d) i++;
  return b[i] === 0x3c; // "<"
}

const UNITS: Record<string, number> = { "": 1, px: 1, pt: 96 / 72, pc: 16, in: 96, cm: 96 / 2.54, mm: 96 / 25.4, q: 96 / 101.6, em: 16, ex: 8 };

/** An absolute length in CSS px; undefined for percentages and invalid lengths. */
function svgLength(s: string | undefined): number | undefined {
  const m = /^\s*([+-]?(?:\d+\.?\d*|\.\d+)(?:[eE][+-]?\d+)?)\s*([A-Za-z%]*)\s*$/.exec(s ?? "");
  if (!m) return undefined;
  const v = Number(m[1]) * (UNITS[m[2].toLowerCase()] ?? NaN);
  return v > 0 ? v : undefined;
}

/**
 * The size of an SVG image in CSS px, which IMAGE_SUB's source rectangle and
 * a pattern's cells measure (docs/spec.md §6.2): the width and height of its
 * root, a missing one following from the other and the view box's
 * proportions, else the view box's size, else 300 × 150.
 */
export function svgSize(data: Uint8Array): { width: number; height: number } {
  const label = data[0] === 0xff && data[1] === 0xfe ? "utf-16le" : data[0] === 0xfe && data[1] === 0xff ? "utf-16be" : "utf-8";
  const head = new TextDecoder(label).decode(data.subarray(0, 65536));
  const tag = /<(?:[\w.-]+:)?svg\b[^>]*>/.exec(head)?.[0] ?? "";
  const attr = (name: string) => {
    const m = new RegExp(`\\s${name}\\s*=\\s*(?:"([^"]*)"|'([^']*)')`).exec(tag);
    return m ? m[1] ?? m[2] : undefined;
  };
  let w = svgLength(attr("width")), h = svgLength(attr("height"));
  const vb = (attr("viewBox") ?? "").split(/[\s,]+/).filter(Boolean).map(Number);
  const [vw, vh] = vb.length === 4 && vb[2] > 0 && vb[3] > 0 ? [vb[2], vb[3]] : [0, 0];
  if (w === undefined || h === undefined) {
    if (w !== undefined && vw) h = (w * vh) / vw;
    else if (h !== undefined && vw) w = (h * vw) / vh;
    else if (vw) [w, h] = [vw, vh];
    w ??= 300;
    h ??= 150;
  }
  return { width: w, height: h };
}

/** Rasters larger than this many pixels are drawn smaller (and scaled up). */
const MAX_PIXELS = 4096 * 4096;
/** A raster up to this much larger than needed is used as it is. */
const SLACK = 1.25;
/** Rasters kept per image, the most recently used. */
const KEEP = 4;

export interface Raster {
  bitmap: ImageBitmap;
  /** The scale it was drawn for: about as many raster pixels per image pixel (bitmap.width / width exactly). */
  k: number;
}

/**
 * An SVG image part and its rasters. A draw asks for the scale it needs;
 * when no raster fits, it gets the closest one (or none) and the scale it
 * missed, to draw (draw()) before drawing again.
 */
export class VectorImage {
  readonly width: number;
  readonly height: number;
  private rasters: Raster[] = [];
  private pending = new Map<number, Promise<boolean>>();
  private failed = new Set<number>();

  readonly data: Uint8Array;

  constructor(readonly hash: Hash, data: Uint8Array) {
    // a view into a larger buffer (a document read whole) would send all of it to the page
    this.data = data.byteLength === data.buffer.byteLength ? data : data.slice();
    ({ width: this.width, height: this.height } = svgSize(data));
  }

  /** The scale of the raster for drawing at scale (device pixels per image pixel): bounded in size. */
  private fit(scale: number): number {
    const max = Math.sqrt(MAX_PIXELS / (this.width * this.height));
    return Math.min(Math.max(scale, 1 / Math.min(this.width, this.height)), max);
  }

  /**
   * A raster for drawing at scale, and the scale of the raster that is
   * missing when none fits (the closest one is returned meanwhile).
   */
  raster(scale: number): { raster?: Raster; missing?: number } {
    const k = this.fit(scale);
    const i = this.rasters.findIndex((r) => r.k >= k && r.k <= k * SLACK);
    if (i >= 0) {
      const [r] = this.rasters.splice(i, 1);
      this.rasters.push(r); // most recently used last
      return { raster: r };
    }
    let best: Raster | undefined;
    for (const r of this.rasters) if (!best || Math.abs(Math.log(r.k / k)) < Math.abs(Math.log(best.k / k))) best = r;
    return { raster: best, missing: this.failed.has(k) ? undefined : k };
  }

  /** The raster used last. */
  latest(): ImageBitmap | undefined {
    return this.rasters.at(-1)?.bitmap;
  }

  /** Draw the raster of scale k (once, however many draws missed it); false when it cannot be drawn. */
  draw(k: number, rasterize: SvgRasterizer | undefined): Promise<boolean> {
    if (this.rasters.some((r) => r.k === k)) return Promise.resolve(true);
    let p = this.pending.get(k);
    if (p) return p;
    p = (async () => {
      const w = Math.max(1, Math.round(this.width * k)), h = Math.max(1, Math.round(this.height * k));
      try {
        if (!rasterize) throw new Error("no SVG rasterizer (workers need the page to draw SVG images)");
        this.rasters.push({ bitmap: await rasterize(this.hash, this.data, w, h), k });
        while (this.rasters.length > KEEP) this.rasters.shift()!.bitmap.close();
        return true;
      } catch (e) {
        this.failed.add(k);
        console.warn(`bdf: SVG image ${this.hash}: ${(e as Error).message ?? e}`);
        return false;
      } finally {
        this.pending.delete(k);
      }
    })();
    this.pending.set(k, p);
    return p;
  }
}

/**
 * A rasterizer that decodes SVG with an image element, where there is a DOM
 * (undefined in workers). It keeps the last images it decoded, by hash.
 */
export function domSvgRasterizer(): SvgRasterizer | undefined {
  if (typeof document === "undefined" || typeof Image === "undefined") return undefined;
  const decoded = new Map<Hash, Promise<HTMLImageElement>>();
  const load = (hash: Hash, data: Uint8Array) => {
    let p = decoded.get(hash);
    if (p) {
      decoded.delete(hash); // most recently used last
    } else {
      p = (async () => {
        const url = URL.createObjectURL(new Blob([data as BlobPart], { type: "image/svg+xml" }));
        const img = new Image();
        img.src = url;
        try {
          await img.decode();
        } finally {
          URL.revokeObjectURL(url);
        }
        return img;
      })();
      p.catch(() => decoded.delete(hash));
    }
    decoded.set(hash, p);
    if (decoded.size > 32) decoded.delete(decoded.keys().next().value!);
    return p;
  };
  return async (hash, data, width, height) => {
    const img = await load(hash, data);
    // drawImage lays an SVG image out at the destination size (createImageBitmap's resize does not, for one without a size)
    if (typeof OffscreenCanvas !== "undefined") {
      const canvas = new OffscreenCanvas(width, height);
      canvas.getContext("2d")!.drawImage(img, 0, 0, width, height);
      return canvas.transferToImageBitmap();
    }
    const canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;
    canvas.getContext("2d")!.drawImage(img, 0, 0, width, height);
    return createImageBitmap(canvas);
  };
}

import {
  type ObjectPart, type OpSink, type Glyph, type Paint, type Rect, walk, UseLimits, BdfFormatError,
  BLEND_NAMES, LINE_CAPS, LINE_JOINS, TEXT_ALIGNS, TEXT_BASELINES, TEXT_DIRECTIONS, FILL_RULES, REPEATS, SMOOTHING_QUALITIES, PaintKind, MaskKind,
} from "@bdf/core";
import { ResourceCache, fontString } from "./resources.js";

/** Any 2D context: on-screen canvas or OffscreenCanvas. */
export type Ctx2D = CanvasRenderingContext2D | OffscreenCanvasRenderingContext2D;

const colorCache = new Map<number, string>();

/** CSS color for a packed 0xRRGGBBAA value. */
export function cssColor(rgba: number): string {
  let s = colorCache.get(rgba);
  if (s === undefined) {
    const a = rgba & 0xff;
    const r = (rgba >>> 24) & 0xff, g = (rgba >>> 16) & 0xff, b = (rgba >>> 8) & 0xff;
    s = a === 255 ? `#${(rgba >>> 8).toString(16).padStart(6, "0")}` : `rgba(${r},${g},${b},${(a / 255).toFixed(4)})`;
    colorCache.set(rgba, s);
  }
  return s;
}

export interface RenderOptions {
  /** Relative tolerance before applying advance correction (default 0.005). */
  advanceTolerance?: number;
  /** Factory for temporary canvases used by GROUP_BEGIN and MASK_BEGIN. */
  createCanvas?: (w: number, h: number) => OffscreenCanvas | HTMLCanvasElement;
}

/** A rectangle as its edges. */
interface Box {
  x0: number;
  y0: number;
  x1: number;
  y1: number;
}

/** The size of a CSS font, in px (10 when it cannot be told). */
function fontSize(font: string): number {
  const m = /(\d*\.?\d+(?:e[+-]?\d+)?)px/i.exec(font);
  return m ? Number(m[1]) : 10;
}

/** The text widths a renderer keeps before it starts over. */
const MAX_WIDTHS = 200_000;

/**
 * How far from its anchor, in font sizes, the ink of a run may reach
 * besides its advance (on either side, for every alignment): for the
 * baseline (top, middle, alphabetic, bottom), accents and stacked
 * combining marks, and italic overhang.
 */
const INK_REACH = 4;

/**
 * The drawing state that decides what text can be skipped: the visible
 * box in the current coordinates (none once they cannot be told, or a
 * shadow or a filter spreads the ink), and the font size.
 */
interface Cull {
  box: Box | undefined;
  size: number;
}

/**
 * Groups and soft masks open at once: each has a canvas of its own, as large
 * as the part of it on the canvas it is drawn on.
 */
export const MAX_GROUP_DEPTH = 64;

/**
 * How far from a canvas the ink of a group still shows on it, in device
 * pixels: the shadow and the filter that the group is composited with spread
 * it. A page is drawn in bands and regions, so what lies beyond the edge of a
 * canvas belongs to the page all the same.
 */
const MAX_INK_REACH = 4096;

function inkReach(ctx: Ctx2D): number {
  // the blur of a shadow is twice its standard deviation, which reaches three times itself
  let reach = Math.max(Math.abs(ctx.shadowOffsetX), Math.abs(ctx.shadowOffsetY)) + 1.5 * ctx.shadowBlur;
  const filter = "filter" in ctx ? (ctx as CanvasRenderingContext2D).filter : "";
  if (filter && filter !== "none") {
    // the lengths of a blur or a drop shadow: as far as all of them together, three times over
    for (const m of filter.matchAll(/-?\d*\.?\d+/g)) reach += 3 * Math.abs(Number(m[0]));
  }
  return Number.isFinite(reach) ? Math.min(Math.ceil(reach), MAX_INK_REACH) : MAX_INK_REACH;
}

/** Let the pixels of a temporary canvas go now, not when it is collected. */
function release(canvas: OffscreenCanvas | HTMLCanvasElement): void {
  canvas.width = canvas.height = 0;
}

// FILTER: the filter functions of CSS and the colors of a drop shadow. A
// url() would name a filter of the page around the canvas, or a file.
const FILTER_FUNCTIONS = new Set([
  "blur", "brightness", "contrast", "drop-shadow", "grayscale", "hue-rotate", "invert", "opacity", "saturate", "sepia",
  "rgb", "rgba", "hsl", "hsla", "hwb", "lab", "lch", "oklab", "oklch", "color",
]);

/** Whether a FILTER string is a list of filter functions and nothing else. */
export function filterAllowed(css: string): boolean {
  // no escapes, strings nor anything else a function name or a url could hide in
  if (!/^[a-zA-Z0-9\s.,%#()+\/-]*$/.test(css)) return false;
  for (const m of css.matchAll(/([a-zA-Z-]*)\(/g)) if (!FILTER_FUNCTIONS.has(m[1].toLowerCase())) return false;
  return true;
}

function defaultCreateCanvas(w: number, h: number): OffscreenCanvas | HTMLCanvasElement {
  if (typeof OffscreenCanvas !== "undefined") return new OffscreenCanvas(w, h);
  const c = document.createElement("canvas");
  c.width = w;
  c.height = h;
  return c;
}

interface Group {
  ctx: Ctx2D;
  canvas: OffscreenCanvas | HTMLCanvasElement;
  alpha: number;
  blend: number;
  dx: number;
  dy: number;
  cull: Cull; // that of ctx, where drawing resumes
}

/** A soft mask being drawn for a group (MASK_BEGIN … MASK_END). */
interface Mask {
  ctx: Ctx2D; // the context drawing resumes on
  canvas: OffscreenCanvas | HTMLCanvasElement;
  mctx: Ctx2D;
  kind: number;
  transfer: Uint8Array | undefined;
  group: Group | undefined;
  cull: Cull; // that of ctx
}

/** Put a context back into the Canvas 2D initial drawing state (docs/spec.md §8). */
export function resetState(ctx: Ctx2D): void {
  ctx.fillStyle = "#000000";
  ctx.strokeStyle = "#000000";
  ctx.lineWidth = 1;
  ctx.lineCap = "butt";
  ctx.lineJoin = "miter";
  ctx.miterLimit = 10;
  ctx.setLineDash([]);
  ctx.lineDashOffset = 0;
  ctx.globalAlpha = 1;
  ctx.globalCompositeOperation = "source-over";
  ctx.shadowColor = "rgba(0,0,0,0)";
  ctx.shadowBlur = 0;
  ctx.shadowOffsetX = 0;
  ctx.shadowOffsetY = 0;
  ctx.font = "10px sans-serif";
  ctx.textAlign = "left";
  ctx.textBaseline = "alphabetic";
  ctx.direction = "inherit";
  if ("letterSpacing" in ctx) (ctx as CanvasRenderingContext2D).letterSpacing = "0px";
  if ("filter" in ctx) (ctx as CanvasRenderingContext2D).filter = "none";
  ctx.imageSmoothingEnabled = true;
}

/**
 * Executes object instructions against a Canvas 2D context.
 * All referenced resources must be loaded first (ResourceCache.prepare).
 */
export class CanvasRenderer implements OpSink {
  private ctx!: Ctx2D;
  private obj!: ObjectPart;
  private groups: Group[] = [];
  private masks: Mask[] = [];
  private readonly tol: number;
  private readonly createCanvas: (w: number, h: number) => OffscreenCanvas | HTMLCanvasElement;
  /** What text can be skipped now, and as of each SAVE not restored yet. */
  private cull: Cull = { box: undefined, size: 10 };
  private culls: Cull[] = [];
  /**
   * Widths that measureText gave, by font (with the letter spacing and
   * direction) and text: the advance correction measures every run each
   * time it is drawn, and measuring takes longer than drawing.
   */
  private widths = new Map<string, Map<string, number>>();
  private widthCount = 0;
  /** The limits of the objects drawn with USE, how deep the object drawn now is, and whether it is drawn again. */
  private limits = new UseLimits();
  private depth = 0;
  private reused = false;

  constructor(readonly res: ResourceCache, opts: RenderOptions = {}) {
    this.tol = opts.advanceTolerance ?? 0.005;
    this.createCanvas = opts.createCanvas ?? defaultCreateCanvas;
  }

  /**
   * Draw an object with the context's current transform and clip. Top-level
   * objects (pages, tiles) start from the initial drawing state; USE'd children
   * inherit the state of their parent. visible is the part of the object
   * that can be seen, in its coordinates: text wholly outside it is skipped
   * (a tile of a sheet has thousands of runs, and a region shows a few).
   * limits counts the objects drawn with USE across the objects of a render
   * (default: of this object alone); past them BdfFormatError is thrown.
   */
  draw(ctx: Ctx2D, obj: ObjectPart, reset = true, visible?: Rect, limits = new UseLimits()): void {
    const prevCtx = this.ctx, prevObj = this.obj, prevCull = this.cull, prevCulls = this.culls;
    const prevLimits = this.limits, prevDepth = this.depth, prevReused = this.reused;
    this.ctx = ctx;
    this.obj = obj;
    this.limits = limits;
    this.depth = 0;
    this.reused = false;
    ctx.save();
    if (reset) resetState(ctx);
    this.cull = {
      box: visible && { x0: visible.x, y0: visible.y, x1: visible.x + visible.w, y1: visible.y + visible.h },
      size: fontSize(ctx.font),
    };
    this.culls = [];
    try {
      walk(obj, this);
    } finally {
      // Unwind groups and masks left open by a malformed stream.
      for (const m of this.masks.splice(0)) release(m.canvas);
      while (this.groups.length) this.groupEnd();
      ctx.restore();
      this.ctx = prevCtx;
      this.obj = prevObj;
      this.cull = prevCull;
      this.culls = prevCulls;
      this.limits = prevLimits;
      this.depth = prevDepth;
      this.reused = prevReused;
    }
  }

  /** Map the visible box into the coordinates a transform sets up (their bounding box when it rotates or skews). */
  private transformCull(a: number, b: number, c: number, d: number, e: number, f: number) {
    const v = this.cull.box;
    if (!v) return;
    const det = a * d - b * c;
    if (!det || !Number.isFinite(det)) {
      this.cull = { ...this.cull, box: undefined };
      return;
    }
    // the inverse of the transform, for the corners
    const ia = d / det, ib = -b / det, ic = -c / det, id = a / det;
    let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
    for (const [px, py] of [[v.x0, v.y0], [v.x1, v.y0], [v.x0, v.y1], [v.x1, v.y1]]) {
      const x = ia * (px - e) + ic * (py - f), y = ib * (px - e) + id * (py - f);
      x0 = Math.min(x0, x); y0 = Math.min(y0, y); x1 = Math.max(x1, x); y1 = Math.max(y1, y);
    }
    this.cull = { ...this.cull, box: { x0, y0, x1, y1 } };
  }

  /** Whether a run with its anchor at x, y and that advance cannot reach the visible box. */
  private unseen(x: number, y: number, advance: number): boolean {
    const v = this.cull.box;
    if (!v) return false;
    const m = this.cull.size * INK_REACH;
    if (y - m > v.y1 || y + m < v.y0) return true;
    // an unknown advance: the run may be of any length
    return advance > 0 && (x - advance - m > v.x1 || x + advance + m < v.x0);
  }

  /** measureText's width of a text in the current font, remembered while the font is loaded. */
  private width(text: string): number {
    const ctx = this.ctx;
    let key = ctx.font;
    const ls = (ctx as CanvasRenderingContext2D).letterSpacing;
    if (ls && ls !== "0px") key += ` ${ls}`;
    if (ctx.direction === "rtl") key += " rtl";
    let widths = this.widths.get(key);
    if (!widths) {
      // a font still loading is measured in another for now
      if (!this.res.fontLoaded(ctx.font)) return ctx.measureText(text).width;
      this.widths.set(key, (widths = new Map()));
    }
    let w = widths.get(text);
    if (w === undefined) {
      w = ctx.measureText(text).width;
      if (++this.widthCount > MAX_WIDTHS) {
        this.widths.clear();
        this.widths.set(key, (widths = new Map()));
        this.widthCount = 1;
      }
      widths.set(text, w);
    }
    return w;
  }

  private paint(index: number): CanvasGradient | CanvasPattern | string {
    const p: Paint | undefined = this.obj.paints[index];
    if (!p) throw new Error(`bdf: bad paint ref ${index}`);
    const ctx = this.ctx;
    const c = p.coords;
    let g: CanvasGradient;
    switch (p.kind) {
      case PaintKind.LINEAR: g = ctx.createLinearGradient(c[0], c[1], c[2], c[3]); break;
      case PaintKind.RADIAL: g = ctx.createRadialGradient(c[0], c[1], c[2], c[3], c[4], c[5]); break;
      case PaintKind.CONIC: g = ctx.createConicGradient(c[0], c[1], c[2]); break;
      case PaintKind.PATTERN: {
        const hash = this.obj.images[p.image];
        const m = p.matrix;
        let img: ImageBitmap | undefined, kx = 1, ky = 1;
        const vec = this.res.vector(hash);
        if (vec) {
          // an SVG cell as large as it is drawn: device pixels per image pixel
          const t = ctx.getTransform();
          const a = t.a * m[0] + t.c * m[1], b = t.b * m[0] + t.d * m[1], c = t.a * m[2] + t.c * m[3], d = t.b * m[2] + t.d * m[3];
          img = this.res.svgRaster(vec, Math.max(Math.hypot(a, b), Math.hypot(c, d)))?.bitmap;
          if (img) [kx, ky] = [img.width / vec.width, img.height / vec.height];
        } else {
          img = this.res.image(hash);
        }
        if (!img) return "rgba(0,0,0,0)"; // an SVG image not drawn yet (the page is drawn again), or an image too large
        const pat = ctx.createPattern(img, REPEATS[p.repeat] ?? "repeat");
        if (!pat) throw new Error("bdf: createPattern failed");
        if (!(m[0] === 1 && m[1] === 0 && m[2] === 0 && m[3] === 1 && m[4] === 0 && m[5] === 0 && kx === 1 && ky === 1)) {
          pat.setTransform(new DOMMatrix([m[0], m[1], m[2], m[3], m[4], m[5]]).scale(1 / kx, 1 / ky));
        }
        return pat;
      }
      default: throw new Error(`bdf: unknown paint kind ${p.kind}`);
    }
    for (const s of p.stops) g.addColorStop(Math.min(1, Math.max(0, s.offset)), cssColor(s.color));
    return g;
  }

  // --- state ---
  save() {
    this.culls.push(this.cull);
    this.ctx.save();
  }
  restore() {
    // one RESTORE too many restores the state draw() started from, whose box is not known here
    this.cull = this.culls.pop() ?? { ...this.cull, box: undefined };
    this.ctx.restore();
  }
  transform(a: number, b: number, c: number, d: number, e: number, f: number) {
    this.ctx.transform(a, b, c, d, e, f);
    this.transformCull(a, b, c, d, e, f);
  }
  translate(x: number, y: number) {
    this.ctx.translate(x, y);
    this.transformCull(1, 0, 0, 1, x, y);
  }
  scale(x: number, y: number) {
    this.ctx.scale(x, y);
    this.transformCull(x, 0, 0, y, 0, 0);
  }
  clipPath(path: number, rule: number) { this.ctx.clip(this.res.path(this.obj, path), FILL_RULES[rule]); }
  clipRect(x: number, y: number, w: number, h: number) {
    this.ctx.beginPath();
    this.ctx.rect(x, y, w, h);
    this.ctx.clip();
  }

  // --- style ---
  fillColor(rgba: number) { this.ctx.fillStyle = cssColor(rgba); }
  fillPaint(paint: number) { this.ctx.fillStyle = this.paint(paint); }
  strokeColor(rgba: number) { this.ctx.strokeStyle = cssColor(rgba); }
  strokePaint(paint: number) { this.ctx.strokeStyle = this.paint(paint); }
  line(width: number, cap: number, join: number, miter: number) {
    const ctx = this.ctx;
    ctx.lineWidth = width;
    ctx.lineCap = LINE_CAPS[cap] ?? "butt";
    ctx.lineJoin = LINE_JOINS[join] ?? "miter";
    ctx.miterLimit = miter;
  }
  dash(segments: Float32Array, offset: number) {
    this.ctx.setLineDash(Array.from(segments));
    this.ctx.lineDashOffset = offset;
  }
  alpha(a: number) { this.ctx.globalAlpha = a; }
  blend(mode: number) { this.ctx.globalCompositeOperation = BLEND_NAMES[mode] ?? "source-over"; }
  shadow(rgba: number, blur: number, dx: number, dy: number) {
    const ctx = this.ctx;
    // Canvas shadows ignore the transform; SHADOW is in units, so scale it
    // by the current zoom (the offset keeps its page direction).
    const m = ctx.getTransform();
    const s = Math.sqrt(Math.abs(m.a * m.d - m.b * m.c)) || 1;
    ctx.shadowColor = cssColor(rgba);
    ctx.shadowBlur = blur * s;
    ctx.shadowOffsetX = dx * s;
    ctx.shadowOffsetY = dy * s;
    // a shadow off the text is ink away from it
    if ((rgba & 0xff) !== 0 && (blur !== 0 || dx !== 0 || dy !== 0)) this.cull = { ...this.cull, box: undefined };
  }
  filter(css: string) {
    // anything but filter functions is as a filter that is not supported: left out
    if (!filterAllowed(css)) css = "none";
    if ("filter" in this.ctx) (this.ctx as CanvasRenderingContext2D).filter = css;
    if (css && css !== "none") this.cull = { ...this.cull, box: undefined }; // a blur or a drop shadow spreads the ink
  }
  font(font: number, size: number) {
    const f = this.obj.fonts[font];
    if (!f) throw new Error(`bdf: bad font ref ${font}`);
    this.ctx.font = fontString(f, size);
    this.cull = { ...this.cull, size: Math.abs(size) };
  }
  textStyle(align: number, baseline: number, dir: number, letterSpacing: number) {
    const ctx = this.ctx;
    ctx.textAlign = TEXT_ALIGNS[align] ?? "left";
    ctx.textBaseline = TEXT_BASELINES[baseline] ?? "alphabetic";
    ctx.direction = TEXT_DIRECTIONS[dir] ?? "inherit";
    if ("letterSpacing" in ctx) (ctx as CanvasRenderingContext2D).letterSpacing = `${letterSpacing}px`;
  }

  // --- shapes ---
  fillRect(x: number, y: number, w: number, h: number) { this.ctx.fillRect(x, y, w, h); }
  strokeRect(x: number, y: number, w: number, h: number) { this.ctx.strokeRect(x, y, w, h); }
  fillPath(path: number, rule: number) { this.ctx.fill(this.res.path(this.obj, path), FILL_RULES[rule]); }
  strokePath(path: number) { this.ctx.stroke(this.res.path(this.obj, path)); }
  fillPathAt(path: number, rule: number, x: number, y: number) {
    const ctx = this.ctx;
    ctx.translate(x, y);
    ctx.fill(this.res.path(this.obj, path), FILL_RULES[rule]);
    ctx.translate(-x, -y);
  }
  fillPathRun(rule: number, glyphs: Glyph[]) {
    const ctx = this.ctx;
    const fr = FILL_RULES[rule];
    for (const g of glyphs) {
      ctx.translate(g.x, g.y);
      ctx.fill(this.res.path(this.obj, g.path), fr);
      ctx.translate(-g.x, -g.y);
    }
  }
  clearRect(x: number, y: number, w: number, h: number) { this.ctx.clearRect(x, y, w, h); }

  // --- text ---
  private text(text: string, x: number, y: number, advance: number, stroke: boolean) {
    const ctx = this.ctx;
    if (advance > 0) {
      const measured = this.width(text);
      if (measured > 0 && Math.abs(measured - advance) / advance > this.tol) {
        ctx.save();
        ctx.translate(x, y);
        ctx.scale(advance / measured, 1);
        if (stroke) ctx.strokeText(text, 0, 0); else ctx.fillText(text, 0, 0);
        ctx.restore();
        return;
      }
    }
    if (stroke) ctx.strokeText(text, x, y); else ctx.fillText(text, x, y);
  }
  fillText(text: string, x: number, y: number, advance: number) {
    if (!this.unseen(x, y, advance)) this.text(text, x, y, advance, false);
  }
  // always drawn: the ink of a stroke reaches as far as its line is wide
  strokeText(text: string, x: number, y: number, advance: number) { this.text(text, x, y, advance, true); }

  // --- images ---
  image(img: number, x: number, y: number, w: number, h: number) {
    const hash = this.obj.images[img];
    const vec = this.res.vector(hash);
    if (!vec) {
      const bmp = this.res.image(hash);
      if (bmp) this.ctx.drawImage(bmp, x, y, w, h);
      return;
    }
    const r = this.res.svgRaster(vec, this.deviceScale(w / vec.width, h / vec.height));
    if (r) this.ctx.drawImage(r.bitmap, x, y, w, h);
  }
  imageSub(img: number, sx: number, sy: number, sw: number, sh: number, dx: number, dy: number, dw: number, dh: number) {
    const hash = this.obj.images[img];
    const vec = this.res.vector(hash);
    if (!vec) {
      const bmp = this.res.image(hash);
      if (bmp) this.ctx.drawImage(bmp, sx, sy, sw, sh, dx, dy, dw, dh);
      return;
    }
    const r = this.res.svgRaster(vec, this.deviceScale(dw / sw, dh / sh));
    if (!r) return;
    // the source rectangle is in the image's pixels (spec §6.2), the raster has more or fewer
    const kx = r.bitmap.width / vec.width, ky = r.bitmap.height / vec.height;
    this.ctx.drawImage(r.bitmap, sx * kx, sy * ky, sw * kx, sh * ky, dx, dy, dw, dh);
  }
  /** Device pixels one image pixel covers, drawn sx by sy units per image pixel. */
  private deviceScale(sx: number, sy: number): number {
    const m = this.ctx.getTransform();
    return Math.max(Math.hypot(m.a, m.b) * Math.abs(sx), Math.hypot(m.c, m.d) * Math.abs(sy)) || 1;
  }
  smoothing(enabled: boolean, quality: number) {
    this.ctx.imageSmoothingEnabled = enabled;
    this.ctx.imageSmoothingQuality = SMOOTHING_QUALITIES[quality] ?? "low";
  }

  // --- composition ---
  use(obj: number) { this.useAt(obj, 0, 0); }
  useAt(obj: number, x: number, y: number) {
    const hash = this.obj.objects[obj];
    if (hash === undefined) throw new Error(`bdf: bad object ref ${obj}`);
    const child = this.res.object(hash);
    const again = this.limits.enter(hash);
    this.limits.descend(this.depth);
    const ctx = this.ctx;
    ctx.save();
    const cull = this.cull, depth = this.culls.length;
    if (x !== 0 || y !== 0) this.translate(x, y);
    const parent = this.obj, reused = this.reused;
    this.obj = child;
    this.depth++;
    this.reused = reused || again;
    try {
      walk(child, this, this.reused ? this.limits.count : undefined);
    } finally {
      this.obj = parent;
      this.depth--;
      this.reused = reused;
      ctx.restore();
      // a child that leaves a SAVE open leaves ctx in its state
      this.cull = this.culls.length === depth ? cull : { ...cull, box: undefined };
      this.culls.length = Math.min(this.culls.length, depth);
    }
  }
  groupBegin(alpha: number, blend: number, x: number, y: number, w: number, h: number) {
    if (this.groups.length + this.masks.length >= MAX_GROUP_DEPTH) throw new BdfFormatError("groups nested too deep");
    const ctx = this.ctx;
    const m = ctx.getTransform();
    // Device-space bounds of the group rectangle.
    const pts = [m.transformPoint({ x, y }), m.transformPoint({ x: x + w, y }), m.transformPoint({ x, y: y + h }), m.transformPoint({ x: x + w, y: y + h })];
    let x0 = Math.floor(Math.min(...pts.map((p) => p.x))), y0 = Math.floor(Math.min(...pts.map((p) => p.y)));
    let x1 = Math.ceil(Math.max(...pts.map((p) => p.x))), y1 = Math.ceil(Math.max(...pts.map((p) => p.y)));
    // Only the part on the canvas it is drawn on, and as far around it as
    // its shadow or filter reach: the rest of it is never seen, and the
    // rectangle is as large as the document says.
    const on = ctx.canvas as { width: number; height: number } | undefined;
    if (typeof on?.width === "number" && typeof on.height === "number") {
      const reach = inkReach(ctx);
      x0 = Math.max(x0, -reach);
      y0 = Math.max(y0, -reach);
      x1 = Math.min(x1, on.width + reach);
      y1 = Math.min(y1, on.height + reach);
    }
    const cw = Math.max(1, x1 - x0), ch = Math.max(1, y1 - y0);
    const canvas = this.createCanvas(cw, ch);
    const gctx = canvas.getContext("2d") as Ctx2D;
    gctx.setTransform(m.a, m.b, m.c, m.d, m.e - x0, m.f - y0);
    gctx.font = ctx.font;
    gctx.fillStyle = ctx.fillStyle;
    gctx.strokeStyle = ctx.strokeStyle;
    gctx.lineWidth = ctx.lineWidth;
    this.groups.push({ ctx, canvas, alpha, blend, dx: x0, dy: y0, cull: this.cull });
    this.ctx = gctx;
  }
  groupEnd() {
    const g = this.groups.pop();
    if (!g) return;
    const ctx = g.ctx;
    this.ctx = ctx;
    this.cull = g.cull;
    ctx.save();
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.globalAlpha = g.alpha;
    ctx.globalCompositeOperation = BLEND_NAMES[g.blend] ?? "source-over";
    ctx.drawImage(g.canvas, g.dx, g.dy);
    ctx.restore();
    release(g.canvas);
  }
  maskBegin(kind: number, backdrop: number, transfer: Uint8Array) {
    if (this.groups.length + this.masks.length >= MAX_GROUP_DEPTH) throw new BdfFormatError("groups nested too deep");
    // The mask covers the innermost group's canvas and starts from the
    // initial drawing state with the current transform.
    const group = this.groups[this.groups.length - 1];
    const w = group?.canvas.width ?? 1, h = group?.canvas.height ?? 1;
    const canvas = this.createCanvas(w, h);
    const mctx = canvas.getContext("2d") as Ctx2D;
    if (kind === MaskKind.LUMINOSITY) {
      mctx.fillStyle = cssColor((backdrop | 0xff) >>> 0);
      mctx.fillRect(0, 0, w, h);
    }
    resetState(mctx);
    mctx.setTransform(this.ctx.getTransform());
    mctx.save();
    this.masks.push({ ctx: this.ctx, canvas, mctx, kind, transfer: transfer.length === 256 ? transfer : undefined, group, cull: this.cull });
    this.ctx = mctx;
    this.cull = { ...this.cull, size: 10 };
  }
  maskEnd() {
    const m = this.masks.pop();
    if (!m) return;
    this.ctx = m.ctx;
    this.cull = m.cull;
    const g = m.group;
    if (!g || this.groups[this.groups.length - 1] !== g) return release(m.canvas);
    const { mctx, canvas } = m;
    mctx.restore();
    const w = canvas.width, h = canvas.height;
    const lum = m.kind === MaskKind.LUMINOSITY, tr = m.transfer;
    if ((lum || tr) && w > 0 && h > 0) {
      // Luminosity to alpha (the weights of the PDF non-separable blend
      // modes), then the transfer function.
      const img = mctx.getImageData(0, 0, w, h);
      const d = img.data;
      for (let i = 0; i < d.length; i += 4) {
        let v = lum ? Math.round(0.3 * d[i] + 0.59 * d[i + 1] + 0.11 * d[i + 2]) : d[i + 3];
        if (tr) v = tr[v];
        d[i + 3] = v;
        d[i] = d[i + 1] = d[i + 2] = 0;
      }
      mctx.putImageData(img, 0, 0);
    }
    // Keep the group's pixels where the mask is: the masked result replaces
    // the group canvas (a clip left on the group context cannot limit it).
    mctx.save();
    mctx.setTransform(1, 0, 0, 1, 0, 0);
    mctx.globalAlpha = 1;
    mctx.globalCompositeOperation = "source-in";
    mctx.drawImage(g.canvas, 0, 0);
    mctx.restore();
    release(g.canvas);
    g.canvas = canvas;
  }

  // --- meta ---
  link() {}
  mark() {}
  ext() {}
}

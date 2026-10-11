import {
  type BdfDocument, type ObjectPart, type OpSink, type Glyph, type Paint, type Font, type PathData, type Rect, type Hash, type View, type Page,
  walk, UseLimits, tileSize, BLEND_NAMES, Verb, PaintKind, FontKind, MaskKind,
} from "@bdfkit/core";
import { imageSize } from "./resources.js";
import { isSvg, svgSize } from "./svg.js";
import { tilesIn } from "./page.js";

// The pages of a document as SVG: a file that a browser displays by itself,
// for where no viewer runs. This is the writer of the Go package
// vector/svgbdf, with the fonts and images in the SVG as data: URLs (its
// Options.Inline): what a viewer saves is one file. The two write the same
// elements; see the Go package for what SVG has no element for.

export interface SvgOptions {
  /** Fill color behind the page, as #rrggbb (default white; null for none). */
  background?: string | null;
}

/** The SVG of a page or a region, and what could not be drawn as the viewer draws it. */
export interface SvgResult {
  svg: string;
  warnings: string[];
}

/** An affine transform as Canvas 2D keeps it: x' = a x + c y + e, y' = b x + d y + f. */
type Matrix = [number, number, number, number, number, number];

const IDENTITY: Matrix = [1, 0, 0, 1, 0, 0];

/** m × n: n applies first. */
function mul(m: Matrix, n: Matrix): Matrix {
  return [
    m[0] * n[0] + m[2] * n[1], m[1] * n[0] + m[3] * n[1],
    m[0] * n[2] + m[2] * n[3], m[1] * n[2] + m[3] * n[3],
    m[0] * n[4] + m[2] * n[5] + m[4], m[1] * n[4] + m[3] * n[5] + m[5],
  ];
}

const finite = (m: Matrix) => m.every(Number.isFinite);
const isIdentity = (m: Matrix) => m[0] === 1 && m[1] === 0 && m[2] === 0 && m[3] === 1 && m[4] === 0 && m[5] === 0;
const maxScale = (m: Matrix) => Math.max(Math.hypot(m[0], m[1]), Math.hypot(m[2], m[3]));

/** How many decimals a length in the units of the page is written with. */
const PAGE_DIGITS = 3;

/** The decimals of the coordinates of what is drawn under m, to be as exact on the page as PAGE_DIGITS. */
function digitsOf(m: Matrix): number {
  const s = maxScale(m);
  if (!(s > 0) || !Number.isFinite(s)) return PAGE_DIGITS;
  return Math.min(Math.max(Math.ceil(Math.log10(s)) + PAGE_DIGITS, 0), 9);
}

/** A number with at most d decimals, as short as SVG reads it. */
function num(v: number, d: number): string {
  if (!Number.isFinite(v)) return "0";
  let s = v.toFixed(d);
  if (d > 0) s = s.replace(/\.?0+$/, "");
  if (s === "-0") return "0";
  if (s.startsWith("0.")) return s.slice(1);
  if (s.startsWith("-0.")) return "-" + s.slice(2);
  return s;
}

/** A number that scales others (a matrix entry, an opacity), with six significant digits. */
function factor(v: number): string {
  if (v === 0 || !Number.isFinite(v)) return "0";
  return num(v, Math.min(Math.max(5 - Math.floor(Math.log10(Math.abs(v))), 0), 12));
}

/** The transform attribute of what is drawn under m, with a space before it. */
function transformOf(m: Matrix): string {
  if (isIdentity(m)) return "";
  if (m[0] === 1 && m[1] === 0 && m[2] === 0 && m[3] === 1) return ` transform="translate(${num(m[4], PAGE_DIGITS)} ${num(m[5], PAGE_DIGITS)})"`;
  return ` transform="matrix(${factor(m[0])} ${factor(m[1])} ${factor(m[2])} ${factor(m[3])} ${num(m[4], PAGE_DIGITS)} ${num(m[5], PAGE_DIGITS)})"`;
}

/** Writes the path data of a Path2D: the commands of Canvas, which start a sub-path where SVG would have an error. */
class PathWriter {
  out = "";
  x = 0;
  y = 0;
  private sx = 0;
  private sy = 0;
  open = false;
  bad = false;

  constructor(readonly d: number) {}

  cmd(c: string, ...v: number[]): void {
    this.out += c;
    v.forEach((f, i) => {
      if (!Number.isFinite(f)) this.bad = true;
      const s = num(f, this.d);
      // a separator only where two numbers would read as one
      this.out += i > 0 && s[0] !== "-" ? " " + s : s;
    });
  }
  moveTo(x: number, y: number): void {
    this.cmd("M", x, y);
    this.x = this.sx = x;
    this.y = this.sy = y;
    this.open = true;
  }
  lineTo(x: number, y: number): void {
    if (!this.open) return this.moveTo(x, y);
    this.cmd("L", x, y);
    this.x = x;
    this.y = y;
  }
  close(): void {
    if (!this.open) return;
    this.out += "Z";
    this.x = this.sx;
    this.y = this.sy;
  }
  rect(x: number, y: number, w: number, h: number): void {
    this.moveTo(x, y);
    this.cmd("h", w);
    this.cmd("v", h);
    this.cmd("h", -w);
    this.out += "Z";
    if (!Number.isFinite(w + h)) this.bad = true;
  }
  /** An elliptical arc with Canvas semantics: a line from the current point to its start. */
  ellipse(cx: number, cy: number, rx: number, ry: number, rot: number, a0: number, a1: number, ccw: boolean): void {
    if (rx < 0 || ry < 0) return;
    const tau = 2 * Math.PI;
    let sweep: number;
    if (!ccw) {
      if (a1 - a0 >= tau) sweep = tau;
      else if ((sweep = (a1 - a0) % tau) < 0) sweep += tau;
    } else {
      if (a0 - a1 >= tau) sweep = -tau;
      else if ((sweep = -((a0 - a1) % tau)) > 0) sweep -= tau;
    }
    const sin = Math.sin(rot), cos = Math.cos(rot);
    const at = (a: number): [number, number] => {
      const x = rx * Math.cos(a), y = ry * Math.sin(a);
      return [cx + x * cos - y * sin, cy + x * sin + y * cos];
    };
    let [x, y] = at(a0);
    this.lineTo(x, y);
    if (sweep === 0 || rx === 0 || ry === 0 || Number.isNaN(sweep)) {
      [x, y] = at(a0 + sweep);
      this.lineTo(x, y);
      return;
    }
    // arcs of a third of a turn at most: one of a half or more has two readings
    const n = Math.ceil(Math.abs(sweep) / (tau / 3));
    const deg = (rot * 180) / Math.PI;
    for (let i = 1; i <= n; i++) {
      [x, y] = at(a0 + (sweep * i) / n);
      this.cmd("A", rx, ry);
      this.out += ` ${num(deg, 3)} 0 ${sweep > 0 ? 1 : 0} ${num(x, this.d)} ${num(y, this.d)}`;
      if (!Number.isFinite(x + y + deg)) this.bad = true;
    }
    this.x = x;
    this.y = y;
  }
  /** A line and an arc of radius r tangent to the lines from the current point to (x1, y1) and on to (x2, y2). */
  arcTo(x1: number, y1: number, x2: number, y2: number, r: number): void {
    if (r < 0) return;
    if (!this.open) this.moveTo(x1, y1);
    const x0 = this.x, y0 = this.y;
    if ((x0 === x1 && y0 === y1) || (x1 === x2 && y1 === y2) || r === 0) return this.lineTo(x1, y1);
    const ax = x0 - x1, ay = y0 - y1, bx = x2 - x1, by = y2 - y1;
    const la = Math.hypot(ax, ay), lb = Math.hypot(bx, by);
    const cross = ax * by - ay * bx;
    if (Math.abs(cross) < 1e-12 * la * lb) return this.lineTo(x1, y1); // collinear
    const theta = Math.acos(Math.max(-1, Math.min(1, (ax * bx + ay * by) / (la * lb))));
    const d = r / Math.tan(theta / 2);
    const t0x = x1 + (ax / la) * d, t0y = y1 + (ay / la) * d;
    const t1x = x1 + (bx / lb) * d, t1y = y1 + (by / lb) * d;
    // centre: along the bisector
    const bisx = ax / la + bx / lb, bisy = ay / la + by / lb;
    const bl = Math.hypot(bisx, bisy), h = r / Math.sin(theta / 2);
    const cx = x1 + (bisx / bl) * h, cy = y1 + (bisy / bl) * h;
    this.lineTo(t0x, t0y);
    this.ellipse(cx, cy, r, r, 0, Math.atan2(t0y - cy, t0x - cx), Math.atan2(t1y - cy, t1x - cx), cross > 0);
  }
  /** A rectangle with corners rounded by r (Canvas roundRect with one radius). */
  roundRect(x: number, y: number, w: number, h: number, r: number): void {
    if (w < 0) [x, w] = [x + w, -w];
    if (h < 0) [y, h] = [y + h, -h];
    r = Math.max(0, Math.min(r, Math.min(w, h) / 2));
    if (r === 0 || Number.isNaN(r)) return this.rect(x, y, w, h);
    const q = Math.PI / 2;
    this.moveTo(x + r, y);
    this.lineTo(x + w - r, y);
    this.ellipse(x + w - r, y + r, r, r, 0, -q, 0, false);
    this.lineTo(x + w, y + h - r);
    this.ellipse(x + w - r, y + h - r, r, r, 0, 0, q, false);
    this.lineTo(x + r, y + h);
    this.ellipse(x + r, y + h - r, r, r, 0, q, 2 * q, false);
    this.lineTo(x, y + r);
    this.ellipse(x + r, y + r, r, r, 0, 2 * q, 3 * q, false);
    this.out += "Z";
    this.x = this.sx = x;
    this.y = this.sy = y;
  }
}

/** The d attribute of a path with d decimals; undefined for a path with nothing to draw or a coordinate that is no number. */
export function svgPathData(p: PathData, d: number): string | undefined {
  const w = new PathWriter(d);
  const a = p.args;
  let i = 0;
  for (const v of p.verbs) {
    switch (v) {
      case Verb.MOVE: w.moveTo(a[i], a[i + 1]); i += 2; break;
      case Verb.LINE: w.lineTo(a[i], a[i + 1]); i += 2; break;
      case Verb.QUAD:
        if (!w.open) w.moveTo(a[i], a[i + 1]);
        w.cmd("Q", a[i], a[i + 1], a[i + 2], a[i + 3]);
        w.x = a[i + 2];
        w.y = a[i + 3];
        i += 4;
        break;
      case Verb.CUBIC:
        if (!w.open) w.moveTo(a[i], a[i + 1]);
        w.cmd("C", a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4], a[i + 5]);
        w.x = a[i + 4];
        w.y = a[i + 5];
        i += 6;
        break;
      case Verb.CLOSE: w.close(); break;
      case Verb.RECT: w.rect(a[i], a[i + 1], a[i + 2], a[i + 3]); i += 4; break;
      case Verb.ELLIPSE: w.ellipse(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4], a[i + 5], a[i + 6], a[i + 7] !== 0); i += 8; break;
      case Verb.ARC_TO: w.arcTo(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4]); i += 5; break;
      case Verb.ROUND_RECT: w.roundRect(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4]); i += 5; break;
    }
  }
  return w.bad || w.out === "" ? undefined : w.out;
}

function rectData(x: number, y: number, w: number, h: number, d: number): string | undefined {
  const p = new PathWriter(d);
  p.rect(x, y, w, h);
  return p.bad ? undefined : p.out;
}

/** A string as XML holds it; what XML has no character for becomes U+FFFD, and the white space Canvas draws as a space is one. */
function escape(s: string): string {
  // eslint-disable-next-line no-control-regex
  if (!/[\u0000-\u001f\ud800-\udfff￾￿&<>"]/.test(s)) return s;
  let out = "";
  for (const ch of s) {
    const c = ch.codePointAt(0)!;
    if (c === 0x09 || c === 0x0a || c === 0x0d || c === 0x0c || c === 0x0b) out += " ";
    else if (c < 0x20 || c === 0xfffe || c === 0xffff || (c >= 0xd800 && c <= 0xdfff)) out += "�";
    else if (ch === "&") out += "&amp;";
    else if (ch === "<") out += "&lt;";
    else if (ch === ">") out += "&gt;";
    else if (ch === '"') out += "&quot;";
    else out += ch;
  }
  return out;
}

function hexColor(rgba: number): string {
  const v = (rgba >>> 8).toString(16).padStart(6, "0");
  return v[0] === v[1] && v[2] === v[3] && v[4] === v[5] ? `#${v[0]}${v[2]}${v[4]}` : `#${v}`;
}

function base64(bytes: Uint8Array): string {
  let s = "";
  for (let i = 0; i < bytes.length; i += 0x8000) s += String.fromCharCode(...bytes.subarray(i, i + 0x8000));
  return btoa(s);
}

const tag = (b: Uint8Array, at: number, s: string) => at + s.length <= b.length && [...s].every((c, i) => b[at + i] === c.charCodeAt(0));

/** The media type of an image part; undefined for a format browsers do not show. */
function imageType(b: Uint8Array): string | undefined {
  if (tag(b, 0, "\x89PNG\r\n\x1a\n")) return "image/png";
  if (b[0] === 0xff && b[1] === 0xd8 && b[2] === 0xff) return "image/jpeg";
  if (tag(b, 0, "RIFF") && tag(b, 8, "WEBP")) return "image/webp";
  if (tag(b, 4, "ftyp") && (tag(b, 8, "avif") || tag(b, 8, "avis"))) return "image/avif";
  if (tag(b, 0, "GIF87a") || tag(b, 0, "GIF89a")) return "image/gif";
  if (tag(b, 0, "BM")) return "image/bmp";
  return isSvg(b) ? "image/svg+xml" : undefined;
}

/**
 * Whether the EXIF orientation of a JPEG or PNG image turns it on its side,
 * as browsers show it: its width and height are then the other way round.
 */
function turned(b: Uint8Array): boolean {
  let exif: Uint8Array | undefined;
  const dv = new DataView(b.buffer, b.byteOffset, b.byteLength);
  if (b[0] === 0xff && b[1] === 0xd8) {
    for (let pos = 2; pos + 4 <= b.length && b[pos] === 0xff; ) {
      const marker = b[pos + 1];
      if (marker === 0xda || marker === 0xd9) break; // start of scan, end of image
      const n = dv.getUint16(pos + 2);
      if (n < 2 || pos + 2 + n > b.length) break;
      if (marker === 0xe1 && tag(b, pos + 4, "Exif\0\0")) {
        exif = b.subarray(pos + 10, pos + 2 + n);
        break;
      }
      pos += 2 + n;
    }
  } else if (tag(b, 0, "\x89PNG\r\n\x1a\n")) {
    for (let pos = 8; pos + 12 <= b.length; ) {
      const n = dv.getUint32(pos);
      if (pos + 12 + n > b.length) break;
      if (tag(b, pos + 4, "eXIf")) {
        exif = b.subarray(pos + 8, pos + 8 + n);
        if (tag(exif, 0, "Exif\0\0")) exif = exif.subarray(6);
        break;
      }
      pos += 12 + n;
    }
  }
  if (!exif || exif.length < 8) return false;
  const little = tag(exif, 0, "II*\0");
  if (!little && !tag(exif, 0, "MM\0*")) return false;
  const ev = new DataView(exif.buffer, exif.byteOffset, exif.byteLength);
  const ifd = ev.getUint32(4, little);
  if (ifd + 2 > exif.length) return false;
  const count = ev.getUint16(ifd, little);
  for (let i = 0; i < count; i++) {
    const e = ifd + 2 + 12 * i;
    if (e + 12 > exif.length) break;
    // Orientation, SHORT
    if (ev.getUint16(e, little) === 0x0112 && ev.getUint16(e + 2, little) === 3) {
      const o = ev.getUint16(e + 8, little);
      return o >= 5 && o <= 8;
    }
  }
  return false;
}

/** The size the header of an AVIF image states: that of the largest of its image items. */
function avifSize(b: Uint8Array): { width: number; height: number } | undefined {
  const dv = new DataView(b.buffer, b.byteOffset, b.byteLength);
  const end = Math.min(b.length, 1 << 16);
  let width = 0, height = 0;
  for (let i = 0; i + 16 <= end; i++) {
    if (b[i] !== 0x69 || !tag(b, i, "ispe")) continue;
    const w = dv.getUint32(i + 8), h = dv.getUint32(i + 12);
    if (w * h > width * height) [width, height] = [w, h];
    i += 15;
  }
  return width > 0 && height > 0 ? { width, height } : undefined;
}

/** The size of an image in pixels, as browsers show it; undefined when its header does not tell. */
function shownSize(b: Uint8Array, type: string): { width: number; height: number } | undefined {
  if (type === "image/svg+xml") return svgSize(b);
  if (type === "image/avif") return avifSize(b);
  const size = imageSize(b);
  return size && size.width > 0 && size.height > 0 && turned(b) ? { width: size.height, height: size.width } : size;
}

function fontType(b: Uint8Array): string {
  if (tag(b, 0, "wOF2")) return "font/woff2";
  if (tag(b, 0, "wOFF")) return "font/woff";
  if (tag(b, 0, "OTTO")) return "font/otf";
  return tag(b, 0, "ttcf") ? "font/collection" : "font/ttf";
}

/** A fill or stroke style: a color, or a gradient or pattern of an object. */
interface PaintRef {
  color: number;
  paint?: Paint;
  obj?: ObjectPart; // whose images a pattern refers to
}

/** The drawing state of a context (what save and restore keep). */
interface State {
  m: Matrix;
  clip: number; // the clip path, 0 for none
  fill: PaintRef;
  stroke: PaintRef;
  width: number;
  miter: number;
  cap: number;
  join: number;
  dash: Float32Array | undefined;
  dashOffset: number;
  alpha: number;
  blend: number;
  shadow: number;
  shadowBlur: number; // in units of the page, like the offsets
  shadowDX: number;
  shadowDY: number;
  font: Font | undefined;
  size: number;
  align: number;
  baseline: number;
  dir: number;
  spacing: number;
  smooth: boolean;
}

/** The Canvas 2D initial state with a transform and clip. */
function initialState(m: Matrix, clip: number): State {
  return {
    m, clip, fill: { color: 0x000000ff }, stroke: { color: 0x000000ff }, width: 1, miter: 10, cap: 0, join: 0,
    dash: undefined, dashOffset: 0, alpha: 1, blend: 0, shadow: 0, shadowBlur: 0, shadowDX: 0, shadowDY: 0,
    font: undefined, size: 10, align: 0, baseline: 0, dir: 0, spacing: 0, smooth: true,
  };
}

const shadowActive = (s: State) => (s.shadow & 0xff) !== 0 && (s.shadowBlur > 0 || s.shadowDX !== 0 || s.shadowDY !== 0);

/** A fill that waits for the stroke of the same outline, which makes one element of the two. */
interface Shape { geom: string; rule: string; tf: string; clip: number; fill: string }

/**
 * What a canvas is to the other renderers: the page, or the content of a
 * group or soft mask. Its elements have the transform they are drawn under
 * each, so that the groups around them need none.
 */
class Context {
  stack: State[] = [];
  buf = "";
  /** The clip path of the group open for the elements being written (0: none is). */
  clip = 0;
  /** The attributes of the group open inside it for elements that share them ("": none is). */
  attrs = "";
  pend: Shape | undefined;
  /** An element of the context has a blend mode, which must not reach what lies under the context. */
  blended = false;
  /** Elements that erase what buf holds (destination-out), until the next thing is drawn. */
  erase: string | undefined;
  done = false;
  constructor(public st: State) {}
}

interface Group { parent: Context; ctx: Context; alpha: number; blend: number; mask: string }
interface SoftMask { parent: Context; ctx: Context; kind: number; backdrop: number; transfer: Uint8Array | undefined; group: Group | undefined }

/** An image as an SVG refers to it. */
interface ImageInfo { hash: Hash; url: string; w: number; h: number }

/** What the objects of a region refer to, loaded before it is drawn. */
interface Loaded {
  fonts: Map<Hash, string>; // the data: URL of each
  images: Map<Hash, ImageInfo>;
  paths: Map<Hash, PathData[]>;
}

// Limits of what the drawing of a page takes (see the Go package).
const MAX_REUSED = 1 << 21, MAX_LAYERS = 64, MAX_STATES = 1024;
/** The cell of a pattern in a direction it is not repeated in: further than a page reaches. */
const UNREPEATED = 1 << 15;

const MODE_OVER = 0, MODE_ERASE = 1, MODE_KEEP = 2, MODE_UNDER = 3;
// the composite operations of BLEND_NAMES
const BLEND_LUMINOSITY = 15, BLEND_DESTINATION_OVER = 16, BLEND_DESTINATION_IN = 17, BLEND_DESTINATION_OUT = 18, BLEND_LIGHTER = 25;

const CAP_NAMES = ["", "round", "square"], JOIN_NAMES = ["", "round", "bevel"];
const BASELINE_NAMES = ["", "text-before-edge", "central", "text-after-edge", "hanging", "ideographic"];
const STYLE_NAMES = ["", "italic", "oblique"];

/** Runs the instructions of objects into one SVG. */
class SvgSink implements OpSink {
  ctx: Context;
  private obj!: ObjectPart;
  private groups: Group[] = [];
  private masks: SoftMask[] = [];
  private depth = 0;
  private limits = new UseLimits(MAX_REUSED);
  private reused = false;
  /** The groups and soft masks that are open but not drawn: those past the limit. */
  private skipped = 0;
  /** The SAVEs of each context that kept no state. */
  private unsaved = new Map<Context, number>();

  private defs = "";
  private css = "";
  private links = "";
  private ids = 0;
  private clips = new Map<string, number>();
  private paints = new Map<Paint, string>();
  private glyphs = new Map<string, string>();
  private filters = new Map<string, string>();
  private fonts = new Map<Hash, string>();
  private symbols = new Map<Hash, string>();
  readonly warnings: string[] = [];

  /** x, y, w, h: the part of the page the SVG shows, in units (its view box). */
  constructor(private doc: BdfDocument, private res: Loaded, private opts: SvgOptions, readonly x: number, readonly y: number, readonly w: number, readonly h: number) {
    this.ctx = new Context(initialState(IDENTITY, 0));
  }

  private warn(s: string): void {
    if (!this.warnings.includes(s)) this.warnings.push(s);
  }
  private id(prefix: string): string {
    return prefix + ++this.ids;
  }
  /** The attributes of a mask or filter that reaches all the SVG shows. */
  private region(): string {
    return ` x="${num(this.x, PAGE_DIGITS)}" y="${num(this.y, PAGE_DIGITS)}" width="${num(this.w, PAGE_DIGITS)}" height="${num(this.h, PAGE_DIGITS)}"`;
  }

  // --- writing elements ---

  /** Write the fill that waited for a stroke. */
  private flush(c: Context): void {
    const p = c.pend;
    if (!p) return;
    c.pend = undefined;
    this.write(c, p.clip, `<${p.geom}${p.rule}${p.tf}${p.fill}/>`);
  }
  /** Apply the elements that erase to what the context holds: they become a mask of what was drawn before them. */
  private settle(c: Context): void {
    const e = c.erase;
    if (e === undefined) return;
    c.erase = undefined;
    if (c.buf === "") return;
    const id = this.id("m");
    this.defs += `<mask id="${id}" maskUnits="userSpaceOnUse"${this.region()}><path d="${rectData(this.x, this.y, this.w, this.h, PAGE_DIGITS)}" fill="#fff"/>${e}</mask>\n`;
    c.buf = `<g mask="url(#${id})">\n${c.buf}</g>\n`;
  }
  private closeAttrs(c: Context): void {
    if (c.attrs === "") return;
    c.attrs = "";
    c.buf += "</g>\n";
  }
  /** Write what waits and close the groups open for the elements being written. */
  private closeAll(c: Context): void {
    this.flush(c);
    this.closeAttrs(c);
    if (c.clip !== 0) {
      c.clip = 0;
      c.buf += "</g>\n";
    }
  }
  private clipAttr(clip: number): string {
    return clip === 0 ? "" : ` clip-path="url(#c${clip})"`;
  }
  /**
   * Open the group of a clip path for the elements that follow. A clip path
   * is in the coordinates of what refers to it, which for an element with a
   * transform are not those of the page: the group has none.
   */
  private within(c: Context, clip: number): void {
    this.settle(c);
    if (c.clip === clip) return;
    this.closeAttrs(c);
    if (c.clip !== 0) c.buf += "</g>\n";
    c.clip = clip;
    if (clip !== 0) c.buf += `<g${this.clipAttr(clip)}>\n`;
  }
  /** Write an element drawn over what is there. */
  private write(c: Context, clip: number, el: string): void {
    this.within(c, clip);
    this.closeAttrs(c);
    c.buf += el + "\n";
  }
  /** Write an element into a group of the attributes it shares with its neighbours. */
  private writeIn(c: Context, clip: number, attrs: string, el: string): void {
    this.flush(c);
    this.within(c, clip);
    if (c.attrs !== attrs) {
      this.closeAttrs(c);
      // elements with nothing to share need no group
      c.attrs = attrs;
      if (attrs !== "") c.buf += `<g${attrs}>\n`;
    }
    c.buf += el + "\n";
  }

  /** How a composite operation is drawn, and the blend mode of CSS for it ("" for none). */
  private mode(blend: number): [number, string] {
    if (blend === 0) return [MODE_OVER, ""];
    if (blend <= BLEND_LUMINOSITY) return [MODE_OVER, BLEND_NAMES[blend]];
    if (blend === BLEND_LIGHTER) return [MODE_OVER, "plus-lighter"];
    if (blend === BLEND_DESTINATION_OUT) return [MODE_ERASE, ""];
    if (blend === BLEND_DESTINATION_IN) return [MODE_KEEP, ""];
    if (blend === BLEND_DESTINATION_OVER) return [MODE_UNDER, ""];
    if (blend < BLEND_NAMES.length) this.warn(`the composite operation ${BLEND_NAMES[blend]} is drawn as source-over`);
    return [MODE_OVER, ""];
  }

  /**
   * Draw an element as the state says: over what is there, with its blend
   * mode and shadow, or as the composite operation makes of it and what is
   * there. black: the element is drawn in black with the alpha it has.
   */
  private compose(c: Context, mode: number, blend: string, el: string, black: boolean): void {
    const st = c.st;
    switch (mode) {
      case MODE_ERASE:
        this.closeAll(c);
        if (!black) el = `<g filter="url(#${this.blackFilter()})">${el}</g>`;
        c.erase = (c.erase ?? "") + `<g${this.clipAttr(st.clip)}>${el}</g>`;
        return;
      case MODE_KEEP: {
        this.closeAll(c);
        this.settle(c);
        if (c.buf === "") return;
        const id = this.id("m");
        this.defs += `<mask id="${id}" maskUnits="userSpaceOnUse"${this.region()} style="mask-type:alpha"><g${this.clipAttr(st.clip)}>${el}</g></mask>\n`;
        c.buf = `<g mask="url(#${id})">\n${c.buf}</g>\n`;
        return;
      }
      case MODE_UNDER:
        this.closeAll(c);
        this.settle(c);
        c.buf = `<g${this.clipAttr(st.clip)}>${el}</g>\n${c.buf}`;
        return;
    }
    if (shadowActive(st)) el = `<g filter="url(#${this.shadowFilter(st)})">${el}</g>`;
    if (blend === "") {
      this.flush(c);
      this.write(c, st.clip, el);
      return;
    }
    // what blends is a child of the context: a group around it would be what it blends with
    this.closeAll(c);
    this.settle(c);
    c.blended = true;
    c.buf += `<g${this.clipAttr(st.clip)} style="mix-blend-mode:${blend}">${el}</g>\n`;
  }

  /** The filter that keeps the alpha of what it is applied to and makes it black. */
  private blackFilter(): string {
    return this.defineFilter("black", `<feColorMatrix values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0"/>`);
  }
  /** Define a filter that reaches all the SVG shows, once. */
  private defineFilter(key: string, primitives: string): string {
    let id = this.filters.get(key);
    if (id === undefined) {
      id = this.id("f");
      this.filters.set(key, id);
      this.defs += `<filter id="${id}" filterUnits="userSpaceOnUse"${this.region()} color-interpolation-filters="sRGB">${primitives}</filter>\n`;
    }
    return id;
  }
  /** The filter of the shadow of the state, applied to a group without a transform. */
  private shadowFilter(st: State): string {
    let prim = `<feDropShadow dx="${num(st.shadowDX, PAGE_DIGITS)}" dy="${num(st.shadowDY, PAGE_DIGITS)}" stdDeviation="${num(st.shadowBlur / 2, PAGE_DIGITS)}" flood-color="${hexColor(st.shadow)}"`;
    const a = st.shadow & 0xff;
    if (a !== 0xff) prim += ` flood-opacity="${factor(a / 255)}"`;
    return this.defineFilter(prim, prim + "/>");
  }

  /** The attributes of a fill or a stroke (name) in a style. */
  private paintAttrs(name: string, p: PaintRef, alpha: number, mode: number): string {
    let value: string, opacity = alpha;
    if (mode === MODE_ERASE || mode === MODE_KEEP) {
      // only its alpha counts
      if (!p.paint) opacity *= (p.color & 0xff) / 255;
      value = "#000";
    } else if (p.paint) {
      value = this.paintServer(p);
    } else {
      value = hexColor(p.color);
      opacity *= (p.color & 0xff) / 255;
    }
    let s = name === "fill" && value === "#000" ? "" : ` ${name}="${value}"`;
    if (opacity < 1) s += ` ${name}-opacity="${factor(Math.max(opacity, 0))}"`;
    return s;
  }
  private fillAttrs(st: State, mode: number): string {
    return this.paintAttrs("fill", st.fill, st.alpha, mode);
  }
  private strokeAttrs(st: State, mode: number): string {
    let s = this.paintAttrs("stroke", st.stroke, st.alpha, mode);
    if (st.width !== 1) s += ` stroke-width="${factor(st.width)}"`;
    if (CAP_NAMES[st.cap]) s += ` stroke-linecap="${CAP_NAMES[st.cap]}"`;
    if (JOIN_NAMES[st.join]) s += ` stroke-linejoin="${JOIN_NAMES[st.join]}"`;
    else if (st.miter !== 4) s += ` stroke-miterlimit="${factor(st.miter)}"`;
    if (st.dash && st.dash.length > 0) {
      s += ` stroke-dasharray="${Array.from(st.dash, (v) => factor(v)).join(" ")}"`;
      if (st.dashOffset !== 0) s += ` stroke-dashoffset="${factor(st.dashOffset)}"`;
    }
    return s;
  }

  /** Fill or stroke an outline under m; a stroke of the outline that was filled just before is the same element. */
  private shape(geom: string, rule: number, m: Matrix, stroke: boolean): void {
    const c = this.ctx, st = c.st;
    if (!finite(m)) return;
    const [mode, blend] = this.mode(st.blend);
    const tf = transformOf(m);
    const plain = mode === MODE_OVER && blend === "" && !shadowActive(st);
    if (stroke) {
      const p = c.pend;
      if (p && plain && p.geom === geom && p.tf === tf && p.clip === st.clip) {
        c.pend = undefined;
        this.write(c, st.clip, `<${geom}${p.rule}${tf}${p.fill}${this.strokeAttrs(st, mode)}/>`);
        return;
      }
      this.compose(c, mode, blend, `<${geom}${tf} fill="none"${this.strokeAttrs(st, mode)}/>`, true);
      return;
    }
    const fr = rule === 1 ? ` fill-rule="evenodd"` : "";
    const fill = this.fillAttrs(st, mode);
    if (plain) {
      this.flush(c);
      c.pend = { geom, rule: fr, tf, clip: st.clip, fill };
      return;
    }
    this.compose(c, mode, blend, `<${geom}${fr}${tf}${fill}/>`, true);
  }

  private pathOf(i: number): PathData | undefined {
    const e = this.obj.paths[i];
    if (!e) return undefined;
    return "inline" in e ? e.inline : this.res.paths.get(e.hash)?.[e.index];
  }
  private outline(i: number, rule: number, stroke: boolean): void {
    const p = this.pathOf(i);
    const m = this.ctx.st.m;
    const data = p && svgPathData(p, digitsOf(m));
    if (data !== undefined) this.shape(`path d="${data}"`, rule, m, stroke);
  }
  private rect(x: number, y: number, w: number, h: number, stroke: boolean): void {
    const m = this.ctx.st.m;
    const data = rectData(x, y, w, h, digitsOf(m));
    if (data !== undefined) this.shape(`path d="${data}"`, 0, m, stroke);
  }

  /** The id of a path that is drawn in many places, defined at its first. */
  private glyph(i: number, m: Matrix): string {
    const e = this.obj.paths[i];
    if (!e) return "";
    let key: string;
    if ("inline" in e) {
      // the paths of an object are its own: the object is numbered
      let n = this.inlineIds.get(this.obj);
      if (n === undefined) this.inlineIds.set(this.obj, (n = this.inlineIds.size + 1));
      key = `${n}/${i}`;
    } else {
      key = `${e.hash}/${e.index}`;
    }
    let id = this.glyphs.get(key);
    if (id === undefined) {
      id = "";
      const p = this.pathOf(i);
      // a digit more than its first place needs: others may be larger
      const data = p && svgPathData(p, Math.min(digitsOf(m) + 1, 9));
      if (data !== undefined) {
        id = this.id("a");
        this.defs += `<path id="${id}" d="${data}"/>\n`;
      }
      this.glyphs.set(key, id);
    }
    return id;
  }
  private inlineIds = new Map<ObjectPart, number>();

  /** Fill paths at places: the outlines of glyphs, each defined once and referred to. */
  private placed(rule: number, glyphs: Glyph[]): void {
    const c = this.ctx, st = c.st;
    if (!finite(st.m)) return;
    const [mode, blend] = this.mode(st.blend);
    let attrs = transformOf(st.m) + this.fillAttrs(st, mode);
    if (rule === 1) attrs += ` fill-rule="evenodd"`;
    const digits = digitsOf(st.m);
    const plain = mode === MODE_OVER && blend === "" && !shadowActive(st);
    let all = "";
    for (const g of glyphs) {
      const id = this.glyph(g.path, st.m);
      if (id === "" || !Number.isFinite(g.x + g.y)) continue;
      let el = `<use href="#${id}"`;
      if (g.x !== 0) el += ` x="${num(g.x, digits)}"`;
      if (g.y !== 0) el += ` y="${num(g.y, digits)}"`;
      el += "/>";
      if (plain) this.writeIn(c, st.clip, attrs, el);
      else all += el;
    }
    // one drawing, as the shadow and the blend mode of one are
    if (all !== "") this.compose(c, mode, blend, `<g${attrs}>${all}</g>`, true);
  }

  /** Cut the clip of the state to an outline under the transform of the state; a clip path refers to the one it narrows. */
  private clip(data: string | undefined, rule: number): void {
    const st = this.ctx.st;
    let el = "";
    if (data !== undefined && data !== "" && finite(st.m)) {
      el = `<path d="${data}"${transformOf(st.m)}${rule === 1 ? ` clip-rule="evenodd"` : ""}/>`;
    }
    const key = st.clip + el;
    let id = this.clips.get(key);
    if (id === undefined) {
      id = this.clips.size + 1;
      this.clips.set(key, id);
      this.defs += `<clipPath id="c${id}"${this.clipAttr(st.clip)}>${el}</clipPath>\n`;
    }
    st.clip = id;
  }

  // --- paints ---

  /** The value of fill or stroke for a gradient or pattern, defined at its first use. */
  private paintServer(p: PaintRef): string {
    let v = this.paints.get(p.paint!);
    if (v === undefined) {
      v = this.definePaint(p);
      this.paints.set(p.paint!, v);
    }
    return v;
  }
  private definePaint(p: PaintRef): string {
    const pt = p.paint!;
    const c = pt.coords;
    if (!c.every(Number.isFinite)) return "none";
    const at = (i: number) => factor(c[i]);
    let open: string, end: string;
    switch (pt.kind) {
      case PaintKind.LINEAR:
        if (c.length < 4 || (c[0] === c[2] && c[1] === c[3])) return "none"; // Canvas paints nothing with it
        open = `gradientUnits="userSpaceOnUse" x1="${at(0)}" y1="${at(1)}" x2="${at(2)}" y2="${at(3)}">`;
        end = "linearGradient";
        break;
      case PaintKind.RADIAL:
        if (c.length < 6 || c[2] < 0 || c[5] < 0 || (c[0] === c[3] && c[1] === c[4] && c[2] === c[5])) return "none";
        // the circles of Canvas: the end circle, and the start circle as the focal one
        open = `gradientUnits="userSpaceOnUse" cx="${at(3)}" cy="${at(4)}" r="${at(5)}" fx="${at(0)}" fy="${at(1)}"${c[2] !== 0 ? ` fr="${at(2)}"` : ""}>`;
        end = "radialGradient";
        break;
      case PaintKind.CONIC:
        this.warn("conic gradients are not drawn: what they paint has the color of their first stop");
        return pt.stops.length === 0 ? "none" : hexColor(pt.stops[0].color);
      case PaintKind.PATTERN:
        return this.definePattern(p);
      default:
        return "none";
    }
    if (pt.stops.length === 0) return "none";
    // in the order of their offsets, as Canvas keeps them
    const stops = pt.stops.map((s) => ({ offset: Math.min(1, Math.max(0, s.offset)), color: s.color })).sort((a, b) => a.offset - b.offset);
    const id = this.id("g");
    let s = `<${end} id="${id}" ${open}`;
    for (const stop of stops) {
      const a = stop.color & 0xff;
      s += `<stop offset="${num(stop.offset, 5)}" stop-color="${hexColor(stop.color)}"${a !== 0xff ? ` stop-opacity="${factor(a / 255)}"` : ""}/>`;
    }
    this.defs += `${s}</${end}>\n`;
    return `url(#${id})`;
  }
  private definePattern(p: PaintRef): string {
    const pt = p.paint!;
    const info = this.res.images.get(p.obj!.images[pt.image]);
    if (!info) return "none";
    if (info.w === 0) {
      this.warn("a pattern of an image whose size cannot be told is not drawn");
      return "none";
    }
    const m = Array.from(pt.matrix) as Matrix;
    if (m.length !== 6 || !finite(m)) return "none";
    let w = info.w, h = info.h;
    if (pt.repeat === 2 || pt.repeat === 3) w = Math.max(w, UNREPEATED);
    if (pt.repeat === 1 || pt.repeat === 3) h = Math.max(h, UNREPEATED);
    const id = this.id("g");
    const tf = isIdentity(m) ? "" : ` patternTransform="matrix(${m.map(factor).join(" ")})"`;
    // the image first: it may define a symbol
    const el = this.imageElement(info, 0, 0, info.w, info.h, 3, "");
    this.defs += `<pattern id="${id}" patternUnits="userSpaceOnUse" width="${num(w, 3)}" height="${num(h, 3)}"${tf}>${el}</pattern>\n`;
    return `url(#${id})`;
  }

  /** The font-family of a font: the family of its embedded font, declared at its first use, before the families the document names. */
  private family(f: Font): string {
    if (f.kind === FontKind.EMBEDDED && f.hash) {
      let name = this.fonts.get(f.hash);
      if (name === undefined) {
        name = "";
        const url = this.res.fonts.get(f.hash);
        if (url !== undefined) {
          name = `bdf-${f.hash.slice(0, 10)}`;
          this.css += `@font-face{font-family:"${name}";src:url("${url}")}\n`;
        }
        this.fonts.set(f.hash, name);
      }
      if (name !== "") return f.family ? `${name},${f.family}` : name;
    }
    return f.family || "sans-serif";
  }

  // --- text ---

  /** A FILL_TEXT or STROKE_TEXT as a text element, fitted to its advance as the viewer scales a run. */
  private text(s: string, x: number, y: number, advance: number, stroke: boolean): void {
    const c = this.ctx, st = c.st;
    if (!(st.size > 0) || !Number.isFinite(st.size) || s === "" || !finite(st.m) || !Number.isFinite(x + y)) return;
    const [mode, blend] = this.mode(st.blend);
    const f: Font = st.font ?? { kind: FontKind.SYSTEM, family: "sans-serif", weight: 400, style: 0 };
    let attrs = `${transformOf(st.m)} font-family="${escape(this.family(f))}" font-size="${factor(st.size)}"`;
    if (f.weight !== 0 && f.weight !== 400) attrs += ` font-weight="${f.weight}"`;
    if (STYLE_NAMES[f.style]) attrs += ` font-style="${STYLE_NAMES[f.style]}"`;
    attrs += stroke ? ` fill="none"${this.strokeAttrs(st, mode)}` : this.fillAttrs(st, mode);
    const rtl = st.dir === 2;
    if (rtl) attrs += ` direction="rtl"`;
    // start and end are the sides of the direction
    if (st.align === 2) attrs += ` text-anchor="middle"`;
    else if ((st.align === 1 && !rtl) || (st.align === 0 && rtl) || st.align === 4) attrs += ` text-anchor="end"`;
    if (BASELINE_NAMES[st.baseline]) attrs += ` dominant-baseline="${BASELINE_NAMES[st.baseline]}"`;
    if (st.spacing !== 0 && Number.isFinite(st.spacing)) attrs += ` letter-spacing="${factor(st.spacing)}"`;
    const digits = digitsOf(st.m);
    let pos = "";
    if (x !== 0) pos += ` x="${num(x, digits)}"`;
    if (y !== 0) pos += ` y="${num(y, digits)}"`;
    if (advance > 0 && Number.isFinite(advance)) pos += ` textLength="${num(advance, digits)}" lengthAdjust="spacingAndGlyphs"`;
    const text = escape(s);
    // Canvas draws every space. The attribute is read on the text element itself: Chrome takes it from no element around it.
    if (text.startsWith(" ") || text.endsWith(" ") || text.includes("  ")) pos += ` xml:space="preserve"`;
    const body = `>${text}</text>`;
    if (mode === MODE_OVER && blend === "" && !shadowActive(st)) this.writeIn(c, st.clip, attrs, `<text${pos}${body}`);
    else this.compose(c, mode, blend, `<text${attrs}${pos}${body}`, true);
  }

  // --- images ---

  /** An image, or the part src (in image pixels) of it, in the rectangle (x, y, w, h). */
  private drawImage(i: number, src: [number, number, number, number] | undefined, x: number, y: number, w: number, h: number): void {
    const c = this.ctx, st = c.st;
    if (w === 0 || h === 0 || !finite(st.m) || !Number.isFinite(x + y + w + h)) return;
    const info = this.res.images.get(this.obj.images[i]);
    if (!info) return;
    // a rectangle of a negative size mirrors the image
    let m = st.m;
    if (w < 0 || h < 0) {
      let sx = 1, sy = 1;
      if (w < 0) [sx, x, w] = [-1, -x, -w];
      if (h < 0) [sy, y, h] = [-1, -y, -h];
      m = mul(m, [sx, 0, 0, sy, 0, 0]);
    }
    let attrs = "";
    if (st.alpha < 1) attrs += ` opacity="${factor(st.alpha)}"`;
    if (!st.smooth) attrs += ` style="image-rendering:pixelated"`;
    const digits = digitsOf(m);
    let el: string;
    if (!src) {
      el = this.imageElement(info, x, y, w, h, digits, transformOf(m) + attrs);
    } else {
      let [sx, sy, sw, sh] = src;
      if (sw < 0) [sx, sw] = [sx + sw, -sw];
      if (sh < 0) [sy, sh] = [sy + sh, -sh];
      if (!(sw > 0) || !(sh > 0) || !Number.isFinite(sx + sy + sw + sh)) return;
      if (info.w === 0) {
        this.warn("a part of an image whose size cannot be told is not drawn");
        return;
      }
      // the part is in the pixels of the image: a view box of them
      el = `<g${transformOf(m)}${attrs}><svg`;
      if (x !== 0) el += ` x="${num(x, digits)}"`;
      if (y !== 0) el += ` y="${num(y, digits)}"`;
      el += ` width="${num(w, digits)}" height="${num(h, digits)}" viewBox="${num(sx, 3)} ${num(sy, 3)} ${num(sw, 3)} ${num(sh, 3)}" preserveAspectRatio="none">` +
        `${this.imageElement(info, 0, 0, info.w, info.h, 3, "")}</svg></g>`;
    }
    const [mode, blend] = this.mode(st.blend);
    this.compose(c, mode, blend, el, false);
  }

  /** The element of an image stretched into a rectangle; the image is in the SVG once, as a symbol. */
  private imageElement(info: ImageInfo, x: number, y: number, w: number, h: number, digits: number, attrs: string): string {
    let box = "";
    if (x !== 0) box += ` x="${num(x, digits)}"`;
    if (y !== 0) box += ` y="${num(y, digits)}"`;
    box += ` width="${num(w, digits)}" height="${num(h, digits)}"`;
    if (info.w === 0) return `<image href="${escape(info.url)}"${box} preserveAspectRatio="none"${attrs}/>`;
    let id = this.symbols.get(info.hash);
    if (id === undefined) {
      id = this.id("i");
      this.symbols.set(info.hash, id);
      const size = ` width="${num(info.w, 3)}" height="${num(info.h, 3)}"`;
      this.defs += `<symbol id="${id}" viewBox="0 0 ${num(info.w, 3)} ${num(info.h, 3)}" preserveAspectRatio="none"><image href="${escape(info.url)}"${size} preserveAspectRatio="none"/></symbol>\n`;
    }
    return `<use href="#${id}"${box}${attrs}/>`;
  }

  // --- the instructions ---

  private skip(): boolean {
    return this.skipped > 0;
  }

  save(): void {
    const c = this.ctx;
    if (this.skip() || c.stack.length >= MAX_STATES) {
      if (!this.skip()) this.warn(`more than ${MAX_STATES} states are saved at a time: the rest are not kept`);
      this.unsaved.set(c, (this.unsaved.get(c) ?? 0) + 1);
      return;
    }
    c.stack.push({ ...c.st });
  }
  restore(): void {
    const c = this.ctx;
    const n = this.unsaved.get(c) ?? 0;
    if (n > 0) {
      this.unsaved.set(c, n - 1);
      return;
    }
    const st = c.stack.pop();
    if (st) c.st = st;
  }
  transform(a: number, b: number, c: number, d: number, e: number, f: number): void {
    if (!this.skip()) this.ctx.st.m = mul(this.ctx.st.m, [a, b, c, d, e, f]);
  }
  translate(x: number, y: number): void {
    if (!this.skip()) this.ctx.st.m = mul(this.ctx.st.m, [1, 0, 0, 1, x, y]);
  }
  scale(x: number, y: number): void {
    if (!this.skip()) this.ctx.st.m = mul(this.ctx.st.m, [x, 0, 0, y, 0, 0]);
  }
  clipPath(path: number, rule: number): void {
    if (this.skip()) return;
    const p = this.pathOf(path);
    if (p) this.clip(svgPathData(p, digitsOf(this.ctx.st.m)), rule);
  }
  clipRect(x: number, y: number, w: number, h: number): void {
    if (!this.skip()) this.clip(rectData(x, y, w, h, digitsOf(this.ctx.st.m)), 0);
  }
  fillColor(rgba: number): void {
    if (!this.skip()) this.ctx.st.fill = { color: rgba };
  }
  strokeColor(rgba: number): void {
    if (!this.skip()) this.ctx.st.stroke = { color: rgba };
  }
  private paint(i: number): PaintRef {
    const p = this.obj.paints[i];
    return p ? { color: 0, paint: p, obj: this.obj } : { color: 0 };
  }
  fillPaint(paint: number): void {
    if (!this.skip()) this.ctx.st.fill = this.paint(paint);
  }
  strokePaint(paint: number): void {
    if (!this.skip()) this.ctx.st.stroke = this.paint(paint);
  }
  line(width: number, cap: number, join: number, miter: number): void {
    if (this.skip()) return;
    const st = this.ctx.st;
    if (width > 0 && Number.isFinite(width)) st.width = width;
    st.cap = cap;
    st.join = join;
    if (miter > 0 && Number.isFinite(miter)) st.miter = miter;
  }
  dash(segments: Float32Array, offset: number): void {
    if (this.skip()) return;
    const st = this.ctx.st;
    // a list with a value that is no length is ignored, as Canvas ignores it
    if (segments.every((v) => v >= 0 && Number.isFinite(v))) st.dash = segments.some((v) => v > 0) ? segments : undefined;
    if (Number.isFinite(offset)) st.dashOffset = offset;
  }
  alpha(a: number): void {
    if (!this.skip() && a >= 0 && a <= 1) this.ctx.st.alpha = a;
  }
  blend(mode: number): void {
    if (!this.skip() && mode < BLEND_NAMES.length) this.ctx.st.blend = mode;
  }
  shadow(rgba: number, blur: number, dx: number, dy: number): void {
    if (this.skip()) return;
    const st = this.ctx.st;
    const m = st.m;
    let s = Math.sqrt(Math.abs(m[0] * m[3] - m[1] * m[2]));
    if (s === 0 || !Number.isFinite(s)) s = 1;
    st.shadow = rgba;
    // a value that is no number is ignored, as Canvas ignores it
    if (Number.isFinite(blur * s)) st.shadowBlur = Math.max(0, blur * s);
    if (Number.isFinite(dx * s)) st.shadowDX = dx * s;
    if (Number.isFinite(dy * s)) st.shadowDY = dy * s;
  }
  filter(css: string): void {
    if (!this.skip() && css !== "" && css !== "none") this.warn(`FILTER ${JSON.stringify(css)} is not applied`);
  }
  font(font: number, size: number): void {
    if (this.skip()) return;
    const f = this.obj.fonts[font];
    if (f) this.ctx.st.font = f;
    this.ctx.st.size = size;
  }
  textStyle(align: number, baseline: number, dir: number, letterSpacing: number): void {
    if (this.skip()) return;
    const st = this.ctx.st;
    st.align = align;
    st.baseline = baseline;
    st.dir = dir;
    st.spacing = letterSpacing;
  }
  fillRect(x: number, y: number, w: number, h: number): void {
    if (!this.skip()) this.rect(x, y, w, h, false);
  }
  strokeRect(x: number, y: number, w: number, h: number): void {
    if (!this.skip()) this.rect(x, y, w, h, true);
  }
  fillPath(path: number, rule: number): void {
    if (!this.skip()) this.outline(path, rule, false);
  }
  strokePath(path: number): void {
    if (!this.skip()) this.outline(path, 0, true);
  }
  fillPathAt(path: number, rule: number, x: number, y: number): void {
    if (!this.skip()) this.placed(rule, [{ path, x, y }]);
  }
  fillPathRun(rule: number, glyphs: Glyph[]): void {
    if (!this.skip()) this.placed(rule, glyphs);
  }
  clearRect(x: number, y: number, w: number, h: number): void {
    if (this.skip()) return;
    // what erases, whatever the state draws with
    const c = this.ctx, saved = c.st;
    c.st = { ...saved, fill: { color: 0x000000ff }, alpha: 1, blend: BLEND_DESTINATION_OUT, shadow: 0 };
    this.rect(x, y, w, h, false);
    c.st = saved;
  }
  fillText(text: string, x: number, y: number, advance: number): void {
    if (!this.skip()) this.text(text, x, y, advance, false);
  }
  strokeText(text: string, x: number, y: number, advance: number): void {
    if (!this.skip()) this.text(text, x, y, advance, true);
  }
  image(img: number, x: number, y: number, w: number, h: number): void {
    if (!this.skip()) this.drawImage(img, undefined, x, y, w, h);
  }
  imageSub(img: number, sx: number, sy: number, sw: number, sh: number, dx: number, dy: number, dw: number, dh: number): void {
    if (!this.skip()) this.drawImage(img, [sx, sy, sw, sh], dx, dy, dw, dh);
  }
  smoothing(enabled: boolean, _quality: number): void {
    if (!this.skip()) this.ctx.st.smooth = enabled;
  }
  use(obj: number): void {
    this.useAt(obj, 0, 0);
  }
  useAt(obj: number, x: number, y: number): void {
    if (this.skip()) return;
    const hash = this.obj.objects[obj];
    if (hash === undefined) return;
    let child: ObjectPart;
    try {
      child = this.doc.objectSync(hash);
    } catch {
      return;
    }
    if (this.depth >= 64) {
      this.warn("objects nested too deep");
      return;
    }
    const again = this.limits.enter(hash);
    const c = this.ctx;
    const saved = { ...c.st }, depth = c.stack.length, unsaved = this.unsaved.get(c) ?? 0;
    if (x !== 0 || y !== 0) c.st.m = mul(c.st.m, [1, 0, 0, 1, x, y]);
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
      // USE is wrapped in an implicit save and restore, of the context it draws on (which is gone if the object ended its group)
      if (!c.done) {
        c.stack.length = Math.min(depth, c.stack.length);
        this.unsaved.set(c, unsaved);
        c.st = saved;
      }
    }
  }
  /** Whether one more group or soft mask may be open; one past the limit is left out with what is in it. */
  private layer(): boolean {
    if (this.skip() || this.groups.length + this.masks.length >= MAX_LAYERS) {
      if (!this.skip()) this.warn(`more than ${MAX_LAYERS} groups and soft masks are open: what is in the rest of them is not drawn`);
      this.skipped++;
      return false;
    }
    return true;
  }
  groupBegin(alpha: number, blend: number, _x: number, _y: number, _w: number, _h: number): void {
    if (!this.layer()) return;
    const c = this.ctx;
    this.flush(c);
    const st = c.st;
    const gs = initialState(st.m, 0);
    gs.fill = st.fill;
    gs.stroke = st.stroke;
    gs.width = st.width;
    gs.font = st.font;
    gs.size = st.size;
    const g: Group = { parent: c, ctx: new Context(gs), alpha, blend, mask: "" };
    this.groups.push(g);
    this.ctx = g.ctx;
  }
  groupEnd(): void {
    if (this.skipped > 0) {
      this.skipped--;
      return;
    }
    const g = this.groups.pop();
    if (!g) return;
    // masks begun inside the group end with it
    while (this.masks.length > 0 && this.masks[this.masks.length - 1].group === g) this.masks.pop()!.ctx.done = true;
    this.closeAll(g.ctx);
    this.settle(g.ctx);
    g.ctx.done = true;
    const c = g.parent;
    this.ctx = c;
    if (g.ctx.buf === "" || !(g.alpha > 0)) return;
    let attrs = "";
    if (g.alpha < 1) attrs += ` opacity="${factor(g.alpha)}"`;
    if (g.mask !== "") attrs += ` mask="url(#${g.mask})"`;
    const [mode, blend] = this.mode(g.blend);
    if (mode === MODE_OVER && blend === "") {
      // a group that changes nothing is its content, unless that blends: with the group, not with what is under it
      this.flush(c);
      let open = "", end = "";
      if (attrs !== "" || g.ctx.blended) {
        if (g.ctx.blended) attrs += ` style="isolation:isolate"`;
        open = `<g${attrs}>\n`;
        end = "</g>\n";
      }
      this.within(c, c.st.clip);
      this.closeAttrs(c);
      c.buf += open + g.ctx.buf + end;
      return;
    }
    // a group is drawn without the shadow of the state
    const shadow = c.st.shadow;
    c.st.shadow = 0;
    this.compose(c, mode, blend, `<g${attrs}>\n${g.ctx.buf}</g>`, false);
    c.st.shadow = shadow;
  }
  maskBegin(kind: number, backdrop: number, transfer: Uint8Array): void {
    if (!this.layer()) return;
    const c = this.ctx;
    this.flush(c);
    const mc = new Context(initialState(c.st.m, 0));
    this.masks.push({ parent: c, ctx: mc, kind, backdrop, transfer: transfer.length === 256 ? transfer : undefined, group: this.groups[this.groups.length - 1] });
    this.ctx = mc;
  }
  maskEnd(): void {
    if (this.skipped > 0) {
      this.skipped--;
      return;
    }
    const m = this.masks.pop();
    if (!m) return;
    this.closeAll(m.ctx);
    this.settle(m.ctx);
    m.ctx.done = true;
    this.ctx = m.parent;
    const g = m.group;
    if (!g || this.groups[this.groups.length - 1] !== g) return;
    const lum = m.kind === MaskKind.LUMINOSITY;
    let content = m.ctx.buf;
    if (lum) content = `<path d="${rectData(this.x, this.y, this.w, this.h, PAGE_DIGITS)}" fill="${hexColor(m.backdrop)}"/>\n${content}`;
    if (m.transfer) {
      // the transfer function maps the value of the mask, which for a mask of luminance is made a gray of the weights of BDF first
      const t = Array.from(m.transfer, (v) => num(v / 255, 4)).join(" ");
      let prim = `<feComponentTransfer><feFuncA type="table" tableValues="${t}"/></feComponentTransfer>`;
      if (lum) {
        const table = ` type="table" tableValues="${t}"/>`;
        prim = `<feColorMatrix values=".3 .59 .11 0 0 .3 .59 .11 0 0 .3 .59 .11 0 0 0 0 0 1 0"/><feComponentTransfer><feFuncR${table}<feFuncG${table}<feFuncB${table}</feComponentTransfer>`;
      }
      const id = this.id("f");
      this.defs += `<filter id="${id}" filterUnits="userSpaceOnUse"${this.region()} color-interpolation-filters="sRGB">${prim}</filter>\n`;
      content = `<g filter="url(#${id})">\n${content}</g>\n`;
    }
    const id = this.id("m");
    this.defs += `<mask id="${id}" maskUnits="userSpaceOnUse"${this.region()}${lum ? "" : ` style="mask-type:alpha"`}>\n${content}</mask>\n`;
    g.mask = id;
  }
  link(x: number, y: number, w: number, h: number, url: string): void {
    if (this.skip() || this.masks.length > 0) return;
    // the links of a file that stands alone: those that leave it
    if (!/^(https?:\/\/|mailto:)/i.test(url)) return;
    const m = this.ctx.st.m;
    const data = finite(m) ? rectData(x, y, w, h, digitsOf(m)) : undefined;
    if (data !== undefined) this.links += `<a href="${escape(url)}"><path d="${data}"${transformOf(m)} fill="none" pointer-events="all"/></a>\n`;
  }
  mark(): void {}
  ext(): void {}

  // --- pages ---

  /** Draw a top-level object (a page layer or a sheet tile) from the initial state with a transform and a clip. */
  drawTop(hash: Hash, m: Matrix, clip: number): void {
    let o: ObjectPart;
    try {
      o = this.doc.objectSync(hash);
    } catch {
      return;
    }
    const top = this.ctx;
    top.st = initialState(m, clip);
    top.stack = [];
    this.unsaved.clear();
    this.reused = this.limits.enter(hash);
    this.obj = o;
    this.depth = 0;
    try {
      walk(o, this, this.reused ? this.limits.count : undefined);
    } catch (e) {
      // what was drawn before the object ended stays: a file to save is better short of an object than none
      this.warn(e instanceof Error ? e.message : String(e));
    }
    // unwind groups and masks left open by a malformed stream
    this.masks.length = 0;
    this.skipped = 0;
    while (this.groups.length > 0) this.groupEnd();
    this.ctx = top;
    this.flush(top);
  }

  /** Fill a rectangle with the background of the pages. */
  background(x: number, y: number, w: number, h: number): void {
    const bg = this.opts.background === undefined ? "#fff" : this.opts.background;
    if (bg === null || !/^#[0-9a-fA-F]{3,8}$/.test(bg)) return;
    this.write(this.ctx, this.ctx.st.clip, `<path d="${rectData(x, y, w, h, PAGE_DIGITS)}" fill="${bg}"/>`);
  }

  /** The clip path of a rectangle under m. */
  clipRectAt(m: Matrix, x: number, y: number, w: number, h: number): number {
    this.ctx.st = initialState(m, 0);
    this.clip(rectData(x, y, w, h, digitsOf(m)), 0);
    return this.ctx.st.clip;
  }

  /** An element over the background, under all the objects draw. */
  element(el: string): void {
    this.write(this.ctx, 0, el);
  }

  finish(): string {
    const root = this.ctx;
    this.closeAll(root);
    this.settle(root);
    let s = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="${num(this.x, PAGE_DIGITS)} ${num(this.y, PAGE_DIGITS)} ${num(this.w, PAGE_DIGITS)} ${num(this.h, PAGE_DIGITS)}"` +
      ` width="${num(this.w, PAGE_DIGITS)}pt" height="${num(this.h, PAGE_DIGITS)}pt"${root.blended ? ` style="isolation:isolate"` : ""}>\n`;
    if (this.css !== "" || this.defs !== "") {
      s += "<defs>\n";
      if (this.css !== "") s += `<style>\n${this.css}</style>\n`;
      s += `${this.defs}</defs>\n`;
    }
    return `${s}${root.buf}${this.links}</svg>\n`;
  }
}

/**
 * Writes the pages, the continuous layout and the sheets of a document as
 * SVG files that stand alone: their fonts and images are in them.
 */
export class SvgExporter {
  private fonts = new Map<Hash, string>();
  private images = new Map<Hash, ImageInfo>();
  private paths = new Map<Hash, PathData[]>();

  constructor(readonly doc: BdfDocument, readonly options: SvgOptions = {}) {}

  /** Load the objects of a region and what they refer to. */
  private async prepare(hashes: Hash[]): Promise<Loaded> {
    await Promise.all(hashes.map((h) => this.doc.ensure(h, async (e, bytes) => {
      switch (e.t) {
        case "font":
          if (!this.fonts.has(e.h)) this.fonts.set(e.h, `data:${fontType(bytes)};base64,${base64(bytes)}`);
          break;
        case "img": {
          if (this.images.has(e.h)) break;
          const type = imageType(bytes);
          if (!type) break;
          const size = shownSize(bytes, type);
          this.images.set(e.h, { hash: e.h, url: `data:${type};base64,${base64(bytes)}`, w: size?.width ?? 0, h: size?.height ?? 0 });
          break;
        }
        case "path":
          if (!this.paths.has(e.h)) this.paths.set(e.h, await this.doc.pathCollection(e.h));
          break;
      }
    })));
    return { fonts: this.fonts, images: this.images, paths: this.paths };
  }

  private sink(res: Loaded, r: Rect): SvgSink {
    return new SvgSink(this.doc, res, this.options, r.x, r.y, r.w, r.h);
  }

  /** A page of a fixed or flow view, or the part region of it (in page units). */
  async page(page: Page, region: Rect = { x: 0, y: 0, w: page.w, h: page.h }): Promise<SvgResult> {
    const res = await this.prepare(page.layers.map((l) => l.obj));
    const d = this.sink(res, region);
    // the SVG shows no more than the region: a clip only where that is more than the page
    let clip = 0;
    if (region.x < 0 || region.y < 0 || region.x + region.w > page.w || region.y + region.h > page.h) clip = d.clipRectAt(IDENTITY, 0, 0, page.w, page.h);
    d.background(0, 0, page.w, page.h);
    for (const l of page.layers) d.drawTop(l.obj, IDENTITY, clip);
    return { svg: d.finish(), warnings: d.warnings };
  }

  /**
   * The part of the continuous layout of a view inside viewport (units from
   * its top): the bodies of its pages stacked, as the viewer scrolls them
   * (the strips of a scroll view, without gaps).
   */
  async continuous(view: View, viewport: Rect): Promise<SvgResult> {
    const pages = view.pages ?? [];
    const gap = view.kind === "scroll" ? 0 : view.continuous?.gap ?? 0;
    const shown: [Page, number][] = [];
    let y = 0;
    for (const p of pages) {
      const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
      if (y < viewport.y + viewport.h && y + b.h > viewport.y) shown.push([p, y]);
      y += b.h + gap;
    }
    const res = await this.prepare(shown.flatMap(([p]) => p.layers.map((l) => l.obj)));
    const d = this.sink(res, viewport);
    for (const [p, top] of shown) {
      const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
      const clip = d.clipRectAt([1, 0, 0, 1, 0, top], 0, 0, b.w, b.h);
      d.background(0, top, b.w, b.h);
      for (const l of p.layers) {
        if (l.role === "body" || l.role === "annotation") d.drawTop(l.obj, [1, 0, 0, 1, -b.x, top - b.y], clip);
      }
    }
    return { svg: d.finish(), warnings: d.warnings };
  }

  /** The region of a sheet inside viewport (units): the gridlines when the view asks for them, and the tiles. */
  async sheet(view: View, viewport: Rect): Promise<SvgResult> {
    const tile = tileSize(view);
    const range = {
      tx0: Math.floor(viewport.x / tile), ty0: Math.floor(viewport.y / tile),
      tx1: Math.floor((viewport.x + viewport.w - 1e-6) / tile), ty1: Math.floor((viewport.y + viewport.h - 1e-6) / tile),
    };
    const keys = tilesIn(view, [range]);
    const res = await this.prepare(keys.map(([, , h]) => h));
    const d = this.sink(res, viewport);
    d.background(viewport.x, viewport.y, viewport.w, viewport.h);
    if (view.gridlines) {
      const lines = gridlines(view, viewport);
      if (lines !== "") d.element(`<path d="${lines}" fill="none" stroke="#d9d9d9" vector-effect="non-scaling-stroke" shape-rendering="crispEdges"/>`);
    }
    for (const [tx, ty, h] of keys) {
      const m: Matrix = [1, 0, 0, 1, tx * tile, ty * tile];
      d.drawTop(h, m, d.clipRectAt(m, 0, 0, tile, tile));
    }
    return { svg: d.finish(), warnings: d.warnings };
  }
}

/** The lines of a sheet region drawn at most: a document decides how many rows and columns a sheet has. */
const MAX_GRID_LINES = 1 << 16;

/** The path of the lines between the cells of a sheet inside a region. */
function gridlines(view: View, vp: Rect): string {
  const sum = (runs?: [number, number][]) => (runs ?? []).reduce((acc, [n, s]) => acc + n * s, 0);
  const x1 = Math.min(vp.x + vp.w, sum(view.cols)), y1 = Math.min(vp.y + vp.h, sum(view.rows));
  let out = "";
  const walkRuns = (runs: [number, number][] | undefined, from: number, to: number, line: (pos: number) => void) => {
    let pos = 0, lines = 0, last = NaN;
    for (const [n, size] of runs ?? []) {
      if (!(n >= 1)) continue;
      for (let i = 0; i < Math.floor(n); i++) {
        pos += size;
        if (pos > to) return;
        if (!(pos >= from) || pos === last) {
          if (size === 0) break; // so are the rest of the run
          continue;
        }
        if (++lines > MAX_GRID_LINES) return;
        line(pos);
        last = pos;
      }
    }
  };
  walkRuns(view.cols, vp.x, x1, (x) => { out += `M${num(x, PAGE_DIGITS)} ${num(vp.y, PAGE_DIGITS)}V${num(y1, PAGE_DIGITS)}`; });
  walkRuns(view.rows, vp.y, y1, (y) => { out += `M${num(vp.x, PAGE_DIGITS)} ${num(y, PAGE_DIGITS)}H${num(x1, PAGE_DIGITS)}`; });
  return out;
}

/** The rectangle of a sheet that has content, from A1: as far as its tiles reach, within its rows and columns. */
export function sheetBounds(view: View): Rect {
  const tile = tileSize(view);
  let w = 0, h = 0;
  for (const key of Object.keys(view.tiles ?? {})) {
    const m = /^(\d+),(\d+)$/.exec(key);
    if (!m || `${Number(m[1])},${Number(m[2])}` !== key) continue;
    w = Math.max(w, (Number(m[1]) + 1) * tile);
    h = Math.max(h, (Number(m[2]) + 1) * tile);
  }
  const sum = (runs?: [number, number][]) => (runs ?? []).reduce((acc, [n, s]) => acc + n * s, 0);
  const sw = sum(view.cols), sh = sum(view.rows);
  if (sw > 0 && sh > 0) [w, h] = [Math.min(w, sw), Math.min(h, sh)];
  return { x: 0, y: 0, w, h };
}

/**
 * What saving a view as one file covers: page of a fixed or flow view (its
 * index), all the strips of a scroll view stacked, the cells of a sheet.
 * target is what BdfWorkerClient.svg and image take; undefined for a view
 * with nothing to save.
 */
export function saveTarget(view: View, page: number): { target: number | Rect; width: number; height: number } | undefined {
  if (view.kind === "sheet") {
    const b = sheetBounds(view);
    return b.w > 0 && b.h > 0 ? { target: b, width: b.w, height: b.h } : undefined;
  }
  const pages = view.pages ?? [];
  if (view.kind === "scroll") {
    let width = 0, height = 0;
    for (const p of pages) {
      const b = p.body ?? p;
      width = Math.max(width, b.w);
      height += b.h;
    }
    return width > 0 && height > 0 ? { target: { x: 0, y: 0, w: width, h: height }, width, height } : undefined;
  }
  const p = pages[page];
  return p && p.w > 0 && p.h > 0 ? { target: page, width: p.w, height: p.h } : undefined;
}

/**
 * The scale an image of width × height units is drawn at when scale is
 * asked for: no more than a canvas of browsers holds (a long scroll view is
 * drawn smaller).
 */
export function imageScale(width: number, height: number, scale: number): number {
  const side = 16384, pixels = 1 << 26;
  return Math.min(scale, side / Math.max(width, height), Math.sqrt(pixels / (width * height)));
}

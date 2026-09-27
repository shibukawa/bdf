// The book layouts of the demo viewer. A view is shown one page at a time,
// or as two facing pages with the first page alone like a book's cover,
// fitted to the stage. Pages are turned with the keys, the page buttons, or
// by pulling a corner or the outer edge of a page.
//
// The pages are the viewer's usual page elements (a bitmap under a text
// layer), so text is selected, copied and searched as in the scrolled
// layout. Only the corners and the outer edges are handles, and those sit
// mostly in the page margins. While a page turns, a WebGL canvas takes the
// place of the page elements and curls the page (pageflip.ts). Its textures
// are the bitmaps of the spreads before and after, rendered ahead of time.

import type { Page, TextContent } from "@bdf/core";
import { FlipRenderer, foldAt, reach, type Leaf, type Point, type Rect } from "./pageflip.js";

export type BookLayout = "single" | "spread";

export interface BookOptions {
  pages: Page[];
  layout: BookLayout;
  /** Pages turn from left to right, as in a book read right to left. */
  rtl: boolean;
  /** Times the size that fits the stage. */
  zoom: number;
  /** The page to open at. */
  start: number;
  /** Accessible name of a page element. */
  label(index: number): string;
  /** A page's bitmap, scale device px per unit. */
  bitmap(index: number, scale: number): Promise<ImageBitmap>;
  content(index: number): Promise<TextContent>;
  /** What goes over a page's bitmap (its text layer), scale CSS px per unit. */
  layers(index: number, content: TextContent, scale: number): HTMLElement[];
  /** A page still to come (a PDF converted a page at a time). */
  pending(index: number): boolean;
  failed(index: number): boolean;
  /** The pages shown have changed. */
  onTurn(pages: number[]): void;
  onError(e: unknown): void;
}

/** The pages of a spread before and beyond the spine; null: none. */
interface Spread { first: number | null; second: number | null }

type Slot = keyof Spread;
type Dir = "next" | "prev";
type Part = "top" | "bottom" | "edge";

/**
 * A page turning from the spread shown to the one before or after it. It is
 * worked out in its own frame, where the sheet lies flat beyond the spine and
 * turns over to the other side. For the next page that is the forward frame.
 * Two pages back, the first page turns back, which is the same seen in a
 * mirror. A page alone going back is the page before turning in reverse.
 * p is where the sheet's corner is now: rest while from is shown, done once
 * the page has turned to `to`.
 */
interface Turn {
  from: number;
  to: number;
  dir: Dir;
  /** The frame of the turn is the mirror image of the screen. */
  mirror: boolean;
  sheet: Rect;
  front: number;
  back: number | null;
  before?: Leaf;
  under?: Leaf;
  corner: Point;
  /** Pulled by the edge: the fold stays upright. */
  edge: boolean;
  rest: Point;
  done: Point;
  p: Point;
  /** Peeled back while the mouse is over a corner. */
  hover?: boolean;
  drag?: Drag;
  anim?: Anim;
}

interface Drag {
  id: number;
  zone: HTMLElement;
  start: Point;
  /** Where the corner was when the drag started. */
  anchor: Point;
  /** How far the corner moves per px of the pointer: a page alone turns over within its own width. */
  kx: number;
  moved: boolean;
  trail: { x: number; t: number }[];
}

interface Anim { from: Point; to: Point; t0: number; dur: number; lift: number; ease: (t: number) => number; then: () => void }

/** Space around the pages, CSS px. */
const PAD = 24;
/** How far a press moves before it is a drag, CSS px. */
const TAP = 4;
/** Half the width of an edge handle, CSS px. */
const EDGE = 12;
const easeOut = (t: number) => 1 - (1 - t) ** 3;
const easeInOut = (t: number) => (t < 0.5 ? 4 * t ** 3 : 1 - (-2 * t + 2) ** 3 / 2);
const smoothstep = (a: number, b: number, x: number) => {
  const t = Math.min(1, Math.max(0, (x - a) / (b - a)));
  return t * t * (3 - 2 * t);
};
const dpr = () => window.devicePixelRatio || 1;

/** The spreads of n pages: one each, or the first alone beyond the spine and then pairs. */
function spreadsOf(n: number, layout: BookLayout): Spread[] {
  if (layout === "single") return Array.from({ length: n }, (_, i) => ({ first: null, second: i }));
  const out: Spread[] = [{ first: null, second: n > 0 ? 0 : null }];
  for (let i = 1; i < n; i += 2) out.push({ first: i, second: i + 1 < n ? i + 1 : null });
  return out;
}

/**
 * A handle (a box on the screen) drawn back from the boxes of text it covers:
 * outward (away 1: to the right), or for a corner up (top) or down (bottom),
 * whichever leaves more of it. Text lies within the page, which the handle
 * sticks out of, so something is left.
 */
function keepOff(z: Rect, boxes: Rect[], away: 1 | -1, part: Part): Rect {
  const hit = boxes.filter((b) => b.x < z.x + z.w && b.x + b.w > z.x && b.y < z.y + z.h && b.y + b.h > z.y);
  if (!hit.length) return z;
  const GAP = 4;
  let across: Rect;
  if (away > 0) {
    const x0 = Math.max(...hit.map((b) => b.x + b.w)) + GAP;
    across = { ...z, x: x0, w: z.x + z.w - x0 };
  } else {
    across = { ...z, w: Math.min(...hit.map((b) => b.x)) - GAP - z.x };
  }
  if (part === "edge") return across;
  let along: Rect;
  if (part === "top") {
    along = { ...z, h: Math.min(...hit.map((b) => b.y)) - GAP - z.y };
  } else {
    const y0 = Math.max(...hit.map((b) => b.y + b.h)) + GAP;
    along = { ...z, y: y0, h: z.y + z.h - y0 };
  }
  const area = (r: Rect) => Math.max(0, r.w) * Math.max(0, r.h);
  return area(along) > area(across) ? along : across;
}

export class Book {
  readonly el = document.createElement("div");
  /** CSS px per unit. */
  scale = 0;
  private readonly spreads: Spread[];
  private at: number;
  private readonly leaves = document.createElement("div");
  private readonly zones = document.createElement("div");
  private readonly canvas = document.createElement("canvas");
  /** null: no WebGL 2, pages change without turning. */
  private renderer: FlipRenderer | null | undefined;
  private readonly bitmaps = new Map<number, ImageBitmap>();
  private readonly texts = new Map<number, TextContent>();
  private readonly asked = new Set<number>();
  private readonly askedText = new Set<number>();
  /** Bumped when the scale changes: bitmaps rendered before are dropped. */
  private epoch = 0;
  private spine = 0;
  private top = 0;
  private slotH = 0;
  private width = 0;
  private height = 0;
  /** The pages are larger than the stage (zoomed in), which scrolls: no turning. */
  private overflow = false;
  private stageSize = "";
  private turn: Turn | undefined;
  private frame = 0;
  private resizeTimer: ReturnType<typeof setTimeout> | undefined;
  private readonly observer: ResizeObserver;
  /** Text came in under the handles during a turn: place them again when it ends. */
  private zonesStale = false;
  /** A press on a handle while pages change without turning. */
  private click: { id: number; dir: Dir } | undefined;
  /** A page to focus once it is shown (the target of a link). */
  private focusPage = -1;
  private dead = false;

  constructor(private readonly stage: HTMLElement, private readonly o: BookOptions) {
    this.spreads = spreadsOf(o.pages.length, o.layout);
    this.at = this.spreadOf(o.start);
    this.el.className = "book";
    this.leaves.className = "leaves";
    this.zones.className = "zones";
    this.canvas.setAttribute("aria-hidden", "true");
    this.el.append(this.leaves, this.canvas, this.zones);
    stage.appendChild(this.el);
    this.layout();
    this.showSpread();
    stage.addEventListener("keydown", this.onKey);
    // a selection dragged over a handle goes on as over the margin around it
    this.leaves.addEventListener("pointerdown", (e) => {
      if (e.button !== 0) return;
      this.el.classList.add("selecting");
      const end = () => {
        this.el.classList.remove("selecting");
        removeEventListener("pointerup", end, true);
        removeEventListener("pointercancel", end, true);
      };
      addEventListener("pointerup", end, true);
      addEventListener("pointercancel", end, true);
    });
    this.observer = new ResizeObserver(() => {
      clearTimeout(this.resizeTimer);
      this.resizeTimer = setTimeout(() => this.relayout(), 100);
    });
    this.observer.observe(stage);
  }

  /** The pages shown. */
  get pages(): number[] {
    return this.pagesOf(this.at);
  }

  /** Whether there is a spread before (-1) or after (+1) the one shown. */
  canTurn(delta: number): boolean {
    const k = this.at + delta;
    return k >= 0 && k < this.spreads.length;
  }

  /** Turn to the spread before (-1) or after (+1). */
  turnBy(delta: 1 | -1) {
    if (this.turn?.drag) return;
    this.finish();
    if (!this.canTurn(delta)) return;
    const t = this.animates() ? this.begin(delta > 0 ? "next" : "prev", "bottom", 0) : undefined;
    if (!t) return this.jump(this.at + delta);
    this.animate(t.done, 650, this.lift(t), easeInOut, () => this.land(true));
  }

  /** Show the spread of a page: turning to it when it is next to this one and animate is set. */
  go(page: number, animate = false, focus = false) {
    const k = this.spreadOf(page);
    if (focus) this.focusPage = page;
    if (k === this.at && !this.turn) return this.focusShown();
    if (animate && Math.abs(k - this.at) === 1 && !this.turn?.drag) return this.turnBy(k > this.at ? 1 : -1);
    this.jump(k);
  }

  /** A page of a stream came in (or failed). */
  arrived(index: number) {
    const el = this.pageEl(index);
    if (el) {
      el.removeAttribute("aria-busy");
      if (this.o.failed(index)) this.fail(el);
    }
    if (this.near(index)) {
      this.want(index);
      this.wantText(index);
    }
  }

  destroy() {
    this.dead = true;
    cancelAnimationFrame(this.frame);
    clearTimeout(this.resizeTimer);
    this.observer.disconnect();
    this.stage.removeEventListener("keydown", this.onKey);
    for (const b of this.bitmaps.values()) b.close();
    this.bitmaps.clear();
    this.renderer?.destroy();
    this.el.remove();
  }

  private spreadOf(page: number): number {
    const i = Math.min(Math.max(0, page), this.o.pages.length - 1);
    return this.o.layout === "single" ? Math.max(0, i) : Math.floor((i + 1) / 2);
  }

  private pagesOf(k: number): number[] {
    const s = this.spreads[k];
    return s ? [s.first, s.second].filter((i): i is number => i !== null) : [];
  }

  /** Whether a page belongs to the spread shown or those next to it. */
  private near(index: number): boolean {
    const k = this.spreadOf(index);
    return Math.abs(k - this.at) <= 1;
  }

  /** Fit the pages to the stage; true when the scale changed. */
  private layout(): boolean {
    const pages = this.o.pages;
    const maxW = pages.reduce((m, p) => Math.max(m, p.w), 1);
    const maxH = pages.reduce((m, p) => Math.max(m, p.h), 1);
    const cols = this.o.layout === "spread" ? 2 : 1;
    // the stage's outer size: the same with or without scroll bars
    const ow = this.stage.offsetWidth, oh = this.stage.offsetHeight;
    this.stageSize = `${ow}x${oh}`;
    const fit = Math.max(0.02, Math.min((ow - 2 * PAD) / (cols * maxW), (oh - 2 * PAD) / maxH));
    const scale = fit * this.o.zoom;
    const cw = this.stage.clientWidth, ch = this.stage.clientHeight;
    this.width = Math.max(cw, Math.ceil(cols * maxW * scale + 2 * PAD));
    this.height = Math.max(ch, Math.ceil(maxH * scale + 2 * PAD));
    this.overflow = this.width > cw || this.height > ch;
    this.slotH = maxH * scale;
    this.top = (this.height - this.slotH) / 2;
    // a page alone lies beyond the spine: the spine is its edge on the side it turns to
    this.spine = cols === 2 ? this.width / 2 : this.width / 2 + ((this.o.rtl ? 1 : -1) * maxW * scale) / 2;
    this.el.style.width = `${this.width}px`;
    this.el.style.height = `${this.height}px`;
    this.renderer?.resize(this.width, this.height, dpr());
    const changed = scale !== this.scale;
    this.scale = scale;
    return changed;
  }

  private relayout() {
    if (this.dead || `${this.stage.offsetWidth}x${this.stage.offsetHeight}` === this.stageSize) return;
    this.finish();
    if (this.layout()) {
      this.epoch++;
      this.asked.clear();
      for (const b of this.bitmaps.values()) b.close();
      this.bitmaps.clear();
      this.renderer?.forgetAll();
    }
    this.showSpread();
  }

  /** A page's place in the forward frame (the sheet turns from x >= spine to the other side). */
  private rectOf(page: number, slot: Slot): Rect {
    const p = this.o.pages[page];
    const w = p.w * this.scale, h = p.h * this.scale;
    return { x: slot === "first" ? this.spine - w : this.spine, y: this.top + (this.slotH - h) / 2, w, h };
  }

  /** The mirror image at the spine. */
  private flip(r: Rect): Rect {
    return { ...r, x: 2 * this.spine - r.x - r.w };
  }

  private mirror(p: Point): Point {
    return { x: 2 * this.spine - p.x, y: p.y };
  }

  /** From the forward frame to the book element and back (the mirror image when right to left). */
  private screen(r: Rect): Rect {
    return this.o.rtl ? this.flip(r) : r;
  }

  /** A pointer's place in the frame of a turn (mirrored or not). */
  private point(e: PointerEvent, mirror: boolean): Point {
    const r = this.el.getBoundingClientRect();
    const x = e.clientX - r.left, y = e.clientY - r.top;
    return { x: mirror ? 2 * this.spine - x : x, y };
  }

  private pageEl(index: number): HTMLDivElement | null {
    return this.leaves.querySelector<HTMLDivElement>(`.page[data-index="${index}"]`);
  }

  /** Lay out the page elements of the spread shown, and its handles. */
  private showSpread() {
    const s = this.spreads[this.at];
    const hadFocus = this.leaves.contains(document.activeElement);
    // bitmaps of the scale before stand in until the new ones come
    const old = new Map<number, HTMLCanvasElement>();
    for (const c of this.leaves.querySelectorAll<HTMLCanvasElement>(".page > canvas")) old.set(Number((c.parentElement as HTMLElement).dataset.index), c);
    this.leaves.replaceChildren();
    for (const slot of ["first", "second"] as const) {
      const i = s[slot];
      if (i === null) continue;
      const r = this.screen(this.rectOf(i, slot));
      const el = document.createElement("div");
      el.className = "page";
      el.dataset.index = String(i);
      el.setAttribute("role", "group");
      el.setAttribute("aria-label", this.o.label(i));
      el.tabIndex = -1; // target of page links
      el.style.cssText = `left: ${r.x}px; top: ${r.y}px; width: ${r.w}px; height: ${r.h}px`;
      this.leaves.appendChild(el);
      if (this.o.failed(i)) {
        this.fail(el);
        continue;
      }
      if (this.o.pending(i)) el.setAttribute("aria-busy", "true");
      const bmp = this.bitmaps.get(i);
      const stale = old.get(i);
      if (bmp) this.paint(el, bmp);
      else if (stale) {
        stale.style.width = `${r.w}px`;
        stale.style.height = `${r.h}px`;
        el.prepend(stale);
      }
      const c = this.texts.get(i);
      if (c) el.append(...this.o.layers(i, c, this.scale));
    }
    this.placeZones();
    this.prefetch();
    // the focus stays on the pages, where the keys turn them
    if (hadFocus && this.focusPage < 0) this.focusPage = this.pages[0] ?? -1;
    this.focusShown();
  }

  private focusShown() {
    if (this.focusPage < 0) return;
    const el = this.pageEl(this.focusPage);
    if (!el) return;
    this.focusPage = -1;
    el.focus({ preventScroll: true });
  }

  private fail(el: HTMLElement) {
    const p = document.createElement("p");
    p.className = "pageError";
    p.textContent = "This page could not be converted.";
    el.append(p);
  }

  private paint(el: HTMLElement, bmp: ImageBitmap) {
    const i = Number(el.dataset.index);
    const p = this.o.pages[i];
    const canvas = document.createElement("canvas");
    canvas.width = bmp.width;
    canvas.height = bmp.height;
    canvas.style.width = `${p.w * this.scale}px`;
    canvas.style.height = `${p.h * this.scale}px`;
    canvas.setAttribute("aria-hidden", "true");
    canvas.getContext("2d")!.drawImage(bmp, 0, 0);
    el.querySelector(":scope > canvas")?.remove();
    el.prepend(canvas);
  }

  /** Render the spreads next to the one shown ahead of time, and drop those further away. */
  private prefetch() {
    const order = [this.at, this.at + 1, this.at - 1];
    for (const k of order) for (const i of this.pagesOf(k)) this.want(i);
    for (const k of order) for (const i of this.pagesOf(k)) this.wantText(i);
    const keep = new Set([-2, -1, 0, 1, 2].flatMap((d) => this.pagesOf(this.at + d)));
    for (const [i, bmp] of this.bitmaps) {
      if (keep.has(i)) continue;
      bmp.close();
      this.bitmaps.delete(i);
      this.renderer?.forget(i);
    }
    for (const i of [...this.texts.keys()]) if (!keep.has(i)) this.texts.delete(i);
  }

  private want(i: number) {
    if (this.bitmaps.has(i) || this.asked.has(i) || this.o.pending(i) || this.o.failed(i)) return;
    this.asked.add(i);
    const epoch = this.epoch;
    this.o.bitmap(i, this.scale * dpr()).then((bmp) => {
      if (this.dead || epoch !== this.epoch) return bmp.close();
      this.asked.delete(i);
      this.bitmaps.set(i, bmp);
      const el = this.pageEl(i);
      if (el) this.paint(el, bmp);
      if (this.turn && this.renderer && this.turnPages(this.turn).includes(i)) {
        this.renderer.upload(i, bmp);
        this.schedule();
      }
    }, (e) => {
      if (this.dead || epoch !== this.epoch) return;
      this.asked.delete(i);
      this.o.onError(e);
    });
  }

  private wantText(i: number) {
    if (this.texts.has(i) || this.askedText.has(i) || this.o.pending(i) || this.o.failed(i)) return;
    this.askedText.add(i);
    this.o.content(i).then((c) => {
      this.askedText.delete(i);
      if (this.dead) return;
      this.texts.set(i, c);
      const el = this.pageEl(i);
      if (!el || el.querySelector(".bdfTextLayer")) return;
      el.append(...this.o.layers(i, c, this.scale));
      // the handles keep off its text (not while one is held: the turn places them when it ends)
      if (this.turn) this.zonesStale = true;
      else this.placeZones();
    }, (e) => {
      this.askedText.delete(i);
      if (!this.dead) this.o.onError(e);
    });
  }

  private jump(k: number) {
    this.finish();
    if (k === this.at) return this.focusShown();
    this.at = k;
    this.showSpread();
    this.o.onTurn(this.pages);
  }

  // Handles: the corners and the outer edge of the page to turn, straddling
  // its edge so that most of each lies in the margin around the page. Where
  // text or a link lies under one (slides have text up to their corners), the
  // handle draws back from it, so that it can still be selected and clicked.

  private placeZones() {
    this.zonesStale = false;
    this.zones.replaceChildren();
    const s = this.spreads[this.at];
    if (this.canTurn(1) && s.second !== null) {
      const r = this.rectOf(s.second, "second");
      this.addZones("next", r.x + r.w, 1, r, s.second);
    }
    if (this.canTurn(-1)) {
      // the outer edge of the page before the spine; a page alone, its edge at the spine
      const page = s.first ?? s.second!;
      const r = this.rectOf(page, s.first !== null ? "first" : "second");
      this.addZones("prev", r.x, -1, r, page);
    }
  }

  private addZones(dir: Dir, x: number, out: 1 | -1, r: Rect, page: number) {
    const c = Math.min(72, Math.max(36, Math.min(r.w, r.h) * 0.12), r.h / 3);
    const text = this.textBoxes(page);
    // outward on the screen
    const away = this.o.rtl ? (-out as 1 | -1) : out;
    const add = (part: Part, x0: number, x1: number, y0: number, y1: number) => {
      const z = keepOff(this.screen({ x: Math.min(x0, x1), y: y0, w: Math.abs(x1 - x0), h: y1 - y0 }), text, away, part);
      if (z.w < 6 || z.h < 6) return;
      const el = document.createElement("div");
      el.dataset.dir = dir;
      el.dataset.part = part;
      el.style.cssText = `left: ${z.x}px; top: ${z.y}px; width: ${z.w}px; height: ${z.h}px`;
      el.title = dir === "next" ? "next page" : "previous page";
      el.onpointerdown = (e) => this.press(e, el);
      el.onpointermove = (e) => this.move(e);
      el.onpointerup = (e) => this.release(e, false);
      el.onpointercancel = (e) => this.release(e, true);
      // after pointerup too, when the drag has already ended
      el.onlostpointercapture = (e) => this.release(e, false);
      el.onpointerenter = (e) => this.enter(e, el);
      el.onpointerleave = () => this.leave();
      this.zones.appendChild(el);
    };
    const inside = x - out * c * 0.65, outside = x + out * c * 0.35;
    add("top", inside, outside, r.y - c * 0.35, r.y + c * 0.65);
    add("bottom", inside, outside, r.y + r.h - c * 0.65, r.y + r.h + c * 0.35);
    add("edge", x - out * EDGE, x + out * EDGE, r.y + c * 0.65, r.y + r.h - c * 0.65);
  }

  /** The boxes of a shown page's text and links, in CSS px of the book element. */
  private textBoxes(page: number): Rect[] {
    const el = this.pageEl(page);
    if (!el) return [];
    const o = this.el.getBoundingClientRect();
    const boxes: Rect[] = [];
    for (const n of el.querySelectorAll(".bdfTextLayer span, .bdfTextLayer a")) {
      const b = n.getBoundingClientRect();
      if (b.width > 0 && b.height > 0) boxes.push({ x: b.left - o.left, y: b.top - o.top, w: b.width, h: b.height });
    }
    return boxes;
  }

  private animates(): boolean {
    if (this.overflow || matchMedia("(prefers-reduced-motion: reduce)").matches) return false;
    if (this.renderer === undefined) {
      this.renderer = FlipRenderer.create(this.canvas) ?? null;
      this.renderer?.resize(this.width, this.height, dpr());
    }
    return !!this.renderer?.ok;
  }

  // Turning: begin puts the canvas in place of the page elements, land puts
  // the page elements of the spread it ends on back.

  private begin(dir: Dir, part: Part, y: number): Turn | undefined {
    const from = this.at, to = this.at + (dir === "next" ? 1 : -1);
    if (!this.renderer || !this.spreads[to]) return undefined;
    const single = this.o.layout === "single";
    const leaf = (i: number | null, slot: Slot, mirrored: boolean): Leaf | undefined => {
      if (i === null) return undefined;
      const r = this.rectOf(i, slot);
      return { rect: mirrored ? this.flip(r) : r, page: i };
    };
    let scene: Pick<Turn, "mirror" | "sheet" | "front" | "back" | "before" | "under">;
    let reverse = false;
    if (dir === "next" || single) {
      // the second page of the earlier spread turns over the spine
      const A = this.spreads[Math.min(from, to)], B = this.spreads[Math.max(from, to)];
      if (A.second === null) return undefined;
      scene = { mirror: this.o.rtl, sheet: this.rectOf(A.second, "second"), front: A.second, back: single ? null : B.first, before: leaf(A.first, "first", false), under: leaf(B.second, "second", false) };
      reverse = dir === "prev";
    } else {
      // the first page of this spread turns back, seen in a mirror
      const A = this.spreads[to], B = this.spreads[from];
      if (B.first === null) return undefined;
      scene = { mirror: !this.o.rtl, sheet: this.flip(this.rectOf(B.first, "first")), front: B.first, back: A.second, before: leaf(B.second, "second", true), under: leaf(A.first, "first", true) };
    }
    const sheet = scene.sheet;
    const corner = { x: sheet.x + sheet.w, y: part === "top" ? sheet.y : part === "bottom" ? sheet.y + sheet.h : Math.min(sheet.y + sheet.h, Math.max(sheet.y, y)) };
    const rest = reverse ? this.mirror(corner) : corner, done = reverse ? corner : this.mirror(corner);
    const t: Turn = { from, to, dir, ...scene, corner, edge: part === "edge", rest, done, p: rest };
    this.turn = t;
    for (const i of this.turnPages(t)) {
      const b = this.bitmaps.get(i);
      if (b) this.renderer.upload(i, b);
    }
    this.draw();
    this.el.classList.add("turning");
    return t;
  }

  private turnPages(t: Turn): number[] {
    return [t.front, t.back, t.before?.page, t.under?.page].filter((i): i is number => i != null);
  }

  /** End a turn on the spread it turns to, or back on the one it started from. */
  private land(turned: boolean) {
    const t = this.turn;
    if (!t) return;
    this.turn = undefined;
    cancelAnimationFrame(this.frame);
    this.frame = 0;
    const k = turned ? t.to : t.from;
    if (k !== this.at) {
      this.at = k;
      this.showSpread();
      this.o.onTurn(this.pages);
    } else if (this.zonesStale) this.placeZones();
    this.el.classList.remove("turning");
  }

  /** Bring a turn to its end at once: where it was going, or back where it began. */
  private finish() {
    const t = this.turn;
    if (!t) return;
    const a = t.anim;
    t.anim = undefined;
    t.drag = undefined;
    if (a) a.then();
    if (this.turn === t) this.land(false);
  }

  /** How far a turn from rest rises (or falls, from a top corner) on the way, CSS px. */
  private lift(t: Turn): number {
    return (t.corner.y === t.sheet.y ? 1 : -1) * t.sheet.h * 0.12;
  }

  private animate(to: Point, dur: number, lift: number, ease: (t: number) => number, then: () => void) {
    const t = this.turn!;
    t.anim = { from: { ...t.p }, to, t0: performance.now(), dur: Math.max(1, dur), lift, ease, then };
    this.schedule();
  }

  private schedule() {
    if (!this.frame) this.frame = requestAnimationFrame((now) => this.tick(now));
  }

  private tick(now: number) {
    this.frame = 0;
    const t = this.turn;
    if (!t || this.dead) return;
    const a = t.anim;
    if (a) {
      const k = Math.min(1, Math.max(0, (now - a.t0) / a.dur));
      const e = a.ease(k);
      t.p = reach(t.sheet, this.spine, t.corner, t.edge, {
        x: a.from.x + (a.to.x - a.from.x) * e,
        y: a.from.y + (a.to.y - a.from.y) * e + a.lift * Math.sin(Math.PI * e),
      });
      if (k >= 1) {
        t.p = a.to;
        t.anim = undefined;
        this.draw();
        a.then();
        return;
      }
      this.schedule();
    }
    this.draw();
  }

  private draw() {
    const t = this.turn;
    if (!t || !this.renderer) return;
    const fold = foldAt(t.sheet, this.spine, t.corner, t.p);
    this.renderer.draw({
      spine: this.spine,
      mirror: t.mirror,
      sheet: t.sheet,
      front: t.front,
      back: t.back,
      before: t.before,
      under: t.under,
      // a page alone fades away as it lies down beside the page
      backAlpha: this.o.layout === "single" ? 1 - smoothstep(0.7, 1, fold?.progress ?? 0) : 1,
    }, fold);
  }

  // Pointer input on the handles.

  private press(e: PointerEvent, zone: HTMLElement) {
    if (e.button !== 0 || this.turn?.drag) return;
    e.preventDefault(); // no text selection from a handle
    const dir = zone.dataset.dir as Dir, part = zone.dataset.part as Part;
    let t = this.turn;
    if (t && !(t.hover && t.dir === dir)) {
      this.finish();
      t = undefined;
    }
    // a turn ended under the pointer, on a spread with handles of its own
    if (!zone.isConnected) return;
    zone.setPointerCapture(e.pointerId);
    if (!t && this.animates()) t = this.begin(dir, part, this.point(e, false).y);
    if (!t) {
      // no turning: a click changes the page
      this.click = { id: e.pointerId, dir };
      return;
    }
    t.anim = undefined;
    t.hover = false;
    const q = this.point(e, t.mirror);
    t.drag = { id: e.pointerId, zone, start: q, anchor: { ...t.p }, kx: this.o.layout === "single" ? 2 : 1, moved: false, trail: [{ x: q.x, t: e.timeStamp }] };
    zone.classList.add("grabbing");
  }

  private move(e: PointerEvent) {
    const t = this.turn, d = t?.drag;
    if (!t || !d || e.pointerId !== d.id) return;
    const q = this.point(e, t.mirror);
    if (!d.moved && Math.hypot(q.x - d.start.x, q.y - d.start.y) < TAP) return;
    d.moved = true;
    d.trail.push({ x: q.x, t: e.timeStamp });
    while (d.trail.length > 2 && e.timeStamp - d.trail[0].t > 100) d.trail.shift();
    t.p = reach(t.sheet, this.spine, t.corner, t.edge, { x: d.anchor.x + (q.x - d.start.x) * d.kx, y: d.anchor.y + (q.y - d.start.y) });
    this.schedule();
  }

  private release(e: PointerEvent, cancelled: boolean) {
    const c = this.click;
    if (c?.id === e.pointerId) {
      this.click = undefined;
      if (!cancelled) this.jump(this.at + (c.dir === "next" ? 1 : -1));
      return;
    }
    const t = this.turn, d = t?.drag;
    if (!t || !d || e.pointerId !== d.id) return;
    t.drag = undefined;
    d.zone.classList.remove("grabbing");
    let turned: boolean;
    if (cancelled) turned = false;
    else if (!d.moved) turned = true; // a click turns the page
    else {
      // a flick turns where it goes, otherwise the side of the spine the corner is on
      const side = Math.sign(t.done.x - this.spine);
      const first = d.trail[0], last = d.trail[d.trail.length - 1];
      const v = last.t > first.t ? ((last.x - first.x) / (last.t - first.t)) * d.kx : 0;
      turned = Math.abs(v) > 0.35 ? Math.sign(v) === side : Math.sign(t.p.x - this.spine) === side;
    }
    const to = turned ? t.done : t.rest;
    if (!d.moved && !cancelled) return this.animate(to, 650, t.edge ? 0 : this.lift(t), easeInOut, () => this.land(true));
    const dist = Math.hypot(to.x - t.p.x, to.y - t.p.y);
    this.animate(to, 120 + (380 * dist) / (2 * t.sheet.w), 0, easeOut, () => this.land(turned));
  }

  /** The mouse over a corner peels the page back a little. */
  private enter(e: PointerEvent, zone: HTMLElement) {
    const dir = zone.dataset.dir as Dir, part = zone.dataset.part as Part;
    // a page alone would bring the page before in from beside it
    if (e.pointerType !== "mouse" || e.buttons || this.turn || part === "edge" || (dir === "prev" && this.o.layout === "single") || !this.animates()) return;
    const t = this.begin(dir, part, 0);
    if (!t) return;
    t.hover = true;
    const peel = Math.min(t.sheet.w, t.sheet.h) * 0.1;
    this.animate({ x: t.rest.x - peel, y: t.rest.y + (part === "top" ? 0.7 : -0.7) * peel }, 180, 0, easeOut, () => {});
  }

  private leave() {
    const t = this.turn;
    if (!t?.hover || t.drag) return;
    t.hover = false;
    this.animate(t.rest, 160, 0, easeOut, () => this.land(false));
  }

  private readonly onKey = (e: KeyboardEvent) => {
    if (e.defaultPrevented || e.altKey || e.ctrlKey || e.metaKey) return;
    const back = this.o.rtl ? 1 : -1;
    const keys: Record<string, number> = { PageDown: 1, PageUp: -1, ArrowRight: -back, ArrowLeft: back, " ": e.shiftKey ? -1 : 1 };
    if (e.key === "Home" || e.key === "End") {
      e.preventDefault();
      this.go(e.key === "Home" ? 0 : this.o.pages.length - 1, false, true);
      return;
    }
    const d = keys[e.key];
    if (!d) return;
    e.preventDefault();
    this.turnBy(d as 1 | -1);
  };
}

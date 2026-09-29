// A small viewer to start from: it shows a bdf document in an element of the
// page, its pages as wide as the element. The worker of @bdfkit/render decodes
// and draws; this file places the bitmaps, the selectable text over them,
// the tabs of the views and the search hits.
//
//   const viewer = new MiniViewer(document.getElementById("viewer"), { worker: "lib/worker.js" });
//   await viewer.open("/files/report.bdf");       // fetched by ranges, a page at a time
//   await viewer.search("revenue");               // highlights the hits and shows the first
//
// It draws pages (slides, drawings, the pages of a PDF or a Word document),
// one-column documents (Markdown, HTML) and sheets (Excel, CSV: the cells
// with their headers and frozen panes). What it leaves to the demo viewer
// (examples/viewer): pages turned like a book's, the music of a score,
// selecting the cells of a sheet, and pages shown while they are converted.
import { dcValues, type Manifest, type View } from "@bdfkit/core";
import {
  BdfWorkerClient, BdfWorkerError, buildTextLayer, installCopyHandler, TEXT_LAYER_CSS, RUN_ATTR,
  type HitRect, type OpenSource,
} from "@bdfkit/render";

export interface MiniViewerOptions {
  /** URL of the rendering worker: packages/render/src/worker.ts, bundled as a module. */
  worker: string;
  /** Pages are as wide as the viewer, but no larger than this many CSS pixels a unit (default 1.5). */
  maxScale?: number;
  /** Asks for the password of an encrypted document; undefined gives up (default: window.prompt). */
  password?: (wrong: boolean) => Promise<string | undefined>;
  /** Called when a page could not be drawn (default: console.error). */
  onError?: (e: unknown) => void;
  /** Called when the view shown changes, and when pages come into view (the first one in view). */
  onChange?: (state: { view: View; page: number }) => void;
}

type Box = { x: number; y: number; w: number; h: number };

/** The side of the bands a one-column view is drawn in, in units. */
const BAND = 800;
/** Pages are drawn when they come this near the visible area; their text from further away. */
const BITMAP_MARGIN = "400px";
const TEXT_MARGIN = "200% 0px";
/** The row and column headers of a sheet, in CSS pixels. */
const HEADER = { w: 40, h: 20 };

export const MINI_VIEWER_CSS = `
.bdfMini { display: grid; grid-template-rows: minmax(0, 1fr) auto; min-height: 0; background: #e9ecef; color: #222; font: 13px system-ui, sans-serif; }
/* the room of the scroll bar is kept, so that the width pages are fitted to does not change as they come */
.bdfMini-stage { position: relative; overflow: auto; scrollbar-gutter: stable; }
.bdfMini-stage:focus-visible { outline: 2px solid #1a73e8; outline-offset: -2px; }
.bdfMini-pages { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 12px; box-sizing: border-box; width: fit-content; min-width: 100%; }
.bdfMini-page { position: relative; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,.3); }
.bdfMini-page canvas { display: block; }
.bdfMini-layer { position: absolute; inset: 0; overflow: hidden; }
.bdfMini-hits { pointer-events: none; forced-color-adjust: none; }
.bdfMini-hit { position: absolute; background: rgba(255, 210, 0, .45); outline: 1px solid rgba(200, 140, 0, .8); border-radius: 2px; }
.bdfMini-hit.current { background: rgba(255, 120, 0, .5); outline-color: #c04000; }
.bdfMini-sheet { position: relative; }
.bdfMini-sheet canvas { position: sticky; top: 0; left: 0; display: block; }
.bdfMini-tabs { display: flex; overflow-x: auto; background: #f3f3f3; border-top: 1px solid #c8c8c8; scrollbar-width: none; }
.bdfMini-tabs:empty { display: none; }
.bdfMini-tabs button { flex: none; max-width: 16em; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; margin: 0 -1px 0 0; padding: 6px 14px;
  border: 1px solid #c8c8c8; border-top: 0; background: #e6e6e6; color: #333; font: inherit; cursor: pointer; }
.bdfMini-tabs button[aria-selected="true"] { background: #fff; color: #1d6b40; font-weight: 600; box-shadow: inset 0 -2px 0 #1d6b40; }
.bdfMini-message { margin: 0; padding: 24px; color: #555; text-align: center; }
`;

/** Row or column sizes of a sheet, from the manifest's run-length list. */
class Axis {
  readonly count: number;
  readonly total: number;
  constructor(private runs: [number, number][] = []) {
    this.count = runs.reduce((a, [n]) => a + n, 0);
    this.total = runs.reduce((a, [n, s]) => a + n * s, 0);
  }
  /** Where entry i starts. */
  pos(i: number): number {
    let p = 0;
    for (const [n, s] of this.runs) {
      if (i <= n) return p + i * s;
      p += n * s;
      i -= n;
    }
    return p;
  }
  /** Call fn for the entries with a size that overlap [from, to). */
  each(from: number, to: number, fn: (i: number, start: number, size: number) => void) {
    let i = 0, p = 0;
    for (const [n, s] of this.runs) {
      if (s > 0 && p + n * s > from) {
        for (let k = Math.max(0, Math.floor((from - p) / s)); k < n; k++) {
          const start = p + k * s;
          if (start >= to) return;
          fn(i + k, start, s);
        }
      }
      p += n * s;
      i += n;
    }
  }
}

/** Column letters: A … Z, AA … */
function columnLabel(c: number): string {
  let s = "";
  for (c++; c > 0; c = Math.floor((c - 1) / 26)) s = String.fromCharCode(65 + ((c - 1) % 26)) + s;
  return s;
}

export class MiniViewer {
  readonly client: BdfWorkerClient;
  manifest?: Manifest;
  /** The view shown. */
  view?: View;
  /** The reader's zoom, over the width of the viewer (1: pages as wide as it). */
  zoom = 1;
  /** Hits of the last search, their rectangles, and the one shown. */
  hits: HitRect[][] = [];
  hit = -1;

  private stage: HTMLDivElement;
  private tabs: HTMLDivElement;
  /** Bumped as another view, zoom or document is shown: work started before is dropped. */
  private generation = 0;
  /** CSS pixels a unit of the view shown. */
  private scale = 1;
  private query = "";
  /** The width of the viewer the view shown was laid out for, and the first page in view. */
  private width = 0;
  private page = 0;
  private redrawSheet?: () => void;
  private resize: ResizeObserver;
  private uninstallCopy: () => void;

  constructor(readonly host: HTMLElement, private options: MiniViewerOptions) {
    if (!document.getElementById("bdfMiniCSS")) {
      const style = document.createElement("style");
      style.id = "bdfMiniCSS";
      style.textContent = TEXT_LAYER_CSS + MINI_VIEWER_CSS;
      document.head.append(style);
    }
    this.client = new BdfWorkerClient(new Worker(options.worker, { type: "module" }));
    host.classList.add("bdfMini");
    this.stage = document.createElement("div");
    this.stage.className = "bdfMini-stage";
    this.stage.tabIndex = 0;
    this.stage.setAttribute("role", "tabpanel");
    this.tabs = document.createElement("div");
    this.tabs.className = "bdfMini-tabs";
    this.tabs.setAttribute("role", "tablist");
    host.replaceChildren(this.stage, this.tabs);
    // copy takes the text of the selected runs, with the document's spaces and line breaks
    this.uninstallCopy = installCopyHandler(this.stage);
    this.stage.addEventListener("click", (e) => this.followLink(e));
    // pages are as wide as the viewer: laid out again when its width changes
    this.resize = new ResizeObserver(() => {
      if (this.view && host.clientWidth !== this.width) this.show(this.view.id, this.page);
    });
    this.resize.observe(host);
  }

  /**
   * Open a document: a URL (a single-file bdf, read by ranges; the
   * directory of a split one when it ends with a slash), or what the worker
   * opens (a buffer, or a server's segments). Returns its manifest, or
   * undefined when it is encrypted and no password opened it.
   */
  async open(source: string | OpenSource, password?: string): Promise<Manifest | undefined> {
    this.generation++;
    this.view = this.manifest = undefined;
    this.clearSearch();
    this.tabs.replaceChildren();
    this.message("loading…");
    const src: OpenSource = typeof source !== "string" ? source
      : source.endsWith("/") ? { kind: "split", base: new URL(source, location.href).href }
      : { kind: "single", url: new URL(source, location.href).href, range: true };
    const ask = this.options.password ?? (async (wrong) => prompt(wrong ? "Wrong password. Try again." : "This document is encrypted. Its password:") ?? undefined);
    let manifest: Manifest | undefined;
    try {
      manifest = await this.client.open(src, password);
    } catch (e) {
      if (!(e instanceof BdfWorkerError) || (e.code !== "password-required" && e.code !== "wrong-password")) {
        this.message(`The document could not be opened: ${(e as Error).message ?? e}`);
        throw e;
      }
      // the worker keeps the locked document: passwords are tried until one opens it
      for (let wrong = e.code === "wrong-password"; !manifest; wrong = true) {
        const pw = await ask(wrong);
        if (pw === undefined) {
          this.message("The document is encrypted.");
          return undefined;
        }
        manifest = await this.client.unlock(pw).catch((err) => {
          if (err instanceof BdfWorkerError && err.code === "wrong-password") return undefined;
          throw err;
        });
      }
    }
    this.manifest = manifest;
    this.tabs.replaceChildren(...(manifest.views.length > 1 ? manifest.views : []).map((v) => {
      const b = document.createElement("button");
      b.type = "button";
      b.textContent = v.title || v.id;
      b.dataset.id = v.id;
      b.setAttribute("role", "tab");
      b.onclick = () => this.show(v.id);
      return b;
    }));
    this.show(manifest.views[0].id);
    return manifest;
  }

  /** Show a view of the document (its first page, or the page given). */
  show(id: string, page = 0) {
    const v = this.manifest?.views.find((v) => v.id === id);
    if (!v) return;
    if (v !== this.view) this.clearSearch();
    this.view = v;
    this.generation++;
    this.redrawSheet = undefined;
    for (const b of this.tabs.querySelectorAll("button")) b.setAttribute("aria-selected", String(b.dataset.id === id));
    this.stage.onscroll = null;
    this.stage.replaceChildren();
    this.stage.scrollTo(0, 0);
    this.width = this.host.clientWidth;
    this.page = page;
    if (v.kind === "sheet") this.showSheet(v);
    else if (v.kind === "scroll") this.showColumn(v);
    else this.showPages(v);
    if (page > 0) this.goToPage(page);
    this.options.onChange?.({ view: v, page });
  }

  /** Zoom in or out: 1 shows pages as wide as the viewer. */
  setZoom(zoom: number) {
    this.zoom = Math.min(4, Math.max(0.25, zoom));
    if (this.view) this.show(this.view.id);
  }

  /** Scroll to a page (0-based) of the view shown. */
  goToPage(index: number) {
    const el = this.stage.querySelector<HTMLElement>(`.bdfMini-page[data-index="${index}"]`);
    el?.scrollIntoView({ block: "start" });
    el?.focus({ preventScroll: true });
  }

  /**
   * Search the view shown and show the first hit (the last one with delta
   * -1); the same query again goes on to the next hit. Returns how many
   * hits there are. An empty query clears the highlights.
   */
  async search(query: string, delta = 1): Promise<number> {
    const v = this.view;
    if (!v) return 0;
    if (query !== this.query) {
      const found = query ? await this.client.search(v.id, query, { limit: 500 }) : [];
      const rects = found.length ? await this.client.locate(v.id, found) : [];
      if (v !== this.view) return 0;
      this.query = query;
      this.hits = rects;
      this.hit = delta < 0 ? 0 : -1;
    }
    this.step(delta);
    return this.hits.length;
  }

  /** Show the next hit (delta 1) or the one before (-1). */
  step(delta: number) {
    const v = this.view, n = this.hits.length;
    if (!v) return;
    if (n) this.hit = (this.hit + delta + n) % n;
    this.redrawSheet?.();
    for (const el of this.stage.querySelectorAll<HTMLElement>(".bdfMini-page")) {
      el.querySelector(".bdfMini-hits")?.replaceWith(this.hitLayer(v, el));
    }
    const r = this.hits[this.hit]?.[0];
    if (!r) return;
    if (v.kind === "sheet") {
      this.stage.scrollTo({ left: Math.max(0, r.x * this.scale - this.stage.clientWidth / 2), top: Math.max(0, r.y * this.scale - this.stage.clientHeight / 2) });
    } else if (v.kind === "scroll") {
      this.stage.scrollTo({ top: Math.max(0, this.inColumn(v, r).y * this.scale - this.stage.clientHeight / 2) });
    } else {
      const el = this.stage.querySelector<HTMLElement>(`.bdfMini-page[data-index="${r.a}"]`);
      if (el) this.stage.scrollTo({ top: Math.max(0, el.offsetTop + r.y * this.scale - this.stage.clientHeight / 3) });
    }
  }

  /** Let the worker and the document go. */
  destroy() {
    this.generation++;
    this.resize.disconnect();
    this.uninstallCopy();
    this.client.terminate();
    this.host.replaceChildren();
  }

  /** Drawing failed: the page stays blank, and the console says why. */
  private failed = (e: unknown) => {
    (this.options.onError ?? console.error)(e);
  };

  private clearSearch() {
    this.query = "";
    this.hits = [];
    this.hit = -1;
  }

  private message(text: string) {
    const p = document.createElement("p");
    p.className = "bdfMini-message";
    p.setAttribute("role", "status");
    p.textContent = text;
    this.stage.replaceChildren(p);
  }

  private get language(): string {
    return dcValues(this.manifest?.meta?.dc?.language)[0] ?? "";
  }

  /** CSS pixels a unit for something width units wide: as wide as the viewer, times the zoom. */
  private fit(width: number, margin: number): number {
    const room = Math.max(120, this.stage.clientWidth - margin);
    return Math.min(this.options.maxScale ?? 1.5, room / Math.max(1, width)) * this.zoom;
  }

  /** Call render once for each element as it comes near the visible area. */
  private near(margin: string, render: (el: HTMLElement) => void): IntersectionObserver {
    const observer = new IntersectionObserver((entries) => {
      for (const e of entries) {
        if (!e.isIntersecting) continue;
        observer.unobserve(e.target);
        render(e.target as HTMLElement);
      }
    }, { root: this.stage, rootMargin: margin });
    return observer;
  }

  /** Put a bitmap on a page, under its text. */
  private place(el: HTMLElement, bmp: ImageBitmap, w: number, h: number) {
    const canvas = document.createElement("canvas");
    canvas.width = bmp.width;
    canvas.height = bmp.height;
    canvas.style.width = `${w * this.scale}px`;
    canvas.style.height = `${h * this.scale}px`;
    canvas.setAttribute("aria-hidden", "true");
    canvas.getContext("2d")!.drawImage(bmp, 0, 0);
    bmp.close();
    el.prepend(canvas);
  }

  private showPages(v: View) {
    const gen = this.generation, pages = v.pages ?? [];
    this.scale = this.fit(Math.max(...pages.map((p) => p.w), 1), 24);
    const dpr = window.devicePixelRatio || 1;
    const live = () => gen === this.generation;
    const bitmaps = this.near(BITMAP_MARGIN, (el) => {
      const i = Number(el.dataset.index);
      this.client.page(v.id, i, this.scale * dpr).then((bmp) => (live() ? this.place(el, bmp, pages[i].w, pages[i].h) : bmp.close())).catch(this.failed);
    });
    const texts = this.near(TEXT_MARGIN, (el) => {
      this.client.content(v.id, Number(el.dataset.index)).then((content) => {
        if (live()) el.append(buildTextLayer(content, this.scale, { lang: this.language }), this.hitLayer(v, el));
      }).catch(this.failed);
    });
    // the first page in view, for onChange
    const inView = new Set<number>();
    const visible = new IntersectionObserver((entries) => {
      for (const e of entries) {
        const i = Number((e.target as HTMLElement).dataset.index);
        if (e.isIntersecting) inView.add(i); else inView.delete(i);
      }
      if (!live() || !inView.size) return;
      this.page = Math.min(...inView);
      this.options.onChange?.({ view: v, page: this.page });
    }, { root: this.stage });
    const list = document.createElement("div");
    list.className = "bdfMini-pages";
    pages.forEach((p, i) => {
      const el = document.createElement("div");
      el.className = "bdfMini-page";
      el.dataset.index = String(i);
      el.setAttribute("role", "group");
      el.setAttribute("aria-label", pages.length === 1 && v.title ? v.title : `Page ${i + 1} of ${pages.length}`);
      el.tabIndex = -1; // target of page links
      el.style.width = `${p.w * this.scale}px`;
      el.style.height = `${p.h * this.scale}px`;
      list.append(el);
      bitmaps.observe(el);
      texts.observe(el);
      visible.observe(el);
    });
    this.stage.append(list);
  }

  /** The strips of a one-column view, one above the other. */
  private column(v: View) {
    let height = 0, width = 0;
    const tops: number[] = [];
    for (const p of v.pages ?? []) {
      const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
      tops.push(height);
      height += b.h;
      width = Math.max(width, b.w);
    }
    return { width, height, tops };
  }

  /** A rectangle of a strip in the coordinates of the column. */
  private inColumn(v: View, r: HitRect): Box {
    const p = v.pages![r.a];
    const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
    return { x: r.x - b.x, y: this.column(v).tops[r.a] + r.y - b.y, w: r.w, h: r.h };
  }

  /** A view of one long column: drawn in bands as they come into view. */
  private showColumn(v: View) {
    const gen = this.generation, { width, height } = this.column(v);
    this.scale = this.fit(width, 24);
    const dpr = window.devicePixelRatio || 1;
    const live = () => gen === this.generation;
    const band = (el: HTMLElement): Box => {
      const y = Number(el.dataset.y);
      return { x: 0, y, w: width, h: Math.min(BAND, height - y) };
    };
    const bitmaps = this.near(BITMAP_MARGIN, (el) => {
      const b = band(el);
      this.client.continuous(v.id, b, this.scale * dpr).then((bmp) => (live() ? this.place(el, bmp, b.w, b.h) : bmp.close())).catch(this.failed);
    });
    const texts = this.near(TEXT_MARGIN, (el) => {
      this.client.continuousContent(v.id, band(el)).then((content) => {
        if (live()) el.append(buildTextLayer(content, this.scale, { lang: this.language }), this.hitLayer(v, el));
      }).catch(this.failed);
    });
    const list = document.createElement("div");
    list.className = "bdfMini-pages";
    list.style.gap = "0";
    for (let y = 0; y < height; y += BAND) {
      const el = document.createElement("div");
      el.className = "bdfMini-page";
      el.dataset.y = String(y);
      el.style.width = `${width * this.scale}px`;
      el.style.height = `${Math.min(BAND, height - y) * this.scale}px`;
      el.style.boxShadow = "none";
      list.append(el);
      bitmaps.observe(el);
      texts.observe(el);
    }
    this.stage.append(list);
  }

  /** The hits on a page or a band, over its text. */
  private hitLayer(v: View, el: HTMLElement): HTMLDivElement {
    const layer = document.createElement("div");
    layer.className = "bdfMini-layer bdfMini-hits";
    const band = el.dataset.y === undefined ? undefined : Number(el.dataset.y);
    this.hits.forEach((rects, i) => {
      for (const r of rects) {
        let box: Box = r;
        if (band === undefined) {
          if (r.a !== Number(el.dataset.index)) continue;
        } else {
          box = this.inColumn(v, r);
          if (box.y + box.h < band || box.y > band + BAND) continue;
          box = { ...box, y: box.y - band };
        }
        const d = document.createElement("div");
        d.className = i === this.hit ? "bdfMini-hit current" : "bdfMini-hit";
        d.style.cssText = `left: ${box.x * this.scale}px; top: ${box.y * this.scale}px; width: ${box.w * this.scale}px; height: ${box.h * this.scale}px`;
        layer.append(d);
      }
    });
    return layer;
  }

  /**
   * A sheet: a scroll area of its size with one canvas that stays in place
   * and shows the headers and the cells in view, in the frozen panes of
   * the view. Each pane is drawn by the worker when the scroll settles on
   * it; the last picture moves with the scroll meanwhile.
   */
  private showSheet(v: View) {
    const gen = this.generation, stage = this.stage;
    const cols = new Axis(v.cols), rows = new Axis(v.rows);
    const k = this.scale = (this.options.maxScale ? Math.min(this.options.maxScale, 1) : 1) * this.zoom;
    const fw = cols.pos(Math.min(v.freeze?.cols ?? 0, cols.count)), fh = rows.pos(Math.min(v.freeze?.rows ?? 0, rows.count));
    const wrap = document.createElement("div");
    wrap.className = "bdfMini-sheet";
    wrap.style.width = `${HEADER.w + cols.total * k}px`;
    wrap.style.height = `${HEADER.h + rows.total * k}px`;
    const canvas = document.createElement("canvas");
    canvas.setAttribute("role", "img");
    canvas.setAttribute("aria-label", v.title || "sheet");
    wrap.append(canvas);
    stage.append(wrap);

    /** The panes in view: the region of the sheet each shows and where it goes on the canvas (CSS px), with its picture. */
    type Pane = { src: Box; dx: number; dy: number; bitmap?: ImageBitmap; drawn?: Box };
    let panes: Pane[] = [];
    const paint = () => {
      const d = window.devicePixelRatio || 1;
      const vw = Math.min(stage.clientWidth, HEADER.w + cols.total * k), vh = Math.min(stage.clientHeight, HEADER.h + rows.total * k);
      canvas.width = Math.ceil(vw * d);
      canvas.height = Math.ceil(vh * d);
      canvas.style.width = `${vw}px`;
      canvas.style.height = `${vh}px`;
      const ctx = canvas.getContext("2d")!;
      ctx.scale(d, d);
      ctx.fillStyle = "#fff";
      ctx.fillRect(0, 0, vw, vh);
      for (const p of panes) {
        ctx.save();
        ctx.beginPath();
        ctx.rect(p.dx, p.dy, p.src.w * k, p.src.h * k);
        ctx.clip();
        // the picture of the region the pane showed last, where that region is now
        if (p.bitmap && p.drawn) ctx.drawImage(p.bitmap, p.dx + (p.drawn.x - p.src.x) * k, p.dy + (p.drawn.y - p.src.y) * k, p.drawn.w * k, p.drawn.h * k);
        this.hits.forEach((rects, i) => {
          ctx.fillStyle = i === this.hit ? "rgba(255,120,0,.5)" : "rgba(255,210,0,.45)";
          for (const r of rects) ctx.fillRect(p.dx + (r.x - p.src.x) * k, p.dy + (r.y - p.src.y) * k, r.w * k, r.h * k);
        });
        ctx.restore();
      }
      // the headers: frozen columns and rows, then the scrolled ones
      const sx = stage.scrollLeft, sy = stage.scrollTop;
      ctx.fillStyle = "#f3f4f6";
      ctx.fillRect(0, 0, vw, HEADER.h);
      ctx.fillRect(0, 0, HEADER.w, vh);
      ctx.font = "11px system-ui, sans-serif";
      ctx.textAlign = "center";
      ctx.textBaseline = "middle";
      ctx.strokeStyle = "#c8ccd2";
      ctx.fillStyle = "#444";
      const cell = (x: number, y: number, w: number, h: number, label: string) => {
        ctx.fillText(label, x + w / 2, y + h / 2);
        ctx.strokeRect(Math.round(x) + 0.5, Math.round(y) + 0.5, Math.round(w), Math.round(h));
      };
      const clipped = (x: number, y: number, draw: () => void) => {
        ctx.save();
        ctx.beginPath();
        ctx.rect(x, y, vw - x, vh - y);
        ctx.clip();
        draw();
        ctx.restore();
      };
      clipped(HEADER.w, 0, () => cols.each(0, fw, (c, start, size) => cell(HEADER.w + start * k, 0, size * k, HEADER.h, columnLabel(c))));
      clipped(HEADER.w + fw * k, 0, () => cols.each(fw + sx / k, fw + (sx + vw) / k, (c, start, size) => cell(HEADER.w + start * k - sx, 0, size * k, HEADER.h, columnLabel(c))));
      clipped(0, HEADER.h, () => rows.each(0, fh, (r, start, size) => cell(0, HEADER.h + start * k, HEADER.w, size * k, String(r + 1))));
      clipped(0, HEADER.h + fh * k, () => rows.each(fh + sy / k, fh + (sy + vh) / k, (r, start, size) => cell(0, HEADER.h + start * k - sy, HEADER.w, size * k, String(r + 1))));
      ctx.fillStyle = "#e5e7eb";
      ctx.fillRect(0, 0, HEADER.w, HEADER.h);
    };

    let asked = 0;
    const draw = () => {
      if (gen !== this.generation) return;
      const mx = fw + stage.scrollLeft / k, my = fh + stage.scrollTop / k;
      const pw = Math.max(0, Math.min((stage.clientWidth - HEADER.w) / k - fw, cols.total - mx));
      const ph = Math.max(0, Math.min((stage.clientHeight - HEADER.h) / k - fh, rows.total - my));
      const regions = [
        { src: { x: mx, y: my, w: pw, h: ph }, dx: HEADER.w + fw * k, dy: HEADER.h + fh * k },
        { src: { x: mx, y: 0, w: pw, h: fh }, dx: HEADER.w + fw * k, dy: HEADER.h },
        { src: { x: 0, y: my, w: fw, h: ph }, dx: HEADER.w, dy: HEADER.h + fh * k },
        { src: { x: 0, y: 0, w: fw, h: fh }, dx: HEADER.w, dy: HEADER.h },
      ];
      // each pane keeps its picture until the new one comes
      panes = regions.map((r, i) => ({ ...r, bitmap: panes[i]?.bitmap, drawn: panes[i]?.drawn }));
      paint();
      const turn = ++asked;
      const d = window.devicePixelRatio || 1;
      panes.forEach((p) => {
        if (p.src.w <= 0 || p.src.h <= 0) return;
        const src = p.src;
        this.client.sheet(v.id, src, k * d).then((bmp) => {
          if (gen !== this.generation || turn !== asked) return bmp.close();
          p.bitmap?.close();
          p.bitmap = bmp;
          p.drawn = src;
          paint();
        }).catch(this.failed);
      });
    };
    let frame = 0;
    stage.onscroll = () => {
      paint();
      cancelAnimationFrame(frame);
      frame = requestAnimationFrame(draw);
    };
    this.redrawSheet = paint;
    draw();
  }

  /** Links to a page ("#page=N") scroll there; links to a view ("#view=ID") show it. */
  private followLink(e: MouseEvent) {
    const a = (e.target as Element).closest(`a[${RUN_ATTR.page}], a[${RUN_ATTR.view}]`);
    if (!a || !this.view) return;
    e.preventDefault();
    const page = Number(a.getAttribute(RUN_ATTR.page) ?? "1") - 1;
    const id = a.getAttribute(RUN_ATTR.view);
    if (id !== null) this.show(id, page);
    else this.goToPage(page);
  }
}

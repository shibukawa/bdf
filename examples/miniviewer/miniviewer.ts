// A small viewer to start from: it shows a bdf document in an element of the
// page, its pages as wide as the element. The worker of @bdfkit/render decodes
// and draws; this file places the bitmaps, the selectable text over them,
// the tabs of the views and the search hits.
//
//   const viewer = new MiniViewer(document.getElementById("viewer"), { worker: "lib/worker.js" });
//   await viewer.open("/files/report.bdf");       // fetched by ranges, a page at a time
//   await viewer.search("revenue");               // highlights the hits and shows the first
//
// The find shortcut (Cmd+F, Ctrl+F) opens a find bar of the viewer's own
// while the focus is in it: the browser finds the text of the pages shown
// only, the bar that of the whole view.
//
// It draws pages (slides, drawings, the pages of a PDF or a Word document),
// one-column documents (Markdown, HTML) and sheets (Excel, CSV: the cells
// with their headers and frozen panes). The card of an audio file plays, with
// the browser's own controls under it; the line of its lyrics that is sung
// is marked when the file says when each is, and a click on a line plays
// from it. What it leaves to the demo viewer (examples/viewer): pages turned
// like a book's, the music of a score, selecting the cells of a sheet, and
// pages shown while they are converted.
import { dcValues, tileSize, type Manifest, type View } from "@bdfkit/core";
import {
  AudioPlayer, BdfWorkerClient, BdfWorkerError, FindBar, FindHits, buildTextLayer, findKey, installCopyHandler, TEXT_LAYER_CSS, RUN_ATTR,
  type FindBarLabels, type HitRect, type OpenSource,
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
  /**
   * Whose the browser's find shortcut is (Cmd+F on macOS, Ctrl+F elsewhere).
   * The browser finds the text the page holds, which is that of the pages
   * shown and those near them; the viewer's find bar searches the whole
   * view. "focus" (the default): the viewer's while the focus is in it;
   * "page": the viewer's wherever the focus is, for a page that is nothing
   * but the viewer; false: the browser's, for a host with a search box of
   * its own (see search) or a document whose pages come as they are read
   * (a search covers the pages that came).
   */
  find?: "focus" | "page" | false;
  /** The words of the find bar (default: English). */
  findLabels?: Partial<FindBarLabels>;
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
.bdfMini { position: relative; display: grid; grid-template-rows: minmax(0, 1fr) auto auto; min-height: 0; background: #e9ecef; color: #222; font: 13px system-ui, sans-serif; }
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
/* the recording of an audio file: the browser's controls, and the line of the lyrics that is sung */
.bdfMini-audio { display: flex; padding: 6px 12px; background: #f3f3f3; border-top: 1px solid #c8c8c8; }
.bdfMini-audio[hidden] { display: none; }
.bdfMini-audio audio { flex: 1; min-width: 0; height: 36px; }
.bdfMini-line { position: absolute; z-index: 2; box-sizing: border-box; border-left: 3px solid rgba(26, 115, 232, .8); border-radius: 3px; background: rgba(26, 115, 232, .14); pointer-events: none; forced-color-adjust: none; }
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
  /** The hits of the last search (they are located as their pages are shown), and the one shown. */
  found?: FindHits;
  hit = -1;

  private stage: HTMLDivElement;
  private tabs: HTMLDivElement;
  /** The controls of the recording a view plays (spec §4.4), its player, and the line it is at. */
  private audioBar: HTMLDivElement;
  private audio?: { view: View; player?: AudioPlayer };
  private line: HTMLDivElement;
  /** Bumped as another view, zoom or document is shown: work started before is dropped. */
  private generation = 0;
  /** CSS pixels a unit of the view shown. */
  private scale = 1;
  /** Bumped by each search and step: a hit located for an earlier one is not gone to. */
  private turn = 0;
  private bar?: FindBar;
  /** Where the find shortcut is listened for. */
  private keys?: HTMLElement | Document;
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
    this.audioBar = document.createElement("div");
    this.audioBar.className = "bdfMini-audio";
    this.audioBar.hidden = true;
    this.line = document.createElement("div");
    this.line.className = "bdfMini-line";
    this.line.setAttribute("aria-hidden", "true");
    host.replaceChildren(this.stage, this.audioBar, this.tabs);
    if (options.find !== false) {
      this.bar = new FindBar({
        labels: options.findLabels,
        onFind: (query, step) => { this.search(query, step).catch(this.failed); },
        onClose: () => {
          this.clearSearch();
          this.markHits();
          this.stage.focus({ preventScroll: true });
        },
      });
      host.append(this.bar.element);
      this.keys = options.find === "page" ? document : host;
      this.keys.addEventListener("keydown", this.onKey as EventListener);
    }
    // copy takes the text of the selected runs, with the document's spaces and line breaks
    this.uninstallCopy = installCopyHandler(this.stage);
    this.stage.addEventListener("click", (e) => this.followLink(e) || this.playFrom(e));
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
    this.bar?.reset();
    this.dropAudio();
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
    const other = v !== this.view;
    if (other) this.clearSearch();
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
    // a view shown again (another width or zoom) goes on playing: the line is put on its new page
    if (this.audio?.view !== v) this.loadAudio(v);
    else this.followAudio();
    this.options.onChange?.({ view: v, page });
    // the find bar goes on with its query in the view that came
    if (other && this.bar?.shown) this.search(this.bar.query, 0).catch(this.failed);
  }

  /** The recording a view plays, if it plays one: the browser's controls under the pages. */
  private loadAudio(v: View) {
    this.dropAudio();
    if (!v.play?.audio) return;
    const audio: { view: View; player?: AudioPlayer } = this.audio = { view: v };
    this.client.audio(v.id).then((data) => {
      if (this.audio !== audio || !data) return;
      const el = document.createElement("audio");
      el.controls = true;
      el.setAttribute("aria-label", v.title || "audio");
      const player = audio.player = new AudioPlayer(data.blob, "", data.cues, { element: el });
      // the line that is sung: the element says where it is a few times a second
      player.onUpdate = () => this.followAudio();
      el.addEventListener("timeupdate", () => this.followAudio());
      player.onError = this.failed;
      this.audioBar.replaceChildren(el);
      this.audioBar.hidden = false;
    }).catch(this.failed);
  }

  private dropAudio() {
    this.audio?.player?.dispose();
    this.audio = undefined;
    this.line.remove();
    this.audioBar.replaceChildren();
    this.audioBar.hidden = true;
  }

  /** Mark the line of the page the recording is at, and bring it into view when it changes. */
  private followAudio() {
    const player = this.audio?.player;
    const c = player && player.state !== "stopped" ? player.cursorAt(player.position) : null;
    const el = c && this.stage.querySelector<HTMLElement>(`.bdfMini-page[data-index="${c.page}"]`);
    if (!c || !el) return this.line.remove();
    const s = this.scale, room = Math.min(8, c.x * s);
    const moved = this.line.parentElement !== el || this.line.dataset.system !== String(c.system);
    this.line.dataset.system = String(c.system);
    this.line.style.cssText = `left:${c.x * s - room}px;top:${c.y * s}px;width:${(c.w ?? 0) * s + 2 * room}px;height:${c.h * s}px`;
    if (this.line.parentElement !== el) el.append(this.line);
    if (!moved) return;
    const top = el.offsetTop + c.y * s, view = this.stage;
    if (top < view.scrollTop || top + c.h * s > view.scrollTop + view.clientHeight) view.scrollTo({ top: Math.max(0, top - view.clientHeight / 4), behavior: "smooth" });
  }

  /** A click on a line that has a time plays from it; false when the click is on none. */
  private playFrom(e: MouseEvent): boolean {
    const player = this.audio?.player, el = (e.target as Element).closest<HTMLElement>(".bdfMini-page[data-index]");
    if (!player?.cues || !el || !(getSelection()?.isCollapsed ?? true)) return false;
    const r = el.getBoundingClientRect(), page = Number(el.dataset.index);
    const x = (e.clientX - r.left) / this.scale, y = (e.clientY - r.top) / this.scale;
    const i = player.cues.systems.findIndex((s) => s.page === page && x >= s.x && x <= s.x + s.w && y >= s.y && y < s.y + s.h);
    const t = i < 0 ? null : player.timeAt(i);
    if (t === null) return false;
    player.seek(t);
    return true;
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
   * Search the view shown, all of it, and show a hit: of another query, the
   * first one on the page being read or after it (delta -1: the last one
   * before it); of the same query again, the next one (delta -1: the one
   * before; 0: the same). Returns how many hits there are. An empty query
   * clears the highlights.
   */
  async search(query: string, delta = 1): Promise<number> {
    const v = this.view;
    if (!v) return 0;
    const turn = ++this.turn;
    if (query !== this.found?.query) {
      const found = await FindHits.search(this.client, v, query);
      // another search, view or document came meanwhile
      if (turn !== this.turn) return found.length;
      this.found = found;
      const n = found.length, start = found.from(this.reading(v));
      this.hit = !n ? -1 : delta < 0 ? (start - 1 + n) % n : start;
    } else {
      const n = this.found.length;
      if (n) this.hit = (this.hit + delta + n) % n;
    }
    const found = this.found;
    this.report();
    this.markHits();
    // where the hit is on its page is known once the page is read
    const r = this.hit < 0 ? undefined : (await found.locateHit(this.hit))[0];
    if (turn !== this.turn || !r) return found.length;
    this.markHits();
    if (v.kind === "sheet") {
      this.stage.scrollTo({ left: Math.max(0, r.x * this.scale - this.stage.clientWidth / 2), top: Math.max(0, r.y * this.scale - this.stage.clientHeight / 2) });
    } else if (v.kind === "scroll") {
      this.stage.scrollTo({ top: Math.max(0, this.inColumn(v, r).y * this.scale - this.stage.clientHeight / 2) });
    } else {
      const el = this.stage.querySelector<HTMLElement>(`.bdfMini-page[data-index="${r.a}"]`);
      if (el) this.stage.scrollTo({ top: Math.max(0, el.offsetTop + r.y * this.scale - this.stage.clientHeight / 3) });
    }
    return found.length;
  }

  /** Show the next hit (delta 1) or the one before (-1). */
  step(delta: number) {
    if (this.found) this.search(this.found.query, delta).catch(this.failed);
  }

  /** Show the find bar, as the find shortcut does (nothing without a document, or with find: false). */
  openFind() {
    const bar = this.bar;
    if (!bar || !this.view) return;
    bar.open();
    // the bar opens with the query it had: its hits are shown again
    if (bar.query && !this.found) this.search(bar.query, 0).catch(this.failed);
  }

  /** The find shortcuts: the bar opens, and while it is shown the next and previous match keys are its own. */
  private onKey = (e: KeyboardEvent) => {
    const bar = this.bar, key = findKey(e);
    if (!bar || !key || e.defaultPrevented || !this.view) return;
    if (key === "open") {
      // pressed again in the field, the shortcut is the browser's: its own search stays at hand
      if (bar.shown && bar.focused) return;
      e.preventDefault();
      this.openFind();
    } else if (bar.shown && bar.query) {
      e.preventDefault();
      this.search(bar.query, key === "next" ? 1 : -1).catch(this.failed);
    }
  };

  /** Where a search starts: the page being read, or the strip of a one-column view at the top of the visible area. */
  private reading(v: View): number {
    if (v.kind !== "scroll") return this.page;
    const y = this.stage.scrollTop / this.scale, tops = this.column(v).tops;
    let i = 0;
    while (i + 1 < tops.length && tops[i + 1] <= y) i++;
    return i;
  }

  /** Tell the find bar what the search found: the count, and for screen readers where the hit shown is. */
  private report() {
    const found = this.found, hit = found?.hits[this.hit];
    if (!found?.query.trim()) { this.bar?.result(); return; }
    const paged = this.view?.kind !== "sheet" && this.view?.kind !== "scroll";
    const detail = hit && `${paged ? `page ${hit.segments[0].a + 1}: ` : ""}${hit.context.replace(/ ⏎ /g, " ")}`;
    this.bar?.result({ index: this.hit, total: found.length, more: found.more, detail });
  }

  /** Draw the hits anew on a sheet and on the pages and bands that show their text. */
  private markHits() {
    this.redrawSheet?.();
    for (const el of this.stage.querySelectorAll<HTMLElement>(".bdfMini-page")) {
      if (el.querySelector(".bdfMini-hits")) this.markPage(el);
    }
  }

  /** Draw the hits on a page or a band, over its text; those not located yet are drawn once they are. */
  private markPage(el: HTMLElement) {
    const v = this.view, found = this.found;
    if (!v) return;
    const draw = () => {
      const layer = this.hitLayer(v, el), old = el.querySelector(".bdfMini-hits");
      if (old) old.replaceWith(layer); else el.append(layer);
    };
    draw();
    if (!found) return;
    let places: number[];
    if (el.dataset.y === undefined) places = [Number(el.dataset.index)];
    else {
      // the strips a band shows
      const y = Number(el.dataset.y), { tops } = this.column(v);
      places = tops.map((_, i) => i).filter((i) => tops[i] < y + BAND && (tops[i + 1] ?? Infinity) > y);
    }
    found.locate(places).then((some) => {
      if (some && found === this.found && el.isConnected) draw();
    }).catch(this.failed);
  }

  /** Let the worker and the document go. */
  destroy() {
    this.generation++;
    this.dropAudio();
    this.resize.disconnect();
    this.uninstallCopy();
    this.keys?.removeEventListener("keydown", this.onKey as EventListener);
    this.client.terminate();
    this.host.replaceChildren();
  }

  /** Drawing failed: the page stays blank, and the console says why. */
  private failed = (e: unknown) => {
    (this.options.onError ?? console.error)(e);
  };

  private clearSearch() {
    this.turn++;
    this.found = undefined;
    this.hit = -1;
    this.bar?.result();
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
        if (!live()) return;
        el.append(buildTextLayer(content, this.scale, { lang: this.language }));
        this.markPage(el);
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
        if (!live()) return;
        el.append(buildTextLayer(content, this.scale, { lang: this.language }));
        this.markPage(el);
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

  /** The located hits on a page or a band, to go over its text. */
  private hitLayer(v: View, el: HTMLElement): HTMLDivElement {
    const layer = document.createElement("div");
    layer.className = "bdfMini-layer bdfMini-hits";
    const band = el.dataset.y === undefined ? undefined : Number(el.dataset.y);
    this.found?.each((r, i) => {
      let box: Box = r;
      if (band === undefined) {
        if (r.a !== Number(el.dataset.index)) return;
      } else {
        box = this.inColumn(v, r);
        if (box.y + box.h < band || box.y > band + BAND) return;
        box = { ...box, y: box.y - band };
      }
      const d = document.createElement("div");
      d.className = i === this.hit ? "bdfMini-hit current" : "bdfMini-hit";
      d.style.cssText = `left: ${box.x * this.scale}px; top: ${box.y * this.scale}px; width: ${box.w * this.scale}px; height: ${box.h * this.scale}px`;
      layer.append(d);
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
        this.found?.each((r, i) => {
          ctx.fillStyle = i === this.hit ? "rgba(255,120,0,.5)" : "rgba(255,210,0,.45)";
          ctx.fillRect(p.dx + (r.x - p.src.x) * k, p.dy + (r.y - p.src.y) * k, r.w * k, r.h * k);
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

    /** Locate the hits of the tiles in view, and draw them once they are. */
    const mark = () => {
      const found = this.found, tile = tileSize(v), tiles: string[] = [];
      if (!found?.length) return;
      for (const { src } of panes) {
        if (src.w <= 0 || src.h <= 0) continue;
        for (let ty = Math.floor(src.y / tile); ty * tile < src.y + src.h; ty++) {
          for (let tx = Math.floor(src.x / tile); tx * tile < src.x + src.w; tx++) tiles.push(`${tx},${ty}`);
        }
      }
      found.locate(tiles).then((some) => {
        if (some && found === this.found && gen === this.generation) paint();
      }).catch(this.failed);
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
      mark();
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
    this.redrawSheet = () => {
      paint();
      mark();
    };
    draw();
  }

  /** Links to a page ("#page=N") scroll there; links to a view ("#view=ID") show it. */
  /** Follow a link to a page or a view; false when the click is on none. */
  private followLink(e: MouseEvent): boolean {
    const a = (e.target as Element).closest(`a[${RUN_ATTR.page}], a[${RUN_ATTR.view}]`);
    if (!a || !this.view) return false;
    e.preventDefault();
    const page = Number(a.getAttribute(RUN_ATTR.page) ?? "1") - 1;
    const id = a.getAttribute(RUN_ATTR.view);
    if (id !== null) this.show(id, page);
    else this.goToPage(page);
    return true;
  }
}

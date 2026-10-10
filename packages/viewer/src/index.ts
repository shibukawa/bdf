import { tileSize, type Manifest, type View } from "@bdfkit/core";
import { createRenderWorker, FindBar, FindHits, findKey, type BdfWorkerClient, type FindBarLabels, type HitRect, type OpenSource } from "@bdfkit/render";

export type ViewerMode = "embedded" | "lightbox";
export type PageLayout = "single" | "spread" | "continuous";
export type MusicMode = "score" | "piano-roll";
export interface MusicRenderContext {
  stage: HTMLDivElement;
  renderer: BdfWorkerClient;
  view: View;
  zoom: number;
  isCurrent(): boolean;
}
export type MusicRenderer = (context: MusicRenderContext) => Promise<void>;

export interface ViewerOptions {
  mode?: ViewerMode;
  layout?: PageLayout;
  musicMode?: MusicMode;
  zoom?: number;
  /** Supply a renderer to share or configure it. The viewer does not terminate a supplied renderer. */
  renderer?: BdfWorkerClient;
  /** Optional piano-roll renderer from @bdfkit/viewer/music. Omit it to leave music code out of the app bundle. */
  music?: MusicRenderer;
  /**
   * Whose the browser's find shortcut is (Cmd+F on macOS, Ctrl+F elsewhere).
   * Pages are drawn as bitmaps, in which the browser's own search finds
   * nothing; the viewer's find bar searches the text of the whole view,
   * the pages not drawn yet too. "focus": the viewer's while the focus is
   * in it (the default when embedded); "page": the viewer's wherever the
   * focus is while the viewer is shown (the default for a lightbox);
   * false: the browser's, and no find bar: search with find() from the
   * host's own controls.
   */
  find?: "focus" | "page" | false;
  /** The words of the find bar (default: English). */
  findLabels?: Partial<FindBarLabels>;
}

/** What a search found: the detail of the "findchange" event. */
export interface FindState {
  query: string;
  /** The match shown (from 0), -1 when there is none. */
  index: number;
  total: number;
  /** The view has more matches than total. */
  more: boolean;
}

export interface FileOptions { name?: string; password?: string; fonts?: string; format?: string; pages?: string; params?: Record<string, string> }
export interface FileConverter {
  convert(data: ArrayBuffer | Uint8Array, options?: FileOptions): Promise<{ bdf: Uint8Array }>;
  open?(data: ArrayBuffer | Uint8Array, options?: FileOptions): Promise<{ bdf: Uint8Array; pages: number; stream?: number }>;
  page?(stream: number, index: number): Promise<{ bdf: Uint8Array }>;
  finish?(stream: number): Promise<{ bdf: Uint8Array }>;
  close?(stream: number): Promise<unknown>;
}

/**
 * A document surface: page placement, scrolling and rendering. The host owns
 * its toolbar, file picker, view tabs and lightbox close control.
 */
export class Viewer extends EventTarget {
  readonly stage: HTMLDivElement;
  readonly surface: HTMLDivElement;
  readonly renderer: BdfWorkerClient;
  private owned: boolean;
  private observer?: IntersectionObserver;
  private resize?: ResizeObserver;
  private generation = 0;
  private sheetFrame = 0;
  private documentSerial = 0;
  private activeStream?: { converter: FileConverter; id: number };
  private manifest?: Manifest;
  private current?: View;
  private layout: PageLayout;
  private musicMode: MusicMode;
  private mode: ViewerMode;
  private zoom: number;
  private musicRenderer?: MusicRenderer;
  /** The hits of the search shown (they are located as their pages are drawn), and the one shown. */
  private found?: FindHits;
  private hit = -1;
  /** Bumped by each search: a hit located for an earlier one is not gone to. */
  private findTurn = 0;
  /** Pages came (a streamed conversion) since the search: it is run again before it goes on. */
  private findStale = false;
  /** The document's pages come as they are read (segments): a search covers those that came. */
  private partial = false;
  private bar?: FindBar;
  private findKeys?: HTMLElement | Document;
  /** The page frames near the visible area: their hits are drawn. */
  private near = new Set<HTMLElement>();
  private sheetMark?: () => void;

  constructor(readonly host: HTMLElement, options: ViewerOptions = {}) {
    super();
    this.mode = options.mode ?? "embedded";
    this.layout = options.layout ?? "single";
    this.musicMode = options.musicMode ?? "score";
    this.zoom = options.zoom ?? 1;
    this.musicRenderer = options.music;
    if (this.musicMode === "piano-roll" && !this.musicRenderer) throw new Error("piano roll requires @bdfkit/viewer/music");
    this.renderer = options.renderer ?? createRenderWorker();
    this.owned = !options.renderer;
    this.surface = document.createElement("div");
    this.surface.style.cssText = this.mode === "lightbox"
      ? "position:fixed;inset:0;z-index:2147483647;background:#000b;display:none;padding:3vmin;box-sizing:border-box"
      : "width:100%;height:100%;min-height:0;box-sizing:border-box";
    this.stage = document.createElement("div");
    this.stage.style.cssText = "width:100%;height:100%;overflow:auto;box-sizing:border-box;overscroll-behavior:contain";
    this.stage.setAttribute("role", "region");
    this.stage.setAttribute("aria-label", "Document pages");
    // the keys go to the pages once they are clicked: scrolling, and the find shortcut
    this.stage.tabIndex = 0;
    // what the find bar is placed in: the area of the pages
    const box = document.createElement("div");
    box.style.cssText = "position:relative;width:100%;height:100%;min-height:0";
    box.append(this.stage);
    this.surface.append(box);
    host.append(this.surface);
    const find = options.find ?? (this.mode === "lightbox" ? "page" : "focus");
    if (find) {
      this.bar = new FindBar({
        labels: options.findLabels,
        onFind: (query, step) => { void this.find(query, step).catch((error) => this.fail(error)); },
        onClose: () => {
          this.dropFind();
          this.stage.focus({ preventScroll: true });
        },
      });
      box.append(this.bar.element);
      this.findKeys = find === "page" ? document : this.surface;
      this.findKeys.addEventListener("keydown", this.onFindKey as EventListener);
    }
    this.resize = new ResizeObserver(() => {
      if (this.current?.kind === "sheet") this.scheduleSheet();
    });
    this.resize.observe(this.stage);
  }

  /** In lightbox mode, the host decides when to open and close the surface. */
  show(): void { if (this.mode === "lightbox") this.surface.style.display = "block"; }
  hide(): void { if (this.mode === "lightbox") this.surface.style.display = "none"; }

  async open(source: OpenSource, password?: string): Promise<Manifest> {
    ++this.documentSerial;
    this.cancelStream();
    const generation = ++this.generation;
    this.clear();
    this.dropFind();
    this.bar?.reset();
    await this.renderer.close();
    this.partial = source.kind === "segments";
    const manifest = await this.renderer.open(source, password);
    if (generation !== this.generation) return manifest;
    this.manifest = manifest;
    this.setView(manifest.views[0]?.id);
    this.dispatchEvent(new CustomEvent("documentchange", { detail: manifest }));
    return manifest;
  }

  /** Show the outline immediately for streaming formats, then fill its pages. */
  async openFile(converter: FileConverter, data: ArrayBuffer | Uint8Array, options: FileOptions = {}): Promise<Manifest> {
    const streamed = converter.open && converter.page && converter.finish && converter.close;
    if (!streamed) {
      const result = await converter.convert(data, options);
      return this.open({ kind: "buffer", buffer: toBuffer(result.bdf) });
    }
    const opened = await converter.open!(data, options);
    const manifest = await this.open({ kind: "buffer", buffer: toBuffer(opened.bdf) });
    if (opened.pages > 0 && opened.stream !== undefined) {
      const serial = this.documentSerial;
      this.activeStream = { converter, id: opened.stream };
      void this.fillPages(converter, opened.stream, opened.pages, manifest.views[0]?.id, serial);
    }
    return manifest;
  }

  private async fillPages(converter: FileConverter, stream: number, pages: number, viewId: string | undefined, serial: number): Promise<void> {
    if (!viewId) { this.cancelStream(); return; }
    try {
      for (let index = 0; index < pages; index++) {
        if (serial !== this.documentSerial) return;
        const part = await converter.page!(stream, index);
        if (serial !== this.documentSerial) return;
        await this.renderer.addPage(viewId, index, toBuffer(part.bdf));
        // the text of the page came with it
        this.findStale = true;
        const frame = this.stage.querySelector<HTMLElement>(`[data-bdf-page="${index}"]`);
        const canvas = frame?.querySelector("canvas");
        if (canvas) {
          canvas.dataset.loaded = "";
          const area = this.stage.getBoundingClientRect(), rect = frame!.getBoundingClientRect();
          if (this.current?.id === viewId && rect.bottom >= area.top - 600 && rect.top <= area.bottom + 600) {
            void this.renderPage(this.current, index, canvas, this.generation);
          }
        }
        this.dispatchEvent(new CustomEvent("conversionprogress", { detail: { page: index, pages } }));
      }
      if (serial !== this.documentSerial) return;
      const whole = await converter.finish!(stream);
      if (serial !== this.documentSerial) return;
      const manifest = await this.renderer.replace({ kind: "buffer", buffer: toBuffer(whole.bdf) });
      if (serial === this.documentSerial) {
        this.manifest = manifest;
        this.current = manifest.views.find((view) => view.id === this.current?.id);
        this.dispatchEvent(new CustomEvent("conversiondone", { detail: manifest }));
        // a search that ran meanwhile covered the pages converted by then: its count and marks follow, the pages stay where they are
        if (this.found?.query && this.findStale) void this.runFind(this.found.query, 0, false).catch((error) => this.fail(error));
      }
    } catch (error) {
      if (serial === this.documentSerial) this.dispatchEvent(new CustomEvent("error", { detail: error }));
    } finally {
      if (this.activeStream?.id === stream) this.activeStream = undefined;
      void converter.close!(stream).catch(() => {});
    }
  }

  private cancelStream(): void {
    if (!this.activeStream) return;
    const { converter, id } = this.activeStream;
    this.activeStream = undefined;
    void converter.close?.(id)?.catch(() => {});
  }

  get document(): Manifest | undefined { return this.manifest; }
  get view(): View | undefined { return this.current; }
  get pageLayout(): PageLayout { return this.layout; }

  setView(id?: string): void {
    const view = this.manifest?.views.find((item) => item.id === id);
    const other = view?.id !== this.current?.id;
    if (other) this.dropFind();
    if (!view) { this.clear(); this.current = undefined; return; }
    this.current = view;
    this.draw();
    this.dispatchEvent(new CustomEvent("viewchange", { detail: view }));
    // the find bar goes on with its query in the view that came
    if (other && this.bar?.shown && this.bar.query) void this.find(this.bar.query, 0).catch((error) => this.fail(error));
  }
  setLayout(layout: PageLayout): void { this.layout = layout; this.draw(); }
  setMusicMode(mode: MusicMode): void {
    if (mode === "piano-roll" && !this.musicRenderer) throw new Error("piano roll requires @bdfkit/viewer/music");
    this.musicMode = mode;
    if (this.current && !this.searchable(this.current)) this.clearFind();
    this.draw();
  }
  setZoom(zoom: number): void {
    if (!Number.isFinite(zoom) || zoom <= 0 || zoom > 8) throw new RangeError("zoom must be in (0, 8]");
    this.zoom = zoom;
    this.draw();
  }

  scrollToPage(index: number): void {
    this.stage.querySelector<HTMLElement>(`[data-bdf-page="${index}"]`)?.scrollIntoView({ block: "start" });
  }

  /**
   * Search the text of the view shown, all of it, and show a match: of
   * another query, the first one on the page being read or after it (step
   * -1: the last one before it); of the same query again, the next one
   * (step -1: the one before; 0: the same one). The matches are marked on
   * the pages and the one shown is scrolled to. An empty query clears the
   * search. A "findchange" event carries what was found, as the result does.
   */
  find(query: string, step: -1 | 0 | 1 = 1): Promise<FindState> {
    return this.runFind(query, step, true);
  }

  private async runFind(query: string, step: -1 | 0 | 1, reveal: boolean): Promise<FindState> {
    const view = this.current;
    if (!view || !this.searchable(view)) return { query, index: -1, total: 0, more: false };
    const turn = ++this.findTurn;
    let found = this.found;
    // a document whose pages come as they are read is searched anew each time: more of them may have come
    if (!found || query !== found.query || this.findStale || this.partial) {
      // the same query over more pages goes on from the match it was at
      const at = found?.query === query ? found.hits[this.hit]?.segments[0] : undefined;
      this.findStale = false;
      found = await FindHits.search(this.renderer, view, query);
      // another search, view or document came meanwhile
      if (turn !== this.findTurn) return this.findState(found, -1);
      this.found = found;
      const n = found.length, start = found.from(this.reading());
      const same = at ? found.hits.findIndex(({ segments: [s] }) => s.a === at.a && s.b === at.b && s.ordinal === at.ordinal && s.start === at.start) : -1;
      this.hit = !n ? -1 : same >= 0 ? (same + step + n) % n : step < 0 ? (start - 1 + n) % n : start;
    } else if (found.length) {
      this.hit = (this.hit + step + found.length) % found.length;
    }
    const state = this.findState(found, this.hit);
    const hit = found.hits[this.hit];
    this.bar?.result(query.trim() ? {
      ...state,
      detail: hit && `${view.kind === "sheet" || view.kind === "scroll" ? "" : `page ${hit.segments[0].a + 1}: `}${hit.context.replace(/ ⏎ /g, " ")}`,
    } : undefined);
    this.markHits();
    this.dispatchEvent(new CustomEvent("findchange", { detail: state }));
    // where the match is on its page is known once the page is read
    const rect = this.hit < 0 || !reveal ? undefined : (await found.locateHit(this.hit))[0];
    if (turn !== this.findTurn || !rect) return state;
    this.markHits();
    this.reveal(view, rect);
    return state;
  }

  /** Show the find bar, as the find shortcut does (nothing with find: false, or when there is no text to search). */
  openFind(): void {
    const bar = this.bar;
    if (!bar || !this.current || !this.searchable(this.current)) return;
    bar.open();
    // the bar opens with the query it had: its matches are shown again
    if (bar.query && !this.found) void this.find(bar.query, 0).catch((error) => this.fail(error));
  }

  /** End the search: the matches are no longer marked, and the find bar closes. */
  clearFind(): void {
    this.dropFind();
    this.bar?.reset();
  }

  private findState(found: FindHits, index: number): FindState {
    return { query: found.query, index, total: found.length, more: found.more };
  }

  /** Pages have text to search; a piano roll (see draw) has none. */
  private searchable(view: View): boolean {
    return !(this.musicMode === "piano-roll" && view.play?.seq && this.musicRenderer);
  }

  private fail(error: unknown): void {
    this.dispatchEvent(new CustomEvent("error", { detail: error }));
  }

  /** Forget the search and its marks, and tell the host when there was one; the find bar stays as it is. */
  private dropFind(): void {
    const searched = !!this.found?.query;
    ++this.findTurn;
    this.found = undefined;
    this.hit = -1;
    this.findStale = false;
    this.bar?.result();
    this.markHits();
    if (searched) this.dispatchEvent(new CustomEvent("findchange", { detail: { query: "", index: -1, total: 0, more: false } satisfies FindState }));
  }

  /** The find shortcuts: the bar opens, and while it is shown the next and previous match keys are its own. */
  private onFindKey = (e: KeyboardEvent): void => {
    const bar = this.bar, key = findKey(e);
    if (!bar || !key || e.defaultPrevented || !this.current || !this.searchable(this.current)) return;
    // a lightbox that is not shown has no part in the page
    if (this.mode === "lightbox" && this.surface.style.display === "none") return;
    if (key === "open") {
      // pressed again in the field, the shortcut is the browser's: its own search stays at hand
      if (bar.shown && bar.focused) return;
      e.preventDefault();
      this.openFind();
    } else if (bar.shown && bar.query) {
      e.preventDefault();
      void this.find(bar.query, key === "next" ? 1 : -1).catch((error) => this.fail(error));
    }
  };

  /** Where a search starts: the first page that reaches into the visible area. */
  private reading(): number {
    const top = this.stage.getBoundingClientRect().top;
    let page = Infinity;
    for (const frame of this.near) {
      if (frame.getBoundingClientRect().bottom > top) page = Math.min(page, Number(frame.dataset.bdfPage));
    }
    return Number.isFinite(page) ? page : 0;
  }

  /** Mark the matches anew on the pages near the visible area, or on the sheet. */
  private markHits(): void {
    for (const frame of this.near) this.markFrame(frame);
    this.sheetMark?.();
  }

  private hitBox(rect: HitRect, current: boolean): HTMLDivElement {
    const box = document.createElement("div");
    box.style.cssText = `position:absolute;left:${rect.x * this.zoom}px;top:${rect.y * this.zoom}px;width:${rect.w * this.zoom}px;height:${rect.h * this.zoom}px;border-radius:2px;`
      + (current ? "background:rgba(255,120,0,.5);outline:1px solid #c04000" : "background:rgba(255,210,0,.45);outline:1px solid rgba(200,140,0,.8)");
    return box;
  }

  /** A layer of marks over a page or a sheet; the colors stay as they are in forced colors, where a background would hide the text. */
  private hitLayer(): HTMLDivElement {
    const layer = document.createElement("div");
    layer.dataset.bdfHits = "";
    layer.style.cssText = "position:absolute;left:0;top:0;width:100%;height:100%;pointer-events:none;forced-color-adjust:none";
    return layer;
  }

  /** Mark the matches on a page; those not located yet are marked once they are. */
  private markFrame(frame: HTMLElement): void {
    frame.querySelector("[data-bdf-hits]")?.remove();
    const found = this.found;
    if (!found?.length) return;
    const page = Number(frame.dataset.bdfPage);
    const layer = this.hitLayer();
    found.each((rect, i) => { if (rect.a === page) layer.append(this.hitBox(rect, i === this.hit)); });
    if (layer.childElementCount) frame.append(layer);
    found.locate([page]).then((some) => {
      if (some && found === this.found && this.near.has(frame)) this.markFrame(frame);
    }).catch((error) => this.fail(error));
  }

  /** Bring a located match into view. */
  private reveal(view: View, rect: HitRect): void {
    const stage = this.stage;
    if (view.kind === "sheet") {
      stage.scrollTo({ left: Math.max(0, rect.x * this.zoom - stage.clientWidth / 2), top: Math.max(0, rect.y * this.zoom - stage.clientHeight / 2) });
      return;
    }
    const frame = stage.querySelector<HTMLElement>(`[data-bdf-page="${rect.a}"]`);
    if (!frame) return;
    const area = stage.getBoundingClientRect(), at = frame.getBoundingClientRect();
    const x = at.left - area.left + rect.x * this.zoom, y = at.top - area.top + rect.y * this.zoom;
    // sideways only when the match is out of sight (a page wider than the viewer)
    const aside = x < 0 || x + rect.w * this.zoom > stage.clientWidth;
    stage.scrollTo({
      top: Math.max(0, stage.scrollTop + y - stage.clientHeight / 3),
      left: aside ? Math.max(0, stage.scrollLeft + x - stage.clientWidth / 2) : stage.scrollLeft,
    });
  }

  private clear(): void {
    this.observer?.disconnect();
    this.observer = undefined;
    cancelAnimationFrame(this.sheetFrame);
    this.stage.onscroll = null;
    this.sheetDraw = undefined;
    this.sheetMark = undefined;
    this.near.clear();
    this.stage.replaceChildren();
  }

  private draw(): void {
    this.clear();
    const view = this.current;
    if (!view) return;
    const generation = ++this.generation;
    if (view.kind === "sheet") this.drawSheet(view, generation);
    else if (this.musicMode === "piano-roll" && view.play?.seq && this.musicRenderer) {
      void this.musicRenderer({ stage: this.stage, renderer: this.renderer, view, zoom: this.zoom, isCurrent: () => generation === this.generation })
        .catch((error) => this.dispatchEvent(new CustomEvent("error", { detail: error })));
    }
    else this.drawPages(view, generation);
  }

  private drawPages(view: View, generation: number): void {
    const pages = view.pages ?? [];
    const spread = this.layout === "spread" && view.kind !== "scroll";
    const gap = this.layout === "continuous" || view.kind === "scroll" ? 0 : 16;
    this.stage.style.background = "#e8e8e8";
    this.observer = new IntersectionObserver((entries) => {
      for (const entry of entries) {
        const canvas = entry.target.querySelector("canvas") as HTMLCanvasElement;
        const frame = entry.target as HTMLElement;
        const index = Number(frame.dataset.bdfPage);
        if (entry.isIntersecting) {
          void this.renderPage(view, index, canvas, generation);
          if (!this.near.has(frame)) {
            this.near.add(frame);
            this.markFrame(frame);
          }
        } else {
          if (canvas.width) { canvas.width = canvas.height = 0; canvas.dataset.loaded = ""; }
          this.near.delete(frame);
          frame.querySelector("[data-bdf-hits]")?.remove();
        }
      }
    }, { root: this.stage, rootMargin: "600px" });
    let row: HTMLDivElement | undefined;
    pages.forEach((page, index) => {
      if (!spread || index % 2 === 0) {
        row = document.createElement("div");
        // as wide as its pages when they are wider than the viewer: centered pages that overflow could not be scrolled to on the left
        row.style.cssText = `display:flex;justify-content:center;align-items:flex-start;gap:${gap}px;margin-bottom:${gap}px;width:max-content;min-width:100%;`;
        if (spread && view.direction === "rtl") row.style.flexDirection = "row-reverse";
        this.stage.append(row);
      }
      const frame = document.createElement("div");
      frame.dataset.bdfPage = String(index);
      frame.style.cssText = `position:relative;flex:none;width:${page.w * this.zoom}px;height:${page.h * this.zoom}px;background:white;`;
      const canvas = document.createElement("canvas");
      canvas.style.cssText = "display:block;width:100%;height:100%";
      frame.append(canvas);
      row!.append(frame);
      this.observer!.observe(frame);
    });
  }

  private async renderPage(view: View, index: number, canvas: HTMLCanvasElement, generation: number): Promise<void> {
    if (canvas.dataset.loaded) return;
    canvas.dataset.loaded = "loading";
    const request = String(Number(canvas.dataset.renderRequest ?? "0") + 1);
    canvas.dataset.renderRequest = request;
    try {
      const bitmap = await this.renderer.page(view.id, index, this.zoom * devicePixelRatio);
      if (generation !== this.generation || !canvas.isConnected || canvas.dataset.renderRequest !== request) { bitmap.close(); return; }
      canvas.width = bitmap.width;
      canvas.height = bitmap.height;
      canvas.getContext("2d")!.drawImage(bitmap, 0, 0);
      bitmap.close();
      canvas.dataset.loaded = "ready";
    } catch (error) {
      canvas.dataset.loaded = "";
      this.dispatchEvent(new CustomEvent("error", { detail: error }));
    }
  }

  private drawSheet(view: View, generation: number): void {
    this.stage.style.background = "white";
    const extent = (runs: [number, number][] = []) => runs.reduce((sum, [count, size]) => sum + count * size, 0);
    const spacer = document.createElement("div");
    spacer.style.cssText = `position:relative;width:${Math.max(this.stage.clientWidth, extent(view.cols) * this.zoom)}px;height:${Math.max(this.stage.clientHeight, extent(view.rows) * this.zoom)}px;`;
    const canvas = document.createElement("canvas");
    canvas.style.cssText = "position:sticky;left:0;top:0;display:block";
    spacer.append(canvas);
    this.stage.append(spacer);
    this.stage.onscroll = () => this.scheduleSheet();
    // the matches of the sheet, over its cells; those of the tiles in view are located as they are scrolled to
    const mark = () => {
      spacer.querySelector("[data-bdf-hits]")?.remove();
      const found = this.found;
      if (!found?.length) return;
      const layer = this.hitLayer();
      found.each((rect, i) => layer.append(this.hitBox(rect, i === this.hit)));
      if (layer.childElementCount) spacer.append(layer);
    };
    const locate = () => {
      const found = this.found;
      if (!found?.length) return;
      const tile = tileSize(view), tiles: string[] = [];
      const x0 = this.stage.scrollLeft / this.zoom, y0 = this.stage.scrollTop / this.zoom;
      const x1 = x0 + this.stage.clientWidth / this.zoom, y1 = y0 + this.stage.clientHeight / this.zoom;
      for (let ty = Math.floor(y0 / tile); ty * tile < y1; ty++) {
        for (let tx = Math.floor(x0 / tile); tx * tile < x1; tx++) tiles.push(`${tx},${ty}`);
      }
      found.locate(tiles).then((some) => {
        if (some && found === this.found && generation === this.generation) mark();
      }).catch((error) => this.fail(error));
    };
    this.sheetMark = () => { mark(); locate(); };
    let lastRequest = 0;
    const draw = async () => {
      locate();
      const request = ++lastRequest;
      const width = this.stage.clientWidth, height = this.stage.clientHeight;
      if (width <= 0 || height <= 0) return;
      const viewport = { x: this.stage.scrollLeft / this.zoom, y: this.stage.scrollTop / this.zoom, w: width / this.zoom, h: height / this.zoom };
      try {
        const bitmap = await this.renderer.sheet(view.id, viewport, this.zoom * devicePixelRatio);
        if (generation !== this.generation || request !== lastRequest) { bitmap.close(); return; }
        canvas.width = bitmap.width;
        canvas.height = bitmap.height;
        canvas.style.width = `${width}px`;
        canvas.style.height = `${height}px`;
        canvas.getContext("2d")!.drawImage(bitmap, 0, 0);
        bitmap.close();
      } catch (error) { this.dispatchEvent(new CustomEvent("error", { detail: error })); }
    };
    this.sheetDraw = draw;
    mark();
    this.scheduleSheet();
  }
  private sheetDraw?: () => Promise<void>;
  private scheduleSheet(): void {
    if (this.sheetFrame || !this.sheetDraw) return;
    this.sheetFrame = requestAnimationFrame(() => { this.sheetFrame = 0; void this.sheetDraw?.(); });
  }

  destroy(): void {
    ++this.documentSerial;
    this.cancelStream();
    ++this.generation;
    this.clear();
    this.resize?.disconnect();
    this.findKeys?.removeEventListener("keydown", this.onFindKey as EventListener);
    this.surface.remove();
    if (this.owned) this.renderer.terminate();
  }
}

function toBuffer(bytes: Uint8Array): ArrayBuffer {
  return bytes.buffer.slice(bytes.byteOffset, bytes.byteOffset + bytes.byteLength) as ArrayBuffer;
}

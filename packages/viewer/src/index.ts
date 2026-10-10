import { type Manifest, type View } from "@bdfkit/core";
import { createRenderWorker, type BdfWorkerClient, type OpenSource } from "@bdfkit/render";

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
    this.surface.append(this.stage);
    host.append(this.surface);
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
    await this.renderer.close();
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
    if (!view) { this.clear(); this.current = undefined; return; }
    this.current = view;
    this.draw();
    this.dispatchEvent(new CustomEvent("viewchange", { detail: view }));
  }
  setLayout(layout: PageLayout): void { this.layout = layout; this.draw(); }
  setMusicMode(mode: MusicMode): void {
    if (mode === "piano-roll" && !this.musicRenderer) throw new Error("piano roll requires @bdfkit/viewer/music");
    this.musicMode = mode;
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

  private clear(): void {
    this.observer?.disconnect();
    this.observer = undefined;
    cancelAnimationFrame(this.sheetFrame);
    this.stage.onscroll = null;
    this.sheetDraw = undefined;
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
        const index = Number((entry.target as HTMLElement).dataset.bdfPage);
        if (entry.isIntersecting) void this.renderPage(view, index, canvas, generation);
        else if (canvas.width) { canvas.width = canvas.height = 0; canvas.dataset.loaded = ""; }
      }
    }, { root: this.stage, rootMargin: "600px" });
    let row: HTMLDivElement | undefined;
    pages.forEach((page, index) => {
      if (!spread || index % 2 === 0) {
        row = document.createElement("div");
        row.style.cssText = `display:flex;justify-content:center;align-items:flex-start;gap:${gap}px;margin-bottom:${gap}px;`;
        if (spread && view.direction === "rtl") row.style.flexDirection = "row-reverse";
        this.stage.append(row);
      }
      const frame = document.createElement("div");
      frame.dataset.bdfPage = String(index);
      frame.style.cssText = `flex:none;width:${page.w * this.zoom}px;height:${page.h * this.zoom}px;background:white;`;
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
    spacer.style.cssText = `width:${Math.max(this.stage.clientWidth, extent(view.cols) * this.zoom)}px;height:${Math.max(this.stage.clientHeight, extent(view.rows) * this.zoom)}px;`;
    const canvas = document.createElement("canvas");
    canvas.style.cssText = "position:sticky;left:0;top:0;display:block";
    spacer.append(canvas);
    this.stage.append(spacer);
    this.stage.onscroll = () => this.scheduleSheet();
    let lastRequest = 0;
    const draw = async () => {
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
    this.surface.remove();
    if (this.owned) this.renderer.terminate();
  }
}

function toBuffer(bytes: Uint8Array): ArrayBuffer {
  return bytes.buffer.slice(bytes.byteOffset, bytes.byteOffset + bytes.byteLength) as ArrayBuffer;
}

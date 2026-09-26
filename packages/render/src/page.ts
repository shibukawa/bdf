import type { BdfDocument, View, Page, Rect } from "@bdf/core";
import { ResourceCache } from "./resources.js";
import { CanvasRenderer, type Ctx2D, type RenderOptions } from "./canvas.js";

export interface PageRenderOptions extends RenderOptions {
  /** Device pixels per unit. */
  scale: number;
  /** Only draw layers whose role is in this set (default: all). */
  roles?: Iterable<string>;
  /** Fill color behind the page (default white; null for transparent). */
  background?: string | null;
}

/** Renders pages, continuous flow regions and sheet regions of one document. */
export class PageRenderer {
  readonly res: ResourceCache;
  readonly renderer: CanvasRenderer;

  constructor(readonly doc: BdfDocument, opts: RenderOptions = {}, fontSet?: FontFaceSet) {
    this.res = new ResourceCache(doc, fontSet);
    this.renderer = new CanvasRenderer(this.res, opts);
  }

  /** Let the fonts and images go: the document is not drawn any more. */
  dispose(): void {
    this.res.dispose();
  }

  /** Load everything a page needs. */
  async preparePage(page: Page): Promise<void> {
    await Promise.all(page.layers.map((l) => this.res.prepare(l.obj)));
  }

  /**
   * Draw one page into ctx. The canvas must be at least page.w*scale by page.h*scale;
   * the page origin is placed at (dx, dy) device pixels.
   */
  async renderPage(ctx: Ctx2D, page: Page, opts: PageRenderOptions, dx = 0, dy = 0): Promise<void> {
    await this.preparePage(page);
    this.drawPageSync(ctx, page, opts, dx, dy);
  }

  drawPageSync(ctx: Ctx2D, page: Page, opts: PageRenderOptions, dx = 0, dy = 0): void {
    const roles = opts.roles ? new Set(opts.roles) : undefined;
    ctx.save();
    ctx.setTransform(opts.scale, 0, 0, opts.scale, dx, dy);
    ctx.beginPath();
    ctx.rect(0, 0, page.w, page.h);
    ctx.clip();
    if (opts.background !== null) {
      ctx.fillStyle = opts.background ?? "#ffffff";
      ctx.fillRect(0, 0, page.w, page.h);
    }
    for (const layer of page.layers) {
      if (roles && !roles.has(layer.role)) continue;
      this.renderer.draw(ctx, this.res.object(layer.obj));
    }
    ctx.restore();
  }

  /**
   * Continuous layout of a flow view: body rectangles stacked vertically
   * (for a scroll view, its strips, without gaps).
   * Returns the y offset of each page's body and the total height, in units.
   */
  continuousLayout(view: View): { offsets: number[]; width: number; height: number } {
    const gap = view.kind === "scroll" ? 0 : view.continuous?.gap ?? 0;
    const offsets: number[] = [];
    let y = 0;
    let width = 0;
    for (const p of view.pages ?? []) {
      const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
      offsets.push(y);
      y += b.h + gap;
      width = Math.max(width, b.w);
    }
    return { offsets, width, height: Math.max(0, y - gap) };
  }

  /**
   * Draw the part of the continuous layout that intersects viewport (units) into ctx,
   * with viewport's top-left at device (0,0).
   */
  async renderContinuous(ctx: Ctx2D, view: View, viewport: Rect, opts: PageRenderOptions): Promise<void> {
    const pages = view.pages ?? [];
    const { offsets } = this.continuousLayout(view);
    const roles = opts.roles ?? ["body", "annotation"];
    const visible: number[] = [];
    pages.forEach((p, i) => {
      const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
      if (offsets[i] < viewport.y + viewport.h && offsets[i] + b.h > viewport.y) visible.push(i);
    });
    await Promise.all(visible.map((i) => this.preparePage(pages[i])));
    for (const i of visible) {
      const p = pages[i];
      const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
      ctx.save();
      ctx.setTransform(opts.scale, 0, 0, opts.scale, -viewport.x * opts.scale, (offsets[i] - viewport.y) * opts.scale);
      ctx.beginPath();
      ctx.rect(0, 0, b.w, b.h);
      ctx.clip();
      if (opts.background !== null) {
        ctx.fillStyle = opts.background ?? "#ffffff";
        ctx.fillRect(0, 0, b.w, b.h);
      }
      ctx.translate(-b.x, -b.y);
      const roleSet = new Set(roles);
      for (const layer of p.layers) {
        if (!roleSet.has(layer.role)) continue;
        this.renderer.draw(ctx, this.res.object(layer.obj));
      }
      ctx.restore();
    }
  }

  /** Total size of a sheet in units, from its row/column runs. */
  sheetSize(view: View): { width: number; height: number } {
    const sum = (runs?: [number, number][]) => (runs ?? []).reduce((acc, [n, s]) => acc + n * s, 0);
    return { width: sum(view.cols), height: sum(view.rows) };
  }

  /**
   * Draw the region of a sheet that intersects viewport (units) into ctx,
   * with viewport's top-left at device (0,0). Draws gridlines when the view asks for them;
   * row/column headers and frozen panes are the viewer's job.
   */
  async renderSheet(ctx: Ctx2D, view: View, viewport: Rect, opts: PageRenderOptions): Promise<void> {
    const tile = view.tile ?? 2048;
    const tiles = view.tiles ?? {};
    const tx0 = Math.floor(viewport.x / tile), ty0 = Math.floor(viewport.y / tile);
    const tx1 = Math.floor((viewport.x + viewport.w - 1e-6) / tile), ty1 = Math.floor((viewport.y + viewport.h - 1e-6) / tile);
    const keys: [number, number, string][] = [];
    for (let ty = ty0; ty <= ty1; ty++) for (let tx = tx0; tx <= tx1; tx++) {
      const h = tiles[`${tx},${ty}`];
      if (h) keys.push([tx, ty, h]);
    }
    await Promise.all(keys.map(([, , h]) => this.res.prepare(h)));

    const s = opts.scale;
    ctx.save();
    ctx.setTransform(s, 0, 0, s, -viewport.x * s, -viewport.y * s);
    if (opts.background !== null) {
      ctx.fillStyle = opts.background ?? "#ffffff";
      ctx.fillRect(viewport.x, viewport.y, viewport.w, viewport.h);
    }
    if (view.gridlines) this.drawGridlines(ctx, view, viewport, s);
    for (const [tx, ty, h] of keys) {
      ctx.save();
      ctx.translate(tx * tile, ty * tile);
      ctx.beginPath();
      ctx.rect(0, 0, tile, tile);
      ctx.clip();
      this.renderer.draw(ctx, this.res.object(h));
      ctx.restore();
    }
    ctx.restore();
  }

  private drawGridlines(ctx: Ctx2D, view: View, vp: Rect, scale: number): void {
    ctx.save();
    ctx.strokeStyle = "#d9d9d9";
    ctx.lineWidth = 1 / scale;
    ctx.beginPath();
    const walkRuns = (runs: [number, number][] | undefined, from: number, to: number, line: (pos: number) => void) => {
      let pos = 0;
      for (const [n, size] of runs ?? []) {
        for (let i = 0; i < n; i++) {
          pos += size;
          if (pos > to) return;
          if (pos >= from) line(pos);
        }
      }
    };
    const { width, height } = this.sheetSize(view);
    const x1 = Math.min(vp.x + vp.w, width), y1 = Math.min(vp.y + vp.h, height);
    const snap = (v: number) => Math.round(v * scale) / scale + 0.5 / scale;
    walkRuns(view.cols, vp.x, x1, (x) => { ctx.moveTo(snap(x), vp.y); ctx.lineTo(snap(x), y1); });
    walkRuns(view.rows, vp.y, y1, (y) => { ctx.moveTo(vp.x, snap(y)); ctx.lineTo(x1, snap(y)); });
    ctx.stroke();
    ctx.restore();
  }
}

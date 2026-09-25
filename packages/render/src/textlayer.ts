// Selectable text over a rendered page (pdf.js style): one transparent,
// absolutely positioned span per text run, so the browser's own selection
// works on top of the bitmap. Copying goes through selectionText(), which
// rebuilds the text from the runs' separators (spaces at line breaks, line
// breaks at paragraph/cell/box boundaries) instead of the browser's
// serialization of unrelated spans.
import { Sep, type TextRun } from "@bdf/core";
import { fontString } from "./resources.js";

export interface TextLayerOptions {
  /** Class of the layer element (default "bdfTextLayer"). */
  className?: string;
  /** Include ALT_TEXT runs (text that was not drawn as glyphs) at their anchor; default false. */
  altText?: boolean;
  /** Width measurement in CSS px for a CSS font string (default: a shared 2D context). */
  measure?: (font: string, text: string) => number;
}

let sharedMeasure: ((font: string, text: string) => number) | undefined;
function defaultMeasure(): (font: string, text: string) => number {
  if (!sharedMeasure) {
    const ctx = document.createElement("canvas").getContext("2d")!;
    sharedMeasure = (font, text) => {
      ctx.font = font;
      return ctx.measureText(text).width;
    };
  }
  return sharedMeasure;
}

/** Attribute names on the run spans. */
export const RUN_ATTR = { ordinal: "data-bdf-ordinal", sep: "data-bdf-sep" } as const;

/**
 * Build the text layer of one page (or band). runs are the extractor's runs
 * for that page, in reading order, in page units; scale is CSS px per unit.
 * The element must be placed over the page (position: absolute; inset: 0)
 * and its spans need `position: absolute; transform-origin: 0 0;
 * white-space: pre; color: transparent` — see TEXT_LAYER_CSS.
 */
export function buildTextLayer(runs: TextRun[], scale: number, opts: TextLayerOptions = {}): HTMLDivElement {
  const layer = document.createElement("div");
  layer.className = opts.className ?? "bdfTextLayer";
  const measure = opts.measure ?? defaultMeasure();
  for (const r of runs) {
    if (!r.text) continue;
    if (r.altText && !opts.altText) continue;
    const span = document.createElement("span");
    span.textContent = r.text;
    span.setAttribute(RUN_ATTR.ordinal, String(r.ordinal));
    span.setAttribute(RUN_ATTR.sep, String(r.sep));
    if (r.font) {
      const font = fontString(r.font, r.size);
      span.style.font = font;
      // The main thread may not have the embedded font: stretch the fallback
      // rendering to the advance the worker measured with the real one.
      const w = measure(font, r.text);
      const sx = r.advance > 0 && w > 0 ? r.advance / w : 1;
      const width = r.advance > 0 ? r.advance : w;
      const anchor = r.align === 1 ? -width : r.align === 2 ? -width / 2 : 0;
      const m = r.matrix;
      span.style.transform = `matrix(${m[0] * scale}, ${m[1] * scale}, ${m[2] * scale}, ${m[3] * scale}, ${r.x * scale}, ${r.y * scale}) translate(${anchor}px, ${-r.size * 0.8}px) scaleX(${sx})`;
    } else {
      // ALT_TEXT without a font: invisible but selectable at its anchor.
      span.style.font = `${Math.max(1, r.size) || 10}px sans-serif`;
      span.style.transform = `translate(${r.x * scale}px, ${r.y * scale}px)`;
    }
    layer.appendChild(span);
  }
  return layer;
}

/** CSS that a text layer needs; the class names follow TextLayerOptions.className. */
export const TEXT_LAYER_CSS = `
.bdfTextLayer { position: absolute; inset: 0; overflow: hidden; line-height: 1; }
.bdfTextLayer span { position: absolute; color: transparent; white-space: pre; transform-origin: 0 0; cursor: text; }
.bdfTextLayer span::selection { background: rgba(0, 120, 255, .35); }
`;

/** A selected piece of a run. */
export interface SelectedRun {
  ordinal: number;
  sep: number;
  text: string;
  /** Element holding the run; layers differ when the selection spans pages. */
  layer: Element | null;
}

/**
 * The runs (or parts of runs) a DOM range covers, in document order: every
 * run span under root that intersects the range, clipped to it.
 */
export function selectedRuns(range: Range, root: ParentNode = document): SelectedRun[] {
  const out: SelectedRun[] = [];
  for (const span of root.querySelectorAll<HTMLElement>(`[${RUN_ATTR.ordinal}]`)) {
    const node = span.firstChild;
    if (!node || !range.intersectsNode(span)) continue;
    const full = node.textContent ?? "";
    let start = 0, end = full.length;
    // The range starts inside this run when its first character lies before the
    // range start; the start offset is only meaningful within the text node.
    if (range.comparePoint(node, 0) < 0) start = range.startContainer === node ? range.startOffset : 0;
    if (range.comparePoint(node, full.length) > 0) end = range.endContainer === node ? range.endOffset : full.length;
    if (end <= start) continue;
    out.push({ ordinal: Number(span.getAttribute(RUN_ATTR.ordinal)), sep: Number(span.getAttribute(RUN_ATTR.sep)), text: full.slice(start, end), layer: span.parentElement });
  }
  return out;
}

/**
 * Join selected runs into plain text: Sep.SPACE becomes a space, Sep.BREAK a
 * line break, and a change of layer (page) a blank line. The first run's
 * own separator is dropped.
 */
export function joinRuns(parts: SelectedRun[]): string {
  let s = "";
  for (let i = 0; i < parts.length; i++) {
    const p = parts[i];
    if (i > 0) {
      if (p.layer !== parts[i - 1].layer) s += "\n\n";
      else if (p.sep === Sep.BREAK) s += "\n";
      else if (p.sep === Sep.SPACE) s += " ";
    }
    s += p.text;
  }
  return s;
}

/** Plain text of the current selection within root (empty when nothing is selected there). */
export function selectionText(sel: Selection | null = typeof getSelection === "function" ? getSelection() : null, root: ParentNode = document): string {
  if (!sel || sel.rangeCount === 0 || sel.isCollapsed) return "";
  let s = "";
  for (let i = 0; i < sel.rangeCount; i++) {
    const t = joinRuns(selectedRuns(sel.getRangeAt(i), root));
    if (t) s = s ? `${s}\n${t}` : t;
  }
  return s;
}

/**
 * Make copy (and cut) inside container put the selection's run text on the
 * clipboard instead of the browser's serialization. Returns a function that
 * removes the handler.
 */
export function installCopyHandler(container: HTMLElement): () => void {
  const onCopy = (e: ClipboardEvent) => {
    const text = selectionText(getSelection(), container);
    if (!text || !e.clipboardData) return;
    e.clipboardData.setData("text/plain", text);
    e.preventDefault();
  };
  container.addEventListener("copy", onCopy);
  container.addEventListener("cut", onCopy);
  return () => {
    container.removeEventListener("copy", onCopy);
    container.removeEventListener("cut", onCopy);
  };
}

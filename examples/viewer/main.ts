// Demo viewer: everything is decoded and rendered in a worker; the main thread
// only places bitmaps and a selectable, accessible text layer. Files opened
// or dropped on the page are converted into bdf in another worker, by the Go
// converters built as wasm (examples/viewer/site.mjs builds them). A PDF is
// converted a page at a time: its pages are shown sized at once and drawn as
// they are converted, those near the visible area first. Pages are scrolled
// through, or shown one or two at a time and turned like a book's (book.ts).
// The cells of a sheet are selected as in a spreadsheet, and copied as
// tab-separated values and an HTML table.
import { dcValues, type Manifest, type View, type SearchHit, type TextContent } from "@bdf/core";
import {
  BdfWorkerClient, BdfWorkerError, buildTextLayer, installCopyHandler, internalLink, tableCells, cellClipboard, TEXT_LAYER_CSS, RUN_ATTR,
  type HitRect, type OpenSource, type TextLayerOptions, type CellText, type CellRange, type CellClipboard,
} from "@bdf/render";
import { ConverterClient, ConvertError, sniff, type Opened } from "./convert.js";
import { Book } from "./book.js";

/** The document shown when the URL has no ?src= (set by the build); "" shows the start page. */
declare const DEFAULT_SRC: string;

const params = new URLSearchParams(location.search);
const src = params.get("src") ?? DEFAULT_SRC;
const client = new BdfWorkerClient(new Worker("./worker.js", { type: "module" }));
/** The converter worker, started with the first file that needs converting. */
let converter: ConverterClient | undefined;
/**
 * Converter modules, one for PDF, one for the Office formats, one for HTML
 * and Markdown and one for the images browsers display by themselves, and
 * the fonts the Office converters lay text out with.
 */
const MODULES = { pdf: "bdf-pdf.wasm", office: "bdf-office.wasm", web: "bdf-web.wasm", image: "bdf-image.wasm" };
const FONTS = "fonts/";

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const stage = $<HTMLDivElement>("stage");
const tabs = $<HTMLSpanElement>("tabs");
const zoomInput = $<HTMLInputElement>("zoom");
const layoutSelect = $<HTMLSelectElement>("layout");
const animateBox = $<HTMLInputElement>("animate");
const status = $<HTMLSpanElement>("status");
const timing = $<HTMLSpanElement>("timing");
const hitsBox = $<HTMLSpanElement>("hits");
const progress = $<HTMLProgressElement>("progress");
/** The start page (kept: the stage drops it when a document opens). */
const landing = $<HTMLDivElement>("landing");

let manifest: Manifest;
/** The view shown; undefined while no document is open. */
let current: View | undefined;
let zoom = 1;
/** Bumped by show(): work started for an earlier view or zoom is dropped. */
let generation = 0;
/** The view shown one or two pages at a time. */
let book: Book | undefined;
/** The reader picked the layout (the menu, or ?layout=); otherwise each document opens in its own. */
let layoutChosen = false;

/** Search state for the current view; pages is how many were converted when it searched (streaming). */
const found = { query: "", hits: [] as SearchHit[], rects: [] as HitRect[][], index: -1, pages: -1 };

/** Where each page of a streamed view is. */
const enum PageState { Pending, Converting, Done, Failed }

/**
 * A conversion running a page at a time: its outline is open, and the pages
 * of its first view come in as the converter gets to them.
 */
interface Streaming {
  /** The converter's stream. */
  id: number;
  /** The open (see opening) that started it. */
  token: number;
  /** The view whose pages come in. */
  view: string;
  state: PageState[];
  /** Pages not converted yet. */
  left: number;
  warnings: string[];
}
/** The conversion the document shown comes from, while pages are still coming. */
let streaming: Streaming | undefined;
/** Pages of the view shown that are in the visible area: converted first. */
const visible = new Set<number>();
let visibility: IntersectionObserver | undefined;
/** Whether a page of a view is still to come. */
const pending = (v: View, index: number) => streaming?.view === v.id && streaming.state[index] !== PageState.Done;
/** Shown sheets: redraw (with the search highlights), scroll a hit into view, and copy the selected cells (true when it did). */
const sheetViews = new WeakMap<View, { redraw: () => void; reveal: (r: { x: number; y: number }) => void; copy: (e: ClipboardEvent) => boolean }>();
/** The cells selected in a sheet: from the active cell, where the selection started, to the cell it reaches (0-based). */
interface CellSelection { anchor: { r: number; c: number }; focus: { r: number; c: number } }
/** The cells selected in each sheet, kept while another view is shown, as a spreadsheet keeps them. */
const cellSelections = new WeakMap<View, CellSelection>();
const dpr = () => window.devicePixelRatio || 1;

/** Announced to screen readers (a polite live region). */
function setStatus(s: string) {
  status.textContent = s;
  status.title = s; // cut to one line
}
/** Render timings: shown only, as they change on every scroll. */
function setTiming(s: string) { timing.textContent = s; }
const showError = (e: unknown) => setStatus(`error: ${(e as Error).message ?? e}`);
/** A failure of work started for a view no longer shown is expected: its document may be gone. */
const unlessStale = (gen: number) => (e: unknown) => { if (gen === generation) showError(e); };

const style = document.createElement("style");
style.textContent = TEXT_LAYER_CSS;
document.head.appendChild(style);

// Copy puts the selected runs' text (with the document's spaces and line
// breaks) on the clipboard, across pages; cells selected in a table as
// tab-separated values and an HTML table. In a sheet, the selected cells.
installCopyHandler(stage);
document.addEventListener("copy", (e) => { if (current) sheetViews.get(current)?.copy(e); });

/**
 * Write cells to the clipboard once they are ready, after the copy event
 * that asked for them: as one item whose data comes later where the
 * browser takes that (Safari takes nothing else after an event), otherwise
 * the text alone.
 */
function writeClipboard(data: Promise<CellClipboard>): Promise<void> {
  if (typeof ClipboardItem !== "undefined" && navigator.clipboard?.write) {
    const blob = (type: string, pick: (d: CellClipboard) => string) => data.then((d) => new Blob([pick(d)], { type }));
    return navigator.clipboard.write([new ClipboardItem({ "text/plain": blob("text/plain", (d) => d.text), "text/html": blob("text/html", (d) => d.html) })]);
  }
  return data.then((d) => navigator.clipboard.writeText(d.text));
}

// Links to a page ("#page=N") scroll there and move the focus to it; links
// to a view ("#view=ID", as between draw.io pages) switch to that view first.
stage.addEventListener("click", (e) => {
  const a = (e.target as Element).closest(`a[${RUN_ATTR.page}], a[${RUN_ATTR.view}]`);
  if (!a) return;
  e.preventDefault();
  const page = a.getAttribute(RUN_ATTR.page);
  const id = a.getAttribute(RUN_ATTR.view);
  if (id !== null) {
    const v = manifest.views.find((v) => v.id === id);
    if (!v) return;
    show(v);
    focusTab(v);
    // the pages are laid out by show: scroll once they are in place
    if (page !== null) requestAnimationFrame(() => goToPage(Number(page) - 1));
    return;
  }
  goToPage(Number(page) - 1);
});

/** Text layer options: the document's language (the UI around it is English). */
const layerOptions = (): TextLayerOptions => ({ lang: dcValues(manifest.meta?.dc?.language)[0] ?? "" });

/** A text layer, which screen readers read ahead of the bitmaps: built over a wider margin. */
const TEXT_MARGIN = "300% 0px";
const BITMAP_MARGIN = "400px";

/** Call render once for each element as it comes near the visible area. */
function onNear(margin: string, render: (el: HTMLDivElement) => void): IntersectionObserver {
  const observer = new IntersectionObserver((entries) => {
    for (const e of entries) {
      if (!e.isIntersecting) continue;
      observer.unobserve(e.target);
      render(e.target as HTMLDivElement);
    }
  }, { root: stage, rootMargin: margin });
  return observer;
}

/** Put a bitmap on a page or band element, under its text layer. */
function placeBitmap(el: HTMLElement, bmp: ImageBitmap, w: number, h: number) {
  const canvas = document.createElement("canvas");
  canvas.width = bmp.width;
  canvas.height = bmp.height;
  canvas.style.width = `${w * zoom}px`;
  canvas.style.height = `${h * zoom}px`;
  canvas.setAttribute("aria-hidden", "true");
  canvas.getContext("2d")!.drawImage(bmp, 0, 0);
  bmp.close();
  el.querySelector("canvas")?.remove();
  el.prepend(canvas);
}

/** Ask for the password of an encrypted document; undefined when the reader cancels. */
function askPassword(message: string, error: boolean): Promise<string | undefined> {
  const dialog = $<HTMLDialogElement>("pwDialog");
  const input = $<HTMLInputElement>("pw");
  const text = $("pwMessage");
  text.textContent = message;
  text.toggleAttribute("data-error", error);
  input.value = "";
  dialog.returnValue = "";
  $("pwCancel").onclick = () => dialog.close("cancel");
  dialog.showModal();
  return new Promise((resolve) => {
    dialog.onclose = () => resolve(dialog.returnValue === "ok" ? input.value : undefined);
  });
}

/** Why a worker call failed, when it is about a password. */
const errorCode = (e: unknown) => (e instanceof BdfWorkerError || e instanceof ConvertError ? e.code : undefined);

/**
 * Run first; when it needs a password, ask the reader for one and run retry
 * with it, as many times as it takes. undefined when the reader cancels.
 */
async function withPassword<T>(first: () => Promise<T>, retry: (password: string) => Promise<T>, busy: string): Promise<T | undefined> {
  try {
    return await first();
  } catch (e) {
    if (errorCode(e) !== "password-required") throw e;
  }
  let message = "This document is encrypted. Enter its password to open it.";
  let error = false;
  for (;;) {
    setStatus("encrypted document: waiting for the password");
    const password = await askPassword(message, error);
    if (password === undefined) return undefined;
    setStatus(busy);
    try {
      return await retry(password);
    } catch (e) {
      if (errorCode(e) !== "wrong-password") throw e;
      message = "Wrong password. Try again.";
      error = true;
    }
  }
}

/** Bumped by each open: a file opened while another is still loading wins. */
let opening = 0;

/**
 * Open a document and show its first view. An encrypted one stays locked in
 * the worker while the reader is asked for its password. With stream, the
 * document is the outline of a conversion whose pages are still to come.
 */
async function load(source: OpenSource, name?: string, token = ++opening, stream?: Streaming) {
  closeDocument();
  setStatus("loading…");
  const opened = await withPassword(() => client.open(source), (password) => client.unlock(password), "unlocking…");
  if (token !== opening) return;
  if (!opened) {
    setStatus("encrypted document: not opened");
    return;
  }
  manifest = opened;
  // a book opens as facing pages
  if (!layoutChosen) layoutSelect.value = manifest.meta?.source === "epub" ? "spread" : "pages";
  if (stream) {
    stream.view = manifest.views[0].id;
    streaming = stream;
  }
  document.title = `${dcValues(manifest.meta?.dc?.title)[0] ?? name ?? "BDF"} – viewer`;
  $<HTMLInputElement>("q").value = "";
  manifest.views.forEach((v, i) => {
    const b = document.createElement("button");
    b.textContent = v.title || v.id;
    b.title = `${v.title || v.id} (${v.kind}, ${describe(v)})`;
    b.onclick = () => show(v);
    b.id = `tab-${i}`;
    b.dataset.id = v.id;
    b.setAttribute("role", "tab");
    b.setAttribute("aria-controls", stage.id);
    tabs.appendChild(b);
  });
  // "#view=ID" in the address opens that view
  const start = internalLink(location.hash);
  show(manifest.views.find((v) => v.id === start?.view) ?? manifest.views[0]);
}

/**
 * Open a file: a bdf document as it is, anything else converted into one
 * first (asking for the password of an encrypted file).
 */
async function openFile(name: string, data: ArrayBuffer) {
  const token = ++opening;
  closeDocument();
  // a view named in the address belongs to the document shown before
  history.replaceState(null, "", location.pathname + location.search);
  setWarnings([]);
  setDownload();
  setTiming("");
  const kind = sniff(new Uint8Array(data), name);
  if (kind === "bdf") return load({ kind: "buffer", buffer: data }, name, token);
  const busy = `converting ${name}…`;
  setStatus(busy);
  converter ??= new ConverterClient(new Worker("./convert-worker.js"));
  const conv = converter;
  const module = new URL(MODULES[kind], location.href).href;
  // only the Office converters lay text out with the font directory (PDFs embed their fonts, images have no text)
  const fonts = kind === "office" ? new URL(FONTS, location.href).href : undefined;
  let t0 = 0; // of the last attempt: the reader's typing is not part of the conversion
  const open = (password?: string) => {
    t0 = performance.now();
    return conv.open(module, data, { fonts, password, name });
  };
  let res: Opened | undefined;
  try {
    res = await withPassword(() => open(), open, busy);
  } catch (e) {
    if (errorCode(e) !== "unknown-format") throw e;
    throw new Error(`${name} is not in a format this page converts`);
  }
  if (token !== opening) {
    if (res?.stream) conv.close(res.stream).catch(() => {});
    return;
  }
  if (!res) {
    setStatus("encrypted file: not opened");
    return;
  }
  setWarnings(res.warnings);
  const buffer = res.bdf.buffer as ArrayBuffer;
  if (!res.stream) {
    // converted whole
    const took = `${name} converted in ${(performance.now() - t0).toFixed(0)} ms: ${res.summary}`;
    setDownload(res.bdf, name);
    await load({ kind: "buffer", buffer }, name, token);
    if (token === opening) setStatus(took);
    return;
  }
  const st: Streaming = { id: res.stream, token, view: "", state: new Array<PageState>(res.pages).fill(PageState.Pending), left: res.pages, warnings: res.warnings };
  await load({ kind: "buffer", buffer }, name, token, st).finally(() => {
    // the outline did not open, or another file did meanwhile
    if (streaming !== st) conv.close(st.id).catch(() => {});
  });
  if (streaming !== st) return;
  setStatus(`${name}: ${res.pages} ${res.pages === 1 ? "page" : "pages"}, converted as they come into view`);
  showProgress(st);
  convertPages(st, name, t0).catch(openFailed);
}

/**
 * Convert the pages of a stream one by one, each time the one nearest to
 * the visible area, and put them in the document shown; then swap in the
 * finished document, which is also the one to download.
 */
async function convertPages(st: Streaming, name: string, t0: number) {
  const conv = converter!;
  const live = () => streaming === st;
  const arrived = (i: number, state: PageState) => {
    st.state[i] = state;
    st.left--;
    showProgress(st);
    pageArrived(i);
  };
  const failed = (i: number) => (e: unknown) => {
    if (!live()) return;
    setStatus(`page ${i + 1} could not be converted: ${(e as Error).message ?? e}`);
    arrived(i, PageState.Failed);
  };
  // the page being put into the document while the converter goes on with the next
  let adding: Promise<void> = Promise.resolve();
  for (let i = nextPage(st); i >= 0; i = nextPage(st)) {
    st.state[i] = PageState.Converting;
    const t1 = performance.now();
    const page = await conv.page(st.id, i).catch(failed(i));
    await adding;
    if (!live()) return;
    if (!page) continue;
    if (page.warnings.length) setWarnings(st.warnings = [...st.warnings, ...page.warnings]);
    adding = client.addPage(st.view, i, page.bdf.buffer as ArrayBuffer).then(() => {
      if (!live()) return;
      setTiming(`page ${i + 1} converted in ${(performance.now() - t1).toFixed(0)} ms`);
      arrived(i, PageState.Done);
    }, failed(i));
  }
  await adding;
  if (!live()) return;
  streaming = undefined;
  showProgress();
  const lost = st.state.filter((s) => s === PageState.Failed).length;
  if (lost) {
    // the whole document would fail on the same pages: nothing to download
    conv.close(st.id).catch(() => {});
    setStatus(`${name}: ${lost} of ${st.state.length} pages could not be converted`);
    return;
  }
  setTiming("finishing the document…");
  try {
    const res = await conv.finish(st.id);
    if (opening !== st.token) return;
    setDownload(res.bdf, name);
    setWarnings(res.warnings);
    await client.replace({ kind: "buffer", buffer: res.bdf.buffer as ArrayBuffer });
    if (opening !== st.token) return;
    setStatus(`${name} converted in ${(performance.now() - t0).toFixed(0)} ms: ${res.summary}`);
  } catch (e) {
    if (opening === st.token) setStatus(`${name}: the document could not be finished: ${(e as Error).message ?? e}`);
  } finally {
    conv.close(st.id).catch(() => {});
    if (opening === st.token) setTiming("");
  }
}

/**
 * The page to convert next: the first visible one still to come, otherwise
 * the nearest one below the visible pages, then above them.
 */
function nextPage(st: Streaming): number {
  const vis = [...visible].sort((a, b) => a - b);
  for (const i of vis) if (st.state[i] === PageState.Pending) return i;
  const lo = vis[0] ?? 0, hi = vis.at(-1) ?? -1;
  for (let d = 1; d <= st.state.length; d++) {
    if (st.state[hi + d] === PageState.Pending) return hi + d;
    if (st.state[lo - d] === PageState.Pending) return lo - d;
  }
  return st.state.indexOf(PageState.Pending);
}

/** The pages converted so far, next to the status; hidden without a stream. */
function showProgress(st?: Streaming) {
  progress.hidden = !st;
  if (!st) return;
  progress.max = st.state.length;
  progress.value = st.state.length - st.left;
  progress.title = `${progress.value} of ${progress.max} pages converted`;
}

/** A page came in: draw what was waiting for it, or say that it failed. */
function pageArrived(index: number) {
  if (!current || !streaming || current.id !== streaming.view || continuous(current)) return;
  if (book) return book.arrived(index);
  const el = stage.querySelector<HTMLDivElement>(`.page[data-index="${index}"]`);
  if (!el) return;
  el.removeAttribute("aria-busy");
  const gen = generation;
  if (streaming.state[index] === PageState.Failed) {
    const p = document.createElement("p");
    p.className = "pageError";
    p.textContent = "This page could not be converted.";
    el.append(p);
    return;
  }
  if (el.dataset.waitBitmap !== undefined) {
    delete el.dataset.waitBitmap;
    renderPage(current, index, el, gen).catch(unlessStale(gen));
  }
  if (el.dataset.waitText !== undefined) {
    delete el.dataset.waitText;
    renderText(current, index, el, gen).catch(unlessStale(gen));
  }
}

/** Take down the document shown, before another one opens: work started for it is dropped. */
function closeDocument() {
  if (streaming) converter?.close(streaming.id).catch(() => {});
  streaming = undefined;
  showProgress();
  current = undefined;
  generation++;
  book?.destroy();
  book = undefined;
  stage.onscroll = stage.onkeydown = stage.onfocus = null;
  stage.replaceChildren();
  tabs.replaceChildren();
  $<HTMLButtonElement>("prevView").disabled = $<HTMLButtonElement>("nextView").disabled = true;
  $("layoutBox").hidden = $("pageNav").hidden = $("animateBox").hidden = true;
  found.query = ""; found.hits = []; found.rects = []; found.index = -1; found.pages = -1; hitsBox.textContent = "";
}

/** A file that did not open leaves the start page, when no other document is open. */
function openFailed(e: unknown) {
  showError(e);
  if (current || landing.isConnected) return;
  landing.hidden = false;
  stage.replaceChildren(landing);
}

/** Open a file the reader chose or dropped. */
function openLocal(file: File) {
  file.arrayBuffer().then((data) => openFile(file.name, data)).catch(openFailed);
}

let downloadURL: string | undefined;

/** Offer the converted document for download (none when bdf is absent). */
function setDownload(bdf?: Uint8Array, name = "") {
  const a = $<HTMLAnchorElement>("download");
  if (downloadURL) URL.revokeObjectURL(downloadURL);
  downloadURL = undefined;
  a.hidden = !bdf;
  if (!bdf) return;
  // the Blob copies the bytes, which then go to the worker
  downloadURL = URL.createObjectURL(new Blob([bdf as BlobPart], { type: "application/octet-stream" }));
  a.href = downloadURL;
  a.download = `${name.replace(/\.[^.]*$/, "")}.bdf`;
}

/** The warnings of the conversion, in a disclosure next to the status. */
function setWarnings(list: string[]) {
  const box = $<HTMLDetailsElement>("warnings");
  box.hidden = list.length === 0;
  box.open = false;
  box.querySelector("summary")!.textContent = `${list.length} ${list.length === 1 ? "warning" : "warnings"}`;
  box.querySelector("ul")!.replaceChildren(...list.map((w) => {
    const li = document.createElement("li");
    li.textContent = w;
    return li;
  }));
}

/** The start page: a file picker, drag and drop, and the samples published with the page. */
async function showLanding() {
  landing.hidden = false;
  setStatus("no document open");
  const res = await fetch("samples/index.json").catch(() => undefined);
  if (!res?.ok) return;
  const samples = (await res.json()) as { name: string; label: string }[];
  const box = landing.querySelector<HTMLDivElement>("#samples")!;
  box.querySelector(".samples")!.replaceChildren(...samples.map((s) => {
    const b = document.createElement("button");
    b.type = "button";
    b.textContent = s.label;
    b.title = s.name;
    b.onclick = () => {
      setStatus(`fetching ${s.name}…`);
      fetch(`samples/${s.name}`)
        .then((r) => (r.ok ? r.arrayBuffer() : Promise.reject(new Error(`${s.name}: HTTP ${r.status}`))))
        .then((data) => openFile(s.name, data))
        .catch(openFailed);
    };
    return b;
  }));
  box.hidden = false;
}

function init() {
  // tablist keys: arrows, Home and End move to a view and show it
  tabs.onkeydown = (e) => {
    const all = [...tabs.querySelectorAll<HTMLButtonElement>("[role=tab]")];
    const i = all.findIndex((b) => b.dataset.id === current?.id);
    const keys: Record<string, number> = { ArrowRight: i + 1, ArrowLeft: i - 1, Home: 0, End: all.length - 1 };
    const to = keys[e.key];
    if (to === undefined) return;
    e.preventDefault();
    const b = all[(to + all.length) % all.length];
    b.focus();
    show(manifest.views.find((v) => v.id === b.dataset.id)!);
  };
  $("prevView").onclick = () => step(-1);
  $("nextView").onclick = () => step(1);
  const q = $<HTMLInputElement>("q");
  // nothing to search, zoom or lay out before a document is open
  const search = (query: string, step: number) => { if (current) runSearch(query, step).catch(showError); };
  q.onkeydown = (e) => { if (e.key === "Enter") search(q.value, e.shiftKey ? -1 : 1); };
  q.oninput = () => { if (!q.value) search("", 0); };
  $("next").onclick = () => search(q.value, 1);
  $("prev").onclick = () => search(q.value, -1);
  zoomInput.oninput = () => {
    zoom = Number(zoomInput.value);
    const pct = `${Math.round(zoom * 100)}%`;
    $("zoomv").textContent = pct;
    zoomInput.setAttribute("aria-valuetext", pct);
    if (current) show(current);
  };
  // "?layout=spread" opens documents in that layout
  const start = params.get("layout");
  if (start && [...layoutSelect.options].some((o) => o.value === start)) {
    layoutSelect.value = start;
    layoutChosen = true;
  }
  layoutSelect.onchange = () => {
    layoutChosen = true;
    if (current) show(current);
  };
  $("prevPage").onclick = () => book?.turnBy(-1);
  $("nextPage").onclick = () => book?.turnBy(1);
  // pages curl as they turn, unless the reader asks for less motion
  animateBox.checked = !matchMedia("(prefers-reduced-motion: reduce)").matches;
  animateBox.onchange = () => book?.setAnimate(animateBox.checked);

  // files: the picker, and drag and drop anywhere on the page
  const picker = $<HTMLInputElement>("file");
  picker.onchange = () => {
    const file = picker.files?.[0];
    picker.value = "";
    if (file) openLocal(file);
  };
  $("openFile").onclick = () => picker.click();
  landing.querySelector<HTMLButtonElement>("#chooseFile")!.onclick = () => picker.click();
  const carriesFiles = (e: DragEvent) => e.dataTransfer?.types.includes("Files") ?? false;
  let depth = 0; // dragenter and dragleave fire for every element crossed
  const dragging = (on: boolean) => document.body.classList.toggle("dragging", on);
  window.addEventListener("dragenter", (e) => { if (carriesFiles(e)) dragging(++depth > 0); });
  window.addEventListener("dragleave", (e) => { if (carriesFiles(e)) dragging(--depth > 0); });
  window.addEventListener("dragover", (e) => {
    if (!carriesFiles(e)) return;
    e.preventDefault();
    e.dataTransfer!.dropEffect = "copy";
  });
  window.addEventListener("drop", (e) => {
    if (!carriesFiles(e)) return;
    e.preventDefault();
    depth = 0;
    dragging(false);
    const file = e.dataTransfer!.files[0];
    if (file) openLocal(file);
  });
}

function main() {
  init();
  if (!src) return showLanding();
  const source: OpenSource = src.endsWith("/") ? { kind: "split", base: new URL(src, location.href).href } : { kind: "single", url: new URL(src, location.href).href, range: params.has("range") };
  return load(source);
}

function show(v: View) {
  // a layout or zoom change goes on from the page in view
  const keep = current === v ? pageInView() : 0;
  if (current !== v) {
    found.query = ""; found.hits = []; found.rects = []; found.index = -1; found.pages = -1; hitsBox.textContent = "";
    setStatus(describe(v));
  }
  current = v;
  for (const b of tabs.querySelectorAll<HTMLButtonElement>("[role=tab]")) {
    const selected = b.dataset.id === v.id;
    b.setAttribute("aria-selected", String(selected));
    b.tabIndex = selected ? 0 : -1;
    if (selected) {
      stage.setAttribute("aria-labelledby", b.id);
      b.scrollIntoView({ block: "nearest", inline: "nearest" });
    }
  }
  const i = manifest.views.indexOf(v);
  $<HTMLButtonElement>("prevView").disabled = i <= 0;
  $<HTMLButtonElement>("nextView").disabled = i >= manifest.views.length - 1;
  // the address names the view, for a link to it or a reload
  if (manifest.views.length > 1) history.replaceState(null, "", `#view=${encodeURIComponent(v.id)}`);
  $("layoutBox").hidden = v.kind === "sheet" || v.kind === "scroll";
  layoutSelect.querySelector<HTMLOptionElement>("[value=continuous]")!.disabled = v.kind !== "flow";
  generation++;
  book?.destroy();
  book = undefined;
  $("pageNav").hidden = $("animateBox").hidden = true;
  stage.onscroll = stage.onkeydown = stage.onfocus = null;
  stage.replaceChildren();
  stage.scrollTop = 0;
  visibility?.disconnect();
  visible.clear();
  if (v.kind === "sheet") showSheet(v);
  else if (continuous(v)) showContinuous(v);
  else if (inBook(v)) showBook(v, keep);
  else {
    showPages(v);
    if (keep > 0) stage.querySelector(`.page[data-index="${keep}"]`)?.scrollIntoView({ block: "start" });
  }
}

/** The first page in view (in the book layouts, of the spread shown). */
function pageInView(): number {
  if (book) return book.pages[0] ?? 0;
  return visible.size ? Math.min(...visible) : 0;
}

/** Show the view before (delta -1) or after (+1) the current one. */
function step(delta: number) {
  const i = current ? manifest.views.indexOf(current) : -1;
  const v = manifest.views[i + delta];
  if (!v) return;
  show(v);
  focusTab(v);
}

/** Move the focus to the tab of a view. */
function focusTab(v: View) {
  tabs.querySelector<HTMLButtonElement>(`[role=tab][data-id="${CSS.escape(v.id)}"]`)?.focus();
}

/** Whether a view is shown as one continuous scroll: a scroll view always, a flow view on request. */
const continuous = (v: View) => v.kind === "scroll" || (v.kind === "flow" && layoutSelect.value === "continuous");
/** Whether a view of pages is shown one or two pages at a time. */
const inBook = (v: View) => (v.kind === "fixed" || v.kind === "flow") && ["single", "spread", "spread-rtl"].includes(layoutSelect.value);

/** What a view holds, for the status line. */
function describe(v: View): string {
  if (v.kind === "scroll") return "one column without pages";
  if (v.kind !== "sheet") {
    const n = v.pages?.length ?? 0;
    return `${n} ${n === 1 ? "page" : "pages"}`;
  }
  const count = (runs?: [number, number][]) => (runs ?? []).reduce((a, [n]) => a + n, 0);
  return `${count(v.rows)} rows, ${count(v.cols)} columns`;
}

/**
 * Pages: each a labelled group. Bitmaps are rendered when a page comes near
 * the visible area, text layers from further away, so that a screen reader
 * can read on: moving its cursor scrolls the pages, which builds the next
 * layers. The whole document is not loaded up front (split or range sources
 * would fetch every page).
 */
function showPages(v: View) {
  const gen = generation;
  const list = document.createElement("div");
  list.className = "pages";
  const pagesOf = v.pages ?? [];
  const noun = manifest.meta?.source === "pptx" ? "Slide" : "Page";
  // a page still to come is drawn when it comes in (pageArrived)
  const bitmaps = onNear(BITMAP_MARGIN, (el) => {
    const i = Number(el.dataset.index);
    if (pending(v, i)) el.dataset.waitBitmap = "";
    else renderPage(v, i, el, gen).catch(unlessStale(gen));
  });
  const texts = onNear(TEXT_MARGIN, (el) => {
    const i = Number(el.dataset.index);
    if (pending(v, i)) el.dataset.waitText = "";
    else renderText(v, i, el, gen).catch(unlessStale(gen));
  });
  // which pages are visible, so that a stream converts them first
  visibility = new IntersectionObserver((entries) => {
    for (const e of entries) {
      const i = Number((e.target as HTMLElement).dataset.index);
      if (e.isIntersecting) visible.add(i); else visible.delete(i);
    }
  }, { root: stage });
  pagesOf.forEach((p, i) => {
    const el = document.createElement("div");
    el.className = "page";
    el.dataset.index = String(i);
    el.setAttribute("role", "group");
    // a view of one page (a draw.io page) is named by its title
    el.setAttribute("aria-label", pagesOf.length === 1 && v.title ? v.title : `${noun} ${i + 1} of ${pagesOf.length}`);
    if (pending(v, i)) el.setAttribute("aria-busy", "true");
    el.tabIndex = -1; // target of page links
    el.style.width = `${p.w * zoom}px`;
    el.style.height = `${p.h * zoom}px`;
    list.appendChild(el);
    bitmaps.observe(el);
    texts.observe(el);
    visibility!.observe(el);
    if (streaming?.view === v.id && streaming.state[i] === PageState.Failed) pageArrived(i);
  });
  stage.appendChild(list);
}

async function renderPage(v: View, index: number, el: HTMLDivElement, gen: number) {
  const page = v.pages![index];
  const t0 = performance.now();
  const bmp = await client.page(v.id, index, zoom * dpr());
  if (gen !== generation) return bmp.close();
  placeBitmap(el, bmp, page.w, page.h);
  setTiming(`page ${index + 1} rendered in ${(performance.now() - t0).toFixed(0)} ms`);
}

async function renderText(v: View, index: number, el: HTMLDivElement, gen: number) {
  const content = await client.content(v.id, index);
  if (gen !== generation) return;
  el.append(buildTextLayer(content, zoom, layerOptions()), highlightLayer(index));
}

/**
 * Pages one or two at a time, fitted to the stage (times the zoom), with the
 * page buttons in the header. The pages shown are the visible ones, which a
 * stream converts first. A view of a book bound on the right (direction
 * "rtl", spec §4.1) is laid out and turned right to left; the menu can make
 * any view do so.
 */
function showBook(v: View, start: number) {
  const gen = generation;
  const pagesOf = v.pages ?? [];
  const noun = manifest.meta?.source === "pptx" ? "Slide" : "Page";
  const rtl = layoutSelect.value === "spread-rtl" || v.direction === "rtl";
  const b = new Book(stage, {
    pages: pagesOf,
    layout: layoutSelect.value === "single" ? "single" : "spread",
    rtl,
    zoom,
    animate: animateBox.checked,
    start,
    label: (i) => (pagesOf.length === 1 && v.title ? v.title : `${noun} ${i + 1} of ${pagesOf.length}`),
    bitmap: (i, scale) => client.page(v.id, i, scale),
    content: (i) => client.content(v.id, i),
    layers: (i, content, scale) => [buildTextLayer(content, scale, layerOptions()), highlightLayer(i, scale)],
    pending: (i) => pending(v, i),
    failed: (i) => streaming?.view === v.id && streaming.state[i] === PageState.Failed,
    onTurn: (pages) => turned(b, pages, noun),
    onError: unlessStale(gen),
  });
  book = b;
  // the buttons point the way the pages turn
  $("pageNav").classList.toggle("rtl", rtl);
  $("prevPage").textContent = rtl ? "›" : "‹";
  $("nextPage").textContent = rtl ? "‹" : "›";
  $("pageNav").hidden = $("animateBox").hidden = false;
  turned(b, b.pages, noun);
}

/** The pages of a book shown: the page number next to its buttons, and what a stream converts first. */
function turned(b: Book, pages: number[], noun: string) {
  visible.clear();
  for (const i of pages) visible.add(i);
  const n = current?.pages?.length ?? 0;
  const range = pages.length > 1 ? `${pages[0] + 1}–${pages[pages.length - 1] + 1}` : `${(pages[0] ?? 0) + 1}`;
  const shown = document.createElement("span");
  shown.setAttribute("aria-hidden", "true");
  shown.textContent = `${range} / ${n}`;
  const spoken = document.createElement("span");
  spoken.className = "sr-only";
  spoken.textContent = `${noun}${pages.length > 1 ? "s" : ""} ${range} of ${n}`;
  $("pageNum").replaceChildren(shown, spoken);
  $<HTMLButtonElement>("prevPage").disabled = !b.canTurn(-1);
  $<HTMLButtonElement>("nextPage").disabled = !b.canTurn(1);
}

/** Follow a link to a page (0-based): scroll to it and focus it. */
function goToPage(index: number) {
  if (!current) return;
  if (book) return book.go(index, true, true);
  if (continuous(current)) {
    const { top, band } = continuousPosition(current, index);
    stage.scrollTo({ top: top * zoom });
    stage.querySelector<HTMLElement>(`.page[data-y="${band}"]`)?.focus({ preventScroll: true });
    return;
  }
  const el = stage.querySelector<HTMLElement>(`.page[data-index="${index}"]`);
  el?.scrollIntoView({ block: "start" });
  el?.focus({ preventScroll: true });
}

/** Search: ask the worker for hits, locate them once, then highlight and step through them. */
async function runSearch(query: string, step: number) {
  const v = current;
  if (!v) return;
  // while pages are coming in, a search covers those converted so far: search again once there are more
  const converted = streaming?.view === v.id ? streaming.state.length - streaming.left : -1;
  if (query !== found.query || converted !== found.pages) {
    // the same query over more pages goes on from the hit it was at
    const at = query === found.query ? found.hits[found.index]?.segments[0] : undefined;
    found.query = query;
    found.pages = converted;
    const hits = query ? await client.search(v.id, query, { limit: 500, context: 30 }) : [];
    const rects = hits.length ? await client.locate(v.id, hits) : [];
    if (v !== current) return;
    found.hits = hits;
    found.rects = rects;
    found.index = -1;
    if (at) {
      const same = hits.findIndex(({ segments: [s] }) => s.a === at.a && s.b === at.b && s.ordinal === at.ordinal && s.start === at.start);
      if (same >= 0) found.index = same;
    }
  }
  if (found.hits.length) found.index = (found.index + step + found.hits.length) % found.hits.length;
  showHitCount(query);
  // refresh highlights on everything that is rendered
  for (const el of stage.querySelectorAll<HTMLDivElement>(".page[data-index]")) {
    const old = el.querySelector(".hlLayer");
    if (old) old.replaceWith(highlightLayer(Number(el.dataset.index), book?.scale ?? zoom));
  }
  sheetViews.get(v)?.redraw();
  for (const el of stage.querySelectorAll<HTMLDivElement>(".page[data-y]")) {
    const old = el.querySelector(".hlLayer");
    if (old) old.replaceWith(bandHighlightLayer(v, Number(el.dataset.y), el));
  }
  // scroll to the current hit
  const rect = found.rects[found.index]?.[0];
  if (!rect) return;
  if (v.kind === "sheet") {
    sheetViews.get(v)?.reveal(rect);
  } else if (continuous(v)) {
    const y = continuousRect(v, rect).y;
    stage.scrollTo({ top: Math.max(0, y * zoom - stage.clientHeight / 2), behavior: "smooth" });
  } else if (book) {
    book.go(rect.a, true);
  } else {
    const el = stage.querySelector<HTMLDivElement>(`.page[data-index="${rect.a}"]`);
    el?.scrollIntoView({ block: "center", behavior: "smooth" });
  }
}

/**
 * "3 of 12" next to the search box. It is a live region, and the highlights
 * are only drawn, so screen readers also get the current hit's page and
 * its surrounding text.
 */
function showHitCount(query: string) {
  if (!query) { hitsBox.textContent = ""; return; }
  const partial = found.pages >= 0 ? ` (in ${found.pages} of ${current?.pages?.length ?? 0} pages converted)` : "";
  if (!found.hits.length) { hitsBox.textContent = `no matches${partial}`; return; }
  const hit = found.hits[found.index];
  // the strips of a scroll view are not pages
  const page = current?.kind === "sheet" || current?.kind === "scroll" ? "" : `, page ${hit.segments[0].a + 1}`;
  const detail = document.createElement("span");
  detail.className = "sr-only";
  detail.textContent = `${page}: ${hit.context.replace(/ ⏎ /g, " ")}`;
  hitsBox.replaceChildren(`${found.index + 1} of ${found.hits.length}${partial}`, detail);
}

/** Highlight rectangles of the hits on one page; scale is CSS px per unit. */
function highlightLayer(pageIndex: number, scale = zoom): HTMLDivElement {
  const layer = document.createElement("div");
  layer.className = "textLayer hlLayer";
  found.rects.forEach((rects, hi) => {
    for (const r of rects) {
      if (r.a !== pageIndex) continue;
      const d = document.createElement("div");
      d.className = hi === found.index ? "hl current" : "hl";
      d.style.left = `${r.x * scale}px`;
      d.style.top = `${r.y * scale}px`;
      d.style.width = `${r.w * scale}px`;
      d.style.height = `${r.h * scale}px`;
      layer.appendChild(d);
    }
  });
  return layer;
}

const BAND = 800;

/** Continuous layout: stacked body rectangles (the strips of a scroll view, without gaps). */
function continuousSize(v: View) {
  const gap = v.kind === "scroll" ? 0 : v.continuous?.gap ?? 0;
  let height = 0, width = 0;
  const tops: number[] = [];
  for (const p of v.pages ?? []) {
    const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
    tops.push(height);
    height += b.h + gap;
    width = Math.max(width, b.w);
  }
  return { width, height: Math.max(0, height - gap), tops };
}

/** A rectangle of a page (a hit) in continuous coordinates. */
function continuousRect(v: View, r: { a: number; x: number; y: number; w: number; h: number }) {
  const p = v.pages![r.a];
  const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
  return { x: r.x - b.x, y: continuousSize(v).tops[r.a] + r.y - b.y, w: r.w, h: r.h };
}

/** Highlight rectangles of the hits in one band of the continuous layout. */
function bandHighlightLayer(v: View, y0: number, el: HTMLElement): HTMLDivElement {
  const layer = document.createElement("div");
  layer.className = "textLayer hlLayer";
  const h = el.offsetHeight / zoom;
  found.rects.forEach((rects, hi) => {
    for (const r of rects) {
      const c = continuousRect(v, r);
      if (c.y + c.h < y0 || c.y > y0 + h) continue;
      const d = document.createElement("div");
      d.className = hi === found.index ? "hl current" : "hl";
      d.style.left = `${c.x * zoom}px`;
      d.style.top = `${(c.y - y0) * zoom}px`;
      d.style.width = `${c.w * zoom}px`;
      d.style.height = `${c.h * zoom}px`;
      layer.appendChild(d);
    }
  });
  return layer;
}

/** Where a page starts in the continuous layout, and the band that holds it. */
function continuousPosition(v: View, index: number) {
  const top = continuousSize(v).tops[index] ?? 0;
  return { top, band: Math.floor(top / BAND) * BAND };
}

/** Continuous flow: body rectangles stacked, rendered in 800-unit bands on demand (text layers further ahead). */
function showContinuous(v: View) {
  const gen = generation;
  const { width, height } = continuousSize(v);
  const list = document.createElement("div");
  list.className = "pages";
  const viewportOf = (el: HTMLElement) => {
    const y = Number(el.dataset.y);
    return { x: 0, y, w: width, h: Math.min(BAND, height - y) };
  };
  const bitmaps = onNear(BITMAP_MARGIN, (el) => {
    const vp = viewportOf(el);
    client.continuous(v.id, vp, zoom * dpr()).then((bmp) => (gen === generation ? placeBitmap(el, bmp, vp.w, vp.h) : bmp.close())).catch(unlessStale(gen));
  });
  const texts = onNear(TEXT_MARGIN, (el) => {
    client.continuousContent(v.id, viewportOf(el)).then((c) => {
      if (gen === generation) el.append(buildTextLayer(c, zoom, layerOptions()), bandHighlightLayer(v, Number(el.dataset.y), el));
    }).catch(unlessStale(gen));
  });
  for (let y = 0; y < height; y += BAND) {
    const el = document.createElement("div");
    el.className = "page";
    el.dataset.y = String(y);
    el.tabIndex = -1; // target of page links
    el.style.width = `${width * zoom}px`;
    el.style.height = `${Math.min(BAND, height - y) * zoom}px`;
    el.style.boxShadow = "none";
    list.appendChild(el);
    bitmaps.observe(el);
    texts.observe(el);
  }
  list.style.gap = "0";
  stage.appendChild(list);
}

/** Row or column sizes of a sheet, from the manifest's run-length list. */
class SheetAxis {
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
  /** The size of entry i (0: hidden). */
  size(i: number): number {
    for (const [n, s] of this.runs) {
      if (i < n) return s;
      i -= n;
    }
    return 0;
  }
  /** The entry with a size at position p (the first or the last one outside). */
  at(p: number): number {
    let i = 0, start = 0, last = 0;
    for (const [n, s] of this.runs) {
      if (s > 0) {
        if (p < start + n * s) return i + Math.max(0, Math.floor((p - start) / s));
        last = i + n - 1;
      }
      start += n * s;
      i += n;
    }
    return last;
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

const SHEET_HEADER = 20;

/** A cell's name: B12. */
const cellName = (r: number, c: number) => `${columnLabel(c)}${r + 1}`;

/**
 * Sheet: a scroll area the size of the sheet with one sticky canvas that
 * shows the column and row headers and the visible region, split into the
 * frozen panes of the view (each rendered by the worker as a region of its
 * own). A text layer (a table of the cells, for screen readers) covers the
 * visible region and a screen around it; the cells of the frozen panes are
 * pinned with a translation that follows the scroll.
 *
 * Cells are selected as in a spreadsheet, not as text: a press on a cell
 * and a drag (scrolling at the edges), Shift with a press to extend, the
 * headers for whole columns and rows (the corner for all), the arrow keys
 * (with Shift to extend), Ctrl or Cmd+A, and Escape. A merged cell is
 * selected whole. Copy puts the cells on the clipboard as tab-separated
 * values and an HTML table, which spreadsheets paste as cells; hidden rows
 * and columns are left out, and a selection of whole rows or columns ends
 * with the last cell with text.
 */
function showSheet(v: View) {
  const cols = new SheetAxis(v.cols), rows = new SheetAxis(v.rows);
  const fc = Math.min(v.freeze?.cols ?? 0, cols.count), fr = Math.min(v.freeze?.rows ?? 0, rows.count);
  const fw = cols.pos(fc), fh = rows.pos(fr); // frozen extent in units
  const hw = Math.max(32, String(rows.count).length * 7 + 14), hh = SHEET_HEADER; // headers in CSS px
  const wrap = document.createElement("div");
  wrap.className = "sheet";
  wrap.style.width = `${hw + cols.total * zoom}px`;
  wrap.style.height = `${hh + rows.total * zoom}px`;
  const canvas = document.createElement("canvas");
  canvas.setAttribute("aria-hidden", "true");
  const layerHost = document.createElement("div");
  layerHost.className = "sheetText";
  layerHost.style.cssText = `left: ${hw}px; top: ${hh}px; width: ${cols.total * zoom}px; height: ${rows.total * zoom}px`;
  // The page's selection is on this element while cells are selected: copy
  // then fires in every browser (Safari fires it only with a selection) and
  // the browser's Copy is offered. It says which cells.
  const hold = document.createElement("div");
  hold.className = "cellHold";
  hold.setAttribute("aria-hidden", "true");
  wrap.append(canvas, layerHost);
  stage.append(wrap, hold);

  const gen = generation;
  const sheet = { rows: rows.count, cols: cols.count, headerRows: fr, headerCols: fc };
  // cells of the frozen panes stay where they are while the sheet scrolls
  const pin = (span: HTMLSpanElement, r: { x: number; y: number }) => {
    const x = r.x < fw, y = r.y < fh;
    if (!x && !y) return;
    span.style.translate = `${x ? "var(--sx)" : "0px"} ${y ? "var(--sy)" : "0px"}`;
    span.style.zIndex = "2"; // above the other runs (TEXT_LAYER_CSS)
  };
  /** The text layer's content, the rectangles it covers, and its cells by position. */
  let loaded: { content: TextContent; rects: { x: number; y: number; w: number; h: number }[]; cells: Map<string, CellText> } | undefined;
  /** Merged cells seen so far, by position. */
  const merges = new Map<string, CellText>();
  let covered: { x: number; y: number; w: number; h: number } | undefined;
  let textTimer: ReturnType<typeof setTimeout> | undefined;
  const updateText = (vp: { x: number; y: number; w: number; h: number }) => {
    const inside = covered && vp.x >= covered.x && vp.y >= covered.y && vp.x + vp.w <= covered.x + covered.w && vp.y + vp.h <= covered.y + covered.h;
    if (inside) return;
    clearTimeout(textTimer);
    textTimer = setTimeout(() => {
      if (gen !== generation) return;
      const region = { x: Math.max(0, vp.x - vp.w), y: Math.max(0, vp.y - vp.h), w: vp.w * 3, h: vp.h * 3 };
      // the frozen panes along the region
      const frozen = [{ x: region.x, y: 0, w: region.w, h: fh }, { x: 0, y: region.y, w: fw, h: region.h }, { x: 0, y: 0, w: fw, h: fh }];
      client.sheetContent(v.id, [region, ...frozen]).then((content: TextContent) => {
        if (gen !== generation) return;
        covered = region;
        const cells = new Map<string, CellText>();
        for (const c of tableCells(content)) {
          cells.set(`${c.row},${c.col}`, c);
          if (c.rows > 1 || c.cols > 1) merges.set(`${c.row},${c.col}`, c);
        }
        loaded = { content, rects: [region, ...frozen], cells };
        const layer = buildTextLayer(content, zoom, { ...layerOptions(), sheet, onSpan: pin, selectable: false });
        layerHost.replaceChildren(layer);
      }).catch(unlessStale(gen));
    }, 150);
  };

  // --- the cells selected ---
  const sel = () => cellSelections.get(v);
  const lastRow = rows.count - 1, lastCol = cols.count - 1;
  /** The merged cell at a position, or the cell itself. */
  const cellAt = (r: number, c: number): CellRange => {
    for (const m of merges.values()) {
      if (r >= m.row && r < m.row + m.rows && c >= m.col && c < m.col + m.cols) return { row0: m.row, col0: m.col, row1: m.row + m.rows - 1, col1: m.col + m.cols - 1 };
    }
    return { row0: r, col0: c, row1: r, col1: c };
  };
  /** The rectangle a selection covers, grown to take in the merged cells it cuts. */
  const rangeOf = (s: CellSelection): CellRange => {
    const g = {
      row0: Math.min(s.anchor.r, s.focus.r), col0: Math.min(s.anchor.c, s.focus.c),
      row1: Math.max(s.anchor.r, s.focus.r), col1: Math.max(s.anchor.c, s.focus.c),
    };
    for (let grew = true; grew;) {
      grew = false;
      for (const m of merges.values()) {
        if (m.row > g.row1 || m.row + m.rows - 1 < g.row0 || m.col > g.col1 || m.col + m.cols - 1 < g.col0) continue;
        const r0 = Math.min(g.row0, m.row), c0 = Math.min(g.col0, m.col), r1 = Math.max(g.row1, m.row + m.rows - 1), c1 = Math.max(g.col1, m.col + m.cols - 1);
        if (r0 !== g.row0 || c0 !== g.col0 || r1 !== g.row1 || c1 !== g.col1) {
          Object.assign(g, { row0: r0, col0: c0, row1: r1, col1: c1 });
          grew = true;
        }
      }
    }
    return g;
  };
  /** Where a rectangle of cells lies, in sheet units. */
  const boxOf = (g: CellRange) => {
    const x = cols.pos(g.col0), y = rows.pos(g.row0);
    return { x, y, w: cols.pos(g.col1 + 1) - x, h: rows.pos(g.row1 + 1) - y };
  };
  const allRows = (g: CellRange) => g.row0 === 0 && g.row1 === lastRow;
  const allCols = (g: CellRange) => g.col0 === 0 && g.col1 === lastCol;
  /** A1 notation: B2:D5, whole columns B:D, whole rows 2:5. */
  const rangeName = (g: CellRange) => {
    if (allRows(g) && allCols(g)) return "all cells";
    if (allRows(g)) return `${columnLabel(g.col0)}:${columnLabel(g.col1)}`;
    if (allCols(g)) return `${g.row0 + 1}:${g.row1 + 1}`;
    const a = cellName(g.row0, g.col0), b = cellName(g.row1, g.col1);
    return a === b ? a : `${a}:${b}`;
  };
  /** What the status says of a selection: a cell's name and text, or the rectangle and its size. */
  const describeSelection = (s: CellSelection) => {
    const g = rangeOf(s), a = cellAt(s.anchor.r, s.anchor.c);
    if (g.row0 === a.row0 && g.col0 === a.col0 && g.row1 === a.row1 && g.col1 === a.col1) {
      const text = loaded?.cells.get(`${a.row0},${a.col0}`)?.text;
      return text ? `${cellName(a.row0, a.col0)}: ${text}` : cellName(a.row0, a.col0);
    }
    return `${rangeName(g)} selected, ${g.row1 - g.row0 + 1} × ${g.col1 - g.col0 + 1} cells`;
  };

  /** The selected cells for the clipboard: made when a selection settles, from the text layer's content when it has them. */
  let clip: { key: string; data?: CellClipboard; promise: Promise<CellClipboard> } | undefined;
  const clipboardFor = (g: CellRange) => {
    const key = `${g.row0},${g.col0},${g.row1},${g.col1}`;
    if (clip?.key === key) return clip;
    const box = boxOf(g);
    const opts = { trim: allRows(g) || allCols(g), skipRow: (r: number) => rows.size(r) === 0, skipCol: (c: number) => cols.size(c) === 0 };
    const inside = (r: { x: number; y: number; w: number; h: number }) => box.x >= r.x && box.y >= r.y && box.x + box.w <= r.x + r.w && box.y + box.h <= r.y + r.h;
    const entry: { key: string; data?: CellClipboard; promise: Promise<CellClipboard> } = { key, promise: Promise.resolve({ text: "", html: "" }) };
    if (loaded?.rects.some(inside)) {
      const data = cellClipboard(tableCells(loaded.content), g, opts);
      entry.data = data;
      entry.promise = Promise.resolve(data);
    } else {
      entry.promise = client.sheetContent(v.id, box).then((c) => (entry.data = cellClipboard(tableCells(c), g, opts)));
      entry.promise.catch(() => { if (clip === entry) clip = undefined; }); // tried again on copy, which reports the error
    }
    clip = entry;
    return entry;
  };

  /** Put the page's selection on the holding element (the name of the cells in it). */
  const holdSelection = () => {
    const s = sel();
    if (!s) return;
    hold.textContent = rangeName(rangeOf(s));
    getSelection()?.selectAllChildren(hold);
  };
  const select = (s: CellSelection) => {
    cellSelections.set(v, s);
    if (bitmaps.length) paint();
  };
  /** A selection made: hold the page's selection, get the cells ready to copy and say what is selected. */
  const settle = () => {
    const s = sel();
    if (!s) return;
    holdSelection();
    clipboardFor(rangeOf(s));
    setStatus(describeSelection(s));
  };
  const clear = () => {
    cellSelections.delete(v);
    const ds = getSelection();
    if (ds && hold.contains(ds.anchorNode)) ds.removeAllRanges();
    if (bitmaps.length) paint();
    setStatus(describe(v));
  };

  /** Where a point of the window lies on the sheet: its zone, and the cell under it (clamped into the cells when clamp). */
  type Zone = "cells" | "cols" | "rows" | "all";
  const hit = (x: number, y: number, clamp = false) => {
    const b = stage.getBoundingClientRect();
    let vx = x - b.left - stage.clientLeft, vy = y - b.top - stage.clientTop;
    const zone: Zone = vy < hh ? (vx < hw ? "all" : "cols") : vx < hw ? "rows" : "cells";
    if (clamp) {
      vx = Math.min(Math.max(vx, hw), stage.clientWidth - 1);
      vy = Math.min(Math.max(vy, hh), stage.clientHeight - 1);
    }
    // the frozen panes do not scroll
    const ux = vx - hw < fw * zoom ? (vx - hw) / zoom : (vx - hw + stage.scrollLeft) / zoom;
    const uy = vy - hh < fh * zoom ? (vy - hh) / zoom : (vy - hh + stage.scrollTop) / zoom;
    return { zone, r: rows.at(uy), c: cols.at(ux) };
  };
  /** A press: a cell, whole columns or rows, or all; with extend (Shift) from the active cell. */
  const press = (zone: Zone, at: { r: number; c: number }, extend: boolean) => {
    const s = extend ? sel() : undefined;
    if (zone === "all") select({ anchor: { r: 0, c: 0 }, focus: { r: lastRow, c: lastCol } });
    else if (zone === "cols") select({ anchor: { r: 0, c: s?.anchor.c ?? at.c }, focus: { r: lastRow, c: at.c } });
    else if (zone === "rows") select({ anchor: { r: s?.anchor.r ?? at.r, c: 0 }, focus: { r: at.r, c: lastCol } });
    else select({ anchor: s?.anchor ?? at, focus: at });
  };

  let drag: { zone: Zone; id: number; x: number; y: number } | undefined;
  let tap: { id: number; x: number; y: number } | undefined;
  let scrollFrame = 0;
  /** Move the end of the selection being dragged to the cell under the pointer. */
  const dragTo = () => {
    const s = sel();
    if (!drag || drag.zone === "all" || !s) return;
    const h = hit(drag.x, drag.y, true);
    const focus = drag.zone === "cols" ? { r: lastRow, c: h.c } : drag.zone === "rows" ? { r: h.r, c: lastCol } : { r: h.r, c: h.c };
    if (focus.r !== s.focus.r || focus.c !== s.focus.c) select({ anchor: s.anchor, focus });
  };
  /** While the pointer is past an edge of the cells, scroll that way (faster further out). */
  const autoscroll = () => {
    if (scrollFrame) return;
    const tick = () => {
      scrollFrame = 0;
      if (!drag) return;
      const b = stage.getBoundingClientRect();
      const vx = drag.x - b.left - stage.clientLeft, vy = drag.y - b.top - stage.clientTop;
      const past = (p: number, lo: number, hi: number) => Math.max(-40, Math.min(40, p < lo ? p - lo : p > hi ? p - hi : 0));
      const dx = drag.zone === "rows" || drag.zone === "all" ? 0 : past(vx, hw, stage.clientWidth);
      const dy = drag.zone === "cols" || drag.zone === "all" ? 0 : past(vy, hh, stage.clientHeight);
      if (!dx && !dy) return;
      stage.scrollBy(dx, dy);
      dragTo();
      scrollFrame = requestAnimationFrame(tick);
    };
    scrollFrame = requestAnimationFrame(tick);
  };
  wrap.addEventListener("pointerdown", (e) => {
    if (e.button !== 0) return;
    const h = hit(e.clientX, e.clientY);
    // a link in a cell is followed
    if (h.zone === "cells" && (e.target as Element).closest("a")) return;
    // a finger scrolls; a tap selects
    if (e.pointerType === "touch") {
      tap = { id: e.pointerId, x: e.clientX, y: e.clientY };
      return;
    }
    // the press also focuses the stage, for the keys (the sheet's text cannot be selected: index.html)
    press(h.zone, h, e.shiftKey);
    drag = { zone: h.zone, id: e.pointerId, x: e.clientX, y: e.clientY };
    wrap.setPointerCapture(e.pointerId);
  });
  // a link of a cell scrolled under the headers is not followed from them
  const underHeader = (e: MouseEvent) => hit(e.clientX, e.clientY).zone !== "cells" && !!(e.target as Element).closest("a");
  wrap.addEventListener("mousedown", (e) => { if (underHeader(e)) e.preventDefault(); });
  wrap.addEventListener("click", (e) => { if (underHeader(e)) e.preventDefault(); }, true);
  wrap.addEventListener("pointermove", (e) => {
    if (!drag || e.pointerId !== drag.id) return;
    drag.x = e.clientX;
    drag.y = e.clientY;
    dragTo();
    autoscroll();
  });
  const release = (e: PointerEvent) => {
    if (tap?.id === e.pointerId) {
      const t = tap;
      tap = undefined;
      if (e.type !== "pointerup" || Math.hypot(e.clientX - t.x, e.clientY - t.y) > 10) return;
      const h = hit(e.clientX, e.clientY);
      if (h.zone === "cells" && (e.target as Element).closest("a")) return;
      press(h.zone, h, false);
      settle();
      return;
    }
    if (!drag || e.pointerId !== drag.id) return;
    drag = undefined;
    cancelAnimationFrame(scrollFrame);
    scrollFrame = 0;
    settle();
  };
  wrap.addEventListener("pointerup", release);
  wrap.addEventListener("pointercancel", release);

  /** The next cell with a size from at along d (past a merged cell), or at itself at the edge. */
  const step = (at: { r: number; c: number }, [dr, dc]: [number, number], leaveMerge: boolean) => {
    let { r, c } = at;
    if (leaveMerge) {
      const m = cellAt(r, c);
      r = dr > 0 ? m.row1 : dr < 0 ? m.row0 : r;
      c = dc > 0 ? m.col1 : dc < 0 ? m.col0 : c;
    }
    do {
      r += dr;
      c += dc;
    } while (r >= 0 && r <= lastRow && c >= 0 && c <= lastCol && (dr ? rows.size(r) : cols.size(c)) === 0);
    return r >= 0 && r <= lastRow && c >= 0 && c <= lastCol ? { r, c } : at;
  };
  /** Scroll so that a cell is in view (cells of a frozen pane always are). */
  const bring = (at: { r: number; c: number }) => {
    const g = boxOf(cellAt(at.r, at.c));
    const pw = (stage.clientWidth - hw) / zoom - fw, ph = (stage.clientHeight - hh) / zoom - fh; // the scrolled pane, in units
    let left = stage.scrollLeft, top = stage.scrollTop;
    if (at.c >= fc) {
      const x = fw + left / zoom;
      if (g.x < x || g.w > pw) left = (g.x - fw) * zoom;
      else if (g.x + g.w > x + pw) left = (g.x + g.w - fw - pw) * zoom;
    }
    if (at.r >= fr) {
      const y = fh + top / zoom;
      if (g.y < y || g.h > ph) top = (g.y - fh) * zoom;
      else if (g.y + g.h > y + ph) top = (g.y + g.h - fh - ph) * zoom;
    }
    if (left !== stage.scrollLeft || top !== stage.scrollTop) stage.scrollTo({ left, top });
  };
  const KEYS: Record<string, [number, number]> = { ArrowUp: [-1, 0], ArrowDown: [1, 0], ArrowLeft: [0, -1], ArrowRight: [0, 1] };
  stage.onkeydown = (e) => {
    if (e.target !== stage || current !== v) return;
    const mod = e.metaKey || e.ctrlKey;
    if (mod && !e.shiftKey && !e.altKey && e.key.toLowerCase() === "a") {
      e.preventDefault();
      press("all", { r: 0, c: 0 }, false);
      return settle();
    }
    if (e.key === "Escape" && sel()) {
      e.preventDefault();
      return clear();
    }
    const d = KEYS[e.key];
    if (!d || mod || e.altKey) return;
    e.preventDefault();
    const s = sel();
    if (!s) {
      // the first key selects the top left cell in view
      const at = { r: fr > 0 ? 0 : rows.at(stage.scrollTop / zoom), c: fc > 0 ? 0 : cols.at(stage.scrollLeft / zoom) };
      select({ anchor: at, focus: at });
    } else if (e.shiftKey) {
      select({ anchor: s.anchor, focus: step(s.focus, d, false) });
    } else {
      const at = step(s.anchor, d, true);
      select({ anchor: at, focus: at });
    }
    bring(sel()!.focus);
    settle();
  };
  // back to the sheet with the keys: its cells are what copy takes again
  stage.onfocus = () => {
    const ds = getSelection();
    if (current === v && sel() && ds && (ds.isCollapsed || !wrap.contains(ds.anchorNode))) holdSelection();
  };

  type Pane = { src: { x: number; y: number; w: number; h: number }; dx: number; dy: number };
  let panes: Pane[] = [];
  let bitmaps: ImageBitmap[] = [];
  const paint = () => {
    const d = dpr();
    const vw = Math.min(stage.clientWidth, hw + cols.total * zoom), vh = Math.min(stage.clientHeight, hh + rows.total * zoom);
    canvas.width = Math.ceil(vw * d);
    canvas.height = Math.ceil(vh * d);
    canvas.style.width = `${vw}px`;
    canvas.style.height = `${vh}px`;
    const ctx = canvas.getContext("2d")!;
    ctx.fillStyle = "#fff";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    panes.forEach((p, i) => ctx.drawImage(bitmaps[i], p.dx * d, p.dy * d, p.src.w * zoom * d, p.src.h * zoom * d));
    // search hits and the selected cells, clipped to each pane
    ctx.save();
    ctx.scale(d, d);
    const s = sel();
    const g = s && rangeOf(s);
    for (const p of panes) {
      ctx.save();
      ctx.beginPath();
      ctx.rect(p.dx, p.dy, p.src.w * zoom, p.src.h * zoom);
      ctx.clip();
      found.rects.forEach((rects, hi) => {
        ctx.fillStyle = hi === found.index ? "rgba(255,120,0,.5)" : "rgba(255,210,0,.45)";
        for (const r of rects) ctx.fillRect(p.dx + (r.x - p.src.x) * zoom, p.dy + (r.y - p.src.y) * zoom, r.w * zoom, r.h * zoom);
      });
      if (s && g) drawSelection(ctx, p, g, cellAt(s.anchor.r, s.anchor.c));
      ctx.restore();
    }
    drawHeaders(ctx, vw, vh, g);
    ctx.restore();
  };
  /** The selected cells in a pane: shaded but for the active cell, and outlined. */
  const drawSelection = (ctx: CanvasRenderingContext2D, p: Pane, g: CellRange, active: CellRange) => {
    // in CSS px of the canvas, cut to the pane (a whole column is millions of pixels long)
    const place = (r: CellRange) => {
      const b = boxOf(r);
      const x0 = Math.max(p.dx - 4, p.dx + (b.x - p.src.x) * zoom), y0 = Math.max(p.dy - 4, p.dy + (b.y - p.src.y) * zoom);
      const x1 = Math.min(p.dx + p.src.w * zoom + 4, p.dx + (b.x + b.w - p.src.x) * zoom), y1 = Math.min(p.dy + p.src.h * zoom + 4, p.dy + (b.y + b.h - p.src.y) * zoom);
      return { x: x0, y: y0, w: x1 - x0, h: y1 - y0 };
    };
    const b = place(g), a = place(active);
    if (b.w <= 0 || b.h <= 0) return;
    ctx.beginPath();
    ctx.rect(b.x, b.y, b.w, b.h);
    if (a.w > 0 && a.h > 0) ctx.rect(a.x, a.y, a.w, a.h);
    ctx.fillStyle = "rgba(26, 115, 232, .12)";
    ctx.fill("evenodd");
    ctx.strokeStyle = "#1a73e8";
    ctx.lineWidth = 2;
    ctx.strokeRect(b.x, b.y, b.w, b.h);
  };
  const drawHeaders = (ctx: CanvasRenderingContext2D, vw: number, vh: number, g: CellRange | undefined) => {
    const sx = stage.scrollLeft, sy = stage.scrollTop;
    ctx.fillStyle = "#f3f4f6";
    ctx.fillRect(0, 0, vw, hh);
    ctx.fillRect(0, 0, hw, vh);
    ctx.font = "11px system-ui, sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.strokeStyle = "#c8ccd2";
    ctx.lineWidth = 1;
    // the headers of the selected cells are shaded, darker when whole columns or rows are selected
    const header = (x: number, y: number, w: number, h: number, label: string, on: 0 | 1 | 2) => {
      if (on) {
        ctx.fillStyle = on === 2 ? "#c2d6f6" : "#dfe8f7";
        ctx.fillRect(x, y, w, h);
      }
      ctx.fillStyle = "#444";
      ctx.fillText(label, x + w / 2, y + h / 2);
      ctx.strokeRect(Math.round(x) + 0.5, Math.round(y) + 0.5, Math.round(w), Math.round(h));
    };
    const colOn = (c: number) => (!g || c < g.col0 || c > g.col1 ? 0 : allRows(g) ? 2 : 1);
    const rowOn = (r: number) => (!g || r < g.row0 || r > g.row1 ? 0 : allCols(g) ? 2 : 1);
    // frozen columns, then the scrolled ones clipped to their pane
    const colsIn = (from: number, to: number, shift: number, clipX: number) => {
      ctx.save();
      ctx.beginPath();
      ctx.rect(clipX, 0, vw - clipX, hh);
      ctx.clip();
      cols.each(from, to, (c, start, size) => header(hw + start * zoom - shift, 0, size * zoom, hh, columnLabel(c), colOn(c)));
      ctx.restore();
    };
    colsIn(0, fw, 0, hw);
    colsIn(fw + sx / zoom, fw + sx / zoom + (vw - hw) / zoom, sx, hw + fw * zoom);
    const rowsIn = (from: number, to: number, shift: number, clipY: number) => {
      ctx.save();
      ctx.beginPath();
      ctx.rect(0, clipY, hw, vh - clipY);
      ctx.clip();
      rows.each(from, to, (r, start, size) => header(0, hh + start * zoom - shift, hw, size * zoom, String(r + 1), rowOn(r)));
      ctx.restore();
    };
    rowsIn(0, fh, 0, hh);
    rowsIn(fh + sy / zoom, fh + sy / zoom + (vh - hh) / zoom, sy, hh + fh * zoom);
    ctx.fillStyle = g && allRows(g) && allCols(g) ? "#c2d6f6" : "#e5e7eb";
    ctx.fillRect(0, 0, hw, hh);
    // the edges of the frozen panes
    ctx.strokeStyle = "#8a9099";
    ctx.beginPath();
    if (fc > 0) { ctx.moveTo(hw + fw * zoom + 0.5, 0); ctx.lineTo(hw + fw * zoom + 0.5, vh); }
    if (fr > 0) { ctx.moveTo(0, hh + fh * zoom + 0.5); ctx.lineTo(vw, hh + fh * zoom + 0.5); }
    ctx.stroke();
  };

  let pending = false;
  const draw = async () => {
    if (pending) return;
    pending = true;
    await new Promise((r) => requestAnimationFrame(r));
    pending = false;
    const vw = stage.clientWidth, vh = stage.clientHeight;
    const sx = stage.scrollLeft, sy = stage.scrollTop;
    layerHost.style.setProperty("--sx", `${sx}px`);
    layerHost.style.setProperty("--sy", `${sy}px`);
    // the scrolled region starts where the frozen panes end
    const mx = fw + sx / zoom, my = fh + sy / zoom;
    const pw = Math.max(0, Math.min((vw - hw) / zoom - fw, cols.total - mx)), ph = Math.max(0, Math.min((vh - hh) / zoom - fh, rows.total - my));
    const next: Pane[] = [
      { src: { x: mx, y: my, w: pw, h: ph }, dx: hw + fw * zoom, dy: hh + fh * zoom },
      { src: { x: mx, y: 0, w: pw, h: fh }, dx: hw + fw * zoom, dy: hh },
      { src: { x: 0, y: my, w: fw, h: ph }, dx: hw, dy: hh + fh * zoom },
      { src: { x: 0, y: 0, w: fw, h: fh }, dx: hw, dy: hh },
    ].filter((p) => p.src.w > 0 && p.src.h > 0);
    const t0 = performance.now();
    const bmps = await Promise.all(next.map((p) => client.sheet(v.id, p.src, zoom * dpr())));
    if (gen !== generation) return bmps.forEach((b) => b.close());
    bitmaps.forEach((b) => b.close());
    panes = next;
    bitmaps = bmps;
    paint();
    updateText(next[0]?.src ?? { x: 0, y: 0, w: 1, h: 1 });
    setTiming(`sheet region ${Math.round(mx)},${Math.round(my)} rendered in ${(performance.now() - t0).toFixed(0)} ms`);
  };
  sheetViews.set(v, {
    redraw: () => { if (bitmaps.length) paint(); },
    reveal: (r) => {
      const pw = stage.clientWidth - hw - fw * zoom, ph = stage.clientHeight - hh - fh * zoom;
      stage.scrollTo({
        left: r.x < fw ? stage.scrollLeft : Math.max(0, (r.x - fw) * zoom - pw / 2),
        top: r.y < fh ? stage.scrollTop : Math.max(0, (r.y - fh) * zoom - ph / 2),
        behavior: "smooth",
      });
    },
    copy: (e) => {
      const s = sel(), ds = getSelection();
      if (!s || !ds || ds.isCollapsed || !hold.contains(ds.anchorNode)) return false;
      e.preventDefault();
      const g = rangeOf(s);
      const c = clipboardFor(g);
      const done = () => setStatus(`copied ${rangeName(g)}`);
      if (c.data && e.clipboardData) {
        e.clipboardData.setData("text/plain", c.data.text);
        e.clipboardData.setData("text/html", c.data.html);
        done();
      } else {
        setStatus(`copying ${rangeName(g)}…`);
        writeClipboard(c.promise).then(done, unlessStale(gen));
      }
      return true;
    },
  });
  const redraw = () => draw().catch(unlessStale(gen));
  stage.onscroll = () => { if (current === v) redraw(); };
  // the visible region changes with the window too
  const resize = new ResizeObserver(() => { if (current === v && gen === generation) redraw(); else resize.disconnect(); });
  resize.observe(stage);
  redraw();
  // the cells selected when the sheet was last shown, ready to copy
  const s = sel();
  if (s) {
    holdSelection();
    clipboardFor(rangeOf(s));
  }
}

main().catch(showError);

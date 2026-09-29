// Demo viewer: everything is decoded and rendered in a worker; the main thread
// only places bitmaps and a selectable, accessible text layer. Files opened
// or dropped on the page are converted into bdf in another worker, by the Go
// converters built as wasm (examples/common/build.mjs builds them). A PDF is
// converted a page at a time: its pages are shown sized at once and drawn as
// they are converted, those near the visible area first. Pages are scrolled
// through, or shown one or two at a time and turned like a book's (book.ts).
// A score's music plays with Web Audio, a bar on the pages following it.
// The cells of a sheet are selected as in a spreadsheet, and copied as
// tab-separated values and an HTML table. A pinch on the stage changes the
// zoom, not the page's.
//
// The address picks the document: ?src= a bdf document (a file, or the
// directory of a split one, which ends with a slash; &range reads a file by
// ranges), ?file= a file to convert (samples/basic.docx); ?layout= the
// layout it opens in.
import { dcValues, type Manifest, type View, type SearchHit, type TextContent, type NoteEvent } from "@bdf/core";
import {
  BdfWorkerClient, BdfWorkerError, MusicPlayer, buildTextLayer, installCopyHandler, internalLink, tableCells, cellClipboard, TEXT_LAYER_CSS, RUN_ATTR,
  type Cursor, type HitRect, type OpenSource, type TextLayerOptions, type CellText, type CellRange, type CellClipboard,
} from "@bdf/render";
import { ConverterClient, ConvertError, sniff, type Opened } from "../common/convert.js";
import { Book } from "./book.js";

/** The document shown when the URL has no ?src= (set by the build); "" shows the start page. */
declare const DEFAULT_SRC: string;
/** Where lib/ (the workers and the converter modules), fonts/ and samples/ are, from the page (set by the build). */
declare const SITE_ROOT: string;

const params = new URLSearchParams(location.search);
const src = params.get("src") ?? DEFAULT_SRC;
const siteURL = (path: string) => new URL(path, location.href).href;
/** What the page shares with the other pages of the site. */
const shared = (path: string) => siteURL(SITE_ROOT + path);
const client = new BdfWorkerClient(new Worker(shared("lib/worker.js"), { type: "module" }));
/** The converter worker, started with the first file that needs converting. */
let converter: ConverterClient | undefined;
/**
 * Converter modules, one for PDF, one for the Office formats, one for HTML
 * and Markdown and one for the images browsers display by themselves, and
 * the fonts the Office converters lay text out with.
 */
const MODULES = { pdf: "lib/bdf-pdf.wasm", office: "lib/bdf-office.wasm", web: "lib/bdf-web.wasm", image: "lib/bdf-image.wasm" };
const FONTS = "fonts/";
/** The converter worker, started when it is first needed. */
const converterWorker = () => (converter ??= new ConverterClient(new Worker(shared("lib/convert-worker.js"))));

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

/** The music of the view shown (spec §4.4), once the worker has sent it. */
let player: MusicPlayer | undefined;
/** Whether the score pages or the note timeline is shown. */
let musicMode: "score" | "piano-roll" = "score";
let scoreScroll: { left: number; top: number } | undefined;
let pianoRollScroll: { left: number; top: number } | undefined;
let rollPlayhead: HTMLDivElement | undefined;
let rollTempoScale = 1;
/** Bumped as the view with music changes: music that comes for another view is dropped. */
let music = 0;
/** The bar on the page where the music is. */
const cursor = document.createElement("div");
cursor.className = "playCursor";
cursor.setAttribute("aria-hidden", "true");
/** The system the cursor was last put on: coming to another brings it into view. */
let followed: { page: number; system: number } | undefined;
/** Whether the cursor was across the visible part of the stage (zoomed in, it scrolls along). */
let cursorInView = true;
let playFrame = 0;

/** Search state for the current view; pages is how many were converted when it searched (streaming). */
const found = { query: "", hits: [] as SearchHit[], rects: [] as HitRect[][], index: -1, pages: -1 };

/** The bitmaps of the pages shown scrolled (showPages), for the pages of a stream that come in. */
let pageBitmaps: Bitmaps | undefined;

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

/**
 * The bitmaps of the pages (or bands) of a view: render is called for each
 * as it comes near the visible area, and the bitmap it puts there is let go
 * when the page is some screens away, to be drawn again when the page comes
 * back. A long document scrolled through would otherwise keep a canvas for
 * every page (8 MB for an A4 page on a screen of twice the density). The
 * text layers stay: a selection may run over pages that are far away.
 */
function nearBitmaps(render: (el: HTMLDivElement) => void) {
  const near = onNear(BITMAP_MARGIN, render);
  const far = new WeakSet<Element>();
  // three screens around the visible area, and well beyond where the bitmaps are drawn
  const margin = (px: number) => Math.max(3 * px, 1200);
  const keep = new IntersectionObserver((entries) => {
    for (const e of entries) {
      if (e.isIntersecting) {
        far.delete(e.target);
        continue;
      }
      far.add(e.target);
      const canvas = e.target.querySelector<HTMLCanvasElement>(":scope > canvas");
      if (!canvas) continue;
      canvas.width = canvas.height = 0; // its pixels go now, not when it is collected
      canvas.remove();
      near.observe(e.target);
    }
  }, { root: stage, rootMargin: `${margin(stage.clientHeight)}px ${margin(stage.clientWidth)}px` });
  return {
    observe(el: HTMLElement) {
      near.observe(el);
      keep.observe(el);
    },
    /** Whether a bitmap that comes for an element is not put there: the page has gone far meanwhile, and is drawn when it comes back. */
    gone(el: HTMLElement, bmp: ImageBitmap): boolean {
      if (!far.has(el)) return false;
      bmp.close();
      near.observe(el);
      return true;
    },
  };
}
type Bitmaps = ReturnType<typeof nearBitmaps>;

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
 * Returns the password that unlocked the document (none when it is not
 * encrypted), or undefined when it did not open or another one opened
 * meanwhile.
 */
async function load(source: OpenSource, name?: string, token = ++opening, stream?: Streaming): Promise<{ password?: string } | undefined> {
  closeDocument();
  setStatus("loading…");
  let password: string | undefined;
  const unlock = (pw: string) => client.unlock(pw).then((m) => { password = pw; return m; });
  const opened = await withPassword(() => client.open(source), unlock, "unlocking…");
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
  return { password };
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
  const base = name.replace(/\.[^.]*$/, "");
  const kind = sniff(new Uint8Array(data), name);
  if (kind === "bdf") {
    await load({ kind: "buffer", buffer: data }, name, token);
    return;
  }
  const busy = `converting ${name}…`;
  setStatus(busy);
  const conv = converterWorker();
  const module = shared(MODULES[kind]);
  // only the Office converters lay text out with the font directory (PDFs embed their fonts, images have no text)
  const fonts = kind === "office" ? shared(FONTS) : undefined;
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
    setDownload(bdfFile(res.bdf), base);
    if (!(await load({ kind: "buffer", buffer }, name, token))) return;
    setStatus(took);
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
  convertPages(st, name, base, t0).catch(openFailed);
}

/**
 * Convert the pages of a stream one by one, each time the one nearest to
 * the visible area, and put them in the document shown; then swap in the
 * finished document, which is also the one to download.
 */
async function convertPages(st: Streaming, name: string, base: string, t0: number) {
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
    setDownload(bdfFile(res.bdf), base);
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
    renderPage(current, index, el, gen, pageBitmaps).catch(unlessStale(gen));
  }
  if (el.dataset.waitText !== undefined) {
    delete el.dataset.waitText;
    renderText(current, index, el, gen).catch(unlessStale(gen));
  }
}

/** Take down the document shown, before another one opens: work started for it is dropped. */
function closeDocument() {
  stopMusic();
  if (streaming) converter?.close(streaming.id).catch(() => {});
  streaming = undefined;
  showProgress();
  current = undefined;
  generation++;
  book?.destroy();
  book = undefined;
  sheetChunks?.drop();
  sheetChunks = undefined;
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

/** A converted document as a file; the Blob copies the bytes, which then go to the worker. */
const bdfFile = (bdf: Uint8Array) => new Blob([bdf as BlobPart], { type: "application/octet-stream" });

/** The object URL of the download link: revoked when replaced. */
let downloadURL = "";

/** Offer the converted document for download (none without file); base is the name without its extension. */
function setDownload(file?: Blob, base = "") {
  const a = $<HTMLAnchorElement>("download");
  if (downloadURL) URL.revokeObjectURL(downloadURL);
  downloadURL = "";
  a.hidden = !file;
  if (!file) {
    a.removeAttribute("href");
    return;
  }
  a.href = downloadURL = URL.createObjectURL(file);
  a.download = `${base}.bdf`;
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

/** Fetch a file by its URL and open it, converted: a sample, or the file ?file= names. */
function openURL(url: string, name = decodeURIComponent(new URL(url).pathname.split("/").pop() ?? "")) {
  setStatus(`fetching ${name}…`);
  return fetch(url)
    .then((r) => (r.ok ? r.arrayBuffer() : Promise.reject(new Error(`${name}: HTTP ${r.status}`))))
    .then((data) => openFile(name, data))
    .catch(openFailed);
}

/** The start page: a file picker, drag and drop, and the samples published with the page. */
async function showLanding() {
  landing.hidden = false;
  setStatus("no document open");
  const res = await fetch(shared("samples/index.json")).catch(() => undefined);
  if (!res?.ok) return;
  const samples = (await res.json()) as { name: string; label: string }[];
  const box = landing.querySelector<HTMLDivElement>("#samples")!;
  box.querySelector(".samples")!.replaceChildren(...samples.map((s) => {
    const b = document.createElement("button");
    b.type = "button";
    b.textContent = s.label;
    b.title = s.name;
    b.onclick = () => openURL(shared(`samples/${s.name}`), s.name);
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
  zoomInput.oninput = () => setZoom(Number(zoomInput.value));
  initPinch();
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
  initMusic();
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

/** The zoom slider's range and step: a pinch ends on a step of it too. */
const ZOOM_MIN = Number(zoomInput.min), ZOOM_MAX = Number(zoomInput.max), ZOOM_STEP = Number(zoomInput.step);

/** Put a zoom in the slider's label (and what screen readers say of it). */
function showZoom(z: number) {
  const pct = `${Math.round(z * 100)}%`;
  $("zoomv").textContent = pct;
  zoomInput.setAttribute("aria-valuetext", pct);
}

/** Zoom to z (from the slider, or a pinch) and show the view anew. */
function setZoom(z: number) {
  zoom = z;
  zoomInput.value = String(z);
  showZoom(z);
  if (current) show(current);
}

/** Safari's pinch events (a trackpad's, and on iOS a touch screen's too): not in the DOM typings. */
interface GestureEvent extends UIEvent { readonly scale: number; readonly clientX: number; readonly clientY: number }

/**
 * A pinch on the stage zooms the view, not the page: two fingers on a touch
 * screen, and on a trackpad ctrl+wheel (Chrome, Edge, Firefox; a mouse wheel
 * with ctrl too) or gesture events (Safari). While it goes on the stage is
 * only scaled, and the label shows the zoom it comes to; when it ends the
 * zoom takes the nearest step of the slider and the view is shown anew,
 * with the point pinched about still under the fingers.
 */
function initPinch() {
  interface Pinch {
    /** What it comes from: iOS sends gesture events along with the touches, which are left alone. */
    by: "touch" | "wheel" | "gesture";
    view: View;
    /** The zoom it started from, and the scale over it so far. */
    from: number;
    scale: number;
    /** Where it started and where it is now (the fingers pan too), client px. */
    x0: number; y0: number; x: number; y: number;
    /** The page (or band) where it started, and where on it (0–1). */
    anchor?: { key: string; fx: number; fy: number };
    /** Where it started in the stage's content, for a view with no page there (a sheet). */
    cx: number; cy: number;
  }
  let pinch: Pinch | undefined;
  const begin = (by: Pinch["by"], x: number, y: number): Pinch => {
    const s = stage.getBoundingClientRect();
    const page = document.elementsFromPoint(x, y)
      .map((el) => el.closest<HTMLElement>(".page[data-index], .page[data-y]"))
      .find((el) => el && stage.contains(el));
    let anchor: Pinch["anchor"];
    if (page) {
      const r = page.getBoundingClientRect();
      const key = page.dataset.index !== undefined ? `.page[data-index="${page.dataset.index}"]` : `.page[data-y="${page.dataset.y}"]`;
      anchor = { key, fx: (x - r.left) / r.width, fy: (y - r.top) / r.height };
    }
    return {
      by, view: current!, from: zoom, scale: 1, x0: x, y0: y, x, y, anchor,
      cx: stage.scrollLeft + x - s.left - stage.clientLeft, cy: stage.scrollTop + y - s.top - stage.clientTop,
    };
  };
  const update = (p: Pinch, scale: number, x: number, y: number) => {
    p.scale = Math.min(ZOOM_MAX / p.from, Math.max(ZOOM_MIN / p.from, scale));
    p.x = x;
    p.y = y;
    showZoom(p.from * p.scale);
    for (const el of [...stage.children] as HTMLElement[]) {
      el.style.transformOrigin = `${p.cx - el.offsetLeft}px ${p.cy - el.offsetTop}px`;
      el.style.transform = `translate(${x - p.x0}px, ${y - p.y0}px) scale(${p.scale})`;
    }
  };
  const end = (p: Pinch) => {
    pinch = undefined;
    for (const el of [...stage.children] as HTMLElement[]) el.style.transform = el.style.transformOrigin = "";
    const z = Number((Math.round(Math.min(ZOOM_MAX, Math.max(ZOOM_MIN, p.from * p.scale)) / ZOOM_STEP) * ZOOM_STEP).toFixed(2));
    if (p.view !== current || z === zoom) return showZoom(zoom);
    setZoom(z);
    const el = p.anchor && stage.querySelector<HTMLElement>(p.anchor.key);
    if (el) {
      const r = el.getBoundingClientRect();
      stage.scrollLeft += r.left + p.anchor!.fx * r.width - p.x;
      stage.scrollTop += r.top + p.anchor!.fy * r.height - p.y;
    } else {
      const s = stage.getBoundingClientRect(), k = z / p.from;
      stage.scrollLeft = p.cx * k - (p.x - s.left - stage.clientLeft);
      stage.scrollTop = p.cy * k - (p.y - s.top - stage.clientTop);
    }
  };

  // two fingers: what the first went down on (a page's corner, a cell) forgets its drag or tap
  const down = new Map<number, EventTarget>();
  stage.addEventListener("pointerdown", (e) => { if (e.pointerType === "touch") down.set(e.pointerId, e.target!); }, true);
  const lift = (e: PointerEvent) => { if (e.isTrusted) down.delete(e.pointerId); };
  addEventListener("pointerup", lift, true);
  addEventListener("pointercancel", lift, true);
  const span = (t: TouchList) => ({
    d: Math.max(1, Math.hypot(t[1].clientX - t[0].clientX, t[1].clientY - t[0].clientY)),
    x: (t[0].clientX + t[1].clientX) / 2, y: (t[0].clientY + t[1].clientY) / 2,
  });
  let d0 = 1;
  stage.addEventListener("touchstart", (e) => {
    if (e.touches.length !== 2 || !current || pinch?.by === "touch" || pinch?.by === "wheel") return;
    e.preventDefault();
    const s = span(e.touches);
    d0 = s.d;
    pinch = begin("touch", s.x, s.y);
    for (const [pointerId, target] of down) target.dispatchEvent(new PointerEvent("pointercancel", { pointerId, pointerType: "touch", bubbles: true }));
    down.clear();
  }, { passive: false });
  stage.addEventListener("touchmove", (e) => {
    if (pinch?.by !== "touch") return;
    e.preventDefault();
    if (e.touches.length < 2) return;
    const s = span(e.touches);
    update(pinch, s.d / d0, s.x, s.y);
  }, { passive: false });
  const touchEnd = (e: TouchEvent) => { if (pinch?.by === "touch" && e.touches.length < 2) end(pinch); };
  stage.addEventListener("touchend", touchEnd);
  stage.addEventListener("touchcancel", touchEnd);

  // a trackpad's pinch comes as small steps (Chrome's scale is e^(-deltaY/100)), a wheel's as large ones: a notch is about ×1.28
  let wheelTimer = 0;
  stage.addEventListener("wheel", (e) => {
    if (!e.ctrlKey || !current) return;
    e.preventDefault();
    if (pinch && pinch.by !== "wheel") return;
    const p = (pinch ??= begin("wheel", e.clientX, e.clientY));
    const dy = e.deltaY * (e.deltaMode === WheelEvent.DOM_DELTA_LINE ? 16 : e.deltaMode === WheelEvent.DOM_DELTA_PAGE ? stage.clientHeight : 1);
    update(p, p.scale * Math.exp(-Math.max(-25, Math.min(25, dy)) / 100), p.x0, p.y0);
    clearTimeout(wheelTimer);
    wheelTimer = setTimeout(() => { if (pinch === p) end(p); }, 250);
  }, { passive: false });

  stage.addEventListener("gesturestart", (e) => {
    e.preventDefault();
    if (!current || pinch) return;
    const g = e as GestureEvent;
    pinch = begin("gesture", g.clientX, g.clientY);
  });
  stage.addEventListener("gesturechange", (e) => {
    e.preventDefault();
    if (pinch?.by === "gesture") update(pinch, (e as GestureEvent).scale, pinch.x0, pinch.y0);
  });
  stage.addEventListener("gestureend", (e) => {
    e.preventDefault();
    if (pinch?.by === "gesture") end(pinch);
  });
}

async function main() {
  init();
  // a file to convert: the samples are named from the site's root (samples/basic.docx)
  const file = params.get("file");
  if (file) return openURL(/^samples\//.test(file) ? shared(file) : siteURL(file));
  if (!src) return showLanding();
  const source: OpenSource = src.endsWith("/") ? { kind: "split", base: siteURL(src) } : { kind: "single", url: siteURL(src), range: params.has("range") };
  await load(source);
}

function show(v: View) {
  // a layout or zoom change goes on from the page in view
  const keep = current === v ? pageInView() : 0;
  const priorRollScroll = current === v && musicMode === "piano-roll" && stage.querySelector(".pianoRoll")
    ? { left: stage.scrollLeft, top: stage.scrollTop } : undefined;
  if (current !== v) {
    found.query = ""; found.hits = []; found.rects = []; found.index = -1; found.pages = -1; hitsBox.textContent = "";
    setStatus(describe(v));
    stopMusic();
    musicMode = "score";
    scoreScroll = pianoRollScroll = undefined;
    if (v.play) loadMusic(v);
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
  if (sheetChunks?.view !== v) {
    sheetChunks?.drop();
    sheetChunks = undefined;
  }
  $("pageNav").hidden = $("animateBox").hidden = true;
  stage.onscroll = stage.onkeydown = stage.onfocus = null;
  stage.replaceChildren();
  stage.scrollTop = stage.scrollLeft = 0;
  visibility?.disconnect();
  visible.clear();
  if (musicMode === "piano-roll" && player) {
    showPianoRoll();
    const at = priorRollScroll ?? pianoRollScroll;
    if (at) stage.scrollTo(at);
  } else if (v.kind === "sheet") showSheet(v);
  else if (continuous(v)) showContinuous(v);
  else if (inBook(v)) showBook(v, keep);
  else {
    showPages(v);
    if (keep > 0) stage.querySelector(`.page[data-index="${keep}"]`)?.scrollIntoView({ block: "start" });
  }
  // the cursor goes on the new layout, which is brought to it
  followed = undefined;
  if (player) follow();
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
  const bitmaps = nearBitmaps((el) => {
    const i = Number(el.dataset.index);
    if (pending(v, i)) el.dataset.waitBitmap = "";
    else renderPage(v, i, el, gen, bitmaps).catch(unlessStale(gen));
  });
  pageBitmaps = bitmaps;
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

async function renderPage(v: View, index: number, el: HTMLDivElement, gen: number, bitmaps?: Bitmaps) {
  const page = v.pages![index];
  const t0 = performance.now();
  const bmp = await client.page(v.id, index, zoom * dpr());
  if (gen !== generation) return bmp.close();
  if (bitmaps?.gone(el, bmp)) return;
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
    // while the music plays or pauses, a tap on a system plays from there
    tap: (i, x, y) => !!player && player.state !== "stopped" && playFrom(i, x, y),
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
  // the pages were made anew: the cursor goes back on (while playing, the next frame puts it)
  if (player && !player.playing) follow();
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

// Music (spec §4.4): a view of a score can play. Its controls are in the
// header; Space plays and pauses. While it plays or pauses, a bar the height
// of the system shows where the music is (in the scrolled pages and in the
// book layouts), and the pages follow it: the stage scrolls to the system the
// music comes to, the book turns to its page. A click on a system plays from
// there.

/** The controls, the keys, and clicks on the pages. */
function initMusic() {
  $("play").onclick = togglePlay;
  $("stopPlay").onclick = () => player?.stop();
  $("musicMode").onclick = toggleMusicMode;
  $("metronome").onclick = toggleMetronome;
  const speedPercent = $<HTMLInputElement>("speedPercent");
  speedPercent.onchange = applySpeedPercent;
  speedPercent.onkeydown = (e) => { if (e.key === "Enter") { e.preventDefault(); applySpeedPercent(); } };
  // Space plays and pauses (before the stage scrolls or the book turns), unless typing or on a control
  window.addEventListener("keydown", (e) => {
    if (e.key !== " " || !player || e.defaultPrevented || e.altKey || e.ctrlKey || e.metaKey || e.shiftKey) return;
    const t = e.target as HTMLElement;
    if (t.isContentEditable || t.closest?.("input, textarea, select, button, summary, dialog")) return;
    e.preventDefault();
    if (!e.repeat) togglePlay();
  }, true);
  // a click on a system of the scrolled pages plays from there (the book layouts: see showBook);
  // not the end of a drag, nor a click that clears a selection
  let press: { x: number; y: number; selected: boolean } | undefined;
  stage.addEventListener("pointerdown", (e) => {
    press = { x: e.clientX, y: e.clientY, selected: !(getSelection()?.isCollapsed ?? true) };
  });
  stage.addEventListener("click", (e) => {
    const p = press;
    press = undefined;
    if (!player || book || !p || p.selected || e.button !== 0 || Math.hypot(e.clientX - p.x, e.clientY - p.y) >= 4) return;
    const target = e.target as Element;
    const el = target.closest<HTMLElement>(".page[data-index]");
    if (!el || target.closest("a") || !(getSelection()?.isCollapsed ?? true)) return;
    const r = el.getBoundingClientRect();
    playFrom(Number(el.dataset.index), (e.clientX - r.left) / zoom, (e.clientY - r.top) / zoom);
  });
}

/** Ask the worker for the music of a view, and make its player. */
function loadMusic(v: View) {
  const token = music;
  $("playBox").hidden = false;
  $<HTMLButtonElement>("play").disabled = $<HTMLButtonElement>("stopPlay").disabled = true;
  $<HTMLButtonElement>("musicMode").disabled = $<HTMLButtonElement>("metronome").disabled = true;
  $<HTMLInputElement>("speedPercent").disabled = true;
  $("playTime").textContent = "";
  client.play(v.id).then((data) => {
    if (token !== music || !data) return;
    player = new MusicPlayer(new Uint8Array(data.seq), data.cues);
    player.onUpdate = musicChanged;
    musicChanged();
  }).catch((e) => {
    if (token === music) setStatus(`the music could not be read: ${(e as Error).message ?? e}`);
  });
}

/** Take the music of the view shown down: another view or document is shown. */
function stopMusic() {
  music++;
  player?.dispose();
  player = undefined;
  cancelAnimationFrame(playFrame);
  playFrame = 0;
  cursor.remove();
  rollPlayhead = undefined;
  followed = undefined;
  $("playBox").hidden = true;
}

function togglePlay() {
  const p = player;
  if (!p) return;
  if (p.playing) return p.pause();
  try {
    p.play().catch(showError);
  } catch (e) {
    showError(e); // no Web Audio
  }
}

/** Switch between the engraved page and its time-by-pitch note view. */
function toggleMusicMode() {
  if (!player || !current) return;
  if (musicMode === "score") {
    scoreScroll = { left: stage.scrollLeft, top: stage.scrollTop };
    musicMode = "piano-roll";
    show(current);
    if (pianoRollScroll) stage.scrollTo(pianoRollScroll);
  } else {
    pianoRollScroll = { left: stage.scrollLeft, top: stage.scrollTop };
    musicMode = "score";
    show(current);
    if (scoreScroll) stage.scrollTo(scoreScroll);
  }
  const mode = $<HTMLButtonElement>("musicMode");
  mode.textContent = musicMode === "score" ? "Piano roll" : "Score";
  mode.title = musicMode === "score" ? "Show piano roll" : "Show score pages";
}

function toggleMetronome() {
  if (player) player.setMetronome(!player.metronomeEnabled);
}

function applySpeedPercent() {
  const p = player;
  const input = $<HTMLInputElement>("speedPercent");
  const value = Number(input.value);
  if (!p || !input.value || !Number.isFinite(value)) {
    if (p) input.value = String(p.speedPercent);
    return;
  }
  p.setSpeedPercent(value);
  input.value = String(p.speedPercent);
}

/** The player started, paused, stopped or jumped: the controls and the cursor. */
function musicChanged() {
  const p = player;
  if (!p) return;
  const b = $<HTMLButtonElement>("play");
  b.disabled = false;
  b.textContent = p.playing ? "❚❚" : "▶";
  b.title = p.playing ? "pause (Space)" : "play (Space)";
  b.setAttribute("aria-label", p.playing ? "pause" : "play");
  $<HTMLButtonElement>("stopPlay").disabled = p.state === "stopped";
  const mode = $<HTMLButtonElement>("musicMode");
  mode.disabled = false;
  mode.textContent = musicMode === "score" ? "Piano roll" : "Score";
  mode.title = musicMode === "score" ? "Show piano roll" : "Show score pages";
  const metronome = $<HTMLButtonElement>("metronome");
  metronome.disabled = !p.metronomeAvailable;
  metronome.textContent = p.metronomeEnabled ? "Metronome on" : "Metronome off";
  metronome.setAttribute("aria-pressed", String(p.metronomeEnabled));
  metronome.title = p.metronomeAvailable ? "Toggle metronome" : "Unavailable for SMPTE-timed MIDI";
  const speedPercent = $<HTMLInputElement>("speedPercent");
  speedPercent.disabled = false;
  speedPercent.value = String(p.speedPercent);
  const openingBpm = $("openingBpm");
  openingBpm.hidden = !p.bpmAvailable;
  openingBpm.textContent = p.bpmAvailable ? `${Number(p.bpm.toFixed(1))} BPM at start` : "";
  if (musicMode === "piano-roll" && rollTempoScale !== p.tempoScale) {
    const at = { left: stage.scrollLeft, top: stage.scrollTop };
    showPianoRoll();
    stage.scrollTo(at);
  }
  follow();
}

/** The time and the cursor, every frame while the music plays. */
function follow() {
  cancelAnimationFrame(playFrame);
  playFrame = 0;
  const p = player;
  if (!p) return;
  const at = p.position;
  $("playTime").textContent = `${minutes(at)} / ${minutes(p.duration)}`;
  if (musicMode === "piano-roll") {
    cursor.remove();
    if (rollPlayhead) {
      const x = pianoRollX(at);
      rollPlayhead.style.left = `${x}px`;
      if (p.playing) {
        const visibleX = x - stage.scrollLeft;
        if (visibleX < 24 || visibleX > stage.clientWidth - 24) {
          stage.scrollTo({ left: Math.max(0, x - stage.clientWidth / 3), top: stage.scrollTop, behavior: "auto" });
        }
      }
    }
  } else placeCursor(p.state === "stopped" ? null : p.cursorAt(at));
  if (p.playing) playFrame = requestAnimationFrame(follow);
}

const ROLL_LABEL_WIDTH = 72;
let rollPixelsPerSecond = 64;

function pianoRollX(seconds: number): number { return ROLL_LABEL_WIDTH + seconds * rollPixelsPerSecond; }

/** Draw the sequence as a scrollable time-by-pitch grid. Clicking seeks the same player as the score. */
function showPianoRoll() {
  const p = player;
  if (!p) return;
  const notes: NoteEvent[] = [];
  let low = 127, high = 0;
  for (const e of p.sequence.events) {
    if (e.type !== "note") continue;
    notes.push(e);
    low = Math.min(low, e.key);
    high = Math.max(high, e.key);
  }
  if (!notes.length) { low = 48; high = 84; }
  low = Math.max(0, Math.floor(low / 12) * 12 - 1);
  high = Math.min(127, Math.ceil(high / 12) * 12 + 1);
  const rowHeight = Math.max(12, Math.min(24, 14 * zoom));
  const rulerHeight = 30;
  const height = rulerHeight + (high - low + 1) * rowHeight;
  const duration = Math.max(p.duration, 1);
  rollPixelsPerSecond = Math.min(64 * zoom, (14000 - ROLL_LABEL_WIDTH) / duration);
  rollTempoScale = p.tempoScale;
  const width = Math.ceil(Math.max(stage.clientWidth, Math.min(14000, ROLL_LABEL_WIDTH + duration * rollPixelsPerSecond + 32)));
  const dpr = Math.min(window.devicePixelRatio || 1, 2);
  const canvas = document.createElement("canvas");
  canvas.width = Math.round(width * dpr);
  canvas.height = Math.round(height * dpr);
  canvas.style.width = `${width}px`;
  canvas.style.height = `${height}px`;
  canvas.setAttribute("role", "img");
  canvas.setAttribute("aria-label", `Piano roll with ${notes.length} notes from ${keyName(high)} to ${keyName(low)}. Click to seek.`);
  canvas.tabIndex = 0;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;
  ctx.scale(dpr, dpr);
  ctx.fillStyle = "#fff";
  ctx.fillRect(0, 0, width, height);
  ctx.font = "11px system-ui, sans-serif";
  ctx.textBaseline = "middle";
  const isBlack = (key: number) => [1, 3, 6, 8, 10].includes(key % 12);
  const yFor = (key: number) => rulerHeight + (high - key) * rowHeight;
  for (let key = low; key <= high; key++) {
    const y = yFor(key);
    ctx.fillStyle = isBlack(key) ? "#edf0f4" : "#fff";
    ctx.fillRect(0, y, width, rowHeight);
    ctx.fillStyle = isBlack(key) ? "#40454d" : "#f5f6f7";
    ctx.fillRect(0, y, ROLL_LABEL_WIDTH, rowHeight);
    if (key % 12 === 0) {
      ctx.fillStyle = "#343a40";
      ctx.fillText(keyName(key), 8, y + rowHeight / 2);
    }
    ctx.strokeStyle = "#d8dde3";
    ctx.beginPath(); ctx.moveTo(0, y + rowHeight); ctx.lineTo(width, y + rowHeight); ctx.stroke();
  }
  ctx.fillStyle = "#f5f6f7";
  ctx.fillRect(0, 0, width, rulerHeight);
  const gridSeconds = Math.max(1, Math.ceil(40 / rollPixelsPerSecond));
  const labelEvery = gridSeconds * Math.max(1, Math.ceil(5 / gridSeconds));
  for (let second = 0; second <= duration; second += gridSeconds) {
    const x = pianoRollX(second);
    const major = second % labelEvery === 0;
    ctx.strokeStyle = major ? "#9ca6b2" : "#e2e6eb";
    ctx.beginPath(); ctx.moveTo(x, rulerHeight); ctx.lineTo(x, height); ctx.stroke();
    if (major) {
      ctx.fillStyle = "#3d4650";
      ctx.fillText(minutes(second), x + 3, rulerHeight / 2);
    }
  }
  for (const note of notes) {
    const x = pianoRollX(note.time * p.tempoScale);
    const right = Math.min(width, x + Math.max(2, note.duration * p.tempoScale * rollPixelsPerSecond));
    if (right <= ROLL_LABEL_WIDTH) continue;
    const y = yFor(note.key) + 2;
    ctx.fillStyle = `hsl(${(note.channel * 47 + 205) % 360} 62% 52% / .82)`;
    ctx.fillRect(Math.max(ROLL_LABEL_WIDTH, x), y, right - Math.max(ROLL_LABEL_WIDTH, x), rowHeight - 4);
  }
  ctx.fillStyle = "#59636f";
  ctx.fillText(`${notes.length} notes`, 8, rulerHeight / 2);
  canvas.addEventListener("click", (e) => {
    if (player !== p) return;
    const x = e.clientX - canvas.getBoundingClientRect().left;
    if (x >= ROLL_LABEL_WIDTH) p.seek(Math.max(0, Math.min(p.duration, (x - ROLL_LABEL_WIDTH) / rollPixelsPerSecond)));
  });
  const roll = document.createElement("div");
  roll.className = "pianoRoll";
  roll.style.width = `${width}px`;
  roll.style.height = `${height}px`;
  roll.append(canvas);
  rollPlayhead = document.createElement("div");
  rollPlayhead.className = "rollPlayhead";
  rollPlayhead.setAttribute("aria-hidden", "true");
  roll.append(rollPlayhead);
  stage.replaceChildren(roll);
}

const NOTE_NAMES = ["C", "C♯", "D", "D♯", "E", "F", "F♯", "G", "G♯", "A", "A♯", "B"];
function keyName(key: number): string { return `${NOTE_NAMES[key % 12]}${Math.floor(key / 12) - 1}`; }

/** Seconds as m:ss (h:mm:ss from an hour). */
function minutes(seconds: number): string {
  const s = Math.floor(seconds);
  const ss = String(s % 60).padStart(2, "0");
  return s < 3600 ? `${Math.floor(s / 60)}:${ss}` : `${Math.floor(s / 3600)}:${String(Math.floor(s / 60) % 60).padStart(2, "0")}:${ss}`;
}

/**
 * Put the bar where the music is, the height of its system, on its page.
 * Coming to another system brings it into view: the book turns to its page,
 * the scrolled pages scroll to it.
 */
function placeCursor(c: Cursor | null) {
  const v = current;
  if (!c || !v?.pages?.[c.page] || continuous(v)) {
    cursor.remove();
    if (!c) followed = undefined;
    return;
  }
  const moved = !followed || followed.page !== c.page || followed.system !== c.system;
  followed = { page: c.page, system: c.system };
  if (book && moved && !book.pages.includes(c.page)) book.go(c.page, true);
  const el = stage.querySelector<HTMLElement>(`.page[data-index="${c.page}"]`);
  if (!el) {
    cursor.remove();
    return;
  }
  const scale = book?.scale ?? zoom;
  if (cursor.parentElement !== el) el.append(cursor);
  cursor.style.left = `${c.x * scale}px`;
  cursor.style.top = `${c.y * scale}px`;
  cursor.style.height = `${c.h * scale}px`;
  if (!book) reveal(el, c, moved);
}

/**
 * Scroll the stage to the system the music has come to when it is not all
 * in view, and along with the cursor when it runs off the side (zoomed in).
 */
function reveal(el: HTMLElement, c: Cursor, moved: boolean) {
  const s = stage.getBoundingClientRect(), r = el.getBoundingClientRect();
  const top = r.top - s.top + c.y * zoom, h = c.h * zoom;
  const x = r.left - s.left + c.x * zoom;
  const vw = stage.clientWidth, vh = stage.clientHeight;
  const across = x >= 16 && x <= vw - 16;
  let left = stage.scrollLeft, down = stage.scrollTop;
  // a quarter down the stage, or centered when it is tall
  if (moved && (top < 0 || top + h > vh)) down += top - Math.max(0, Math.min(vh / 4, (vh - h) / 2));
  if (!across && (moved || cursorInView)) left += x - vw / 4;
  cursorInView = across;
  if (left !== stage.scrollLeft || down !== stage.scrollTop) stage.scrollTo({ left, top: down, behavior: "smooth" });
}

/** Play from the place of a system clicked on (page units); false when no system of the music is there. */
function playFrom(page: number, x: number, y: number): boolean {
  const p = player;
  if (!p?.cues) return false;
  const pad = 4;
  const i = p.cues.systems.findIndex((s) => s.page === page && x >= s.x - pad && x <= s.x + s.w + pad && y >= s.y - pad && y <= s.y + s.h + pad);
  const t = i < 0 ? null : p.timeAt(i, x);
  if (t === null) return false;
  p.seek(t);
  return true;
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
  const bitmaps = nearBitmaps((el) => {
    const vp = viewportOf(el);
    client.continuous(v.id, vp, zoom * dpr()).then((bmp) => {
      if (gen !== generation) return bmp.close();
      if (!bitmaps.gone(el, bmp)) placeBitmap(el, bmp, vp.w, vp.h);
    }).catch(unlessStale(gen));
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

/** The side of the chunks a sheet is drawn in, in CSS px. */
const SHEET_CHUNK = 512;
/** Bytes of chunk bitmaps kept (width × height × 4), unless those in view take more. */
const SHEET_CHUNK_BUDGET = 192 * 1024 * 1024;
/** Chunks asked of the worker at once. */
const SHEET_CHUNK_REQUESTS = 4;

/**
 * The chunks a sheet is drawn in at one scale: squares of SHEET_CHUNK CSS px
 * from the sheet's origin, each a region the worker draws once. Scrolling
 * puts the chunks drawn already on the canvas and asks only for those that
 * come into view: drawing the whole visible region again at every step of
 * a scroll took the worker, and the transfer of its bitmap, longer than a
 * frame. The chunks are whole device pixels, so that they meet without
 * seams. The least recently shown go beyond the budget.
 */
class SheetChunks {
  /** A chunk's side in device pixels, and in sheet units. */
  readonly px: number;
  readonly unit: number;
  /** Bitmaps by key, least recently shown first. */
  private bitmaps = new Map<string, ImageBitmap>();
  private bytes = 0;
  private loading = new Set<string>();
  private dropped = false;

  constructor(readonly view: View, readonly scale: number, d: number) {
    this.px = Math.max(1, Math.round(SHEET_CHUNK * d));
    this.unit = this.px / scale;
  }

  static key(cx: number, cy: number) { return `${cx},${cy}`; }

  /** The chunks a region of the sheet (in units) lies in, row by row. */
  within(r: { x: number; y: number; w: number; h: number }, out: [number, number][] = []): [number, number][] {
    if (r.w <= 0 || r.h <= 0) return out;
    const u = this.unit;
    for (let cy = Math.max(0, Math.floor(r.y / u)); cy * u < r.y + r.h; cy++) {
      for (let cx = Math.max(0, Math.floor(r.x / u)); cx * u < r.x + r.w; cx++) out.push([cx, cy]);
    }
    return out;
  }

  /** A chunk drawn already, now the most recently shown. */
  get(cx: number, cy: number): ImageBitmap | undefined {
    const key = SheetChunks.key(cx, cy);
    const b = this.bitmaps.get(key);
    if (b) {
      this.bitmaps.delete(key);
      this.bitmaps.set(key, b);
    }
    return b;
  }

  /**
   * Ask the worker for the first of the chunks wanted that are neither
   * drawn nor being drawn, SHEET_CHUNK_REQUESTS at a time; arrived runs as
   * each comes (the chunks in view are those wanted first, and kept).
   */
  fetch(wanted: [number, number][], keep: number, arrived: (ms: number) => void, failed: (e: unknown) => void) {
    for (const [cx, cy] of wanted) {
      if (this.loading.size >= SHEET_CHUNK_REQUESTS) return;
      const key = SheetChunks.key(cx, cy);
      if (this.bitmaps.has(key) || this.loading.has(key)) continue;
      this.loading.add(key);
      const t0 = performance.now();
      client.sheet(this.view.id, { x: cx * this.unit, y: cy * this.unit, w: this.unit, h: this.unit }, this.scale).then((b) => {
        this.loading.delete(key);
        if (this.dropped) return b.close();
        this.bitmaps.set(key, b);
        this.bytes += b.width * b.height * 4;
        this.trim(Math.max(SHEET_CHUNK_BUDGET, 2 * keep * this.px * this.px * 4));
        arrived(performance.now() - t0);
      }, (e) => {
        this.loading.delete(key);
        failed(e);
      });
    }
  }

  private trim(budget: number) {
    for (const [key, b] of this.bitmaps) {
      if (this.bytes <= budget) return;
      this.bitmaps.delete(key);
      this.bytes -= b.width * b.height * 4;
      b.close();
    }
  }

  /** Let the bitmaps go: the sheet is not shown at this scale any more. */
  drop() {
    this.dropped = true;
    for (const b of this.bitmaps.values()) b.close();
    this.bitmaps.clear();
    this.bytes = 0;
  }
}

/** The chunks of the sheet shown last, kept while it is shown again at the same scale. */
let sheetChunks: SheetChunks | undefined;

/** A cell's name: B12. */
const cellName = (r: number, c: number) => `${columnLabel(c)}${r + 1}`;

/**
 * Sheet: a scroll area the size of the sheet with one sticky canvas that
 * shows the column and row headers and the visible region, split into the
 * frozen panes of the view (each put together from the chunks the worker
 * drew, see SheetChunks). A text layer (a table of the cells, for screen
 * readers) covers the visible region and a screen around it; the cells of
 * the frozen panes are pinned with a translation that follows the scroll.
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
  // Cells of the frozen panes stay where they are while the sheet scrolls:
  // each is moved by the scroll offset. (A custom property on the layer
  // would restyle all its runs at every step of a scroll.)
  let pinned: { span: HTMLSpanElement; x: boolean; y: boolean }[] = [];
  let building: typeof pinned = [];
  /** The scroll offset the pinned cells are moved by (read once: reading it after a move lays the page out again). */
  let scrolled = { x: 0, y: 0 };
  const place = (p: (typeof pinned)[number]) => {
    p.span.style.translate = `${p.x ? scrolled.x : 0}px ${p.y ? scrolled.y : 0}px`;
  };
  const pin = (span: HTMLSpanElement, r: { x: number; y: number }) => {
    const x = r.x < fw, y = r.y < fh;
    if (!x && !y) return;
    const p = { span, x, y };
    place(p);
    span.style.zIndex = "2"; // above the other runs (TEXT_LAYER_CSS)
    building.push(p);
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
        building = [];
        const layer = buildTextLayer(content, zoom, { ...layerOptions(), sheet, onSpan: pin, selectable: false });
        pinned = building;
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
      // a rectangle of too many cells throws: copy reports it, as it does what the worker fails at
      try {
        entry.data = cellClipboard(tableCells(loaded.content), g, opts);
        entry.promise = Promise.resolve(entry.data);
      } catch (e) {
        entry.promise = Promise.reject(e);
      }
    } else {
      entry.promise = client.sheetContent(v.id, box).then((c) => (entry.data = cellClipboard(tableCells(c), g, opts)));
    }
    entry.promise.catch(() => { if (clip === entry) clip = undefined; }); // tried again on copy, which reports the error
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
    if (panes.length) paint();
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
    if (panes.length) paint();
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
  const paint = () => {
    const d = dpr();
    const vw = Math.min(stage.clientWidth, hw + cols.total * zoom), vh = Math.min(stage.clientHeight, hh + rows.total * zoom);
    const cw = Math.ceil(vw * d), ch = Math.ceil(vh * d);
    // a new size clears the canvas and its state (the same one would, too, at a cost on every scroll)
    if (canvas.width !== cw || canvas.height !== ch) {
      canvas.width = cw;
      canvas.height = ch;
      canvas.style.width = `${vw}px`;
      canvas.style.height = `${vh}px`;
    }
    const ctx = canvas.getContext("2d")!;
    ctx.fillStyle = "#fff";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    // the chunks drawn already, at whole device pixels; those still to come show the gridlines meanwhile
    const chunks = sheetChunks;
    if (chunks?.view === v) {
      const k = chunks.scale, px = chunks.px, u = chunks.unit;
      ctx.strokeStyle = "#d9d9d9"; // as the worker draws them (PageRenderer)
      ctx.lineWidth = 1;
      for (const p of panes) {
        const x0 = Math.round(p.dx * d), y0 = Math.round(p.dy * d);
        ctx.save();
        ctx.beginPath();
        ctx.rect(x0, y0, Math.round((p.dx + p.src.w * zoom) * d) - x0, Math.round((p.dy + p.src.h * zoom) * d) - y0);
        ctx.clip();
        const bx = Math.round(p.dx * d - p.src.x * k), by = Math.round(p.dy * d - p.src.y * k);
        ctx.beginPath();
        for (const [cx, cy] of chunks.within(p.src)) {
          const b = chunks.get(cx, cy);
          if (b) {
            ctx.drawImage(b, bx + cx * px, by + cy * px);
            continue;
          }
          if (!v.gridlines) continue;
          const top = by + cy * px, left = bx + cx * px;
          cols.each(cx * u, (cx + 1) * u, (_, start, size) => {
            if (start + size > (cx + 1) * u) return;
            const x = Math.round(bx + (start + size) * k) + 0.5;
            ctx.moveTo(x, top);
            ctx.lineTo(x, top + px);
          });
          rows.each(cy * u, (cy + 1) * u, (_, start, size) => {
            if (start + size > (cy + 1) * u) return;
            const y = Math.round(by + (start + size) * k) + 0.5;
            ctx.moveTo(left, y);
            ctx.lineTo(left + px, y);
          });
        }
        ctx.stroke();
        ctx.restore();
      }
    }
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

  /**
   * The chunks the panes show (the frozen ones first: a few, with the
   * column names), then those a chunk around the scrolled region (drawn
   * ahead of a scroll).
   */
  const wanted = (chunks: SheetChunks) => {
    const shown: [number, number][] = [];
    for (const p of [...panes].reverse()) chunks.within(p.src, shown);
    const main = panes[0]?.src, u = chunks.unit;
    const around = main ? chunks.within({ x: main.x - u, y: main.y - u, w: main.w + 2 * u, h: main.h + 2 * u }) : [];
    return { list: [...shown, ...around], shown: shown.length };
  };
  let painting = false;
  const repaint = () => {
    if (painting) return;
    painting = true;
    requestAnimationFrame(() => {
      painting = false;
      if (gen === generation) paint();
    });
  };
  /** Ask for the chunks that are wanted and not drawn yet; paint each as it comes. */
  const fetchChunks = () => {
    const chunks = sheetChunks;
    if (gen !== generation || chunks?.view !== v) return;
    const { list, shown } = wanted(chunks);
    chunks.fetch(list, shown, (ms) => {
      if (gen !== generation) return;
      repaint();
      fetchChunks();
      setTiming(`sheet chunk rendered in ${ms.toFixed(0)} ms`);
    }, unlessStale(gen));
  };

  let pending = false;
  const draw = async () => {
    if (pending) return;
    pending = true;
    await new Promise((r) => requestAnimationFrame(r));
    pending = false;
    if (gen !== generation) return;
    const vw = stage.clientWidth, vh = stage.clientHeight;
    const sx = stage.scrollLeft, sy = stage.scrollTop;
    scrolled = { x: sx, y: sy };
    pinned.forEach(place);
    // the scrolled region starts where the frozen panes end
    const mx = fw + sx / zoom, my = fh + sy / zoom;
    const pw = Math.max(0, Math.min((vw - hw) / zoom - fw, cols.total - mx)), ph = Math.max(0, Math.min((vh - hh) / zoom - fh, rows.total - my));
    const next: Pane[] = [
      { src: { x: mx, y: my, w: pw, h: ph }, dx: hw + fw * zoom, dy: hh + fh * zoom },
      { src: { x: mx, y: 0, w: pw, h: fh }, dx: hw + fw * zoom, dy: hh },
      { src: { x: 0, y: my, w: fw, h: ph }, dx: hw, dy: hh + fh * zoom },
      { src: { x: 0, y: 0, w: fw, h: fh }, dx: hw, dy: hh },
    ].filter((p) => p.src.w > 0 && p.src.h > 0);
    panes = next;
    // the chunks of another view or scale (a zoom, or a screen of another resolution) are drawn anew
    const scale = zoom * dpr();
    if (sheetChunks?.view !== v || sheetChunks.scale !== scale) {
      sheetChunks?.drop();
      sheetChunks = new SheetChunks(v, scale, dpr());
    }
    paint();
    fetchChunks();
    updateText(next[0]?.src ?? { x: 0, y: 0, w: 1, h: 1 });
  };
  sheetViews.set(v, {
    redraw: () => { if (panes.length) paint(); },
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

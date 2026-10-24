// Shared by the administrator's pages (topology.js, inboxes.js,
// admin-jobs.js, deliveries.js, audit.js): the API with the session's
// cookie (the sign-in at /), toasts, small builders of elements, the
// names of the kinds and states in the page's language (i18n.js), and
// the nav's own tab.

import { _, locale, sentences } from "./i18n.js";

export { _, locale, sentences };

export const $ = (id) => document.getElementById(id);

// status shows a message for a while, as a toast: over a dialog open
// (the lightbox), in it, as the page's are under it.
export function status(text, bad) {
  if (!text) return;
  const t = document.createElement("div");
  t.role = "status";
  t.className = "alert shadow-lg " + (bad ? "alert-error" : "alert-success");
  t.textContent = text;
  const over = document.querySelector("dialog.modal[open]");
  let host = $("toasts");
  if (over) {
    host = over.querySelector(".toast");
    if (!host) {
      host = document.createElement("div");
      host.className = "toast toast-end z-50";
      over.append(host);
    }
  }
  host.append(t);
  setTimeout(() => t.remove(), bad ? 6000 : 3000);
}

// api calls the API with the session's cookie; a session that ended
// sends the page back to the sign-in.
export async function api(method, path, body, raw) {
  const res = await fetch(path, {
    method,
    headers: body && !raw ? { "Content-Type": "application/json" } : {},
    body: body ? (raw ? body : JSON.stringify(body)) : undefined,
  });
  if (res.status === 401) {
    location.href = "/";
    throw new Error(_("signed out"));
  }
  if (res.status === 204) return null;
  const out = await res.json();
  if (!res.ok) throw new Error(out.error || res.statusText);
  return out;
}

// el makes an element: its text, and its properties, but data-* and
// aria-*, which are attributes (a property of that name would be lost on
// querySelector).
export function el(tag, text, attrs = {}) {
  const e = document.createElement(tag);
  if (text !== undefined) e.textContent = text;
  for (const [k, v] of Object.entries(attrs)) {
    if (k.startsWith("data-") || k.startsWith("aria-")) e.setAttribute(k, v);
    else e[k] = v;
  }
  return e;
}

// changed tells the page that something changed: it loads again.
export function changed() {
  document.dispatchEvent(new CustomEvent("pm:changed"));
}

export function onChange(fn) {
  document.addEventListener("pm:changed", fn);
}

// button runs fn, disabled meanwhile, and tells the page something
// changed (unless fn says false: nothing was done); an error is a toast.
export function button(text, fn, className = "btn btn-sm") {
  const b = el("button", text, { type: "button", className });
  b.addEventListener("click", async () => {
    b.disabled = true;
    try {
      if ((await fn()) !== false) changed();
    } catch (e) {
      status(e.message, true);
    } finally {
      b.disabled = false;
    }
  });
  return b;
}

export function since(t) {
  if (!t) return "—";
  return new Date(t).toLocaleString(locale);
}

export function badge(text, cls = "badge-ghost") {
  return el("span", text, { className: "badge " + cls });
}

// KINDS are the kinds' names, as the palette says them (pages.yaml).
export const KINDS = { upload: _("Uploader"), cups: _("CUPS server (source)"), "cups-out": _("CUPS server (sink)"), queue: _("Queue"), class: _("Class"), filter: _("Filter"), assert: _("Assert"), branch: _("Branch"), template: _("Template"), convert: _("Convert"), sign: _("Sign"), inbox: _("Inbox"), archive: _("Archive"), display: _("Display"), printer: _("Printer"), sink: _("Sink"), download: _("Downloader"), render: _("Render"), notify: _("Notify"), smtp: _("Mail (SMTP)"), s3: "S3", "s3-trigger": _("S3 Trigger"), http: "HTTP (push)" };
export const KIND_BADGE = {
  upload: "badge-accent badge-soft",
  cups: "badge-error badge-soft",
  sink: "badge-success badge-soft",
  download: "badge-accent badge-soft",
  filter: "badge-secondary badge-soft",
  assert: "badge-secondary badge-soft",
  branch: "badge-secondary badge-soft",
  template: "badge-secondary badge-soft",
  convert: "badge-secondary badge-soft",
  sign: "badge-secondary badge-soft",
  queue: "badge-primary badge-soft",
  class: "badge-primary badge-soft",
  inbox: "badge-warning badge-soft",
  archive: "badge-neutral badge-soft",
  display: "badge-info badge-soft",
  printer: "badge-neutral badge-soft",
  "cups-out": "badge-neutral badge-soft",
};

export function kindBadge(kind) {
  return badge(KINDS[kind] || kind, KIND_BADGE[kind] || "badge-ghost");
}

// The kinds of payload that flow along the edges (kinds.go): a .bdf
// document, data for a template, any file as received, or a stocked
// file, which an inbox or an archive keeps and sends to a downloader
// alone. Ports, edges and the borders of the nodes are colored by them
// (canvas.js), as are the jobs' badges: blue, green, amber, cyan.
export const PAYLOADS = { bdf: _(".bdf document"), data: _("data (JSON, YAML, CSV, Excel, Parquet)"), any: _("any file as received"), stocked: _("stocked file (any format, kept by an inbox or an archive)") };
export const PAYLOAD_COLOR = { bdf: "#2563eb", data: "#16a34a", any: "#d97706", stocked: "#0891b2" };

// payloadBadge tells what a job carries: its kind and format.
export function payloadBadge(j) {
  const kind = j.kind || "bdf";
  const b = badge(kind === "bdf" ? ".bdf" : kind + " · " + (j.format || ""), "badge-sm badge-outline whitespace-nowrap");
  b.style.borderColor = PAYLOAD_COLOR[kind] || PAYLOAD_COLOR.any;
  b.style.color = PAYLOAD_COLOR[kind] || PAYLOAD_COLOR.any;
  b.title = j.format ? _("{payload}, received as {format}", { payload: PAYLOADS[kind] || kind, format: j.format }) : PAYLOADS[kind] || kind;
  return b;
}

// sender is who sent a job, as verified, and the name they gave.
export function sender(j) {
  const by = j.submitted_by || {};
  const who = by.name || by.id || by.kind || "";
  return j.claimed_user && j.claimed_user !== who ? _("{who} (as {claimed})", { who, claimed: j.claimed_user }) : who;
}

export function legsOf(j) {
  return (j.legs || [])
    .map((l) => l.node + ":" + l.state + (l.member ? "(" + l.member + ")" : "") + (l.error ? " ! " + l.error : ""))
    .join(" → ");
}

export const ENDED = ["completed", "canceled", "aborted"];

// STATES are the states of a job, as the page says them; the badge's
// title is the state itself, as the API and the manual write it.
export const STATES = { processing: _("processing"), held: _("held"), completed: _("completed"), canceled: _("canceled"), aborted: _("aborted") };

export function stateBadge(state) {
  const cls = state === "completed" ? "badge-success badge-soft" : state === "held" ? "badge-warning" : state === "processing" ? "badge-info" : "badge-ghost";
  return el("span", STATES[state] || state, { className: "badge whitespace-nowrap " + cls, title: state });
}

// pagesText says a count of pages.
export const pagesText = (n) => (n === 1 ? _("{n} page", { n }) : _("{n} pages", { n }));

// cancellable tells a job the administrator may still cancel: one not
// ended, or one whose slip is still open on a board.
export function cancellable(j) {
  return !ENDED.includes(j.state) || (j.legs || []).some((l) => l.kind === "display" && l.state === "shown");
}

// The lightbox shows a document large over the page: every page, one
// under another, with the download beside its title and a bar to zoom
// under it: − and +, the zoom (a part of the box's width), Fit (each page
// whole in the box, as it opens) and Width (as wide as the box). The keys
// +, − and 0, Ctrl or ⌘ with the wheel, and a pinch zoom too, around
// where they point; a page larger than the box is dragged to move it. A
// page is drawn as wide as it is shown, for the screen's pixels (?w=, in
// steps up to the server's largest), the first two at once and the others
// as they are scrolled to; zoomed in, it is drawn again, sharper. What may
// be done with the document (an inbox's Release and Discard) is at the
// end of the bar. Escape, ✕ or a click beside it closes it.
let lightbox;

// ZOOMS are the steps of − and +, as a part of the box's width.
const ZOOMS = [0.25, 0.33, 0.5, 0.67, 0.75, 1, 1.25, 1.5, 2, 3, 4];
const MIN_ZOOM = 0.1;
const MAX_ZOOM = 4;
// WIDTHS are the widths a page is drawn at: the first as wide as it is
// shown (MaxPageWidth, docflow/submit.go, is the last).
const WIDTHS = [400, 800, 1200, 1600, 2000, 2400];

addEventListener("resize", () => lightbox?.open && lightbox.relayout?.());

// preview opens the lightbox on a document: its title and the line under
// it, its count of pages, page(n, width) the URL of a page, href its
// download, and actions(close) the buttons of what may be done with it;
// one with no pages (not a .bdf document) says its format.
export function preview(doc) {
  if (!lightbox) {
    lightbox = el("dialog", undefined, { className: "modal" });
    document.body.append(lightbox);
  }
  const close = () => lightbox.close();
  const box = el("div", undefined, { className: "modal-box flex max-h-[94vh] w-11/12 max-w-5xl flex-col gap-0 p-0" + (doc.pages ? " h-[94vh]" : "") });
  const head = el("div", undefined, { className: "flex items-center gap-2 border-b border-base-300 px-4 py-3" });
  const titles = el("div", undefined, { className: "min-w-0 flex-1" });
  titles.append(el("h3", doc.title, { className: "truncate text-lg font-semibold", title: doc.title }), el("p", doc.sub, { className: "text-xs text-base-content/60" }));
  const x = el("button", "✕", { type: "button", className: "btn btn-square btn-ghost btn-sm", title: _("Close"), "aria-label": _("Close") });
  x.addEventListener("click", close);
  head.append(titles, el("a", _("Download"), { className: "btn btn-primary btn-sm", href: doc.href }), x);
  const bar = el("div", undefined, { className: "flex flex-wrap items-center gap-2 border-b border-base-300 px-4 py-2" });
  const pages = el("div", undefined, { className: "min-h-0 flex-1 overflow-auto bg-base-200 p-4 outline-none", tabIndex: 0 });
  // the browser's own pinch would zoom the whole page
  pages.style.touchAction = "pan-x pan-y";
  const imgs = [];
  for (let n = 1; n <= doc.pages; n++) {
    const fig = el("figure", undefined, { className: "mx-auto mb-4 flex w-fit flex-col items-center gap-1" });
    const img = el("img", undefined, { alt: _("page {n}", { n }), className: "max-w-none rounded bg-white shadow-md", draggable: false, loading: n <= 2 ? "eager" : "lazy" });
    fig.append(img, el("figcaption", `${n} / ${doc.pages}`, { className: "text-xs text-base-content/60" }));
    pages.append(fig);
    imgs.push(img);
  }
  if (!doc.pages) pages.append(el("p", doc.format ? _("{format}: no pages to show. Download it to open it.", { format: doc.format }) : _("No pages to show: download it to open it."), { className: "py-8 text-center text-sm text-base-content/60" }));

  // zoom: null, each page whole in the box; "width", as wide as it; else
  // a part of its width. ratio is a page's height as a part of its width,
  // of the last page drawn: a page not drawn yet holds that place.
  let zoom = null;
  let ratio = 1.3;
  let sharpening = 0;
  const room = () => ({ w: Math.max(80, pages.clientWidth - 32), h: Math.max(80, pages.clientHeight - 32 - 20) });
  const ratioOf = (img) => (img.naturalWidth ? img.naturalHeight / img.naturalWidth : ratio);
  const widthOf = (img) => {
    const { w, h } = room();
    return zoom === null ? Math.min(w, h / ratioOf(img)) : zoom === "width" ? w : zoom * w;
  };
  const size = (img) => {
    const w = widthOf(img);
    img.style.width = `${w}px`;
    img.style.height = `${w * ratioOf(img)}px`;
  };
  // anchor is the page at a point of the box (its middle by default), or
  // the nearest, and where in it; keep scrolls it back under the point
  function anchor(px, py) {
    const b = pages.getBoundingClientRect();
    const x = b.left + (px ?? b.width / 2);
    const y = b.top + (py ?? b.height / 2);
    let best = null;
    let far = Infinity;
    for (const img of imgs) {
      const r = img.getBoundingClientRect();
      const d = y < r.top ? r.top - y : y > r.bottom ? y - r.bottom : 0;
      if (d < far) [best, far] = [{ img, r }, d];
    }
    return best && best.r.width ? { img: best.img, fx: (x - best.r.left) / best.r.width, fy: (y - best.r.top) / best.r.height, x, y } : null;
  }
  function keep(a) {
    if (!a) return;
    const r = a.img.getBoundingClientRect();
    pages.scrollLeft += r.left + a.fx * r.width - a.x;
    pages.scrollTop += r.top + a.fy * r.height - a.y;
  }
  // shown is the zoom of the page in view, as a part of the box's width
  const shown = () => {
    const a = anchor();
    return a ? widthOf(a.img) / room().w : 1;
  };
  function layout() {
    for (const img of imgs) size(img);
    level.textContent = `${Math.round(shown() * 100)}%`;
    fit.classList.toggle("btn-active", zoom === null);
    wide.classList.toggle("btn-active", zoom === "width");
    pages.style.cursor = pages.scrollWidth > pages.clientWidth || pages.scrollHeight > pages.clientHeight ? "grab" : "";
    clearTimeout(sharpening);
    sharpening = setTimeout(sharpen, 150);
  }
  // sharpen draws each page again where it is shown wider than it was
  // drawn; one drawn already is swapped once the sharper one is there
  function sharpen() {
    const dpr = Math.min(devicePixelRatio || 1, 2);
    for (const [i, img] of imgs.entries()) {
      const need = parseFloat(img.style.width) * dpr;
      const want = WIDTHS.find((w) => w >= need) || WIDTHS[WIDTHS.length - 1];
      if (want <= Number(img.dataset.w || 0)) continue;
      img.dataset.w = String(want);
      const url = doc.page(i + 1, want);
      if (!img.naturalWidth) {
        img.src = url; // not drawn yet: drawn as it is scrolled to
      } else {
        const sharper = new Image();
        sharper.onload = () => (img.src = url);
        sharper.src = url;
      }
    }
  }
  function setZoom(z, px, py) {
    const a = anchor(px, py);
    zoom = z === null || z === "width" ? z : Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, z));
    layout();
    keep(a);
  }
  const step = (up) => {
    const z = shown();
    setZoom(up ? ZOOMS.find((s) => s > z + 0.005) || MAX_ZOOM : [...ZOOMS].reverse().find((s) => s < z - 0.005) || MIN_ZOOM);
  };
  for (const img of imgs) {
    // a page drawn takes its own height: what is in view stays there
    // (the top, before anything was scrolled)
    img.addEventListener("load", () => {
      const a = pages.scrollTop > 0 ? anchor() : null;
      ratio = ratioOf(img);
      for (const o of imgs) if (o === img || !o.naturalWidth) size(o);
      keep(a);
    });
  }

  const zoomer = el("div", undefined, { className: "join" });
  const out = el("button", "−", { type: "button", className: "btn btn-sm join-item", title: _("Zoom out (−)"), "aria-label": _("Zoom out") });
  const level = el("span", "", { className: "btn btn-sm join-item pointer-events-none min-w-16 tabular-nums", "aria-live": "polite" });
  const into = el("button", "+", { type: "button", className: "btn btn-sm join-item", title: _("Zoom in (+)"), "aria-label": _("Zoom in") });
  out.addEventListener("click", () => step(false));
  into.addEventListener("click", () => step(true));
  zoomer.append(out, level, into);
  const fit = el("button", _("Fit"), { type: "button", className: "btn btn-sm", title: _("Each page whole in the box (0)") });
  const wide = el("button", _("Width"), { type: "button", className: "btn btn-sm", title: _("As wide as the box") });
  fit.addEventListener("click", () => setZoom(null));
  wide.addEventListener("click", () => setZoom("width"));
  if (doc.pages) bar.append(zoomer, fit, wide);
  const acts = doc.actions ? doc.actions(close) : [];
  if (acts.length) {
    const end = el("div", undefined, { className: "ms-auto flex flex-wrap items-center gap-2" });
    end.append(...acts);
    bar.append(end);
  }

  // Ctrl or ⌘ with the wheel (a trackpad's pinch, in most browsers) zooms
  // where it points; Safari's pinch is a gesture, a phone's two fingers
  let pinch = null;
  pages.addEventListener("wheel", (ev) => {
    if (!(ev.ctrlKey || ev.metaKey) || !doc.pages || pinch) return;
    ev.preventDefault();
    const b = pages.getBoundingClientRect();
    const dy = ev.deltaMode === 1 ? ev.deltaY * 33 : ev.deltaY;
    setZoom(shown() * Math.min(1.25, Math.max(0.8, Math.exp(-dy * 0.01))), ev.clientX - b.left, ev.clientY - b.top);
  }, { passive: false });
  if ("ongesturechange" in window) {
    pages.addEventListener("gesturestart", (ev) => {
      ev.preventDefault();
      pinch = { z: shown() };
    });
    pages.addEventListener("gesturechange", (ev) => {
      if (!pinch) return;
      ev.preventDefault();
      const b = pages.getBoundingClientRect();
      setZoom(pinch.z * ev.scale, ev.clientX - b.left, ev.clientY - b.top);
    });
    pages.addEventListener("gestureend", () => (pinch = null));
  } else {
    const apart = (t) => Math.hypot(t[0].clientX - t[1].clientX, t[0].clientY - t[1].clientY);
    pages.addEventListener("touchstart", (ev) => {
      if (ev.touches.length === 2) pinch = { d: apart(ev.touches), z: shown() };
    }, { passive: true });
    pages.addEventListener("touchmove", (ev) => {
      if (!pinch || ev.touches.length !== 2) return;
      ev.preventDefault();
      const b = pages.getBoundingClientRect();
      const [t0, t1] = ev.touches;
      setZoom((pinch.z * apart(ev.touches)) / pinch.d, (t0.clientX + t1.clientX) / 2 - b.left, (t0.clientY + t1.clientY) / 2 - b.top);
    }, { passive: false });
    pages.addEventListener("touchend", (ev) => {
      if (ev.touches.length < 2) pinch = null;
    });
  }
  // a mouse drags what is larger than the box (not on its scroll bars)
  let drag = null;
  pages.addEventListener("pointerdown", (ev) => {
    const b = pages.getBoundingClientRect();
    if (ev.pointerType !== "mouse" || ev.button !== 0 || !pages.style.cursor || ev.clientX - b.left >= pages.clientWidth || ev.clientY - b.top >= pages.clientHeight) return;
    drag = { x: ev.clientX, y: ev.clientY, left: pages.scrollLeft, top: pages.scrollTop };
    pages.setPointerCapture(ev.pointerId);
    pages.style.cursor = "grabbing";
    ev.preventDefault();
    pages.focus({ preventScroll: true });
  });
  pages.addEventListener("pointermove", (ev) => {
    if (!drag) return;
    pages.scrollLeft = drag.left - (ev.clientX - drag.x);
    pages.scrollTop = drag.top - (ev.clientY - drag.y);
  });
  const drop = () => {
    if (!drag) return;
    drag = null;
    pages.style.cursor = "grab";
  };
  pages.addEventListener("pointerup", drop);
  pages.addEventListener("pointercancel", drop);
  // the keys, but those of a field and the browser's own (with a modifier)
  lightbox.onkeydown = (ev) => {
    if (!doc.pages || ev.ctrlKey || ev.metaKey || ev.altKey || ev.target.closest?.("input, textarea, select")) return;
    if (ev.key === "+" || ev.key === "=") step(true);
    else if (ev.key === "-" || ev.key === "_") step(false);
    else if (ev.key === "0") setZoom(null);
    else return;
    ev.preventDefault();
  };
  lightbox.relayout = () => {
    const a = anchor();
    layout();
    keep(a);
  };

  const beside = el("form", undefined, { method: "dialog", className: "modal-backdrop" });
  beside.append(el("button", _("Close")));
  box.append(head);
  if (bar.childElementCount) box.append(bar);
  box.append(pages);
  lightbox.replaceChildren(box, beside);
  lightbox.showModal();
  if (doc.pages) {
    layout();
    clearTimeout(sharpening);
    sharpen();
  }
  x.focus();
}

// thumb is the first page of a document, which opens it in the
// lightbox.
function thumb(doc, width, className) {
  const b = el("button", undefined, { type: "button", className: "cursor-pointer " + className, title: _("Show the pages, with the download") });
  b.append(el("img", undefined, { src: doc.page(1, width), alt: _("page 1 of {doc}", { doc: doc.title }), className: "aspect-[1/1.3] w-full rounded bg-white object-cover object-top shadow-sm", loading: "lazy" }));
  b.addEventListener("click", () => preview(doc));
  return b;
}

// The sections of what is at a node now, on the topology (the sidebar of
// a node, looking, and its Usage): what a sink keeps, what a display's
// board shows, what an inbox holds (worked there: Release and Discard),
// what an archive wrote lately, and what a downloader lists. Each is made
// at once and read in the background; its refresh reads it again in
// place (as the activity changes), keeping what it shows until the new
// is there.

// liveSection is a section of a title, read by fill (the elements to
// show, or null to keep what it shows); top stays above them.
function liveSection(title, fill, top) {
  const box = el("div", undefined, { className: "mt-3 flex flex-col gap-2 border-t border-base-300 pt-3" });
  const body = el("div", undefined, { className: "flex flex-col gap-2" });
  box.append(el("h4", title, { className: "font-semibold" }));
  if (top) box.append(top);
  box.append(body);
  let reading = null;
  let again = false;
  box.refresh = () => {
    if (reading) {
      again = true;
      return reading;
    }
    reading = fill()
      .then((parts) => parts && body.replaceChildren(...parts))
      .catch((e) => status(e.message, true))
      .finally(() => {
        reading = null;
        if (again) {
          again = false;
          box.refresh();
        }
      });
    return reading;
  };
  box.refresh();
  return box;
}

const note = (text) => el("p", text, { className: "text-xs text-base-content/60" });

// fresh tells whether what a section read differs from what it read
// last (else it keeps what it shows: no flicker, nothing redrawn).
function fresh() {
  let last = null;
  return (v) => {
    const s = JSON.stringify(v);
    if (s === last) return false;
    last = s;
    return true;
  };
}

// card is a document in a section's grid: its first page (a click shows
// them all, with the download), its title, the line under it and badges.
function card(doc, badges = []) {
  const b = el("button", undefined, { type: "button", className: "flex min-w-0 cursor-pointer flex-col gap-1 rounded-box bg-base-200 p-2 text-left text-xs hover:bg-base-300", title: doc.pages ? _("Show the pages, with the download") : _("Show it, with the download") });
  if (doc.pages) b.append(el("img", undefined, { src: doc.page(1, 240), alt: _("page 1 of {doc}", { doc: doc.title }), className: "aspect-[1/1.3] w-full rounded bg-white object-cover object-top shadow-sm", loading: "lazy" }));
  else b.append(el("span", doc.format || "?", { className: "flex aspect-[1/1.3] w-full items-center justify-center rounded bg-base-100 font-mono text-base-content/60" }));
  b.append(el("span", doc.title, { className: "truncate font-semibold" }), el("span", doc.sub, { className: "text-base-content/60" }));
  if (badges.length) {
    const row = el("span", undefined, { className: "flex flex-wrap gap-1" });
    row.append(...badges);
    b.append(row);
  }
  b.addEventListener("click", () => preview(doc));
  return b;
}

function grid(cards) {
  const g = el("div", undefined, { className: "grid grid-cols-2 gap-2" });
  g.append(...cards);
  return g;
}

// keptSection shows the jobs a sink keeps that the viewer sees, newest
// first, with the first page of each: a click shows them all.
export function keptSection(name) {
  const isNew = fresh();
  return liveSection(_("Kept here"), async () => {
    const jobs = await api("GET", `/api/nodes/${encodeURIComponent(name)}/latest`);
    if (!isNew(jobs)) return null;
    if (!jobs.length) return [note(_("Nothing yet."))];
    return [grid(jobs.map((j) => card({ title: `#${j.number} ${j.name || ""}`, sub: `${pagesText(j.pages)} · ${since(j.received)}`, pages: j.pages || 0, format: j.format, page: (n, w) => `/api/jobs/${j.id}/pages/${n}.png?w=${w}`, href: `/api/jobs/${j.id}/document` })))];
  });
}

// SLIP_STATES are the states of a slip, as the board says them.
const SLIP_STATES = { new: _("New"), opened: _("Opened"), done: _("Done"), canceled: _("Canceled") };
const SLIP_BADGE = { new: "badge-warning", opened: "badge-info badge-soft", done: "badge-success badge-soft", canceled: "badge-ghost" };
// ENDED_SHOWN is how many of the slips done or canceled are shown.
const ENDED_SHOWN = 6;

// boardSection shows what a display's board shows that the viewer sees
// (every slip, to the role topology): the slips open, oldest first, as
// the devices list them, and the latest done or canceled, with the first
// page of each: a click shows them all.
export function boardSection(name) {
  const isNew = fresh();
  return liveSection(_("On the board"), async () => {
    const b = await api("GET", `/api/nodes/${encodeURIComponent(name)}/board`);
    if (!isNew(b)) return null;
    if (!b.open.length && !b.ended.length) return [note(_("Nothing on the board that you may see."))];
    const slip = (s) => {
      const at = `/api/slips/${encodeURIComponent(s.id)}`;
      const badges = [el("span", SLIP_STATES[s.state] || s.state, { className: "badge badge-xs whitespace-nowrap " + (SLIP_BADGE[s.state] || "badge-ghost") })];
      if ((s.printed || []).length) badges.push(el("span", _("Printed on paper"), { className: "badge badge-xs badge-outline whitespace-nowrap" }));
      return card({ title: `#${s.job} ${s.name || ""}`, sub: `${pagesText(s.pages)} · ${since(s.received)}`, pages: s.pages || 0, page: (n, w) => `${at}/pages/${n}.png?w=${w}`, href: `${at}/document` }, badges);
    };
    const out = [el("p", _("Open ({n})", { n: b.open.length }), { className: "text-xs font-semibold" })];
    out.push(b.open.length ? grid(b.open.map(slip)) : note(_("None open.")));
    if (b.ended.length) {
      out.push(el("p", _("Done and canceled ({n})", { n: b.ended.length }), { className: "text-xs font-semibold" }), grid(b.ended.slice(0, ENDED_SHOWN).map(slip)));
      if (b.ended.length > ENDED_SHOWN) out.push(note(_("The latest {n} of them.", { n: ENDED_SHOWN })));
    }
    return out;
  });
}

let me = null; // /api/me, read once: whether the viewer has the role topology

// workerField is the name releases and discards are recorded under, as
// the Inboxes page keeps it (bdf-pm-worker).
function workerField() {
  const label = el("label", undefined, { className: "input input-sm w-full" });
  const input = el("input", undefined, { autocomplete: "off", placeholder: _("name (for the record)") });
  try {
    input.value = localStorage.getItem("bdf-pm-worker") || "";
  } catch {}
  input.addEventListener("change", () => {
    try {
      localStorage.setItem("bdf-pm-worker", input.value);
    } catch {}
  });
  label.append(el("span", _("Worker"), { className: "label" }), input);
  label.input = input;
  return label;
}

// ITEMS_ENDED is how many of the items released or discarded are listed.
const ITEMS_ENDED = 5;

// heldSection shows what an inbox holds, oldest first, to those who work
// it (its API, as the Inboxes page reads it): each with its first page
// (a click shows them all, with Release and Discard beside the download),
// the nodes it goes to (the inbox's, every one chosen at first; any typed
// by the role topology when it names none), Release and Discard, under
// the worker's name; and the latest released or discarded. What was
// chosen of an item stays while the section is read again. Anyone else is
// told who works it.
export function heldSection(name) {
  const base = `/api/inbox/${encodeURIComponent(name)}`;
  const worker = workerField();
  const chosen = new Map(); // by job: {off} the nodes not chosen, or {typed}
  let etag = "";
  return liveSection(_("Held here"), async () => {
    const res = await fetch(base, { headers: etag ? { "If-None-Match": etag } : {} });
    if (res.status === 304) return null;
    if (res.status === 401) {
      location.href = "/";
      return null;
    }
    if (res.status === 403) {
      worker.hidden = true;
      return [note(_("You do not work this inbox: the members of its group, its owner and the role topology release what it holds."))];
    }
    const v = await res.json();
    if (!res.ok) throw new Error(v.error || res.statusText);
    etag = res.headers.get("ETag") || "";
    worker.hidden = false;
    me ||= await api("GET", "/api/me");
    const out = [];
    if (!v.held.length) out.push(note(_("Nothing held")));
    const list = el("ul", undefined, { className: "flex flex-col gap-2" });
    for (const it of v.held) {
      const at = `${base}/${it.job}`;
      const doc = { title: `#${it.job} ${it.name || ""}`, sub: `${it.pages ? pagesText(it.pages) + " · " : ""}${sender(it)} · ${since(it.received)}`, pages: it.pages || 0, page: (n, w) => `${at}/pages/${n}.png?w=${w}`, href: `${at}/document` };
      const kept = chosen.get(it.job) || {};
      const pick = el("div", undefined, { className: "flex flex-wrap items-center gap-x-3 gap-y-1" });
      let to = null;
      if (v.to.length) {
        pick.append(el("span", "→", { className: "text-base-content/60" }));
        const boxes = v.to.map((node) => {
          const cb = el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", checked: !(kept.off || []).includes(node), value: node });
          cb.addEventListener("change", () => {
            chosen.set(it.job, { off: boxes.filter((b) => !b.checked).map((b) => b.value) });
            sync();
          });
          const l = el("label", undefined, { className: "inline-flex items-center gap-1 font-mono" });
          l.append(cb, node);
          pick.append(l);
          return cb;
        });
        to = () => boxes.filter((b) => b.checked).map((b) => b.value);
      } else if (me.topology) {
        const free = el("input", undefined, { className: "input input-xs w-full", placeholder: _("nodes (comma-separated)"), value: kept.typed || "" });
        free.addEventListener("input", () => {
          chosen.set(it.job, { typed: free.value });
          sync();
        });
        pick.append(free);
        to = () => free.value.split(",").map((x) => x.trim()).filter(Boolean);
      } else {
        pick.append(note(_("This inbox names no nodes: the role topology releases what it holds.")));
      }
      const release = async () => {
        const nodes = to();
        await api("POST", `${at}/release`, { to: nodes, worker: worker.input.value });
        status(_("#{n} released to {nodes}", { n: it.job, nodes: nodes.join(", ") }));
      };
      const discard = async () => {
        if (!confirm(_("Discard #{n} {name}? It goes on to no node.", { n: it.job, name: it.name || "" }))) return false;
        await api("POST", `${at}/discard`, { worker: worker.input.value });
        status(_("#{n} discarded", { n: it.job }));
      };
      const releaseButton = (then) => {
        const b = button(_("Release"), async () => {
          if ((await release()) === false) return false;
          then?.();
        }, "btn btn-xs btn-primary");
        b.disabled = !to || !to().length;
        return b;
      };
      const discardButton = (then) => button(_("Discard"), async () => {
        if ((await discard()) === false) return false;
        then?.();
      }, "btn btn-xs btn-ghost text-error");
      const rel = releaseButton();
      const sync = () => (rel.disabled = !to || !to().length);
      doc.actions = (close) => [releaseButton(close), discardButton(close)];
      const li = el("li", undefined, { className: "flex items-start gap-2 rounded-box bg-base-200 p-2 text-xs" });
      if (it.pages) li.append(thumb(doc, 120, "w-14 shrink-0"));
      else {
        const open = el("button", _("Document"), { type: "button", className: "link w-14 shrink-0 text-left" });
        open.addEventListener("click", () => preview(doc));
        li.append(open);
      }
      const text = el("div", undefined, { className: "flex min-w-0 flex-1 flex-col gap-1" });
      const title = el("button", doc.title, { type: "button", className: "link truncate text-left font-semibold", title: _("Show the pages, with the download") });
      title.addEventListener("click", () => preview(doc));
      const acts = el("div", undefined, { className: "flex flex-wrap gap-1" });
      acts.append(rel, discardButton());
      text.append(title, el("span", doc.sub, { className: "text-base-content/60" }), pick, acts);
      li.append(text);
      list.append(li);
    }
    for (const job of chosen.keys()) if (!v.held.some((it) => it.job === job)) chosen.delete(job);
    if (v.held.length) out.push(list);
    if (v.ended.length) {
      const details = el("details", undefined, { className: "text-xs" });
      details.append(el("summary", _("Released or discarded: {n}", { n: v.ended.length }), { className: "cursor-pointer text-base-content/60" }));
      const ended = el("ul", undefined, { className: "mt-1 flex flex-col gap-1" });
      for (const it of v.ended.slice(0, ITEMS_ENDED)) {
        const what = it.state === "released" ? `→ ${it.released.to.join(", ")} (${it.released.worker || ""} ${since(it.released.at)})` : _("discarded ({who} {when})", { who: it.discarded.worker || "", when: since(it.discarded.at) });
        ended.append(el("li", `#${it.job} ${it.name || ""} ${what}`));
      }
      details.append(ended);
      out.push(details);
    }
    return out;
  }, worker);
}

// archivedSection shows the latest documents an archive wrote that the
// viewer sees, newest first, with the first page of each: a click shows
// them all, drawn from the archive's store, with the download (which is
// in the audit log).
export function archivedSection(name) {
  const base = `/api/nodes/${encodeURIComponent(name)}/archived`;
  const isNew = fresh();
  return liveSection(_("Archived here"), async () => {
    const es = await api("GET", base);
    if (!isNew(es)) return null;
    if (!es.length) return [note(_("Nothing you may see yet."))];
    const cards = es.map((e) => {
      const at = `${base}/${encodeURIComponent(e.key)}`;
      const pages = e.kind === "bdf" ? e.pages || 0 : 0;
      return card({ title: `#${e.job} ${e.name || ""}`, sub: `${pages ? pagesText(pages) + " · " : ""}${since(e.archived)}`, pages, format: e.format || e.kind, page: (n, w) => `${at}/pages/${n}.png?w=${w}`, href: `${at}/document` });
    });
    return [grid(cards), note(_("The latest it wrote; its store keeps every one."))];
  });
}

// listedSection shows the entries of a downloader the viewer may read:
// what an archive before it wrote (or an inbox holds), with the owner
// each is listed for; a click shows one, with its download.
export function listedSection(name) {
  const at = `/api/download/${encodeURIComponent(name)}`;
  const isNew = fresh();
  return liveSection(_("Listed here"), async () => {
    const entries = await api("GET", `${at}?limit=20`);
    if (!isNew(entries)) return null;
    if (!entries.length) return [note(_("Nothing you may read yet."))];
    const list = el("ul", undefined, { className: "flex flex-col gap-1 text-xs" });
    for (const e of entries) {
      const pages = e.kind === "bdf" ? e.pages || 0 : 0;
      const doc = { title: `#${e.job} ${e.name || ""}`, sub: `${pages ? pagesText(pages) + " · " : ""}${e.format || e.kind} · ${since(e.received)}`, pages, format: e.format || e.kind, page: (n, w) => `${at}/${e.job}/pages/${n}.png?w=${w}`, href: `${at}/${e.job}/document` };
      const li = el("li", undefined, { className: "flex items-start gap-2 rounded-box bg-base-200 p-2" });
      if (pages) li.append(thumb(doc, 96, "w-10 shrink-0"));
      const text = el("div", undefined, { className: "flex min-w-0 flex-1 flex-wrap items-center gap-x-2" });
      const open = el("button", doc.title, { type: "button", className: "link text-left font-semibold", title: pages ? _("Show the pages, with the download") : _("Show it, with the download") });
      open.addEventListener("click", () => preview(doc));
      text.append(open);
      if (e.owner) text.append(el("span", e.owner, { className: "badge badge-ghost badge-xs" }));
      text.append(el("span", doc.sub, { className: "text-base-content/60" }));
      for (const [k, v] of Object.entries(e.fields || {})) text.append(el("span", `${k}: ${v}`, { className: "font-mono text-base-content/60" }));
      li.append(text);
      list.append(li);
    }
    return [list];
  });
}

// the nav's own tab
for (const a of document.querySelectorAll("#admin-nav a")) {
  if (a.getAttribute("href") === location.pathname) a.classList.add("tab-active");
}

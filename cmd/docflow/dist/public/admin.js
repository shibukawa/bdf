// Shared by the administrator's pages (topology.js, inboxes.js,
// admin-jobs.js, deliveries.js, audit.js): the API with the session's
// cookie (the sign-in at /), toasts, small builders of elements, the
// names of the kinds and states in the page's language (i18n.js), and
// the nav's own tab.

import { _, locale, sentences } from "./i18n.js";

export { _, locale, sentences };

export const $ = (id) => document.getElementById(id);

// status shows a message for a while, as a toast.
export function status(text, bad) {
  if (!text) return;
  const t = document.createElement("div");
  t.role = "status";
  t.className = "alert shadow-lg " + (bad ? "alert-error" : "alert-success");
  t.textContent = text;
  $("toasts").append(t);
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

// button runs fn and tells the page something changed; an error is a
// toast.
export function button(text, fn, className = "btn btn-sm") {
  const b = el("button", text, { type: "button", className });
  b.addEventListener("click", async () => {
    try {
      await fn();
      changed();
    } catch (e) {
      status(e.message, true);
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
export const KINDS = { endpoint: _("Endpoint"), cups: _("CUPS server (source)"), "cups-out": _("CUPS server (sink)"), queue: _("Queue"), class: _("Class"), filter: _("Filter"), assert: _("Assert"), branch: _("Branch"), template: _("Template"), convert: _("Convert"), sign: _("Sign"), inbox: _("Inbox"), archive: _("Archive"), display: _("Display"), printer: _("Printer"), sink: _("Sink"), export: "HTTP (pull)", render: _("Render"), notify: _("Notify"), smtp: _("Mail (SMTP)"), s3: "S3", "s3-trigger": _("S3 Trigger"), http: "HTTP (push)" };
export const KIND_BADGE = {
  endpoint: "badge-accent badge-soft",
  cups: "badge-error badge-soft",
  sink: "badge-success badge-soft",
  export: "badge-accent badge-soft",
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
// file, which an inbox or an archive keeps and sends to HTTP (pull)
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
// under another, each as tall as the screen at most, with the download
// beside its title. A page is drawn as wide as the box shows it (?w=,
// for the screen's pixels), the first two at once and the others as they
// are scrolled to. Escape, ✕ or a click beside it closes it.
let lightbox;

// preview opens the lightbox on a document: its title and the line under
// it, its count of pages, page(n, width) the URL of a page, href its
// download; one with no pages (not a .bdf document) says its format.
export function preview(doc) {
  if (!lightbox) {
    lightbox = el("dialog", undefined, { className: "modal" });
    document.body.append(lightbox);
  }
  const box = el("div", undefined, { className: "modal-box flex max-h-[94vh] w-11/12 max-w-5xl flex-col gap-0 p-0" });
  const head = el("div", undefined, { className: "flex items-center gap-2 border-b border-base-300 px-4 py-3" });
  const titles = el("div", undefined, { className: "min-w-0 flex-1" });
  titles.append(el("h3", doc.title, { className: "truncate text-lg font-semibold", title: doc.title }), el("p", doc.sub, { className: "text-xs text-base-content/60" }));
  const close = el("button", "✕", { type: "button", className: "btn btn-square btn-ghost btn-sm", title: _("Close"), "aria-label": _("Close") });
  close.addEventListener("click", () => lightbox.close());
  head.append(titles, el("a", _("Download"), { className: "btn btn-primary btn-sm", href: doc.href }), close);
  const pages = el("div", undefined, { className: "flex flex-col items-center gap-4 overflow-y-auto bg-base-200 p-4" });
  const width = Math.round(Math.min(innerWidth * 0.92, 1024) * Math.min(devicePixelRatio || 1, 2));
  for (let n = 1; n <= doc.pages; n++) {
    const fig = el("figure", undefined, { className: "flex max-w-full flex-col items-center gap-1" });
    // until it is drawn, a page holds the place of a portrait one
    const img = el("img", undefined, { src: doc.page(n, width), alt: _("page {n}", { n }), className: "aspect-[1/1.3] h-[80vh] max-w-full rounded bg-white shadow-md", loading: n <= 2 ? "eager" : "lazy" });
    img.addEventListener("load", () => {
      img.classList.replace("aspect-[1/1.3]", "w-auto");
      img.classList.replace("h-[80vh]", "max-h-[80vh]");
    }, { once: true });
    fig.append(img, el("figcaption", `${n} / ${doc.pages}`, { className: "text-xs text-base-content/60" }));
    pages.append(fig);
  }
  if (!doc.pages) pages.append(el("p", _("{format}: no pages to show. Download it to open it.", { format: doc.format || "?" }), { className: "py-8 text-sm text-base-content/60" }));
  const beside = el("form", undefined, { method: "dialog", className: "modal-backdrop" });
  beside.append(el("button", _("Close")));
  box.append(head, pages);
  lightbox.replaceChildren(box, beside);
  lightbox.showModal();
  close.focus();
}

// thumb is the first page of a document, which opens it in the
// lightbox.
function thumb(doc, width, className) {
  const b = el("button", undefined, { type: "button", className: "cursor-pointer " + className, title: _("Show the pages, with the download") });
  b.append(el("img", undefined, { src: doc.page(1, width), alt: _("page 1 of {doc}", { doc: doc.title }), className: "aspect-[1/1.3] w-full rounded bg-white object-cover object-top shadow-sm", loading: "lazy" }));
  b.addEventListener("click", () => preview(doc));
  return b;
}

// keptSection lists the jobs a sink keeps that the viewer sees, newest
// first, with the first page of each: a click shows them all.
export async function keptSection(name) {
  const box = el("div", undefined, { className: "mt-3 flex flex-col gap-2 border-t border-base-300 pt-3" });
  box.append(el("h4", _("Kept here"), { className: "font-semibold" }));
  const jobs = await api("GET", `/api/nodes/${encodeURIComponent(name)}/latest`);
  if (!jobs.length) box.append(el("p", _("Nothing yet."), { className: "text-xs text-base-content/60" }));
  const grid = el("div", undefined, { className: "grid grid-cols-2 gap-2" });
  for (const j of jobs) {
    const doc = { title: `#${j.number} ${j.name || ""}`, sub: `${pagesText(j.pages)} · ${since(j.received)}`, pages: j.pages || 0, page: (n, w) => `/api/jobs/${j.id}/pages/${n}.png?w=${w}`, href: `/api/jobs/${j.id}/document` };
    const card = thumb(doc, 240, "flex flex-col gap-1 rounded-box bg-base-200 p-2 text-left text-xs hover:bg-base-300");
    card.append(el("span", doc.title, { className: "truncate font-semibold" }), el("span", doc.sub, { className: "text-base-content/60" }));
    grid.append(card);
  }
  box.append(grid);
  return box;
}

// listedSection shows the entries of an export the viewer may read:
// what an archive before it wrote (or an inbox holds), with the owner
// each is listed for; a click shows one, with its download.
export async function listedSection(name) {
  const box = el("div", undefined, { className: "mt-3 flex flex-col gap-2 border-t border-base-300 pt-3" });
  box.append(el("h4", _("Listed here"), { className: "font-semibold" }));
  const at = `/out/${encodeURIComponent(name)}`;
  const entries = await api("GET", `${at}?limit=20`);
  if (!entries.length) box.append(el("p", _("Nothing you may read yet."), { className: "text-xs text-base-content/60" }));
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
  box.append(list);
  return box;
}

// the nav's own tab
for (const a of document.querySelectorAll("#admin-nav a")) {
  if (a.getAttribute("href") === location.pathname) a.classList.add("tab-active");
}

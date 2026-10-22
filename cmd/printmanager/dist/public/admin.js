// Shared by the administrator's pages (topology.js, inboxes.js,
// admin-jobs.js, deliveries.js, audit.js): the API with the session's
// cookie (the sign-in at /), toasts, small builders of elements, and the
// nav's own tab.

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
    throw new Error("signed out");
  }
  if (res.status === 204) return null;
  const out = await res.json();
  if (!res.ok) throw new Error(out.error || res.statusText);
  return out;
}

export function el(tag, text, attrs = {}) {
  const e = document.createElement(tag);
  if (text !== undefined) e.textContent = text;
  Object.assign(e, attrs);
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
  return new Date(t).toLocaleString();
}

export function badge(text, cls = "badge-ghost") {
  return el("span", text, { className: "badge " + cls });
}

export const KINDS = { endpoint: "Endpoint", cups: "CUPS server", queue: "Queue", class: "Class", filter: "Filter", template: "Template", convert: "Convert", inbox: "Inbox", archive: "Archive", display: "Display", printer: "Printer", sink: "Sink", export: "Export" };
export const KIND_BADGE = {
  endpoint: "badge-accent badge-soft",
  cups: "badge-error badge-soft",
  sink: "badge-success badge-soft",
  export: "badge-accent badge-soft",
  filter: "badge-secondary badge-soft",
  template: "badge-secondary badge-soft",
  convert: "badge-secondary badge-soft",
  queue: "badge-primary badge-soft",
  class: "badge-primary badge-soft",
  inbox: "badge-warning badge-soft",
  archive: "badge-neutral badge-soft",
  display: "badge-info badge-soft",
  printer: "badge-neutral badge-soft",
};

export function kindBadge(kind) {
  return badge(KINDS[kind] || kind, KIND_BADGE[kind] || "badge-ghost");
}

// The kinds of payload that flow along the edges (kinds.go): a .bdf
// document, data for a template, or any file as received. Ports and
// edges are colored by them.
export const PAYLOADS = { bdf: ".bdf document", data: "data (JSON, YAML, CSV, Excel, Parquet)", any: "any file" };
export const PAYLOAD_COLOR = { bdf: "#2563eb", data: "#16a34a", any: "#64748b" };

// payloadBadge tells what a job carries: its kind and format.
export function payloadBadge(j) {
  const kind = j.kind || "bdf";
  const b = badge(kind === "bdf" ? ".bdf" : kind + " · " + (j.format || ""), "badge-sm badge-outline whitespace-nowrap");
  b.style.borderColor = PAYLOAD_COLOR[kind] || PAYLOAD_COLOR.any;
  b.style.color = PAYLOAD_COLOR[kind] || PAYLOAD_COLOR.any;
  b.title = (PAYLOADS[kind] || kind) + (j.format ? ", received as " + j.format : "");
  return b;
}

// sender is who sent a job, as verified, and the name they gave.
export function sender(j) {
  const by = j.submitted_by || {};
  const who = by.name || by.id || by.kind || "";
  return j.claimed_user && j.claimed_user !== who ? `${who} (as ${j.claimed_user})` : who;
}

export function legsOf(j) {
  return (j.legs || [])
    .map((l) => l.node + ":" + l.state + (l.member ? "(" + l.member + ")" : "") + (l.error ? " ! " + l.error : ""))
    .join(" → ");
}

export const ENDED = ["completed", "canceled", "aborted"];

export function stateBadge(state) {
  const cls = state === "completed" ? "badge-success badge-soft" : state === "held" ? "badge-warning" : state === "processing" ? "badge-info" : "badge-ghost";
  return el("span", state, { className: "badge whitespace-nowrap " + cls });
}

// cancellable tells a job the administrator may still cancel: one not
// ended, or one whose slip is still open on a board.
export function cancellable(j) {
  return !ENDED.includes(j.state) || (j.legs || []).some((l) => l.kind === "display" && l.state === "shown");
}

// keptSection lists the jobs a sink keeps that the viewer sees, newest
// first, with the first page of each.
export async function keptSection(name) {
  const box = el("div", undefined, { className: "mt-3 flex flex-col gap-2 border-t border-base-300 pt-3" });
  box.append(el("h4", "Kept here", { className: "font-semibold" }));
  const jobs = await api("GET", `/api/nodes/${encodeURIComponent(name)}/latest`);
  if (!jobs.length) box.append(el("p", "Nothing yet.", { className: "text-xs text-base-content/60" }));
  const grid = el("div", undefined, { className: "grid grid-cols-2 gap-2" });
  for (const j of jobs) {
    const card = el("a", undefined, { className: "flex flex-col gap-1 rounded-box bg-base-200 p-2 text-xs hover:bg-base-300", href: `/api/jobs/${j.id}/document`, title: `#${j.number} ${j.name || ""}` });
    card.append(el("img", undefined, { src: `/api/jobs/${j.id}/pages/1.png?w=240`, alt: `page 1 of #${j.number}`, className: "aspect-[1/1.3] w-full rounded bg-white object-cover object-top shadow-sm", loading: "lazy" }));
    card.append(el("span", `#${j.number} ${j.name || ""}`, { className: "truncate font-semibold" }));
    card.append(el("span", `${j.pages} page${j.pages === 1 ? "" : "s"} · ${since(j.received)}`, { className: "text-base-content/60" }));
    grid.append(card);
  }
  box.append(grid);
  return box;
}

// listedSection shows the entries of an export the viewer may read:
// what an archive before it wrote, with the owner each is listed for.
export async function listedSection(name) {
  const box = el("div", undefined, { className: "mt-3 flex flex-col gap-2 border-t border-base-300 pt-3" });
  box.append(el("h4", "Listed here", { className: "font-semibold" }));
  const entries = await api("GET", `/out/${encodeURIComponent(name)}?limit=20`);
  if (!entries.length) box.append(el("p", "Nothing you may read yet.", { className: "text-xs text-base-content/60" }));
  const list = el("ul", undefined, { className: "flex flex-col gap-1 text-xs" });
  for (const e of entries) {
    const li = el("li", undefined, { className: "flex flex-wrap items-center gap-x-2 rounded-box bg-base-200 px-2 py-1" });
    li.append(el("a", `#${e.job} ${e.name || ""}`, { className: "link font-semibold", href: `/out/${encodeURIComponent(name)}/${e.job}/document`, title: "download the document from the archive" }));
    if (e.owner) li.append(el("span", e.owner, { className: "badge badge-ghost badge-xs" }));
    li.append(el("span", `${e.pages ? e.pages + " page" + (e.pages === 1 ? "" : "s") + " · " : ""}${e.format || e.kind} · ${since(e.received)}`, { className: "text-base-content/60" }));
    for (const [k, v] of Object.entries(e.fields || {})) li.append(el("span", `${k}: ${v}`, { className: "font-mono text-base-content/60" }));
    list.append(li);
  }
  box.append(list);
  return box;
}

// the nav's own tab
for (const a of document.querySelectorAll("#admin-nav a")) {
  if (a.getAttribute("href") === location.pathname) a.classList.add("tab-active");
}

// Everyone's view of the topology: the nodes and the edges as the
// administrator laid them out, with the activity of the jobs the person
// may see (their own and their groups'): counts at the nodes, and the
// jobs of the last minutes running along the edges, polled every few
// seconds with the activity's revision. A node shows the jobs sent to it
// that the person sees, and a form to print there when they may. Nothing
// is changed here.

import { $, api, el, kindBadge, KINDS, legsOf, sender, since, stateBadge, status } from "./admin.js";
import { createCanvas } from "./canvas.js";

const POLL_MS = 4000;

let topo = [];
let dests = []; // the nodes the person may send to
let jobs = []; // the jobs the person sees, latest first
let live = null;
let etag = "";
let selected = null;
let fitted = false;

const canvas = createCanvas($("canvas"), {
  onSelect(name) {
    selected = name;
    renderSidebar();
    $("sidebar").scrollTop = 0;
  },
});

async function load() {
  try {
    const [t, l, d, j] = await Promise.all([api("GET", "/api/topology"), api("GET", "/api/layout"), api("GET", "/api/destinations"), api("GET", "/api/my/jobs?limit=200")]);
    topo = t;
    dests = d;
    jobs = j;
    canvas.setData(topo, l.positions || {});
    if (!fitted) {
      canvas.fit();
      fitted = true;
    }
    if (selected && !topo.some((n) => n.name === selected)) selected = null;
    renderSidebar();
  } catch (e) {
    status(e.message, true);
  }
}

// poll reads the activity when its revision changed, and the jobs with
// it.
async function poll() {
  try {
    const res = await fetch("/api/topology/live", { headers: etag ? { "If-None-Match": etag } : {} });
    if (res.status === 401) {
      location.href = "/";
      return;
    }
    if (res.status === 200) {
      etag = res.headers.get("ETag") || "";
      live = await res.json();
      canvas.setLive(live);
      jobs = await api("GET", "/api/my/jobs?limit=200");
      renderSidebar();
    } else if (res.status !== 304) {
      throw new Error("live " + res.status);
    }
    $("live").textContent = "Live · " + new Date().toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit", second: "2-digit" });
    $("live").className = "badge badge-success badge-soft badge-sm";
  } catch (e) {
    $("live").textContent = "Not connected: " + e.message;
    $("live").className = "badge badge-error badge-sm";
  }
}

function countsOf(name) {
  const c = (live && live.nodes[name]) || {};
  const out = [];
  if (c.open) out.push([String(c.open), "open"]);
  if (c.held) out.push([String(c.held), "held"]);
  if (c.pending) out.push([String(c.pending), "waiting"]);
  if (c.failed) out.push([String(c.failed), "failed"]);
  out.push([String(c.recent || 0), "in the last hour"]);
  return out;
}

function renderSidebar() {
  const aside = $("sidebar");
  const n = topo.find((x) => x.name === selected);
  if (!n) {
    const p = el("p", undefined, { className: "text-sm text-base-content/60" });
    p.append("Click a node to see the jobs sent there that you may see, and to print there when you may. ");
    if (live && Object.keys(live.nodes).length === 0) p.append("Nothing of yours is in the network right now.");
    aside.replaceChildren(p);
    return;
  }
  const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  head.append(kindBadge(n.kind), el("h3", n.title && n.title !== n.name ? n.title : n.name, { className: "text-lg font-semibold" }));
  if (n.title && n.title !== n.name) head.append(el("span", n.name, { className: "font-mono text-sm text-base-content/60" }));
  const stats = el("div", undefined, { className: "stats stats-horizontal flex-wrap bg-base-200" });
  for (const [v, label] of countsOf(n.name)) {
    const s = el("div", undefined, { className: "stat px-3 py-2" });
    s.append(el("div", v, { className: "stat-value text-xl" }), el("div", label, { className: "stat-title text-xs" }));
    stats.append(s);
  }
  const parts = [head, stats];
  const out = [];
  for (const to of n.to || []) out.push(to);
  for (const m of n.members || []) out.push(m + (n.policy === "failover" ? " (first)" : " (in turn)"));
  if (n.fallback && n.fallback.printer) out.push(n.fallback.printer + " (fallback)");
  if (out.length) parts.push(el("p", "→ " + out.join(", "), { className: "font-mono text-xs text-base-content/60" }));
  const dest = dests.find((d) => d.name === n.name);
  if (dest) {
    const form = el("form", undefined, { className: "flex flex-col gap-2 border-t border-base-300 pt-3" });
    form.append(el("h4", "Print here", { className: "font-semibold" }));
    const file = el("input", undefined, { type: "file", className: "file-input file-input-sm w-full", required: true });
    const name = el("input", undefined, { className: "input input-sm w-full", placeholder: "job name (the file's)" });
    form.append(file, name, el("button", "Print", { type: "submit", className: "btn btn-sm btn-primary w-fit" }));
    form.addEventListener("submit", async (ev) => {
      ev.preventDefault();
      const f = file.files[0];
      if (!f) return;
      try {
        const j = await api("POST", `/api/destinations/${encodeURIComponent(n.name)}/jobs?name=${encodeURIComponent(name.value || f.name)}`, f, true);
        status(`Job #${j.number} is ${j.state}`);
        file.value = "";
        name.value = "";
        etag = "";
        await poll();
      } catch (e) {
        status(e.message, true);
      }
    });
    parts.push(form);
  } else {
    parts.push(el("p", "You may not print here.", { className: "text-xs text-base-content/60"}));
  }
  const mine = jobs.filter((j) => j.dest === n.name || (j.legs || []).some((l) => l.node === n.name));
  const sec = el("div", undefined, { className: "flex flex-col gap-2 border-t border-base-300 pt-3" });
  sec.append(el("h4", `Jobs here (${mine.length})`, { className: "font-semibold" }));
  if (!mine.length) sec.append(el("p", "None of yours.", { className: "text-xs text-base-content/60" }));
  const list = el("ul", undefined, { className: "flex flex-col gap-2" });
  for (const j of mine.slice(0, 20)) {
    const li = el("li", undefined, { className: "rounded-box bg-base-200 p-2 text-sm" });
    const row = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
    row.append(el("span", "#" + j.number, { className: "font-mono" }), el("span", j.name || "", { className: "font-semibold" }), stateBadge(j.state));
    li.append(row);
    li.append(el("div", `${sender(j)} · ${since(j.received)}`, { className: "text-xs text-base-content/60" }));
    li.append(el("div", legsOf(j), { className: "font-mono text-xs text-base-content/60 break-all", title: legsOf(j) }));
    li.append(el("a", "Document", { className: "link text-xs", href: `/api/jobs/${j.id}/document` }));
    list.append(li);
  }
  sec.append(list);
  parts.push(sec);
  aside.replaceChildren(...parts);
}

$("zoom-in").addEventListener("click", () => canvas.zoomBy(1.25));
$("zoom-out").addEventListener("click", () => canvas.zoomBy(0.8));
$("zoom-fit").addEventListener("click", () => canvas.fit());
if (location.hash.length > 1) selected = decodeURIComponent(location.hash.slice(1));
load().then(poll);
setInterval(() => {
  if (!document.hidden) poll();
}, POLL_MS);
setInterval(() => {
  if (!document.hidden) load();
}, 60000);
document.addEventListener("visibilitychange", () => {
  if (!document.hidden) poll();
});

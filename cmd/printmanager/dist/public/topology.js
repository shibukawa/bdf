// The editor of the topology. Nodes of every kind are placed on a canvas
// (SVG, panned and zoomed), connected by dragging from a node's output
// port onto another node, and set up in the sidebar: the fields of the
// node's kind, and what is at it (an endpoint's templates and key, a
// display's devices and enrollment codes). Everything goes through the
// API; the server checks a node when it is put (its edges, its name, no
// cycle).

import { $, api, badge, button, changed, el, kindBadge, KINDS, onChange, since, status } from "./admin.js";

const SVG = "http://www.w3.org/2000/svg";
const NODE_W = 176;
const NODE_H = 54;
const COL = 240;
const ROW = 84;
const PAD = 24;
const GRID = 8;
const FILL = { endpoint: "#fde68a", queue: "#ddd6fe", class: "#c7d2fe", inbox: "#fed7aa", archive: "#e5e7eb", display: "#bfdbfe", printer: "#d1d5db" };
// the kinds with an output port: they send to other nodes (a display to
// its fallback printer)
const OUT_PORT = new Set(["endpoint", "queue", "class", "inbox", "display"]);
const SCHEMES = ["ipp", "ipps", "socket", "lpd"];
const ACCEPTS = { both: "documents, and templates to fill", file: "documents only", template: "templates to fill only" };

let topo = []; // the nodes, as /api/topology lists them
let layout = {}; // name -> {x, y}
let formats = []; // the formats an endpoint may take
let selected = null; // the node in the sidebar
let draft = null; // a node being added: {kind}
let draftPos = null;
let view = { x: PAD, y: PAD, k: 1 }; // the canvas's pan and zoom
let fitted = false;

const byName = () => Object.fromEntries(topo.map((n) => [n.name, n]));

async function load() {
  try {
    const [t, l] = await Promise.all([api("GET", "/api/topology"), api("GET", "/api/layout")]);
    topo = t;
    layout = l.positions || {};
    if (!formats.length) formats = await api("GET", "/api/formats");
    autoLayout();
    $("node-names").replaceChildren(...topo.map((n) => el("option", undefined, { value: n.name })));
    if (selected && !topo.some((n) => n.name === selected)) selected = null;
    if (!fitted) {
      fit();
      fitted = true;
    }
    draw();
    renderSidebar();
  } catch (e) {
    status(e.message, true);
  }
}

// The canvas.

function svgEl(tag, attrs = {}) {
  const e = document.createElementNS(SVG, tag);
  for (const [k, v] of Object.entries(attrs)) e.setAttribute(k, v);
  return e;
}

// edgesOf are the edges from a node: to (a queue's, an endpoint's, an
// inbox's), member (a class's) or fallback (a display's printer).
function edgesOf(n) {
  const out = [];
  for (const to of n.to || []) out.push({ to, kind: "to" });
  for (const to of n.members || []) out.push({ to, kind: "member" });
  if (n.fallback && n.fallback.printer) out.push({ to: n.fallback.printer, kind: "fallback" });
  return out;
}

// autoLayout places the nodes that have no position: in columns by their
// depth from the entries, rows within a column.
function autoLayout() {
  const nodes = byName();
  const incoming = {};
  for (const n of topo) for (const e of edgesOf(n)) incoming[e.to] = (incoming[e.to] || 0) + 1;
  const depth = {};
  const walk = (name, d) => {
    if ((depth[name] ?? -1) >= d) return;
    depth[name] = d;
    if (d > 20) return;
    for (const e of edgesOf(nodes[name] || {})) if (nodes[e.to]) walk(e.to, d + 1);
  };
  for (const n of topo) if (!incoming[n.name]) walk(n.name, 0);
  for (const n of topo) if (depth[n.name] === undefined) walk(n.name, 0);
  const columns = {};
  for (const n of topo) (columns[depth[n.name]] ||= []).push(n.name);
  const taken = new Set(Object.values(layout).map((p) => `${Math.round(p.x)},${Math.round(p.y)}`));
  for (const [d, names] of Object.entries(columns)) {
    let row = 0;
    for (const name of names) {
      if (layout[name]) continue;
      let p;
      do {
        p = { x: d * COL, y: row * ROW };
        row++;
      } while (taken.has(`${p.x},${p.y}`));
      layout[name] = p;
    }
  }
}

function pos(name) {
  return layout[name] || { x: 0, y: 0 };
}

// freeSpot is a place for a new node near a point, clear of the others:
// to the right of them, then on the next row.
function freeSpot(c) {
  const taken = (p) => topo.some((n) => Math.abs(pos(n.name).x - p.x) < NODE_W + GRID && Math.abs(pos(n.name).y - p.y) < NODE_H + GRID);
  const start = { x: snap(c.x - NODE_W / 2), y: snap(c.y - NODE_H / 2) };
  let p = { ...start };
  for (let i = 0; i < 64 && taken(p); i++) {
    p = { x: p.x + COL, y: p.y };
    if (p.x > start.x + 3 * COL) p = { x: start.x, y: p.y + ROW };
  }
  return p;
}

function snap(v) {
  return Math.round(v / GRID) * GRID;
}

function clipText(s, n) {
  return s.length > n ? s.slice(0, n - 1) + "…" : s;
}

// counts are what is at a node now.
function counts(n) {
  const out = [];
  if (n.held) out.push(`${n.held} held`);
  if (n.open) out.push(`${n.open} open`);
  if (n.devices) out.push(`${n.devices} device${n.devices === 1 ? "" : "s"}`);
  if (n.pending) out.push(`${n.pending} waiting`);
  return out;
}

function draw() {
  const svg = $("canvas");
  svg.replaceChildren();
  const defs = svgEl("defs");
  const marker = svgEl("marker", { id: "arrow", viewBox: "0 0 10 10", refX: "10", refY: "5", markerWidth: "8", markerHeight: "8", orient: "auto-start-reverse" });
  marker.append(svgEl("path", { d: "M 0 0 L 10 5 L 0 10 z", fill: "#64748b" }));
  defs.append(marker);
  const world = svgEl("g", { id: "world", transform: `translate(${view.x} ${view.y}) scale(${view.k})` });
  const nodes = byName();
  for (const n of topo) {
    for (const e of edgesOf(n)) {
      if (!nodes[e.to]) continue;
      const a = pos(n.name);
      const b = pos(e.to);
      const x1 = a.x + NODE_W;
      const y1 = a.y + NODE_H / 2;
      const x2 = b.x;
      const y2 = b.y + NODE_H / 2;
      const dx = Math.max(40, Math.abs(x2 - x1) / 2);
      const d = `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`;
      const hit = svgEl("path", { d, fill: "none", stroke: "transparent", "stroke-width": "14", class: "cursor-pointer" });
      const line = svgEl("path", { d, fill: "none", stroke: e.kind === "fallback" ? "#94a3b8" : "#64748b", "stroke-width": "2", "marker-end": "url(#arrow)", "pointer-events": "none" });
      if (e.kind === "member") line.setAttribute("stroke-dasharray", "6 4");
      if (e.kind === "fallback") line.setAttribute("stroke-dasharray", "2 4");
      const title = svgEl("title");
      title.textContent = e.kind === "fallback" ? `${n.name} prints on ${e.to} as a fallback (click to remove)` : `${n.name} → ${e.to} (click to remove)`;
      hit.append(title);
      hit.addEventListener("pointerdown", (ev) => ev.stopPropagation());
      hit.addEventListener("click", () => removeEdge(n, e));
      world.append(hit, line);
    }
  }
  for (const n of topo) world.append(nodeBox(n, pos(n.name), n.name === selected));
  if (draft && draftPos) {
    const g = svgEl("g", { transform: `translate(${draftPos.x} ${draftPos.y})`, class: "cursor-grab", "data-draft": "1" });
    g.append(svgEl("rect", { width: NODE_W, height: NODE_H, rx: "8", fill: FILL[draft.kind] || "#e5e7eb", stroke: "#4f46e5", "stroke-width": "2", "stroke-dasharray": "6 4" }));
    const t = svgEl("text", { x: "10", y: "22", "font-size": "13", "font-weight": "600", fill: "#111827" });
    t.textContent = ($("f-name") && $("f-name").value) || "new " + (KINDS[draft.kind] || draft.kind).toLowerCase();
    const k = svgEl("text", { x: "10", y: "40", "font-size": "11", fill: "#374151" });
    k.textContent = (KINDS[draft.kind] || draft.kind) + " · not saved yet";
    g.append(t, k);
    g.addEventListener("pointerdown", (ev) => startDrag(ev, null));
    world.append(g);
  }
  svg.append(defs, world);
}

// nodeBox draws a node: its title, its kind and what is at it, an input
// port at the left and, for the kinds that send on, an output port at the
// right.
function nodeBox(n, p, isSelected) {
  const g = svgEl("g", { transform: `translate(${p.x} ${p.y})`, class: "cursor-grab", "data-name": n.name });
  g.append(svgEl("rect", { width: NODE_W, height: NODE_H, rx: "8", fill: FILL[n.kind] || "#e5e7eb", stroke: isSelected ? "#4f46e5" : "#94a3b8", "stroke-width": isSelected ? "3" : "1.5" }));
  const name = svgEl("text", { x: "10", y: "21", "font-size": "13", "font-weight": "600", fill: "#111827" });
  name.textContent = clipText(n.title && n.title !== n.name ? n.title : n.name, 22);
  const sub = svgEl("text", { x: "10", y: "39", "font-size": "11", fill: "#374151" });
  const parts = [KINDS[n.kind] || n.kind];
  if (n.title && n.title !== n.name) parts.push(n.name);
  parts.push(...counts(n));
  sub.textContent = clipText(parts.join(" · "), 30);
  const title = svgEl("title");
  title.textContent = `${n.name} (${KINDS[n.kind] || n.kind})${n.title ? ": " + n.title : ""}`;
  g.append(title, name, sub);
  g.append(svgEl("circle", { cx: "0", cy: NODE_H / 2, r: "5", fill: "#f8fafc", stroke: "#64748b", "stroke-width": "1.5" }));
  if (OUT_PORT.has(n.kind)) {
    const port = svgEl("circle", { cx: NODE_W, cy: NODE_H / 2, r: "7", fill: "#4f46e5", stroke: "#fff", "stroke-width": "2", class: "cursor-crosshair", "data-port": n.name });
    const t = svgEl("title");
    t.textContent = n.kind === "display" ? "Drag onto a printer: the fallback printer" : n.kind === "class" ? "Drag onto a node: a member" : "Drag onto a node to connect";
    port.append(t);
    g.append(port);
  }
  g.addEventListener("pointerdown", (ev) => startDrag(ev, n));
  return g;
}

// Pan and zoom: the world is translated by view.x, view.y and scaled by
// view.k. Screen points are relative to the canvas.

function screenPoint(ev) {
  const r = $("canvas").getBoundingClientRect();
  return { x: ev.clientX - r.left, y: ev.clientY - r.top };
}

function toWorld(s) {
  return { x: (s.x - view.x) / view.k, y: (s.y - view.y) / view.k };
}

function zoomAt(s, factor) {
  const k = Math.min(3, Math.max(0.25, view.k * factor));
  view = { x: s.x - (s.x - view.x) * (k / view.k), y: s.y - (s.y - view.y) * (k / view.k), k };
  draw();
}

function canvasSize() {
  const r = $("canvas").getBoundingClientRect();
  return { w: r.width || 800, h: r.height || 400 };
}

function center() {
  const { w, h } = canvasSize();
  return toWorld({ x: w / 2, y: h / 2 });
}

// fit shows every node.
function fit() {
  const { w, h } = canvasSize();
  if (!topo.length) {
    view = { x: PAD, y: PAD, k: 1 };
    return;
  }
  let minX = Infinity;
  let minY = Infinity;
  let maxX = -Infinity;
  let maxY = -Infinity;
  for (const n of topo) {
    const p = pos(n.name);
    minX = Math.min(minX, p.x);
    minY = Math.min(minY, p.y);
    maxX = Math.max(maxX, p.x + NODE_W);
    maxY = Math.max(maxY, p.y + NODE_H);
  }
  const k = Math.min(1.5, Math.max(0.25, Math.min((w - 2 * PAD) / (maxX - minX), (h - 2 * PAD) / (maxY - minY))));
  view = { x: (w - (maxX - minX) * k) / 2 - minX * k, y: (h - (maxY - minY) * k) / 2 - minY * k, k };
  draw();
}

// Dragging: the background pans, a node moves (a click selects it), a
// port draws an edge, two pointers pinch.
let drag = null;
const pointers = new Map();
let pinch = null;

function startDrag(ev, n) {
  if (ev.button !== 0) return;
  ev.preventDefault();
  ev.stopPropagation();
  const s = screenPoint(ev);
  const svg = $("canvas");
  if (n && ev.target.dataset.port) {
    const p = pos(n.name);
    const line = svgEl("path", { fill: "none", stroke: "#4f46e5", "stroke-width": "2", "stroke-dasharray": "4 3", "pointer-events": "none" });
    $("world").append(line);
    drag = { kind: "edge", from: n, line, x1: p.x + NODE_W, y1: p.y + NODE_H / 2 };
  } else {
    const p = n ? pos(n.name) : draftPos;
    const w = toWorld(s);
    drag = { kind: "move", node: n, dx: w.x - p.x, dy: w.y - p.y, moved: false };
  }
  svg.setPointerCapture(ev.pointerId);
}

const canvas = $("canvas");

canvas.addEventListener("pointerdown", (ev) => {
  pointers.set(ev.pointerId, screenPoint(ev));
  if (pointers.size === 2) {
    const [a, b] = [...pointers.values()];
    pinch = { d0: Math.hypot(a.x - b.x, a.y - b.y) || 1, k0: view.k, c0: { x: (a.x + b.x) / 2, y: (a.y + b.y) / 2 }, x0: view.x, y0: view.y };
    if (drag && drag.kind === "edge") drag.line.remove();
    drag = null;
    return;
  }
  if (ev.button !== 0) return;
  const s = screenPoint(ev);
  drag = { kind: "pan", sx: s.x, sy: s.y, x0: view.x, y0: view.y, moved: false };
  canvas.setPointerCapture(ev.pointerId);
});

canvas.addEventListener("pointermove", (ev) => {
  if (pointers.has(ev.pointerId)) pointers.set(ev.pointerId, screenPoint(ev));
  if (pinch && pointers.size === 2) {
    const [a, b] = [...pointers.values()];
    const d = Math.hypot(a.x - b.x, a.y - b.y) || 1;
    const k = Math.min(3, Math.max(0.25, (pinch.k0 * d) / pinch.d0));
    const c = { x: (a.x + b.x) / 2, y: (a.y + b.y) / 2 };
    view = { x: c.x - ((pinch.c0.x - pinch.x0) / pinch.k0) * k, y: c.y - ((pinch.c0.y - pinch.y0) / pinch.k0) * k, k };
    draw();
    return;
  }
  if (!drag) return;
  const s = screenPoint(ev);
  if (drag.kind === "pan") {
    view = { ...view, x: drag.x0 + s.x - drag.sx, y: drag.y0 + s.y - drag.sy };
    drag.moved = drag.moved || Math.abs(s.x - drag.sx) + Math.abs(s.y - drag.sy) > 3;
    draw();
  } else if (drag.kind === "move") {
    const w = toWorld(s);
    const p = { x: snap(w.x - drag.dx), y: snap(w.y - drag.dy) };
    if (drag.node) layout[drag.node.name] = p;
    else draftPos = p;
    drag.moved = true;
    draw();
  } else {
    const w = toWorld(s);
    drag.line.setAttribute("d", `M ${drag.x1} ${drag.y1} L ${w.x} ${w.y}`);
  }
});

async function endDrag(ev) {
  pointers.delete(ev.pointerId);
  if (pointers.size < 2) pinch = null;
  if (!drag) return;
  const d = drag;
  drag = null;
  try {
    switch (d.kind) {
      case "pan":
        if (!d.moved) {
          selected = null;
          draw();
          renderSidebar();
        }
        return;
      case "move":
        if (!d.node) return;
        if (d.moved) await api("PUT", "/api/layout", { positions: { [d.node.name]: layout[d.node.name] } });
        else select(d.node.name);
        return;
      case "edge": {
        d.line.remove();
        const target = document.elementFromPoint(ev.clientX, ev.clientY)?.closest("g[data-name]");
        if (target) await addEdge(d.from, target.dataset.name);
      }
    }
  } catch (e) {
    status(e.message, true);
  }
}
canvas.addEventListener("pointerup", endDrag);
canvas.addEventListener("pointercancel", endDrag);
canvas.addEventListener(
  "wheel",
  (ev) => {
    ev.preventDefault();
    zoomAt(screenPoint(ev), Math.exp(-ev.deltaY * 0.0015));
  },
  { passive: false },
);
$("zoom-in").addEventListener("click", () => zoomAt({ x: canvasSize().w / 2, y: canvasSize().h / 2 }, 1.25));
$("zoom-out").addEventListener("click", () => zoomAt({ x: canvasSize().w / 2, y: canvasSize().h / 2 }, 0.8));
$("zoom-fit").addEventListener("click", fit);
$("refresh").addEventListener("click", load);
window.addEventListener("resize", draw);

// addEdge connects a node to a target through the API, which checks the
// edge (the target exists, no cycle): a class's member, a display's
// fallback printer, or the node's to.
async function addEdge(from, to) {
  if (from.name === to) return;
  const n = await api("GET", "/api/nodes/" + from.name);
  if (n.kind === "display") {
    n.fallback = { ...(n.fallback || {}), printer: to };
  } else {
    const field = n.kind === "class" ? "members" : "to";
    const list = n[field] || [];
    if (list.includes(to)) return;
    n[field] = [...list, to];
  }
  delete n.name;
  await api("PUT", "/api/nodes/" + from.name, n);
  status(`${from.name} → ${to}`);
  await load();
}

async function removeEdge(from, e) {
  const what = e.kind === "fallback" ? `the fallback printer ${e.to} of ${from.name}` : `the edge ${from.name} → ${e.to}`;
  if (!confirm(`Remove ${what}?`)) return;
  try {
    const n = await api("GET", "/api/nodes/" + from.name);
    if (e.kind === "fallback") n.fallback = null;
    else {
      const field = e.kind === "member" ? "members" : "to";
      n[field] = (n[field] || []).filter((x) => x !== e.to);
    }
    delete n.name;
    await api("PUT", "/api/nodes/" + from.name, n);
    await load();
  } catch (err) {
    status(err.message, true);
  }
}

function select(name) {
  selected = name;
  draft = null;
  draftPos = null;
  draw();
  renderSidebar();
  $("sidebar").scrollTop = 0;
}

// Adding a node: the palette makes a draft, placed at the centre of the
// view; the sidebar takes its name and settings, and Save puts it.
$("palette").addEventListener("click", (ev) => {
  const a = ev.target.closest("a[data-kind]");
  if (!a) return;
  ev.preventDefault();
  document.activeElement?.blur();
  draft = { kind: a.dataset.kind, name: "", title: "" };
  draftPos = freeSpot(center());
  selected = null;
  draw();
  renderSidebar();
  $("sidebar").scrollTop = 0;
  $("f-name")?.focus();
});

// The sidebar: the node selected (or the draft), its kind's fields, and
// what is at it.

function field(label, input, hint) {
  const f = el("fieldset", undefined, { className: "fieldset" });
  f.append(el("legend", label, { className: "fieldset-legend" }), input);
  if (hint) f.append(el("p", hint, { className: "label whitespace-normal" }));
  return f;
}

function input(id, value, attrs = {}) {
  // list is read-only as a property: an attribute
  const { list, ...rest } = attrs;
  const e = el("input", undefined, { id, className: "input input-sm w-full", value: value ?? "", ...rest });
  if (list) e.setAttribute("list", list);
  return e;
}

function selectInput(id, options, value) {
  const s = el("select", undefined, { id, className: "select select-sm w-full" });
  for (const [v, label] of options) s.append(el("option", label, { value: v, selected: v === value }));
  return s;
}

const names = (s) => s.split(",").map((x) => x.trim()).filter(Boolean);

function hint(text) {
  return el("p", text, { className: "text-xs text-base-content/60" });
}

function section(title) {
  const d = el("div", undefined, { className: "mt-3 flex flex-col gap-2 border-t border-base-300 pt-3" });
  d.append(el("h4", title, { className: "font-semibold" }));
  return d;
}

// scopeBadge tells who may send to a node, and whose its jobs are.
function scopeBadge(n) {
  const scope = n.scope || "shared";
  const cls = { group: "badge badge-sm badge-soft badge-secondary whitespace-nowrap", personal: "badge badge-sm badge-soft badge-accent whitespace-nowrap" }[scope] || "badge badge-sm badge-ghost whitespace-nowrap";
  return el("span", scope + (n.owner ? ": " + n.owner : ""), { className: cls });
}

// parsePrinterURL splits a printer's URL into the sidebar's fields.
function parsePrinterURL(u) {
  if (!u) return { scheme: "ipp", host: "", port: "", path: "" };
  try {
    const x = new URL(u);
    return { scheme: x.protocol.replace(":", ""), host: x.hostname, port: x.port, path: x.pathname === "/" && !x.search ? "" : x.pathname + x.search };
  } catch {
    return { scheme: "ipp", host: u, port: "", path: "" };
  }
}

// composeURL is the printer's URL of the sidebar's fields.
function composeURL() {
  const scheme = $("f-scheme").value;
  const host = $("f-host").value.trim();
  const port = $("f-port").value.trim();
  let path = $("f-path").value.trim();
  if (!host) return "";
  if (!path && scheme.startsWith("ipp")) path = "/ipp/print";
  if (path && !path.startsWith("/")) path = "/" + path;
  return `${scheme}://${host}${port ? ":" + port : ""}${path}`;
}

function renderSidebar() {
  const aside = $("sidebar");
  const n = draft || topo.find((x) => x.name === selected);
  if (!n) {
    aside.replaceChildren(el("p", "Click a node to edit it here, or add one with “+ Add a node”. Drag from a node's right port onto another node to connect them; a display's port sets its fallback printer.", { className: "text-sm text-base-content/60" }));
    return;
  }
  const kind = n.kind;
  const form = el("form", undefined, { className: "flex flex-col gap-1" });
  const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  head.append(kindBadge(kind));
  if (draft) head.append(el("span", "new", { className: "badge badge-outline" }));
  else head.append(el("h3", n.name, { className: "font-mono text-lg font-semibold" }), scopeBadge(n));
  form.append(head);
  if (draft) {
    const name = input("f-name", "", { placeholder: "intake, both, tokyo, review, wh-a, laser, audit", required: true, pattern: "[A-Za-z0-9][A-Za-z0-9_.\\-]{0,63}", autocomplete: "off" });
    name.addEventListener("input", draw);
    form.append(field("Name", name, "Letters, digits, . _ and -. It is the node's IPP destination (/ipp/print/<name>) and does not change later."));
  }
  form.append(field("Title", input("f-title", n.title || "", { placeholder: "shown instead of the name" })));
  switch (kind) {
    case "endpoint": {
      form.append(field("Forwards to", input("f-to", (n.to || []).join(", "), { placeholder: "nodes, comma-separated", list: "node-names" }), "Every node named gets the job. The port on the canvas adds one too."));
      form.append(field("Accepts", selectInput("f-accept", Object.entries(ACCEPTS), n.accept || "both")));
      const fmts = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
      for (const f of formats) {
        const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
        l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "fmt", value: f, checked: (n.formats || []).includes(f) }), f);
        fmts.append(l);
      }
      form.append(field("Formats taken", fmts, "None checked: every format the converter reads. A document of another format is refused, and the refusal is in the audit log."));
      break;
    }
    case "queue":
      form.append(field("Forwards to", input("f-to", (n.to || []).join(", "), { placeholder: "nodes, comma-separated", list: "node-names", required: true }), "Every node named gets the job: a fan-out when there are several."));
      break;
    case "class":
      form.append(field("Members", input("f-members", (n.members || []).join(", "), { placeholder: "mfp-3f, mfp-4f", list: "node-names", required: true })));
      form.append(field("Policy", selectInput("f-policy", [["round-robin", "in turn (round-robin)"], ["failover", "the first member (failover)"]], n.policy || "round-robin"), "A printer that does not take the job is skipped for the next member."));
      break;
    case "inbox":
      form.append(field("A person releases to", input("f-to", (n.to || []).join(", "), { placeholder: "nodes, comma-separated; empty: any node", list: "node-names" })));
      form.append(field("Keep", input("f-keep", n.keep || "", { placeholder: "14d" }), "How long released and discarded items stay shown."));
      break;
    case "archive":
      form.append(field("Store", input("f-store", n.store || "", { placeholder: "s3://bucket?region=… or bbolt:path", required: true }), "A documents URL: bbolt, s3, r2 or gcs. Each job's .bdf and record are written once; the manager never changes or deletes them."));
      form.append(field("Collection", input("f-collection", n.collection || "", { placeholder: "archive-" + (n.name || "<name>") })));
      break;
    case "display": {
      const printers = topo.filter((x) => x.kind === "printer").map((x) => [x.name, x.name + (x.title ? " · " + x.title : "")]);
      form.append(field("Fallback printer", selectInput("f-fallback", [["", "none"], ...printers], n.fallback ? n.fallback.printer : ""), "Prints the slips nobody opened in time, and every unopened slip while no device is connected. The port on the canvas sets it too."));
      form.append(field("Print after", input("f-after", n.fallback && n.fallback.after ? n.fallback.after : "", { placeholder: "10m; empty: only while no device is connected" })));
      form.append(field("Keep", input("f-keep", n.keep || "", { placeholder: "14d" }), "How long slips done or canceled stay on the devices."));
      form.append(field("Offline after", input("f-offline", n.offline_after || "", { placeholder: "2m" }), "How long a device may be silent and still count as connected."));
      break;
    }
    case "printer": {
      const u = parsePrinterURL(n.url);
      const row = el("div", undefined, { className: "grid grid-cols-[6rem_1fr_5rem] gap-1" });
      row.append(selectInput("f-scheme", SCHEMES.map((s) => [s, s + "://"]), SCHEMES.includes(u.scheme) ? u.scheme : "ipp"), input("f-host", u.host, { placeholder: "host or address", required: true }), input("f-port", u.port, { placeholder: "port", inputMode: "numeric" }));
      form.append(field("Protocol and address", row, "ipp and ipps are IPP printers; socket is port 9100 (raw), lpd a line printer daemon."));
      form.append(field("Path", input("f-path", u.path, { placeholder: "/ipp/print (ipp); a queue name (lpd)" })));
      const preview = el("p", "", { className: "font-mono text-xs text-base-content/60 break-all" });
      const update = () => (preview.textContent = composeURL() || "—");
      form.append(preview);
      setTimeout(update);
      for (const id of ["f-scheme", "f-host", "f-port", "f-path"]) form.querySelector("#" + id).addEventListener("input", update);
      form.append(field("Route", selectInput("f-route", [["direct", "direct: the manager sends to it"], ["branch", "branch: a bdf printd of the site fetches (not yet)"]], n.route || "direct")));
      break;
    }
  }
  form.append(field("Scope", selectInput("f-scope", [["shared", "shared: anyone signed in may send"], ["group", "group: its members send and see the jobs"], ["personal", "personal: one person"]], n.scope || "shared")));
  form.append(field("Owner", input("f-owner", n.owner || "", { placeholder: "a group (a value of the groups claim), or a person's email" })));
  const actions = el("div", undefined, { className: "mt-2 flex flex-wrap gap-2" });
  actions.append(el("button", draft ? "Add" : "Save", { type: "submit", className: "btn btn-sm btn-primary" }));
  if (draft) {
    const cancel = el("button", "Cancel", { type: "button", className: "btn btn-sm btn-ghost" });
    cancel.addEventListener("click", () => {
      draft = null;
      draftPos = null;
      draw();
      renderSidebar();
    });
    actions.append(cancel);
  } else {
    const del = el("button", "Delete", { type: "button", className: "btn btn-sm btn-ghost text-error" });
    del.addEventListener("click", () => remove(n.name));
    actions.append(del);
  }
  form.append(actions);
  form.addEventListener("submit", save);
  aside.replaceChildren(form);
  if (!draft) extras(n, aside);
}

async function save(ev) {
  ev.preventDefault();
  const n = draft || topo.find((x) => x.name === selected);
  if (!n) return;
  const name = draft ? $("f-name").value.trim() : n.name;
  const body = { kind: n.kind, title: $("f-title").value.trim(), scope: $("f-scope").value };
  if ($("f-owner").value.trim()) body.owner = $("f-owner").value.trim();
  const dur = (id) => ($(id).value.trim() ? $(id).value.trim() : undefined);
  switch (n.kind) {
    case "endpoint":
      body.to = names($("f-to").value);
      body.accept = $("f-accept").value;
      body.formats = [...document.querySelectorAll('input[name="fmt"]:checked')].map((c) => c.value);
      break;
    case "queue":
      body.to = names($("f-to").value);
      break;
    case "class":
      body.members = names($("f-members").value);
      body.policy = $("f-policy").value;
      break;
    case "inbox":
      body.to = names($("f-to").value);
      body.keep = dur("f-keep");
      break;
    case "archive":
      body.store = $("f-store").value.trim();
      if ($("f-collection").value.trim()) body.collection = $("f-collection").value.trim();
      break;
    case "display":
      if ($("f-fallback").value) {
        body.fallback = { printer: $("f-fallback").value };
        if (dur("f-after")) body.fallback.after = dur("f-after");
      }
      body.keep = dur("f-keep");
      body.offline_after = dur("f-offline");
      break;
    case "printer":
      body.url = composeURL();
      body.route = $("f-route").value;
      break;
  }
  try {
    const res = await api("PUT", "/api/nodes/" + encodeURIComponent(name), body);
    if (draft) {
      await api("PUT", "/api/layout", { positions: { [name]: draftPos } });
      layout[name] = draftPos;
      draft = null;
      draftPos = null;
      selected = name;
    }
    status(`${name} saved`);
    if (res.key) showKey(name, res.key);
    await load();
    $("sidebar").scrollTop = 0;
  } catch (e) {
    status(e.message, true);
  }
}

async function remove(name) {
  if (!confirm(`Delete ${name}? A node another one sends to stays. An endpoint's templates go with it; a display's slips and devices stay until it is added again.`)) return;
  try {
    await api("DELETE", "/api/nodes/" + encodeURIComponent(name));
    selected = null;
    status(`${name} deleted`);
    await load();
  } catch (e) {
    status(e.message, true);
  }
}

// extras are what is at a saved node: its jobs, an endpoint's entry, key
// and templates, a display's devices.
function extras(n, aside) {
  const links = el("div", undefined, { className: "mt-3 flex flex-col gap-1 border-t border-base-300 pt-3 text-sm" });
  links.append(el("span", "IPP: " + n.ipp, { className: "font-mono text-xs text-base-content/60 break-all" }));
  const row = el("div", undefined, { className: "flex flex-wrap gap-3" });
  row.append(el("a", "Jobs sent here", { className: "link", href: "/admin/jobs?dest=" + encodeURIComponent(n.name) }));
  if (n.kind === "inbox") row.append(el("a", `${n.held || 0} held: open the inbox`, { className: "link", href: "/admin/inboxes" }));
  if (n.pending) row.append(el("a", `${n.pending} deliveries waiting`, { className: "link", href: "/admin/deliveries" }));
  links.append(row);
  aside.append(links);
  if (n.kind === "endpoint") {
    const entry = section("Entry");
    entry.append(hint("Applications POST a document, or JSON naming a template and its data, with the key as a bearer token."));
    entry.append(el("pre", `POST ${location.origin}${n.in}`, { className: "overflow-x-auto rounded-box bg-base-200 p-2 font-mono text-xs break-all whitespace-pre-wrap" }));
    const key = section("Key");
    key.append(hint("Shown once, when the endpoint was made or the key rotated; the manager keeps only its hash. Rotating it makes a new one and the old one stops working."));
    key.append(button("Rotate key", async () => showKey(n.name, (await api("POST", `/api/nodes/${n.name}/key`)).key), "btn btn-sm btn-outline w-fit"));
    const tpls = section("Templates");
    aside.append(entry, key, tpls);
    templatesSection(n, tpls).catch((e) => status(e.message, true));
  }
  if (n.kind === "display") {
    const dev = section("Devices");
    aside.append(dev);
    devicesSection(n, dev).catch((e) => status(e.message, true));
  }
}

async function templatesSection(n, box) {
  box.append(hint(`This endpoint fills them with the data it is sent: POST ${n.in} with JSON naming the template ("template") and its data ("data"). A compiled template (.bdft), or Excel, Word, PowerPoint, a PDF form, pdfme or text.`));
  if (n.accept === "file") box.append(el("p", "This endpoint takes documents only: its templates are not filled until Accepts changes.", { className: "text-xs text-warning" }));
  const ts = await api("GET", `/api/nodes/${n.name}/templates`);
  const table = el("table", undefined, { className: "table table-xs" });
  const tbody = el("tbody");
  for (const t of ts) {
    const tr = el("tr");
    const fields = (t.fields || []).map((f) => f.name).join(", ");
    tr.append(el("td", t.id, { className: "font-mono font-semibold" }), el("td", t.format), el("td", `${(t.fields || []).length} fields`, { title: fields }), el("td", `${Math.max(1, Math.round(t.size / 1024))} KiB`), el("td"));
    tr.lastChild.append(el("a", "File", { className: "btn btn-xs btn-ghost", href: `/api/nodes/${n.name}/templates/${t.id}/file` }), button("Delete", () => api("DELETE", `/api/nodes/${n.name}/templates/${t.id}`), "btn btn-xs btn-ghost text-error"));
    tbody.append(tr);
  }
  if (!ts.length) {
    const tr = el("tr");
    tr.append(el("td", "No templates yet", { colSpan: 5, className: "text-base-content/50" }));
    tbody.append(tr);
  }
  table.append(tbody);
  const wrap = el("div", undefined, { className: "overflow-x-auto" });
  wrap.append(table);
  box.append(wrap);
  const form = el("form", undefined, { className: "flex flex-wrap gap-1" });
  const id = el("input", undefined, { className: "input input-sm w-28", placeholder: "id (invoice)", required: true, pattern: "[A-Za-z0-9][A-Za-z0-9_\\-]{0,63}" });
  const file = el("input", undefined, { type: "file", className: "file-input file-input-sm min-w-0 flex-1", required: true });
  form.append(id, file, el("button", "Upload", { type: "submit", className: "btn btn-sm" }));
  form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    const f = file.files[0];
    if (!f) return;
    try {
      await api("PUT", `/api/nodes/${n.name}/templates/${encodeURIComponent(id.value)}?name=${encodeURIComponent(f.name)}`, f, true);
      status(`${id.value} uploaded`);
      changed();
    } catch (e) {
      status(e.message, true);
    }
  });
  box.append(form);
}

async function devicesSection(n, box) {
  box.append(hint(`${n.open || 0} open slip${n.open === 1 ? "" : "s"} on the board. Only the devices enrolled here see its documents; a device keeps a key of its own.`));
  const ds = await api("GET", `/api/nodes/${n.name}/devices`);
  const table = el("table", undefined, { className: "table table-xs" });
  const tbody = el("tbody");
  for (const d of ds) {
    const tr = el("tr");
    tr.append(el("td", d.name, { className: "font-semibold" }), el("td", "seen " + since(d.last_seen), { className: "text-base-content/60" }), el("td"));
    tr.lastChild.append(button("Remove", () => api("DELETE", "/api/devices/" + d.id), "btn btn-xs btn-ghost text-error"));
    tbody.append(tr);
  }
  if (!ds.length) {
    const tr = el("tr");
    tr.append(el("td", "No devices yet", { colSpan: 3, className: "text-base-content/50" }));
    tbody.append(tr);
  }
  table.append(tbody);
  box.append(table);
  const add = el("button", "Add a device (QR code)", { type: "button", className: "btn btn-sm btn-primary w-fit" });
  add.addEventListener("click", async () => {
    try {
      const e = await api("POST", `/api/nodes/${n.name}/enrollments`);
      $("qr").innerHTML = e.qr; // SVG the manager drew
      $("enroll-code").textContent = e.code.replace(/(.{5})/, "$1-");
      $("enroll-url").textContent = e.url;
      $("enroll-url").href = e.url;
      $("enroll-expires").textContent = `Valid once, until ${since(e.expires)}`;
      $("enroll-modal").showModal();
    } catch (err) {
      status(err.message, true);
    }
  });
  box.append(add);
}

function showKey(name, key) {
  $("key-name").textContent = name;
  $("key-value").textContent = key;
  $("key-usage").textContent = `curl -X POST -H "Authorization: Bearer ${key}" --data-binary @document.pdf "${location.origin}/in/${name}?name=document"`;
  $("key-modal").showModal();
}

document.addEventListener("keydown", (ev) => {
  const typing = ["INPUT", "TEXTAREA", "SELECT"].includes(document.activeElement?.tagName);
  if (ev.key === "Escape" && draft) {
    draft = null;
    draftPos = null;
    draw();
    renderSidebar();
  } else if (ev.key === "Delete" && selected && !typing) {
    remove(selected);
  }
});

// a link from another page names a node: /admin/topology#<name>
if (location.hash.length > 1) selected = decodeURIComponent(location.hash.slice(1));
onChange(load);
load();

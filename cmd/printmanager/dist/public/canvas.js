// The canvas of the topology, shared by the editor (topology.js) and
// everyone's view (topology-view.js): nodes as boxes with what is at them,
// edges as arrows, panned by dragging the background and zoomed by the
// wheel or a pinch. In the editor, dragging a node moves it and dragging
// from its right port onto another node connects them; a click selects.
// The flows of the last minutes (jobs going from a node to the next) run
// as dots along the edges.

import { KINDS } from "./admin.js";

export const NODE_W = 176;
export const NODE_H = 54;
const COL = 240;
const ROW = 84;
const PAD = 24;
const GRID = 8;
const SVG = "http://www.w3.org/2000/svg";
const FILL = { endpoint: "#fde68a", cups: "#fecaca", queue: "#ddd6fe", class: "#c7d2fe", filter: "#fbcfe8", inbox: "#fed7aa", archive: "#e5e7eb", display: "#bfdbfe", printer: "#d1d5db", sink: "#bbf7d0" };
// the kinds with an output port: they send to other nodes (a display to
// its fallback printer)
const OUT_PORT = new Set(["endpoint", "queue", "class", "filter", "inbox", "display"]);
// a flow older than this is not run (the page was away)
const FLOW_AGE_MS = 2 * 60 * 1000;

function svgEl(tag, attrs = {}) {
  const e = document.createElementNS(SVG, tag);
  for (const [k, v] of Object.entries(attrs)) e.setAttribute(k, v);
  return e;
}

// edgesOf are the edges from a node: to (a queue's, an endpoint's, an
// inbox's, a filter's rules and otherwise), member (a class's) or
// fallback (a display's printer).
export function edgesOf(n) {
  const out = [];
  const seen = new Set();
  const to = (name, kind) => {
    if (seen.has(kind + ":" + name)) return;
    seen.add(kind + ":" + name);
    out.push({ to: name, kind });
  };
  for (const x of n.to || []) to(x, "to");
  for (const r of n.rules || []) for (const x of r.to || []) to(x, "to");
  for (const x of n.otherwise || []) to(x, "to");
  for (const x of n.members || []) to(x, "member");
  if (n.fallback && n.fallback.printer) to(n.fallback.printer, "fallback");
  return out;
}

function clipText(s, n) {
  return s.length > n ? s.slice(0, n - 1) + "…" : s;
}

// createCanvas draws a topology in an svg. opts: editable (nodes move
// and connect), onSelect(name or null), onMove(name, pos), onConnect(node,
// toName), onRemoveEdge(node, edge), draftLabel() (the name typed for a
// node being added).
export function createCanvas(svg, opts = {}) {
  const { editable = false, onSelect = () => {}, onMove = () => {}, onConnect = () => {}, onRemoveEdge = () => {}, draftLabel = () => "" } = opts;
  let topo = []; // the nodes, as /api/topology lists them
  let layout = {}; // name -> {x, y}
  let live = null; // /api/topology/live
  let selected = null;
  let draft = null; // {kind} of a node being added
  let draftPos = null;
  let view = { x: PAD, y: PAD, k: 1 };
  let world = null; // the <g> of the world
  let flowsLayer = null;
  const shown = new Set(); // the flows run already
  let drag = null;
  const pointers = new Map();
  let pinch = null;

  const byName = () => Object.fromEntries(topo.map((n) => [n.name, n]));

  // autoLayout places the nodes that have no position: in columns by
  // their depth from the entries, rows within a column.
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

  const pos = (name) => layout[name] || { x: 0, y: 0 };
  const snap = (v) => Math.round(v / GRID) * GRID;

  // freeSpot is a place for a new node near a point, clear of the
  // others: to the right of them, then on the next row.
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

  // counts are what is at a node: the live activity when there is some
  // (the viewer's), else what the topology says (an administrator's).
  function counts(n) {
    const l = live ? live.nodes[n.name] || {} : n;
    const out = [];
    if (l.held) out.push(`${l.held} held`);
    if (l.open) out.push(`${l.open} open`);
    if (l.pending) out.push(`${l.pending} waiting`);
    if (l.failed) out.push(`${l.failed} failed`);
    if (l.kept) out.push(`${l.kept} kept`);
    if (!live && n.devices) out.push(`${n.devices} device${n.devices === 1 ? "" : "s"}`);
    if (live && l.recent) out.push(`${l.recent}/h`);
    return out;
  }

  function edgePath(from, to) {
    const a = pos(from);
    const b = pos(to);
    const x1 = a.x + NODE_W;
    const y1 = a.y + NODE_H / 2;
    const x2 = b.x;
    const y2 = b.y + NODE_H / 2;
    const dx = Math.max(40, Math.abs(x2 - x1) / 2);
    return `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`;
  }

  function draw() {
    svg.replaceChildren();
    const defs = svgEl("defs");
    const marker = svgEl("marker", { id: "arrow", viewBox: "0 0 10 10", refX: "10", refY: "5", markerWidth: "8", markerHeight: "8", orient: "auto-start-reverse" });
    marker.append(svgEl("path", { d: "M 0 0 L 10 5 L 0 10 z", fill: "#64748b" }));
    defs.append(marker);
    world = svgEl("g", { id: "world", transform: `translate(${view.x} ${view.y}) scale(${view.k})` });
    const nodes = byName();
    for (const n of topo) {
      for (const e of edgesOf(n)) {
        if (!nodes[e.to]) continue;
        const d = edgePath(n.name, e.to);
        const hit = svgEl("path", { d, fill: "none", stroke: "transparent", "stroke-width": "14", class: editable ? "cursor-pointer" : "" });
        const line = svgEl("path", { d, fill: "none", stroke: e.kind === "fallback" ? "#94a3b8" : "#64748b", "stroke-width": "2", "marker-end": "url(#arrow)", "pointer-events": "none" });
        if (e.kind === "member") line.setAttribute("stroke-dasharray", "6 4");
        if (e.kind === "fallback") line.setAttribute("stroke-dasharray", "2 4");
        const title = svgEl("title");
        const what = e.kind === "fallback" ? `${n.name} prints on ${e.to} as a fallback` : e.kind === "member" ? `${e.to} is a member of ${n.name}` : `${n.name} → ${e.to}`;
        title.textContent = editable ? what + " (click to remove)" : what;
        hit.append(title);
        if (editable) {
          hit.addEventListener("pointerdown", (ev) => ev.stopPropagation());
          hit.addEventListener("click", () => onRemoveEdge(n, e));
        }
        world.append(hit, line);
      }
    }
    for (const n of topo) world.append(nodeBox(n, pos(n.name), n.name === selected));
    if (draft && draftPos) {
      const g = svgEl("g", { transform: `translate(${draftPos.x} ${draftPos.y})`, class: "cursor-grab", "data-draft": "1" });
      g.append(svgEl("rect", { width: NODE_W, height: NODE_H, rx: "8", fill: FILL[draft.kind] || "#e5e7eb", stroke: "#4f46e5", "stroke-width": "2", "stroke-dasharray": "6 4" }));
      const t = svgEl("text", { x: "10", y: "22", "font-size": "13", "font-weight": "600", fill: "#111827" });
      t.textContent = draftLabel() || "new " + (KINDS[draft.kind] || draft.kind).toLowerCase();
      const k = svgEl("text", { x: "10", y: "40", "font-size": "11", fill: "#374151" });
      k.textContent = (KINDS[draft.kind] || draft.kind) + " · not saved yet";
      g.append(t, k);
      g.addEventListener("pointerdown", (ev) => startDrag(ev, null));
      world.append(g);
    }
    flowsLayer = svgEl("g", { id: "flows", "pointer-events": "none" });
    world.append(flowsLayer);
    svg.append(defs, world);
  }

  // nodeBox draws a node: its title, its kind and what is at it, an input
  // port at the left and, in the editor, an output port at the right for
  // the kinds that send on.
  function nodeBox(n, p, isSelected) {
    const g = svgEl("g", { transform: `translate(${p.x} ${p.y})`, class: editable ? "cursor-grab" : "cursor-pointer", "data-name": n.name });
    g.append(svgEl("rect", { width: NODE_W, height: NODE_H, rx: "8", fill: FILL[n.kind] || "#e5e7eb", stroke: isSelected ? "#4f46e5" : "#94a3b8", "stroke-width": isSelected ? "3" : "1.5" }));
    // a sink shows the first page of the latest document it keeps
    const latest = n.kind === "sink" ? (live ? (live.nodes[n.name] || {}).latest : n.latest) : "";
    const name = svgEl("text", { x: "10", y: "21", "font-size": "13", "font-weight": "600", fill: "#111827" });
    name.textContent = clipText(n.title && n.title !== n.name ? n.title : n.name, latest ? 16 : 22);
    const sub = svgEl("text", { x: "10", y: "39", "font-size": "11", fill: "#374151" });
    const parts = [KINDS[n.kind] || n.kind];
    if (n.title && n.title !== n.name) parts.push(n.name);
    parts.push(...counts(n));
    sub.textContent = clipText(parts.join(" · "), latest ? 22 : 30);
    const title = svgEl("title");
    title.textContent = `${n.name} (${KINDS[n.kind] || n.kind})${n.title ? ": " + n.title : ""}${counts(n).length ? " · " + counts(n).join(", ") : ""}`;
    g.append(title, name, sub);
    const c = live ? live.nodes[n.name] : null;
    if (c && (c.held || c.open || c.pending || c.failed)) {
      // a badge of what waits at the node
      const count = (c.held || 0) + (c.open || 0) + (c.pending || 0) + (c.failed || 0);
      const badge = svgEl("g", { transform: `translate(${NODE_W - 14} -8)` });
      badge.append(svgEl("circle", { r: "10", fill: c.failed ? "#dc2626" : "#f59e0b", stroke: "#fff", "stroke-width": "2" }));
      const t = svgEl("text", { "text-anchor": "middle", y: "4", "font-size": "11", "font-weight": "700", fill: "#fff" });
      t.textContent = count > 99 ? "99+" : String(count);
      badge.append(t);
      g.append(badge);
    }
    if (latest) {
      g.append(svgEl("rect", { x: NODE_W - 48, y: 6, width: 40, height: 42, rx: "3", fill: "#fff", stroke: "#94a3b8" }));
      g.append(svgEl("image", { href: `/api/jobs/${latest}/pages/1.png?w=96`, x: NODE_W - 47, y: 7, width: 38, height: 40, preserveAspectRatio: "xMidYMin slice" }));
    }
    // a CUPS server is an entry: nothing is sent to it from the network
    if (n.kind !== "cups") g.append(svgEl("circle", { cx: "0", cy: NODE_H / 2, r: "5", fill: "#f8fafc", stroke: "#64748b", "stroke-width": "1.5" }));
    if (editable && OUT_PORT.has(n.kind)) {
      const port = svgEl("circle", { cx: NODE_W, cy: NODE_H / 2, r: "7", fill: "#4f46e5", stroke: "#fff", "stroke-width": "2", class: "cursor-crosshair", "data-port": n.name });
      const t = svgEl("title");
      t.textContent = n.kind === "display" ? "Drag onto a printer: the fallback printer" : n.kind === "class" ? "Drag onto a node: a member" : "Drag onto a node to connect";
      port.append(t);
      g.append(port);
    }
    g.addEventListener("pointerdown", (ev) => startDrag(ev, n));
    return g;
  }

  // runFlows runs the flows not run yet as dots along their edges.
  function runFlows() {
    if (!live || !flowsLayer) return;
    const nodes = byName();
    const now = Date.now();
    for (const f of live.flows) {
      const key = `${f.job}:${f.from}>${f.to}`;
      if (shown.has(key)) continue;
      shown.add(key);
      if (!nodes[f.from] || !nodes[f.to] || now - new Date(f.at).getTime() > FLOW_AGE_MS) continue;
      const dot = svgEl("circle", { r: "6", fill: f.state === "error" || f.state === "failed" ? "#dc2626" : "#4f46e5", stroke: "#fff", "stroke-width": "2" });
      const title = svgEl("title");
      title.textContent = `#${f.job} ${f.name || ""} ${f.from} → ${f.to}`;
      dot.append(title);
      const move = svgEl("animateMotion", { dur: "1.6s", fill: "freeze", path: edgePath(f.from, f.to) });
      dot.append(move);
      flowsLayer.append(dot);
      setTimeout(() => dot.remove(), 2200);
    }
    if (shown.size > 2000) shown.clear();
  }

  // Pan and zoom: the world is translated by view.x, view.y and scaled
  // by view.k. Screen points are relative to the canvas.

  function screenPoint(ev) {
    const r = svg.getBoundingClientRect();
    return { x: ev.clientX - r.left, y: ev.clientY - r.top };
  }

  const toWorld = (s) => ({ x: (s.x - view.x) / view.k, y: (s.y - view.y) / view.k });

  function zoomAt(s, factor) {
    const k = Math.min(3, Math.max(0.25, view.k * factor));
    view = { x: s.x - (s.x - view.x) * (k / view.k), y: s.y - (s.y - view.y) * (k / view.k), k };
    draw();
  }

  function size() {
    const r = svg.getBoundingClientRect();
    return { w: r.width || 800, h: r.height || 400 };
  }

  const center = () => {
    const { w, h } = size();
    return toWorld({ x: w / 2, y: h / 2 });
  };

  // fit shows every node.
  function fit() {
    const { w, h } = size();
    if (!topo.length) {
      view = { x: PAD, y: PAD, k: 1 };
      draw();
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
  // port draws an edge, two pointers pinch. Not editable, a node is only
  // selected.
  function startDrag(ev, n) {
    if (ev.button !== 0) return;
    ev.preventDefault();
    ev.stopPropagation();
    const s = screenPoint(ev);
    if (editable && n && ev.target.dataset.port) {
      const p = pos(n.name);
      const line = svgEl("path", { fill: "none", stroke: "#4f46e5", "stroke-width": "2", "stroke-dasharray": "4 3", "pointer-events": "none" });
      world.append(line);
      drag = { kind: "edge", from: n, line, x1: p.x + NODE_W, y1: p.y + NODE_H / 2 };
    } else if (editable || !n) {
      const p = n ? pos(n.name) : draftPos;
      const w = toWorld(s);
      drag = { kind: "move", node: n, dx: w.x - p.x, dy: w.y - p.y, moved: false };
    } else {
      drag = { kind: "click", node: n };
    }
    svg.setPointerCapture(ev.pointerId);
  }

  svg.addEventListener("pointerdown", (ev) => {
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
    svg.setPointerCapture(ev.pointerId);
  });

  svg.addEventListener("pointermove", (ev) => {
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
    } else if (drag.kind === "edge") {
      const w = toWorld(s);
      drag.line.setAttribute("d", `M ${drag.x1} ${drag.y1} L ${w.x} ${w.y}`);
    }
  });

  function endDrag(ev) {
    pointers.delete(ev.pointerId);
    if (pointers.size < 2) pinch = null;
    if (!drag) return;
    const d = drag;
    drag = null;
    switch (d.kind) {
      case "pan":
        if (!d.moved) {
          selected = null;
          draw();
          onSelect(null);
        }
        return;
      case "click":
        selected = d.node.name;
        draw();
        onSelect(d.node.name);
        return;
      case "move":
        if (!d.node) return;
        if (d.moved) onMove(d.node.name, layout[d.node.name]);
        else {
          selected = d.node.name;
          draw();
          onSelect(d.node.name);
        }
        return;
      case "edge": {
        d.line.remove();
        const target = document.elementFromPoint(ev.clientX, ev.clientY)?.closest("g[data-name]");
        if (target) onConnect(d.from, target.dataset.name);
      }
    }
  }
  svg.addEventListener("pointerup", endDrag);
  svg.addEventListener("pointercancel", endDrag);
  svg.addEventListener(
    "wheel",
    (ev) => {
      ev.preventDefault();
      zoomAt(screenPoint(ev), Math.exp(-ev.deltaY * 0.0015));
    },
    { passive: false },
  );
  window.addEventListener("resize", draw);

  return {
    // setData takes the nodes and their positions, and places the nodes
    // without one.
    setData(t, l) {
      topo = t;
      layout = l || {};
      autoLayout();
      if (selected && !topo.some((n) => n.name === selected)) selected = null;
      draw();
    },
    // setLive takes the activity: the counts, and the flows to run.
    setLive(l) {
      live = l;
      draw();
      runFlows();
    },
    setSelected(name) {
      selected = name;
      draft = null;
      draftPos = null;
      draw();
    },
    setDraft(d, p) {
      draft = d;
      draftPos = p;
      selected = null;
      draw();
    },
    draftPos: () => draftPos,
    pos,
    freeSpot,
    center,
    draw,
    fit,
    zoomBy(f) {
      const { w, h } = size();
      zoomAt({ x: w / 2, y: h / 2 }, f);
    },
  };
}

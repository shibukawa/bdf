// The canvas of the topology, shared by the editor (topology.js) and
// everyone's view (topology.js, looking): nodes as boxes with what is at them,
// edges as arrows, panned by dragging the background and zoomed by the
// wheel or a pinch. In the editor, dragging a node moves it and dragging
// from its right port onto another node connects them; a click selects.
// The flows of the last minutes (jobs going from a node to the next) run
// as dots along the edges.
//
// Every node is a white box with a title bar colored by its role (where
// a job comes in, what passes it on, what picks, what makes a new
// payload, where it ends) and an icon of its kind. The graph is typed
// (kinds.go): what flows along an edge is a .bdf document, data or any
// file as received, and the ports and the edges are colored by it, as
// the pins and the wires of a visual editor are. A node that passes on
// what reaches it (a
// queue, a class, a filter, an inbox) is typed by its first edge, on
// either side: a queue in front of a printer carries .bdf, and so says
// its input port; one with no edge yet is not typed (hollow, dashed).
// While an edge is dragged, the nodes that would not take what it
// carries fade, the wire turns red over one, and it snaps to the port of
// one that takes it. An edge that no longer carries what the node after
// it takes (after a change of settings) is drawn broken: dashed, red,
// with a mark on it.

import { KINDS, PAYLOAD_COLOR, PAYLOADS } from "./admin.js";

export const NODE_W = 176;
export const NODE_H = 54;
const COL = 240;
const ROW = 84;
const PAD = 24;
const GRID = 8;
const SVG = "http://www.w3.org/2000/svg";
const BASE_FILL = "#ffffff";
const BORDER = "#cbd5e1";
const HEAD_H = 22; // the title bar
const TAG_Y = 37; // the baseline of the [kind] line
const DESC_Y = 52; // the baseline of the first line of the description
const LINE_H = 13;
const MAX_LINES = 5; // of a description on the canvas
const GAP = 30; // between the rows the nodes are laid out in
// The roles of the kinds, which color a node's title bar: where a job
// comes in from outside (source), what passes it on (control flow),
// what picks by rules (filter), what makes a new payload of it
// (transform), and where it ends (sink).
export const ROLES = { source: "Source", flow: "Control flow", filter: "Filter", transform: "Transform", sink: "Sink" };
export const ROLE_OF = { endpoint: "source", cups: "source", queue: "flow", class: "flow", inbox: "flow", filter: "filter", template: "transform", convert: "transform", display: "sink", printer: "sink", archive: "sink", sink: "sink" };
export const ROLE_COLOR = { source: "#0f766e", flow: "#4338ca", filter: "#be185d", transform: "#c2410c", sink: "#334155" };
// The icons of the kinds: strokes in a 24 by 24 box.
const ICONS = {
  endpoint: ["M14 4h4a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-4", "M3 12h10", "M9 8l4 4-4 4"],
  cups: ["M4 5h16v6H4z", "M4 13h16v6H4z", "M7 8h.01", "M7 16h.01", "M15 8h2", "M15 16h2"],
  queue: ["M12 3l9 5-9 5-9-5z", "M3 13l9 5 9-5", "M3 17l9 5 9-5"],
  class: ["M17 3l4 4-4 4", "M3 11V9a2 2 0 0 1 2-2h16", "M7 21l-4-4 4-4", "M21 13v2a2 2 0 0 1-2 2H3"],
  filter: ["M4 4h16l-6 8v6l-4 2v-8z"],
  template: ["M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z", "M14 3v5h5", "M9 13h6", "M9 17h4"],
  convert: ["M21 12a9 9 0 0 1-15 6.7L3 16", "M3 12a9 9 0 0 1 15-6.7L21 8", "M3 21v-5h5", "M21 3v5h-5"],
  inbox: ["M5.5 5h13l3.5 7v7H2v-7z", "M2 12h5l2 3h6l2-3h5"],
  display: ["M3 4h18v12H3z", "M8 21h8", "M12 16v5"],
  printer: ["M6 9V3h12v6", "M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2", "M6 14h12v7H6z"],
  archive: ["M3 4h18v4H3z", "M5 8v12h14V8", "M10 12h4"],
  sink: ["M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7z", "M12 9a3 3 0 1 0 0 6 3 3 0 0 0 0-6z"],
};
// the kinds with an output port: they send to other nodes (a display to
// its fallback printer)
const OUT_PORT = new Set(["endpoint", "cups", "queue", "class", "filter", "template", "convert", "inbox", "display"]);
// the entries: their input is from outside the network (an endpoint's
// REST, a CUPS server's IPP port), drawn as an arrow coming in
const ENTRY = new Set(["endpoint", "cups"]);
// the kinds of node that send on what reaches them: typed by their first
// edge, on either side (kinds.go: passesOn)
const PASSES_ON = new Set(["queue", "class", "filter", "inbox", "endpoint", "cups"]);
// the kinds of node that send on a .bdf whatever reaches them: a template
// or a convert node the one it makes, a display the one its fallback
// printer prints
const MAKES_BDF = new Set(["template", "convert", "display"]);
const ALL = ["bdf", "data", "any"];
// the colors besides the kinds' (admin.js: PAYLOAD_COLOR): several kinds
// at once, a port or an edge not typed yet, an input that takes any kind
// for good, and an edge that no longer connects
export const MIXED_COLOR = "#7c3aed";
export const UNTYPED_COLOR = "#94a3b8";
export const ANY_KIND_COLOR = "#475569";
export const BROKEN_COLOR = "#dc2626";
// a flow older than this is not run (the page was away)
const FLOW_AGE_MS = 2 * 60 * 1000;

function svgEl(tag, attrs = {}) {
  const e = document.createElementNS(SVG, tag);
  for (const [k, v] of Object.entries(attrs)) e.setAttribute(k, v);
  return e;
}

// iconOf is the icon of a kind of node, as an svg of a size, drawn in
// the current color (color: the role's, white on a title bar).
export function iconOf(kind, size = 14) {
  const svg = svgEl("svg", { width: size, height: size, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", "stroke-width": "2", "stroke-linecap": "round", "stroke-linejoin": "round", "aria-hidden": "true" });
  for (const d of ICONS[kind] || ICONS.queue) svg.append(svgEl("path", { d }));
  return svg;
}

export const roleColor = (kind) => ROLE_COLOR[ROLE_OF[kind] || "flow"];

// edgesOf are the edges from a node: to (a queue's, an endpoint's, a CUPS
// server's, an inbox's, a filter's rules and otherwise), member (a
// class's) or fallback (a display's printer).
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

const asList = (kinds) => (Array.isArray(kinds) ? kinds : kinds ? [kinds] : []);
// narrowed tells kinds that are some, not every one
const narrowed = (list) => list.length > 0 && list.length < ALL.length;

// payloadColor is the color of the kinds of payload a port or an edge
// carries (takes, accepts and emits of /api/topology, kinds.go): one
// kind its color, several purple, none grey (not typed yet).
export function payloadColor(kinds) {
  const list = asList(kinds);
  if (list.length === 1) return PAYLOAD_COLOR[list[0]] || PAYLOAD_COLOR.any;
  if (list.length > 1) return MIXED_COLOR;
  return UNTYPED_COLOR;
}

export function payloadText(kinds) {
  const list = asList(kinds);
  return list.map((k) => PAYLOADS[k] || k).join(" or ");
}

// typeOf tells what a node is typed with: in, the kinds its input port is
// typed with (none: not typed yet, or any kind); out, the kinds its
// output carries (none: not typed yet); accepts, the kinds its input may
// take now (what it takes, within what the nodes after it take); anyKind
// when its input takes every kind for good (a convert node, an archive, a
// sink); dead when nothing fits its input (the nodes after it take
// different kinds). A node that passes on what reaches it is typed by
// what reaches it, else by what the nodes after it take: its first edge.
export function typeOf(n) {
  const takes = n.takes && n.takes !== "any" ? [n.takes] : ALL;
  const accepts = n.accepts ? asList(n.accepts) : takes;
  const emits = asList(n.emits);
  const passes = PASSES_ON.has(n.kind);
  let inn;
  if (ENTRY.has(n.kind)) inn = emits; // what it sends on is what it is: a data endpoint, a document endpoint
  else if (passes) inn = emits.length ? emits : narrowed(accepts) ? accepts : [];
  else inn = narrowed(takes) ? takes : [];
  const out = MAKES_BDF.has(n.kind) ? ["bdf"] : emits.length ? emits : passes && narrowed(accepts) ? accepts : [];
  return { in: inn, out, accepts, anyKind: !passes && takes.length === ALL.length, dead: accepts.length === 0 };
}

// mend says what to put between two nodes whose kinds disagree.
function mend(out, accepts) {
  if (out.includes("any") && accepts.includes("bdf")) return "put a convert node between";
  if (out.includes("data") && accepts.includes("bdf")) return "put a template node between";
  if (out.includes("bdf") && accepts.includes("data")) return "a template node takes data: send documents elsewhere";
  return "";
}

// legend fills an element with what the colors mean.
export function legend(box) {
  const item = (pic, text) => {
    const s = document.createElement("span");
    s.className = "inline-flex items-center gap-1 whitespace-nowrap";
    s.append(pic, text);
    return s;
  };
  const dot = (fill, stroke, dash) => {
    const s = svgEl("svg", { width: "14", height: "14", viewBox: "0 0 14 14" });
    const c = svgEl("circle", { cx: "7", cy: "7", r: "5", fill, stroke, "stroke-width": "1.5" });
    if (dash) c.setAttribute("stroke-dasharray", dash);
    s.append(c);
    return s;
  };
  const line = (color, dash) => {
    const s = svgEl("svg", { width: "26", height: "14", viewBox: "0 0 26 14" });
    s.append(svgEl("path", { d: "M 1 7 L 25 7", stroke: color, "stroke-width": "2", "stroke-dasharray": dash }));
    return s;
  };
  const bar = (fill) => {
    const s = svgEl("svg", { width: "14", height: "14", viewBox: "0 0 14 14" });
    s.append(svgEl("rect", { x: "1", y: "2", width: "12", height: "10", rx: "2", fill }));
    return s;
  };
  const of = { source: "endpoint, CUPS server", flow: "queue, class, inbox", filter: "by rules", transform: "template, convert", sink: "display, printer, archive, sink" };
  box.replaceChildren(...Object.entries(ROLES).map(([role, label]) => item(bar(ROLE_COLOR[role]), `${label} (${of[role]})`)), item(dot(PAYLOAD_COLOR.bdf, "#fff"), ".bdf document"), item(dot(PAYLOAD_COLOR.data, "#fff"), "data (for a template node)"), item(dot(PAYLOAD_COLOR.any, "#fff"), "any file as received (for a convert node)"), item(dot(MIXED_COLOR, "#fff"), "several kinds"), item(dot(ANY_KIND_COLOR, "#fff"), "takes any kind"), item(dot("#f8fafc", UNTYPED_COLOR, "2 2"), "not typed yet: its first edge types it"), item(line(BROKEN_COLOR, "6 4"), "no longer connects: the node after does not take what reaches it"));
}

// createCanvas draws a topology in an svg. opts: editable (nodes move
// and connect), onSelect(name or null), onMove(name, pos), onConnect(node,
// toName), onReject(node, toName, why) (an edge dropped on a node that
// cannot take it), onRemoveEdge(node, edge), draftLabel() and
// draftNote() (the name and the description typed for a node being
// added).
export function createCanvas(svg, opts = {}) {
  const { editable = false, onSelect = () => {}, onMove = () => {}, onConnect = () => {}, onReject = () => {}, onRemoveEdge = () => {}, draftLabel = () => "", draftNote = () => "" } = opts;
  let topo = []; // the nodes, as /api/topology lists them
  let layout = {}; // name -> {x, y}
  let live = null; // /api/topology/live
  let selected = null;
  let draft = null; // {kind} of a node being added
  let draftPos = null;
  let view = { x: PAD, y: PAD, k: 1 };
  let world = null; // the <g> of the world
  let defs = null;
  let flowsLayer = null;
  const shown = new Set(); // the flows run already
  let drag = null;
  const pointers = new Map();
  let pinch = null;

  const byName = () => Object.fromEntries(topo.map((n) => [n.name, n]));

  // A node's description is written in the node, as the description of
  // a C4 element is, wrapped over the lines of the node (the first
  // MAX_LINES of them): the node is as tall as they need. The text is
  // measured at the canvas's font with a 2D context.
  let ctx2d = null;
  let family = "";
  function textWidth(s, size) {
    if (ctx2d === null) {
      try {
        ctx2d = document.createElement("canvas").getContext("2d") || false;
        family = getComputedStyle(svg).fontFamily || "sans-serif";
      } catch {
        ctx2d = false;
      }
    }
    if (!ctx2d) return s.length * size * 0.55;
    ctx2d.font = `${size}px ${family}`;
    return ctx2d.measureText(s).width;
  }

  // wrap breaks a text into lines no wider than width: at a space when
  // there is one, else between a Japanese character and its neighbour
  // (a word in Latin letters stays whole, and a line does not start
  // with a closing mark), else anywhere (a long word); at most max
  // lines, the last ending in … when there was more.
  const wrapped = new Map();
  const isCJK = (ch) => /[\u3000-\u30ff\u3400-\u9fff\uf900-\ufaff\uff00-\uffef]/.test(ch || "");
  const noLineStart = "、。，．）」』】〕〉》〙〗！？：；ーぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ";
  function wrap(text, width, max) {
    const key = `${width}|${text}`;
    if (wrapped.has(key)) return wrapped.get(key);
    const lines = [];
    let more = false;
    const paras = text
      .split(/\r?\n/)
      .map((p) => p.replace(/\s+/g, " ").trim())
      .filter(Boolean);
    for (const s of paras) {
      let i = 0;
      while (i < s.length) {
        if (lines.length >= max) {
          more = true;
          break;
        }
        let n = 1;
        while (i + n < s.length && textWidth(s.slice(i, i + n + 1), 11) <= width) n++;
        let cut = n;
        if (i + n < s.length) {
          const sp = s.slice(i, i + n).lastIndexOf(" ");
          if (sp > 0) cut = sp + 1;
          else {
            for (let k = n; k > 0; k--) {
              if ((isCJK(s[i + k - 1]) || isCJK(s[i + k])) && !noLineStart.includes(s[i + k])) {
                cut = k;
                break;
              }
            }
          }
        }
        lines.push(s.slice(i, i + cut).trimEnd());
        i += cut;
        while (s[i] === " ") i++;
      }
      if (more) break;
    }
    if (more && lines.length) {
      let last = lines[lines.length - 1];
      while (last.length && textWidth(last + "…", 11) > width) last = last.slice(0, -1);
      lines[lines.length - 1] = last + "…";
    }
    if (wrapped.size > 500) wrapped.clear();
    wrapped.set(key, lines);
    return lines;
  }

  // latestOf is the latest document a sink keeps, whose first page it
  // shows.
  const latestOf = (n) => (n.kind === "sink" ? (live ? (live.nodes[n.name] || {}).latest : n.latest) || "" : "");

  // linesOf are the lines of a node's description on the canvas.
  const linesOf = (n) => (n.description ? wrap(n.description, NODE_W - 20 - (latestOf(n) ? 44 : 0), MAX_LINES) : []);
  const heightOfLines = (count) => (count ? DESC_Y + (count - 1) * LINE_H + 11 : NODE_H);
  const heightOf = (n) => heightOfLines(linesOf(n).length);
  let heights = {}; // by name, as drawn last

  // autoLayout places the nodes that have no position: in columns by
  // their depth from the entries, rows within a column, each as tall as
  // its node.
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
    const placed = () => Object.entries(layout).map(([name, p]) => ({ x: p.x, y: p.y, h: heightOf(nodes[name] || {}) }));
    for (const [d, names] of Object.entries(columns)) {
      const x = d * COL;
      let y = 0;
      for (const name of names) {
        if (layout[name]) continue;
        const h = heightOf(nodes[name]);
        const hits = (y) => placed().some((q) => Math.abs(q.x - x) < NODE_W + GRID && y < q.y + q.h + GRID && y + h + GRID > q.y);
        while (hits(y)) y += GRID;
        layout[name] = { x, y };
        y += h + GAP;
      }
    }
  }

  const pos = (name) => layout[name] || { x: 0, y: 0 };
  const snap = (v) => Math.round(v / GRID) * GRID;

  // freeSpot is a place for a new node near a point, clear of the
  // others: to the right of them, then on the next row.
  function freeSpot(c) {
    const taken = (p) => topo.some((n) => Math.abs(pos(n.name).x - p.x) < NODE_W + GRID && p.y < pos(n.name).y + heightOf(n) + GRID && p.y + NODE_H + GRID > pos(n.name).y);
    const start = { x: snap(c.x - NODE_W / 2), y: snap(c.y - NODE_H / 2) };
    let p = { ...start };
    for (let i = 0; i < 64 && taken(p); i++) {
      p = { x: p.x + COL, y: p.y };
      if (p.x > start.x + 3 * COL) p = { x: start.x, y: p.y + ROW };
    }
    return p;
  }

  // reaches tells a path from a node to another along the edges.
  function reaches(from, to) {
    const nodes = byName();
    const seen = new Set([from]);
    const stack = [from];
    while (stack.length) {
      const name = stack.pop();
      for (const e of edgesOf(nodes[name] || {})) {
        if (e.to === to) return true;
        if (!seen.has(e.to)) {
          seen.add(e.to);
          stack.push(e.to);
        }
      }
    }
    return false;
  }

  // canConnect tells whether an edge from a node to another may be made:
  // not to itself, not to a CUPS server (its input is computers'), a
  // display's only to a printer (its fallback), not twice (already), no
  // path back (a cycle), and what the edge would carry is taken by the
  // node after it (the server checks the same when the node is put).
  function canConnect(a, b) {
    if (a.name === b.name) return { ok: false, why: "a node does not send to itself" };
    if (b.kind === "cups") return { ok: false, why: `${b.name} is a CUPS server: it takes what computers print there, not what other nodes send` };
    if (a.kind === "display" && b.kind !== "printer") return { ok: false, why: `a display's port sets its fallback printer; ${b.name} is a ${(KINDS[b.kind] || b.kind).toLowerCase()}` };
    if (edgesOf(a).some((e) => e.to === b.name)) return { ok: false, already: true, why: `${a.name} sends to ${b.name} already` };
    if (reaches(b.name, a.name)) return { ok: false, why: `a cycle: ${b.name} leads back to ${a.name}` };
    const out = typeOf(a).out;
    const tb = typeOf(b);
    const bad = out.filter((k) => !tb.accepts.includes(k));
    if (bad.length) {
      const takes = tb.accepts.length ? `takes ${payloadText(tb.accepts)}` : "takes nothing now: the nodes after it take different kinds";
      const fix = mend(out, tb.accepts);
      return { ok: false, why: `${b.name} ${takes}, but ${a.name} carries ${payloadText(out)}${fix ? ": " + fix : ""}` };
    }
    return { ok: true };
  }

  const compatOf = (a) => Object.fromEntries(topo.filter((x) => x.name !== a.name).map((x) => [x.name, canConnect(a, x)]));

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

  // curve is an edge from a point to another: a cubic, and its middle.
  function curve(x1, y1, x2, y2) {
    const dx = Math.max(40, Math.abs(x2 - x1) / 2);
    const c1 = { x: x1 + dx, y: y1 };
    const c2 = { x: x2 - dx, y: y2 };
    return { d: `M ${x1} ${y1} C ${c1.x} ${c1.y}, ${c2.x} ${c2.y}, ${x2} ${y2}`, mid: { x: (x1 + 3 * c1.x + 3 * c2.x + x2) / 8, y: (y1 + 3 * c1.y + 3 * c2.y + y2) / 8 } };
  }

  const outPoint = (name) => ({ x: pos(name).x + NODE_W, y: pos(name).y + (heights[name] ?? NODE_H) / 2 });
  const inPoint = (name) => ({ x: pos(name).x, y: pos(name).y + (heights[name] ?? NODE_H) / 2 });

  function edgePath(from, to) {
    const a = outPoint(from);
    const b = inPoint(to);
    return curve(a.x, a.y, b.x, b.y).d;
  }

  // arrow is the arrowhead of an edge of a color: one marker per color
  // in the defs of the drawing (a marker does not take the stroke of its
  // path everywhere), hollow for a class's member.
  const markers = new Set();
  function arrow(color, hollow) {
    const id = "arrow-" + color.replace("#", "") + (hollow ? "-h" : "");
    if (!markers.has(id)) {
      markers.add(id);
      const m = svgEl("marker", { id, viewBox: "0 0 10 10", refX: "10", refY: "5", markerWidth: "8", markerHeight: "8", orient: "auto-start-reverse" });
      m.append(svgEl("path", { d: "M 0 0 L 10 5 L 0 10 z", fill: hollow ? "#fff" : color, stroke: color, "stroke-width": hollow ? "1.5" : "0" }));
      defs.append(m);
    }
    return `url(#${id})`;
  }

  // connecting is the edge being dragged, if one is.
  const connecting = () => (drag && drag.kind === "edge" ? drag : null);

  function draw() {
    svg.replaceChildren();
    defs = svgEl("defs");
    markers.clear();
    heights = Object.fromEntries(topo.map((n) => [n.name, heightOf(n)]));
    world = svgEl("g", { id: "world", transform: `translate(${view.x} ${view.y}) scale(${view.k})` });
    const nodes = byName();
    for (const n of topo) {
      const t = typeOf(n);
      for (const e of edgesOf(n)) {
        if (!nodes[e.to]) continue;
        const a = outPoint(n.name);
        const b = inPoint(e.to);
        const { d, mid } = curve(a.x, a.y, b.x, b.y);
        const hit = svgEl("path", { d, fill: "none", stroke: "transparent", "stroke-width": "14", class: editable ? "cursor-pointer" : "" });
        // the edge carries what the node sends on (or is typed with, by
        // what is after it): colored by that; broken when the node after
        // it does not take what reaches it along it
        const broken = (n.broken || {})[e.to];
        const kinds = e.kind === "fallback" ? ["bdf"] : t.out;
        const color = broken ? BROKEN_COLOR : payloadColor(kinds);
        const line = svgEl("path", { d, fill: "none", stroke: color, "stroke-width": broken ? "2.5" : kinds.length ? "2" : "1.5", "marker-end": arrow(color, e.kind === "member"), "pointer-events": "none" });
        if (broken) line.setAttribute("stroke-dasharray", "8 5");
        else if (e.kind === "fallback") {
          line.setAttribute("stroke-dasharray", "1 5");
          line.setAttribute("stroke-linecap", "round");
        }
        const title = svgEl("title");
        let what;
        if (broken) {
          const fix = mend(asList(broken), typeOf(nodes[e.to]).accepts);
          what = `${n.name} → ${e.to} no longer connects: ${e.to} takes ${payloadText(nodes[e.to].takes)}, but ${payloadText(broken)} reaches it here${fix ? ": " + fix : ""}, or remove the edge`;
        } else if (e.kind === "fallback") what = `${n.name} prints on ${e.to} as a fallback`;
        else if (e.kind === "member") what = `${e.to} is a member of ${n.name}: ${kinds.length ? payloadText(kinds) : "not typed yet"}`;
        else what = `${n.name} → ${e.to}: ${kinds.length ? payloadText(kinds) : "not typed yet"}`;
        title.textContent = editable ? what + " (click to remove)" : what;
        hit.append(title);
        if (editable) {
          hit.addEventListener("pointerdown", (ev) => ev.stopPropagation());
          hit.addEventListener("click", () => onRemoveEdge(n, e));
        }
        world.append(hit, line);
        if (broken) {
          // a mark on the broken edge, where it is seen even when long
          const mark = svgEl("g", { transform: `translate(${mid.x} ${mid.y})`, "pointer-events": "none" });
          mark.append(svgEl("circle", { r: "8", fill: BROKEN_COLOR, stroke: "#fff", "stroke-width": "2" }));
          const x = svgEl("text", { "text-anchor": "middle", y: "4", "font-size": "11", "font-weight": "700", fill: "#fff" });
          x.textContent = "!";
          mark.append(x);
          world.append(mark);
        }
      }
    }
    for (const n of topo) world.append(nodeBox(n, pos(n.name), n.name === selected));
    if (draft && draftPos) {
      const g = svgEl("g", { transform: `translate(${draftPos.x} ${draftPos.y})`, class: "cursor-grab", "data-draft": "1" });
      const kind = KINDS[draft.kind] || draft.kind;
      const lines = draftNote() ? wrap(draftNote(), NODE_W - 20, MAX_LINES) : [];
      box(g, draft.kind, draftLabel() || "new " + kind.toLowerCase(), kind, "not saved yet", lines, heightOfLines(lines.length), 22, 30);
      g.firstChild.setAttribute("stroke", "#4f46e5");
      g.firstChild.setAttribute("stroke-width", "2");
      g.firstChild.setAttribute("stroke-dasharray", "6 4");
      g.addEventListener("pointerdown", (ev) => startDrag(ev, null));
      world.append(g);
    }
    flowsLayer = svgEl("g", { id: "flows", "pointer-events": "none" });
    world.append(flowsLayer);
    // the edge being dragged stays on top, through a redraw (the
    // activity polled)
    const c = connecting();
    if (c) world.append(c.line, c.label);
    svg.append(defs, world);
  }

  // portTitle tells what an input port takes, as it is typed now.
  function portTitle(n, t) {
    if (t.dead) return "Nothing fits here: the nodes after it take different kinds";
    if (ENTRY.has(n.kind)) return `Other nodes may send ${t.accepts.length === ALL.length ? "any kind" : payloadText(t.accepts)} here`;
    if (t.in.length) return `Takes ${payloadText(t.in)}` + (PASSES_ON.has(n.kind) && !t.anyKind ? (asList(n.emits).length ? " (what reaches it)" : " (what is after it takes)") : "");
    if (t.anyKind) return "Takes any kind";
    return "Not typed yet: the first edge, in or out, types it";
  }

  // nodeBox draws a node: its title, its kind and what is at it, an input
  // port at the left (an entry has an arrow from outside instead of, or
  // besides, it) and, in the editor, an output port at the right for the
  // kinds that send on. The border and the ports are colored by what the
  // node is typed with. While an edge is dragged, a node that would not
  // take it fades, and the port of one that would is lit.
  // box draws the body of a node into a group, as a C4 element: the
  // white box of a height, the title bar in the role's color with the
  // kind's icon and the node's name, the kind as a tag ([Printer]) with
  // the title and what is at the node after it, and the lines of the
  // description.
  function box(g, kind, title, tag, rest, lines, h, titleMax, restMax) {
    g.append(svgEl("rect", { width: NODE_W, height: h, rx: "8", fill: BASE_FILL, stroke: BORDER, "stroke-width": "1.5" }));
    g.append(svgEl("path", { d: `M 0 8 a 8 8 0 0 1 8 -8 h ${NODE_W - 16} a 8 8 0 0 1 8 8 v ${HEAD_H - 8} h ${-NODE_W} z`, fill: roleColor(kind) }));
    const icon = iconOf(kind, 14);
    icon.setAttribute("x", "8");
    icon.setAttribute("y", "4");
    icon.setAttribute("color", "#fff");
    const name = svgEl("text", { x: "27", y: "15.5", "font-size": "12", "font-weight": "600", fill: "#fff" });
    name.textContent = clipText(title, titleMax);
    const sub = svgEl("text", { x: "10", y: TAG_Y, "font-size": "10.5", fill: "#374151" });
    const t = svgEl("tspan", { fill: "#6b7280", "font-weight": "500" });
    t.textContent = `[${tag}]`;
    sub.append(t);
    if (rest) sub.append(" · " + clipText(rest, Math.max(6, restMax - tag.length - 3)));
    g.append(icon, name, sub);
    lines.forEach((line, i) => {
      const l = svgEl("text", { x: "10", y: DESC_Y + i * LINE_H, "font-size": "11", fill: "#4b5563" });
      l.textContent = line;
      g.append(l);
    });
  }

  function nodeBox(n, p, isSelected) {
    const t = typeOf(n);
    const c = connecting();
    const cm = c && c.from.name !== n.name ? c.compat[n.name] : null;
    const faded = cm && !cm.ok && !cm.already;
    const g = svgEl("g", { transform: `translate(${p.x} ${p.y})`, class: editable ? "cursor-grab" : "cursor-pointer", "data-name": n.name });
    if (faded) g.setAttribute("opacity", "0.35");
    const lines = linesOf(n);
    const h = heightOfLines(lines.length);
    if (isSelected) g.append(svgEl("rect", { x: "-4", y: "-4", width: NODE_W + 8, height: h + 8, rx: "11", fill: "none", stroke: "#4f46e5", "stroke-width": "2", opacity: "0.7" }));
    // a sink shows the first page of the latest document it keeps
    const latest = latestOf(n);
    // the name is the title bar's; the kind, the title and what is at
    // the node are the line under it
    const parts = [];
    if (n.title && n.title !== n.name) parts.push(n.title);
    parts.push(...counts(n));
    const first = g.childElementCount;
    box(g, n.kind, n.name, KINDS[n.kind] || n.kind, parts.join(" · "), lines, h, latest ? 14 : 20, latest ? 24 : 32);
    // a node nothing fits any more is marked
    if (t.dead) {
      g.children[first].setAttribute("stroke", BROKEN_COLOR);
      g.children[first].setAttribute("stroke-width", "2");
    }
    const title = svgEl("title");
    const typed = [];
    if (t.in.length || t.anyKind) typed.push("takes " + (t.anyKind ? "any kind" : payloadText(t.in)));
    if (t.out.length) typed.push("sends on " + payloadText(t.out));
    title.textContent = `${n.name} (${KINDS[n.kind] || n.kind}, ${ROLES[ROLE_OF[n.kind]] || ""})${n.title ? ": " + n.title : ""}${typed.length ? " · " + typed.join(", ") : ""}${counts(n).length ? " · " + counts(n).join(", ") : ""}${n.description ? "\n" + n.description : ""}`;
    g.append(title);
    const l = live ? live.nodes[n.name] : null;
    if (l && (l.held || l.open || l.pending || l.failed)) {
      // a badge of what waits at the node
      const count = (l.held || 0) + (l.open || 0) + (l.pending || 0) + (l.failed || 0);
      const badge = svgEl("g", { transform: `translate(${NODE_W - 14} -8)` });
      badge.append(svgEl("circle", { r: "10", fill: l.failed ? "#dc2626" : "#f59e0b", stroke: "#fff", "stroke-width": "2" }));
      const x = svgEl("text", { "text-anchor": "middle", y: "4", "font-size": "11", "font-weight": "700", fill: "#fff" });
      x.textContent = count > 99 ? "99+" : String(count);
      badge.append(x);
      g.append(badge);
    }
    if (latest) {
      g.append(svgEl("rect", { x: NODE_W - 46, y: 4, width: 40, height: 46, rx: "3", fill: "#fff", stroke: "#94a3b8" }));
      g.append(svgEl("image", { href: `/api/jobs/${latest}/pages/1.png?w=96`, x: NODE_W - 45, y: 5, width: 38, height: 44, preserveAspectRatio: "xMidYMin slice" }));
    }
    if (ENTRY.has(n.kind)) {
      // what an entry receives comes from outside: computers print to a
      // CUPS server's port, applications POST to an endpoint
      const inArrow = svgEl("path", { d: `M -30 ${h / 2} L -8 ${h / 2}`, fill: "none", stroke: "#64748b", "stroke-width": "2", "stroke-dasharray": "3 3", "marker-end": arrow("#64748b") });
      const x = svgEl("title");
      x.textContent = n.kind === "cups" ? `IPP${n.port ? " on port " + n.port : ""}: computers print here, and the job goes on along the edges` : `REST${n.port ? " on port " + n.port : ""}: applications POST here, and the job goes on along the edges`;
      inArrow.append(x);
      g.append(inArrow);
    }
    // a CUPS server takes nothing from the other nodes: no input port.
    // The ports are colored by the kind of payload they take and send
    // on; one not typed yet is hollow
    if (n.kind !== "cups") {
      if (cm && cm.ok) g.append(svgEl("circle", { cx: "0", cy: h / 2, r: "14", fill: c.color, opacity: "0.25" }));
      const input = svgEl("circle", { cx: "0", cy: h / 2, r: cm && cm.ok ? "8" : "6", fill: "#f8fafc", stroke: "#fff", "stroke-width": "2" });
      if (t.dead) input.setAttribute("stroke", BROKEN_COLOR);
      else if (t.in.length) input.setAttribute("fill", payloadColor(t.in));
      else if (t.anyKind) input.setAttribute("fill", ANY_KIND_COLOR);
      else {
        input.setAttribute("stroke", UNTYPED_COLOR);
        input.setAttribute("stroke-dasharray", "2 2");
      }
      const x = svgEl("title");
      x.textContent = portTitle(n, t);
      input.append(x);
      g.append(input);
    }
    if (OUT_PORT.has(n.kind)) {
      const port = svgEl("circle", { cx: NODE_W, cy: h / 2, r: editable ? "7" : "5", fill: t.out.length ? payloadColor(t.out) : "#f8fafc", stroke: "#fff", "stroke-width": "2" });
      if (!t.out.length) {
        port.setAttribute("stroke", UNTYPED_COLOR);
        port.setAttribute("stroke-dasharray", "2 2");
      }
      if (editable) {
        port.setAttribute("class", "cursor-crosshair");
        port.setAttribute("data-port", n.name);
      }
      const x = svgEl("title");
      const sends = t.out.length ? `${asList(n.emits).length ? "Sends on" : "Carries"} ${payloadText(t.out)}` : "Not typed yet: the first edge types it";
      x.textContent = editable ? (n.kind === "display" ? "Drag onto a printer: the fallback printer. " : n.kind === "class" ? "Drag onto a node: a member. " : "Drag onto a node to connect. ") + sends : sends;
      port.append(x);
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
      maxY = Math.max(maxY, p.y + heightOf(n));
    }
    const k = Math.min(1.5, Math.max(0.25, Math.min((w - 2 * PAD) / (maxX - minX), (h - 2 * PAD) / (maxY - minY))));
    view = { x: (w - (maxX - minX) * k) / 2 - minX * k, y: (h - (maxY - minY) * k) / 2 - minY * k, k };
    draw();
  }

  // wireLabel is the label that follows a dragged edge: the node under
  // it and whether it takes the edge.
  function wireLabel() {
    const g = svgEl("g", { "pointer-events": "none", visibility: "hidden" });
    g.append(svgEl("rect", { rx: "4", height: "20", y: "-14", fill: "#111827", opacity: "0.85" }), svgEl("text", { x: "6", "font-size": "11", fill: "#fff" }));
    return g;
  }

  function setLabel(g, text, p, color) {
    if (!text) {
      g.setAttribute("visibility", "hidden");
      return;
    }
    const t = g.lastChild;
    t.textContent = text;
    t.setAttribute("fill", color);
    g.setAttribute("visibility", "visible");
    let width = text.length * 6.2 + 12;
    try {
      width = t.getComputedTextLength() + 12;
    } catch {}
    g.firstChild.setAttribute("width", width);
    // to the left of the pointer when it would leave the canvas on the right
    const left = view.x + (p.x + 14 + width) * view.k > size().w;
    g.setAttribute("transform", `translate(${left ? p.x - 14 - width : p.x + 14} ${p.y + 18})`);
  }

  // nodeAt is the node drawn under a pointer event.
  const nodeAt = (ev) => document.elementFromPoint(ev.clientX, ev.clientY)?.closest("g[data-name]")?.dataset.name || null;

  // Dragging: the background pans, a node moves (a click selects it), a
  // port draws an edge, two pointers pinch. Not editable, a node is only
  // selected.
  function startDrag(ev, n) {
    if (ev.button !== 0) return;
    ev.preventDefault();
    ev.stopPropagation();
    const s = screenPoint(ev);
    if (editable && n && ev.target.dataset.port) {
      const p = outPoint(n.name);
      const out = typeOf(n).out;
      const color = out.length ? payloadColor(out) : UNTYPED_COLOR;
      const line = svgEl("path", { fill: "none", stroke: color, "stroke-width": "2.5", "stroke-dasharray": "4 3", "pointer-events": "none" });
      drag = { kind: "edge", from: n, line, label: wireLabel(), color, x1: p.x, y1: p.y, compat: compatOf(n), over: null };
      draw(); // the nodes that would not take it fade
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
      const was = drag;
      drag = null;
      if (was && was.kind === "edge") draw();
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
      // the wire follows the pointer, snaps to the port of a node that
      // takes it, and turns red over one that does not
      const w = toWorld(s);
      const name = nodeAt(ev);
      const cm = name && name !== drag.from.name ? drag.compat[name] : null;
      let end = w;
      let color = drag.color;
      let dash = "4 3";
      let text = "";
      if (cm && cm.ok) {
        end = inPoint(name);
        dash = "";
        text = `→ ${name}`;
      } else if (cm) {
        color = cm.already ? UNTYPED_COLOR : BROKEN_COLOR;
        text = (cm.already ? "" : "✕ ") + cm.why;
      }
      drag.line.setAttribute("d", curve(drag.x1, drag.y1, end.x, end.y).d);
      drag.line.setAttribute("stroke", color);
      drag.line.setAttribute("stroke-dasharray", dash);
      setLabel(drag.label, text, w, cm && cm.ok ? "#fff" : cm && !cm.already ? "#fecaca" : "#fff");
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
        const name = nodeAt(ev);
        draw(); // the wire goes, the nodes come back
        if (!name || name === d.from.name) return;
        const cm = d.compat[name] || canConnect(d.from, byName()[name]);
        if (cm.ok) onConnect(d.from, name);
        else if (!cm.already) onReject(d.from, name, cm.why);
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
      const c = connecting();
      if (c) c.compat = compatOf(c.from);
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

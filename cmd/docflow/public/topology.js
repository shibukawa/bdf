// The topology: everyone's page, and the editor of a person with the
// role topology. The nodes and the edges are drawn on the canvas
// (canvas.js: panned and zoomed, with the activity of the jobs the
// person may see running along the edges, polled with its revision).
// Editing: nodes of every kind are placed, connected by dragging from a
// node's output port onto another node, and set up in the sidebar (the
// fields of the node's kind under Settings; how it is reached from
// outside and what is at it under Usage: a CUPS server's printer, an
// uploader's request and key, the API of an inbox and of a downloader, a
// display's devices and enrollment codes);
// everything goes through the API, and the server checks a node when it
// is put (its edges, its name, no cycle, and that every node takes the
// kinds of payload that reach it along a new edge; a change of settings
// may leave an edge broken, which the canvas draws dashed and the sidebar
// explains). A node is put before the nodes after it, with no edges yet,
// and renamed when its name is changed: its edges and what is kept under
// its name follow. Looking: a node shows the jobs sent to it that the
// person sees, and a form to print there when they may. Both show what is
// at a node now, read again as the activity changes (atSection): what a
// sink keeps, a display's board, what an inbox holds (released or
// discarded there by those who work it), what an archive wrote lately
// and what a downloader lists, each document's pages in the lightbox.

import { $, _, api, archivedSection, boardSection, button, changed, el, heldSection, keptSection, kindBadge, KINDS, legsOf, listedSection, locale, onChange, PAYLOADS, sender, sentences, since, stateBadge, STATES, status } from "./admin.js";
import { createCanvas, iconOf, legend, payloadColor, PRINTS, ROLE_OF, roleColor, typeOf } from "./canvas.js";

const SCHEMES = ["ipp", "ipps", "socket", "lpd"];
const ACCEPTS = { file: _("documents (PDF, Office, images, .bdf …)"), data: _("data (JSON, YAML, CSV, Excel, Parquet)") };
const OUTPUTS = { bdf: _("the .bdf document made of it"), file: _("the file as received (a convert node makes the document later)") };
const KINDS_OF_PAYLOAD = { "": _("payload: any kind"), bdf: _("payload: .bdf document"), data: _("payload: data"), any: _("payload: any other file") };
const SUBJECTS = { email: _("email: the person's email, as the IdP verified it"), sub: _("sub: the person's subject at the IdP"), id: _("id: the person's ID here"), group: _("group: any of the person's groups") };
const POLL_MS = 4000;

let topo = []; // the nodes, as /api/topology lists them
let formats = []; // the formats an entry may take
let me = null; // /api/me: whether this process has ports of its own
let apps = []; // the apps, for the allow-lists of entries, inboxes and downloaders
let kinds = {}; // GET /api/kinds by name: the settings each requires
// bonjourFound is what Bonjour found last, by kind (a printer's: every
// IPP printer; a CUPS server (sink)'s: the shared queues), that the
// sidebar's AirPrint lists are made of before they are searched again.
const bonjourFound = new Map();
const bonjourAsked = new Map(); // the searches going on, by kind

// bonjourSearch asks what Bonjour finds for a kind, one search at a time.
function bonjourSearch(kind) {
  if (!bonjourAsked.has(kind)) {
    const p = api("POST", `/api/kinds/${kind}/discover`, {})
      .then((list) => (bonjourFound.set(kind, list), list))
      .finally(() => bonjourAsked.delete(kind));
    bonjourAsked.set(kind, p);
  }
  return bonjourAsked.get(kind);
}
let selected = null; // the node in the sidebar
// drafts are the nodes being added, not put yet: {id, kind, pos, kept}
// (kept: their form when the sidebar left them, shown again when they
// are clicked); draft is the one in the sidebar. A draft is put by its
// Add, and goes with Cancel (or Escape).
let drafts = [];
let draft = null;
let drafted = 0; // the drafts made, for their ids
// edits are the forms of saved nodes whose settings were changed in the
// sidebar and not saved, by name, kept when the sidebar left them:
// shown again when the node is, until Save or Revert.
const edits = new Map();
// shown is what the sidebar's form is of: {name, draft, base (the form
// as the saved node makes it, to tell a change), touched (the ids of the
// fields left), all (every required field empty is marked), restored
// (it shows edits kept), dirty, marked}, or null.
let shown = null;
let fitted = false;
let etag = ""; // the activity's revision
let live = null; // the activity, as /api/topology/live last told it
// editable: the server rendered the editor's controls (the palette) for
// a person with the role topology; anyone else looks, and prints where
// they may. The API checks the role itself.
const editable = !!$("palette");
let dests = []; // looking: the nodes the person may send to
let jobs = []; // looking: the jobs the person sees, latest first

const canvas = createCanvas($("canvas"), {
  editable,
  onSelect(name) {
    show(name, null);
  },
  onSelectDraft(id) {
    show(null, drafts.find((d) => d.id === id) || null);
  },
  onMove(name, p) {
    api("PUT", "/api/layout", { positions: { [name]: p } }).catch((e) => status(e.message, true));
  },
  onConnect(from, to, port) {
    addEdge(from, to, port).catch((e) => status(e.message, true));
  },
  onReject(from, to, why) {
    status(`${from.name} → ${to}: ${why}`, true);
  },
  onRemoveEdge: removeEdge,
  draftLabel: (d) => typed(d, "f-title").trim() || typed(d, "f-name").trim(),
  draftNote: (d) => typed(d, "f-description"),
  marksOf,
});

// typed is a field of a draft as typed: in the sidebar when the draft is
// there, else as it was kept.
function typed(d, id) {
  if (shown && shown.draft === d) return shown.form.querySelector("#" + id)?.value || "";
  return (d.kept && d.kept.raw[id]) || "";
}

// marksOf tells the canvas what is not saved at a node or a draft (by
// its id): the settings changed, and the required fields left empty (in
// the sidebar, those marked; of one left, every one, as it shows them
// when it is opened again).
function marksOf(key) {
  const d = drafts.find((x) => x.id === key);
  if (d) return shown && shown.draft === d ? { missing: shown.marked } : { missing: d.kept ? d.kept.missing : 0 };
  if (shown && shown.name === key && !shown.draft) return shown.dirty ? { edited: true, missing: shown.marked } : null;
  const kept = edits.get(key);
  return kept ? { edited: true, missing: kept.missing } : null;
}

// show opens a node (by name) or a draft in the sidebar, keeping what
// was there first.
function show(name, d) {
  keepShown();
  selected = d ? null : name;
  draft = d;
  renderSidebar();
  syncCanvas();
  $("sidebar").scrollTop = 0;
}

// syncCanvas tells the canvas the drafts and what is selected.
function syncCanvas() {
  canvas.setDrafts(drafts, draft ? draft.id : null);
  if (!draft) canvas.setSelected(selected);
}

async function load() {
  try {
    const [t, l] = await Promise.all([api("GET", "/api/topology"), api("GET", "/api/layout")]);
    topo = t;
    if (editable) {
      if (!formats.length) formats = await api("GET", "/api/formats");
      if (!me) me = await api("GET", "/api/me");
      if (!Object.keys(kinds).length) kinds = Object.fromEntries((await api("GET", "/api/kinds")).map((k) => [k.name, k]));
      apps = await api("GET", "/api/apps");
      $("node-names").replaceChildren(...topo.map((n) => el("option", undefined, { value: n.name })));
    } else {
      [dests, jobs] = await Promise.all([api("GET", "/api/destinations"), api("GET", "/api/my/jobs?limit=200")]);
    }
    // the edits of a node deleted (elsewhere) go with it
    for (const name of edits.keys()) if (!topo.some((n) => n.name === name)) edits.delete(name);
    if (selected && !topo.some((n) => n.name === selected)) {
      selected = null;
      if (shown && !shown.draft) shown = null;
    }
    canvas.setData(topo, l.positions || {});
    if (!fitted) {
      canvas.fit();
      fitted = true;
    }
    renderSidebar();
  } catch (e) {
    status(e.message, true);
  }
}

// poll reads the activity (the jobs the person sees; every job, for an
// administrator) when its revision changed, and, looking, the jobs with
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
      if (!editable) {
        jobs = await api("GET", "/api/my/jobs?limit=200");
        renderSidebar();
      } else if (at && at.box.isConnected) {
        at.box.refresh();
      }
    } else if (res.status !== 304) {
      throw new Error("live " + res.status);
    }
    $("live").textContent = _("Live · {time}", { time: new Date().toLocaleTimeString(locale, { hour: "2-digit", minute: "2-digit", second: "2-digit" }) });
    $("live").className = "badge badge-success badge-soft badge-sm";
  } catch (e) {
    $("live").textContent = _("Not connected: {why}", { why: e.message });
    $("live").className = "badge badge-error badge-sm";
  }
}

$("zoom-in").addEventListener("click", () => canvas.zoomBy(1.25));
$("zoom-out").addEventListener("click", () => canvas.zoomBy(0.8));
$("zoom-fit").addEventListener("click", () => canvas.fit());
if (editable) $("refresh").addEventListener("click", load);
legend($("legend"));
// the palette's kinds get their icons, in their roles' colors
for (const a of document.querySelectorAll("#palette a[data-kind]")) a.prepend(kindIcon(a.dataset.kind, 18));

// kindIcon is the icon of a kind, in its role's color.
function kindIcon(kind, size) {
  const i = iconOf(kind, size);
  i.style.color = roleColor(kind);
  i.classList.add("shrink-0");
  return i;
}

// brokenEdges names the edges that no longer carry what the node after
// them takes.
const brokenEdges = (t) => t.flatMap((n) => Object.entries(n.broken || {}).map(([to, kinds]) => `${n.name} → ${to} (${kinds.map((k) => PAYLOADS[k] || k).join(", ")})`));

// either says kinds of payload as alternatives: a .bdf document or data.
const either = (kinds) => kinds.map((k) => PAYLOADS[k] || k).join(_(" or "));

// addEdge connects a node to a target through the API, which checks the
// edge (the target exists, no cycle): a class's member, a display's
// fallback printer, a branch's case (by its port) or default, an
// assert's on_error (its fail port), or the node's to.
async function addEdge(from, to, port) {
  if (from.name === to) return;
  const n = await api("GET", "/api/nodes/" + from.name);
  const arm = n.kind === "branch" && port && port.startsWith("case:") ? (n.cases || [])[Number(port.slice(5))] : null;
  if (n.kind === "display") {
    n.fallback = { ...(n.fallback || {}), printer: to };
  } else if (n.kind === "filter") {
    if ((n.otherwise || []).includes(to)) return;
    n.otherwise = [...(n.otherwise || []), to];
  } else if (n.kind === "branch" || (n.kind === "assert" && port === "error")) {
    const field = n.kind === "assert" ? "on_error" : "default";
    const list = arm ? arm.to || [] : n[field] || [];
    if (list.includes(to)) return;
    if (arm) arm.to = [...list, to];
    else n[field] = [...list, to];
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
  const ask = e.kind === "fallback" ? _("Remove the fallback printer {to} of {from}?", { to: e.to, from: from.name }) : _("Remove the edge {from} → {to}?", { from: from.name, to: e.to });
  if (!confirm(ask)) return;
  try {
    const n = await api("GET", "/api/nodes/" + from.name);
    if (e.kind === "fallback") n.fallback = null;
    else if (n.kind === "branch" && e.port && e.port.startsWith("case:")) {
      const arm = (n.cases || [])[Number(e.port.slice(5))];
      if (arm) arm.to = (arm.to || []).filter((x) => x !== e.to);
    } else if (n.kind === "branch" && e.port === "default") n.default = (n.default || []).filter((x) => x !== e.to);
    else if (n.kind === "filter") {
      n.otherwise = (n.otherwise || []).filter((x) => x !== e.to);
      n.rules = (n.rules || []).map((r) => ({ ...r, to: (r.to || []).filter((x) => x !== e.to) })).filter((r) => r.to.length);
    } else {
      const field = e.kind === "member" ? "members" : e.kind === "original" ? "original" : e.kind === "error" ? "on_error" : "to";
      n[field] = (n[field] || []).filter((x) => x !== e.to);
    }
    delete n.name;
    await api("PUT", "/api/nodes/" + from.name, n);
    await load();
  } catch (err) {
    status(err.message, true);
  }
}

// The kinds of the adapters (docflow/io, docflow/render) are set as JSON
// in the sidebar: their settings, and what they mean.
const ADAPTER_SETTINGS = {
  render: _("to, format (pdf or png), page (of a png)."),
  notify: _("via (slack, mail or http), text, credential (slack: the webhook's URL), url (http; mail: smtp://host:587), from, recipients, mail_subject, to."),
  smtp: _("url (smtp://host:587, smtps://host:465), credential (user:password), from, recipients, mail_subject, text."),
  s3: _("url (s3://bucket?region=…&endpoint=…), credential (access-key:secret-key), path (the key; {{meta.$job}}/{{meta.$name}} by default)."),
  http: _("url, method (POST or PUT), body (direct: the document as the body, the metadata as X-Meta-* headers; multipart: form-data, the metadata as JSON in the part meta and the document in the part file), headers, credential (the Authorization: a token as Bearer)."),
};

// SETTINGS are the fields of a node an adapter's kind may have.
const SETTINGS = ["to", "url", "credential", "format", "page", "via", "text", "from", "recipients", "mail_subject", "path", "method", "body", "headers", "accept", "output", "formats", "meta"];

function settingsOf(n, kind) {
  const out = {};
  for (const k of SETTINGS) if (n[k] !== undefined && n[k] !== "" && !(Array.isArray(n[k]) && !n[k].length)) out[k] = n[k];
  if (!Object.keys(out).length) {
    out.to = [];
    if (kind === "notify") out.via = "slack";
  }
  return out;
}

// Adding a node: the palette makes a draft, placed at the centre of the
// view clear of the others; the sidebar takes its name and settings, and
// Add puts it. A draft left before it is added stays on the canvas, with
// what was typed for it, until it is added or canceled.
if (editable) $("palette").addEventListener("click", (ev) => {
  const a = ev.target.closest("a[data-kind]");
  if (!a) return;
  ev.preventDefault();
  document.activeElement?.blur();
  const d = { id: "+" + ++drafted, kind: a.dataset.kind, pos: canvas.freeSpot(canvas.center()), kept: null };
  drafts.push(d);
  // the sidebar first: the draft is drawn with its title and description
  show(null, d);
  $("f-name")?.focus();
});

// dropDraft cancels a draft: it goes from the canvas, with what was
// typed for it.
function dropDraft(d) {
  drafts = drafts.filter((x) => x !== d);
  if (draft === d) {
    draft = null;
    shown = null;
  }
  renderSidebar();
  syncCanvas();
}

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

// picker is the nodes a node may send to, as checkboxes (name: the
// checkboxes' name; except the node itself and the entries), each with a
// dot of what it takes; one that would not take what the node carries
// says so.
function picker(name, chosen, except, carried) {
  const box = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
  const from = draft || topo.find((x) => x.name === selected) || {};
  const carries = carried || typeOf(from).out;
  const names = topo.filter((x) => x.name !== except && x.kind !== "cups").map((x) => x.name);
  for (const x of names) {
    const target = topo.find((t) => t.name === x);
    const t = typeOf(target);
    const bad = carries.filter((k) => !t.accepts.includes(k));
    const title = bad.length
      ? t.accepts.length
        ? _("{node} takes {takes}, but this carries {carries}", { node: x, takes: either(t.accepts), carries: either(bad) })
        : _("{node} takes nothing now, but this carries {carries}", { node: x, carries: either(bad) })
      : t.accepts.length === 3
        ? _("may be sent any kind")
        : _("may be sent {kinds}", { kinds: either(t.accepts) });
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" + (bad.length ? " text-error" : ""), title });
    const dot = el("span", undefined, { className: "inline-block size-2.5 rounded-full border border-base-100" });
    dot.style.background = t.dead ? "#dc2626" : t.in.length ? payloadColor(t.in) : t.anyKind ? "#475569" : "#f8fafc";
    if (!t.in.length && !t.anyKind && !t.dead) dot.style.borderColor = "#94a3b8";
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name, value: x, checked: chosen.includes(x) }), dot, bad.length ? _("{node} (would not take it)", { node: x }) : x);
    box.append(l);
  }
  // one chosen that is not there any more stays chosen, so that it is not lost unseen
  for (const x of chosen) {
    if (names.includes(x)) continue;
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm text-error" });
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name, value: x, checked: true }), _("{node} (gone)", { node: x }));
    box.append(l);
  }
  if (!names.length) box.append(el("span", _("No other node yet"), { className: "text-xs text-base-content/60" }));
  return box;
}

// picked are the values checked of the checkboxes of a name in a form.
const picked = (form, name) => [...form.querySelectorAll(`input[name="${name}"]:checked`)].map((c) => c.value);

// ruleRow is a rule of a filter in the sidebar.
function ruleRow(r = { match: {}, to: [] }) {
  const row = el("div", undefined, { className: "rule grid grid-cols-2 gap-1 rounded-box bg-base-200 p-2" });
  const m = r.match || {};
  row.append(
    el("input", undefined, { className: "input input-xs", placeholder: _("name matches (RE2)"), value: m.name || "", "data-k": "name" }),
    el("input", undefined, { className: "input input-xs", placeholder: _("user matches"), value: m.user || "", "data-k": "user" }),
    el("input", undefined, { className: "input input-xs", placeholder: _("format matches (pdf, bdf …)"), value: m.format || "", "data-k": "format" }),
    el("div", undefined, { className: "flex gap-1" }),
  );
  row.lastChild.append(el("input", undefined, { className: "input input-xs w-full", placeholder: _("pages ≥"), value: m.min_pages || "", inputMode: "numeric", "data-k": "min_pages" }), el("input", undefined, { className: "input input-xs w-full", placeholder: _("pages ≤"), value: m.max_pages || "", inputMode: "numeric", "data-k": "max_pages" }));
  const meta = el("input", undefined, { className: "input input-xs", placeholder: _("metadata matches: key=RE2, key2=RE2"), value: Object.entries(m.meta || {}).map(([k, v]) => k + "=" + v).join(", "), "data-k": "meta", title: _("The sender's keys, the entry's, and $job, $entry, $sender, $file, $kind, $format, $pages …") });
  const kind = el("select", undefined, { className: "select select-xs", "data-k": "kind" });
  for (const [v, label] of Object.entries(KINDS_OF_PAYLOAD)) kind.append(el("option", label, { value: v, selected: v === (m.kind || "") }));
  row.append(meta, kind);
  const to = el("input", undefined, { className: "input input-xs col-span-2", placeholder: _("→ nodes, comma-separated"), value: (r.to || []).join(", "), "data-k": "to" });
  to.setAttribute("list", "node-names");
  const drop = el("button", _("remove the rule"), { type: "button", className: "btn btn-xs btn-ghost col-span-2 justify-self-end text-error" });
  drop.addEventListener("click", () => row.remove());
  row.append(to, drop);
  return row;
}

// rulesOf reads the filter's rules from its form.
function rulesOf(form) {
  const out = [];
  for (const row of form.querySelectorAll("#f-rules .rule")) {
    const v = (k) => row.querySelector(`[data-k="${k}"]`).value.trim();
    const match = {};
    if (v("name")) match.name = v("name");
    if (v("user")) match.user = v("user");
    if (v("format")) match.format = v("format");
    if (v("min_pages")) match.min_pages = Number(v("min_pages"));
    if (v("max_pages")) match.max_pages = Number(v("max_pages"));
    if (v("kind")) match.kind = v("kind");
    const meta = {};
    for (const pair of v("meta").split(",")) {
      const i = pair.indexOf("=");
      if (i > 0) meta[pair.slice(0, i).trim()] = pair.slice(i + 1).trim();
    }
    if (Object.keys(meta).length) match.meta = meta;
    const to = names(v("to"));
    if (!Object.keys(match).length && !to.length) continue;
    out.push({ match, to });
  }
  return out;
}

// checksBox is the checks of an assert, or of a branch's case, in the
// sidebar: a row each, and a button that adds one. An assert's say what
// they tell a sender who fails them (message).
function checksBox(checks, assert) {
  const rows = el("div", undefined, { className: "checks flex flex-col gap-1" });
  if (assert) rows.id = "f-checks";
  for (const c of checks) rows.append(checkRow(c, assert));
  const add = el("button", _("+ a check"), { type: "button", className: "btn btn-xs w-fit" });
  add.addEventListener("click", () => rows.append(checkRow({}, assert)));
  const box = el("div", undefined, { className: "flex flex-col gap-1" });
  box.append(rows, add);
  return box;
}

function checkRow(c, assert) {
  const row = el("div", undefined, { className: "check grid grid-cols-2 gap-1 rounded-box border border-base-300 bg-base-100 p-2" });
  const need = el("select", undefined, { className: "select select-xs", "data-k": "optional" });
  need.append(el("option", _("required"), { value: "", selected: !c.optional }), el("option", _("optional: may be absent"), { value: "1", selected: !!c.optional }));
  const num = (k, placeholder) => el("input", undefined, { className: "input input-xs w-full", placeholder, value: c[k] ?? "", inputMode: "decimal", "data-k": k });
  const range = el("div", undefined, { className: "col-span-2 flex gap-1" });
  range.append(num("min", _("a number ≥ (min)")), num("max", _("≤ (max)")));
  row.append(
    el("input", undefined, { className: "input input-xs", placeholder: _("key: doc-type, ipp.job-account-id, $pages …"), value: c.key || "", "data-k": "key", autocomplete: "off", title: _("The sender's keys, the entry's, and $job, $entry, $sender, $file, $kind, $format, $pages, $size …") }),
    need,
    el("input", undefined, { className: "input input-xs col-span-2 font-mono", placeholder: _("matches (RE2): {example}", { example: "^\\d{3}-\\d{4}$" }), value: c.pattern || "", "data-k": "pattern", autocomplete: "off" }),
    range,
    el("input", undefined, { className: "input input-xs col-span-2", placeholder: _("one of: invoice, order (comma-separated)"), value: (c.enum || []).join(", "), "data-k": "enum", autocomplete: "off" }),
  );
  if (assert) row.append(el("input", undefined, { className: "input input-xs col-span-2", placeholder: _("what a failure says (instead of the check's own words)"), value: c.message || "", "data-k": "message", maxLength: 255 }));
  const drop = el("button", _("remove the check"), { type: "button", className: "btn btn-xs btn-ghost col-span-2 justify-self-end text-error" });
  drop.addEventListener("click", () => row.remove());
  row.append(drop);
  return row;
}

// checksOf reads the checks of a box of rows; a bound that is not a
// number is refused here (lenient: kept as typed, as a form left is).
function checksOf(box, lenient) {
  const out = [];
  for (const row of box.querySelectorAll(":scope > .check")) {
    const raw = (k) => row.querySelector(`[data-k="${k}"]`)?.value ?? "";
    const v = (k) => raw(k).trim();
    const c = { key: v("key") };
    if (!c.key) continue;
    if (v("optional")) c.optional = true;
    if (raw("pattern")) c.pattern = raw("pattern");
    for (const k of ["min", "max"]) {
      if (v(k) === "") continue;
      c[k] = Number(v(k));
      if (!Number.isFinite(c[k]) && lenient) c[k] = v(k);
      else if (!Number.isFinite(c[k])) throw new Error(_("{key}: {bound} {value} is not a number", { key: c.key, bound: k, value: v(k) }));
    }
    if (v("enum")) c.enum = names(v("enum"));
    if (v("message")) c.message = v("message");
    out.push(c);
  }
  return out;
}

// caseRow is a case of a branch in the sidebar: its name (its port's
// label), its checks and the nodes it sends to; ↑ moves it before the one
// above it, as the cases are tried in order.
function caseRow(c = { when: [], to: [] }) {
  const row = el("div", undefined, { className: "case flex flex-col gap-1 rounded-box bg-base-200 p-2" });
  const head = el("div", undefined, { className: "flex items-center gap-1" });
  const up = el("button", "↑", { type: "button", className: "btn btn-xs btn-ghost", title: _("Try this case before the one above it") });
  up.addEventListener("click", () => row.previousElementSibling && row.parentElement.insertBefore(row, row.previousElementSibling));
  const drop = el("button", _("remove the case"), { type: "button", className: "btn btn-xs btn-ghost text-error" });
  drop.addEventListener("click", () => row.remove());
  head.append(el("input", undefined, { className: "input input-xs flex-1", placeholder: _("name: invoice (its port's label; empty: case <n>)"), value: c.name || "", "data-k": "name", maxLength: 64, autocomplete: "off" }), up, drop);
  const to = el("input", undefined, { className: "input input-xs", placeholder: _("→ nodes, comma-separated (or drag from its port)"), value: (c.to || []).join(", "), "data-k": "to" });
  to.setAttribute("list", "node-names");
  row.append(head, checksBox(c.when || [], false), to);
  return row;
}

// casesOf reads the branch's cases from its form: a case with no name,
// check or node is left out.
function casesOf(form, lenient) {
  const out = [];
  for (const row of form.querySelectorAll("#f-cases > .case")) {
    const c = { name: row.querySelector('[data-k="name"]').value.trim(), when: checksOf(row.querySelector(".checks"), lenient), to: names(row.querySelector(':scope > [data-k="to"]').value) };
    if (!c.name && !c.when.length && !c.to.length) continue;
    out.push(c);
  }
  return out;
}

function selectInput(id, options, value) {
  const s = el("select", undefined, { id, className: "select select-sm w-full" });
  for (const [v, label] of options) s.append(el("option", label, { value: v, selected: v === value }));
  return s;
}

// metaRows are the metadata an entry adds to every job it takes, as
// rows of key and value in the sidebar.
function metaRow(k = "", v = "") {
  const row = el("div", undefined, { className: "meta grid grid-cols-[1fr_1fr_auto] gap-1" });
  row.append(el("input", undefined, { className: "input input-xs", placeholder: _("key (site)"), value: k, "data-k": "key" }), el("input", undefined, { className: "input input-xs", placeholder: _("value (tokyo)"), value: v, "data-k": "value" }));
  const drop = el("button", "×", { type: "button", className: "btn btn-xs btn-ghost text-error" });
  drop.addEventListener("click", () => row.remove());
  row.append(drop);
  return row;
}

function metaField(n) {
  const rows = el("div", undefined, { id: "f-meta", className: "flex flex-col gap-1" });
  for (const [k, v] of Object.entries(n.meta || {})) rows.append(metaRow(k, typeof v === "string" ? v : JSON.stringify(v)));
  const add = el("button", _("+ a key"), { type: "button", className: "btn btn-xs w-fit" });
  add.addEventListener("click", () => rows.append(metaRow()));
  const box = el("div", undefined, { className: "flex flex-col gap-1" });
  box.append(rows, add);
  return field(_("Metadata added"), box, _("Keys this entry adds to every job it takes (a site, a department): they win over the sender's (X-Meta-* headers, ?meta.<key>=), and filters route by them. Keys starting with $ are the manager's."));
}

// metaOf reads the metadata rows: a value that is JSON (a number, true,
// an object) is kept as that, else as text.
function metaOf(form) {
  const out = {};
  for (const row of form.querySelectorAll("#f-meta .meta")) {
    const k = row.querySelector('[data-k="key"]').value.trim();
    const v = row.querySelector('[data-k="value"]').value;
    if (!k) continue;
    try {
      out[k] = /^[-\d\[{"tfn]/.test(v.trim()) ? JSON.parse(v) : v;
    } catch {
      out[k] = v;
    }
  }
  return out;
}

// formatsField is the formats an entry takes, or a convert node converts.
function formatsField(n, label, what) {
  const fmts = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
  for (const f of formats) {
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "fmt", value: f, checked: (n.formats || []).includes(f) }), f);
    fmts.append(l);
  }
  return field(label, fmts, what);
}

// entryFields are what an uploader or an S3 trigger takes and sends on:
// documents or data, the formats, and whether a document goes on as the
// .bdf made of it or as received. A CUPS server has none: what computers
// print is a document, made .bdf.
function entryFields(form, n) {
  const accept = selectInput("f-accept", Object.entries(ACCEPTS), n.accept || "file");
  const output = selectInput("f-output", Object.entries(OUTPUTS), n.output || "bdf");
  const sync = () => (output.disabled = accept.value === "data");
  accept.addEventListener("change", sync);
  setTimeout(sync);
  form.append(field(_("Accepts"), accept, _("Data is for a template node after this entry: the edges leaving it carry data.")));
  form.append(field(_("Sends on"), output, _("Of a document: converted here (printers and displays take .bdf), or as received for a convert node or an archive after this. Data always goes on as it is.")));
  form.append(formatsField(n, _("Formats taken"), _("None checked: every format the converter reads (and json, yaml for data). A file of another format is refused, and the refusal is in the audit log.")));
}

// kindsLine tells what a node takes and sends on, as the API inferred
// it from the entries, what may be sent to it now (what the nodes after
// it take), and the edges from it that no longer connect.
function kindsLine(n) {
  const box = el("div", undefined, { className: "flex flex-col gap-1 text-xs text-base-content/60" });
  const parts = [];
  if (n.takes) parts.push(_("takes {kinds}", { kinds: PAYLOADS[n.takes] || n.takes }));
  const acc = n.accepts || [];
  if (n.accepts && acc.length < 3 && !(acc.length === 1 && acc[0] === n.takes)) parts.push(acc.length ? _("may be sent {kinds} only (what is after it takes)", { kinds: either(acc) }) : _("nothing fits it now: the nodes after it take different kinds"));
  if (n.emits && n.emits.length) parts.push(_("sends on {kinds}", { kinds: either(n.emits) }));
  if (parts.length) box.append(el("p", parts.join(" · ")));
  for (const [to, kinds] of Object.entries(n.broken || {})) {
    const target = topo.find((x) => x.name === to);
    const takes = target ? PAYLOADS[target.takes] || target.takes : _("something else");
    box.append(el("p", _("⚠ → {to} no longer connects: {to} takes {takes}, but {kinds} reaches it from here. Put a convert or a template node between, change what this node takes, or remove the edge; a job gets an error there meanwhile.", { to, takes, kinds: either(kinds) }), { className: "text-error" }));
  }
  return box.childElementCount ? box : null;
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
  const label = { shared: _("shared"), group: _("group"), personal: _("personal") }[scope] || scope;
  return el("span", label + (n.owner ? ": " + n.owner : ""), { className: cls });
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

// composeURL is the printer's URL of the fields of its form.
function composeURL(form) {
  const f = (id) => form.querySelector("#" + id);
  const scheme = f("f-scheme").value;
  const host = f("f-host").value.trim();
  const port = f("f-pport").value.trim();
  let path = f("f-path").value.trim();
  if (!host) return "";
  if (!path && scheme.startsWith("ipp")) path = "/ipp/print";
  if (path && !path.startsWith("/")) path = "/" + path;
  return `${scheme}://${host}${port ? ":" + port : ""}${path}`;
}

// connectFields are how a printer or a CUPS server (sink) is reached,
// put before the fields of its address (manual): the address typed, or
// picked from what Bonjour finds on this server's local network (POST
// /api/kinds/<kind>/discover with no URL: a printer's, every IPP printer,
// AirPrint, IPP Everywhere and Mopria printers alike, each named with
// what it speaks; a CUPS server (sink)'s, the queues computers and CUPS
// servers share), whose instance name the node keeps (bonjour) to
// find it again where its address changes. The fields of the way not
// chosen are hidden and not required. changed is told the URL picked
// (null: Manual chosen), to fill the address's fields with, so that
// Manual shows it. A serverless deployment has no local network: the
// address is typed there.
function connectFields(form, n, manual, changed) {
  if (me && me.ports === false) return;
  const queues = n.kind === "cups-out";
  const via = selectInput("f-via", [["manual", _("Manual: the address typed")], ["bonjour", _("Bonjour: AirPrint, IPP Everywhere, Mopria")]], n.bonjour ? "bonjour" : "manual");
  form.append(field(_("Connect"), via));
  const pick = el("select", undefined, { id: "f-bonjour", className: "select select-sm w-full", required: true });
  pick.dataset.need = queues ? _("Required: choose a queue found, or search again.") : _("Required: choose a printer found, or search again.");
  const again = el("button", _("Search again"), { type: "button", className: "btn btn-sm w-fit" });
  const said = el("span", "", { className: "text-xs text-base-content/60" });
  const row = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  row.append(again, said);
  const box = el("div", undefined, { className: "flex flex-col gap-1" });
  box.append(pick, row);
  const air = field(queues ? _("Shared queue") : _("Printer"), box, queues
    ? _("The queues that Macs and CUPS servers on this server's local network share, as Bonjour finds them. The node keeps the name: when the address stops answering, it is looked for again.")
    : _("The printers on this server's local network, as Bonjour finds them (AirPrint, IPP Everywhere and Mopria printers, each named with what it speaks, and the queues computers share). The node prints to any of them the same way, PDF first. It keeps the name: when the address stops answering, it is looked for again."));
  form.append(air);
  const opt = (name, label, url) => {
    const o = el("option", label, { value: name });
    o.dataset.url = url || "";
    return o;
  };
  // the options: what was found, and the one chosen when it is not there
  const fill = (list) => {
    const was = pick.selectedOptions[0];
    const chosen = pick.value || n.bonjour || "";
    const url = (was && was.value === chosen && was.dataset.url) || (chosen === n.bonjour ? n.url : "");
    pick.replaceChildren(opt("", _("— choose —"), ""));
    for (const x of list || []) pick.append(opt(x.name, [x.name, x.model, x.note, speaks(x), !queues && x.queue && _("shared queue")].filter(Boolean).join(" · "), x.url));
    if (chosen && ![...pick.options].some((o) => o.value === chosen)) pick.append(opt(chosen, list ? _("{name} (not found now)", { name: chosen }) : chosen, url));
    pick.value = chosen;
  };
  fill(bonjourFound.get(n.kind));
  let tried = false;
  const search = async () => {
    tried = again.disabled = true;
    said.className = "text-xs text-base-content/60";
    said.textContent = _("Looking on the local network …");
    try {
      const list = await bonjourSearch(n.kind);
      fill(list);
      said.textContent = list.length ? _("Found: {n}", { n: list.length }) : queues ? _("No shared queue found.") : _("No printer found.");
    } catch (e) {
      said.className = "text-xs text-error";
      said.textContent = e.message;
    } finally {
      again.disabled = false;
    }
  };
  again.addEventListener("click", search);
  const sync = () => {
    const on = via.value === "bonjour";
    air.hidden = !on;
    pick.disabled = !on;
    for (const f of manual) {
      f.hidden = on;
      for (const x of f.querySelectorAll("input, select, button")) x.disabled = on;
    }
    // searched once by itself; again when asked
    if (on && !tried && !bonjourFound.has(n.kind)) search();
  };
  via.addEventListener("change", () => {
    sync();
    changed(via.value === "bonjour" ? pickedURL(form) : null);
  });
  pick.addEventListener("change", () => changed(pickedURL(form)));
  sync();
}

// speaks names what a printer Bonjour found speaks, of AirPrint, IPP
// Everywhere and Mopria ("" for none of them).
function speaks(x) {
  return [x.airprint && "AirPrint", x.everywhere && "IPP Everywhere", x.mopria && "Mopria"].filter(Boolean).join(", ");
}

// pickedURL is the URL of what was picked from Bonjour when the form
// says so ("" when nothing is), null when it says Manual.
function pickedURL(form) {
  const via = form.querySelector("#f-via");
  if (!via || via.value !== "bonjour") return null;
  const o = form.querySelector("#f-bonjour").selectedOptions[0];
  return (o && o.dataset.url) || "";
}

// cupsOutFields are a CUPS server (sink)'s: its server, the queue there
// (Find queues asks the server for them: CUPS-Get-Printers, through POST
// /api/kinds/cups-out/discover, before the node is saved), or a queue
// shared on the local network (connectFields), and the credential the
// server asks for.
function cupsOutFields(form, n) {
  let path = "/printers/";
  let queue = "";
  const u = parsePrinterURL(n.url);
  const m = /^\/(printers|classes)\/([^/]+)$/.exec(u.path);
  if (m) [path, queue] = ["/" + m[1] + "/", m[2]];
  const row = el("div", undefined, { className: "grid grid-cols-[6rem_1fr_5rem] gap-1" });
  row.append(selectInput("f-cscheme", [["ipp", "ipp://"], ["ipps", "ipps://"]], u.scheme === "ipps" ? "ipps" : "ipp"), input("f-chost", u.host, { placeholder: _("the CUPS server's host"), required: true, autocomplete: "off" }), input("f-cport", u.port, { placeholder: "631", inputMode: "numeric" }));
  const server = field(_("Server"), row, _("The CUPS server the jobs are printed to over IPP, as computers print to it; ipps over TLS."));
  const q = input("f-queue", queue, { placeholder: _("a printer or a class there: laser-3f"), required: true, autocomplete: "off" });
  q.dataset.path = path;
  q.addEventListener("input", () => (q.dataset.path = q.value.includes("/") ? "" : found.get(q.value.trim()) || "/printers/"));
  const found = new Map(); // queue name -> its path, as the server listed it
  const list = el("div", undefined, { className: "flex flex-col gap-1" });
  const find = el("button", _("Find queues"), { type: "button", className: "btn btn-sm w-fit" });
  find.addEventListener("click", async () => {
    list.replaceChildren(el("span", _("Asking the server …"), { className: "text-xs text-base-content/60" }));
    try {
      const qs = await api("POST", "/api/kinds/cups-out/discover", { url: `${$("f-cscheme").value}://${$("f-chost").value.trim()}:${$("f-cport").value.trim() || 631}`, credential: $("f-credential").value.trim() });
      list.replaceChildren();
      if (!qs.length) list.append(el("span", _("The server lists no queue."), { className: "text-xs text-base-content/60" }));
      for (const x of qs) {
        found.set(x.name, x.path.slice(0, x.path.length - x.name.length));
        const b = el("button", undefined, { type: "button", className: "btn btn-xs btn-ghost h-auto justify-start py-1 text-left font-normal" });
        const about = `${x.info || x.make_model || ""} · ${x.state}`;
        b.append(el("span", x.name, { className: "font-mono font-semibold" }), el("span", " " + (x.class ? _("class · {about}", { about }) : about) + (x.accepting ? "" : _(", not accepting jobs")), { className: "text-base-content/60" }));
        b.addEventListener("click", () => {
          q.value = x.name;
          q.dataset.path = found.get(x.name);
          update();
        });
        list.append(b);
      }
    } catch (e) {
      list.replaceChildren(el("span", e.message, { className: "text-xs text-error" }));
    }
  });
  const box = el("div", undefined, { className: "flex flex-col gap-1" });
  box.append(q, find, list);
  const queueField = field(_("Queue"), box, _("Its name on the server; Find queues lists them (the server, and the credential below if it asks for one, first)."));
  const preview = el("p", "", { className: "font-mono text-xs text-base-content/60 break-all" });
  const update = () => (preview.textContent = (pickedURL(form) ?? cupsQueueURL(form)) || "—");
  connectFields(form, n, [server, queueField], (url) => {
    const x = url && parsePrinterURL(url);
    const m = x && /^\/(printers|classes)\/([^/]+)$/.exec(x.path);
    if (m) {
      const f = (id) => form.querySelector("#" + id);
      f("f-cscheme").value = x.scheme === "ipps" ? "ipps" : "ipp";
      f("f-chost").value = x.host;
      f("f-cport").value = x.port;
      q.value = m[2];
      q.dataset.path = "/" + m[1] + "/";
    }
    update();
  });
  form.append(server, queueField);
  form.append(field(_("Credential"), input("f-credential", n.credential || "", { placeholder: _("a name: the secret is DOCFLOW_CRED_<name> (user:password)"), autocomplete: "off" }), _("None: the server asks for no user and password.")));
  form.append(preview);
  setTimeout(update);
  for (const x of [row, q]) x.addEventListener("input", update);
  row.addEventListener("change", update);
}

// cupsQueueURL is the URL of the queue of a CUPS server (sink)'s fields
// in its form.
function cupsQueueURL(form) {
  const f = (id) => form.querySelector("#" + id);
  const host = f("f-chost").value.trim();
  const q = f("f-queue");
  const name = q.value.trim().replace(/^\/+/, "");
  if (!host || !name) return "";
  return `${f("f-cscheme").value}://${host}:${f("f-cport").value.trim() || 631}${q.dataset.path || "/"}${name}`;
}

// AT_NODE are the sections of what is at a node now, by its kind.
const AT_NODE = { sink: keptSection, display: boardSection, inbox: heldSection, archive: archivedSection, download: listedSection };
// at is the section of the node in the sidebar: {name, kind, box}.
let at = null;

// atSection is the section of what is at a node now, or null: the one
// shown already while the sidebar is of the same node, read again (what
// was chosen in it stays), else a new one.
function atSection(n) {
  const make = AT_NODE[n.kind];
  if (!make) return null;
  if (at && at.name === n.name && at.kind === n.kind) at.box.refresh();
  else at = { name: n.name, kind: n.kind, box: make(n.name) };
  return at.box;
}

// Looking: the sidebar of a node is its counts, what is at it now, a form
// to print there when the person may, and the jobs there they see.

function countsOf(name) {
  const c = (live && live.nodes[name]) || {};
  const out = [];
  if (c.open) out.push([String(c.open), _("open")]);
  if (c.held) out.push([String(c.held), _("held")]);
  if (c.pending) out.push([String(c.pending), _("waiting")]);
  if (c.failed) out.push([String(c.failed), _("failed")]);
  out.push([String(c.recent || 0), _("in the last hour")]);
  return out;
}

function viewSidebar() {
  const aside = $("sidebar");
  const n = topo.find((x) => x.name === selected);
  if (!n) {
    const p = el("p", undefined, { className: "text-sm text-base-content/60" });
    p.append(sentences(_("Click a node to see the jobs sent there that you may see, and to print there when you may."), live && Object.keys(live.nodes).length === 0 && _("Nothing of yours is in the network right now.")));
    aside.replaceChildren(p, el("a", _("The manual: this page, and what each node does"), { className: "link mt-2 block text-sm", href: "/help/topology#view", target: "_blank", rel: "noopener" }));
    return;
  }
  const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  head.append(kindIcon(n.kind, 20), kindBadge(n.kind), el("h3", n.title && n.title !== n.name ? n.title : n.name, { className: "text-lg font-semibold" }));
  if (n.title && n.title !== n.name) head.append(el("span", n.name, { className: "font-mono text-sm text-base-content/60" }));
  const stats = el("div", undefined, { className: "stats stats-horizontal flex-wrap bg-base-200" });
  for (const [v, label] of countsOf(n.name)) {
    const s = el("div", undefined, { className: "stat px-3 py-2" });
    s.append(el("div", v, { className: "stat-value text-xl" }), el("div", label, { className: "stat-title text-xs" }));
    stats.append(s);
  }
  const parts = [head];
  if (n.description) parts.push(el("p", n.description, { className: "whitespace-pre-wrap text-sm text-base-content/80" }));
  parts.push(stats);
  const out = [];
  for (const to of n.to || []) out.push(to);
  for (const [i, r] of (n.rules || []).entries()) out.push(_("{nodes} (rule {n})", { nodes: (r.to || []).join(", "), n: i + 1 }));
  for (const to of n.otherwise || []) out.push(_("{node} (otherwise)", { node: to }));
  for (const [i, c] of (n.cases || []).entries()) for (const to of c.to || []) out.push(`${to} (${c.name || "case " + (i + 1)})`);
  for (const to of n.default || []) out.push(`${to} (default)`);
  for (const to of n.original || []) out.push(_("{node} (original)", { node: to }));
  for (const m of n.members || []) out.push(n.policy === "failover" ? _("{node} (first)", { node: m }) : _("{node} (in turn)", { node: m }));
  if (n.fallback && n.fallback.printer) out.push(_("{node} (fallback)", { node: n.fallback.printer }));
  if (out.length) parts.push(el("p", "→ " + out.join(", "), { className: "font-mono text-xs text-base-content/60" }));
  const k = kindsLine(n);
  if (k) parts.push(k);
  const now = atSection(n);
  if (now) parts.push(now);
  if (n.kind === "cups") parts.push(printerSetup(n));
  const dest = dests.find((d) => d.name === n.name);
  if (dest) {
    const form = el("form", undefined, { className: "flex flex-col gap-2 border-t border-base-300 pt-3" });
    form.append(el("h4", _("Print here"), { className: "font-semibold" }));
    const file = el("input", undefined, { type: "file", className: "file-input file-input-sm w-full", required: true });
    const name = el("input", undefined, { className: "input input-sm w-full", placeholder: _("job name (the file's)") });
    form.append(file, name, el("button", _("Print"), { type: "submit", className: "btn btn-sm btn-primary w-fit" }));
    form.addEventListener("submit", async (ev) => {
      ev.preventDefault();
      const f = file.files[0];
      if (!f) return;
      try {
        const j = await api("POST", `/api/destinations/${encodeURIComponent(n.name)}/jobs?name=${encodeURIComponent(name.value || f.name)}`, f, true);
        status(_("Job #{n} is {state}", { n: j.number, state: STATES[j.state] || j.state }));
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
    // documents come in at a source: a person prints at one of them
    parts.push(el("p", ROLE_OF[n.kind] === "source" ? _("You may not print here.") : _("Documents come in at a source: print at an uploader or a CUPS server."), { className: "text-xs text-base-content/60" }));
  }
  const mine = jobs.filter((j) => j.dest === n.name || (j.legs || []).some((l) => l.node === n.name));
  const sec = el("div", undefined, { className: "flex flex-col gap-2 border-t border-base-300 pt-3" });
  sec.append(el("h4", _("Jobs here ({n})", { n: mine.length }), { className: "font-semibold" }));
  if (!mine.length) sec.append(el("p", _("None of yours."), { className: "text-xs text-base-content/60" }));
  const list = el("ul", undefined, { className: "flex flex-col gap-2" });
  for (const j of mine.slice(0, 20)) {
    const li = el("li", undefined, { className: "rounded-box bg-base-200 p-2 text-sm" });
    const row = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
    row.append(el("span", "#" + j.number, { className: "font-mono" }), el("span", j.name || "", { className: "font-semibold" }), stateBadge(j.state));
    li.append(row);
    li.append(el("div", `${sender(j)} · ${since(j.received)}`, { className: "text-xs text-base-content/60" }));
    li.append(el("div", legsOf(j), { className: "font-mono text-xs text-base-content/60 break-all", title: legsOf(j) }));
    const links = el("div", undefined, { className: "flex gap-3" });
    links.append(el("a", _("Document"), { className: "link text-xs", href: `/api/jobs/${j.id}/document` }));
    if (dest) {
      const re = el("button", _("Reprint here"), { type: "button", className: "link text-xs" });
      re.addEventListener("click", async () => {
        try {
          const again = await api("POST", `/api/jobs/${j.id}/reprint`, { to: n.name });
          status(_("Job #{n} is {state}", { n: again.number, state: STATES[again.state] || again.state }));
          etag = "";
          await poll();
        } catch (e) {
          status(e.message, true);
        }
      });
      links.append(re);
    }
    li.append(links);
    list.append(li);
  }
  sec.append(list);
  parts.push(sec);
  aside.replaceChildren(...parts);
}

// tab is the sidebar's tab of a saved node: its settings, or its usage
// (how computers, applications and devices reach it, and what is at it).
let tab = "settings";

// A form left for another node is kept: a draft's always, a saved
// node's when its settings were changed (edits). Opened again, it is made
// of the node as it is now with what was changed in it (the fields that
// differ from the node it was made of: an edge dragged meanwhile is not
// undone), and every required field left empty is marked.

// snapshot is what a form holds: the node's settings as read from it
// (formOf, lenient), the fields as typed (by id), and the number of
// required fields left empty.
function snapshot(n, form) {
  const raw = {};
  for (const e of form.querySelectorAll("input[id], select[id], textarea[id]")) raw[e.id] = e.type === "checkbox" ? e.checked : e.value;
  const q = form.querySelector("#f-queue");
  if (q) raw["f-queue@path"] = q.dataset.path || "";
  return { node: formOf(n, form, true), raw, missing: missingOf(n, form).reduce((c, m) => c + m.count, 0) };
}

const differs = (a, b) => JSON.stringify([a.node, a.raw]) !== JSON.stringify([b.node, b.raw]);

// keepShown keeps the form in the sidebar before it goes.
function keepShown() {
  if (!shown || !shown.form) return;
  const n = shown.draft || topo.find((x) => x.name === shown.name);
  if (!n) return;
  const kept = snapshot(n, shown.form);
  if (shown.draft) {
    if (drafts.includes(shown.draft)) shown.draft.kept = kept;
  } else if (differs(kept, shown.base)) edits.set(shown.name, { ...kept, base: shown.base });
  else edits.delete(shown.name);
}

// restoredForm is the form of a node as it was kept: made of the node
// with the settings changed in it, then each field changed as it was
// typed (and dispatched, for the checks and the previews to follow).
function restoredForm(n, kept) {
  const base = kept.base;
  const node = { ...n };
  for (const [k, v] of Object.entries(kept.node)) if (!base || JSON.stringify(v) !== JSON.stringify(base.node[k])) node[k] = v;
  if (base) for (const k of Object.keys(base.node)) if (!(k in kept.node)) delete node[k];
  const form = settingsForm(node);
  for (const [key, v] of Object.entries(kept.raw)) {
    if (base && base.raw[key] === v) continue;
    const [id, attr] = key.split("@");
    const e = form.querySelector("#" + id);
    if (!e) continue;
    if (attr) e.dataset[attr] = v;
    else if (e.type === "checkbox" ? e.checked !== v : e.value !== v) {
      if (e.type === "checkbox") e.checked = v;
      else e.value = v;
      e.dispatchEvent(new Event(e.tagName === "SELECT" || e.type === "checkbox" ? "change" : "input", { bubbles: true }));
    }
  }
  return form;
}

// sidebarForm is the form of the node or the draft opened in the
// sidebar, as it was kept or as the node makes it, with the required
// fields marked; the form there before is kept first.
function sidebarForm(n) {
  const same = shown && (draft ? shown.draft === draft : !shown.draft && shown.name === n.name);
  keepShown();
  const kept = draft ? draft.kept : edits.get(n.name);
  let form = settingsForm(n);
  const base = draft ? null : snapshot(n, form);
  if (kept) form = restoredForm(n, kept);
  // the same node made again (the topology read again) keeps what was
  // marked; one opened again shows what was kept of it
  const restored = !draft && !!kept && (same ? shown.restored : true);
  shown = { name: draft ? null : n.name, draft, form, base, touched: same ? shown.touched : new Set(), all: same ? shown.all : !!kept, restored, dirty: false, marked: 0 };
  if (shown.restored) {
    const note = el("div", undefined, { role: "status", className: "alert alert-warning alert-soft mb-1 flex flex-wrap items-center gap-2 px-3 py-2 text-sm" });
    const revert = el("button", _("Revert"), { type: "button", className: "btn btn-xs" });
    revert.addEventListener("click", () => {
      edits.delete(n.name);
      shown = null;
      renderSidebar();
      canvas.draw();
    });
    note.append(el("span", _("Changed, not saved: Save keeps the changes, Revert goes back to the node as it is."), { className: "flex-1" }), revert);
    form.prepend(note);
  }
  // a field left is marked when it is required and empty; what is not
  // saved is drawn on the canvas as it is typed
  form.addEventListener("focusout", (ev) => {
    if (!shown || shown.form !== form) return;
    if (ev.target.id) shown.touched.add(ev.target.id);
    refreshShown();
  });
  form.addEventListener("input", refreshShown);
  form.addEventListener("change", refreshShown);
  // a row added or removed (a rule, a check, a key) changes the settings
  form.addEventListener("click", (ev) => {
    if (ev.target.closest('button[type="button"]')) setTimeout(refreshShown);
  });
  refreshShown();
  return form;
}

// refreshShown marks the required fields of the form in the sidebar and
// tells the canvas what is not saved there.
function refreshShown() {
  if (!shown) return;
  const n = shown.draft || topo.find((x) => x.name === shown.name);
  if (!n) return;
  shown.marked = markRequired(n, shown.form);
  if (!shown.draft) shown.dirty = differs(snapshot(n, shown.form), shown.base);
  canvas.draw();
}

// Required fields: a field of the attribute required (the kind's Check
// refuses a node without it), and the settings (JSON) of an adapter's
// kind that its Check requires (GET /api/kinds). One left empty is
// marked in the warning color, with what it needs under it, once it was
// left (blur), or every one when the form is sent or opened again after
// it was left; a value that is wrong stays the error's red (liveCheck).

// missingOf are the required fields of a form left empty: {box, text,
// count}, with what each needs; the settings (JSON) of an adapter are
// one box, of the keys it lacks.
function missingOf(n, form) {
  const out = [];
  for (const box of form.querySelectorAll("[required]")) if (!box.disabled && !box.value.trim()) out.push({ box, text: box.dataset.need || _("Required: the node is not added or saved without it."), count: 1 });
  const box = form.querySelector("#f-settings");
  const keys = box ? missingSettings(n.kind, box.value) : [];
  if (keys.length) out.push({ box, text: _("Required in the settings: {keys}.", { keys: keys.join(", ") }), count: keys.length });
  return out;
}

// missingSettings are the keys of a kind's settings (JSON) its Check
// requires and the JSON lacks (of those whose if it fits); none when it
// is not JSON, which is wrong instead.
function missingSettings(kind, text) {
  let v;
  try {
    v = JSON.parse(text || "{}");
  } catch {
    return [];
  }
  if (!v || typeof v !== "object" || Array.isArray(v)) return [];
  const empty = (x) => x === undefined || x === null || (typeof x === "string" && !x.trim()) || (Array.isArray(x) && !x.length);
  const fits = (r) => Object.entries(r.if || {}).every(([k, want]) => String(v[k] ?? "").trim().toLowerCase() === want);
  return ((kinds[kind] || {}).required || []).filter((r) => fits(r) && empty(v[r.key])).map((r) => r.key);
}

// WARN is the class of a field marked required and empty, by its tag.
const WARN = { INPUT: "input-warning", SELECT: "select-warning", TEXTAREA: "textarea-warning" };

// markRequired marks the required fields of a form left empty that are
// to be marked (left, or all), and the settings (JSON) that do not parse;
// it returns how many are marked.
function markRequired(n, form) {
  const missing = new Map(missingOf(n, form).map((m) => [m.box, m]));
  const json = form.querySelector("#f-settings");
  let marked = 0;
  for (const box of form.querySelectorAll("[required], #f-settings")) {
    const due = shown.all || shown.touched.has(box.id);
    let m = due ? missing.get(box) : null;
    let wrong = "";
    if (box === json && due) {
      try {
        JSON.parse(json.value || "{}");
      } catch (e) {
        wrong = _("Settings: {why}", { why: e.message });
      }
    }
    box.classList.toggle(WARN[box.tagName], !!m && !wrong);
    if (box === json) box.classList.toggle("textarea-error", !!wrong);
    let msg = form.querySelector(`[data-need-of="${box.id}"]`);
    if ((m || wrong) && !msg) {
      // under the field, or under the row it is in
      msg = el("p", "", { id: box.id + "-need", className: "label whitespace-normal" });
      msg.dataset.needOf = box.id;
      (box.parentElement.classList.contains("grid") ? box.parentElement : box).after(msg);
    }
    if (msg) {
      msg.textContent = wrong || (m ? m.text : "");
      msg.className = "label whitespace-normal " + (wrong ? "text-error" : "text-warning");
      msg.hidden = !msg.textContent;
    }
    if (m || wrong) {
      box.setAttribute("aria-invalid", "true");
      box.setAttribute("aria-describedby", msg.id);
      if (m) marked += m.count;
    } else if (msg && box.getAttribute("aria-describedby") === msg.id) {
      box.removeAttribute("aria-describedby");
      box.setAttribute("aria-invalid", box.validity.customError ? "true" : "false");
    }
  }
  return marked;
}

// legendOf is the legend of the field a box is in.
const legendOf = (box) => box.closest("fieldset")?.querySelector("legend")?.textContent || box.id;

function renderSidebar() {
  if (!editable) return viewSidebar();
  const aside = $("sidebar");
  const n = draft || topo.find((x) => x.name === selected);
  if (!n) {
    keepShown();
    shown = null;
    aside.replaceChildren(
      el("p", _("Click a node to edit it here, or add one with “+ Add a node”: it is placed at once, and connected when the nodes after it are there. Drag from a node's right port onto another node to connect them; a display's port sets its fallback printer, a branch has a port for each case and one for default, an assert one for what passes and one for what fails. A node's title bar is colored by its role (source, control flow, filter, transform, sink), with an icon of its kind. The ports and the edges are colored by what flows: blue a .bdf document, green data (for a template node), amber any file as received (for a convert node), purple several kinds. A queue, a class, a filter or an inbox is typed by its first edge, on either side: in front of a printer it carries .bdf, and takes .bdf alone; one with no edge yet is hollow. While you drag an edge, the nodes that would not take it fade and the wire turns red over them. An edge drawn dashed and red no longer connects: a change of settings left it carrying what the node after it does not take."), { className: "text-sm text-base-content/60" }),
      el("p", _("A node left before it is added stays on the canvas, dashed, with what was typed for it; a node changed and not saved is ringed in amber. Click either to go on with it. A required field left empty is marked in amber once it was left, and a badge on the canvas counts them."), { className: "mt-2 text-sm text-base-content/60" }),
      el("a", _("The manual: the editor, and every kind of node"), { className: "link mt-2 block text-sm", href: "/help/topology#editor", target: "_blank", rel: "noopener" }),
    );
    return;
  }
  const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  head.append(kindIcon(n.kind, 20), kindBadge(n.kind));
  if (draft) head.append(el("span", _("new"), { className: "badge badge-outline" }));
  else {
    head.append(el("h3", n.title || n.name, { className: "text-lg font-semibold" }));
    if (n.title && n.title !== n.name) head.append(el("span", n.name, { className: "font-mono text-sm text-base-content/60" }));
    head.append(scopeBadge(n));
  }
  // the manual's section of the kind (in the language of the pages)
  head.append(el("a", _("Help"), { className: "link link-hover ml-auto text-xs text-base-content/60", href: "/help/nodes#" + n.kind, target: "_blank", rel: "noopener", title: _("The manual: {kind} nodes", { kind: KINDS[n.kind] || n.kind }) }));
  const form = sidebarForm(n);
  if (draft) {
    aside.replaceChildren(head, form);
    return;
  }
  const settings = el("div");
  settings.append(form);
  if (n.kind === "template") {
    const tpls = section(_("Templates"));
    settings.append(tpls);
    templatesSection(n, tpls).catch((e) => status(e.message, true));
  }
  const usage = el("div", undefined, { className: "flex flex-col gap-2" });
  usagePanel(n, usage);
  const panels = { settings, usage };
  const tabs = el("div", undefined, { role: "tablist", className: "tabs tabs-border mt-1 mb-2" });
  const showTab = () => {
    for (const t of tabs.children) {
      t.classList.toggle("tab-active", t.dataset.tab === tab);
      t.setAttribute("aria-selected", String(t.dataset.tab === tab));
    }
    for (const [id, p] of Object.entries(panels)) p.style.display = id === tab ? "" : "none";
  };
  for (const [id, label] of [["settings", _("Settings")], ["usage", _("Usage")]]) {
    const t = el("button", label, { type: "button", role: "tab", className: "tab gap-1" });
    t.dataset.tab = id;
    // what waits at an inbox or a display, worked under Usage
    const waiting = id === "usage" ? (n.kind === "inbox" ? n.held : n.kind === "display" ? n.open : 0) : 0;
    if (waiting) t.append(el("span", String(waiting), { className: "badge badge-xs badge-warning", title: n.kind === "inbox" ? _("{n} held", { n: waiting }) : _("{n} open", { n: waiting }) }));
    t.addEventListener("click", () => {
      tab = id;
      showTab();
    });
    tabs.append(t);
  }
  showTab();
  aside.replaceChildren(head, tabs, settings, usage);
}

// NAME_RE is the shape of a node's name, as the server takes it (nameRE,
// model.go): it is in the URLs, and the records are kept under it.
const NAME_RE = /^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$/;

// nameError tells what is wrong with a name typed for a node (self: the
// node's name now; "" for one being added), or "" when it is fine.
function nameError(v, self) {
  if (!v) return _("A name is needed: other nodes and the URLs call the node by it.");
  if (/\s/.test(v)) return _("No spaces: letters, digits, . _ and - only.");
  const bad = [...v].find((c) => !/[A-Za-z0-9_.-]/.test(c));
  if (bad) return _("“{char}” is not taken: letters, digits, . _ and - only.", { char: bad });
  if (!NAME_RE.test(v)) return v.length > 64 ? _("64 characters at most ({n} now).", { n: v.length }) : _("It starts with a letter or a digit.");
  const other = topo.find((x) => x.name === v && x.name !== self);
  if (other) return _("Taken: another node ({other}) is named {name}.", { other: (KINDS[other.kind] || other.kind) + (other.title ? ", " + other.title : ""), name: v });
  return "";
}

// renameNote says what renaming a node changes: where it is reached.
function renameNote(n, to) {
  const at = [];
  if (n.kind === "cups") at.push(_("/printers/{name} on its port", { name: to }), `/ipp/print/${to}`);
  if (n.api) at.push(`/api/${n.kind}/${to}`);
  const where = at.length && _("It is reached at {where} from then on: computers and applications set up with {was} must be set up again.", { where: at.join(", "), was: n.name });
  return sentences(_("Saving renames it to {to}.", { to }), where, _("The nodes that send to it follow; the jobs that reached it keep {was}.", { was: n.name }));
}

// liveCheck checks a field as it is typed: it turns red, with what is
// wrong under it, and the form is not sent until it is right (test gives
// the error, or ""). note gives what is said under it meanwhile, and
// whether it is a warning: [text, warn]. A field left empty is not
// wrong: a required one is marked as such (markRequired).
function liveCheck(box, test, note) {
  const msg = el("p", "", { className: "label whitespace-normal" });
  let touched = !!box.value;
  const run = () => {
    const err = touched && box.value ? test(box.value) : "";
    box.classList.toggle("input-error", !!err);
    box.setCustomValidity(err);
    box.setAttribute("aria-invalid", err ? "true" : "false");
    const [text, warn] = err ? [err, false] : note(box.value);
    msg.textContent = text;
    msg.className = "label whitespace-normal" + (err ? " text-error" : warn ? " text-warning" : "");
  };
  box.addEventListener("input", () => {
    touched = true;
    run();
  });
  run();
  return msg;
}

// nameField is a node's name, checked as it is typed. Changing a saved
// node's name renames it when it is saved (POST /api/nodes/<name>/rename).
function nameField(n) {
  const self = draft ? "" : n.name;
  const box = input("f-name", self, { placeholder: "intake, both, tokyo, review, wh-a, laser, audit", required: true, autocomplete: "off", spellcheck: false, maxLength: 80 });
  box.dataset.need = _("A name is needed: other nodes and the URLs call the node by it.");
  const f = field(_("Name"), box);
  f.append(liveCheck(box, (v) => nameError(v, self), (v) => (!draft && v !== n.name ? [renameNote(n, v), true] : [_("The system's name: letters, digits, . _ and -, no spaces. Other nodes and the URLs call the node by it; the title is what people read."), false])));
  return f;
}

// nextPort is the first port from 6310 that no node has: the server gives
// a CUPS server put without a port the same.
function nextPort(self) {
  const used = new Set(topo.filter((x) => x.name !== self).map((x) => x.port).filter(Boolean));
  let p = 6310;
  while (used.has(p)) p++;
  return p;
}

// portField is the port of a CUPS server of its own (its IPP; what
// applications call is at the API, on the manager's port): a node being
// added gets the first port free from 6310; none where this process has
// no ports (serverless).
function portField(n, what, empty) {
  const none = me && me.ports === false;
  const self = draft ? "" : n.name;
  const box = input("f-port", n.port || (draft && !none ? nextPort(self) : ""), { placeholder: none ? "" : empty, inputMode: "numeric", disabled: none, autocomplete: "off" });
  const f = field(_("Port (this process only)"), box);
  f.append(liveCheck(box, (v) => {
    if (!v.trim()) return "";
    const p = Number(v);
    if (!/^\d+$/.test(v.trim()) || p < 1024 || p > 65535) return _("A port is a number from 1024 to 65535.");
    const other = topo.find((x) => x.port === p && x.name !== self);
    return other ? _("Port {port} is {node}'s: {free} is free.", { port: p, node: other.name, free: nextPort(self) }) : "";
  }, () => [none ? _("Not available here (serverless): no node has a port of its own.") : what, false]));
  return f;
}

// appsField is the apps a node lists, which use it with their own
// credentials: an entry's (an uploader, a CUPS server) send to it, an
// inbox's work it.
function appsField(n, what) {
  const boxes = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
  for (const a of apps) {
    if (a.public) continue; // a frontend uses the API as a person
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "app", value: a.id, checked: (n.apps || []).includes(a.id) }), a.id + (a.title ? " · " + a.title : ""));
    boxes.append(l);
  }
  if (!boxes.childElementCount) boxes.append(el("a", _("No apps yet: register one"), { className: "link text-sm", href: "/admin/apps" }));
  return field(_("Apps allowed"), boxes, what);
}

// settingsForm is the form of a node's settings: its title and name, a
// description, and the fields of its kind.
function settingsForm(n) {
  const kind = n.kind;
  // the form marks the required fields itself (markRequired): no
  // bubbles of the browser's
  const form = el("form", undefined, { className: "flex flex-col gap-1", noValidate: true });
  const title = input("f-title", n.title || "", { placeholder: _("what people call it: 3F laser, invoices of the ERP …"), maxLength: 255 });
  form.append(field(_("Title"), title, _("Shown on the node, in any characters.")));
  form.append(nameField(n));
  const note = el("textarea", undefined, { id: "f-description", className: "textarea textarea-sm w-full", rows: 3, value: n.description || "", placeholder: _("what it is for: the invoices of the ERP, PDF only; printed on the 3rd floor …") });
  form.append(field(_("Description"), note, _("A note shown in the node, under its [kind], as in a C4 diagram: the first lines on the canvas, all of it here.")));
  if (!draft) {
    const k = kindsLine(n);
    if (k) form.append(k);
  }
  // nothing fails at a branch: it sends on, or drops; an assert's
  // failures are with its checks
  if (kind !== "download" && kind !== "branch" && kind !== "assert") form.append(field(_("On error"), picker("on_error", n.on_error || [], n.name), _("Where the job goes, with what reached this node, when its send or step fails for good (a notify node); none chosen: the manager's default, or the job stops here.")));
  if (ADAPTER_SETTINGS[kind]) {
    const box = el("textarea", undefined, { id: "f-settings", className: "textarea textarea-sm w-full font-mono", rows: 8, value: JSON.stringify(settingsOf(n, kind), null, 2) });
    form.append(field(_("Settings (JSON)"), box, sentences(ADAPTER_SETTINGS[kind], _("Credentials are names: the secret is the server's DOCFLOW_CRED_<name>, or a file. Text, recipients, paths and headers may use {{meta.<key>}}."))));
  }
  const later = (s) => sentences(s, _("None yet is fine: connect it when the nodes after it are there (a job that reaches it meanwhile gets an error)."));
  switch (kind) {
    case "upload": {
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), later(_("Every node chosen gets the job. The port on the canvas adds one too."))));
      entryFields(form, n);
      form.append(metaField(n));
      form.append(appsField(n, _("Apps send here with their token (POST /oauth/token), within their scopes; the uploader's own key works as well.")));
      break;
    }
    case "s3-trigger":
      form.append(field(_("Bucket"), input("f-url", n.url || "", { placeholder: _("s3://bucket?region=ap-northeast-1 (&endpoint=… for MinIO, R2)"), required: true, autocomplete: "off" }), _("The bucket whose objects it takes. It does not listen to the bucket: an S3 event's reference is POSTed to it (Usage).")));
      form.append(field(_("Credential"), input("f-credential", n.credential || "", { placeholder: _("a name: the secret is DOCFLOW_CRED_<name> (access-key:secret-key)"), autocomplete: "off" }), _("None: AWS's own (the environment, the role).")));
      form.append(field(_("Prefix"), input("f-path", n.path || "", { placeholder: _("incoming/ (the keys it takes are under it)"), autocomplete: "off" }), _("A key elsewhere, or another bucket, is refused (403).")));
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), later(_("Every node chosen gets the object read."))));
      entryFields(form, n);
      form.append(metaField(n));
      break;
    case "cups":
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), later(_("Every node chosen gets what computers print here, made .bdf: a printer makes it a relay, a sink shows the pages, a filter keeps only what fits. The port on the canvas adds one too."))));
      form.append(metaField(n));
      form.append(portField(n, _("Computers print to ipp://<this host>:<port>/printers/{name}, the one printer of this CUPS server (Usage). Empty: the first port free from 6310.", { name: n.name || "<name>" }), _("empty: the first free from 6310")));
      form.append(field(_("TLS certificate (file)"), input("f-tls-cert", n.tls_cert || "", { placeholder: _("/etc/bdf/cups.crt; empty: plain ipp") })));
      form.append(field(_("TLS key (file)"), input("f-tls-key", n.tls_key || "", { placeholder: "/etc/bdf/cups.key" })));
      {
        // advertised unless it says not to (dnssd false), where nodes have ports
        const none = !!(me && me.ports === false);
        const bonjour = el("label", undefined, { className: "inline-flex items-center gap-2 text-sm" });
        bonjour.append(el("input", undefined, { id: "f-dnssd", type: "checkbox", className: "checkbox checkbox-sm", checked: !none && n.dnssd !== false, disabled: none }), _("Bonjour: advertised on the local network, where Macs, iPhones and iPads find it as an AirPrint printer, Android phones and tablets as a Mopria printer, and Windows and Linux computers as an IPP Everywhere printer"));
        form.append(bonjour);
      }
      form.append(appsField(n, _("Apps print here with their id and secret as IPP Basic, within their scopes; computers print as the submit tokens allow.")));
      break;
    case "sink":
      form.append(field(_("Latest"), input("f-latest", n.latest || "", { placeholder: "10", inputMode: "numeric" }), _("How many of the latest documents it keeps; an older one goes when a newer comes. Their first pages show in Usage and on the canvas.")));
      break;
    case "download": {
      form.append(field(_("Owner's key"), input("f-by", n.by || "", { placeholder: _("customer_id (a metadata key); empty: a person reads what they could see of the job"), autocomplete: "off" }), _("The metadata key whose value says whose a document is: a person reads the documents whose value is them, under the subject below. A backend reads every one.")));
      form.append(field(_("Subject"), selectInput("f-subject", Object.entries(SUBJECTS), n.subject || "email"), _("What of a person is compared with the owner's value.")));
      form.append(field(_("Fields"), input("f-fields", (n.fields || []).join(", "), { placeholder: _("invoice_no, doc_type (metadata keys, comma-separated)") }), _("Metadata copied into the entries listed, for the reader to show or pick by.")));
      form.append(field(_("Entries"), input("f-entries", n.entries || "", { placeholder: "1000", inputMode: "numeric" }), _("How many of the latest entries are listed for one owner, and for everyone; an older entry is still read by its job's number.")));
      const appBoxes = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
      for (const a of apps) {
        if (a.public) continue;
        const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
        l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "app", value: a.id, checked: (n.apps || []).includes(a.id) }), a.id + (a.title ? " · " + a.title : ""));
        appBoxes.append(l);
      }
      if (!appBoxes.childElementCount) appBoxes.append(el("a", _("No backend apps yet: register one"), { className: "link text-sm", href: "/admin/apps" }));
      form.append(field(_("Apps allowed"), appBoxes, _("Backends that read every entry with their token (the scope download). A frontend (a public client) reads as the person using it, with a token of theirs: it needs no listing here.")));
      break;
    }
    case "queue":
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), later(_("Every node chosen gets the job: a fan-out when there are several."))));
      break;
    case "template": {
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), later(_("The form filled (.bdf) goes to every node chosen."))));
      form.append(field(_("The data as it came to"), picker("original", n.original || [], n.name, n.reaches || []), _("What reached the node goes on as it was to every node chosen (an archive of originals, say), with the same trace as the document made: each archive's record names the other's object.")));
      const ids = n.templates || [];
      const options = [["", ids.length ? _("— named by the metadata below, or none yet —") : _("— none yet: upload one once the node is added —")], ...ids.map((id) => [id, id])];
      if (n.template && !ids.includes(n.template)) options.push([n.template, _("{id} (not uploaded yet)", { id: n.template })]);
      form.append(field(_("Template"), selectInput("f-template", options, n.template || ""), draft ? _("Add the node first, then upload its templates (Settings, below) and choose one here.") : _("One of the templates uploaded below, filled with the data that reaches the node.")));
      form.append(field(_("Or the metadata key naming it"), input("f-template-from", n.template_from || "", { placeholder: _("doc_type: the sender's X-Meta-Doc-Type names the template") }), _("A job whose metadata lacks the key, or names a template that is not here, fails here (no retry); so does one while the node names no template.")));
      form.append(field(_("Name of a list"), input("f-list-name", n.list_name || "", { placeholder: _("the file's name without its extension (orders for orders.csv)") }), _("Data that is a list of records (a CSV file, a sheet, a JSON array) is the value of this name for the template; an object is the data itself.")));
      const strict = el("label", undefined, { className: "inline-flex items-center gap-2 text-sm" });
      strict.append(el("input", undefined, { id: "f-strict", type: "checkbox", className: "checkbox checkbox-sm", checked: !!n.strict }), _("Strict: refuse data that lacks a name the template looks up (else the tag is left empty, with a warning)"));
      form.append(strict);
      break;
    }
    case "convert":
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), later(_("The .bdf document made goes to every node chosen; a .bdf that reaches the node goes on as it is."))));
      form.append(field(_("The file as it came to"), picker("original", n.original || [], n.name, n.reaches || []), _("What reached the node goes on as it was to every node chosen (an archive of originals, say), with the same trace as the document made: each archive's record names the other's object.")));
      form.append(formatsField(n, _("Formats converted"), _("None checked: every format the converter reads. A file of another format fails here (no retry); JSON and YAML have no document of their own: fill a template with them.")));
      break;
    case "sign":
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), later(_("The signed .bdf goes to every node chosen."))));
      form.append(field(_("Role"), input("f-sign-role", n.sign_role || "", { placeholder: "witness" }), _("The role the signature says Docflow signs in. The signature also says the node, the job, who sent it (as Docflow verified it) and when it was received; bdf verify shows them. The key is the manager's (docflow.sign_key_file); without one the node fails its jobs.")));
      break;
    case "filter": {
      const rules = el("div", undefined, { id: "f-rules", className: "flex flex-col gap-2" });
      for (const r of n.rules || []) rules.append(ruleRow(r));
      const add = el("button", _("+ a rule"), { type: "button", className: "btn btn-xs w-fit" });
      add.addEventListener("click", () => rules.append(ruleRow()));
      const box = el("div", undefined, { className: "flex flex-col gap-2" });
      box.append(rules, add);
      form.append(field(_("Rules, in order"), box, _("The first rule that fits the job sends it to its nodes. name, user (the sender), format and the values of the metadata are regular expressions; pages are bounds; the kind of payload is .bdf, data or any other file. What reaches the filter goes on as it is.")));
      form.append(field(_("Otherwise"), picker("otherwise", n.otherwise || [], n.name), _("Where a job no rule fits goes; none: it ends here (dropped).")));
      break;
    }
    case "assert":
      form.append(field(_("Passes to"), picker("to", n.to || [], n.name), later(_("A job that passes every check goes on to every node chosen, as it reached the assert. The pass port on the canvas adds one too."))));
      form.append(field(_("Checks"), checksBox(n.checks || [], true), _("Every check must hold. The key is the sender's (X-Meta-*, ?meta.<key>=, IPP's ipp.<attribute>), the entry's, or the manager's ($pages, $format, $sender …); it must have a value unless optional. Then, those given: the value matches the regular expression (RE2, anywhere: ^…$ for all of it), is a number from min to max, and is one of the values listed. Right after an entry (or a queue a job is sent to), a job that fails is refused to its sender (REST 422, IPP client-error-bad-request) with what failed; further on, it is aborted here.")));
      form.append(field(_("On fail"), picker("on_error", n.on_error || [], n.name), _("Where a job that fails a check goes, with what reached the assert and what failed ({{meta.$error}}): a notify node, say. None chosen: the manager's default, or the job stops here. The fail port on the canvas adds one too.")));
      break;
    case "branch": {
      const cases = el("div", undefined, { id: "f-cases", className: "flex flex-col gap-2" });
      for (const c of n.cases || []) cases.append(caseRow(c));
      const add = el("button", _("+ a case"), { type: "button", className: "btn btn-xs w-fit" });
      add.addEventListener("click", () => cases.append(caseRow()));
      const box = el("div", undefined, { className: "flex flex-col gap-2" });
      box.append(cases, add);
      form.append(field(_("Cases, in order"), box, _("As a switch: the first case whose checks all hold sends the job to its nodes; the others are not tried. Each case is a port of its own on the canvas, labelled with its name: drag from it to add a node. What reaches the branch goes on as it is.")));
      form.append(field(_("Default"), picker("default", n.default || [], n.name), _("Where a job no case fits goes; none: it ends here (dropped). The default port on the canvas adds one too.")));
      break;
    }
    case "class":
      form.append(field(_("Members"), input("f-members", (n.members || []).join(", "), { placeholder: "mfp-3f, mfp-4f", list: "node-names" }), _("None yet is fine: add the members when they are there (a job that reaches the class meanwhile gets an error).")));
      form.append(field(_("Policy"), selectInput("f-policy", [["round-robin", _("in turn (round-robin)")], ["failover", _("the first member (failover)")]], n.policy || "round-robin"), _("A printer that does not take the job is skipped for the next member.")));
      break;
    case "inbox":
      form.append(field(_("Releases to"), picker("to", n.to || [], n.name), _("The nodes an item may be released to; none chosen: any node, which only the role topology releases to. A downloader chosen is not one to release to: it lists each job while the inbox holds it, for applications to fetch.")));
      form.append(field(_("Keep"), input("f-keep", n.keep || "", { placeholder: "14d" }), _("How long released and discarded items stay shown.")));
      form.append(appsField(n, _("Apps that work the inbox with their token (the scope inbox): they list what it holds, see its pages, release and discard. The people who work it are its scope's, below.")));
      break;
    case "archive":
      form.append(field(_("Forwards to"), picker("to", n.to || [], n.name), _("Once written, the job goes on to every node chosen: a downloader after the archive lists the document for applications and people to fetch. The port on the canvas adds one too.")));
      form.append(field(_("Store"), input("f-store", n.store || "", { placeholder: _("s3://bucket?region=… or bbolt:path"), required: true }), _("A documents URL: bbolt, s3, r2 or gcs. Each job's .bdf and record are written once; the manager never changes or deletes them.")));
      form.append(field(_("Collection"), input("f-collection", n.collection || "", { placeholder: "archive-" + (n.name || "<name>") })));
      break;
    case "display": {
      const printers = topo.filter((x) => PRINTS.has(x.kind)).map((x) => [x.name, (x.title || x.name) + (x.title ? " · " + x.name : "")]);
      form.append(field(_("Fallback printer"), selectInput("f-fallback", [["", _("none")], ...printers], n.fallback ? n.fallback.printer : ""), _("Prints the slips nobody opened in time, and every unopened slip while no device is connected. The port on the canvas sets it too.")));
      form.append(field(_("Print after"), input("f-after", n.fallback && n.fallback.after ? n.fallback.after : "", { placeholder: _("10m; empty: only while no device is connected") })));
      form.append(field(_("Keep"), input("f-keep", n.keep || "", { placeholder: "14d" }), _("How long slips done or canceled stay on the devices.")));
      form.append(field(_("Offline after"), input("f-offline", n.offline_after || "", { placeholder: "2m" }), _("How long a device may be silent and still count as connected.")));
      break;
    }
    case "cups-out":
      cupsOutFields(form, n);
      break;
    case "printer": {
      const u = parsePrinterURL(n.url);
      const row = el("div", undefined, { className: "grid grid-cols-[6rem_1fr_5rem] gap-1" });
      row.append(selectInput("f-scheme", SCHEMES.map((s) => [s, s + "://"]), SCHEMES.includes(u.scheme) ? u.scheme : "ipp"), input("f-host", u.host, { placeholder: _("host or address"), required: true }), input("f-pport", u.port, { placeholder: _("port"), inputMode: "numeric" }));
      const address = field(_("Protocol and address"), row, _("ipp and ipps are IPP printers; socket is port 9100 (raw), lpd a line printer daemon."));
      const path = field(_("Path"), input("f-path", u.path, { placeholder: _("/ipp/print (ipp); a queue name (lpd)") }));
      const preview = el("p", "", { className: "font-mono text-xs text-base-content/60 break-all" });
      const update = () => (preview.textContent = (pickedURL(form) ?? composeURL(form)) || "—");
      connectFields(form, n, [address, path], (url) => {
        if (url) {
          const x = parsePrinterURL(url);
          const f = (id) => form.querySelector("#" + id);
          f("f-scheme").value = SCHEMES.includes(x.scheme) ? x.scheme : "ipp";
          f("f-host").value = x.host;
          f("f-pport").value = x.port;
          f("f-path").value = x.path;
        }
        update();
      });
      form.append(address, path, preview);
      setTimeout(update);
      for (const id of ["f-scheme", "f-host", "f-pport", "f-path"]) form.querySelector("#" + id).addEventListener("input", update);
      form.append(field(_("Route"), selectInput("f-route", [["direct", _("direct: the manager sends to it")], ["branch", _("branch: a bdf printd of the site fetches (not yet)")]], n.route || "direct")));
      form.append(field(_("Preview (a fake printer)"), input("f-preview", n.preview || "", { placeholder: "http://127.0.0.1:8650/api/printers/eve" }), _("The printer on a tools/fakeprinter farm's API: what it received is shown in Usage. For development and demonstrations.")));
      break;
    }
  }
  const scopes = {
    download: [["shared", _("shared: anyone signed in may read their own")], ["group", _("group: its members may read their own")], ["personal", _("personal: one person")]],
    inbox: [["shared", _("shared: anyone signed in may send; the role topology releases")], ["group", _("group: its members send, see the jobs, release and discard")], ["personal", _("personal: one person sends and releases")]],
  }[kind] || [["shared", _("shared: anyone signed in may send")], ["group", _("group: its members send and see the jobs")], ["personal", _("personal: one person")]];
  form.append(field(_("Scope"), selectInput("f-scope", scopes, n.scope || "shared")));
  form.append(field(_("Owner"), input("f-owner", n.owner || "", { placeholder: _("a group (a value of the groups claim), or a person's email") })));
  const actions = el("div", undefined, { className: "mt-2 flex flex-wrap gap-2" });
  actions.append(el("button", draft ? _("Add") : _("Save"), { type: "submit", className: "btn btn-sm btn-primary" }));
  if (draft) {
    const cancel = el("button", _("Cancel"), { type: "button", className: "btn btn-sm btn-ghost" });
    cancel.addEventListener("click", () => dropDraft(draft));
    actions.append(cancel);
  } else {
    const del = el("button", _("Delete"), { type: "button", className: "btn btn-sm btn-ghost text-error" });
    del.addEventListener("click", () => remove(n.name));
    actions.append(del);
  }
  form.append(actions);
  form.addEventListener("submit", save);
  return form;
}

// formOf reads the settings of a node from its form, as PUT /api/nodes
// takes them; JSON that does not parse and a bound that is not a number
// are refused (an error), or, lenient (what is kept of a form left),
// kept as they were typed.
function formOf(n, form, lenient) {
  const f = (id) => form.querySelector("#" + id);
  const body = { kind: n.kind, title: f("f-title").value.trim(), description: f("f-description").value.trim(), scope: f("f-scope").value };
  if (f("f-owner").value.trim()) body.owner = f("f-owner").value.trim();
  if (f("f-settings")) {
    try {
      Object.assign(body, JSON.parse(f("f-settings").value || "{}"));
    } catch (e) {
      if (!lenient) throw new Error(_("Settings: {why}", { why: e.message }));
    }
  }
  if (form.querySelector('input[name="on_error"]')) body.on_error = picked(form, "on_error");
  const dur = (id) => (f(id).value.trim() ? f(id).value.trim() : undefined);
  const entry = () => {
    body.to = picked(form, "to");
    body.accept = f("f-accept").value;
    body.output = f("f-output").value;
    body.formats = [...form.querySelectorAll('input[name="fmt"]:checked')].map((c) => c.value);
    body.meta = metaOf(form);
  };
  switch (n.kind) {
    case "upload":
      entry();
      body.apps = [...form.querySelectorAll('input[name="app"]:checked')].map((c) => c.value);
      break;
    case "s3-trigger":
      entry();
      body.url = f("f-url").value.trim();
      if (dur("f-credential")) body.credential = dur("f-credential");
      if (dur("f-path")) body.path = dur("f-path");
      break;
    case "cups":
      // what computers print is a document, made .bdf: no accept, output or formats
      body.to = picked(form, "to");
      body.meta = metaOf(form);
      body.apps = [...form.querySelectorAll('input[name="app"]:checked')].map((c) => c.value);
      if (dur("f-port")) body.port = Number(dur("f-port"));
      if (dur("f-tls-cert")) body.tls_cert = dur("f-tls-cert");
      if (dur("f-tls-key")) body.tls_key = dur("f-tls-key");
      if (!f("f-dnssd").checked && !f("f-dnssd").disabled) body.dnssd = false;
      break;
    case "template":
      body.to = picked(form, "to");
      if (f("f-template").value) body.template = f("f-template").value;
      if (dur("f-template-from")) body.template_from = dur("f-template-from");
      if (dur("f-list-name")) body.list_name = dur("f-list-name");
      if (f("f-strict").checked) body.strict = true;
      body.original = picked(form, "original");
      break;
    case "convert":
      body.to = picked(form, "to");
      body.original = picked(form, "original");
      body.formats = [...form.querySelectorAll('input[name="fmt"]:checked')].map((c) => c.value);
      break;
    case "sign":
      body.to = picked(form, "to");
      if (dur("f-sign-role")) body.sign_role = dur("f-sign-role");
      break;
    case "sink":
      if (dur("f-latest")) body.latest = Number(dur("f-latest"));
      break;
    case "download":
      if (dur("f-by")) body.by = dur("f-by");
      body.subject = f("f-subject").value;
      body.fields = names(f("f-fields").value);
      if (dur("f-entries")) body.entries = Number(dur("f-entries"));
      body.apps = [...form.querySelectorAll('input[name="app"]:checked')].map((c) => c.value);
      break;
    case "queue":
      body.to = picked(form, "to");
      break;
    case "filter":
      body.rules = rulesOf(form);
      body.otherwise = picked(form, "otherwise");
      break;
    case "assert":
    case "branch":
      if (n.kind === "assert") {
        body.to = picked(form, "to");
        body.checks = checksOf(f("f-checks"), lenient);
      } else {
        body.cases = casesOf(form, lenient);
        body.default = picked(form, "default");
      }
      break;
    case "class":
      body.members = names(f("f-members").value);
      body.policy = f("f-policy").value;
      break;
    case "inbox":
      body.to = picked(form, "to");
      body.keep = dur("f-keep");
      body.apps = [...form.querySelectorAll('input[name="app"]:checked')].map((c) => c.value);
      break;
    case "archive":
      body.to = picked(form, "to");
      body.store = f("f-store").value.trim();
      if (f("f-collection").value.trim()) body.collection = f("f-collection").value.trim();
      break;
    case "display":
      if (f("f-fallback").value) {
        body.fallback = { printer: f("f-fallback").value };
        if (dur("f-after")) body.fallback.after = dur("f-after");
      }
      body.keep = dur("f-keep");
      body.offline_after = dur("f-offline");
      break;
    case "cups-out":
      body.url = pickedURL(form) ?? cupsQueueURL(form);
      if (pickedURL(form) !== null && f("f-bonjour").value) body.bonjour = f("f-bonjour").value;
      if (dur("f-credential")) body.credential = dur("f-credential");
      break;
    case "printer":
      body.url = pickedURL(form) ?? composeURL(form);
      if (pickedURL(form) !== null && f("f-bonjour").value) body.bonjour = f("f-bonjour").value;
      body.route = f("f-route").value;
      if (dur("f-preview")) body.preview = dur("f-preview");
      break;
  }
  return body;
}

async function save(ev) {
  ev.preventDefault();
  const form = ev.target;
  const n = draft || topo.find((x) => x.name === selected);
  if (!n || !shown) return;
  // the required fields left empty are marked, every one, and the first
  // has the focus: nothing is sent
  const missing = missingOf(n, form);
  if (missing.length) {
    shown.all = true;
    refreshShown();
    missing[0].box.focus();
    status(_("Fill in what is marked: {fields}.", { fields: [...new Set(missing.map((m) => legendOf(m.box)))].join(", ") }), true);
    return;
  }
  // a field that is wrong (a name taken, a port another node's) is said
  const wrong = [...form.elements].find((e) => e.validity && e.validity.customError);
  if (wrong) {
    status(wrong.validationMessage, true);
    wrong.focus();
    return;
  }
  const name = form.querySelector("#f-name").value;
  const bad = nameError(name, draft ? "" : n.name);
  if (bad) {
    status(bad, true);
    form.querySelector("#f-name").focus();
    return;
  }
  const rename = !draft && name !== n.name;
  if (rename && !confirm(_("Rename {was} to {name}?", { was: n.name, name }) + "\n\n" + renameNote(n, name))) return;
  let body;
  try {
    body = formOf(n, form, false);
  } catch (e) {
    status(e.message, true);
    return;
  }
  let renamed = false;
  try {
    const was = new Set(brokenEdges(topo));
    if (rename) {
      // the node, its edges and what is kept under its name move first;
      // its settings are saved under the new name after
      await api("POST", `/api/nodes/${encodeURIComponent(n.name)}/rename`, { name });
      renamed = true;
      edits.delete(n.name);
      selected = name;
      shown.name = name;
      canvas.setSelected(name);
      status(_("{was} renamed {name}", { was: n.name, name }));
    }
    const res = await api("PUT", "/api/nodes/" + encodeURIComponent(name), body);
    // what was typed is saved: the form is made of the node again
    shown = null;
    edits.delete(name);
    if (draft) {
      const d = draft;
      drafts = drafts.filter((x) => x !== d);
      draft = null;
      selected = name;
      syncCanvas();
      await api("PUT", "/api/layout", { positions: { [name]: d.pos } });
    }
    status(_("{name} saved", { name }));
    if (res.key) showKey(name, res.kind, res.key);
    await load();
    // a change of settings may leave an edge carrying what the node
    // after it does not take: said here, and drawn dashed
    const broke = brokenEdges(topo).filter((e) => !was.has(e));
    if (broke.length) status(_("Saved, but {edges} no longer connects: the node after does not take what reaches it (the dashed edge). Put a convert or a template node between, or remove the edge.", { edges: broke.join("; ") }), true);
    $("sidebar").scrollTop = 0;
  } catch (e) {
    status(renamed ? _("Renamed {name}, but its settings were not saved: {why}", { name, why: e.message }) : e.message, true);
    if (renamed) await load();
  }
}

async function remove(name) {
  if (!confirm(_("Delete {name}? A node another one sends to stays. A template node's templates go with it; a display's slips and devices stay until it is added again.", { name }))) return;
  try {
    await api("DELETE", "/api/nodes/" + encodeURIComponent(name));
    edits.delete(name);
    if (shown && shown.name === name) shown = null;
    selected = null;
    status(_("{name} deleted", { name }));
    await load();
  } catch (e) {
    status(e.message, true);
  }
}

// usagePanel is how a saved node is used: what is at it now (first: a
// sink's, a display's, an inbox's, an archive's; a downloader's after
// its readers), and how computers, applications and devices reach it from
// outside (a CUPS server, an uploader, an S3 trigger, an inbox, a
// downloader, a display).
function usagePanel(n, box) {
  const now = atSection(n);
  if (now && n.kind !== "download") box.append(now);
  switch (n.kind) {
    case "cups":
      box.append(printerSetup(n));
      break;
    case "upload":
    case "s3-trigger":
      entryUsage(n, box);
      break;
    case "inbox": {
      const sec = section(_("Work it"));
      sec.append(hint(_("People decide on the Inboxes page; applications through the API: what the inbox holds (the items, oldest first, and the nodes it releases to), an item's thumbnail, pages and document, and the decision. An app the inbox lists sends its token (POST /oauth/token, the scope inbox); a frontend the token it got for a person who works the inbox; the pages use their session.")));
      const base = `${location.origin}/api/inbox/${n.name}`;
      sec.append(pre(`GET  ${base}
GET  ${base}/<job>
GET  ${base}/<job>/thumbnail?size=256
GET  ${base}/<job>/pages/<n>.png?w=800
GET  ${base}/<job>/document
POST ${base}/<job>/release   {"to": [${releaseTo(n).map((x) => JSON.stringify(x)).join(", ")}], "worker": "…"}
POST ${base}/<job>/discard   {"worker": "…"}`));
      const allowed = n.apps || [];
      sec.append(hint(allowed.length ? _("Apps that work it with their own token: {apps}.", { apps: allowed.join(", ") }) : _("No app works it yet (Settings → Apps allowed).")));
      box.append(sec);
      break;
    }
    case "download": {
      const sec = section(_("Readers"));
      sec.append(hint(_("Applications and people fetch what is listed here: the entries newest first (?owner= one owner's, ?limit=), one entry, its record (the archive's, or of an inbox's the job), and the document itself, read from the inbox or the archive before it. A backend sends its token (POST /oauth/token, client credentials, the scope download); a frontend sends the token it got for the person (a public client, /oauth/authorize with PKCE); the pages use their session.")));
      const base = `${location.origin}/api/download/${n.name}`;
      sec.append(pre(`GET ${base}
GET ${base}/<job>
GET ${base}/<job>/record
GET ${base}/<job>/document
GET ${base}/<job>/pages/<n>.png?w=`));
      box.append(sec, now);
      break;
    }
    case "printer":
    case "cups-out": {
      const sec = section(_("Where it prints"));
      sec.append(hint(n.kind === "cups-out"
        ? sentences(_("This system prints each job to this queue of the CUPS server over IPP, as a computer does: PDF (or PWG Raster), as the queue says it supports; the server's filters and drivers do the rest."), n.credential && _("It gives the user and password of the credential {name}.", { name: n.credential }))
        : _("This system prints each job on this printer itself: over IPP (PDF or PWG Raster, as the printer says it supports), or raw over 9100 or LPD.")));
      sec.append(pre(n.url || "—"));
      if (n.bonjour) sec.append(hint(_("Picked with Bonjour as {name}: when this address stops answering, it is looked for again by the name on the local network, and printed to where it is now.", { name: n.bonjour })));
      box.append(sec);
      break;
    }
    case "display": {
      const dev = section(_("Devices"));
      box.append(dev);
      devicesSection(n, dev).catch((e) => status(e.message, true));
      break;
    }
    default:
      box.append(hint(_("Other nodes send to it along the edges: it is not reached from outside. Computers print to a CUPS server (source); applications post to an uploader, and a bucket's events to an S3 trigger.")));
  }
  // the first section has no rule above it, under the tabs
  box.firstElementChild?.classList.remove("mt-3", "border-t", "pt-3");
  const row = el("div", undefined, { className: "mt-3 flex flex-wrap gap-3 border-t border-base-300 pt-3 text-sm" });
  if (n.kind !== "download") row.append(el("a", _("Jobs sent here"), { className: "link", href: "/admin/jobs?dest=" + encodeURIComponent(n.name) }));
  if (n.kind === "inbox") row.append(el("a", _("{n} held: open the inbox", { n: n.held || 0 }), { className: "link", href: "/inboxes#" + encodeURIComponent(n.name) }));
  if (n.pending) row.append(el("a", n.pending === 1 ? _("{n} delivery waiting", { n: n.pending }) : _("{n} deliveries waiting", { n: n.pending }), { className: "link", href: "/admin/deliveries" }));
  // what a node publishes at the API, for an application's developer
  if (n.api) row.append(el("a", "OpenAPI", { className: "link", href: "/api/openapi.json?node=" + encodeURIComponent(n.name), download: n.name + ".openapi.json" }));
  box.append(row);
  if (n.kind === "printer" && n.preview) {
    const sec = section(_("Preview"));
    box.append(sec);
    previewSection(n, sec).catch((e) => status(e.message, true));
  }
}

function pre(text) {
  return el("pre", text, { className: "overflow-x-auto rounded-box bg-base-200 p-2 font-mono text-xs break-all whitespace-pre-wrap" });
}

// printerSetup tells how a computer adds a CUPS server node as a printer:
// at its own port, as a printer of a CUPS server (/printers/<name>), or
// at the manager's own server (/ipp/print/<name>: on a serverless
// deployment, the only one). Looking, the port is not known: the
// manager's own server is said.
function printerSetup(n) {
  const sec = section(_("Add it as a printer"));
  const own = !!n.port;
  const scheme = own ? (n.tls_cert ? "ipps" : "ipp") : location.protocol === "https:" ? "ipps" : "ipp";
  const host = own ? `${location.hostname}:${n.port}` : location.host;
  const queue = own ? `printers/${n.name}` : ((n.addresses || {}).ipp || `/ipp/print/${n.name}`).replace(/^\//, "");
  const url = `${scheme}://${host}/${queue}`;
  sec.append(hint(n.to && n.to.length
    ? _("Computers and business systems add it as an IPP printer, as a printer of a CUPS server: what they print goes on to {nodes}, made .bdf.", { nodes: n.to.join(", ") })
    : _("Computers and business systems add it as an IPP printer, as a printer of a CUPS server: what they print goes on to no node yet (connect it in Settings), made .bdf.")));
  sec.append(pre(url));
  const steps = el("ul", undefined, { className: "flex list-disc flex-col gap-1 pl-5 text-sm" });
  if (own && n.dnssd !== false) {
    const name = n.title || n.name;
    steps.append(
      el("li", _("Macs, iPhones and iPads on the local network find it with Bonjour as {name} @ <this computer>: macOS: System Settings → Printers & Scanners → Add Printer → Default, Use: AirPrint; iPhone and iPad: Print → Printer.", { name })),
      el("li", _("Android phones and tablets on the local network find it as {name} @ <this computer>, a Mopria printer: Print → Select a printer → All printers (the Default Print Service, or the Mopria Print Service).", { name })),
      el("li", _("Windows computers on the local network find it as {name} @ <this computer>, an IPP Everywhere printer: Settings → Bluetooth & devices → Printers & scanners → Add device (the Microsoft IPP Class Driver); when it is not listed, add it by name as below.", { name })),
    );
  }
  steps.append(
    el("li", _("macOS: System Settings → Printers & Scanners → Add Printer → IP: address {host}, protocol IPP (Internet Printing Protocol), queue {queue}, Use: AirPrint. Not Generic PostScript Printer: its PostScript is not taken.", { host, queue })),
    el("li", _("Windows: Settings → Printers & scanners → Add device → Add manually → Select a shared printer by name: {url}.", { url: `${scheme === "ipps" ? "https" : "http"}://${host}/${queue}` })),
    el("li", _("Linux, or a Mac's Terminal (CUPS):")),
  );
  sec.append(steps);
  sec.append(pre(`lpadmin -p ${n.name} -E -v ${url} -m everywhere` + (own ? `\nlpstat -h ${host} -p` : "")));
  const listed = n.apps || [];
  sec.append(hint(sentences(
    _("When the manager has submit tokens (DOCFLOW_SUBMIT_TOKENS), the computer asks for a name and a password: a token's name and its secret."),
    listed.length && _("The apps allowed ({apps}) give their id and secret instead.", { apps: listed.join(", ") }),
    own && _("The manager's own server takes it too, at {url}.", { url: `${location.protocol === "https:" ? "ipps" : "ipp"}://${location.host}${(n.addresses || {}).ipp || ""}` }),
  )));
  return sec;
}

// releaseTo are the nodes an inbox releases to: its to, but the
// downloaders, which list what it holds.
function releaseTo(n) {
  return (n.to || []).filter((x) => (topo.find((t) => t.name === x) || {}).kind !== "download");
}

// entryUsage is how applications send to an uploader, and how a bucket's
// events reach an S3 trigger: the request, and the node's key.
function entryUsage(n, box) {
  const sec = section(n.kind === "upload" ? _("Send to it") : _("Point the bucket's events at it"));
  if (n.kind === "upload") {
    sec.append(hint(n.accept === "data" ? _("Applications POST the data (JSON, YAML, CSV, an Excel workbook or a Parquet file) as the body, with the key as a bearer token; ?name= names the job, X-Meta-* headers and ?meta.<key>= give metadata.") : _("Applications POST a document as the body, with the key as a bearer token; ?name= names the job, X-Meta-* headers and ?meta.<key>= give metadata.")));
    sec.append(pre(`curl -X POST -H "Authorization: Bearer <key>" --data-binary @${n.accept === "data" ? "orders.json" : "document.pdf"} "${location.origin}${n.api}?name=document"`));
    const allowed = n.apps || [];
    sec.append(hint(allowed.length ? _("Apps that send with their own token (POST /oauth/token): {apps}.", { apps: allowed.join(", ") }) : _("No app sends with its own token yet: only the key (Settings → Apps allowed).")));
  } else {
    sec.append(hint(_("It does not listen to the bucket: an S3 event goes to an EventBridge rule (an API destination) or a Lambda, which posts the object's reference here with the node's key. The body is {body}, an S3 notification (Records, one object), or the event EventBridge sends (detail). The node reads the object itself; the job goes on with s3.bucket and s3.key in its metadata.", { body: '{"bucket": …, "key": …}' })));
    sec.append(pre(`curl -X POST -H "Authorization: Bearer <key>" -H "Content-Type: application/json" \\
  --data '{"bucket": "${bucketOf(n.url) || "<bucket>"}", "key": "${n.path ? n.path + "/" : ""}<key>"}' "${location.origin}${n.api}"`));
  }
  const key = section(_("Key"));
  key.append(hint(_("Shown once, when the node was made or the key rotated; the manager keeps only its hash. Rotating it makes a new one and the old one stops working.")));
  key.append(button(_("Rotate key"), async () => showKey(n.name, n.kind, (await api("POST", `/api/nodes/${n.name}/key`)).key), "btn btn-sm btn-outline w-fit"));
  box.append(sec, key);
}

// bucketOf is the bucket of an s3:// URL.
function bucketOf(u) {
  try {
    return new URL(u).hostname;
  } catch {
    return "";
  }
}

// previewSection shows what the printer's fake printer received: the
// latest jobs, and the pages read back of the newest.
async function previewSection(n, box) {
  box.append(hint(_("What the fake printer on the farm received, as it read it back.")));
  const jobs = await api("GET", `/api/nodes/${n.name}/preview`);
  if (!jobs.length) {
    box.append(el("p", _("Nothing received yet."), { className: "text-xs text-base-content/60" }));
    return;
  }
  const latest = jobs[0];
  box.append(el("p", `#${latest.id} ${latest.name || ""} · ${latest.state || ""} · ${since(latest.created)}`, { className: "text-sm font-semibold" }));
  const pages = el("div", undefined, { className: "flex flex-col gap-2" });
  for (let k = 1; k <= (latest.pages || 0); k++) {
    pages.append(el("img", undefined, { src: `/api/nodes/${n.name}/preview/${latest.id}/pages/${k}.png`, alt: _("page {n}", { n: k }), className: "w-full rounded bg-white shadow-sm", loading: "lazy" }));
  }
  if (!latest.pages) pages.append(el("p", _("The farm could not read it back as pages."), { className: "text-xs text-base-content/60" }));
  box.append(pages);
  if (jobs.length > 1) {
    const list = el("ul", undefined, { className: "text-xs text-base-content/60" });
    for (const j of jobs.slice(1)) list.append(el("li", `#${j.id} ${j.name || ""} · ${since(j.created)}`));
    box.append(el("p", _("Before it:"), { className: "text-xs font-semibold" }), list);
  }
}

async function templatesSection(n, box) {
  box.append(hint(_("This node fills one of them (the one chosen above, or the one the metadata names) with the data that reaches it from a data uploader. A compiled template (.bdft), or Excel, Word, PowerPoint, a PDF form, pdfme or text.")));
  const ts = await api("GET", `/api/nodes/${n.name}/templates`);
  const table = el("table", undefined, { className: "table table-xs" });
  const tbody = el("tbody");
  for (const t of ts) {
    const tr = el("tr");
    const fields = (t.fields || []).map((f) => f.name).join(", ");
    tr.append(el("td", t.id, { className: "font-mono font-semibold" }), el("td", t.format), el("td", (t.fields || []).length === 1 ? _("{n} field", { n: 1 }) : _("{n} fields", { n: (t.fields || []).length }), { title: fields }), el("td", `${Math.max(1, Math.round(t.size / 1024))} KiB`), el("td"));
    tr.lastChild.append(el("a", _("File"), { className: "btn btn-xs btn-ghost", href: `/api/nodes/${n.name}/templates/${t.id}/file` }), button(_("Delete"), () => api("DELETE", `/api/nodes/${n.name}/templates/${t.id}`), "btn btn-xs btn-ghost text-error"));
    tbody.append(tr);
  }
  if (!ts.length) {
    const tr = el("tr");
    tr.append(el("td", _("No templates yet"), { colSpan: 5, className: "text-base-content/50" }));
    tbody.append(tr);
  }
  table.append(tbody);
  const wrap = el("div", undefined, { className: "overflow-x-auto" });
  wrap.append(table);
  box.append(wrap);
  const form = el("form", undefined, { className: "flex flex-wrap gap-1" });
  const id = el("input", undefined, { className: "input input-sm w-28", placeholder: _("id (invoice)"), required: true, pattern: "[A-Za-z0-9][A-Za-z0-9_\\-]{0,63}" });
  const file = el("input", undefined, { type: "file", className: "file-input file-input-sm min-w-0 flex-1", required: true });
  form.append(id, file, el("button", _("Upload"), { type: "submit", className: "btn btn-sm" }));
  form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    const f = file.files[0];
    if (!f) return;
    try {
      await api("PUT", `/api/nodes/${n.name}/templates/${encodeURIComponent(id.value)}?name=${encodeURIComponent(f.name)}`, f, true);
      status(_("{id} uploaded", { id: id.value }));
      changed();
    } catch (e) {
      status(e.message, true);
    }
  });
  box.append(form);
}

async function devicesSection(n, box) {
  box.append(hint(sentences(n.open === 1 ? _("{n} open slip on the board.", { n: 1 }) : _("{n} open slips on the board.", { n: n.open || 0 }), _("Only the devices enrolled here see its documents; a device keeps a key of its own."))));
  const ds = await api("GET", `/api/nodes/${n.name}/devices`);
  const table = el("table", undefined, { className: "table table-xs" });
  const tbody = el("tbody");
  for (const d of ds) {
    const tr = el("tr");
    tr.append(el("td", d.name, { className: "font-semibold" }), el("td", _("seen {when}", { when: since(d.last_seen) }), { className: "text-base-content/60" }), el("td"));
    tr.lastChild.append(button(_("Remove"), () => api("DELETE", "/api/devices/" + d.id), "btn btn-xs btn-ghost text-error"));
    tbody.append(tr);
  }
  if (!ds.length) {
    const tr = el("tr");
    tr.append(el("td", _("No devices yet"), { colSpan: 3, className: "text-base-content/50" }));
    tbody.append(tr);
  }
  table.append(tbody);
  box.append(table);
  const add = el("button", _("Add a device (QR code)"), { type: "button", className: "btn btn-sm btn-primary w-fit" });
  add.addEventListener("click", async () => {
    try {
      const e = await api("POST", `/api/nodes/${n.name}/enrollments`);
      $("qr").innerHTML = e.qr; // SVG the manager drew
      $("enroll-code").textContent = e.code.replace(/(.{5})/, "$1-");
      $("enroll-url").textContent = e.url;
      $("enroll-url").href = e.url;
      // a device makes its key with WebCrypto: HTTPS, or localhost
      const u = new URL(e.url);
      $("enroll-insecure").hidden = u.protocol === "https:" || ["localhost", "127.0.0.1", "[::1]"].includes(u.hostname);
      $("enroll-expires").textContent = _("Valid once, until {when}", { when: since(e.expires) });
      $("enroll-modal").showModal();
    } catch (err) {
      status(err.message, true);
    }
  });
  box.append(add);
}

function showKey(name, kind, key) {
  $("key-name").textContent = name;
  $("key-value").textContent = key;
  $("key-usage").textContent = `curl -X POST -H "Authorization: Bearer ${key}" --data-binary @document.pdf "${location.origin}/api/${kind}/${name}?name=document"`;
  $("key-copy").textContent = _("Copy");
  $("key-modal").showModal();
}

// The Copy button puts the key on the clipboard. Without the clipboard
// API (a page that is not a secure context) the key is selected, so
// that the keyboard's copy takes it.
if (editable) $("key-copy").addEventListener("click", async () => {
  const key = $("key-value").textContent;
  let copied = false;
  try {
    await navigator.clipboard.writeText(key);
    copied = true;
  } catch {
    const whole = document.createRange();
    whole.selectNodeContents($("key-value"));
    const sel = getSelection();
    sel.removeAllRanges();
    sel.addRange(whole);
    try {
      copied = document.execCommand("copy");
    } catch {}
  }
  $("key-copy").textContent = copied ? _("Copied") : _("Selected: press Ctrl+C or ⌘C");
  setTimeout(() => ($("key-copy").textContent = _("Copy")), 2500);
});


if (editable) document.addEventListener("keydown", (ev) => {
  const typing = ["INPUT", "TEXTAREA", "SELECT"].includes(document.activeElement?.tagName);
  if (ev.key === "Escape" && draft) {
    dropDraft(draft);
  } else if (ev.key === "Delete" && selected && !typing) {
    remove(selected);
  }
});

// Leaving the page with drafts, or changes not saved, asks first.
if (editable) window.addEventListener("beforeunload", (ev) => {
  keepShown();
  if (!drafts.length && !edits.size) return;
  ev.preventDefault();
  ev.returnValue = "";
});

// a link from another page names a node: /topology#<name>
if (location.hash.length > 1) {
  selected = decodeURIComponent(location.hash.slice(1));
  canvas.setSelected(selected);
}
onChange(load);
load().then(poll);
setInterval(() => {
  if (!document.hidden) poll();
}, POLL_MS);
if (!editable) {
  // looking, the nodes are read again now and then (editing, Refresh)
  setInterval(() => {
    if (!document.hidden) load();
  }, 60000);
}
document.addEventListener("visibilitychange", () => {
  if (!document.hidden) poll();
});

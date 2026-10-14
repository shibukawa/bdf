// The form page: a template of a form (a workbook, a document or a deck
// with tags, a PDF with form fields, a template of pdfme, or a compiled
// template) is read inside the browser by the form
// module (cmd/bdfwasm built with -tags formonly: the Go packages form and
// vector/pdfbdf, as a server uses them). The page makes a form of the JSON
// Schema of the data the tags ask for, fills the template with what is typed
// into it, draws the filled form, and saves it as a PDF file.
import type { Manifest, Page } from "@bdfkit/core";
import { BdfWorkerClient } from "@bdfkit/render";
import { ConverterClient } from "../../examples/common/convert.js";
import { failure, kb, message } from "../intake.js";

/** Where the page finds lib/, form/fonts/ and forms/ (set by the build). */
declare const SITE_ROOT: string;

/** A value of JSON. */
type Json = null | boolean | number | string | Json[] | { [name: string]: Json };
type JsonObject = { [name: string]: Json };

/** What a template says of a name (the x-bdf keyword of the schemas of form.Schema). */
interface Info {
  kind?: "value" | "section" | "image";
  format?: string;
  mode?: "rows" | "slots" | "pages";
  slots?: number;
  code?: string;
  inverted?: boolean;
  at?: string[];
  /** What a reader sees of the choices of an enum, in their order. */
  labels?: string[];
  multiline?: boolean;
}

/** The part of JSON Schema that form.Schema writes. */
interface Schema {
  type?: string | string[];
  format?: string;
  title?: string;
  description?: string;
  properties?: Record<string, Schema>;
  items?: Schema;
  maxItems?: number;
  enum?: Json[];
  default?: Json;
  maxLength?: number;
  "x-bdf"?: Info;
}

/** A sample of the page: a shape of form, as templates of several formats filled with the same data. */
interface Preset {
  id: string;
  label: string;
  about: string;
  data: string;
  files: string[];
  templates: { name: string; label: string; style?: "standard" | "modern" }[];
}

/** The template being filled. */
interface Template {
  name: string;
  base: string;
  bytes: ArrayBuffer;
  schema: Schema;
  schemaText: string;
}

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const url = (path: string) => new URL(SITE_ROOT + path, location.href).href;
const MODULE = url("lib/bdf-form.wasm");
const FONTS = url("form/fonts/");

const status = $("status");
const fieldsBox = $<HTMLFormElement>("fields");
const jsonBox = $<HTMLTextAreaElement>("json");
const preview = $<HTMLDivElement>("preview");
const savePdf = $<HTMLButtonElement>("savePdf");
const saveBdf = $<HTMLAnchorElement>("saveBdf");

function setStatus(text: string) {
  status.textContent = text;
  status.toggleAttribute("data-error", text.startsWith("error"));
}

const converter = new ConverterClient(new Worker(url("lib/convert-worker.js")));
const renderer = new BdfWorkerClient(new Worker(url("lib/worker.js"), { type: "module" }));

let template: Template | undefined;
/** The data that fills the template, and the pictures it names by their file names. */
let data: JsonObject = {};
const pictures = new Map<string, Uint8Array>();
/** Bumped by each template: what was started for the one before is dropped. */
let opening = 0;
/** The filled form: the document, and what is offered of it. */
let filled: Uint8Array | undefined;
let bdfURL = "";

// --- the data, and the form made of its schema ---

const types = (s: Schema): string[] => (Array.isArray(s.type) ? s.type : s.type ? [s.type] : []);
const isList = (s: Schema) => types(s).includes("array");
const isObject = (s: Schema) => types(s).includes("object") && !!s.properties;
const isPicture = (s: Schema) => s["x-bdf"]?.kind === "image" && !s["x-bdf"]?.code;

/** Data with nothing in it yet: every name of the schema, a row in each table. */
function blank(s: Schema): Json {
  if (isList(s)) {
    const items = s.items ?? {};
    return isObject(items) ? [blank(items)] : [];
  }
  if (isObject(s)) return Object.fromEntries(Object.entries(s.properties!).map(([name, p]) => [name, blank(p)]));
  // (what the template has for the name: the value of a field of a PDF form)
  if (s.default !== undefined && s.default !== null && typeof s.default !== "object") return s.default;
  const t = types(s);
  if (t.length === 1 && t[0] === "boolean") return false;
  if (t.length === 1 && t[0] === "number") return null;
  return "";
}

/** What was typed, as a value of a text or a number: a number when it is written as one of JSON (007 stays a text). */
const typed = (text: string): Json => (/^-?(0|[1-9]\d*)(\.\d+)?$/.test(text) && text.length < 16 ? Number(text) : text);

const shown = (v: Json | undefined): string => (v === null || v === undefined || typeof v === "object" ? "" : String(v));

/** Where a name is in the template, for the tooltip of its label. */
function where(s: Schema): string {
  const info = s["x-bdf"];
  const at = info?.at?.length ? `in ${info.at.join(", ")}` : "";
  const format = info?.format ? `shown as ${info.format}` : "";
  return [s.description, at, format].filter(Boolean).join("; ");
}

/** What a list is, after its name. */
function listHint(s: Schema): string {
  const info = s["x-bdf"];
  if (info?.mode === "slots") return `at most ${s.maxItems ?? info.slots} to a page`;
  if (info?.mode === "pages") return "a page for each";
  if (info?.mode === "rows") return "a row for each";
  return "";
}

let nextId = 0;
/** What gets the focus after the form is made again: the id of a control. */
let focusNext = "";

/** A control for a value that is not a list nor an object; set is called with what is typed. */
function control(s: Schema, value: Json | undefined, set: (v: Json) => void, inTable: boolean): HTMLElement {
  const t = types(s);
  const only = (name: string) => t.length === 1 && t[0] === name;
  if (only("boolean")) {
    const box = document.createElement("input");
    box.type = "checkbox";
    box.checked = value === true;
    box.onchange = () => set(box.checked);
    return box;
  }
  if (isPicture(s)) return pictureControl(value, set);
  if (s.enum?.length) {
    // one of the choices of a field, or none
    const choice = document.createElement("select");
    const labels = s["x-bdf"]?.labels ?? [];
    const options = s.enum.map((v, i) => new Option(labels[i] ?? shown(v), shown(v)));
    const now = shown(value);
    if (now !== "" && !s.enum.some((v) => shown(v) === now)) options.push(new Option(now, now));
    choice.append(new Option("", ""), ...options);
    choice.value = now;
    choice.onchange = () => set(choice.value);
    return choice;
  }
  const input = document.createElement("input");
  if (only("number")) {
    input.type = "number";
    input.step = "any";
    input.value = typeof value === "number" ? String(value) : "";
    input.oninput = () => set(input.value === "" || Number.isNaN(input.valueAsNumber) ? null : input.valueAsNumber);
    return input;
  }
  const text = shown(value);
  if (s.format === "date" && /^(\d{4}-\d{2}-\d{2})?$/.test(text)) {
    input.type = "date";
  } else if (s.format === "date-time" && /^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(:\d{2})?)?$/.test(text)) {
    input.type = "datetime-local";
  } else if (!inTable) {
    // a text may have lines
    const area = document.createElement("textarea");
    area.rows = Math.min(6, text.split("\n").length);
    area.value = text;
    if (s.maxLength) area.maxLength = s.maxLength;
    area.oninput = () => {
      area.rows = Math.min(6, area.value.split("\n").length);
      set(t.includes("number") ? typed(area.value) : area.value);
    };
    if (s["x-bdf"]?.code) area.inputMode = /^(ean|upc|itf|jan)/.test(s["x-bdf"].code) ? "numeric" : "text";
    return area;
  }
  input.value = text;
  if (s.maxLength) input.maxLength = s.maxLength;
  input.oninput = () => set(t.includes("number") && input.type === "text" ? typed(input.value) : input.value);
  return input;
}

/** A picture of the data: its file name, and a button that picks a file for it. */
function pictureControl(value: Json | undefined, set: (v: Json) => void): HTMLElement {
  const box = document.createElement("span");
  box.className = "picture";
  const name = document.createElement("input");
  name.value = shown(value).startsWith("data:") ? "" : shown(value);
  name.placeholder = "no picture";
  name.oninput = () => set(name.value);
  const pick = document.createElement("input");
  pick.type = "file";
  pick.accept = "image/png,image/jpeg,image/gif";
  pick.hidden = true;
  pick.onchange = async () => {
    const file = pick.files?.[0];
    pick.value = "";
    if (!file) return;
    pictures.set(file.name, new Uint8Array(await file.arrayBuffer()));
    name.value = file.name;
    set(file.name);
  };
  const button = document.createElement("button");
  button.type = "button";
  button.textContent = "Choose…";
  button.onclick = () => pick.click();
  box.append(name, button, pick);
  return box;
}

/** A labelled control, for a name of an object. */
function field(name: string, s: Schema, value: Json | undefined, set: (v: Json) => void): HTMLElement {
  const title = s.title ?? name;
  if (isList(s)) return list(title, s, Array.isArray(value) ? value : [], set);
  if (isObject(s)) {
    const group = document.createElement("fieldset");
    const legend = document.createElement("legend");
    legend.textContent = title;
    group.append(legend, properties(s, value, set));
    return group;
  }
  const label = document.createElement("label");
  const c = control(s, value, set, false);
  const text = document.createElement("span");
  text.textContent = title;
  text.title = where(s);
  const code = s["x-bdf"]?.code;
  if (code) {
    const hint = document.createElement("small");
    hint.textContent = code === "qr" ? "QR code" : `barcode: ${code}`;
    text.append(" ", hint);
  }
  if (c instanceof HTMLInputElement && c.type === "checkbox") {
    label.className = "check";
    label.append(c, text);
  } else {
    label.className = "field";
    label.append(text, c);
  }
  return label;
}

/** The controls of the names of an object; set is called with the object when one of them changes. */
function properties(s: Schema, value: Json | undefined, set: (v: Json) => void): HTMLElement {
  const obj: JsonObject = value !== null && typeof value === "object" && !Array.isArray(value) ? value : {};
  const box = document.createElement("div");
  box.className = "props";
  // the lists after the values: a table of many rows does not part what is around it in the template
  const all = Object.entries(s.properties ?? {});
  for (const [name, p] of [...all.filter(([, p]) => !isList(p)), ...all.filter(([, p]) => isList(p))]) {
    box.append(field(name, p, obj[name], (v) => {
      obj[name] = v;
      set(obj);
    }));
  }
  return box;
}

/** The items of a list: a table of rows when they are values or objects of values, and otherwise a group for each. */
function list(title: string, s: Schema, items: Json[], set: (v: Json) => void): HTMLElement {
  const item = s.items ?? {};
  const group = document.createElement("fieldset");
  group.className = "list";
  const legend = document.createElement("legend");
  legend.textContent = title;
  legend.title = where(s);
  const hint = listHint(s);
  if (hint) {
    const small = document.createElement("small");
    small.textContent = hint;
    legend.append(" ", small);
  }
  group.append(legend);
  const full = s.maxItems !== undefined && items.length >= s.maxItems;
  const remake = (focus: string) => {
    focusNext = focus;
    set(items);
    buildForm();
  };
  const addId = `add${nextId++}`;
  const remove = (i: number, label: string) => {
    const b = document.createElement("button");
    b.type = "button";
    b.className = "remove";
    b.textContent = "×";
    b.title = b.ariaLabel = `remove ${label}`;
    b.onclick = () => {
      items.splice(i, 1);
      remake(addId);
    };
    return b;
  };
  const props = isObject(item) ? Object.entries(item.properties!) : undefined;
  const flat = !props || props.every(([, p]) => !isList(p) && !isObject(p));
  if (flat) {
    const table = document.createElement("table");
    if (props) {
      const head = table.createTHead().insertRow();
      for (const [name, p] of props) {
        const th = document.createElement("th");
        th.scope = "col";
        th.textContent = p.title ?? name;
        th.title = where(p);
        head.append(th);
      }
      head.append(document.createElement("td"));
    }
    const body = table.createTBody();
    items.forEach((it, i) => {
      const row = body.insertRow();
      const cell = (p: Schema, value: Json | undefined, put: (v: Json) => void, label: string) => {
        const c = control(p, value, put, true);
        (c.querySelector("input") ?? c).setAttribute("aria-label", label);
        row.insertCell().append(c);
      };
      if (props) {
        const obj: JsonObject = it !== null && typeof it === "object" && !Array.isArray(it) ? it : (items[i] = {});
        for (const [name, p] of props) {
          cell(p, obj[name], (v) => {
            obj[name] = v;
            set(items);
          }, `${p.title ?? name} of ${title} ${i + 1}`);
        }
      } else {
        cell(item, it, (v) => {
          items[i] = v;
          set(items);
        }, `${title} ${i + 1}`);
      }
      row.insertCell().append(remove(i, `${title} ${i + 1}`));
    });
    const wrap = document.createElement("div");
    wrap.className = "rows";
    wrap.append(table);
    group.append(wrap);
  } else {
    items.forEach((it, i) => {
      const one = document.createElement("fieldset");
      one.className = "item";
      const l = document.createElement("legend");
      l.append(`${title} ${i + 1} `, remove(i, `${title} ${i + 1}`));
      one.append(l, properties(item, it, (v) => {
        items[i] = v;
        set(items);
      }));
      group.append(one);
    });
  }
  const add = document.createElement("button");
  add.type = "button";
  add.id = addId;
  add.className = "add";
  add.textContent = flat ? "Add a row" : `Add an item to ${title}`;
  add.disabled = full;
  if (full) add.title = `${title} has room for ${s.maxItems}`;
  add.onclick = () => {
    items.push(isObject(item) ? blank(item) : "");
    remake(addId);
  };
  group.append(add);
  return group;
}

/** Make the form of the template's schema, with the data in it. */
function buildForm() {
  if (!template) return;
  nextId = 0;
  const names = Object.keys(template.schema.properties ?? {});
  if (!names.length) {
    const p = document.createElement("p");
    p.className = "note";
    p.textContent = "This template has no tags: it asks for no data. Write {{name}} where a value goes.";
    fieldsBox.replaceChildren(p);
    return;
  }
  fieldsBox.replaceChildren(properties(template.schema, data, (v) => {
    data = v as JsonObject;
    changed();
  }));
  if (focusNext) document.getElementById(focusNext)?.focus();
  focusNext = "";
}

/** The data changed in the form (or, from, in its JSON): fill the template again. */
function changed(from?: "json") {
  if (from !== "json") showJson();
  schedule();
}

function showJson() {
  jsonBox.value = JSON.stringify(data, null, 2);
  jsonBox.removeAttribute("aria-invalid");
  const a = $<HTMLAnchorElement>("saveJson");
  if (a.href) URL.revokeObjectURL(a.href);
  a.href = URL.createObjectURL(new Blob([jsonBox.value + "\n"], { type: "application/json" }));
}

jsonBox.oninput = () => {
  let v: unknown;
  try {
    v = JSON.parse(jsonBox.value);
  } catch (e) {
    jsonBox.setAttribute("aria-invalid", "true");
    setStatus(`error: the data is not JSON: ${message(e)}`);
    return;
  }
  if (v === null || typeof v !== "object" || Array.isArray(v)) {
    jsonBox.setAttribute("aria-invalid", "true");
    setStatus("error: the data is an object of JSON");
    return;
  }
  jsonBox.removeAttribute("aria-invalid");
  data = v as JsonObject;
  changed("json");
};

// --- the workspace tabs, and the tabs within Data ---

function tabGroup(ids: string[], onSelect: (tab: HTMLElement) => void) {
  const tabs = ids.map((id) => $(id));
  const select = (tab: HTMLElement) => {
    if (tab.hasAttribute("disabled")) return;
    for (const t of tabs) {
      const on = t === tab;
      t.setAttribute("aria-selected", String(on));
      t.tabIndex = on ? 0 : -1;
      $(t.getAttribute("aria-controls")!).hidden = !on;
    }
    onSelect(tab);
  };
  for (const tab of tabs) {
    tab.onclick = () => select(tab);
    tab.onkeydown = (e) => {
      const enabled = tabs.filter((t) => !t.hasAttribute("disabled"));
      const step = e.key === "ArrowRight" ? 1 : e.key === "ArrowLeft" ? -1 : 0;
      const to = e.key === "Home" ? enabled[0] : e.key === "End" ? enabled.at(-1)
        : step ? enabled[(enabled.indexOf(tab) + step + enabled.length) % enabled.length] : undefined;
      if (!to) return;
      e.preventDefault();
      select(to);
      to.focus();
    };
  }
  return { tabs, select };
}

const dataTabs = tabGroup(["tabForm", "tabJson", "tabSchema"], (tab) => {
  // The JSON may have changed the data.
  if (tab.id === "tabForm") buildForm();
});
const workspaceTabs = tabGroup(["tabTemplate", "tabData", "tabPreview"], (tab) => {
  if (tab.id === "tabPreview") layout();
});

// --- filling the template ---

let timer = 0;
let busy = false;
let again = false;

/** Fill the template soon: after the typing pauses. */
function schedule() {
  clearTimeout(timer);
  timer = window.setTimeout(() => { fillNow().catch((e) => setStatus(failure(e))); }, 250);
}

async function fillNow() {
  if (busy) {
    again = true;
    return;
  }
  busy = true;
  try {
    do {
      again = false;
      await fillOnce();
    } while (again);
  } finally {
    busy = false;
  }
}

async function fillOnce() {
  const t = template;
  if (!t) return;
  const token = opening;
  const t0 = performance.now();
  // the pictures the data names, of those picked or fetched
  const text = JSON.stringify(data);
  const files: Record<string, Uint8Array> = {};
  for (const [name, bytes] of pictures) if (text.includes(JSON.stringify(name))) files[name] = bytes;
  let res;
  try {
    res = await converter.fill(MODULE, t.bytes, { data: text, files, fonts: FONTS, name: t.name });
  } catch (e) {
    if (token === opening) setStatus(failure(e));
    return;
  }
  if (token !== opening) return;
  filled = res.bdf;
  const took = performance.now() - t0;
  if (bdfURL) URL.revokeObjectURL(bdfURL);
  bdfURL = URL.createObjectURL(new Blob([filled as BlobPart], { type: "application/octet-stream" }));
  saveBdf.href = bdfURL;
  saveBdf.download = `${t.base}.bdf`;
  saveBdf.hidden = false;
  savePdf.disabled = false;
  setWarnings(res.warnings);
  await draw(filled.slice().buffer as ArrayBuffer);
  if (token !== opening) return;
  setStatus(`${t.name} filled in ${took.toFixed(0)} ms: ${res.pages} ${res.pages === 1 ? "page" : "pages"}, ${kb(filled.length)} as a bdf document`);
}

function setWarnings(list: string[]) {
  const box = $<HTMLDetailsElement>("warnings");
  box.hidden = list.length === 0;
  box.querySelector("summary")!.textContent = `${list.length} ${list.length === 1 ? "warning" : "warnings"}`;
  box.querySelector("ul")!.replaceChildren(...list.map((w) => {
    const li = document.createElement("li");
    li.textContent = w;
    return li;
  }));
}

// --- the pages of the filled form ---

interface Frame {
  view: string;
  index: number;
  page: Page;
  box: HTMLDivElement;
  canvas: HTMLCanvasElement;
  /** The drawing its canvas shows, and the one it is being drawn for. */
  drawn: number;
  asked: number;
}

let frames: Frame[] = [];
/** Bumped by each document and each size of the pages: a canvas shows the drawing of its number. */
let drawing = 0;
/** CSS px to a unit of the pages. */
let scale = 1;
let manifest: Manifest | undefined;
const near = new Set<Frame>();
const byBox = new Map<Element, Frame>();
const watcher = new IntersectionObserver((entries) => {
  for (const e of entries) {
    const f = byBox.get(e.target);
    if (!f) continue;
    if (e.isIntersecting) {
      near.add(f);
      paint(f);
    } else {
      near.delete(f);
    }
  }
}, { root: preview, rootMargin: "400px" });

/** Show a filled form: its pages take the places of those of the one before, which stay until they are drawn. */
async function draw(buffer: ArrayBuffer) {
  drawing++;
  manifest = undefined;
  await renderer.close();
  manifest = await renderer.open({ kind: "buffer", buffer });
  layout();
}

/** Fit the widest page inside the preview's padding. */
function fittedScale(widest: number) {
  const style = getComputedStyle(preview);
  const width = preview.clientWidth - parseFloat(style.paddingLeft) - parseFloat(style.paddingRight);
  return Math.min(1.5, Math.max(0.2, width / widest));
}

/** Lay the pages out at the width of the pane, and draw the ones in sight. */
function layout() {
  if (!manifest || !preview.clientWidth) return;
  const pages = manifest.views.flatMap((v) => (v.kind === "fixed" || v.kind === "flow" ? (v.pages ?? []).map((page, index) => ({ view: v.id, index, page })) : []));
  const widest = Math.max(1, ...pages.map((p) => p.page.w));
  scale = fittedScale(widest);
  drawing++;
  // the frames of the pages before are kept, with what they show
  while (frames.length > pages.length) {
    const f = frames.pop()!;
    watcher.unobserve(f.box);
    byBox.delete(f.box);
    near.delete(f);
    f.box.remove();
  }
  pages.forEach((p, i) => {
    let f = frames[i];
    if (!f) {
      const box = document.createElement("div");
      box.className = "page";
      const canvas = document.createElement("canvas");
      box.append(canvas);
      preview.append(box);
      f = frames[i] = { ...p, box, canvas, drawn: 0, asked: 0 };
      byBox.set(box, f);
      watcher.observe(box);
    }
    Object.assign(f, p);
    f.box.style.width = `${p.page.w * scale}px`;
    f.box.style.height = `${p.page.h * scale}px`;
    f.canvas.setAttribute("aria-label", `page ${i + 1} of ${pages.length}`);
    f.canvas.setAttribute("role", "img");
  });
  for (const f of near) paint(f);
}

async function paint(f: Frame) {
  if (!manifest || !preview.clientWidth) return;
  const turn = drawing;
  if (f.drawn === turn || f.asked === turn) return;
  f.asked = turn;
  try {
    const bitmap = await renderer.page(f.view, f.index, scale * (window.devicePixelRatio || 1));
    if (turn !== drawing || !f.box.isConnected) {
      bitmap.close();
      return;
    }
    f.canvas.width = bitmap.width;
    f.canvas.height = bitmap.height;
    f.canvas.getContext("2d")!.drawImage(bitmap, 0, 0);
    bitmap.close();
    f.drawn = turn;
  } catch (e) {
    // (a document that came meanwhile took the place of the one asked of)
    if (turn === drawing) console.error("bdf:", e);
  }
}

let resizing = 0;
new ResizeObserver(() => {
  clearTimeout(resizing);
  resizing = window.setTimeout(() => {
    if (manifest && preview.clientWidth && Math.abs(fittedScale(Math.max(1, ...frames.map((f) => f.page.w))) - scale) > 0.001) layout();
  }, 150);
}).observe(preview);

// --- saving ---

savePdf.onclick = async () => {
  const t = template, doc = filled;
  if (!t || !doc) return;
  savePdf.disabled = true;
  const label = savePdf.textContent;
  savePdf.textContent = "Writing the PDF…";
  try {
    const t0 = performance.now();
    // a plain file, a tagged one, or a PDF/A file of a level
    const kind = $<HTMLSelectElement>("pdfKind").value;
    const res = await converter.pdf(MODULE, doc.slice().buffer as ArrayBuffer, { fonts: FONTS, tagged: kind === "tagged", pdfa: /^\d/.test(kind) ? kind : undefined });
    const file = new Blob([res.pdf as BlobPart], { type: "application/pdf" });
    const a = document.createElement("a");
    a.href = URL.createObjectURL(file);
    a.download = `${t.base}.pdf`;
    a.click();
    setTimeout(() => URL.revokeObjectURL(a.href), 60_000);
    const notes = res.warnings.length ? `; ${res.warnings.length} ${res.warnings.length === 1 ? "warning" : "warnings"}: ${res.warnings.join("; ")}` : "";
    setStatus(`${t.base}.pdf written in ${(performance.now() - t0).toFixed(0)} ms: ${res.pages} ${res.pages === 1 ? "page" : "pages"}, ${res.fonts} embedded ${res.fonts === 1 ? "font" : "fonts"}, ${kb(file.size)}${notes}`);
  } catch (e) {
    setStatus(failure(e));
  } finally {
    savePdf.textContent = label;
    savePdf.disabled = false;
  }
};

// --- templates ---

/** Open a template: read its tags, make the form, and fill it with the sample (or with nothing). */
async function openTemplate(name: string, bytes: ArrayBuffer, sample?: JsonObject) {
  const token = ++opening;
  const base = name.replace(/\.[^.]*$/, "");
  setStatus(`reading ${name}…`);
  let res;
  try {
    res = await converter.fields(MODULE, bytes, { sample: sample ? JSON.stringify(sample) : undefined });
  } catch (e) {
    if (token === opening) setStatus(failure(e));
    return;
  }
  if (token !== opening) return;
  const schema = JSON.parse(res.schema) as Schema;
  template = { name, base, bytes, schema, schemaText: res.schema };
  // a compiled template brings its sample
  const given = sample ?? (res.sample ? (JSON.parse(res.sample) as JsonObject) : undefined);
  data = given ? structuredClone(given) : (blank(schema) as JsonObject);
  filled = undefined;
  savePdf.disabled = true;
  saveBdf.hidden = true;
  setWarnings([]);
  $("schema").textContent = res.schema;
  const a = $<HTMLAnchorElement>("saveSchema");
  if (a.href) URL.revokeObjectURL(a.href);
  a.href = URL.createObjectURL(new Blob([res.schema], { type: "application/schema+json" }));
  a.download = `${base}.schema.json`;
  $<HTMLAnchorElement>("saveJson").download = `${base}.json`;
  for (const tab of workspaceTabs.tabs) tab.removeAttribute("disabled");
  dataTabs.select($("tabForm"));
  workspaceTabs.select($("tabData"));
  $("tabData").focus({ preventScroll: true });
  showJson();
  how(template);
  document.title = `${name} – Forms – BDF`;
  setStatus(`filling ${name}…`);
  await fillNow();
}

/** How a server does the same: with the bdf command, and with the Go packages. */
function how(t: Template) {
  $("cli").textContent = `bdf template schema -json-schema ${t.name} ${t.base}.schema.json   # what the form above is made of
bdf fill -data ${t.base}.json ${t.name} ${t.base}.bdf
bdf pdf ${t.base}.bdf ${t.base}.pdf            # -tagged for a tagged PDF, -pdfa 2a for a PDF/A file`;
  $("go").textContent = `schema, err := form.SchemaOf(template, size, nil) // the JSON Schema of the data
// …
res, err := form.Fill(template, size, &form.Options{Data: data})
// …
_, err = pdfbdf.Write(w, res.Doc, &pdfbdf.Options{}) // Tagged: true, PDFA: "2a"`;
  $("how").hidden = false;
}

async function fetchBytes(path: string): Promise<ArrayBuffer> {
  const r = await fetch(url(path));
  if (!r.ok) throw new Error(`${path}: HTTP ${r.status}`);
  return r.arrayBuffer();
}

/** Open a sample: a template with the data and the pictures of its shape. */
async function openPreset(p: Preset, name: string) {
  setStatus(`fetching ${name}…`);
  try {
    const [bytes, sample] = await Promise.all([fetchBytes(`forms/${name}`), fetchBytes(`forms/${p.data}`).then((b) => JSON.parse(new TextDecoder().decode(b)) as JsonObject)]);
    for (const f of p.files) if (!pictures.has(f)) pictures.set(f, new Uint8Array(await fetchBytes(`forms/${f}`)));
    await openTemplate(name, bytes, sample);
  } catch (e) {
    setStatus(failure(e));
  }
}

function renderPresets(presets: Preset[], style: string) {
  const body = $("presets").querySelector("tbody")!;
  body.replaceChildren();
  for (const p of presets) {
    const templates = p.templates.filter((t) => !t.style || t.style === style);
    if (!templates.length) continue;
    const row = body.insertRow();
    const th = document.createElement("th");
    th.scope = "row";
    th.textContent = p.label;
    const about = row.insertCell();
    about.textContent = p.about;
    const cell = row.insertCell();
    cell.className = "samples";
    for (const t of templates) {
      const entry = document.createElement("span");
      entry.className = "sample-template";
      const b = document.createElement("button");
      b.type = "button";
      b.textContent = t.label;
      b.title = t.name;
      b.ariaLabel = `${p.label}, ${t.style ?? ""} ${t.label}`;
      b.onclick = () => openPreset(p, t.name);
      const source = document.createElement("a");
      source.href = url(`forms/${t.name}`);
      source.download = t.name;
      source.textContent = `Download .${t.name.split(".").pop()}`;
      source.title = `Original template: ${t.name}`;
      source.ariaLabel = `${p.label}, ${t.style ?? ""} ${t.label}: download original template`;
      entry.append(b, source);
      cell.append(entry);
    }
    row.prepend(th);
  }
  $("sampleStyleNote").textContent = style === "modern"
    ? "More open spacing, navy accents and light horizontal rules."
    : "Black and white, with a centered title and a full table grid.";
}

async function listPresets() {
  const res = await fetch(url("forms/index.json"));
  if (!res.ok) return;
  const presets = (await res.json()) as Preset[];
  // "?template=form-rows.xlsx" opens that sample
  const want = new URLSearchParams(location.search).get("template");
  const preset = presets.find((p) => p.templates.some((t) => t.name === want));
  const picker = $<HTMLSelectElement>("sampleStyle");
  picker.value = preset?.templates.find((t) => t.name === want)?.style ?? "standard";
  picker.onchange = () => renderPresets(presets, picker.value);
  renderPresets(presets, picker.value);
  $("sampleStyleControls").hidden = !presets.some((p) => p.templates.some((t) => t.style));
  $("presets").hidden = false;
  if (preset && want) await openPreset(preset, want);
}

function openLocal(file: File) {
  file.arrayBuffer().then((bytes) => openTemplate(file.name, bytes)).catch((e) => setStatus(failure(e)));
}

const picker = $<HTMLInputElement>("file");
picker.onchange = () => {
  const file = picker.files?.[0];
  picker.value = "";
  if (file) openLocal(file);
};
$("choose").onclick = () => picker.click();

// a template dropped anywhere on the page
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

fieldsBox.onsubmit = (e) => e.preventDefault();
listPresets().catch((e) => console.error("bdf:", e));

// The search text page: a file opened or dropped on it is converted into
// bdf inside the browser, and its text for a search index is shown in a
// text area: what Document.SearchText gives a server, from the preview
// module (cmd/bdfwasm built with -tags previewonly).
import { Intake, converted, kb, message, type Source } from "../intake.js";

/** Where the page finds lib/, fonts/ and samples/ (set by the build). */
declare const SITE_ROOT: string;

/** The text of a document, as bdf text writes it. */
interface SearchText {
  meta: { source?: string; dc?: Record<string, unknown> };
  views: { id: string; kind: string; title?: string; pages: { page?: number; text: string }[] }[];
}

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const status = $("status");
const output = $<HTMLTextAreaElement>("output");
const shape = $<HTMLSelectElement>("shape");
const copy = $<HTMLButtonElement>("copy"), save = $<HTMLAnchorElement>("save");

function setStatus(text: string) {
  status.textContent = text;
  status.toggleAttribute("data-error", text.startsWith("error"));
}

/** The document shown, and its text as JSON. */
let shown: { source: Source; json: string; text: SearchText } | undefined;
let url = "";

const intake = new Intake({
  root: SITE_ROOT,
  status: setStatus,
  open: extract,
  picker: $<HTMLInputElement>("file"),
  buttons: [$("choose")],
  samples: $("samples"),
});

shape.onchange = show;
copy.onclick = () => {
  navigator.clipboard.writeText(output.value).then(() => setStatus("copied"), (e) => setStatus(`error: ${message(e)}`));
};

async function extract(s: Source) {
  const token = intake.token;
  shown = undefined;
  show();
  const t0 = performance.now();
  // the worker takes the buffer: each call gets a copy
  const one = (password?: string) => intake.converter.text(intake.preview, s.bdf.slice().buffer as ArrayBuffer, { password });
  const res = await intake.withPassword(one, "extracting…", "This document is encrypted. Enter its password to read its text.", s.password);
  if (!intake.current(token)) return;
  if (!res) {
    setStatus("encrypted document: no text extracted");
    return;
  }
  s.password = res.password;
  const text = JSON.parse(res.value.json) as SearchText;
  shown = { source: s, json: res.value.json, text };
  const pages = text.views.reduce((n, v) => n + v.pages.length, 0);
  const chars = text.views.reduce((n, v) => v.pages.reduce((n, p) => n + p.text.length, n), 0);
  setStatus(`${converted(s)}; the text of ${pages} ${pages === 1 ? "page" : "pages"}, ${chars.toLocaleString("en")} characters, extracted in ${(performance.now() - t0).toFixed(0)} ms`);
  show();
}

/** The text alone: the pages one after another, under the title of their view when there are several. */
function plain(text: SearchText): string {
  const parts: string[] = [];
  for (const v of text.views) {
    for (const p of v.pages) {
      const head = [text.views.length > 1 ? v.title || v.id : "", p.page ? `page ${p.page}` : ""].filter(Boolean).join(", ");
      parts.push(`${head ? `— ${head} —\n` : ""}${p.text}`);
    }
  }
  return parts.join("\n\n");
}

function show() {
  if (url) URL.revokeObjectURL(url);
  url = "";
  const json = shape.value === "json";
  output.value = !shown ? "" : json ? shown.json : plain(shown.text);
  copy.disabled = !shown;
  save.hidden = !shown;
  $("notes").hidden = $("how").hidden = !shown;
  if (!shown) return;
  const s = shown.source;
  const file = new Blob([output.value], { type: json ? "application/json" : "text/plain" });
  save.href = url = URL.createObjectURL(file);
  save.download = `${s.base}-text.${json ? "json" : "txt"}`;
  save.textContent = `Save (${kb(file.size)})`;
  const notes = ["Pages without text are left out; a sheet is one entry, without a page number."];
  if (s.protected) notes.push("The file was password-protected: its text is not.");
  if (s.password !== undefined) notes.push("The document is encrypted: its text is not.");
  $("notes").textContent = notes.join(" ");
  $("cli").textContent = s.format === "bdf" ? `bdf text ${s.name} text.json` : `bdf generate -text text.json ${s.name} out.bdf`;
}

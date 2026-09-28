// The thumbnail page: a file opened or dropped on it is converted into bdf
// inside the browser (only its first page, all a thumbnail shows), and its
// thumbnails are drawn at once at several sizes
// by the preview module (cmd/bdfwasm built with -tags previewonly: the Go
// packages thumbnail and raster, as a server uses them).
import type { Thumbnail, ThumbnailOptions } from "../../examples/common/convert.js";
import { Intake, converted, kb, message, type Source } from "../intake.js";

/** Where the page finds lib/, fonts/ and samples/ (set by the build). */
declare const SITE_ROOT: string;

/** The sizes drawn: the side of a cropped thumbnail, the longer side of a fitted one. */
const SIZES = [64, 128, 256, 512];

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const status = $("status");
const list = $<HTMLUListElement>("thumbs");
const modeBox = $<HTMLSelectElement>("mode"), formatBox = $<HTMLSelectElement>("format"), dpiBox = $<HTMLInputElement>("dpi");

function setStatus(text: string) {
  status.textContent = text;
  status.toggleAttribute("data-error", text.startsWith("error"));
}

/** The document the thumbnails are of, and the password that opens it. */
let source: Source | undefined;
/** Object URLs of the thumbnails shown: revoked when they are replaced. */
let urls: string[] = [];
/** Bumped by each drawing: one started before is dropped. */
let drawing = 0;

const intake = new Intake({
  root: SITE_ROOT,
  status: setStatus,
  open: (s) => {
    source = s;
    return draw();
  },
  picker: $<HTMLInputElement>("file"),
  buttons: [$("choose")],
  samples: $("samples"),
  // the pages after the first are not in the thumbnails: a long PDF converts in milliseconds instead of seconds
  pages: "1",
});

for (const box of [modeBox, formatBox, dpiBox]) box.onchange = () => { draw().catch((e) => setStatus(`error: ${message(e)}`)); };

function options(): Pick<ThumbnailOptions, "mode" | "format" | "sheetDpi"> {
  const dpi = Math.min(288, Math.max(24, Number(dpiBox.value) || 72));
  return { mode: modeBox.value as ThumbnailOptions["mode"], format: formatBox.value as ThumbnailOptions["format"], sheetDpi: dpi };
}

/** Draw the thumbnails of the document at every size, with the options chosen. */
async function draw() {
  const s = source;
  if (!s) return;
  const turn = ++drawing, token = intake.token;
  const live = () => turn === drawing && intake.current(token);
  const opts = options();
  for (const u of urls) URL.revokeObjectURL(u);
  urls = [];
  list.replaceChildren();
  $("notes").hidden = $("how").hidden = true;
  const warnings = new Set<string>();
  const t0 = performance.now();
  let bytes = 0;
  for (const size of SIZES) {
    // the worker takes the buffer: each call gets a copy
    const one = (password?: string) => intake.converter.thumbnail(intake.preview, s.bdf.slice().buffer as ArrayBuffer, { ...opts, size, password, fonts: intake.fonts });
    const res = await intake.withPassword(one, "drawing…", "This document is encrypted. Enter its password to draw its thumbnails.", s.password);
    if (!live()) return;
    if (!res) {
      setStatus("encrypted document: no thumbnails drawn");
      return;
    }
    s.password = res.password;
    const t = res.value;
    bytes += t.image.length;
    for (const w of t.warnings) warnings.add(w);
    list.append(figure(s, t, size));
  }
  setStatus(`${converted(s)}; ${SIZES.length} thumbnails drawn in ${(performance.now() - t0).toFixed(0)} ms (${kb(bytes)} in all)`);
  const notes = [...warnings].map((w) => `Warning: ${w}`);
  if (s.protected) notes.unshift("The file was password-protected: its thumbnails are not.");
  if (s.password !== undefined) notes.unshift("The document is encrypted: its thumbnails are not.");
  $("notes").textContent = notes.join(" ");
  $("notes").hidden = !notes.length;
  how(s, opts);
}

function figure(s: Source, t: Thumbnail, size: number): HTMLLIElement {
  const ext = t.format === "jpeg" ? "jpg" : "png";
  const file = new Blob([t.image as BlobPart], { type: `image/${t.format}` });
  const url = URL.createObjectURL(file);
  urls.push(url);
  const li = document.createElement("li");
  const fig = document.createElement("figure");
  const img = document.createElement("img");
  img.src = url;
  img.width = t.width;
  img.height = t.height;
  img.alt = `thumbnail of ${s.name}, ${t.width} by ${t.height} pixels`;
  const cap = document.createElement("figcaption");
  const title = document.createElement("strong");
  title.textContent = `${size} px`;
  const a = document.createElement("a");
  a.href = url;
  a.download = `${s.base}-${size}.${ext}`;
  a.textContent = "Save";
  cap.append(title, `${t.width} × ${t.height} ${t.format.toUpperCase()}, ${t.mode}, ${kb(file.size)} · `, a);
  fig.append(img, cap);
  li.append(fig);
  return li;
}

/** How a server makes the same thumbnails: with the bdf command, and with the Go packages. */
function how(s: Source, o: ReturnType<typeof options>) {
  const ext = o.format === "jpeg" ? "jpg" : "png";
  const mode = o.mode === "auto" ? "" : ` -mode ${o.mode}`;
  const dpi = o.sheetDpi === 72 ? "" : ` -sheet-dpi ${o.sheetDpi}`;
  const input = s.format === "bdf" ? undefined : s.name;
  // of a file to convert, bdf thumbnail converts the first page only, as this page does
  $("cli").textContent = `bdf thumbnail -size 256${mode}${dpi} ${s.name} thumb.${ext}`;
  const goMode = { auto: "thumbnail.Auto", crop: "thumbnail.Crop", fit: "thumbnail.Fit" }[o.mode ?? "auto"];
  // a document converted for its thumbnails only: its first page (a server that keeps the document converts it all)
  $("go").textContent = `res, err := converter.ConvertFile(${JSON.stringify(input ?? "in.pptx")}, "", &converter.Options{Pages: converter.PageList(1)})
// …
for _, size := range []int{${SIZES.join(", ")}} {
	th, err := thumbnail.Make(res.Doc, &thumbnail.Options{Size: size, Mode: ${goMode}, SheetDPI: ${o.sheetDpi}})
	// …
	err = thumbnail.Encode(w, th.Image, thumbnail.${o.format === "jpeg" ? "JPEG" : "PNG"}) // thumbnail.WebP on a server
}`;
  $("how").hidden = false;
}

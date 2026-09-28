// What the thumbnail and the search text pages share: a file picked, dropped
// on the page or chosen among the samples becomes a bdf document, converted
// inside the browser by the converter modules (cmd/bdfwasm) as the viewer
// converts it, and the page makes of it what a server would, with the
// preview module.
import { ConverterClient, ConvertError, sniff } from "../examples/common/convert.js";

/** A document to make previews of. */
export interface Source {
  /** The name of the file it came from, and that name without its extension. */
  name: string;
  base: string;
  /** The document, as a single-file bdf. */
  bdf: Uint8Array;
  /** The format of the file ("pptx", "pdf", …; "bdf" for a document opened as it is). */
  format: string;
  /** What the converter said of it, and how long it took in milliseconds (0 for a document opened as it is). */
  summary: string;
  ms: number;
  warnings: string[];
  /** The password of the document, when it is encrypted (a bdf document opened as it is). */
  password?: string;
  /** The file needed a password: what is made of it is not protected as it was. */
  protected: boolean;
}

export interface IntakeOptions {
  /** Where lib/ (the workers and the converter modules), fonts/ and samples/ are, from the page. */
  root: string;
  /** What the page is doing, for its status line. */
  status: (text: string) => void;
  /** A document is ready. */
  open: (source: Source) => void | Promise<void>;
  /** The file picker, the buttons that open it, and where the samples go (hidden until they are listed). */
  picker: HTMLInputElement;
  buttons: HTMLElement[];
  samples: HTMLElement;
  /** The pages (slides, sheets) of a file to convert, as bdf generate -pages takes them; all when absent. */
  pages?: string;
}

const MODULES = { pdf: "lib/bdf-pdf.wasm", office: "lib/bdf-office.wasm", web: "lib/bdf-web.wasm", image: "lib/bdf-image.wasm" };

/** Why a call of the converter worker failed, when it says so with a code. */
export const errorCode = (e: unknown) => (e instanceof ConvertError ? e.code : undefined);
export const message = (e: unknown) => (e as Error).message ?? String(e);

export class Intake {
  readonly converter: ConverterClient;
  /** The preview module, and the fonts that text in fonts referred to by name is drawn with. */
  readonly preview: string;
  readonly fonts: string;
  /** Bumped by each file: one opened while another is converted wins. */
  private opening = 0;
  private dialog?: HTMLDialogElement;

  constructor(private options: IntakeOptions) {
    this.converter = new ConverterClient(new Worker(this.url("lib/convert-worker.js")));
    this.preview = this.url("lib/bdf-preview.wasm");
    this.fonts = this.url("fonts/");
    this.install();
  }

  /** The URL of what the pages of the site share. */
  url(path: string): string {
    return new URL(this.options.root + path, location.href).href;
  }

  /** Whether source is still the document of the page (no other file was opened since). */
  current(token: number): boolean {
    return token === this.opening;
  }
  get token(): number {
    return this.opening;
  }

  private install() {
    const { picker, buttons } = this.options;
    picker.onchange = () => {
      const file = picker.files?.[0];
      picker.value = "";
      if (file) this.openLocal(file);
    };
    for (const b of buttons) b.onclick = () => picker.click();
    // files dropped anywhere on the page
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
      if (file) this.openLocal(file);
    });
    this.listSamples().catch(() => {});
    // "?file=samples/basic.docx" opens that file
    const file = new URLSearchParams(location.search).get("file");
    if (file) this.openURL(/^samples\//.test(file) ? this.url(file) : new URL(file, location.href).href);
  }

  private async listSamples() {
    const res = await fetch(this.url("samples/index.json"));
    if (!res.ok) return;
    const samples = (await res.json()) as { name: string; label: string }[];
    const box = this.options.samples;
    box.replaceChildren(...samples.map((s) => {
      const b = document.createElement("button");
      b.type = "button";
      b.textContent = s.label;
      b.title = s.name;
      b.onclick = () => this.openURL(this.url(`samples/${s.name}`), s.name);
      return b;
    }));
    box.hidden = false;
  }

  private fail = (e: unknown) => this.options.status(`error: ${message(e)}`);

  private openLocal(file: File) {
    file.arrayBuffer().then((data) => this.openFile(file.name, data)).catch(this.fail);
  }

  /** Fetch a file by its URL and open it. */
  openURL(url: string, name = decodeURIComponent(new URL(url).pathname.split("/").pop() ?? "")) {
    this.options.status(`fetching ${name}…`);
    return fetch(url)
      .then((r) => (r.ok ? r.arrayBuffer() : Promise.reject(new Error(`${name}: HTTP ${r.status}`))))
      .then((data) => this.openFile(name, data))
      .catch(this.fail);
  }

  /** Open a file: a bdf document as it is, anything else converted into one (asking for the password of a protected file). */
  async openFile(name: string, data: ArrayBuffer) {
    const token = ++this.opening;
    const { status, open } = this.options;
    const base = name.replace(/\.[^.]*$/, "");
    const kind = sniff(new Uint8Array(data), name);
    if (kind === "bdf") {
      await open({ name, base, bdf: new Uint8Array(data), format: "bdf", summary: "", ms: 0, warnings: [], protected: false });
      return;
    }
    const busy = `converting ${name}…`;
    status(busy);
    let t0 = 0; // of the last attempt: the reader's typing is not part of the conversion
    const convert = (password?: string) => {
      t0 = performance.now();
      // only the Office converters lay text out with the font directory (PDFs embed their fonts, images have no text)
      return this.converter.convert(this.url(MODULES[kind]), data, { fonts: kind === "office" ? this.fonts : undefined, password, name, pages: this.options.pages });
    };
    let res;
    try {
      res = await this.withPassword(convert, busy, "This file is protected. Enter its password to convert it.");
    } catch (e) {
      if (errorCode(e) !== "unknown-format") throw e;
      throw new Error(`${name} is not in a format this page converts`);
    }
    if (token !== this.opening) return;
    if (!res) {
      status("protected file: not converted");
      return;
    }
    const { value } = res;
    const ms = performance.now() - t0;
    status(`${name} converted in ${ms.toFixed(0)} ms: ${value.summary}`);
    await open({ name, base, bdf: value.bdf, format: value.format, summary: value.summary, ms, warnings: value.warnings, protected: value.protected });
  }

  /**
   * Run with the password known, if any; when a password is needed, ask
   * the reader for one and run again with it, as many times as it takes.
   * The value comes with the password that worked; undefined when the
   * reader cancels.
   */
  async withPassword<T>(run: (password?: string) => Promise<T>, busy: string, prompt: string, known?: string): Promise<{ value: T; password?: string } | undefined> {
    try {
      return { value: await run(known), password: known };
    } catch (e) {
      const code = errorCode(e);
      if (code !== "password-required" && code !== "wrong-password") throw e;
    }
    let text = prompt, error = false;
    for (;;) {
      this.options.status("waiting for the password");
      const password = await this.askPassword(text, error);
      if (password === undefined) return undefined;
      this.options.status(busy);
      try {
        return { value: await run(password), password };
      } catch (e) {
        if (errorCode(e) !== "wrong-password") throw e;
        text = "Wrong password. Try again.";
        error = true;
      }
    }
  }

  /** Ask for a password; undefined when the reader cancels. */
  private askPassword(text: string, error: boolean): Promise<string | undefined> {
    const dialog = (this.dialog ??= passwordDialog());
    const input = dialog.querySelector("input")!;
    const msg = dialog.querySelector("p")!;
    msg.textContent = text;
    msg.toggleAttribute("data-error", error);
    input.value = "";
    dialog.returnValue = "";
    dialog.querySelector<HTMLButtonElement>("button[type=button]")!.onclick = () => dialog.close("cancel");
    dialog.showModal();
    return new Promise((resolve) => {
      dialog.onclose = () => resolve(dialog.returnValue === "ok" ? input.value : undefined);
    });
  }
}

function passwordDialog(): HTMLDialogElement {
  const dialog = document.createElement("dialog");
  dialog.className = "password";
  dialog.setAttribute("aria-labelledby", "pwTitle");
  dialog.setAttribute("aria-describedby", "pwMessage");
  dialog.innerHTML = `<form method="dialog">
  <h2 id="pwTitle">Password</h2>
  <p id="pwMessage"></p>
  <label>Password <input type="password" autocomplete="off" required></label>
  <div class="buttons"><button value="ok">Open</button><button type="button">Cancel</button></div>
</form>`;
  document.body.append(dialog);
  return dialog;
}

/** What the status line says of the conversion a document came from. */
export const converted = (s: Source) => (s.format === "bdf" ? s.name : `${s.name} converted in ${s.ms.toFixed(0)} ms (${s.summary})`);

/** A size in bytes, for a reader. */
export const kb = (n: number) => (n < 1024 ? `${n} B` : n < 1024 * 1024 ? `${(n / 1024).toFixed(1)} KB` : `${(n / 1024 / 1024).toFixed(1)} MB`);

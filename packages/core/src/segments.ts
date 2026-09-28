import { BdfFormatError } from "./bytes.js";
import { BufferSource } from "./container.js";
import { BdfDocument } from "./document.js";
import type { Segment } from "./types.js";

// A document read from a server a few pages at a time (spec §3.6; the Go
// side is package segment). Each request carries the public key of a key
// pair made for it, whose private key cannot be exported; the server seals
// the segment for that key, and once the segment is open nothing refers to
// the private key any more. What was sent then stays sealed, whatever keys
// are taken later: the server's, the TLS keys, the reader's session.
//
// The request is a POST of JSON: { key, view, page, have }, key the public
// key as an uncompressed P-256 point in base64, have the segments held
// already (their parts are not sent again). The answer is the segment, a
// single-file bdf.

export interface SegmentLoaderOptions {
  /** Passed to fetch with each request (credentials, headers). */
  init?: RequestInit;
  /** The view and page (0-based) whose segment is fetched first: the first page of the first view by default. */
  view?: string;
  page?: number;
  /** A page asked for also fetches the segment this many pages on, unless it has come (default 3; 0: none). */
  ahead?: number;
}

/** A segment request the server refused: status is its HTTP status. */
export class BdfSegmentError extends Error {
  constructor(readonly status: number) {
    super(`bdf: the segment request was answered with ${status}`);
    this.name = "BdfSegmentError";
  }
}

const toBase64 = (b: Uint8Array) => btoa(String.fromCharCode(...b));

export class SegmentLoader {
  /** The segments that came, in the order they came. */
  readonly held: Segment[] = [];
  /** Called as each segment comes, once its pages are in the document. */
  onSegment?: (s: Segment) => void;
  /** The document: the first segment, which the others are put in (BdfDocument.addSegment). */
  doc!: BdfDocument;
  /** Requests on their way, with the pages they are expected to bring. */
  private inflight: { view: string; from: number; to: number; done: Promise<void> }[] = [];
  /** The pages of a segment, once a whole one has come (the last segment of a view is shorter). */
  private size = 0;
  private abort = new AbortController();

  private constructor(readonly url: string, private readonly options: SegmentLoaderOptions) {}

  /** Fetch the first segment; its document is doc. */
  static async open(url: string, options: SegmentLoaderOptions = {}): Promise<SegmentLoader> {
    const l = new SegmentLoader(url, options);
    l.doc = await l.fetch(options.view, options.page ?? 0);
    l.took(l.doc.manifest.segment!);
    return l;
  }

  /** Whether page of view has come. */
  has(view: string, page: number): boolean {
    return this.held.some((s) => s.view === view && s.from <= page && page < s.to);
  }

  /**
   * Fetch the segment that holds page of view unless it came, and start
   * fetching the one a few pages on.
   */
  async ensure(view: string, page: number): Promise<void> {
    await this.load(view, page);
    const ahead = this.options.ahead ?? 3;
    if (ahead > 0) this.prefetch(view, page + ahead);
  }

  /** Start fetching the segment that holds page of view, unless it came; a failure is left for ensure to meet. */
  prefetch(view: string, page: number): void {
    if (page >= (this.doc.view(view).pages?.length ?? 0) || this.has(view, page)) return;
    this.load(view, page).catch(() => {});
  }

  /** Stop the requests on their way. */
  close(): void {
    this.abort.abort();
  }

  private async load(view: string, page: number): Promise<void> {
    const n = this.doc.view(view).pages?.length ?? 0;
    if (!Number.isInteger(page) || page < 0 || page >= n) throw new Error(`bdf: no page ${page} in view ${view}`);
    for (;;) {
      if (this.has(view, page)) return;
      // a request that should bring the page: wait for it rather than ask again
      const on = this.inflight.find((f) => f.view === view && f.from <= page && page < f.to);
      if (on) {
        await on.done;
        continue;
      }
      const from = this.size ? page - (page % this.size) : page;
      const to = this.size ? Math.min(from + this.size, n) : page + 1;
      const f = { view, from, to, done: this.fetch(view, page).then((d) => this.merge(d)) };
      this.inflight.push(f);
      try {
        await f.done;
      } finally {
        this.inflight.splice(this.inflight.indexOf(f), 1);
      }
      if (!this.has(view, page)) throw new BdfFormatError(`the segment of page ${page} does not hold it`);
    }
  }

  private async fetch(view: string | undefined, page: number): Promise<BdfDocument> {
    // a key pair for this request only; its private key cannot be exported
    const keyPair = (await crypto.subtle.generateKey({ name: "ECDH", namedCurve: "P-256" }, false, ["deriveKey"])) as CryptoKeyPair;
    const key = new Uint8Array(await crypto.subtle.exportKey("raw", keyPair.publicKey));
    const init = this.options.init ?? {};
    const res = await fetch(this.url, {
      ...init,
      method: "POST",
      headers: { ...(init.headers as Record<string, string> | undefined), "Content-Type": "application/json" },
      body: JSON.stringify({ key: toBase64(key), view, page, have: this.held }),
      cache: "no-store",
      signal: this.abort.signal,
    });
    if (!res.ok) throw new BdfSegmentError(res.status);
    const doc = await BdfDocument.open(new BufferSource(new Uint8Array(await res.arrayBuffer())), { keyPair });
    const s = doc.manifest.segment;
    if (!s || (view !== undefined && s.view !== view) || !(s.from <= page && page < s.to)) throw new BdfFormatError("the server sent another segment");
    return doc;
  }

  private merge(d: BdfDocument): void {
    this.doc.addSegment(d);
    this.took(d.manifest.segment!);
  }

  private took(s: Segment): void {
    if (!this.has(s.view, s.from) || !this.has(s.view, s.to - 1)) this.held.push({ view: s.view, from: s.from, to: s.to });
    if (s.to < (this.doc.view(s.view).pages?.length ?? 0)) this.size = Math.max(this.size, s.to - s.from);
    this.onSegment?.(s);
  }
}

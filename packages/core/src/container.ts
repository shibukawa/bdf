import { ByteReader, BdfFormatError } from "./bytes.js";
import { FORMAT_VERSION } from "./opcodes.js";
import type { Manifest, PartEntry, Encoding } from "./types.js";

export const HEADER_SIZE = 32;
/** Magic number of the single-file form: "bdf" and a NUL byte. */
export const MAGIC = new Uint8Array([0x62, 0x64, 0x66, 0x00]);
/**
 * The manifest JSON inflates to this much at most: a manifest states no size
 * of its own, so a small file could otherwise ask for any amount of memory.
 */
export const MAX_MANIFEST_SIZE = 256 << 20;

export interface Header {
  version: number;
  flags: number;
  manifestOff: number;
  manifestLen: number;
  manifestEnc: Encoding;
}

export function parseHeader(bytes: Uint8Array): Header {
  const r = new ByteReader(bytes);
  const m = r.bytesN(4);
  if (!MAGIC.every((b, i) => m[i] === b)) throw new BdfFormatError("bad magic");
  const version = r.u16();
  if (version > FORMAT_VERSION) throw new BdfFormatError(`unsupported format version ${version}`);
  const flags = r.u16();
  const manifestOff = r.u64();
  const manifestLen = r.u64();
  const manifestEnc: Encoding = r.u8() === 1 ? "deflate-raw" : "identity";
  return { version, flags, manifestOff, manifestLen, manifestEnc };
}

/**
 * Decode a stored part according to its encoding using DecompressionStream.
 * The data inflates to limit bytes at most (for a part, the size its
 * manifest states); data that goes on past the limit is an error, so a small
 * file cannot ask for more memory than its manifest says it needs (deflate
 * packs 1032 bytes into one at best).
 */
export async function decode(bytes: Uint8Array, enc: Encoding | undefined, limit: number): Promise<Uint8Array> {
  if (!enc || enc === "identity") return bytes;
  if (enc !== "deflate-raw") throw new BdfFormatError(`unknown encoding ${enc}`);
  if (!Number.isInteger(limit) || limit < 0) throw new BdfFormatError(`bad part size ${limit}`);
  const ds = new DecompressionStream("deflate-raw");
  const reader = new Blob([bytes as BlobPart]).stream().pipeThrough(ds).getReader();
  const chunks: Uint8Array[] = [];
  let n = 0;
  for (;;) {
    const { done, value } = await reader.read();
    if (done) break;
    n += value.length;
    if (n > limit) {
      reader.cancel().catch(() => {}); // the rest is not inflated
      throw new BdfFormatError("data inflates to more than its stated size");
    }
    chunks.push(value);
  }
  const out = new Uint8Array(n);
  let at = 0;
  for (const c of chunks) {
    out.set(c, at);
    at += c.length;
  }
  return out;
}

/** Whether n bytes at off lie within size bytes; the values come from the file. */
function inRange(off: number, n: number, size: number): boolean {
  return Number.isInteger(off) && Number.isInteger(n) && off >= 0 && n >= 0 && off <= size && n <= size - off;
}

/**
 * Part names are 32 lowercase hex characters (spec §3.1). They go into the
 * URLs of a split document, so any other name is refused: a manifest could
 * otherwise name any path of the server its document is on.
 */
export function checkHash(h: unknown, what = "part name"): asserts h is string {
  if (typeof h !== "string" || !/^[0-9a-f]{32}$/.test(h)) throw new BdfFormatError(`bad ${what} ${JSON.stringify(h)}`);
}

/** Source of manifest and stored (still encoded) part bytes. */
export interface PartSource {
  manifest(): Promise<Manifest>;
  /** Stored bytes of the part (compressed if enc says so). */
  stored(entry: PartEntry): Promise<Uint8Array>;
}

const utf8 = new TextDecoder();

/** Single-file form held fully in memory. */
export class BufferSource implements PartSource {
  private header: Header;
  private manifestPromise?: Promise<Manifest>;
  constructor(private readonly bytes: Uint8Array) {
    this.header = parseHeader(bytes);
  }
  manifest(): Promise<Manifest> {
    this.manifestPromise ??= (async () => {
      const h = this.header;
      if (!inRange(h.manifestOff, h.manifestLen, this.bytes.length)) throw new BdfFormatError("manifest out of range");
      const raw = this.bytes.subarray(h.manifestOff, h.manifestOff + h.manifestLen);
      return JSON.parse(utf8.decode(await decode(raw, h.manifestEnc, MAX_MANIFEST_SIZE))) as Manifest;
    })();
    return this.manifestPromise;
  }
  async stored(e: PartEntry): Promise<Uint8Array> {
    const base = this.header.manifestOff + this.header.manifestLen;
    if (!inRange(e.off ?? 0, e.len, this.bytes.length - base)) throw new BdfFormatError("part out of range");
    const off = base + (e.off ?? 0);
    return this.bytes.subarray(off, off + e.len);
  }
}

/** Single-file form fetched with HTTP Range requests. */
export class RangeSource implements PartSource {
  private header?: Header;
  private manifestPromise?: Promise<Manifest>;
  /** The whole file, once a server has answered a range request with it. */
  private whole?: Promise<Uint8Array>;
  constructor(private readonly url: string, private readonly init: RequestInit = {}) {}

  private async range(off: number, len: number): Promise<Uint8Array> {
    if (!this.whole) {
      const res = await fetch(this.url, { ...this.init, headers: { ...(this.init.headers as Record<string, string>), Range: `bytes=${off}-${off + len - 1}` } });
      if (res.status === 206) return new Uint8Array(await res.arrayBuffer());
      if (res.status !== 200) throw new Error(`bdf: range request failed with ${res.status}`);
      // The server ignored the range and sent the file: it is kept, and the
      // ranges are taken from it (every part would fetch the file again).
      if (this.whole) {
        res.body?.cancel().catch(() => {}); // another request brought it meanwhile
      } else {
        this.whole = res.arrayBuffer().then((b) => new Uint8Array(b));
        this.whole.catch(() => { this.whole = undefined; });
      }
    }
    const all = await this.whole;
    if (!inRange(off, len, all.length)) throw new BdfFormatError("part out of range");
    return all.subarray(off, off + len);
  }

  manifest(): Promise<Manifest> {
    this.manifestPromise ??= (async () => {
      this.header = parseHeader(await this.range(0, HEADER_SIZE));
      const h = this.header;
      const raw = await this.range(h.manifestOff, h.manifestLen);
      return JSON.parse(utf8.decode(await decode(raw, h.manifestEnc, MAX_MANIFEST_SIZE))) as Manifest;
    })();
    return this.manifestPromise;
  }

  async stored(e: PartEntry): Promise<Uint8Array> {
    await this.manifest();
    const h = this.header!;
    if (!inRange(e.off ?? 0, e.len, Number.MAX_SAFE_INTEGER)) throw new BdfFormatError("part out of range");
    return this.range(h.manifestOff + h.manifestLen + (e.off ?? 0), e.len);
  }
}

/** Split form: manifest.json plus parts/<hash> under a base URL. */
export class SplitSource implements PartSource {
  private manifestPromise?: Promise<Manifest>;
  constructor(private readonly base: string, private readonly init: RequestInit = {}) {
    if (!this.base.endsWith("/")) this.base += "/";
  }
  manifest(): Promise<Manifest> {
    this.manifestPromise ??= fetch(this.base + "manifest.json", this.init).then((res) => {
      if (!res.ok) throw new Error(`bdf: manifest fetch failed with ${res.status}`);
      return res.json() as Promise<Manifest>;
    });
    return this.manifestPromise;
  }
  async stored(e: PartEntry): Promise<Uint8Array> {
    checkHash(e.h); // it goes into the URL
    const res = await fetch(this.base + "parts/" + e.h, this.init);
    if (!res.ok) throw new Error(`bdf: part ${e.h} fetch failed with ${res.status}`);
    return new Uint8Array(await res.arrayBuffer());
  }
}

/** Open a single-file document from a URL, streaming the whole file. */
export async function fetchSingle(url: string, init?: RequestInit): Promise<BufferSource> {
  const res = await fetch(url, init);
  if (!res.ok) throw new Error(`bdf: fetch failed with ${res.status}`);
  return new BufferSource(new Uint8Array(await res.arrayBuffer()));
}

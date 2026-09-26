import { ByteReader, BdfFormatError } from "./bytes.js";
import { FORMAT_VERSION } from "./opcodes.js";
import type { Manifest, PartEntry, Encoding } from "./types.js";

export const HEADER_SIZE = 32;
/** Magic number of the single-file form: "bdf" and a NUL byte. */
export const MAGIC = new Uint8Array([0x62, 0x64, 0x66, 0x00]);

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

/** Decode a stored part according to its encoding using DecompressionStream. */
export async function decode(bytes: Uint8Array, enc: Encoding | undefined): Promise<Uint8Array> {
  if (!enc || enc === "identity") return bytes;
  if (enc !== "deflate-raw") throw new BdfFormatError(`unknown encoding ${enc}`);
  const ds = new DecompressionStream("deflate-raw");
  const stream = new Blob([bytes as BlobPart]).stream().pipeThrough(ds);
  return new Uint8Array(await new Response(stream).arrayBuffer());
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
      const raw = this.bytes.subarray(h.manifestOff, h.manifestOff + h.manifestLen);
      return JSON.parse(utf8.decode(await decode(raw, h.manifestEnc))) as Manifest;
    })();
    return this.manifestPromise;
  }
  async stored(e: PartEntry): Promise<Uint8Array> {
    const base = this.header.manifestOff + this.header.manifestLen;
    const off = base + (e.off ?? 0);
    if (off + e.len > this.bytes.length) throw new BdfFormatError("part out of range");
    return this.bytes.subarray(off, off + e.len);
  }
}

/** Single-file form fetched with HTTP Range requests. */
export class RangeSource implements PartSource {
  private header?: Header;
  private manifestPromise?: Promise<Manifest>;
  constructor(private readonly url: string, private readonly init: RequestInit = {}) {}

  private async range(off: number, len: number): Promise<Uint8Array> {
    const res = await fetch(this.url, { ...this.init, headers: { ...(this.init.headers as Record<string, string>), Range: `bytes=${off}-${off + len - 1}` } });
    if (res.status === 206) return new Uint8Array(await res.arrayBuffer());
    if (res.status === 200) {
      // Server ignored the range: take the slice ourselves.
      const all = new Uint8Array(await res.arrayBuffer());
      return all.subarray(off, off + len);
    }
    throw new Error(`bdf: range request failed with ${res.status}`);
  }

  manifest(): Promise<Manifest> {
    this.manifestPromise ??= (async () => {
      this.header = parseHeader(await this.range(0, HEADER_SIZE));
      const h = this.header;
      const raw = await this.range(h.manifestOff, h.manifestLen);
      return JSON.parse(utf8.decode(await decode(raw, h.manifestEnc))) as Manifest;
    })();
    return this.manifestPromise;
  }

  async stored(e: PartEntry): Promise<Uint8Array> {
    await this.manifest();
    const h = this.header!;
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

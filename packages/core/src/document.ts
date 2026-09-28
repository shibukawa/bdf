import { BdfFormatError } from "./bytes.js";
import { decode, type PartSource } from "./container.js";
import { checkManifest, tileSize } from "./manifest.js";
import { BdfKeyError, BdfPasswordError, SealedSource } from "./crypto.js";
import { decodeCues, type Cues } from "./cues.js";
import { decodeObject, decodePathCollection, objectDeps, UseLimits } from "./object.js";
import { decodeTextIndex, type IndexRun } from "./search.js";
import { extractText } from "./text.js";
import type { Manifest, PartEntry, ObjectPart, PathData, Hash, View } from "./types.js";

export interface OpenOptions {
  /** Password of an encrypted document (spec §3.5). */
  password?: string;
  /**
   * The key pair an encrypted document was sealed for (an ecdh key slot,
   * spec §3.5): a segment a server sealed for one request (SegmentLoader).
   */
  keyPair?: CryptoKeyPair;
}

/** A loaded document: manifest plus a cache of decoded parts. */
export class BdfDocument {
  /** Part entries, with the source that holds each (addPage brings parts of other sources). */
  private entries = new Map<Hash, { entry: PartEntry; source: PartSource }>();
  private parts = new Map<Hash, Promise<Uint8Array>>();
  private objects = new Map<Hash, Promise<ObjectPart>>();
  private pathSets = new Map<Hash, Promise<PathData[]>>();

  private constructor(readonly source: PartSource, readonly manifest: Manifest) {
    checkManifest(manifest);
    for (const e of manifest.parts) this.entries.set(e.h, { entry: e, source });
  }

  /**
   * Open a document. An encrypted one needs options.password: without it
   * BdfPasswordError("required") is thrown, and BdfPasswordError("wrong")
   * when it does not open the document. The source's manifest is cached, so
   * the same source can be opened again with another password. One sealed
   * for a key pair needs options.keyPair, and throws BdfKeyError when it is
   * another. A manifest with part names or numbers out of range (see
   * checkManifest) is refused with BdfFormatError.
   */
  static async open(source: PartSource, options: OpenOptions = {}): Promise<BdfDocument> {
    const manifest = await source.manifest();
    const enc = manifest.encryption;
    if (!enc) return new BdfDocument(source, manifest);
    const secret: string | CryptoKeyPair | undefined = options.keyPair ?? options.password;
    if (secret === undefined) {
      // a password would open nothing: do not ask for one
      if (Array.isArray(enc.keys) && !enc.keys.some((k) => k.type === "password")) throw new BdfKeyError();
      throw new BdfPasswordError("required");
    }
    const sealed = await SealedSource.unlock(source, secret);
    return new BdfDocument(sealed, await sealed.manifest());
  }

  entry(hash: Hash): PartEntry {
    const e = this.entries.get(hash);
    if (!e) throw new Error(`bdf: unknown part ${hash}`);
    return e.entry;
  }

  /**
   * Put the page of a page document (the only page of its only view: a page
   * a streamed conversion returned) in place of page index of a view, and
   * add the parts it brings. The page's objects may use parts added before.
   */
  addPage(viewId: string, index: number, from: BdfDocument): void {
    const pages = this.view(viewId).pages;
    if (!pages || index < 0 || index >= pages.length) throw new Error(`bdf: no page ${index} in view ${viewId}`);
    const views = from.manifest.views;
    if (views.length !== 1 || views[0].pages?.length !== 1) throw new Error("bdf: not a page document");
    this.adopt(from);
    pages[index] = views[0].pages[0];
  }

  /**
   * Put in the pages a segment document carries (spec §3.6): a document with
   * the same views and pages, whose manifest names the pages it carries.
   * Its parts are added to those this document has; its pages may use parts
   * that earlier segments brought.
   */
  addSegment(from: BdfDocument): void {
    const s = from.manifest.segment;
    if (!s) throw new Error("bdf: not a segment document");
    const mine = this.manifest.views, theirs = from.manifest.views;
    if (mine.length !== theirs.length || mine.some((v, i) => v.id !== theirs[i].id || (v.pages?.length ?? 0) !== (theirs[i].pages?.length ?? 0))) {
      throw new Error("bdf: the segment is of another document");
    }
    const src = from.view(s.view).pages ?? [], dst = this.view(s.view).pages ?? [];
    if (!Number.isInteger(s.from) || !Number.isInteger(s.to) || s.from < 0 || s.from >= s.to || s.to > dst.length) {
      throw new BdfFormatError(`segment of pages ${s.from} to ${s.to} out of range`);
    }
    this.adopt(from);
    for (let i = s.from; i < s.to; i++) dst[i] = src[i];
  }

  /** Add the parts of another document that this one does not have, read from its source. */
  private adopt(from: BdfDocument): void {
    for (const e of from.manifest.parts) {
      if (this.entries.has(e.h)) continue;
      this.entries.set(e.h, { entry: e, source: from.source });
      this.manifest.parts.push(e);
    }
  }

  view(id: string): View {
    const v = this.manifest.views.find((v) => v.id === id);
    if (!v) throw new Error(`bdf: unknown view ${id}`);
    return v;
  }

  /** Decoded bytes of a part. */
  part(hash: Hash): Promise<Uint8Array> {
    let p = this.parts.get(hash);
    if (!p) {
      const e = this.entries.get(hash);
      if (!e) throw new Error(`bdf: unknown part ${hash}`);
      // a part inflates to the size its manifest states at most
      p = e.source.stored(e.entry).then((b) => decode(b, e.entry.enc, e.entry.size ?? 0));
      this.parts.set(hash, p);
    }
    return p;
  }

  object(hash: Hash): Promise<ObjectPart> {
    let p = this.objects.get(hash);
    if (!p) {
      p = this.part(hash).then(decodeObject);
      this.objects.set(hash, p);
    }
    return p;
  }

  /** Synchronous access to an object that has already been loaded. */
  objectSync(hash: Hash): ObjectPart {
    const o = this.loadedObjects.get(hash);
    if (!o) throw new Error(`bdf: object ${hash} not loaded`);
    return o;
  }
  private loadedObjects = new Map<Hash, ObjectPart>();

  /**
   * Text index runs of a view: the text index part when present, otherwise
   * built by extracting text from every object of the view (which loads them all).
   */
  async textIndex(view: View): Promise<IndexRun[]> {
    if (view.textIndex) return decodeTextIndex(await this.part(view.textIndex));
    const runs: IndexRun[] = [];
    // A sheet tile keeps the runs whose anchor lies in it (tiles repeat what straddles them).
    const tile = tileSize(view);
    const keep = (r: { x: number; y: number }) => view.kind !== "sheet" || (r.x >= 0 && r.x < tile && r.y >= 0 && r.y < tile);
    const limits = new UseLimits(); // of the whole view
    const add = async (a: number, b: number, hash: Hash) => {
      const obj = await this.ensure(hash);
      let n = 0;
      for (const r of extractText(obj, (h) => this.objectSync(h), undefined, limits)) {
        if (!keep(r)) continue;
        runs.push({ a, b, ordinal: r.ordinal, sep: n++ === 0 ? 2 : r.sep, text: r.text });
      }
    };
    if (view.kind === "sheet") {
      const keys = Object.keys(view.tiles ?? {}).map((k) => k.split(",").map(Number) as [number, number]);
      keys.sort((p, q) => p[1] - q[1] || p[0] - q[0]);
      for (const [x, y] of keys) await add(x, y, view.tiles![`${x},${y}`]);
    } else {
      const pages = view.pages ?? [];
      for (let pi = 0; pi < pages.length; pi++) {
        for (let li = 0; li < pages[pi].layers.length; li++) await add(pi, li, pages[pi].layers[li].obj);
      }
    }
    return runs;
  }

  /**
   * The music of a view (docs/spec.md §4.4): its Standard MIDI File as it is
   * stored, and its cues when it has them; null for a view without play.
   */
  async play(view: View): Promise<{ seq: Uint8Array; cues: Cues | null } | null> {
    if (!view.play) return null;
    const { seq, cues } = view.play;
    const [bytes, decoded] = await Promise.all([this.part(seq), cues ? this.part(cues).then(decodeCues) : null]);
    return { seq: bytes, cues: decoded };
  }

  pathCollection(hash: Hash): Promise<PathData[]> {
    let p = this.pathSets.get(hash);
    if (!p) {
      p = this.part(hash).then(decodePathCollection);
      this.pathSets.set(hash, p);
    }
    return p;
  }

  /**
   * Load an object and everything it references, recursively.
   * Calls onResource for each non-object dependency (font, image, path collection)
   * so a renderer can prepare it.
   */
  async ensure(hash: Hash, onResource?: (entry: PartEntry, bytes: Uint8Array) => Promise<void> | void, seen = new Set<Hash>()): Promise<ObjectPart> {
    const o = await this.object(hash);
    this.loadedObjects.set(hash, o);
    if (seen.has(hash)) return o;
    seen.add(hash);
    await Promise.all(
      objectDeps(o).map(async (dep) => {
        const e = this.entry(dep);
        if (e.t === "obj") {
          await this.ensure(dep, onResource, seen);
        } else if (onResource && !seen.has(dep)) {
          seen.add(dep);
          await onResource(e, await this.part(dep));
        }
      }),
    );
    return o;
  }
}

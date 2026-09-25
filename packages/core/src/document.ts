import { decode, type PartSource } from "./container.js";
import { decodeObject, decodePathCollection, objectDeps } from "./object.js";
import type { Manifest, PartEntry, ObjectPart, PathData, Hash, View } from "./types.js";

/** A loaded document: manifest plus a cache of decoded parts. */
export class BdfDocument {
  private entries = new Map<Hash, PartEntry>();
  private parts = new Map<Hash, Promise<Uint8Array>>();
  private objects = new Map<Hash, Promise<ObjectPart>>();
  private pathSets = new Map<Hash, Promise<PathData[]>>();

  private constructor(readonly source: PartSource, readonly manifest: Manifest) {
    for (const e of manifest.parts) this.entries.set(e.h, e);
  }

  static async open(source: PartSource): Promise<BdfDocument> {
    return new BdfDocument(source, await source.manifest());
  }

  entry(hash: Hash): PartEntry {
    const e = this.entries.get(hash);
    if (!e) throw new Error(`bdf: unknown part ${hash}`);
    return e;
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
      const e = this.entry(hash);
      p = this.source.stored(e).then((b) => decode(b, e.enc));
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
